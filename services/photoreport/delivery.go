package photoreport

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"
	"unicode"

	"github.com/plugimt/transat-backend/models"
)

const (
	// MaxAttempts is how many times one recipient's email is tried before it is left as failed for an admin.
	MaxAttempts = 5
	batchSize   = 25
)

// retryDelay grows with the attempts already made: 2m, 10m, 30m, 2h.
func retryDelay(attempts int) time.Duration {
	steps := []time.Duration{2 * time.Minute, 10 * time.Minute, 30 * time.Minute, 2 * time.Hour}
	if attempts < 1 {
		attempts = 1
	}
	if attempts > len(steps) {
		attempts = len(steps)
	}
	return steps[attempts-1]
}

type claim struct {
	reportID int64
	email    string
	role     string
	language string
	attempts int
}

// DeliverDue sends every email that is due: new ones, retries whose delay has passed, and ones whose sender
// crashed. Safe to run from several goroutines or instances: each row is claimed atomically. It returns the
// number of emails accepted by the email service.
func (s *Service) DeliverDue(ctx context.Context) (int, error) {
	if s.mailer == nil {
		return 0, nil
	}
	// A sender that died on its last attempt must not leave the row "sending" forever.
	if _, err := s.db.ExecContext(ctx, `
		UPDATE photo_report_deliveries
		SET status = 'failed', last_error = 'interrupted'
		WHERE status = 'sending' AND lease_until < now() AND attempts >= $1`, MaxAttempts); err != nil {
		return 0, err
	}

	sent := 0
	for {
		claims, err := s.claimDue(ctx)
		if err != nil {
			return sent, err
		}
		for _, c := range claims {
			if s.deliver(ctx, c) {
				sent++
			}
		}
		if len(claims) < batchSize {
			return sent, nil
		}
	}
}

func (s *Service) claimDue(ctx context.Context) ([]claim, error) {
	rows, err := s.db.QueryContext(ctx, `
		UPDATE photo_report_deliveries d
		SET status = 'sending', attempts = d.attempts + 1, lease_until = now() + make_interval(secs => $2)
		WHERE (d.report_id, d.recipient_email) IN (
			SELECT report_id, recipient_email FROM photo_report_deliveries
			WHERE attempts < $1
			  AND ((status IN ('pending', 'failed') AND next_attempt_at <= now())
			    OR (status = 'sending' AND lease_until < now()))
			ORDER BY next_attempt_at, report_id
			LIMIT $3
			FOR UPDATE SKIP LOCKED)
		RETURNING d.report_id, d.recipient_email, d.recipient_role, d.language, d.attempts`,
		MaxAttempts, s.lease.Seconds(), batchSize)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []claim
	for rows.Next() {
		var c claim
		if err := rows.Scan(&c.reportID, &c.email, &c.role, &c.language, &c.attempts); err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

// deliver sends one claimed email and records the outcome. Only a nil error from the email service marks it sent.
func (s *Service) deliver(ctx context.Context, c claim) bool {
	err := s.send(ctx, c)
	if err == nil {
		if _, uerr := s.db.ExecContext(ctx, `
			UPDATE photo_report_deliveries
			SET status = 'sent', sent_at = now(), lease_until = NULL, last_error = NULL
			WHERE report_id = $1 AND recipient_email = $2`, c.reportID, c.email); uerr != nil {
			log.Printf("Photo report %d: sent but not recorded: %v", c.reportID, uerr)
		}
		return true
	}
	// The error comes from the email providers; the report text never reaches a log line.
	log.Printf("Photo report %d: delivery attempt %d failed", c.reportID, c.attempts)
	if _, uerr := s.db.ExecContext(ctx, `
		UPDATE photo_report_deliveries
		SET status = 'failed', lease_until = NULL, last_error = $3,
		    next_attempt_at = now() + make_interval(secs => $4)
		WHERE report_id = $1 AND recipient_email = $2`,
		c.reportID, c.email, truncateRunes(err.Error(), 500), retryDelay(c.attempts).Seconds()); uerr != nil {
		log.Printf("Photo report %d: failure not recorded: %v", c.reportID, uerr)
	}
	return false
}

func (s *Service) send(ctx context.Context, c claim) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("email sender panicked: %v", r)
		}
	}()
	data, err := s.emailData(ctx, c)
	if err != nil {
		return err
	}
	return s.mailer.SendEmail(models.Email{
		Recipient:  c.email,
		Template:   s.template,
		SubjectKey: "photo_report_email.subject",
		Language:   c.language,
	}, data)
}

// EmailData feeds the template. Free text (comment, names) is only ever rendered by html/template, which
// escapes it; it is never passed through a translation message, which would not.
type EmailData struct {
	SubjectTitle string
	ReportID     int64
	ReasonCode   string
	Comment      string
	PhotoID      string
	PhotoName    string
	PhotoURL     string
	AlbumID      int64
	AlbumTitle   string
	AlbumURL     string
	ClubID       int64
	ClubName     string
	ReporterName string
	ReporterMail string
	ReporterID   int
	ReporterYear int
	Formation    string
	SubmittedAt  string
	AdminURL     string
	IsClubOwner  bool
	HasClub      bool
}

func (s *Service) emailData(ctx context.Context, c claim) (EmailData, error) {
	var d EmailData
	var albumID, clubID, year sql.NullInt64
	var email sql.NullString
	var created time.Time
	err := s.db.QueryRowContext(ctx, `
		SELECT photo_id, photo_name, photo_url, album_id, album_title, album_url, club_id, club_name, reason, comment,
		       reporter_email, reporter_id, reporter_name, reporter_graduation_year, reporter_formation, created_at
		FROM photo_reports WHERE id_photo_reports = $1`, c.reportID).
		Scan(&d.PhotoID, &d.PhotoName, &d.PhotoURL, &albumID, &d.AlbumTitle, &d.AlbumURL, &clubID, &d.ClubName,
			&d.ReasonCode, &d.Comment, &email, &d.ReporterID, &d.ReporterName, &year, &d.Formation, &created)
	if err != nil {
		return d, err
	}
	d.ReportID = c.reportID
	d.AlbumID, d.ClubID, d.ReporterMail = albumID.Int64, clubID.Int64, email.String
	d.ReporterYear = int(year.Int64)
	d.HasClub = clubID.Valid
	d.IsClubOwner = c.role == "club_owner"
	d.SubmittedAt = created.UTC().Format("2006-01-02 15:04 UTC")
	if s.adminURL != "" {
		d.AdminURL = fmt.Sprintf("%s/photo-reports?report=%d", s.adminURL, c.reportID)
	}
	d.SubjectTitle = oneLine(firstNonEmpty(d.AlbumTitle, d.PhotoName), 80)
	return d, nil
}

// RetryFailed gives every exhausted delivery of a report a fresh set of attempts and sends it right away.
func (s *Service) RetryFailed(ctx context.Context, reportID int64) (int, error) {
	res, err := s.db.ExecContext(ctx, `
		UPDATE photo_report_deliveries
		SET status = 'pending', attempts = 0, next_attempt_at = now(), lease_until = NULL
		WHERE report_id = $1 AND status = 'failed'`, reportID)
	if err != nil {
		return 0, err
	}
	n, _ := res.RowsAffected()
	if n > 0 {
		_, err = s.DeliverDue(ctx)
	}
	return int(n), err
}

// oneLine keeps a value usable in an email subject: no line breaks (header injection) and a bounded length.
func oneLine(v string, n int) string {
	v = strings.Map(func(r rune) rune {
		if unicode.IsControl(r) {
			return ' '
		}
		return r
	}, v)
	return truncateRunes(strings.Join(strings.Fields(v), " "), n)
}

func truncateRunes(s string, n int) string {
	r := []rune(s)
	if len(r) <= n {
		return s
	}
	return string(r[:n-1]) + "…"
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if strings.TrimSpace(v) != "" {
			return v
		}
	}
	return ""
}
