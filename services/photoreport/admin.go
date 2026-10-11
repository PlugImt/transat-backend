package photoreport

import (
	"context"
	"database/sql"
	"errors"
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
	       delivery_status, total, sent, failed`

func scanReport(s interface{ Scan(...any) error }, extra ...any) (Report, error) {
	var r Report
	var album, club, year sql.NullInt64
	var public, email sql.NullString
	dest := append([]any{&r.ID, &r.CreatedAt, &r.Reason, &r.Comment, &r.PhotoID, &r.PhotoName, &r.PhotoURL, &r.PhotoAvailable,
		&album, &r.AlbumTitle, &r.AlbumURL, &club, &r.ClubName,
		&r.Reporter.ID, &public, &r.Reporter.Name, &email, &year, &r.Reporter.Formation,
		&r.DeliveryStatus, &r.Recipients, &r.RecipientsSent, &r.RecipientsFailed}, extra...)
	if err := s.Scan(dest...); err != nil {
		return r, err
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
		WHERE ($2 = '' OR reason = $2) AND ($3 = '' OR delivery_status = $3)
		ORDER BY created_at DESC, id_photo_reports DESC
		LIMIT $4 OFFSET $5`, MaxAttempts, f.Reason, f.Delivery, f.Limit, f.Offset)
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
