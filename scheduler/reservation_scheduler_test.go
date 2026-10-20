package scheduler

import (
	"context"
	"testing"
	"time"

	"github.com/plugimt/transat-backend/internal/database/dbtest"
	"github.com/plugimt/transat-backend/models"
	"github.com/plugimt/transat-backend/services"
	"github.com/plugimt/transat-backend/services/push"
	"github.com/plugimt/transat-backend/services/reservations"
)

type recordingNotifier struct{ sent []services.Announcement }

func (n *recordingNotifier) Announce(_ context.Context, a services.Announcement) (push.EnqueueResult, error) {
	n.sent = append(n.sent, a)
	return push.EnqueueResult{ID: int64(len(n.sent))}, nil
}
func (n *recordingNotifier) CanReach(context.Context, string, models.NotificationCategory) (bool, error) {
	return true, nil
}

// The scheduler is what runs the reminders in production: one tick queues what is due, and a second tick, or a
// second instance, queues nothing more.
func TestTheReservationSchedulerQueuesEachReminderOnce(t *testing.T) {
	db := dbtest.Open(t)
	exec := func(q string, args ...any) {
		t.Helper()
		if _, err := db.Exec(q, args...); err != nil {
			t.Fatalf("%s: %v", q, err)
		}
	}
	var club, loan, room int
	exec(`INSERT INTO newf (email, password, first_name, last_name, language) VALUES ('bob@x.test', 'x', 'Bob', 'M', 1)`)
	if err := db.QueryRow(`INSERT INTO clubs (name, picture) VALUES ('Club', 'p') RETURNING id_clubs`).Scan(&club); err != nil {
		t.Fatal(err)
	}
	db.QueryRow(`INSERT INTO reservation_element (name, slot, id_clubs) VALUES ('Câble', FALSE, $1) RETURNING id_reservation_element`, club).Scan(&loan)
	db.QueryRow(`INSERT INTO reservation_element (name, slot, id_clubs) VALUES ('Salle', TRUE, $1) RETURNING id_reservation_element`, club).Scan(&room)

	now := time.Now().UTC()
	// A loan one day old, and a slot that starts in a little under an hour (booked yesterday).
	exec(`INSERT INTO reservation (email, id_reservation_element, start_date) VALUES ('bob@x.test', $1, $2)`, loan, now.Add(-25*time.Hour))
	slotStart := now.Add(59 * time.Minute).Truncate(time.Minute)
	exec(`INSERT INTO reservation (email, id_reservation_element, start_date, end_date, created_at) VALUES ('bob@x.test', $1, $2, $3, now() - interval '1 day')`,
		room, slotStart, slotStart.Add(time.Hour))

	notifier := &recordingNotifier{}
	scheduler := NewReservationScheduler(reservations.New(db, notifier, nil))
	scheduler.tick()
	scheduler.tick()
	other := NewReservationScheduler(reservations.New(db, notifier, nil))
	other.tick()

	kinds := map[string]int{}
	for _, a := range notifier.sent {
		kinds[a.Kind]++
		if a.Category != models.NotificationReservations {
			t.Fatalf("category %q", a.Category)
		}
	}
	if len(notifier.sent) != 2 || kinds["reservation_reminder"] != 2 {
		t.Fatalf("%d reminders queued (%v), want the loan's and the slot's, once each", len(notifier.sent), kinds)
	}

	// The loan is given back before the next milestone: nothing more is ever queued for it.
	exec(`UPDATE reservation SET end_date = now() AT TIME ZONE 'UTC' WHERE id_reservation_element = $1`, loan)
	scheduler.tick()
	if len(notifier.sent) != 2 {
		t.Fatalf("a reminder was queued after the return: %d", len(notifier.sent))
	}

	scheduler.Start()
	scheduler.Stop()
}
