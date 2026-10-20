// Package reservations holds the rules around reservations that go beyond booking: who may cancel what, the
// notifications and emails that follow, the reminders sent by the scheduler, and the waiting list of people
// interested in something that is already reserved.
//
// Times. The reservation tables store UTC wall-clock times in TIMESTAMP columns, and the API speaks the same
// "2006-01-02 15:04:05" UTC format. Everything here passes UTC to the database and converts to Paris time only
// for the text shown to people.
package reservations

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/plugimt/transat-backend/models"
	"github.com/plugimt/transat-backend/services"
	"github.com/plugimt/transat-backend/services/push"
)

// APITimeLayout is the layout of start_date / end_date in the API (UTC wall time).
const APITimeLayout = "2006-01-02 15:04:05"

// WaitlistActionCategory is the interactive notification category the app registers with its two buttons,
// "Réserver" and "Annuler". The notification's navigation carries what the buttons need (see notifyNext).
const WaitlistActionCategory = "reservation_waitlist"

const (
	emailTemplate    = "email_templates/email_template_reservation_cancelled.html"
	emailSubjectKey  = "reservation_cancelled_email.subject"
	maxQueueEntries  = 30
	maxItemNameRunes = 80
	// lockNamespace is the first key of the advisory locks that serialize bookings and waiting-list changes of an item.
	lockNamespace = 7301
)

var (
	ErrItemNotFound   = errors.New("reservation item not found")
	ErrStartRequired  = errors.New("start_date is required for a slot item")
	ErrInvalidStart   = errors.New("start_date must be a future slot on a plain hour")
	ErrNotSlotItem    = errors.New("item is not a per slot item, cannot delete reservation")
	ErrNotFound       = errors.New("no reservation found, or you are not its owner")
	ErrNotAvailable   = errors.New("item is not available")
	ErrStillAvailable = errors.New("item is available: reserve it directly")
	ErrOwnReservation = errors.New("you already hold this reservation")
	ErrAlreadyQueued  = errors.New("you are already on the waiting list")
	ErrNotQueued      = errors.New("you are not on the waiting list")
	ErrQueueFull      = errors.New("too many waiting list entries")
)

// Notifier is the part of the notification service this package needs.
type Notifier interface {
	OutboxEnabled() bool
	Announce(ctx context.Context, a services.Announcement) (push.EnqueueResult, error)
	// CanReach reports whether a push to this user would go to at least one device right now.
	CanReach(ctx context.Context, email string, category models.NotificationCategory) (bool, error)
}

// Mailer sends a templated email.
type Mailer interface {
	SendEmail(mail models.Email, data interface{}) error
}

// Item is a reservable element.
type Item struct {
	ID   int
	Name string
	Slot bool
}

// Service implements the reservation rules on top of the database.
type Service struct {
	db       *sql.DB
	notifier Notifier
	mailer   Mailer
	now      func() time.Time
	// run executes slow side effects (emails). It runs them on their own goroutine unless a test says otherwise.
	run func(func())
}

// New builds the service. The mailer may be nil: no email is then sent.
func New(db *sql.DB, notifier Notifier, mailer Mailer) *Service {
	return &Service{db: db, notifier: notifier, mailer: mailer, now: time.Now, run: func(f func()) { go f() }}
}

// SetBackground replaces how slow side effects run; tests pass a function that runs them inline.
func (s *Service) SetBackground(run func(func())) { s.run = run }

// SetClock replaces the clock; for tests.
func (s *Service) SetClock(now func() time.Time) { s.now = now }

func (s *Service) background(label string, fn func(ctx context.Context) error) {
	s.run(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Reservation %s panicked: %v", label, r)
			}
		}()
		if err := fn(ctx); err != nil {
			log.Printf("Reservation %s failed: %v", label, err)
		}
	})
}

type queryRower interface {
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

func (s *Service) loadItem(ctx context.Context, id int) (Item, error) {
	it := Item{ID: id}
	err := s.db.QueryRowContext(ctx,
		`SELECT COALESCE(name, ''), COALESCE(slot, FALSE) FROM reservation_element WHERE id_reservation_element = $1`, id).
		Scan(&it.Name, &it.Slot)
	if errors.Is(err, sql.ErrNoRows) {
		return it, ErrItemNotFound
	}
	return it, err
}

// IsAdmin reports whether the user has the ADMIN role.
func (s *Service) IsAdmin(ctx context.Context, email string) (bool, error) {
	var admin bool
	err := s.db.QueryRowContext(ctx, `
		SELECT EXISTS(
			SELECT 1 FROM newf_roles nr JOIN roles r ON nr.id_roles = r.id_roles
			WHERE nr.email = $1 AND r.name = 'ADMIN')`, email).Scan(&admin)
	return admin, err
}

// displayName is "First Last" for a user, empty when the account does not exist.
func (s *Service) displayName(ctx context.Context, email string) string {
	var first, last sql.NullString
	if err := s.db.QueryRowContext(ctx, `SELECT first_name, last_name FROM newf WHERE email = $1`, email).Scan(&first, &last); err != nil {
		return ""
	}
	return strings.TrimSpace(first.String + " " + last.String)
}

// userLanguage is the language code of a user, "fr" when unknown.
func (s *Service) userLanguage(ctx context.Context, email string) string {
	var code sql.NullString
	err := s.db.QueryRowContext(ctx, `
		SELECT l.code FROM newf n LEFT JOIN languages l ON l.id_languages = n.language WHERE n.email = $1`, email).Scan(&code)
	if err != nil || !code.Valid || code.String == "" {
		return push.DefaultLanguage
	}
	return code.String
}

func lockItem(ctx context.Context, tx *sql.Tx, itemID int) error {
	_, err := tx.ExecContext(ctx, `SELECT pg_advisory_xact_lock($1::int, $2::int)`, lockNamespace, itemID)
	return err
}

// isFree reports whether nobody holds the item (loan) or the slot right now.
func isFree(ctx context.Context, q queryRower, item Item, start *time.Time) (bool, error) {
	var free bool
	var err error
	if item.Slot {
		err = q.QueryRowContext(ctx, `
			SELECT NOT EXISTS (SELECT 1 FROM reservation WHERE id_reservation_element = $1 AND start_date = $2)`,
			item.ID, start.UTC()).Scan(&free)
	} else {
		err = q.QueryRowContext(ctx, `
			SELECT NOT EXISTS (SELECT 1 FROM reservation WHERE id_reservation_element = $1 AND end_date IS NULL)`,
			item.ID).Scan(&free)
	}
	return free, err
}

// target normalizes the slot of a request: a slot item needs a future slot on a plain hour, a loan item has none.
func (s *Service) target(item Item, start *time.Time) (*time.Time, error) {
	if !item.Slot {
		return nil, nil
	}
	if start == nil {
		return nil, ErrStartRequired
	}
	u := start.UTC()
	if u.Minute() != 0 || u.Second() != 0 || u.Nanosecond() != 0 || !u.After(s.now()) {
		return nil, ErrInvalidStart
	}
	return &u, nil
}

func startArg(start *time.Time) any {
	if start == nil {
		return nil
	}
	return start.UTC()
}

// FormatSlot writes a time for people, in Paris time.
func FormatSlot(lang string, t time.Time) string {
	p := t.In(parisLocation())
	if lang == "en" {
		return p.Format("2 Jan 2006, 15:04")
	}
	return p.Format("02/01/2006 15:04")
}

func parisLocation() *time.Location {
	loc, err := time.LoadLocation("Europe/Paris")
	if err != nil {
		return time.UTC
	}
	return loc
}

func parisDate(t time.Time) string { return t.In(parisLocation()).Format("2006-01-02") }

func truncateRunes(s string, n int) string {
	r := []rune(s)
	if len(r) <= n {
		return s
	}
	return string(r[:n])
}

// navigation opens the item's reservation page, on the day of the slot when there is one.
func navigation(item Item, start *time.Time, extra map[string]string) *models.NavigationTarget {
	params := map[string]string{"title": truncateRunes(item.Name, maxItemNameRunes)}
	if start != nil {
		params["date"] = parisDate(*start)
	}
	for k, v := range extra {
		params[k] = v
	}
	return models.ReservationNavigation(item.ID, params)
}

func usersAudience(email string) push.Audience {
	return push.Audience{NotificationAudience: models.NotificationAudience{Type: models.AudienceUsers, Emails: []string{email}}}
}

func dedupe(parts ...any) string {
	var b strings.Builder
	for i, p := range parts {
		if i > 0 {
			b.WriteByte(':')
		}
		fmt.Fprint(&b, p)
	}
	return b.String()
}
