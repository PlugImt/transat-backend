package photoreport

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/lib/pq"
)

var (
	ErrReportNotFound = errors.New("report not found")
	ErrAlbumNotFound  = errors.New("album not found")
	ErrClubNotFound   = errors.New("club not found")
)

// Report is the administrator's view of a stored report. Only administrators ever receive it.
type Report struct {
	ID               int64      `json:"id"`
	CreatedAt        time.Time  `json:"created_at"`
	Reason           string     `json:"reason"`
	Comment          string     `json:"comment"`
	PhotoID          string     `json:"photo_id"`
	PhotoName        string     `json:"photo_name"`
	PhotoURL         string     `json:"photo_url"`
	PhotoAvailable   bool       `json:"photo_available"`
	AlbumID          *int64     `json:"album_id"`
	AlbumTitle       string     `json:"album_title"`
	AlbumURL         string     `json:"album_url"`
	ClubID           *int64     `json:"club_id"`
	ClubName         string     `json:"club_name"`
	Reporter         Reporter   `json:"reporter"`
	DeliveryStatus   string     `json:"delivery_status"` // sent | pending | failed | no_recipients
	Recipients       int        `json:"recipients"`
	RecipientsSent   int        `json:"recipients_sent"`
	RecipientsFailed int        `json:"recipients_failed"`
	Deliveries       []Delivery `json:"deliveries,omitempty"`
	// pending until someone decides; approved means the photo was removed, rejected that it stays.
	Status       string     `json:"status"`
	DecidedAt    *time.Time `json:"decided_at"`
	DecidedBy    *string    `json:"decided_by"`
	DecisionNote string     `json:"decision_note"`
}

type Reporter struct {
	ID             int     `json:"id"`
	PublicID       *string `json:"public_id"`
	Name           string  `json:"name"`
	Email          *string `json:"email"` // null once the account is deleted
	GraduationYear *int    `json:"graduation_year"`
	Formation      string  `json:"formation"`
}

type Delivery struct {
	Email     string     `json:"email"`
	Role      string     `json:"role"`
	Status    string     `json:"status"`
	Attempts  int        `json:"attempts"`
	SentAt    *time.Time `json:"sent_at"`
	LastError *string    `json:"last_error"`
}

type ListFilter struct {
	Status   string // empty = all
	Reason   string // empty = all
	Delivery string // empty = all
	Limit    int
	Offset   int
}

const reportSelect = `
	WITH r AS (
		SELECT p.*,
		       COALESCE(d.total, 0) AS total, COALESCE(d.sent, 0) AS sent, COALESCE(d.failed, 0) AS failed,
		       CASE WHEN COALESCE(d.total, 0) = 0 THEN 'no_recipients'
		            WHEN COALESCE(d.failed, 0) > 0 THEN 'failed'
		            WHEN d.sent = d.total THEN 'sent'
		            ELSE 'pending' END AS delivery_status,
		       EXISTS (SELECT 1 FROM photos ph WHERE ph.drive_file_id = p.photo_id AND ph.deleted_at IS NULL) AS photo_available
		FROM photo_reports p
		LEFT JOIN LATERAL (
			SELECT count(*) AS total,
			       count(*) FILTER (WHERE status = 'sent') AS sent,
			       count(*) FILTER (WHERE status = 'failed' AND attempts >= $1) AS failed
			FROM photo_report_deliveries WHERE report_id = p.id_photo_reports) d ON true
	)
	SELECT id_photo_reports, created_at, reason, comment, photo_id, photo_name, photo_url, photo_available,
	       album_id, album_title, album_url, club_id, club_name,
	       reporter_id, reporter_public_id::text, reporter_name, reporter_email, reporter_graduation_year, reporter_formation,
	       delivery_status, total, sent, failed, status, decided_at, decided_by, decision_note`

func scanReport(s interface{ Scan(...any) error }, extra ...any) (Report, error) {
	var r Report
	var album, club, year sql.NullInt64
	var public, email, decidedBy sql.NullString
	var decidedAt sql.NullTime
	dest := append([]any{&r.ID, &r.CreatedAt, &r.Reason, &r.Comment, &r.PhotoID, &r.PhotoName, &r.PhotoURL, &r.PhotoAvailable,
		&album, &r.AlbumTitle, &r.AlbumURL, &club, &r.ClubName,
		&r.Reporter.ID, &public, &r.Reporter.Name, &email, &year, &r.Reporter.Formation,
		&r.DeliveryStatus, &r.Recipients, &r.RecipientsSent, &r.RecipientsFailed,
		&r.Status, &decidedAt, &decidedBy, &r.DecisionNote}, extra...)
	if err := s.Scan(dest...); err != nil {
		return r, err
	}
	if decidedAt.Valid {
		r.DecidedAt = &decidedAt.Time
	}
	if decidedBy.Valid {
		r.DecidedBy = &decidedBy.String
	}
	if album.Valid {
		r.AlbumID = &album.Int64
	}
	if club.Valid {
		r.ClubID = &club.Int64
	}
	if public.Valid {
		r.Reporter.PublicID = &public.String
	}
	if email.Valid {
		r.Reporter.Email = &email.String
	}
	if year.Valid {
		y := int(year.Int64)
		r.Reporter.GraduationYear = &y
	}
	return r, nil
}

// List returns the newest reports first, with the total matching the filter.
func (s *Service) List(ctx context.Context, f ListFilter) ([]Report, int, error) {
	if f.Limit <= 0 || f.Limit > 100 {
		f.Limit = 25
	}
	rows, err := s.db.QueryContext(ctx, reportSelect+`, count(*) OVER()
		FROM r
		WHERE ($2 = '' OR reason = $2) AND ($3 = '' OR delivery_status = $3) AND ($6 = '' OR status = $6)
		ORDER BY created_at DESC, id_photo_reports DESC
		LIMIT $4 OFFSET $5`, MaxAttempts, f.Reason, f.Delivery, f.Limit, f.Offset, f.Status)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	out, total := []Report{}, 0
	for rows.Next() {
		var n int
		r, err := scanReport(rows, &n)
		if err != nil {
			return nil, 0, err
		}
		total = n
		out = append(out, r)
	}
	return out, total, rows.Err()
}

// Get returns one report with the state of each of its emails.
func (s *Service) Get(ctx context.Context, id int64) (Report, error) {
	r, err := scanReport(s.db.QueryRowContext(ctx, reportSelect+` FROM r WHERE id_photo_reports = $2`, MaxAttempts, id))
	if errors.Is(err, sql.ErrNoRows) {
		return r, ErrReportNotFound
	}
	if err != nil {
		return r, err
	}
	rows, err := s.db.QueryContext(ctx, `
		SELECT recipient_email, recipient_role, status, attempts, sent_at, last_error
		FROM photo_report_deliveries WHERE report_id = $1 ORDER BY recipient_role, recipient_email`, id)
	if err != nil {
		return r, err
	}
	defer rows.Close()
	r.Deliveries = []Delivery{}
	for rows.Next() {
		var d Delivery
		var sent sql.NullTime
		var lastErr sql.NullString
		if err := rows.Scan(&d.Email, &d.Role, &d.Status, &d.Attempts, &sent, &lastErr); err != nil {
			return r, err
		}
		if sent.Valid {
			d.SentAt = &sent.Time
		}
		if lastErr.Valid {
			d.LastError = &lastErr.String
		}
		r.Deliveries = append(r.Deliveries, d)
	}
	return r, rows.Err()
}

// SetAlbumClub links an album to the club that owns its photos, or unlinks it with a nil clubID.
func (s *Service) SetAlbumClub(ctx context.Context, albumID int64, clubID *int64) error {
	res, err := s.db.ExecContext(ctx, `UPDATE photo_albums SET club_id = $2 WHERE id_photo_albums = $1`, albumID, clubID)
	var pgErr *pq.Error
	if errors.As(err, &pgErr) && pgErr.Code == "23503" {
		return ErrClubNotFound
	}
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return ErrAlbumNotFound
	}
	return nil
}

// Counts gives the number of reports per status, for the tabs of the admin panel.
func (s *Service) Counts(ctx context.Context) (map[string]int, error) {
	out := map[string]int{"pending": 0, "approved": 0, "rejected": 0}
	rows, err := s.db.QueryContext(ctx, `SELECT status, count(*) FROM photo_reports GROUP BY status`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var st string
		var n int
		if err := rows.Scan(&st, &n); err != nil {
			return nil, err
		}
		out[st] = n
	}
	return out, rows.Err()
}

// Decision is the reviewer's answer to a report.
type Decision string

const (
	DecisionApprove Decision = "approve"
	DecisionReject  Decision = "reject"
)

var (
	ErrAlreadyDecided     = errors.New("report already decided")
	ErrRemovalUnavailable = errors.New("photo removal is not configured")
	ErrInvalidDecision    = errors.New("invalid decision")
)

// PhotoRemover deletes a photo from Google Drive (to the drive's trash, which Drive empties after 30 days). It needs
// a service account with write access; the photo sync's read-only account cannot do it.
type PhotoRemover interface {
	Trash(ctx context.Context, driveFileID string) error
}

// SetRemover enables approvals. Without it, reports can still be rejected.
func (s *Service) SetRemover(r PhotoRemover) { s.remover = r }

// CanRemove tells whether approving a report can delete its photo.
func (s *Service) CanRemove() bool { return s.remover != nil }

// Decide closes a pending report. Approving removes the photo from the Drive and hides it from the gallery at once,
// and closes the other pending reports of the same photo with it. Rejecting leaves the photo untouched.
func (s *Service) Decide(ctx context.Context, id int64, decision Decision, note, by string) (Report, error) {
	note, err := CleanComment(note)
	if err != nil {
		return Report{}, err
	}
	if decision != DecisionApprove && decision != DecisionReject {
		return Report{}, ErrInvalidDecision
	}
	if decision == DecisionApprove && s.remover == nil {
		return Report{}, ErrRemovalUnavailable
	}

	var status, photoID string
	err = s.db.QueryRowContext(ctx, `SELECT status, photo_id FROM photo_reports WHERE id_photo_reports = $1`, id).Scan(&status, &photoID)
	if errors.Is(err, sql.ErrNoRows) {
		return Report{}, ErrReportNotFound
	}
	if err != nil {
		return Report{}, err
	}
	if status != "pending" {
		return Report{}, ErrAlreadyDecided
	}

	if decision == DecisionApprove {
		// Drive first: if it refuses, nothing is recorded and the reviewer can try again.
		if err := s.remover.Trash(ctx, photoID); err != nil {
			return Report{}, fmt.Errorf("remove photo from drive: %w", err)
		}
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return Report{}, err
	}
	defer func() { _ = tx.Rollback() }()
	newStatus := "rejected"
	if decision == DecisionApprove {
		newStatus = "approved"
	}
	res, err := tx.ExecContext(ctx, `UPDATE photo_reports SET status = $2, decided_at = now(), decided_by = $3, decision_note = $4
		WHERE id_photo_reports = $1 AND status = 'pending'`, id, newStatus, by, note)
	if err != nil {
		return Report{}, err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return Report{}, ErrAlreadyDecided
	}
	if decision == DecisionApprove {
		if _, err := tx.ExecContext(ctx, `UPDATE photo_reports SET status = 'approved', decided_at = now(), decided_by = $2,
			decision_note = $3 WHERE photo_id = $1 AND status = 'pending'`,
			photoID, by, fmt.Sprintf("Photo supprimée via le signalement #%d", id)); err != nil {
			return Report{}, err
		}
		if _, err := tx.ExecContext(ctx, `UPDATE photos SET deleted_at = now() WHERE drive_file_id = $1 AND deleted_at IS NULL`, photoID); err != nil {
			return Report{}, err
		}
		// The gallery's response caches are keyed by this version.
		if _, err := tx.ExecContext(ctx, `UPDATE photo_sync_state SET data_version = data_version + 1 WHERE id = 1`); err != nil {
			return Report{}, err
		}
	}
	if err := tx.Commit(); err != nil {
		return Report{}, err
	}
	return s.Get(ctx, id)
}

// ApproveAllResult is how far one call of ApproveAll went.
type ApproveAllResult struct {
	// Photos removed from the Drive and the app during this call.
	Photos int `json:"photos"`
	// Reports closed during this call: a photo's other pending reports close with it.
	Approved int `json:"approved"`
	// Reports Drive refused to act on; they stay pending and are skipped by the next calls given them back.
	Failed []int64 `json:"failed"`
	// Pending reports still matching the filter, the failed and skipped ones excepted: call again while above zero.
	Remaining int `json:"remaining"`
}

// ApproveAll approves the pending reports (of one reason when reason is not empty), oldest first, until budget
// runs out: every approval deletes a photo on Drive, so a large backlog takes several calls. Reports whose ids are
// in skip (the failures of earlier calls) are left alone, so a photo Drive keeps refusing cannot stall the rest.
func (s *Service) ApproveAll(ctx context.Context, reason, by string, skip []int64, budget time.Duration) (ApproveAllResult, error) {
	out := ApproveAllResult{Failed: []int64{}}
	if reason != "" {
		if _, err := ParseReason(reason); err != nil {
			return out, err
		}
	}
	if s.remover == nil {
		return out, ErrRemovalUnavailable
	}
	excluded := append([]int64{}, skip...)
	pending := func() ([]int64, error) {
		rows, err := s.db.QueryContext(ctx, `SELECT id_photo_reports FROM photo_reports
			WHERE status = 'pending' AND ($1 = '' OR reason = $1) AND NOT (id_photo_reports = ANY($2))
			ORDER BY created_at, id_photo_reports`, reason, pq.Array(excluded))
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		var ids []int64
		for rows.Next() {
			var id int64
			if err := rows.Scan(&id); err != nil {
				return nil, err
			}
			ids = append(ids, id)
		}
		return ids, rows.Err()
	}

	ids, err := pending()
	if err != nil {
		return out, err
	}
	before := len(ids)
	deadline := time.Now().Add(budget)
	note := "Approuvé avec « Tout approuver »"
	for _, id := range ids {
		if time.Now().After(deadline) || ctx.Err() != nil {
			break
		}
		_, err := s.Decide(ctx, id, DecisionApprove, note, by)
		switch {
		case err == nil:
			out.Photos++
		case errors.Is(err, ErrAlreadyDecided):
			// Closed with an earlier report of the same photo.
		case errors.Is(err, ErrRemovalUnavailable):
			return out, err
		default:
			out.Failed = append(out.Failed, id)
			excluded = append(excluded, id)
		}
	}

	left, err := pending()
	if err != nil {
		return out, err
	}
	out.Remaining = len(left)
	out.Approved = before - len(out.Failed) - out.Remaining
	return out, nil
}
