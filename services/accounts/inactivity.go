package accounts

import (
	"context"
	"log"
	"time"

	"github.com/plugimt/transat-backend/models"
)

const (
	// InactivityLimit is how long an account may stay unused before it is deleted.
	InactivityLimit = 5 * 365 * 24 * time.Hour
	// WarningNotice is how long before the deletion the owner is told by email.
	WarningNotice = 30 * 24 * time.Hour
	// batch keeps one run short; the job runs every day.
	batch = 200
)

// Mailer sends one templated email.
type Mailer interface {
	SendEmail(mail models.Email, data interface{}) error
}

// InactivityEmail feeds email_templates/email_template_account_inactivity.html.
type InactivityEmail struct {
	FirstName    string
	DeletionDate string
}

// PurgeReport says what one run did.
type PurgeReport struct {
	Warned  int
	Deleted int
}

// PurgeInactive warns the accounts that reach 5 years without activity in the next 30 days, then deletes those
// whose warning is at least 30 days old and that stayed unused since. Any sign-in after the warning cancels it,
// because last_activity then passes deletion_warned_at.
func (s *Service) PurgeInactive(ctx context.Context, mailer Mailer) (PurgeReport, error) {
	var rep PurgeReport
	now := s.now().UTC()
	if mailer != nil {
		rows, err := s.db.QueryContext(ctx, `
			SELECT n.email, n.first_name, COALESCE(l.code, 'fr'), COALESCE(n.last_activity, n.creation_date)
			FROM newf n LEFT JOIN languages l ON l.id_languages = n.language
			WHERE n.email <> $1
			  AND COALESCE(n.last_activity, n.creation_date) < $2
			  AND (n.deletion_warned_at IS NULL OR n.deletion_warned_at < COALESCE(n.last_activity, n.creation_date))
			ORDER BY COALESCE(n.last_activity, n.creation_date) LIMIT $3`,
			DeletedAccountEmail, now.Add(-InactivityLimit+WarningNotice), batch)
		if err != nil {
			return rep, err
		}
		type target struct {
			email, first, lang string
			last               time.Time
		}
		var targets []target
		for rows.Next() {
			var t target
			if err := rows.Scan(&t.email, &t.first, &t.lang, &t.last); err != nil {
				_ = rows.Close()
				return rep, err
			}
			targets = append(targets, t)
		}
		_ = rows.Close()
		for _, t := range targets {
			// The account gets the full notice from now, even if its five years end sooner.
			deletion := maxTime(t.last.Add(InactivityLimit), now.Add(WarningNotice))
			err := mailer.SendEmail(models.Email{
				Recipient:  t.email,
				Template:   "email_templates/email_template_account_inactivity.html",
				SubjectKey: "account_inactivity_email.subject",
				Language:   t.lang,
			}, InactivityEmail{FirstName: t.first, DeletionDate: deletion.Format("02/01/2006")})
			if err != nil {
				log.Printf("Accounts: inactivity warning not sent: %v", err)
				continue
			}
			if _, err := s.db.ExecContext(ctx, `UPDATE newf SET deletion_warned_at = $2 WHERE email = $1`, t.email, now); err != nil {
				return rep, err
			}
			rep.Warned++
		}
	}

	due, err := s.dueForDeletion(ctx, now)
	if err != nil {
		return rep, err
	}
	for _, email := range due {
		if err := s.Delete(ctx, email); err != nil {
			log.Printf("Accounts: inactive account not deleted: %v", err)
			continue
		}
		rep.Deleted++
	}
	return rep, nil
}

func (s *Service) dueForDeletion(ctx context.Context, now time.Time) ([]string, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT email FROM newf
		WHERE email <> $1
		  AND COALESCE(last_activity, creation_date) < $2
		  AND deletion_warned_at IS NOT NULL
		  AND deletion_warned_at >= COALESCE(last_activity, creation_date)
		  AND deletion_warned_at <= $3
		ORDER BY email LIMIT $4`,
		DeletedAccountEmail, now.Add(-InactivityLimit), now.Add(-WarningNotice), batch)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []string
	for rows.Next() {
		var e string
		if err := rows.Scan(&e); err != nil {
			return nil, err
		}
		out = append(out, e)
	}
	return out, rows.Err()
}

func maxTime(a, b time.Time) time.Time {
	if a.After(b) {
		return a
	}
	return b
}
