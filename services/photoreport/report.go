// Package photoreport turns a user's report about a photo into a stored review request and emails it to the
// people who can act on it: the owner of the club the photo belongs to, and the administrators.
package photoreport

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/plugimt/transat-backend/models"
)

// Reason is why a photo is reported. The wire values are shared with the app.
type Reason string

const (
	ReasonAppearsInPhoto Reason = "appears_in_photo"
	ReasonInappropriate  Reason = "inappropriate"
	ReasonPoorQuality    Reason = "poor_quality"
	ReasonOther          Reason = "other"
)

var reasons = map[Reason]bool{ReasonAppearsInPhoto: true, ReasonInappropriate: true, ReasonPoorQuality: true, ReasonOther: true}

const (
	// MaxCommentRunes is enforced here, whatever the app does.
	MaxCommentRunes = 1000

	// A user may send a handful of reports in a row (several photos of one album), not a flood.
	burstLimit  = 5
	burstWindow = 10 * time.Minute
	dayLimit    = 20
	dayWindow   = 24 * time.Hour

	// The same report (reporter, photo, reason, comment) inside this window is the same report.
	duplicateWindow = 24 * time.Hour
)

var (
	ErrInvalidReason   = errors.New("invalid report reason")
	ErrInvalidPhotoID  = errors.New("invalid photo id")
	ErrCommentTooLong  = errors.New("comment too long")
	ErrPhotoNotFound   = errors.New("photo not found")
	ErrUnknownReporter = errors.New("unknown reporter")
)

// RateLimitedError carries how long the reporter has to wait.
type RateLimitedError struct{ RetryAfter time.Duration }

func (e *RateLimitedError) Error() string { return "too many reports" }

// Mailer is the part of services.EmailService the reports need: one message to one recipient.
type Mailer interface {
	SendEmail(mail models.Email, data interface{}) error
}

// Service stores reports and delivers them. A nil mailer keeps every report pending until one is configured.
type Service struct {
	db       *sql.DB
	mailer   Mailer
	adminURL string // optional base URL of the admin panel, used for the link in emails
	template string
	lease    time.Duration
}

func NewService(db *sql.DB, mailer Mailer, adminURL string) *Service {
	return &Service{
		db:       db,
		mailer:   mailer,
		adminURL: strings.TrimRight(adminURL, "/"),
		template: "email_templates/email_template_photo_report.html",
		lease:    10 * time.Minute,
	}
}

// Submission is what the client may say. Everything else is derived on the server.
type Submission struct {
	PhotoID       string
	Reason        Reason
	Comment       string
	ReporterEmail string // from the verified token, never from the request body
}

type Receipt struct {
	ID        int64
	Duplicate bool
}

// ParseReason validates a wire value.
func ParseReason(v string) (Reason, error) {
	r := Reason(strings.TrimSpace(v))
	if !reasons[r] {
		return "", ErrInvalidReason
	}
	return r, nil
}

// CleanComment drops control characters (newlines are kept, as line breaks), trims, and enforces the length cap.
func CleanComment(raw string) (string, error) {
	if len(raw) > 8*MaxCommentRunes { // cheap guard before scanning a huge string
		return "", ErrCommentTooLong
	}
	raw = strings.ReplaceAll(strings.ReplaceAll(raw, "\r\n", "\n"), "\r", "\n")
	var b strings.Builder
	for _, r := range raw {
		if r == '\n' || r == '\t' || !unicode.IsControl(r) {
			b.WriteRune(r)
		}
	}
	out := strings.TrimSpace(b.String())
	if utf8.RuneCountInString(out) > MaxCommentRunes {
		return "", ErrCommentTooLong
	}
	return out, nil
}

type photoContext struct {
	photoID, photoName, photoURL string
	albumID                      sql.NullInt64
	albumTitle, albumURL         string
	clubID                       sql.NullInt64
	clubName                     string
}

type reporter struct {
	id         int
	publicID   sql.NullString
	name       string
	graduation sql.NullInt64
	formation  string
}

// Submit records the report and who should receive it, in one transaction. Nothing is emailed here: the
// report is safe in the database before any delivery is attempted.
func (s *Service) Submit(ctx context.Context, in Submission) (Receipt, error) {
	if !reasons[in.Reason] {
		return Receipt{}, ErrInvalidReason
	}
	if l := len(in.PhotoID); l < 10 || l > 100 {
		return Receipt{}, ErrInvalidPhotoID
	}
	comment, err := CleanComment(in.Comment)
	if err != nil {
		return Receipt{}, err
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return Receipt{}, err
	}
	defer func() { _ = tx.Rollback() }()

	// One submission per reporter at a time makes the duplicate and rate checks race-free.
	if _, err := tx.ExecContext(ctx, `SELECT pg_advisory_xact_lock(hashtextextended($1, 36))`, in.ReporterEmail); err != nil {
		return Receipt{}, err
	}

	rep, err := loadReporter(ctx, tx, in.ReporterEmail)
	if err != nil {
		return Receipt{}, err
	}
	pc, err := loadPhoto(ctx, tx, in.PhotoID)
	if err != nil {
		return Receipt{}, err
	}

	key := dedupeKey(in.ReporterEmail, pc.photoID, in.Reason, comment)
	var existing int64
	err = tx.QueryRowContext(ctx, `
		SELECT id_photo_reports FROM photo_reports
		WHERE dedupe_key = $1 AND created_at > now() - make_interval(secs => $2)
		ORDER BY id_photo_reports DESC LIMIT 1`, key, duplicateWindow.Seconds()).Scan(&existing)
	if err == nil {
		return Receipt{ID: existing, Duplicate: true}, nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return Receipt{}, err
	}

	if wait, limited, err := s.rateLimited(ctx, tx, in.ReporterEmail); err != nil {
		return Receipt{}, err
	} else if limited {
		return Receipt{}, &RateLimitedError{RetryAfter: wait}
	}

	recipients, err := resolveRecipients(ctx, tx, pc.clubID, pc.clubName)
	if err != nil {
		return Receipt{}, err
	}

	var id int64
	err = tx.QueryRowContext(ctx, `
		INSERT INTO photo_reports (photo_id, photo_name, photo_url, album_id, album_title, album_url, club_id, club_name,
			reason, comment, reporter_email, reporter_id, reporter_public_id, reporter_name, reporter_graduation_year,
			reporter_formation, dedupe_key)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17)
		RETURNING id_photo_reports`,
		pc.photoID, pc.photoName, pc.photoURL, pc.albumID, pc.albumTitle, pc.albumURL, pc.clubID, pc.clubName,
		string(in.Reason), comment, in.ReporterEmail, rep.id, rep.publicID, rep.name, rep.graduation, rep.formation, key).
		Scan(&id)
	if err != nil {
		return Receipt{}, err
	}
	for _, r := range recipients {
		if _, err := tx.ExecContext(ctx, `
			INSERT INTO photo_report_deliveries (report_id, recipient_email, recipient_role, language)
			VALUES ($1, $2, $3, $4)`, id, r.email, r.role, r.language); err != nil {
			return Receipt{}, err
		}
	}
	if err := tx.Commit(); err != nil {
		return Receipt{}, err
	}
	return Receipt{ID: id}, nil
}

func loadReporter(ctx context.Context, tx *sql.Tx, email string) (reporter, error) {
	var r reporter
	var first, last string
	err := tx.QueryRowContext(ctx, `
		SELECT id_newf, public_id::text, first_name, last_name, graduation_year, COALESCE(formation_name, '')
		FROM newf WHERE email = $1`, email).Scan(&r.id, &r.publicID, &first, &last, &r.graduation, &r.formation)
	if errors.Is(err, sql.ErrNoRows) {
		return r, ErrUnknownReporter
	}
	r.name = strings.TrimSpace(first + " " + last)
	return r, err
}

// loadPhoto reads the photo and, through its album, the club it belongs to. The club comes from the album
// (or its parent event) and never from the request.
func loadPhoto(ctx context.Context, tx *sql.Tx, id string) (photoContext, error) {
	pc := photoContext{}
	err := tx.QueryRowContext(ctx, `
		SELECT p.drive_file_id, p.name, p.web_url,
		       a.id_photo_albums, COALESCE(a.title, ''), COALESCE(a.drive_url, ''),
		       c.id_clubs, COALESCE(c.name, '')
		FROM photos p
		LEFT JOIN photo_albums a ON a.id_photo_albums = p.album_id AND a.deleted_at IS NULL
		LEFT JOIN photo_albums parent ON parent.id_photo_albums = a.parent_album_id
		LEFT JOIN clubs c ON c.id_clubs = COALESCE(a.club_id, parent.club_id)
		WHERE p.drive_file_id = $1 AND p.deleted_at IS NULL AND p.media_type <> 'other'`, id).
		Scan(&pc.photoID, &pc.photoName, &pc.photoURL, &pc.albumID, &pc.albumTitle, &pc.albumURL, &pc.clubID, &pc.clubName)
	if errors.Is(err, sql.ErrNoRows) {
		return pc, ErrPhotoNotFound
	}
	return pc, err
}

func (s *Service) rateLimited(ctx context.Context, tx *sql.Tx, email string) (time.Duration, bool, error) {
	var burst, day int
	var oldestBurst, oldestDay sql.NullTime
	err := tx.QueryRowContext(ctx, `
		SELECT count(*) FILTER (WHERE created_at > now() - make_interval(secs => $2)),
		       min(created_at) FILTER (WHERE created_at > now() - make_interval(secs => $2)),
		       count(*),
		       min(created_at)
		FROM photo_reports
		WHERE reporter_email = $1 AND created_at > now() - make_interval(secs => $3)`,
		email, burstWindow.Seconds(), dayWindow.Seconds()).Scan(&burst, &oldestBurst, &day, &oldestDay)
	if err != nil {
		return 0, false, err
	}
	var wait time.Duration
	if burst >= burstLimit && oldestBurst.Valid {
		wait = time.Until(oldestBurst.Time.Add(burstWindow))
	}
	if day >= dayLimit && oldestDay.Valid {
		wait = max(wait, time.Until(oldestDay.Time.Add(dayWindow)))
	}
	if wait <= 0 && (burst >= burstLimit || day >= dayLimit) {
		wait = time.Minute
	}
	return max(wait, 0), wait > 0, nil
}

type recipient struct{ email, role, language string }

// clubRoleName is the role that marks a club's owners: the same rule the club handler uses.
func clubRoleName(clubName string) string {
	return strings.ToLower(strings.ReplaceAll(clubName, " ", "")) + "_respo"
}

// resolveRecipients returns the club's owners (when the photo has a club) and every administrator, once each.
// A photo without club, or a club without owner, simply reaches the administrators only.
func resolveRecipients(ctx context.Context, tx *sql.Tx, clubID sql.NullInt64, clubName string) ([]recipient, error) {
	ownerRole := ""
	if clubID.Valid {
		ownerRole = clubRoleName(clubName)
	}
	rows, err := tx.QueryContext(ctx, `
		SELECT DISTINCT ON (n.email) n.email, x.role, COALESCE(l.code, 'fr')
		FROM (
			SELECT nr.email, 'club_owner' AS role, 0 AS prio
			FROM newf_roles nr JOIN roles r ON r.id_roles = nr.id_roles
			WHERE $1 <> '' AND r.name = $1
			UNION ALL
			SELECT nr.email, 'admin', 1
			FROM newf_roles nr JOIN roles r ON r.id_roles = nr.id_roles
			WHERE r.name = 'ADMIN'
		) x
		JOIN newf n ON n.email = x.email
		LEFT JOIN languages l ON l.id_languages = n.language
		ORDER BY n.email, x.prio`, ownerRole)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []recipient
	for rows.Next() {
		var r recipient
		if err := rows.Scan(&r.email, &r.role, &r.language); err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	return out, rows.Err()
}

func dedupeKey(email, photoID string, reason Reason, comment string) string {
	norm := strings.Join(strings.Fields(strings.ToLower(comment)), " ")
	sum := sha256.Sum256([]byte(fmt.Sprintf("%s\x00%s\x00%s\x00%s", strings.ToLower(email), photoID, reason, norm)))
	return hex.EncodeToString(sum[:])
}
