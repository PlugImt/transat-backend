// Package faces builds each consenting user's "my photos" album with face recognition.
//
// The photos stay on Google Drive: the worker reads Drive's resized rendition of every indexed photo, sends it to
// the machine-learning container (Immich's, on Railway's private network), and keeps only the face boxes and
// embeddings. A user who opts in sends a selfie, from which one reference embedding is kept (the selfie itself is
// never stored); their album is the photos whose faces are close enough to one of their references.
//
// Biometric data (GDPR art. 9): nothing is matched without an explicit consent row, withdrawing consent deletes
// every reference and match of the user, and so do the account deletion and the inactivity purge (cascade).
package faces

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/lib/pq"
)

// ConsentVersion identifies the consent text the app shows. Changing the text means bumping it: users who agreed
// to an older version must agree again before new references are accepted.
const ConsentVersion = "2026-11-v1"

const (
	MaxReferences     = 3
	referenceMinScore = 0.75
	// A reference face must cover at least this share of the selfie's shorter side.
	referenceMinSize = 0.12
	// When a selfie shows several faces, the largest is used only if it is this much larger than the next one.
	referenceDominance = 3.0
)

var (
	ErrUnavailable    = errors.New("face recognition is not available")
	ErrNoConsent      = errors.New("face recognition consent is required")
	ErrConsentVersion = errors.New("the consent text has changed")
	ErrNoFace         = errors.New("no face found on the photo")
	ErrManyFaces      = errors.New("several faces found on the photo")
	ErrFaceTooSmall   = errors.New("the face is too small or blurry")
	ErrTooManyRefs    = errors.New("too many reference photos")
	ErrNotFound       = errors.New("not found")
)

// Config tunes recognition. MaxDistance is the cosine distance under which a face is taken for the user's.
type Config struct {
	MaxDistance     float32
	MinFaceScore    float32
	MinFacePixels   int // faces smaller than this on the scanned image are ignored (too few pixels to recognise)
	ScanSize        int // long side of the image sent to the model
	ScanBatch       int
	ScanConcurrency int
	TickBudget      time.Duration
	// ScopeFolders limits scanning to these Drive folders and everything below them (FACES_SCOPE_FOLDERS).
	// Empty: the whole gallery.
	ScopeFolders []string
	// ScanDisabled: selfies and albums still work, but no photo is scanned (see FromEnv).
	ScanDisabled bool
}

func DefaultConfig() Config {
	return Config{
		MaxDistance:     0.5,
		MinFaceScore:    0.7,
		MinFacePixels:   24,
		ScanSize:        2048,
		ScanBatch:       40,
		ScanConcurrency: 2,
		TickBudget:      4 * time.Minute,
	}
}

type Service struct {
	db     *sql.DB
	ml     Recognizer  // nil when the ML service is not configured
	images ImageSource // nil when Drive is not configured: references still work, nothing is scanned
	cfg    Config
	wake   chan struct{}
	now    func() time.Time
}

func NewService(db *sql.DB, ml Recognizer, images ImageSource, cfg Config) *Service {
	return &Service{db: db, ml: ml, images: images, cfg: cfg, wake: make(chan struct{}, 1), now: time.Now}
}

// Available tells whether the app should offer the feature.
func (s *Service) Available() bool { return s != nil && s.ml != nil }

// Status is what the app needs to draw the feature's screens.
type Status struct {
	Available             bool       `json:"available"`
	Enabled               bool       `json:"enabled"`
	ConsentVersion        *string    `json:"consent_version"`
	CurrentConsentVersion string     `json:"current_consent_version"`
	ConsentedAt           *time.Time `json:"consented_at"`
	References            int        `json:"references"`
	MaxReferences         int        `json:"max_references"`
	// Matching is true while the library is still being compared with the user's latest reference.
	Matching   bool `json:"matching"`
	PhotoCount int  `json:"photo_count"`
	// Library progress: photos that went through the detector out of those eligible.
	LibraryScanned int `json:"library_scanned"`
	LibraryTotal   int `json:"library_total"`
}

// eligiblePhotos is the gallery's own visibility rule, restricted to images with a thumbnail.
const eligiblePhotos = `p.deleted_at IS NULL AND p.duplicate_of IS NULL AND p.media_type = 'image' AND NOT p.is_cover AND p.has_thumbnail`

// inScope keeps the photos under the configured folders, given as the text[] parameter n (NULL: no limit).
func inScope(n int) string {
	return fmt.Sprintf(`($%[1]d::text[] IS NULL OR p.folder_id IN (
		WITH RECURSIVE tree AS (
			SELECT drive_folder_id FROM photo_folders WHERE drive_folder_id = ANY($%[1]d) AND deleted_at IS NULL
			UNION
			SELECT f.drive_folder_id FROM photo_folders f JOIN tree t ON f.parent_drive_id = t.drive_folder_id
			WHERE f.deleted_at IS NULL)
		SELECT drive_folder_id FROM tree))`, n)
}

// scopeParam is the parameter inScope expects.
func (s *Service) scopeParam() any {
	if len(s.cfg.ScopeFolders) == 0 {
		return nil
	}
	return pq.Array(s.cfg.ScopeFolders)
}

func (s *Service) Status(ctx context.Context, email string) (Status, error) {
	st := Status{Available: s.Available(), CurrentConsentVersion: ConsentVersion, MaxReferences: MaxReferences}
	var version sql.NullString
	var at sql.NullTime
	err := s.db.QueryRowContext(ctx, `SELECT consent_version, consented_at FROM face_consents WHERE email = $1`, email).
		Scan(&version, &at)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return st, err
	}
	if err == nil {
		st.Enabled = true
		st.ConsentVersion = &version.String
		st.ConsentedAt = &at.Time
		if err := s.db.QueryRowContext(ctx, `
			SELECT count(*), COALESCE(bool_or(r.matched_through < (SELECT COALESCE(max(id), 0) FROM photo_faces)), FALSE)
			FROM face_references r WHERE r.email = $1`, email).Scan(&st.References, &st.Matching); err != nil {
			return st, err
		}
		if err := s.db.QueryRowContext(ctx, `
			SELECT count(*) FROM photos p
			WHERE `+eligiblePhotos+` AND EXISTS (
				SELECT 1 FROM face_matches m WHERE m.email = $1 AND m.drive_file_id = p.drive_file_id AND m.rejected_at IS NULL)`,
			email).Scan(&st.PhotoCount); err != nil {
			return st, err
		}
	}
	if st.Available {
		model := s.ml.Model()
		if err := s.db.QueryRowContext(ctx, `
			SELECT count(*), count(s.drive_file_id) FROM photos p
			LEFT JOIN face_photo_scans s ON s.drive_file_id = p.drive_file_id AND s.model = $1 AND s.status <> 'failed'
			WHERE `+eligiblePhotos+` AND `+inScope(2), model, s.scopeParam()).Scan(&st.LibraryTotal, &st.LibraryScanned); err != nil {
			return st, err
		}
	}
	return st, nil
}

// Consent records the user's explicit opt-in to the given consent text.
func (s *Service) Consent(ctx context.Context, email, version string) error {
	if !s.Available() {
		return ErrUnavailable
	}
	if version != ConsentVersion {
		return ErrConsentVersion
	}
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO face_consents (email, consent_version, consented_at) VALUES ($1, $2, $3)
		ON CONFLICT (email) DO UPDATE SET consent_version = EXCLUDED.consent_version, consented_at = EXCLUDED.consented_at`,
		email, version, s.now().UTC())
	return err
}

// Withdraw deletes the consent and, through the cascade, every reference and match of the user. It is idempotent.
func (s *Service) Withdraw(ctx context.Context, email string) error {
	_, err := s.db.ExecContext(ctx, `DELETE FROM face_consents WHERE email = $1`, email)
	return err
}

// AddReference computes a reference embedding from a selfie. The image only lives in memory for this call.
func (s *Service) AddReference(ctx context.Context, email string, image []byte) error {
	if !s.Available() {
		return ErrUnavailable
	}
	var version string
	err := s.db.QueryRowContext(ctx, `SELECT consent_version FROM face_consents WHERE email = $1`, email).Scan(&version)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrNoConsent
	}
	if err != nil {
		return err
	}
	if version != ConsentVersion {
		return ErrConsentVersion
	}
	var n int
	if err := s.db.QueryRowContext(ctx, `SELECT count(*) FROM face_references WHERE email = $1`, email).Scan(&n); err != nil {
		return err
	}
	if n >= MaxReferences {
		return ErrTooManyRefs
	}

	det, err := s.ml.Detect(ctx, image)
	if errors.Is(err, ErrUnreadableImage) {
		return ErrNoFace
	}
	if err != nil {
		return err
	}
	face, err := pickReference(det)
	if err != nil {
		return err
	}

	// Consent may have been withdrawn while the model was running: the insert then fails on the foreign key.
	res, err := s.db.ExecContext(ctx, `
		INSERT INTO face_references (email, model, embedding)
		SELECT email, $2, $3 FROM face_consents WHERE email = $1`, email, s.ml.Model(), encode(face.Embedding))
	if err != nil {
		return err
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return ErrNoConsent
	}
	s.Wake()
	return nil
}

// pickReference keeps the selfie's face: a single clear face, or one clearly larger than any other.
func pickReference(d Detection) (Face, error) {
	var faces []Face
	for _, f := range d.Faces {
		if f.Score >= referenceMinScore {
			faces = append(faces, f)
		}
	}
	if len(faces) == 0 {
		if len(d.Faces) > 0 {
			return Face{}, ErrFaceTooSmall
		}
		return Face{}, ErrNoFace
	}
	area := func(f Face) float32 { return f.Width() * f.Height() }
	sort.Slice(faces, func(i, j int) bool { return area(faces[i]) > area(faces[j]) })
	if len(faces) > 1 && area(faces[0]) < referenceDominance*area(faces[1]) {
		return Face{}, ErrManyFaces
	}
	f := faces[0]
	short := float32(min(d.ImageWidth, d.ImageHeight))
	if f.Width()*float32(d.ImageWidth) < referenceMinSize*short || f.Height()*float32(d.ImageHeight) < referenceMinSize*short {
		return Face{}, ErrFaceTooSmall
	}
	return f, nil
}

// Reject removes a photo from the user's album ("this is not me"); the faces stay unmatched for them.
func (s *Service) Reject(ctx context.Context, email, driveFileID string) error {
	res, err := s.db.ExecContext(ctx, `
		UPDATE face_matches SET rejected_at = $3 WHERE email = $1 AND drive_file_id = $2 AND rejected_at IS NULL`,
		email, driveFileID, s.now().UTC())
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return ErrNotFound
	}
	return nil
}

// Wake asks the worker to run soon (after a new reference), without waiting for the next interval.
func (s *Service) Wake() {
	select {
	case s.wake <- struct{}{}:
	default:
	}
}
