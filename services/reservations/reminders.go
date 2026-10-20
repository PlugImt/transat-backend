package reservations

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	appI18n "github.com/plugimt/transat-backend/i18n"
	"github.com/plugimt/transat-backend/models"
	"github.com/plugimt/transat-backend/services"
)

// Reminder milestones. The value is stored in reservation_reminders.milestone.
const (
	milestoneSlot1h = "slot_1h"
	milestoneLoanD1 = "loan_d1"
	milestoneLoanD5 = "loan_d5"
	milestoneLoanD7 = "loan_d7"
)

const (
	// slotReminderLead is how long before a slot the reminder goes out.
	slotReminderLead = time.Hour
	// slotReminderLatest: a reminder that could only go out less than this before the slot is dropped.
	slotReminderLatest = 30 * time.Minute
	// reminderCatchUp is how late a reminder may still go out, so that a restart or a short outage does not lose
	// it. It also keeps the first run after a deployment from reminding about loans forgotten long ago.
	reminderCatchUp = 24 * time.Hour
)

type loanMilestone struct {
	name  string
	after string // PostgreSQL interval
	days  int
}

var loanMilestones = []loanMilestone{
	{milestoneLoanD1, "1 day", 1},
	{milestoneLoanD5, "5 days", 5},
	{milestoneLoanD7, "7 days", 7},
}

type reminder struct {
	email     string
	itemID    int
	itemName  string
	start     time.Time
	milestone string
}

// Tick does the periodic work: it sends the reminders that are due and forgets the waiting list entries of slots
// that are over. It is safe to call as often as wanted, and from several processes.
func (s *Service) Tick(ctx context.Context) (sent int, err error) {
	if _, err := s.db.ExecContext(ctx,
		`DELETE FROM reservation_waitlist WHERE slot_start IS NOT NULL AND slot_start <= $1`, s.now().UTC()); err != nil {
		log.Printf("Reservation waiting list cleanup failed: %v", err)
	}
	return s.SendDueReminders(ctx)
}

// SendDueReminders queues the reminders that are due now:
//   - slot reservations: one reminder an hour before the slot starts (not if it was booked less than an hour
//     before the start, and not once it is too late to be useful);
//   - loans of items without a slot: a reminder one day, five days and seven days after the loan started, as long
//     as the item has not been returned.
//
// Each reminder is claimed with an insert into reservation_reminders, whose primary key makes a second claim of
// the same milestone impossible; only the claim that inserted the row sends. A reservation that was cancelled or
// returned does not match the claim query, so it gets nothing. A claim whose push could not be queued is
// released, and the next run tries again.
func (s *Service) SendDueReminders(ctx context.Context) (int, error) {
	if s.notifier == nil || !s.notifier.OutboxEnabled() {
		return 0, nil
	}
	now := s.now().UTC()
	var due []reminder

	slot, err := s.claim(ctx, `
		INSERT INTO reservation_reminders (email, id_reservation_element, start_date, milestone)
		SELECT r.email, r.id_reservation_element, r.start_date, $3::varchar
		FROM reservation r
		JOIN reservation_element re ON re.id_reservation_element = r.id_reservation_element
		WHERE re.slot = TRUE
		  AND r.start_date - make_interval(secs => $2::float8) <= $1::timestamp
		  AND r.start_date > $1::timestamp + make_interval(secs => $4::float8)
		  AND (r.created_at IS NULL OR r.created_at <= (r.start_date AT TIME ZONE 'UTC') - make_interval(secs => $2::float8))
		ON CONFLICT DO NOTHING
		RETURNING email, id_reservation_element, start_date`,
		now, slotReminderLead.Seconds(), milestoneSlot1h, slotReminderLatest.Seconds(), milestoneSlot1h)
	if err != nil {
		return 0, err
	}
	due = append(due, slot...)

	for _, m := range loanMilestones {
		loans, err := s.claim(ctx, `
			INSERT INTO reservation_reminders (email, id_reservation_element, start_date, milestone)
			SELECT r.email, r.id_reservation_element, r.start_date, $3::varchar
			FROM reservation r
			JOIN reservation_element re ON re.id_reservation_element = r.id_reservation_element
			WHERE re.slot = FALSE
			  AND r.end_date IS NULL
			  AND r.start_date + $2::interval <= $1::timestamp
			  AND r.start_date + $2::interval > $1::timestamp - make_interval(secs => $4::float8)
			ON CONFLICT DO NOTHING
			RETURNING email, id_reservation_element, start_date`,
			now, m.after, m.name, reminderCatchUp.Seconds(), m.name)
		if err != nil {
			return len(due), err
		}
		due = append(due, loans...)
	}

	sent := 0
	for _, r := range due {
		if err := s.sendReminder(ctx, r); err != nil {
			log.Printf("Reservation reminder %s for item %d could not be queued: %v", r.milestone, r.itemID, err)
			s.release(ctx, r)
			continue
		}
		sent++
	}
	return sent, nil
}

// claim runs one claim query. The last argument only names the milestone for the returned rows.
func (s *Service) claim(ctx context.Context, query string, args ...any) ([]reminder, error) {
	milestone := args[len(args)-1].(string)
	// The item name is looked up after the claim, so the insert stays a single statement.
	rows, err := s.db.QueryContext(ctx, `WITH claimed AS (`+query+`)
		SELECT c.email, c.id_reservation_element, c.start_date, COALESCE(re.name, '')
		FROM claimed c JOIN reservation_element re ON re.id_reservation_element = c.id_reservation_element`, args[:len(args)-1]...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []reminder
	for rows.Next() {
		r := reminder{milestone: milestone}
		if err := rows.Scan(&r.email, &r.itemID, &r.start, &r.itemName); err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	return out, rows.Err()
}

func (s *Service) release(ctx context.Context, r reminder) {
	if _, err := s.db.ExecContext(ctx, `
		DELETE FROM reservation_reminders
		WHERE email = $1 AND id_reservation_element = $2 AND start_date = $3 AND milestone = $4`,
		r.email, r.itemID, r.start.UTC(), r.milestone); err != nil {
		log.Printf("Reservation reminder %s for item %d: claim not released: %v", r.milestone, r.itemID, err)
	}
}

func (s *Service) sendReminder(ctx context.Context, r reminder) error {
	item := Item{ID: r.itemID, Name: r.itemName}
	now := s.now()
	a := services.Announcement{
		Kind:      "reservation_reminder",
		Category:  models.NotificationReservations,
		Audience:  usersAudience(r.email),
		DedupeKey: dedupe("reservation_reminder", r.milestone, r.itemID, r.start.UTC().Unix(), r.email),
	}
	if r.milestone == milestoneSlot1h {
		item.Slot = true
		start := r.start.UTC()
		a.Navigation = navigation(item, &start, nil)
		a.ExpiresAt = start
		a.Compose = func(l string, loc *i18n.Localizer) (string, string) {
			data := map[string]interface{}{"ItemName": item.Name, "Time": r.start.In(parisLocation()).Format("15:04")}
			return appI18n.Localize(loc, "reservation_notification.slot_reminder_title", "Reservation in 1h", nil),
				appI18n.Localize(loc, "reservation_notification.slot_reminder_message", "{{.ItemName}} · {{.Time}}", data)
		}
	} else {
		days := 0
		for _, m := range loanMilestones {
			if m.name == r.milestone {
				days = m.days
			}
		}
		if days == 0 {
			return fmt.Errorf("unknown milestone %q", r.milestone)
		}
		a.Navigation = navigation(item, nil, nil)
		a.ExpiresAt = now.Add(12 * time.Hour)
		a.Compose = func(l string, loc *i18n.Localizer) (string, string) {
			data := map[string]interface{}{"ItemName": item.Name}
			return appI18n.Localize(loc, "reservation_notification.loan_reminder_d"+fmt.Sprint(days)+"_title", "Item still borrowed", nil),
				appI18n.Localize(loc, "reservation_notification.loan_reminder_d"+fmt.Sprint(days)+"_message", "{{.ItemName}} is still in your name.", data)
		}
	}
	_, err := s.notifier.Announce(ctx, a)
	return err
}
