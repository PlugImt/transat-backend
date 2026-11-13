package photobooth

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"log"
	"regexp"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/plugimt/transat-backend/services/photos"
)

const (
	KindPhoto   = "photo"
	KindCollage = "collage"

	MaxPhotos    = 20
	MaxFileBytes = 80 << 20

	// An orphan upload (the booth died between opening the session and confirming it) is looked for on Drive
	// after orphanGrace, and forgotten once Drive's resumable sessions have expired.
	orphanGrace  = 15 * time.Minute
	orphanExpiry = 8 * 24 * time.Hour
)

var (
	sessionIDRe  = regexp.MustCompile(`^[A-Za-z0-9_-]{8,64}$`)
	md5Re        = regexp.MustCompile(`^[0-9a-f]{32}$`)
	layoutRe     = regexp.MustCompile(`^[A-Z0-9_]{0,40}$`)
	shareTokenRe = regexp.MustCompile(`^[A-Za-z0-9_-]{32,64}$`)
)

// Error carries the HTTP status the handler answers with; its message is safe to show to the booth.
type Error struct {
	Status  int
	Message string
}

func (e *Error) Error() string { return e.Message }

func invalid(format string, args ...any) error {
	return &Error{Status: 400, Message: fmt.Sprintf(format, args...)}
}
func conflict(format string, args ...any) error {
	return &Error{Status: 409, Message: fmt.Sprintf(format, args...)}
}

var errNotFound = &Error{Status: 404, Message: "not found"}

// Service keeps the photobooth's shoots, their Drive files and their share links.
type Service struct {
	db           *sql.DB
	drive        Drive
	rootID       string
	shareBaseURL string
	now          func() time.Time

	folderMu sync.Mutex
	shootMu  sync.Mutex
}

func NewService(db *sql.DB, drive Drive, rootID, shareBaseURL string) *Service {
	return &Service{db: db, drive: drive, rootID: rootID, shareBaseURL: strings.TrimRight(shareBaseURL, "/"), now: time.Now}
}

// Event identifies the event a shoot belongs to, as the booth knows it.
type Event struct {
	Name       string
	Date       string // yyyy-mm-dd
	LayoutType string
}

// StartRequest describes the file the booth is about to send.
type StartRequest struct {
	Event
	Kind     string
	Number   int
	Size     int64
	MD5      string
	MimeType string
}

// StartResult tells the booth where to send the bytes.
type StartResult struct {
	FileID    int64  `json:"file_id"`
	UploadURL string `json:"upload_url"`
	FileName  string `json:"file_name"`
}

type shoot struct {
	ID        int64
	EventName string
	EventDate time.Time
	FolderID  string
	Status    string
	Token     sql.NullString
	CreatedAt time.Time
}

func validateEvent(e Event) (string, time.Time, error) {
	name := CleanEventName(e.Name)
	if name == "" {
		return "", time.Time{}, invalid("event name is required")
	}
	date, err := time.Parse("2006-01-02", e.Date)
	if err != nil || date.Year() < 2000 || date.Year() > 2100 {
		return "", time.Time{}, invalid("event date must be yyyy-mm-dd")
	}
	if !layoutRe.MatchString(e.LayoutType) {
		return "", time.Time{}, invalid("invalid layout type")
	}
	return name, date, nil
}

func validateSlot(kind string, number int) error {
	switch {
	case kind == KindCollage && number == 0:
		return nil
	case kind == KindPhoto && number >= 1 && number <= MaxPhotos:
		return nil
	}
	return invalid("invalid kind/number")
}

// StartUpload records the file and opens a Drive upload session for it in the event's Photobooth folder. A file
// already in that slot (a retake) is discarded first.
func (s *Service) StartUpload(ctx context.Context, device, sessionID string, req StartRequest) (StartResult, error) {
	if !sessionIDRe.MatchString(sessionID) {
		return StartResult{}, invalid("invalid session id")
	}
	if err := validateSlot(req.Kind, req.Number); err != nil {
		return StartResult{}, err
	}
	req.MD5 = strings.ToLower(req.MD5)
	if !md5Re.MatchString(req.MD5) {
		return StartResult{}, invalid("md5 must be 32 hex characters")
	}
	if req.Size <= 0 || req.Size > MaxFileBytes {
		return StartResult{}, invalid("size must be between 1 byte and %d MB", MaxFileBytes>>20)
	}
	if req.MimeType == "" {
		req.MimeType = "image/jpeg"
	}
	if req.MimeType != "image/jpeg" {
		return StartResult{}, invalid("only image/jpeg is accepted")
	}

	sh, err := s.ensureShoot(ctx, device, sessionID, req.Event)
	if err != nil {
		return StartResult{}, err
	}
	var fileID int64
	var name string
	var replaced []string
	err = s.inTx(ctx, func(tx *sql.Tx) error {
		var status string
		if err := tx.QueryRowContext(ctx, `SELECT status FROM photobooth_shoots WHERE id = $1 FOR UPDATE`, sh.ID).Scan(&status); err != nil {
			return err
		}
		if status != "open" {
			return conflict("shoot is %s", status)
		}
		rows, err := tx.QueryContext(ctx, `UPDATE photobooth_files SET status = 'trash'
			WHERE shoot_id = $1 AND kind = $2 AND number = $3 AND status IN ('pending', 'uploaded')
			RETURNING COALESCE(drive_file_id, '')`, sh.ID, req.Kind, req.Number)
		if err != nil {
			return err
		}
		replaced, err = scanStrings(rows)
		if err != nil {
			return err
		}
		var previous int
		if err := tx.QueryRowContext(ctx, `SELECT count(*) FROM photobooth_files WHERE shoot_id = $1 AND kind = $2 AND number = $3`,
			sh.ID, req.Kind, req.Number).Scan(&previous); err != nil {
			return err
		}
		name = FileName(sh.CreatedAt, sessionID, req.Kind, req.Number, previous+1)
		return tx.QueryRowContext(ctx, `INSERT INTO photobooth_files (shoot_id, kind, number, file_name, mime_type, size_bytes, md5)
			VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`,
			sh.ID, req.Kind, req.Number, name, req.MimeType, req.Size, req.MD5).Scan(&fileID)
	})
	if err != nil {
		return StartResult{}, err
	}
	s.trashNow(replaced)

	uploadURL, err := s.drive.StartUpload(ctx, sh.FolderID, name, req.MimeType, req.Size)
	if err != nil {
		// Nothing exists on Drive yet: the row can go straight to deleted.
		_, _ = s.db.ExecContext(context.WithoutCancel(ctx), `UPDATE photobooth_files SET status = 'deleted' WHERE id = $1`, fileID)
		if isNotFound(err) {
			s.forgetFolder(ctx, sh.EventDate, sh.EventName)
		}
		return StartResult{}, fmt.Errorf("open drive upload: %w", err)
	}
	return StartResult{FileID: fileID, UploadURL: uploadURL, FileName: name}, nil
}

// ensureShoot returns the booth session's shoot, creating it (and resolving its Drive folder) on first use.
func (s *Service) ensureShoot(ctx context.Context, device, sessionID string, e Event) (shoot, error) {
	name, date, err := validateEvent(e)
	if err != nil {
		return shoot{}, err
	}
	// Serialises the first uploads of a session, which arrive close together, so the folder is resolved once.
	s.shootMu.Lock()
	defer s.shootMu.Unlock()
	sh, err := s.loadShoot(ctx, device, sessionID)
	if err == nil {
		if NameKey(sh.EventName) != NameKey(name) || !sh.EventDate.Equal(date) {
			return shoot{}, conflict("session already belongs to another event")
		}
		return sh, nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return shoot{}, err
	}
	folder, err := s.resolveFolder(ctx, date, name)
	if err != nil {
		return shoot{}, fmt.Errorf("resolve drive folder: %w", err)
	}
	_, err = s.db.ExecContext(ctx, `INSERT INTO photobooth_shoots
		(device_key_id, client_session_id, event_name, event_date, layout_type, drive_folder_id)
		VALUES ($1, $2, $3, $4, $5, $6) ON CONFLICT (device_key_id, client_session_id) DO NOTHING`,
		device, sessionID, name, date, e.LayoutType, folder)
	if err != nil {
		return shoot{}, err
	}
	return s.loadShoot(ctx, device, sessionID)
}

func (s *Service) loadShoot(ctx context.Context, device, sessionID string) (shoot, error) {
	var sh shoot
	err := s.db.QueryRowContext(ctx, `SELECT id, event_name, event_date, drive_folder_id, status, share_token, created_at
		FROM photobooth_shoots WHERE device_key_id = $1 AND client_session_id = $2`, device, sessionID).
		Scan(&sh.ID, &sh.EventName, &sh.EventDate, &sh.FolderID, &sh.Status, &sh.Token, &sh.CreatedAt)
	return sh, err
}

// CompleteUpload checks that the file the booth sent is on Drive exactly as announced, then keeps it. A file
// whose slot was discarded in the meantime (retake, cancelled session) is trashed instead.
func (s *Service) CompleteUpload(ctx context.Context, device, sessionID string, fileID int64, driveFileID string) error {
	if !sessionIDRe.MatchString(sessionID) || !driveIDRe.MatchString(driveFileID) {
		return invalid("invalid identifiers")
	}
	var (
		shootID                         int64
		folder, name, md5, status, prev string
		size                            int64
	)
	err := s.db.QueryRowContext(ctx, `SELECT f.shoot_id, sh.drive_folder_id, f.file_name, f.md5, f.size_bytes, f.status,
			COALESCE(f.drive_file_id, '')
		FROM photobooth_files f JOIN photobooth_shoots sh ON sh.id = f.shoot_id
		WHERE f.id = $1 AND sh.device_key_id = $2 AND sh.client_session_id = $3`, fileID, device, sessionID).
		Scan(&shootID, &folder, &name, &md5, &size, &status, &prev)
	if errors.Is(err, sql.ErrNoRows) {
		return errNotFound
	}
	if err != nil {
		return err
	}
	if prev != "" {
		if prev != driveFileID {
			return conflict("file already completed with another drive file")
		}
		if status == "uploaded" {
			return nil
		}
	}

	info, err := s.drive.GetFile(ctx, driveFileID)
	if errors.Is(err, ErrNotFound) {
		return invalid("drive file not found")
	}
	if err != nil {
		return err
	}
	if !slices.Contains(info.Parents, folder) || info.Name != name || info.Trashed {
		// Not the file this slot announced: never touch it.
		return invalid("drive file does not match the upload")
	}
	if info.Size != size || !strings.EqualFold(info.MD5, md5) {
		_ = s.drive.Trash(ctx, driveFileID)
		return invalid("uploaded content does not match its size or checksum")
	}

	var kept bool
	err = s.inTx(ctx, func(tx *sql.Tx) error {
		var shootStatus string
		if err := tx.QueryRowContext(ctx, `SELECT status FROM photobooth_shoots WHERE id = $1 FOR UPDATE`, shootID).Scan(&shootStatus); err != nil {
			return err
		}
		if err := tx.QueryRowContext(ctx, `SELECT status FROM photobooth_files WHERE id = $1 FOR UPDATE`, fileID).Scan(&status); err != nil {
			return err
		}
		kept = status == "pending" && shootStatus == "open"
		newStatus := "trash"
		if kept {
			newStatus = "uploaded"
		}
		_, err := tx.ExecContext(ctx, `UPDATE photobooth_files SET drive_file_id = $2, status = $3, width = NULLIF($4, 0),
			height = NULLIF($5, 0), uploaded_at = now() WHERE id = $1`, fileID, driveFileID, newStatus, info.Width, info.Height)
		return err
	})
	if err != nil {
		return err
	}
	if !kept {
		s.trashNow([]string{driveFileID})
		return &Error{Status: 410, Message: "this picture was discarded"}
	}
	return nil
}

// DiscardSlot removes the picture in one slot (the guest retook it). Idempotent.
func (s *Service) DiscardSlot(ctx context.Context, device, sessionID, kind string, number int) error {
	if !sessionIDRe.MatchString(sessionID) {
		return invalid("invalid session id")
	}
	if err := validateSlot(kind, number); err != nil {
		return err
	}
	var ids []string
	err := s.inTx(ctx, func(tx *sql.Tx) error {
		var shootID int64
		var status string
		err := tx.QueryRowContext(ctx, `SELECT id, status FROM photobooth_shoots WHERE device_key_id = $1 AND client_session_id = $2
			FOR UPDATE`, device, sessionID).Scan(&shootID, &status)
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}
		if err != nil {
			return err
		}
		if status == "finalized" {
			return conflict("shoot is already shared")
		}
		rows, err := tx.QueryContext(ctx, `UPDATE photobooth_files SET status = 'trash'
			WHERE shoot_id = $1 AND kind = $2 AND number = $3 AND status IN ('pending', 'uploaded')
			RETURNING COALESCE(drive_file_id, '')`, shootID, kind, number)
		if err != nil {
			return err
		}
		ids, err = scanStrings(rows)
		return err
	})
	if err != nil {
		return err
	}
	s.trashNow(ids)
	return nil
}

// Cancel removes every picture of an abandoned session from Drive. Idempotent; unknown sessions are fine.
func (s *Service) Cancel(ctx context.Context, device, sessionID string) error {
	if !sessionIDRe.MatchString(sessionID) {
		return invalid("invalid session id")
	}
	var ids []string
	err := s.inTx(ctx, func(tx *sql.Tx) error {
		var shootID int64
		var status string
		err := tx.QueryRowContext(ctx, `SELECT id, status FROM photobooth_shoots WHERE device_key_id = $1 AND client_session_id = $2
			FOR UPDATE`, device, sessionID).Scan(&shootID, &status)
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}
		if err != nil {
			return err
		}
		if status == "finalized" {
			return conflict("shoot is already shared")
		}
		if _, err := tx.ExecContext(ctx, `UPDATE photobooth_shoots SET status = 'cancelled',
			cancelled_at = COALESCE(cancelled_at, now()) WHERE id = $1`, shootID); err != nil {
			return err
		}
		rows, err := tx.QueryContext(ctx, `UPDATE photobooth_files SET status = 'trash'
			WHERE shoot_id = $1 AND status IN ('pending', 'uploaded') RETURNING COALESCE(drive_file_id, '')`, shootID)
		if err != nil {
			return err
		}
		ids, err = scanStrings(rows)
		return err
	})
	if err != nil {
		return err
	}
	s.trashNow(ids)
	return nil
}

// Finalize checks that the collage and photos 1..photoCount are all on Drive and returns the shoot's permanent
// share link, the same one on every call.
func (s *Service) Finalize(ctx context.Context, device, sessionID string, photoCount int) (string, error) {
	if !sessionIDRe.MatchString(sessionID) {
		return "", invalid("invalid session id")
	}
	if photoCount < 1 || photoCount > MaxPhotos {
		return "", invalid("photo_count must be between 1 and %d", MaxPhotos)
	}
	var token string
	err := s.inTx(ctx, func(tx *sql.Tx) error {
		var shootID int64
		var status string
		var existing sql.NullString
		err := tx.QueryRowContext(ctx, `SELECT id, status, share_token FROM photobooth_shoots
			WHERE device_key_id = $1 AND client_session_id = $2 FOR UPDATE`, device, sessionID).Scan(&shootID, &status, &existing)
		if errors.Is(err, sql.ErrNoRows) {
			return errNotFound
		}
		if err != nil {
			return err
		}
		if status == "finalized" && existing.Valid {
			token = existing.String
			return nil
		}
		if status != "open" {
			return conflict("shoot is %s", status)
		}
		rows, err := tx.QueryContext(ctx, `SELECT kind, number, status FROM photobooth_files
			WHERE shoot_id = $1 AND status IN ('pending', 'uploaded')`, shootID)
		if err != nil {
			return err
		}
		uploaded := map[string]bool{}
		for rows.Next() {
			var kind, st string
			var n int
			if err := rows.Scan(&kind, &n, &st); err != nil {
				_ = rows.Close()
				return err
			}
			uploaded[fmt.Sprintf("%s-%d", kind, n)] = st == "uploaded"
		}
		_ = rows.Close()
		if err := rows.Err(); err != nil {
			return err
		}
		var missing []string
		if !uploaded[KindCollage+"-0"] {
			missing = append(missing, KindCollage)
		}
		for n := 1; n <= photoCount; n++ {
			if !uploaded[fmt.Sprintf("%s-%d", KindPhoto, n)] {
				missing = append(missing, fmt.Sprintf("%s %d", KindPhoto, n))
			}
		}
		if len(missing) > 0 {
			return conflict("not uploaded yet: %s", strings.Join(missing, ", "))
		}
		if token, err = newToken(); err != nil {
			return err
		}
		_, err = tx.ExecContext(ctx, `UPDATE photobooth_shoots SET status = 'finalized', share_token = $2, finalized_at = now()
			WHERE id = $1`, shootID, token)
		return err
	})
	if err != nil {
		return "", err
	}
	return s.shareBaseURL + "/" + token, nil
}

func newToken() (string, error) {
	b := make([]byte, 24) // 192 bits: unguessable, 32 characters in the URL
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

// ---- Public share page ----------------------------------------------------------------------------------

type ShareFile struct {
	ID        int64  `json:"id"`
	Kind      string `json:"kind"`
	Number    int    `json:"number"`
	Name      string `json:"name"`
	SizeBytes int64  `json:"size_bytes"`
	Width     int    `json:"width,omitempty"`
	Height    int    `json:"height,omitempty"`
}

type Share struct {
	EventName string      `json:"event_name"`
	EventDate string      `json:"event_date"`
	TakenAt   time.Time   `json:"taken_at"`
	Files     []ShareFile `json:"files"`
}

// GetShare returns what the share page shows: the collage first, then the photos in order.
func (s *Service) GetShare(ctx context.Context, token string) (Share, error) {
	if !shareTokenRe.MatchString(token) {
		return Share{}, errNotFound
	}
	var sh Share
	var shootID int64
	var date time.Time
	err := s.db.QueryRowContext(ctx, `SELECT id, event_name, event_date, created_at FROM photobooth_shoots
		WHERE share_token = $1 AND status = 'finalized'`, token).Scan(&shootID, &sh.EventName, &date, &sh.TakenAt)
	if errors.Is(err, sql.ErrNoRows) {
		return Share{}, errNotFound
	}
	if err != nil {
		return Share{}, err
	}
	sh.EventDate = date.Format("2006-01-02")
	rows, err := s.db.QueryContext(ctx, `SELECT id, kind, number, file_name, size_bytes, COALESCE(width, 0), COALESCE(height, 0)
		FROM photobooth_files WHERE shoot_id = $1 AND status = 'uploaded'
		ORDER BY CASE kind WHEN 'collage' THEN 0 ELSE 1 END, number`, shootID)
	if err != nil {
		return Share{}, err
	}
	defer rows.Close()
	sh.Files = []ShareFile{}
	for rows.Next() {
		var f ShareFile
		if err := rows.Scan(&f.ID, &f.Kind, &f.Number, &f.Name, &f.SizeBytes, &f.Width, &f.Height); err != nil {
			return Share{}, err
		}
		sh.Files = append(sh.Files, f)
	}
	return sh, rows.Err()
}

// SharedFile is an open stream of one shared picture.
type SharedFile struct {
	Body        io.ReadCloser
	ContentType string
	Size        int64 // -1 when unknown
	Name        string
}

// OpenSharedFile streams a shared picture: the original, or a preview at most previewSize pixels wide (the
// original when Drive has no preview yet).
func (s *Service) OpenSharedFile(ctx context.Context, token string, fileID int64, previewSize int) (SharedFile, error) {
	if !shareTokenRe.MatchString(token) {
		return SharedFile{}, errNotFound
	}
	var driveID, name string
	var size int64
	err := s.db.QueryRowContext(ctx, `SELECT f.drive_file_id, f.file_name, f.size_bytes FROM photobooth_files f
		JOIN photobooth_shoots sh ON sh.id = f.shoot_id
		WHERE f.id = $1 AND sh.share_token = $2 AND sh.status = 'finalized' AND f.status = 'uploaded'`, fileID, token).
		Scan(&driveID, &name, &size)
	if errors.Is(err, sql.ErrNoRows) {
		return SharedFile{}, errNotFound
	}
	if err != nil {
		return SharedFile{}, err
	}
	return s.openDriveFile(ctx, driveID, name, size, previewSize)
}

// openDriveFile streams a picture from Drive: the original, or a preview at most previewSize pixels wide.
func (s *Service) openDriveFile(ctx context.Context, driveID, name string, size int64, previewSize int) (SharedFile, error) {
	if previewSize > 0 {
		if info, err := s.drive.GetFile(ctx, driveID); err == nil {
			if body, ct, err := s.drive.Thumbnail(ctx, info, previewSize); err == nil {
				if ct == "" {
					ct = "image/jpeg"
				}
				return SharedFile{Body: body, ContentType: ct, Size: -1, Name: name}, nil
			}
		} else if errors.Is(err, ErrNotFound) {
			return SharedFile{}, errNotFound
		}
	}
	body, err := s.drive.Download(ctx, driveID)
	if errors.Is(err, ErrNotFound) {
		return SharedFile{}, errNotFound
	}
	if err != nil {
		return SharedFile{}, err
	}
	return SharedFile{Body: body, ContentType: "image/jpeg", Size: size, Name: name}, nil
}

// ---- Drive folders --------------------------------------------------------------------------------------

// resolveFolder finds, or creates, "<school year>/<YYYY.MM.DD - Event>/photobooth_temp" under the root folder.
// Pictures always land in a pending folder, even when the event already has an approved "Photobooth" one.
func (s *Service) resolveFolder(ctx context.Context, date time.Time, eventName string) (string, error) {
	s.folderMu.Lock()
	defer s.folderMu.Unlock()
	key := NameKey(eventName)

	var cached string
	err := s.db.QueryRowContext(ctx, `SELECT photobooth_folder FROM photobooth_event_folders WHERE event_date = $1 AND event_key = $2`,
		date, key).Scan(&cached)
	if err == nil {
		// Once approved (renamed) the folder is published: later pictures need a new pending folder.
		if info, err := s.drive.GetFile(ctx, cached); err == nil && !info.Trashed && photos.IsPendingPhotoboothFolder(info.Name) {
			return cached, nil
		} else if err != nil && !errors.Is(err, ErrNotFound) {
			return "", err
		}
	} else if !errors.Is(err, sql.ErrNoRows) {
		return "", err
	}

	year := SchoolYear(date)
	yearID, err := s.childFolder(ctx, s.rootID, year, func(n string) bool { return IsSchoolYearFolder(n, year) })
	if err != nil {
		return "", err
	}
	eventID, err := s.childFolder(ctx, yearID, EventFolderName(date, eventName), func(n string) bool {
		return MatchesEventFolder(n, date, eventName)
	})
	if err != nil {
		return "", err
	}
	boothID, err := s.childFolder(ctx, eventID, PendingFolderName, photos.IsPendingPhotoboothFolder)
	if err != nil {
		return "", err
	}
	_, err = s.db.ExecContext(ctx, `INSERT INTO photobooth_event_folders (event_date, event_key, event_folder_id, photobooth_folder)
		VALUES ($1, $2, $3, $4) ON CONFLICT (event_date, event_key)
		DO UPDATE SET event_folder_id = EXCLUDED.event_folder_id, photobooth_folder = EXCLUDED.photobooth_folder, updated_at = now()`,
		date, key, eventID, boothID)
	if err != nil {
		return "", err
	}
	log.Printf("Photobooth: Drive folder for %s %q resolved to %s", date.Format("2006-01-02"), eventName, boothID)
	return boothID, nil
}

// childFolder returns the first child folder accepted by match (by name order), creating createName otherwise.
func (s *Service) childFolder(ctx context.Context, parentID, createName string, match func(string) bool) (string, error) {
	children, err := s.drive.ListChildFolders(ctx, parentID)
	if err != nil {
		return "", err
	}
	slices.SortFunc(children, func(a, b Folder) int { return strings.Compare(a.Name, b.Name) })
	for _, c := range children {
		if match(c.Name) {
			return c.ID, nil
		}
	}
	log.Printf("Photobooth: creating Drive folder %q", createName)
	return s.drive.CreateFolder(ctx, parentID, createName)
}

func (s *Service) forgetFolder(ctx context.Context, date time.Time, eventName string) {
	_, _ = s.db.ExecContext(context.WithoutCancel(ctx), `DELETE FROM photobooth_event_folders WHERE event_date = $1 AND event_key = $2`,
		date, NameKey(eventName))
}

// ---- Trash ----------------------------------------------------------------------------------------------

// trashNow removes files from Drive in the background; whatever fails is retried by RunTrash.
func (s *Service) trashNow(driveIDs []string) {
	ids := slices.DeleteFunc(slices.Clone(driveIDs), func(id string) bool { return id == "" })
	if len(ids) == 0 {
		return
	}
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cancel()
		for _, id := range ids {
			if err := s.drive.Trash(ctx, id); err != nil {
				log.Printf("Photobooth: trash %s failed, will retry: %v", id, err)
				continue
			}
			_, _ = s.db.ExecContext(ctx, `UPDATE photobooth_files SET status = 'deleted' WHERE drive_file_id = $1 AND status = 'trash'`, id)
		}
	}()
}

// RunTrash removes from Drive every discarded picture still there, including uploads the booth never
// confirmed (found by their name). Safe to run concurrently with requests.
func (s *Service) RunTrash(ctx context.Context) (int, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT f.id, COALESCE(f.drive_file_id, ''), f.file_name, sh.drive_folder_id, f.created_at
		FROM photobooth_files f JOIN photobooth_shoots sh ON sh.id = f.shoot_id
		WHERE f.status = 'trash' ORDER BY f.id LIMIT 200`)
	if err != nil {
		return 0, err
	}
	type item struct {
		id                    int64
		driveID, name, folder string
		createdAt             time.Time
	}
	var items []item
	for rows.Next() {
		var it item
		if err := rows.Scan(&it.id, &it.driveID, &it.name, &it.folder, &it.createdAt); err != nil {
			_ = rows.Close()
			return 0, err
		}
		items = append(items, it)
	}
	_ = rows.Close()
	if err := rows.Err(); err != nil {
		return 0, err
	}

	done := 0
	for _, it := range items {
		if it.driveID != "" {
			if err := s.drive.Trash(ctx, it.driveID); err != nil {
				log.Printf("Photobooth: trash %s failed: %v", it.driveID, err)
				continue
			}
		} else {
			age := s.now().Sub(it.createdAt)
			if age < orphanGrace {
				continue
			}
			if err := s.trashOrphans(ctx, it.folder, it.name); err != nil {
				log.Printf("Photobooth: orphan lookup for %s failed: %v", it.name, err)
				continue
			}
			if age < orphanExpiry {
				continue // the upload session may still deliver the file: look again next time
			}
		}
		if _, err := s.db.ExecContext(ctx, `UPDATE photobooth_files SET status = 'deleted' WHERE id = $1 AND status = 'trash'`, it.id); err != nil {
			return done, err
		}
		done++
	}
	return done, nil
}

// trashOrphans trashes files with this name in the folder that no live picture refers to.
func (s *Service) trashOrphans(ctx context.Context, folder, name string) error {
	files, err := s.drive.ListChildFilesNamed(ctx, folder, name)
	if err != nil {
		return err
	}
	for _, f := range files {
		var live bool
		if err := s.db.QueryRowContext(ctx, `SELECT EXISTS (SELECT 1 FROM photobooth_files
			WHERE drive_file_id = $1 AND status IN ('pending', 'uploaded'))`, f.ID).Scan(&live); err != nil {
			return err
		}
		if !live {
			if err := s.drive.Trash(ctx, f.ID); err != nil {
				return err
			}
		}
	}
	return nil
}

// ---- helpers --------------------------------------------------------------------------------------------

func (s *Service) inTx(ctx context.Context, fn func(*sql.Tx) error) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	if err := fn(tx); err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}

func scanStrings(rows *sql.Rows) ([]string, error) {
	defer rows.Close()
	var out []string
	for rows.Next() {
		var v string
		if err := rows.Scan(&v); err != nil {
			return nil, err
		}
		out = append(out, v)
	}
	return out, rows.Err()
}
