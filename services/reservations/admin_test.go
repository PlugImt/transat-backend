package reservations

import (
	"context"
	"errors"
	"testing"
	"time"
)

func (f *fx) category(name string, parent *int) int {
	f.t.Helper()
	var id int
	if err := f.db.QueryRow(`INSERT INTO reservation_category (name, id_parent_category, id_clubs) VALUES ($1, $2, $3) RETURNING id_reservation_category`,
		name, parent, f.club).Scan(&id); err != nil {
		f.t.Fatal(err)
	}
	return id
}

func byName(items []AdminItem, name string) AdminItem {
	for _, it := range items {
		if it.Name == name {
			return it
		}
	}
	return AdminItem{}
}

func TestTheCatalogueShowsWhatIsHappeningToEachItem(t *testing.T) {
	f := newFx(t)
	f.user("bob@x.test", "Bob", "Martin", 1)
	f.user("amy@x.test", "Amy", "Roy", 1)
	ctx := context.Background()

	parent := f.category("Salles", nil)
	child := f.category("Etage 1", &parent)
	room := f.item("Salle A", true)
	f.exec(`UPDATE reservation_element SET id_reservation_category = $1, location = 'Bat. 1' WHERE id_reservation_element = $2`, child, room.ID)
	cable := f.item("Câble", false)
	_ = f.item("Projecteur", false)
	broken := f.item("Console", false)
	f.exec(`UPDATE reservation_element SET status = 'maintenance', status_note = 'Écran cassé' WHERE id_reservation_element = $1`, broken.ID)

	f.book("bob@x.test", room, t0)                   // being used now (12:00 to 13:00)
	f.book("amy@x.test", room, t0.Add(3*time.Hour))  // next
	f.book("amy@x.test", room, t0.Add(30*time.Hour)) // later
	f.lend("bob@x.test", cable, t0.Add(-3*24*time.Hour))
	f.exec(`INSERT INTO reservation_waitlist (email, id_reservation_element) VALUES ('amy@x.test', $1)`, cable.ID)

	items, err := f.svc.ListItems(ctx)
	if err != nil || len(items) != 4 {
		t.Fatalf("%d items, %v", len(items), err)
	}

	r := byName(items, "Salle A")
	if r.Availability != "occupied" || r.Current == nil || r.Current.Email != "bob@x.test" || r.Current.FirstName != "Bob" {
		t.Fatalf("room = %+v", r)
	}
	if r.Next == nil || r.Next.Email != "amy@x.test" || !r.Next.StartDate.Equal(t0.Add(3*time.Hour)) || r.UpcomingCount != 2 {
		t.Fatalf("next = %+v, upcoming %d", r.Next, r.UpcomingCount)
	}
	if r.CategoryPath != "Salles / Etage 1" || r.CategoryID == nil || *r.CategoryID != child || r.Location == nil || *r.Location != "Bat. 1" {
		t.Fatalf("path = %q", r.CategoryPath)
	}

	c := byName(items, "Câble")
	if c.Availability != "borrowed" || c.Current == nil || c.Current.ElapsedSeconds != 3*24*3600 || c.WaitlistCount != 1 {
		t.Fatalf("cable = %+v / %+v", c, c.Current)
	}
	if p := byName(items, "Projecteur"); p.Availability != "available" || p.Current != nil || p.CategoryPath != "" {
		t.Fatalf("projector = %+v", p)
	}
	if m := byName(items, "Console"); m.Availability != "maintenance" || m.StatusNote == nil || *m.StatusNote != "Écran cassé" {
		t.Fatalf("console = %+v", m)
	}
}

func TestTheReservationListFiltersEverything(t *testing.T) {
	f := newFx(t)
	f.user("admin@x.test", "Ada", "Lovelace", 1, "ADMIN")
	for _, u := range []struct{ email, first, last string }{
		{"bob@x.test", "Bob", "Martin"}, {"amy@x.test", "Amy", "Roy"}, {"eve@x.test", "Eve", "Durand"},
	} {
		f.user(u.email, u.first, u.last, 1)
	}
	ctx := context.Background()
	room := f.item("Salle A", true)
	cable := f.item("Câble", false)
	laptop := f.item("Portable", false)
	hdmi := f.item("HDMI", false)
	mic := f.item("Micro", false)

	f.lend("bob@x.test", cable, t0.Add(-8*24*time.Hour))  // 8 days: past the 7 day mark
	f.lend("amy@x.test", laptop, t0.Add(-6*24*time.Hour)) // 6 days: past the 5 day mark
	f.lend("eve@x.test", hdmi, t0.Add(-2*24*time.Hour))   // 2 days: past the first day
	f.lend("eve@x.test", mic, t0.Add(-3*time.Hour))       // recent
	f.exec(`INSERT INTO reservation_reminders (email, id_reservation_element, start_date, milestone)
		SELECT email, id_reservation_element, start_date, m FROM reservation, unnest(ARRAY['loan_d1', 'loan_d5']) m
		WHERE email = 'amy@x.test'`)
	f.book("bob@x.test", room, t0.Add(-5*time.Hour)) // over
	f.book("amy@x.test", room, t0)                   // ongoing
	f.book("eve@x.test", room, t0.Add(2*time.Hour))  // upcoming
	f.book("bob@x.test", room, t0.Add(4*time.Hour))  // cancelled by an admin below
	f.book("eve@x.test", room, t0.Add(6*time.Hour))  // cancelled by its owner below
	cancelled := t0.Add(4 * time.Hour)
	if _, err := f.svc.Cancel(ctx, CancelRequest{ItemID: room.ID, Start: &cancelled, ActorEmail: "admin@x.test", AsAdmin: true}); err != nil {
		t.Fatal(err)
	}
	own := t0.Add(6 * time.Hour)
	if _, err := f.svc.Cancel(ctx, CancelRequest{ItemID: room.ID, Start: &own, ActorEmail: "eve@x.test"}); err != nil {
		t.Fatal(err)
	}
	// A loan ended by an administrator, who is not the user: it counts as cancelled.
	f.lend("bob@x.test", f.item("Rallonge", false), t0.Add(-time.Hour))
	if _, err := f.svc.Cancel(ctx, CancelRequest{ItemID: itemID(f, "Rallonge"), ActorEmail: "admin@x.test", AsAdmin: true}); err != nil {
		t.Fatal(err)
	}

	list := func(filter ReservationFilter) ReservationPage {
		t.Helper()
		page, err := f.svc.ListReservations(ctx, filter)
		if err != nil {
			t.Fatal(err)
		}
		return page
	}
	names := func(p ReservationPage) map[string]bool {
		out := map[string]bool{}
		for _, r := range p.Items {
			out[r.ItemName+"/"+r.Email+"/"+r.Status] = true
		}
		return out
	}

	all := list(ReservationFilter{})
	if all.Total != 10 || len(all.Items) != 10 {
		t.Fatalf("%d reservations in all (%v)", all.Total, names(all))
	}

	active := names(list(ReservationFilter{Status: "active"}))
	for _, want := range []string{"Câble/bob@x.test/active", "Portable/amy@x.test/active", "HDMI/eve@x.test/active", "Micro/eve@x.test/active", "Salle A/amy@x.test/ongoing"} {
		if !active[want] {
			t.Fatalf("active is missing %s: %v", want, active)
		}
	}
	if len(active) != 5 {
		t.Fatalf("active = %v", active)
	}
	if got := names(list(ReservationFilter{Status: "upcoming"})); len(got) != 1 || !got["Salle A/eve@x.test/upcoming"] {
		t.Fatalf("upcoming = %v", got)
	}
	if got := names(list(ReservationFilter{Status: "past"})); len(got) != 1 || !got["Salle A/bob@x.test/past"] {
		t.Fatalf("past = %v", got)
	}
	cancelledPage := list(ReservationFilter{Status: "cancelled"})
	if len(cancelledPage.Items) != 3 {
		t.Fatalf("cancelled = %v", names(cancelledPage))
	}
	kinds := map[string]string{}
	for _, r := range cancelledPage.Items {
		if r.CancelKind == nil || r.CancelledBy == nil || r.CancelledAt == nil || r.Cancellable {
			t.Fatalf("cancelled line = %+v", r)
		}
		kinds[r.ItemName+"/"+r.Email] = *r.CancelKind + "/" + *r.CancelledBy
	}
	if kinds["Salle A/bob@x.test"] != "cancelled_by_admin/admin@x.test" || kinds["Salle A/eve@x.test"] != "cancelled/eve@x.test" ||
		kinds["Rallonge/bob@x.test"] != "cancelled_by_admin/admin@x.test" {
		t.Fatalf("kinds = %v", kinds)
	}

	// Loans that last: older than a day (default), 5 days, 7 days.
	if got := names(list(ReservationFilter{Status: "overdue"})); len(got) != 3 {
		t.Fatalf("overdue = %v", got)
	}
	if got := names(list(ReservationFilter{OverdueDays: 5})); len(got) != 2 || !got["Câble/bob@x.test/active"] || !got["Portable/amy@x.test/active"] {
		t.Fatalf("5 days = %v", got)
	}
	seven := list(ReservationFilter{OverdueDays: 7})
	if len(seven.Items) != 1 || seven.Items[0].ItemName != "Câble" || seven.Items[0].OverdueLevel != 7 {
		t.Fatalf("7 days = %+v", seven.Items)
	}
	levels := map[string]int{}
	for _, r := range list(ReservationFilter{Kind: "loan", Status: "active"}).Items {
		levels[r.ItemName] = r.OverdueLevel
	}
	if levels["Câble"] != 7 || levels["Portable"] != 5 || levels["HDMI"] != 1 || levels["Micro"] != 0 {
		t.Fatalf("levels = %v", levels)
	}

	// Reminders already sent are shown with the loan.
	for _, r := range list(ReservationFilter{Kind: "loan", Status: "active", Query: "amy"}).Items {
		if len(r.Reminders) != 2 || r.Reminders[0] != "loan_d1" || r.Reminders[1] != "loan_d5" {
			t.Fatalf("reminders = %v", r.Reminders)
		}
		if r.ElapsedSeconds == nil || *r.ElapsedSeconds != 6*24*3600 {
			t.Fatalf("elapsed = %v", r.ElapsedSeconds)
		}
	}

	// By user (name in either order, or email), by item, by period, by kind.
	for _, q := range []string{"BOB", "martin bob", "bob@x"} {
		if got := list(ReservationFilter{Query: q}); got.Total != 4 {
			t.Fatalf("query %q: %d", q, got.Total)
		}
	}
	if got := list(ReservationFilter{ItemID: room.ID}); got.Total != 5 {
		t.Fatalf("room: %d", got.Total)
	}
	from, to := t0, t0.Add(3*time.Hour)
	if got := names(list(ReservationFilter{ItemID: room.ID, From: &from, To: &to})); len(got) != 2 {
		t.Fatalf("period = %v", got)
	}
	if got := list(ReservationFilter{Kind: "slot"}); got.Total != 5 {
		t.Fatalf("slots: %d", got.Total)
	}

	// Pagination keeps the total, and the order can be reversed.
	page := list(ReservationFilter{Limit: 4, Offset: 4, Ascending: true})
	if page.Total != 10 || len(page.Items) != 4 || page.Offset != 4 {
		t.Fatalf("page = %d/%d", len(page.Items), page.Total)
	}
	first := list(ReservationFilter{Limit: 1, Ascending: true}).Items[0].StartDate
	last := list(ReservationFilter{Limit: 1}).Items[0].StartDate
	if !first.Before(last) {
		t.Fatalf("order: %v then %v", first, last)
	}

	if _, err := f.svc.ListReservations(ctx, ReservationFilter{Status: "nonsense"}); err == nil {
		t.Fatal("an unknown status was accepted")
	}
}

func itemID(f *fx, name string) int {
	var id int
	if err := f.db.QueryRow(`SELECT id_reservation_element FROM reservation_element WHERE name = $1`, name).Scan(&id); err != nil {
		f.t.Fatal(err)
	}
	return id
}

func TestMarkingALoanAsReturnedFreesTheItemAndStopsTheReminders(t *testing.T) {
	f := newFx(t)
	f.user("admin@x.test", "Ada", "Lovelace", 1, "ADMIN")
	f.user("bob@x.test", "Bob", "Martin", 1)
	f.user("amy@x.test", "Amy", "Roy", 1)
	ctx := context.Background()
	cable := f.item("Câble", false)
	room := f.item("Salle", true)
	start := t0.Add(-25 * time.Hour)
	f.lend("bob@x.test", cable, start)
	if _, err := f.svc.Join(ctx, "amy@x.test", cable.ID, nil); err != nil {
		t.Fatal(err)
	}

	if _, err := f.svc.MarkReturned(ctx, room.ID, "admin@x.test"); !errors.Is(err, ErrNotLoanItem) {
		t.Fatalf("a slot item: %v", err)
	}

	res, err := f.svc.MarkReturned(ctx, cable.ID, "admin@x.test")
	if err != nil || res.Owner != "bob@x.test" || !res.NextNotified {
		t.Fatalf("result = %+v, %v", res, err)
	}
	if f.holder(cable, nil) != "" {
		t.Fatal("the item is still held")
	}
	// The one who gave the item back is told nothing; the first of the queue is.
	sent := f.ntf.all()
	if len(sent) != 1 || to(t, sent[0]) != "amy@x.test" || sent[0].Kind != "reservation_waitlist" {
		t.Fatalf("announcements = %+v", sent)
	}
	if len(f.mail.sent) != 0 {
		t.Fatal("an email went out for a return")
	}

	// The reminders of that loan are over, even when their time comes.
	f.ntf.reset()
	f.now = start.Add(7*24*time.Hour + time.Hour)
	if n, err := f.svc.SendDueReminders(ctx); err != nil || n != 0 {
		t.Fatalf("%d reminders after the return, %v", n, err)
	}

	page, _ := f.svc.ListReservations(ctx, ReservationFilter{ItemID: cable.ID})
	if len(page.Items) != 1 || page.Items[0].Status != "past" || page.Items[0].CancelKind == nil || *page.Items[0].CancelKind != "returned_by_admin" ||
		page.Items[0].EndDate == nil || page.Items[0].Cancellable {
		t.Fatalf("history = %+v", page.Items)
	}
	if _, err := f.svc.MarkReturned(ctx, cable.ID, "admin@x.test"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("returning twice: %v", err)
	}
}

func TestTheWaitingListIsVisibleInOrderAndAdministrable(t *testing.T) {
	f := newFx(t)
	for _, u := range []string{"owner", "q1", "q2", "q3"} {
		f.user(u+"@x.test", u, "X", 1)
	}
	ctx := context.Background()
	room := f.item("Salle", true)
	cable := f.item("Câble", false)
	slotA, slotB := t0.Add(26*time.Hour), t0.Add(28*time.Hour)
	f.book("owner@x.test", room, slotA)
	f.book("owner@x.test", room, slotB)
	f.lend("owner@x.test", cable, t0.Add(-time.Hour))
	for _, who := range []string{"q1", "q2", "q3"} {
		if _, err := f.svc.Join(ctx, who+"@x.test", room.ID, &slotA); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := f.svc.Join(ctx, "q3@x.test", room.ID, &slotB); err != nil {
		t.Fatal(err)
	}
	f.svc.Join(ctx, "q2@x.test", cable.ID, nil)

	list, err := f.svc.AdminWaitlist(ctx, room.ID, &slotA)
	if err != nil || len(list.Entries) != 3 || list.Free || list.AwaitingResponse {
		t.Fatalf("list = %+v, %v", list, err)
	}
	for i, want := range []string{"q1", "q2", "q3"} {
		e := list.Entries[i]
		if e.Position != i+1 || e.Email != want+"@x.test" || e.FirstName != want || e.JoinedAt.IsZero() || e.NotifiedAt != nil {
			t.Fatalf("entry %d = %+v", i, e)
		}
	}
	// Without a slot, every upcoming queue of the item, each numbered on its own.
	every, _ := f.svc.AdminWaitlist(ctx, room.ID, nil)
	if len(every.Entries) != 4 || every.Entries[3].SlotStart == nil || !every.Entries[3].SlotStart.Equal(slotB) || every.Entries[3].Position != 1 {
		t.Fatalf("every = %+v", every.Entries)
	}
	if loan, _ := f.svc.AdminWaitlist(ctx, cable.ID, nil); len(loan.Entries) != 1 || loan.Entries[0].Email != "q2@x.test" {
		t.Fatalf("loan = %+v", loan)
	}

	// The slot is not free: the turn cannot be passed.
	if _, err := f.svc.AdminAdvanceWaitlist(ctx, room.ID, &slotA); !errors.Is(err, ErrNotFree) {
		t.Fatalf("advance while occupied: %v", err)
	}

	// Released: q1 holds the baton, and the list says so.
	f.svc.Cancel(ctx, CancelRequest{ItemID: room.ID, Start: &slotA, ActorEmail: "owner@x.test"})
	list, _ = f.svc.AdminWaitlist(ctx, room.ID, &slotA)
	if !list.Free || !list.AwaitingResponse || list.Entries[0].NotifiedAt == nil || list.Entries[1].NotifiedAt != nil {
		t.Fatalf("after release = %+v", list)
	}

	// Pass the turn by hand: q1 is dropped, q2 is told.
	f.ntf.reset()
	next, err := f.svc.AdminAdvanceWaitlist(ctx, room.ID, &slotA)
	if err != nil || !next {
		t.Fatalf("advance: %v, %v", next, err)
	}
	if got := f.ntf.all(); len(got) != 1 || to(t, got[0]) != "q2@x.test" {
		t.Fatalf("announcements = %+v", got)
	}
	list, _ = f.svc.AdminWaitlist(ctx, room.ID, &slotA)
	if len(list.Entries) != 2 || list.Entries[0].Email != "q2@x.test" || list.Entries[0].Position != 1 || list.Entries[0].NotifiedAt == nil {
		t.Fatalf("after advance = %+v", list.Entries)
	}

	// Removing the person who holds the baton passes it on; removing somebody else does not.
	f.ntf.reset()
	if next, err := f.svc.AdminRemoveWaitlistEntry(ctx, list.Entries[1].ID); err != nil || next {
		t.Fatalf("removing q3: %v, %v", next, err)
	}
	f.svc.Join(ctx, "q1@x.test", room.ID, &slotB)
	if len(f.ntf.all()) != 0 {
		t.Fatal("a removal that was not the baton notified somebody")
	}
	if next, err := f.svc.AdminRemoveWaitlistEntry(ctx, list.Entries[0].ID); err != nil || next {
		t.Fatalf("removing q2 with nobody behind: %v, %v", next, err)
	}
	if _, err := f.svc.AdminRemoveWaitlistEntry(ctx, 999999); !errors.Is(err, ErrNotQueued) {
		t.Fatalf("unknown entry: %v", err)
	}
	if _, err := f.svc.AdminWaitlist(ctx, 99999, nil); !errors.Is(err, ErrItemNotFound) {
		t.Fatalf("unknown item: %v", err)
	}
}

func TestAnItemUnderMaintenanceCannotBeQueuedForOrBooked(t *testing.T) {
	f := newFx(t)
	f.user("owner@x.test", "O", "X", 1)
	f.user("q1@x.test", "Q", "X", 1)
	ctx := context.Background()
	cable := f.item("Câble", false)
	f.lend("owner@x.test", cable, t0.Add(-time.Hour))
	if _, err := f.svc.Join(ctx, "q1@x.test", cable.ID, nil); err != nil {
		t.Fatal(err)
	}
	f.exec(`UPDATE reservation_element SET status = 'maintenance' WHERE id_reservation_element = $1`, cable.ID)
	f.exec(`UPDATE reservation SET end_date = $1`, t0)

	if next, err := f.svc.ReleaseOf(ctx, cable.ID, nil); err != nil || next {
		t.Fatalf("a release of an item in maintenance notified somebody: %v, %v", next, err)
	}
	if _, err := f.svc.Reserve(ctx, "q1@x.test", cable.ID, nil); !errors.Is(err, ErrNotBookable) {
		t.Fatalf("reserve: %v", err)
	}
	f.exec(`UPDATE reservation SET end_date = NULL`)
	f.user("q2@x.test", "Q", "Y", 1)
	if _, err := f.svc.Join(ctx, "q2@x.test", cable.ID, nil); !errors.Is(err, ErrNotBookable) {
		t.Fatalf("join: %v", err)
	}
}

func TestTheSummary(t *testing.T) {
	f := newFx(t)
	f.user("admin@x.test", "Ada", "Lovelace", 1, "ADMIN")
	f.user("bob@x.test", "Bob", "Martin", 1)
	f.user("amy@x.test", "Amy", "Roy", 1)
	ctx := context.Background()
	room := f.item("Salle", true)
	cable := f.item("Câble", false)
	hdmi := f.item("HDMI", false)
	broken := f.item("Console", false)
	f.exec(`UPDATE reservation_element SET status = 'inactive' WHERE id_reservation_element = $1`, broken.ID)
	f.lend("bob@x.test", cable, t0.Add(-8*24*time.Hour))
	f.lend("bob@x.test", hdmi, t0.Add(-2*24*time.Hour))
	f.book("bob@x.test", room, t0)
	f.book("amy@x.test", room, t0.Add(2*time.Hour))
	f.book("amy@x.test", room, t0.Add(4*time.Hour))
	cancelled := t0.Add(4 * time.Hour)
	f.svc.Cancel(ctx, CancelRequest{ItemID: room.ID, Start: &cancelled, ActorEmail: "admin@x.test", AsAdmin: true})
	f.svc.Join(ctx, "amy@x.test", cable.ID, nil)
	f.exec(`UPDATE reservation_waitlist SET notified_at = now()`)

	st, err := f.svc.Stats(ctx)
	if err != nil {
		t.Fatal(err)
	}
	want := ReservationStats{ActiveLoans: 2, OngoingSlots: 1, UpcomingSlots: 1, OverdueDay1: 2, OverdueDay5: 1, OverdueDay7: 1,
		WaitlistEntries: 1, AwaitingResponse: 1, Cancelled30Days: 1, Items: 4, UnavailableItems: 1}
	if st != want {
		t.Fatalf("stats = %+v\nwant   %+v", st, want)
	}
}
