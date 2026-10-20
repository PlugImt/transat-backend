package reservations

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"
)

// release ends the open loan of an item and tells the queue, as returning the item does.
func (f *fx) release(it Item) {
	f.t.Helper()
	f.exec(`UPDATE reservation SET end_date = $1 WHERE id_reservation_element = $2 AND end_date IS NULL`, f.now, it.ID)
	if _, err := f.svc.ReleaseOf(context.Background(), it.ID, nil); err != nil {
		f.t.Fatal(err)
	}
}

func (f *fx) tick() {
	f.t.Helper()
	if _, err := f.svc.Tick(context.Background()); err != nil {
		f.t.Fatal(err)
	}
}

func (f *fx) deadlineOf(email string) time.Time {
	f.t.Helper()
	var at sql.NullTime
	if err := f.db.QueryRow(`SELECT turn_expires_at FROM reservation_waitlist WHERE email = $1`, email).Scan(&at); err != nil || !at.Valid {
		f.t.Fatalf("deadline of %s: %v, %v", email, at, err)
	}
	return at.Time.UTC()
}

func (f *fx) queued(email string) bool {
	return f.count(`SELECT count(*) FROM reservation_waitlist WHERE email = $1`, email) == 1
}

func queueUp(f *fx, it Item, start *time.Time, users ...string) {
	for _, u := range users {
		f.user(u+"@x.test", u, "X", 1)
		if _, err := f.svc.Join(context.Background(), u+"@x.test", it.ID, start); err != nil {
			f.t.Fatal(err)
		}
	}
}

func TestAnItemTurnLastsHalfAnHourThenPassesToTheNextUntilTheQueueIsEmpty(t *testing.T) {
	f := newFx(t)
	f.user("owner@x.test", "O", "X", 1)
	cable := f.item("Câble", false)
	f.lend("owner@x.test", cable, t0.Add(-time.Hour))
	queueUp(f, cable, nil, "q1", "q2", "q3")
	f.release(cable)

	if got := f.ntf.all(); len(got) != 1 || to(t, got[0]) != "q1@x.test" {
		t.Fatalf("first notification = %+v", got)
	}
	if got := f.deadlineOf("q1@x.test"); !got.Equal(t0.Add(30 * time.Minute)) {
		t.Fatalf("q1 has until %v", got)
	}

	// Still inside the delay: nothing moves.
	f.now = t0.Add(29 * time.Minute)
	f.tick()
	if len(f.ntf.all()) != 1 || !f.queued("q1@x.test") {
		t.Fatal("the turn was passed before the delay")
	}

	// q1 did not answer: out of the queue, q2 is told with the buttons.
	f.now = t0.Add(31 * time.Minute)
	f.tick()
	got := f.ntf.all()
	if len(got) != 2 || to(t, got[1]) != "q2@x.test" || got[1].ActionCategory != WaitlistActionCategory {
		t.Fatalf("after the first expiry: %+v", got)
	}
	if f.queued("q1@x.test") {
		t.Fatal("q1 kept a place after not answering")
	}
	if _, err := f.svc.Reserve(context.Background(), "q1@x.test", cable.ID, nil); !errors.Is(err, ErrNotQueued) {
		t.Fatalf("q1 answering too late: %v", err)
	}
	if got := f.deadlineOf("q2@x.test"); !got.Equal(t0.Add(61 * time.Minute)) {
		t.Fatalf("q2 has until %v", got)
	}

	// Ticking again at once changes nothing: one notification per turn.
	f.tick()
	if len(f.ntf.all()) != 2 {
		t.Fatal("a turn was notified twice")
	}

	f.now = t0.Add(62 * time.Minute)
	f.tick()
	if got := f.ntf.all(); len(got) != 3 || to(t, got[2]) != "q3@x.test" {
		t.Fatalf("after the second expiry: %+v", got)
	}

	// The last one does not answer either: the queue is empty and nobody is notified.
	f.now = t0.Add(93 * time.Minute)
	f.tick()
	if len(f.ntf.all()) != 3 || f.count(`SELECT count(*) FROM reservation_waitlist`) != 0 {
		t.Fatalf("%d notifications, %d entries", len(f.ntf.all()), f.count(`SELECT count(*) FROM reservation_waitlist`))
	}
}

func TestTheTurnForASlotThatStartsSoonLastsOnlyAQuarterOfAnHour(t *testing.T) {
	f := newFx(t)
	f.user("owner@x.test", "O", "X", 1)
	room := f.item("Salle", true)
	soon, later := t0.Add(time.Hour), t0.Add(26*time.Hour)
	f.book("owner@x.test", room, soon)
	f.book("owner@x.test", room, later)
	queueUp(f, room, &soon, "q1", "q2")
	// The same person may also wait for the far slot.
	if _, err := f.svc.Join(context.Background(), "q1@x.test", room.ID, &later); err != nil {
		t.Fatal(err)
	}
	for _, s := range []time.Time{soon, later} {
		if _, err := f.svc.Cancel(context.Background(), CancelRequest{ItemID: room.ID, Start: &s, ActorEmail: "owner@x.test"}); err != nil {
			t.Fatal(err)
		}
	}

	var near, far time.Time
	if err := f.db.QueryRow(`SELECT turn_expires_at FROM reservation_waitlist WHERE slot_start = $1`, soon).Scan(&near); err != nil {
		t.Fatal(err)
	}
	if err := f.db.QueryRow(`SELECT turn_expires_at FROM reservation_waitlist WHERE slot_start = $1`, later).Scan(&far); err != nil {
		t.Fatal(err)
	}
	if !near.Equal(t0.Add(15*time.Minute)) || !far.Equal(t0.Add(30*time.Minute)) {
		t.Fatalf("near slot until %v, far slot until %v", near, far)
	}

	f.now = t0.Add(16 * time.Minute)
	f.ntf.reset()
	f.tick()
	got := f.ntf.all()
	if len(got) != 1 || to(t, got[0]) != "q2@x.test" {
		t.Fatalf("the near slot did not move on: %+v", got)
	}
	// The far slot still waits for its person.
	if n := f.count(`SELECT count(*) FROM reservation_waitlist WHERE slot_start = $1 AND notified_at IS NOT NULL`, later); n != 1 {
		t.Fatalf("the far slot's turn was touched: %d", n)
	}
}

func TestIfSomebodyBookedMeanwhileTheCascadeStopsAndTheQueueIsKept(t *testing.T) {
	f := newFx(t)
	f.user("owner@x.test", "O", "X", 1)
	f.user("fast@x.test", "F", "X", 1)
	cable := f.item("Câble", false)
	f.lend("owner@x.test", cable, t0.Add(-time.Hour))
	queueUp(f, cable, nil, "q1", "q2")
	f.release(cable)
	f.ntf.reset()

	// First come, first served: somebody takes the item without going through the queue.
	f.lend("fast@x.test", cable, t0.Add(5*time.Minute))
	f.now = t0.Add(31 * time.Minute)
	f.tick()

	if len(f.ntf.all()) != 0 {
		t.Fatalf("the next person was told about an item that is taken: %+v", f.ntf.all())
	}
	if !f.queued("q1@x.test") || !f.queued("q2@x.test") {
		t.Fatal("the queue was not kept")
	}
	if n := f.count(`SELECT count(*) FROM reservation_waitlist WHERE notified_at IS NOT NULL OR turn_expires_at IS NOT NULL`); n != 0 {
		t.Fatalf("%d entries still hold a turn that is moot", n)
	}

	// The next release starts again from the head of the queue.
	f.exec(`UPDATE reservation SET end_date = $1 WHERE email = 'fast@x.test'`, f.now)
	if _, err := f.svc.ReleaseOf(context.Background(), cable.ID, nil); err != nil {
		t.Fatal(err)
	}
	if got := f.ntf.all(); len(got) != 1 || to(t, got[0]) != "q1@x.test" {
		t.Fatalf("after the next release: %+v", got)
	}
}

func TestAnswerGivenInTimeIsNotUndoneByTheExpiry(t *testing.T) {
	f := newFx(t)
	f.user("owner@x.test", "O", "X", 1)
	cable := f.item("Câble", false)
	f.lend("owner@x.test", cable, t0.Add(-time.Hour))
	queueUp(f, cable, nil, "q1", "q2")
	f.release(cable)
	f.ntf.reset()

	f.now = t0.Add(10 * time.Minute)
	if _, err := f.svc.Reserve(context.Background(), "q1@x.test", cable.ID, nil); err != nil {
		t.Fatal(err)
	}
	f.now = t0.Add(45 * time.Minute)
	f.tick()
	if len(f.ntf.all()) != 0 || !f.queued("q2@x.test") {
		t.Fatalf("q2 was told about an item q1 booked: %+v", f.ntf.all())
	}
	if f.holder(cable, nil) != "q1@x.test" {
		t.Fatal("q1 lost the item")
	}
}

func TestWhoeverIsToldLateInTheCascadeGetsAFullDelay(t *testing.T) {
	f := newFx(t)
	f.user("owner@x.test", "O", "X", 1)
	cable := f.item("Câble", false)
	f.lend("owner@x.test", cable, t0.Add(-time.Hour))
	queueUp(f, cable, nil, "q1", "q2", "q3")
	f.release(cable)

	f.now = t0.Add(20 * time.Minute)
	if next, err := f.svc.Remove(context.Background(), "q1@x.test", cable.ID, nil); err != nil || !next {
		t.Fatalf("q1 declines: %v, %v", next, err)
	}
	if got := f.deadlineOf("q2@x.test"); !got.Equal(t0.Add(50 * time.Minute)) {
		t.Fatalf("q2 has until %v", got)
	}

	f.now = t0.Add(31 * time.Minute)
	f.tick()
	if len(f.ntf.all()) != 2 || !f.queued("q2@x.test") {
		t.Fatalf("q2 lost their turn early: %+v", f.ntf.all())
	}
	f.now = t0.Add(51 * time.Minute)
	f.tick()
	if got := f.ntf.all(); len(got) != 3 || to(t, got[2]) != "q3@x.test" || f.queued("q2@x.test") {
		t.Fatalf("notifications = %+v", got)
	}
}
