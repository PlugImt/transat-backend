package reservations

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"html/template"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	goi18n "github.com/nicksnyder/go-i18n/v2/i18n"
	appI18n "github.com/plugimt/transat-backend/i18n"
	"github.com/plugimt/transat-backend/internal/database/dbtest"
	"github.com/plugimt/transat-backend/models"
	"github.com/plugimt/transat-backend/services"
	"github.com/plugimt/transat-backend/services/push"
)

var (
	i18nOnce sync.Once
	i18nErr  error
)

func initI18n(t *testing.T) {
	t.Helper()
	i18nOnce.Do(func() { i18nErr = appI18n.Init() })
	if i18nErr != nil {
		t.Fatal(i18nErr)
	}
}

// fakeNotifier records what would be queued in the delivery pipeline.
type fakeNotifier struct {
	mu          sync.Mutex
	unreachable map[string]bool
	fail        error
	sent        []services.Announcement
}

func (n *fakeNotifier) Announce(_ context.Context, a services.Announcement) (push.EnqueueResult, error) {
	n.mu.Lock()
	defer n.mu.Unlock()
	if n.fail != nil {
		return push.EnqueueResult{}, n.fail
	}
	n.sent = append(n.sent, a)
	return push.EnqueueResult{ID: int64(len(n.sent))}, nil
}

func (n *fakeNotifier) CanReach(_ context.Context, email string, _ models.NotificationCategory) (bool, error) {
	n.mu.Lock()
	defer n.mu.Unlock()
	return !n.unreachable[email], nil
}

func (n *fakeNotifier) all() []services.Announcement {
	n.mu.Lock()
	defer n.mu.Unlock()
	return append([]services.Announcement(nil), n.sent...)
}

func (n *fakeNotifier) reset() {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.sent = nil
}

// to returns the single user an announcement is for.
func to(t *testing.T, a services.Announcement) string {
	t.Helper()
	if a.Audience.Type != models.AudienceUsers || len(a.Audience.Emails) != 1 {
		t.Fatalf("audience = %+v, want exactly one user", a.Audience)
	}
	return a.Audience.Emails[0]
}

func text(a services.Announcement, lang string) (string, string) {
	return a.Compose(lang, appI18n.GetLocalizer(lang))
}

type sentMail struct {
	mail models.Email
	data CancelEmailData
}

type fakeMailer struct {
	mu   sync.Mutex
	sent []sentMail
}

func (m *fakeMailer) SendEmail(mail models.Email, data interface{}) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.sent = append(m.sent, sentMail{mail, data.(CancelEmailData)})
	return nil
}

type fx struct {
	t      *testing.T
	db     *sql.DB
	svc    *Service
	ntf    *fakeNotifier
	mail   *fakeMailer
	now    time.Time
	club   int
	nextID int
}

// t0 is the fixed "now" of the tests.
var t0 = time.Date(2026, 7, 1, 12, 0, 0, 0, time.UTC)

func newFx(t *testing.T) *fx {
	initI18n(t)
	db := dbtest.Open(t)
	f := &fx{t: t, db: db, ntf: &fakeNotifier{unreachable: map[string]bool{}}, mail: &fakeMailer{}, now: t0}
	f.svc = New(db, f.ntf, f.mail)
	f.svc.SetBackground(func(fn func()) { fn() })
	f.svc.SetClock(func() time.Time { return f.now })
	f.exec(`INSERT INTO roles (name) VALUES ('ADMIN') ON CONFLICT (name) DO NOTHING`)
	if err := db.QueryRow(`INSERT INTO clubs (name, picture) VALUES ('Club', 'p') RETURNING id_clubs`).Scan(&f.club); err != nil {
		t.Fatal(err)
	}
	return f
}

func (f *fx) exec(q string, args ...any) {
	f.t.Helper()
	if _, err := f.db.Exec(q, args...); err != nil {
		f.t.Fatalf("%s: %v", q, err)
	}
}

func (f *fx) count(q string, args ...any) int {
	f.t.Helper()
	var n int
	if err := f.db.QueryRow(q, args...).Scan(&n); err != nil {
		f.t.Fatal(err)
	}
	return n
}

// user creates an account (language 1 = fr, 2 = en).
func (f *fx) user(email, first, last string, lang int, roles ...string) {
	f.t.Helper()
	f.exec(`INSERT INTO newf (email, password, first_name, last_name, language) VALUES ($1, 'x', $2, $3, $4)`, email, first, last, lang)
	for _, r := range roles {
		f.exec(`INSERT INTO newf_roles (email, id_roles) SELECT $1, id_roles FROM roles WHERE name = $2`, email, r)
	}
}

func (f *fx) item(name string, slot bool) Item {
	f.t.Helper()
	it := Item{Name: name, Slot: slot}
	if err := f.db.QueryRow(`INSERT INTO reservation_element (name, slot, id_clubs) VALUES ($1, $2, $3) RETURNING id_reservation_element`,
		name, slot, f.club).Scan(&it.ID); err != nil {
		f.t.Fatal(err)
	}
	return it
}

// book inserts a slot reservation made a day before t0.
func (f *fx) book(email string, it Item, start time.Time) {
	f.t.Helper()
	f.exec(`INSERT INTO reservation (email, id_reservation_element, start_date, end_date, created_at) VALUES ($1, $2, $3, $4, $5)`,
		email, it.ID, start, start.Add(time.Hour), t0.Add(-24*time.Hour))
}

// lend inserts an open loan that started at the given time.
func (f *fx) lend(email string, it Item, start time.Time) {
	f.t.Helper()
	f.exec(`INSERT INTO reservation (email, id_reservation_element, start_date, created_at) VALUES ($1, $2, $3, $4)`, email, it.ID, start, start)
}

func (f *fx) holder(it Item, start *time.Time) string {
	f.t.Helper()
	var email sql.NullString
	var err error
	if it.Slot {
		err = f.db.QueryRow(`SELECT email FROM reservation WHERE id_reservation_element = $1 AND start_date = $2`, it.ID, start.UTC()).Scan(&email)
	} else {
		err = f.db.QueryRow(`SELECT email FROM reservation WHERE id_reservation_element = $1 AND end_date IS NULL`, it.ID).Scan(&email)
	}
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		f.t.Fatal(err)
	}
	return email.String
}

func (f *fx) at(d time.Duration) *time.Time {
	t := t0.Add(d)
	return &t
}

// ---------------------------------------------------------------------------------------------------------------
// Admin cancellation

func TestAdminCancelsAThirdPartySlotAndTheUserIsToldWithTheAdminsName(t *testing.T) {
	f := newFx(t)
	f.user("admin@x.test", "Ada", "Lovelace", 1, "ADMIN")
	f.user("bob@x.test", "Bob", "Martin", 2)
	room := f.item("Salle <b>Ampère</b>", true)
	start := t0.Add(26 * time.Hour)
	f.book("bob@x.test", room, start)

	res, err := f.svc.Cancel(context.Background(), CancelRequest{ItemID: room.ID, Start: &start, ActorEmail: "admin@x.test", AsAdmin: true})
	if err != nil {
		t.Fatal(err)
	}
	if !res.ByAdmin || res.Owner != "bob@x.test" {
		t.Fatalf("result = %+v", res)
	}
	if f.holder(room, &start) != "" {
		t.Fatal("the reservation is still there")
	}

	// Push notification to the owner, and only to them, in the category the user can switch off.
	sent := f.ntf.all()
	if len(sent) != 1 {
		t.Fatalf("%d notifications", len(sent))
	}
	a := sent[0]
	if to(t, a) != "bob@x.test" || a.Category != models.NotificationReservations || a.Kind != "reservation_cancelled" {
		t.Fatalf("announcement = %+v", a)
	}
	if a.Navigation == nil || a.Navigation.Type != models.NavigationReservation || a.Navigation.ID != fmt.Sprint(room.ID) {
		t.Fatalf("navigation = %+v", a.Navigation)
	}
	title, body := text(a, "en")
	if title != "Reservation cancelled" || !strings.Contains(body, "Ada Lovelace") || !strings.Contains(body, room.Name) || !strings.Contains(body, "Jul 2026") {
		t.Fatalf("en: %q / %q", title, body)
	}
	_, body = text(a, "fr")
	if !strings.Contains(body, "Ada Lovelace") || !strings.Contains(body, "a annulé votre réservation") || !strings.Contains(body, "02/07/2026 16:00") {
		t.Fatalf("fr: %q", body)
	}

	// Email to the same user, in their language, with the same information.
	if len(f.mail.sent) != 1 {
		t.Fatalf("%d emails", len(f.mail.sent))
	}
	m := f.mail.sent[0]
	if m.mail.Recipient != "bob@x.test" || m.mail.Language != "en" || m.mail.SubjectKey != emailSubjectKey ||
		m.mail.Template != "email_templates/email_template_reservation_cancelled.html" {
		t.Fatalf("mail = %+v", m.mail)
	}
	if m.data.AdminName != "Ada Lovelace" || m.data.ItemName != room.Name || !m.data.HasSlot || m.data.Slot != "2 Jul 2026, 16:00" {
		t.Fatalf("data = %+v", m.data)
	}
}

func TestAnAdminCancellingTheirOwnReservationNotifiesNobody(t *testing.T) {
	f := newFx(t)
	f.user("admin@x.test", "Ada", "Lovelace", 1, "ADMIN")
	room := f.item("Salle", true)
	start := t0.Add(26 * time.Hour)
	f.book("admin@x.test", room, start)

	res, err := f.svc.Cancel(context.Background(), CancelRequest{ItemID: room.ID, Start: &start, ActorEmail: "admin@x.test", AsAdmin: true})
	if err != nil || res.ByAdmin {
		t.Fatalf("result = %+v, %v", res, err)
	}
	if len(f.ntf.all()) != 0 || len(f.mail.sent) != 0 {
		t.Fatal("somebody was notified of their own cancellation")
	}
}

func TestAUserCancelsTheirOwnSlotButNotSomebodyElses(t *testing.T) {
	f := newFx(t)
	f.user("bob@x.test", "Bob", "Martin", 1)
	f.user("carol@x.test", "Carol", "Dupont", 1)
	room := f.item("Salle", true)
	start := t0.Add(26 * time.Hour)
	f.book("bob@x.test", room, start)

	_, err := f.svc.Cancel(context.Background(), CancelRequest{ItemID: room.ID, Start: &start, ActorEmail: "carol@x.test"})
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("a user cancelled somebody else's reservation: %v", err)
	}
	if f.holder(room, &start) != "bob@x.test" {
		t.Fatal("the reservation was removed")
	}
	if len(f.ntf.all()) != 0 || len(f.mail.sent) != 0 {
		t.Fatal("a refused cancellation notified somebody")
	}

	res, err := f.svc.Cancel(context.Background(), CancelRequest{ItemID: room.ID, Start: &start, ActorEmail: "bob@x.test"})
	if err != nil || res.ByAdmin {
		t.Fatalf("owner: %+v, %v", res, err)
	}
	if f.holder(room, &start) != "" {
		t.Fatal("the owner could not cancel their own reservation")
	}
	if len(f.mail.sent) != 0 {
		t.Fatal("the owner got an admin cancellation email for their own action")
	}
}

func TestAnAdministratorOutsideTheAdminRoutesIsAnOrdinaryUser(t *testing.T) {
	f := newFx(t)
	f.user("admin@x.test", "Ada", "Lovelace", 1, "ADMIN")
	f.user("bob@x.test", "Bob", "Martin", 1)
	room := f.item("Salle", true)
	cable := f.item("Câble HDMI", false)
	start := t0.Add(26 * time.Hour)
	f.book("bob@x.test", room, start)
	f.lend("bob@x.test", cable, t0.Add(-48*time.Hour))
	ctx := context.Background()

	if _, err := f.svc.Cancel(ctx, CancelRequest{ItemID: room.ID, Start: &start, ActorEmail: "admin@x.test"}); !errors.Is(err, ErrNotFound) {
		t.Fatalf("an administrator cancelled a slot through the user route: %v", err)
	}
	if _, err := f.svc.Cancel(ctx, CancelRequest{ItemID: cable.ID, ActorEmail: "admin@x.test"}); !errors.Is(err, ErrNotSlotItem) {
		t.Fatalf("an administrator ended a loan through the user route: %v", err)
	}
	if f.holder(room, &start) != "bob@x.test" || f.holder(cable, nil) != "bob@x.test" {
		t.Fatal("a reservation was removed")
	}
	if len(f.ntf.all()) != 0 || len(f.mail.sent) != 0 {
		t.Fatal("a refused cancellation notified somebody")
	}

	// Their own reservation, they can cancel, and it is an ordinary cancellation.
	f.book("admin@x.test", f.item("Salle B", true), start)
	own := f.count(`SELECT id_reservation_element FROM reservation WHERE email = 'admin@x.test'`)
	res, err := f.svc.Cancel(ctx, CancelRequest{ItemID: own, Start: &start, ActorEmail: "admin@x.test"})
	if err != nil || res.ByAdmin {
		t.Fatalf("own cancellation: %+v, %v", res, err)
	}

	// The admin flag does not replace the role.
	if _, err := f.svc.Cancel(ctx, CancelRequest{ItemID: room.ID, Start: &start, ActorEmail: "bob@x.test", AsAdmin: true}); !errors.Is(err, ErrNotAdmin) {
		t.Fatalf("a user without the role cancelled as an administrator: %v", err)
	}
	if f.holder(room, &start) != "bob@x.test" {
		t.Fatal("the reservation was removed")
	}
}

func TestOnlyAnAdminCanEndSomebodysLoanOfAnItemWithoutSlot(t *testing.T) {
	f := newFx(t)
	f.user("admin@x.test", "Ada", "Lovelace", 1, "ADMIN")
	f.user("bob@x.test", "Bob", "Martin", 1)
	f.user("carol@x.test", "Carol", "Dupont", 1)
	cable := f.item("Câble HDMI", false)
	f.lend("bob@x.test", cable, t0.Add(-48*time.Hour))

	if _, err := f.svc.Cancel(context.Background(), CancelRequest{ItemID: cable.ID, ActorEmail: "carol@x.test"}); !errors.Is(err, ErrNotSlotItem) {
		t.Fatalf("a user ended a loan: %v", err)
	}
	if f.holder(cable, nil) != "bob@x.test" {
		t.Fatal("the loan was ended")
	}

	res, err := f.svc.Cancel(context.Background(), CancelRequest{ItemID: cable.ID, ActorEmail: "admin@x.test", AsAdmin: true})
	if err != nil || !res.ByAdmin || res.Owner != "bob@x.test" {
		t.Fatalf("result = %+v, %v", res, err)
	}
	if f.holder(cable, nil) != "" {
		t.Fatal("the loan is still open")
	}
	if n := f.count(`SELECT count(*) FROM reservation WHERE id_reservation_element = $1 AND end_date IS NOT NULL`, cable.ID); n != 1 {
		t.Fatal("the loan should be ended, not deleted")
	}
	sent := f.ntf.all()
	if len(sent) != 1 || to(t, sent[0]) != "bob@x.test" {
		t.Fatalf("notifications = %+v", sent)
	}
	_, body := text(sent[0], "fr")
	if body != "Ada Lovelace a annulé votre réservation de Câble HDMI." {
		t.Fatalf("body = %q", body)
	}
	if len(f.mail.sent) != 1 || f.mail.sent[0].data.HasSlot {
		t.Fatalf("emails = %+v", f.mail.sent)
	}

	// Nothing left to end.
	if _, err := f.svc.Cancel(context.Background(), CancelRequest{ItemID: cable.ID, ActorEmail: "admin@x.test", AsAdmin: true}); !errors.Is(err, ErrNotFound) {
		t.Fatalf("second cancel: %v", err)
	}
}

func TestACancellationNeverFailsBecauseTheNotificationFailed(t *testing.T) {
	f := newFx(t)
	f.user("admin@x.test", "Ada", "Lovelace", 1, "ADMIN")
	f.user("bob@x.test", "Bob", "Martin", 1)
	room := f.item("Salle", true)
	start := t0.Add(26 * time.Hour)
	f.book("bob@x.test", room, start)
	f.ntf.fail = errors.New("queue down")

	if _, err := f.svc.Cancel(context.Background(), CancelRequest{ItemID: room.ID, Start: &start, ActorEmail: "admin@x.test", AsAdmin: true}); err != nil {
		t.Fatal(err)
	}
	if f.holder(room, &start) != "" {
		t.Fatal("the reservation was not cancelled")
	}
	if len(f.mail.sent) != 1 {
		t.Fatal("the email is independent of the push")
	}
}

// ---------------------------------------------------------------------------------------------------------------
// Reminders

func TestASlotReservationIsRemindedOneHourBeforeItStartsOnce(t *testing.T) {
	f := newFx(t)
	f.user("bob@x.test", "Bob", "Martin", 1)
	f.user("amy@x.test", "Amy", "Roy", 1)
	room := f.item("Salle", true)
	soon := t0.Add(time.Hour)      // the reminder is due now
	later := t0.Add(3 * time.Hour) // not due yet
	f.book("bob@x.test", room, soon)
	f.book("amy@x.test", room, later)

	sent, err := f.svc.SendDueReminders(context.Background())
	if err != nil || sent != 1 {
		t.Fatalf("sent %d, %v", sent, err)
	}
	a := f.ntf.all()
	if len(a) != 1 || to(t, a[0]) != "bob@x.test" || a[0].Kind != "reservation_reminder" || a[0].Category != models.NotificationReservations {
		t.Fatalf("announcements = %+v", a)
	}
	title, body := text(a[0], "fr")
	if title != "Réservation dans 1h" || body != "Salle · 15:00" {
		t.Fatalf("text = %q / %q", title, body)
	}
	if !a[0].ExpiresAt.Equal(soon) {
		t.Fatalf("expires at %v, want the start of the slot", a[0].ExpiresAt)
	}
	if a[0].Navigation.Params["date"] != "2026-07-01" {
		t.Fatalf("navigation = %+v", a[0].Navigation)
	}

	// Running again, or from another instance, sends nothing more.
	for i := 0; i < 3; i++ {
		if sent, err := f.svc.SendDueReminders(context.Background()); err != nil || sent != 0 {
			t.Fatalf("run %d sent %d, %v", i, sent, err)
		}
	}

	// An hour later the second slot is due, and only that one.
	f.now = t0.Add(2 * time.Hour)
	if sent, _ := f.svc.SendDueReminders(context.Background()); sent != 1 {
		t.Fatalf("second slot: sent %d", sent)
	}
	if got := f.ntf.all(); len(got) != 2 || to(t, got[1]) != "amy@x.test" {
		t.Fatalf("announcements = %+v", got)
	}
}

func TestNoSlotReminderForACancelledOrLateBookedOrPastReservation(t *testing.T) {
	f := newFx(t)
	f.user("bob@x.test", "Bob", "Martin", 1)
	f.user("amy@x.test", "Amy", "Roy", 1)
	f.user("eve@x.test", "Eve", "Roy", 1)
	f.user("tom@x.test", "Tom", "Roy", 1)
	room := f.item("Salle", true)

	cancelled := t0.Add(time.Hour)
	f.book("bob@x.test", room, cancelled)
	if _, err := f.svc.Cancel(context.Background(), CancelRequest{ItemID: room.ID, Start: &cancelled, ActorEmail: "bob@x.test"}); err != nil {
		t.Fatal(err)
	}
	// Booked 5 minutes ago for 50 minutes from now: the reminder time had already passed when it was booked.
	f.exec(`INSERT INTO reservation (email, id_reservation_element, start_date, end_date, created_at) VALUES ($1, $2, $3, $4, $5)`,
		"amy@x.test", room.ID, t0.Add(50*time.Minute), t0.Add(110*time.Minute), t0.Add(-5*time.Minute))
	// Already over, and about to start in 10 minutes (too late to be worth a reminder).
	f.book("eve@x.test", room, t0.Add(-3*time.Hour))
	f.book("tom@x.test", room, t0.Add(10*time.Minute))

	if sent, err := f.svc.SendDueReminders(context.Background()); err != nil || sent != 0 {
		t.Fatalf("sent %d, %v; announcements %+v", sent, err, f.ntf.all())
	}
}

func TestLoanRemindersAfterOneFiveAndSevenDaysEachOnce(t *testing.T) {
	f := newFx(t)
	f.user("bob@x.test", "Bob", "Martin", 2)
	cable := f.item("HDMI cable", false)
	start := t0
	f.lend("bob@x.test", cable, start)

	type step struct {
		after time.Duration
		want  int
		title string
	}
	steps := []step{
		{23 * time.Hour, 0, ""},
		{25 * time.Hour, 1, "Remember to return the item"},
		{25*time.Hour + 10*time.Minute, 0, ""}, // same milestone: never twice
		{4 * 24 * time.Hour, 0, ""},
		{5*24*time.Hour + time.Hour, 1, "Item borrowed for 5 days"},
		{6 * 24 * time.Hour, 0, ""},
		{7*24*time.Hour + time.Hour, 1, "Item borrowed for 7 days"},
		{7*24*time.Hour + 2*time.Hour, 0, ""},
		{20 * 24 * time.Hour, 0, ""}, // nothing after the seventh day
	}
	total := 0
	for _, s := range steps {
		f.now = start.Add(s.after)
		sent, err := f.svc.SendDueReminders(context.Background())
		if err != nil || sent != s.want {
			t.Fatalf("after %v: sent %d, %v; want %d", s.after, sent, err, s.want)
		}
		if s.want == 1 {
			total++
			a := f.ntf.all()
			last := a[len(a)-1]
			title, body := text(last, "en")
			if to(t, last) != "bob@x.test" || title != s.title || !strings.Contains(body, "HDMI cable") {
				t.Fatalf("after %v: %q / %q", s.after, title, body)
			}
			if last.Navigation.Type != models.NavigationReservation || last.Navigation.ID != fmt.Sprint(cable.ID) {
				t.Fatalf("navigation = %+v", last.Navigation)
			}
		}
	}
	if total != 3 || len(f.ntf.all()) != 3 {
		t.Fatalf("%d reminders in total", len(f.ntf.all()))
	}
	if n := f.count(`SELECT count(*) FROM reservation_reminders`); n != 3 {
		t.Fatalf("%d claims recorded", n)
	}
}

func TestNoLoanReminderOnceTheItemIsReturnedOrCancelled(t *testing.T) {
	f := newFx(t)
	f.user("admin@x.test", "Ada", "Lovelace", 1, "ADMIN")
	f.user("bob@x.test", "Bob", "Martin", 1)
	f.user("amy@x.test", "Amy", "Roy", 1)
	returned := f.item("Rendu", false)
	cancelled := f.item("Annulé", false)
	f.lend("bob@x.test", returned, t0.Add(-25*time.Hour))
	f.lend("amy@x.test", cancelled, t0.Add(-25*time.Hour))
	f.exec(`UPDATE reservation SET end_date = $1 WHERE id_reservation_element = $2`, t0.Add(-time.Hour), returned.ID)
	if _, err := f.svc.Cancel(context.Background(), CancelRequest{ItemID: cancelled.ID, ActorEmail: "admin@x.test", AsAdmin: true}); err != nil {
		t.Fatal(err)
	}
	f.ntf.reset()

	if sent, err := f.svc.SendDueReminders(context.Background()); err != nil || sent != 0 {
		t.Fatalf("sent %d, %v", sent, err)
	}
	if len(f.ntf.all()) != 0 {
		t.Fatal("a reminder went out for something that was returned")
	}
}

func TestLoansForgottenLongAgoAreNotRemindedAfterADeployment(t *testing.T) {
	f := newFx(t)
	f.user("bob@x.test", "Bob", "Martin", 1)
	cable := f.item("Câble", false)
	f.lend("bob@x.test", cable, t0.Add(-90*24*time.Hour))
	if sent, _ := f.svc.SendDueReminders(context.Background()); sent != 0 {
		t.Fatalf("sent %d reminders for a loan 90 days old", sent)
	}
}

func TestAReminderThatCouldNotBeQueuedIsTriedAgainAndNotDuplicated(t *testing.T) {
	f := newFx(t)
	f.user("bob@x.test", "Bob", "Martin", 1)
	cable := f.item("Câble", false)
	f.lend("bob@x.test", cable, t0.Add(-25*time.Hour))

	f.ntf.fail = errors.New("database busy")
	if sent, err := f.svc.SendDueReminders(context.Background()); err != nil || sent != 0 {
		t.Fatalf("sent %d, %v", sent, err)
	}
	if n := f.count(`SELECT count(*) FROM reservation_reminders`); n != 0 {
		t.Fatal("a reminder that was not queued stays claimed")
	}
	f.ntf.fail = nil
	if sent, _ := f.svc.SendDueReminders(context.Background()); sent != 1 {
		t.Fatalf("retry sent %d", sent)
	}
	if sent, _ := f.svc.SendDueReminders(context.Background()); sent != 0 {
		t.Fatalf("a third run sent %d", sent)
	}
}

func TestConcurrentSchedulersSendEachReminderOnce(t *testing.T) {
	f := newFx(t)
	f.user("bob@x.test", "Bob", "Martin", 1)
	cable := f.item("Câble", false)
	f.lend("bob@x.test", cable, t0.Add(-25*time.Hour))
	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if _, err := f.svc.SendDueReminders(context.Background()); err != nil {
				t.Error(err)
			}
		}()
	}
	wg.Wait()
	if n := len(f.ntf.all()); n != 1 {
		t.Fatalf("%d reminders sent by 8 concurrent runs", n)
	}
}

func TestCancellingAndBookingAgainResetsTheSlotReminder(t *testing.T) {
	f := newFx(t)
	f.user("bob@x.test", "Bob", "Martin", 1)
	room := f.item("Salle", true)
	start := t0.Add(time.Hour)
	f.book("bob@x.test", room, start)
	if sent, _ := f.svc.SendDueReminders(context.Background()); sent != 1 {
		t.Fatalf("sent %d", sent)
	}
	if _, err := f.svc.Cancel(context.Background(), CancelRequest{ItemID: room.ID, Start: &start, ActorEmail: "bob@x.test"}); err != nil {
		t.Fatal(err)
	}
	if n := f.count(`SELECT count(*) FROM reservation_reminders`); n != 0 {
		t.Fatal("the claim of a cancelled reservation was kept")
	}
}

func TestTickPurgesTheWaitingListOfSlotsThatAreOver(t *testing.T) {
	f := newFx(t)
	f.user("bob@x.test", "Bob", "Martin", 1)
	room := f.item("Salle", true)
	cable := f.item("Câble", false)
	f.exec(`INSERT INTO reservation_waitlist (email, id_reservation_element, slot_start) VALUES ('bob@x.test', $1, $2)`, room.ID, t0.Add(-time.Hour))
	f.exec(`INSERT INTO reservation_waitlist (email, id_reservation_element, slot_start) VALUES ('bob@x.test', $1, $2)`, room.ID, t0.Add(time.Hour))
	f.exec(`INSERT INTO reservation_waitlist (email, id_reservation_element) VALUES ('bob@x.test', $1)`, cable.ID)
	if _, err := f.svc.Tick(context.Background()); err != nil {
		t.Fatal(err)
	}
	if n := f.count(`SELECT count(*) FROM reservation_waitlist`); n != 2 {
		t.Fatalf("%d entries left, want the future slot and the item", n)
	}
}

// ---------------------------------------------------------------------------------------------------------------
// Waiting list

func TestJoiningTheWaitingList(t *testing.T) {
	f := newFx(t)
	for _, u := range []string{"owner", "q1", "q2"} {
		f.user(u+"@x.test", u, "X", 1)
	}
	room := f.item("Salle", true)
	cable := f.item("Câble", false)
	start := t0.Add(26 * time.Hour)
	f.book("owner@x.test", room, start)
	f.lend("owner@x.test", cable, t0.Add(-time.Hour))
	ctx := context.Background()

	// Only what is already reserved can be waited for.
	free := t0.Add(30 * time.Hour)
	if _, err := f.svc.Join(ctx, "q1@x.test", room.ID, &free); !errors.Is(err, ErrStillAvailable) {
		t.Fatalf("a free slot: %v", err)
	}
	if _, err := f.svc.Join(ctx, "owner@x.test", room.ID, &start); !errors.Is(err, ErrOwnReservation) {
		t.Fatalf("own slot: %v", err)
	}
	if _, err := f.svc.Join(ctx, "q1@x.test", room.ID, nil); !errors.Is(err, ErrStartRequired) {
		t.Fatalf("no slot: %v", err)
	}
	past := t0.Add(-time.Hour)
	if _, err := f.svc.Join(ctx, "q1@x.test", room.ID, &past); !errors.Is(err, ErrInvalidStart) {
		t.Fatalf("past slot: %v", err)
	}
	notHour := start.Add(30 * time.Minute)
	if _, err := f.svc.Join(ctx, "q1@x.test", room.ID, &notHour); !errors.Is(err, ErrInvalidStart) {
		t.Fatalf("half hour: %v", err)
	}
	if _, err := f.svc.Join(ctx, "q1@x.test", 99999, nil); !errors.Is(err, ErrItemNotFound) {
		t.Fatalf("unknown item: %v", err)
	}

	// Queue of a slot, first come first served.
	for want, who := range []string{"q1", "q2"} {
		pos, err := f.svc.Join(ctx, who+"@x.test", room.ID, &start)
		if err != nil || pos != want+1 {
			t.Fatalf("%s: position %d, %v", who, pos, err)
		}
	}
	if _, err := f.svc.Join(ctx, "q1@x.test", room.ID, &start); !errors.Is(err, ErrAlreadyQueued) {
		t.Fatalf("twice: %v", err)
	}
	// Queue of an item on loan; the slot argument is irrelevant.
	if pos, err := f.svc.Join(ctx, "q2@x.test", cable.ID, nil); err != nil || pos != 1 {
		t.Fatalf("loan queue: %d, %v", pos, err)
	}

	st, err := f.svc.Status(ctx, "q2@x.test", room.ID, &start)
	if err != nil || st.Count != 2 || !st.Joined || st.Position != 2 || st.Notified {
		t.Fatalf("status = %+v, %v", st, err)
	}
	if st, _ := f.svc.Status(ctx, "owner@x.test", room.ID, &start); st.Joined || st.Count != 2 {
		t.Fatalf("status of a stranger = %+v", st)
	}
	entries, err := f.svc.List(ctx, "q2@x.test")
	if err != nil || len(entries) != 2 || entries[0].ItemName != "Salle" || entries[0].Position != 2 || entries[1].ItemName != "Câble" || entries[1].Position != 1 {
		t.Fatalf("entries = %+v, %v", entries, err)
	}
	if entries[0].StartDate == nil || *entries[0].StartDate != start.Format(APITimeLayout) || entries[1].StartDate != nil {
		t.Fatalf("start dates = %v / %v", entries[0].StartDate, entries[1].StartDate)
	}
}

func TestOnlyTheFirstInTheQueueIsNotifiedWhenASlotIsReleased(t *testing.T) {
	f := newFx(t)
	for _, u := range []string{"owner", "q1", "q2", "q3"} {
		f.user(u+"@x.test", u, "X", 1)
	}
	room := f.item("Salle", true)
	start := t0.Add(26 * time.Hour)
	f.book("owner@x.test", room, start)
	ctx := context.Background()
	for _, u := range []string{"q1", "q2", "q3"} {
		if _, err := f.svc.Join(ctx, u+"@x.test", room.ID, &start); err != nil {
			t.Fatal(err)
		}
	}
	// Joining is not an event: nobody was told anything.
	if len(f.ntf.all()) != 0 {
		t.Fatal("joining notified somebody")
	}

	res, err := f.svc.Cancel(ctx, CancelRequest{ItemID: room.ID, Start: &start, ActorEmail: "owner@x.test"})
	if err != nil || !res.NextNotified {
		t.Fatalf("result = %+v, %v", res, err)
	}
	sent := f.ntf.all()
	if len(sent) != 1 {
		t.Fatalf("%d notifications, want only the first of the queue", len(sent))
	}
	a := sent[0]
	if to(t, a) != "q1@x.test" || a.Category != models.NotificationReservations || a.Kind != "reservation_waitlist" {
		t.Fatalf("announcement = %+v", a)
	}
	if a.ActionCategory != WaitlistActionCategory {
		t.Fatalf("action category = %q", a.ActionCategory)
	}
	nav := a.Navigation
	if nav.Type != models.NavigationReservation || nav.ID != fmt.Sprint(room.ID) || nav.Params["waitlist"] != "1" ||
		nav.Params["start"] != start.Format(APITimeLayout) || nav.Params["date"] != "2026-07-02" {
		t.Fatalf("navigation = %+v", nav)
	}
	if err := nav.Validate(); err != nil {
		t.Fatal(err)
	}
	title, body := text(a, "fr")
	if title != "Salle est disponible" || !strings.Contains(body, "02/07/2026 16:00") {
		t.Fatalf("fr: %q / %q", title, body)
	}
	if st, _ := f.svc.Status(ctx, "q1@x.test", room.ID, &start); !st.Notified {
		t.Fatal("the first of the queue should hold the baton")
	}
	if st, _ := f.svc.Status(ctx, "q2@x.test", room.ID, &start); st.Notified {
		t.Fatal("the second of the queue was notified")
	}
}

func TestReturningAnItemNotifiesTheFirstInTheQueue(t *testing.T) {
	f := newFx(t)
	for _, u := range []string{"owner", "q1", "q2"} {
		f.user(u+"@x.test", u, "X", 2)
	}
	cable := f.item("HDMI cable", false)
	f.lend("owner@x.test", cable, t0.Add(-time.Hour))
	ctx := context.Background()
	f.svc.Join(ctx, "q1@x.test", cable.ID, nil)
	f.svc.Join(ctx, "q2@x.test", cable.ID, nil)

	// The owner returns the item through the existing endpoint, which ends the loan and calls ReleaseOf.
	f.exec(`UPDATE reservation SET end_date = $1 WHERE id_reservation_element = $2`, t0, cable.ID)
	notified, err := f.svc.ReleaseOf(ctx, cable.ID, nil)
	if err != nil || !notified {
		t.Fatalf("notified = %v, %v", notified, err)
	}
	sent := f.ntf.all()
	if len(sent) != 1 || to(t, sent[0]) != "q1@x.test" {
		t.Fatalf("announcements = %+v", sent)
	}
	title, body := text(sent[0], "en")
	if title != "HDMI cable is available" || !strings.Contains(body, "waiting list") {
		t.Fatalf("en: %q / %q", title, body)
	}
	if _, has := sent[0].Navigation.Params["start"]; has {
		t.Fatal("an item without slots has no start")
	}
}

func TestADirectBookingIsNeverBlockedAndTheNotifiedUserGetsAClearError(t *testing.T) {
	f := newFx(t)
	for _, u := range []string{"owner", "q1", "q2", "fast"} {
		f.user(u+"@x.test", u, "X", 1)
	}
	room := f.item("Salle", true)
	start := t0.Add(26 * time.Hour)
	f.book("owner@x.test", room, start)
	ctx := context.Background()
	f.svc.Join(ctx, "q1@x.test", room.ID, &start)
	f.svc.Join(ctx, "q2@x.test", room.ID, &start)
	if _, err := f.svc.Cancel(ctx, CancelRequest{ItemID: room.ID, Start: &start, ActorEmail: "owner@x.test"}); err != nil {
		t.Fatal(err)
	}

	// Somebody who is not in the queue books the slot first: nothing in the queue prevents it.
	f.exec(`INSERT INTO reservation (email, id_reservation_element, start_date, end_date) VALUES ('fast@x.test', $1, $2, $3)`,
		room.ID, start, start.Add(time.Hour))

	// "Réserver" is now refused, explicitly, and the user keeps their place.
	_, err := f.svc.Reserve(ctx, "q1@x.test", room.ID, &start)
	if !errors.Is(err, ErrNotAvailable) {
		t.Fatalf("error = %v", err)
	}
	if f.holder(room, &start) != "fast@x.test" {
		t.Fatal("the slot changed hands")
	}
	if st, _ := f.svc.Status(ctx, "q1@x.test", room.ID, &start); !st.Joined {
		t.Fatal("the user lost their place after a failed attempt")
	}

	// "Annuler" while the slot is taken: no one else is notified for a slot that is not free.
	f.ntf.reset()
	next, err := f.svc.Remove(ctx, "q1@x.test", room.ID, &start)
	if err != nil || next || len(f.ntf.all()) != 0 {
		t.Fatalf("next = %v, %v, %d notifications", next, err, len(f.ntf.all()))
	}
}

func TestReserveBooksTheSlotAndLeavesTheQueue(t *testing.T) {
	f := newFx(t)
	for _, u := range []string{"owner", "q1", "q2"} {
		f.user(u+"@x.test", u, "X", 1)
	}
	room := f.item("Salle", true)
	start := t0.Add(26 * time.Hour)
	f.book("owner@x.test", room, start)
	ctx := context.Background()
	f.svc.Join(ctx, "q1@x.test", room.ID, &start)
	f.svc.Join(ctx, "q2@x.test", room.ID, &start)

	// Not free yet: refused, and nothing changed.
	if _, err := f.svc.Reserve(ctx, "q1@x.test", room.ID, &start); !errors.Is(err, ErrNotAvailable) {
		t.Fatalf("reserve while taken: %v", err)
	}
	if _, err := f.svc.Cancel(ctx, CancelRequest{ItemID: room.ID, Start: &start, ActorEmail: "owner@x.test"}); err != nil {
		t.Fatal(err)
	}
	f.ntf.reset()

	res, err := f.svc.Reserve(ctx, "q1@x.test", room.ID, &start)
	if err != nil {
		t.Fatal(err)
	}
	if res.User.Email != "q1@x.test" || !res.Start.Equal(start) || res.End == nil || !res.End.Equal(start.Add(time.Hour)) {
		t.Fatalf("result = %+v", res)
	}
	if f.holder(room, &start) != "q1@x.test" {
		t.Fatal("the reservation was not created")
	}
	if st, _ := f.svc.Status(ctx, "q1@x.test", room.ID, &start); st.Joined {
		t.Fatal("the user is still in the queue")
	}
	if st, _ := f.svc.Status(ctx, "q2@x.test", room.ID, &start); !st.Joined || st.Position != 1 {
		t.Fatalf("q2 = %+v", st)
	}
	if len(f.ntf.all()) != 0 {
		t.Fatal("booking notified somebody")
	}
	// Somebody who is not in the queue cannot use the action.
	if _, err := f.svc.Reserve(ctx, "owner@x.test", room.ID, &start); !errors.Is(err, ErrNotQueued) {
		t.Fatalf("a stranger reserved through the queue: %v", err)
	}
}

func TestReserveAnItemOnLoan(t *testing.T) {
	f := newFx(t)
	f.user("owner@x.test", "O", "X", 1)
	f.user("q1@x.test", "Q", "X", 1)
	cable := f.item("Câble", false)
	f.lend("owner@x.test", cable, t0.Add(-time.Hour))
	ctx := context.Background()
	f.svc.Join(ctx, "q1@x.test", cable.ID, nil)
	f.exec(`UPDATE reservation SET end_date = $1`, t0)

	if _, err := f.svc.Reserve(ctx, "q1@x.test", cable.ID, nil); err != nil {
		t.Fatal(err)
	}
	if f.holder(cable, nil) != "q1@x.test" {
		t.Fatal("the loan was not created")
	}
	var created sql.NullTime
	if err := f.db.QueryRow(`SELECT created_at FROM reservation WHERE email = 'q1@x.test'`).Scan(&created); err != nil || !created.Valid {
		t.Fatalf("created_at = %v, %v", created, err)
	}
}

func TestDecliningPassesTheBatonDownTheQueueUntilItIsEmpty(t *testing.T) {
	f := newFx(t)
	for _, u := range []string{"owner", "q1", "q2", "q3", "q4"} {
		f.user(u+"@x.test", u, "X", 1)
	}
	room := f.item("Salle", true)
	start := t0.Add(26 * time.Hour)
	f.book("owner@x.test", room, start)
	ctx := context.Background()
	for _, u := range []string{"q1", "q2", "q3", "q4"} {
		f.svc.Join(ctx, u+"@x.test", room.ID, &start)
	}
	f.svc.Cancel(ctx, CancelRequest{ItemID: room.ID, Start: &start, ActorEmail: "owner@x.test"})
	if got := f.ntf.all(); len(got) != 1 || to(t, got[0]) != "q1@x.test" {
		t.Fatalf("first notification = %+v", got)
	}

	// q3 leaves without holding the baton: nobody is notified, the order of the others is unchanged.
	if next, err := f.svc.Remove(ctx, "q3@x.test", room.ID, &start); err != nil || next {
		t.Fatalf("leaving without the baton: %v, %v", next, err)
	}
	f.ntf.reset()

	// q1 declines: q2 is told at once, and only q2.
	next, err := f.svc.Remove(ctx, "q1@x.test", room.ID, &start)
	if err != nil || !next {
		t.Fatalf("q1 declines: %v, %v", next, err)
	}
	if got := f.ntf.all(); len(got) != 1 || to(t, got[0]) != "q2@x.test" || got[0].ActionCategory != WaitlistActionCategory {
		t.Fatalf("after q1: %+v", got)
	}
	// q2 declines: q4 (q3 left).
	if next, _ := f.svc.Remove(ctx, "q2@x.test", room.ID, &start); !next {
		t.Fatal("q4 was not notified")
	}
	if got := f.ntf.all(); len(got) != 2 || to(t, got[1]) != "q4@x.test" {
		t.Fatalf("after q2: %+v", got)
	}
	// q4 declines: nobody is left.
	if next, err := f.svc.Remove(ctx, "q4@x.test", room.ID, &start); err != nil || next {
		t.Fatalf("q4 declines: %v, %v", next, err)
	}
	if n := f.count(`SELECT count(*) FROM reservation_waitlist`); n != 0 {
		t.Fatalf("%d entries left", n)
	}
	if len(f.ntf.all()) != 2 {
		t.Fatal("a notification went out to an empty queue")
	}
	if _, err := f.svc.Remove(ctx, "q4@x.test", room.ID, &start); !errors.Is(err, ErrNotQueued) {
		t.Fatalf("declining twice: %v", err)
	}
}

func TestPeopleWhoCannotReceiveThePushArePassedOver(t *testing.T) {
	f := newFx(t)
	for _, u := range []string{"owner", "q1", "q2"} {
		f.user(u+"@x.test", u, "X", 1)
	}
	room := f.item("Salle", true)
	start := t0.Add(26 * time.Hour)
	f.book("owner@x.test", room, start)
	ctx := context.Background()
	f.svc.Join(ctx, "q1@x.test", room.ID, &start)
	f.svc.Join(ctx, "q2@x.test", room.ID, &start)
	f.ntf.unreachable["q1@x.test"] = true

	f.svc.Cancel(ctx, CancelRequest{ItemID: room.ID, Start: &start, ActorEmail: "owner@x.test"})
	got := f.ntf.all()
	if len(got) != 1 || to(t, got[0]) != "q2@x.test" {
		t.Fatalf("announcements = %+v", got)
	}
	// q1 keeps their place at the head of the queue.
	if st, _ := f.svc.Status(ctx, "q1@x.test", room.ID, &start); !st.Joined || st.Position != 1 || st.Notified {
		t.Fatalf("q1 = %+v", st)
	}
}

func TestAWaitingListNotificationThatCannotBeQueuedDoesNotHoldTheBaton(t *testing.T) {
	f := newFx(t)
	f.user("owner@x.test", "O", "X", 1)
	f.user("q1@x.test", "Q", "X", 1)
	room := f.item("Salle", true)
	start := t0.Add(26 * time.Hour)
	f.book("owner@x.test", room, start)
	ctx := context.Background()
	f.svc.Join(ctx, "q1@x.test", room.ID, &start)
	f.ntf.fail = errors.New("queue down")
	if _, err := f.svc.Cancel(ctx, CancelRequest{ItemID: room.ID, Start: &start, ActorEmail: "owner@x.test"}); err != nil {
		t.Fatal(err)
	}
	if st, _ := f.svc.Status(ctx, "q1@x.test", room.ID, &start); st.Notified {
		t.Fatal("somebody holds a baton nobody told them about")
	}
}

func TestBookingDirectlyForgetsTheUsersQueueEntry(t *testing.T) {
	f := newFx(t)
	f.user("owner@x.test", "O", "X", 1)
	f.user("q1@x.test", "Q", "X", 1)
	cable := f.item("Câble", false)
	f.lend("owner@x.test", cable, t0.Add(-time.Hour))
	ctx := context.Background()
	f.svc.Join(ctx, "q1@x.test", cable.ID, nil)
	if err := f.svc.ForgetInterest(ctx, "q1@x.test", cable.ID, nil); err != nil {
		t.Fatal(err)
	}
	if n := f.count(`SELECT count(*) FROM reservation_waitlist`); n != 0 {
		t.Fatal("the entry stayed")
	}
}

func TestTheQueueIsLimitedPerUser(t *testing.T) {
	f := newFx(t)
	f.user("owner@x.test", "O", "X", 1)
	f.user("q1@x.test", "Q", "X", 1)
	ctx := context.Background()
	for i := 0; i < maxQueueEntries; i++ {
		it := f.item(fmt.Sprintf("Item %d", i), false)
		f.lend("owner@x.test", it, t0.Add(-time.Hour))
		if _, err := f.svc.Join(ctx, "q1@x.test", it.ID, nil); err != nil {
			t.Fatal(err)
		}
	}
	extra := f.item("One more", false)
	f.lend("owner@x.test", extra, t0.Add(-time.Hour))
	if _, err := f.svc.Join(ctx, "q1@x.test", extra.ID, nil); !errors.Is(err, ErrQueueFull) {
		t.Fatalf("error = %v", err)
	}
}

// ---------------------------------------------------------------------------------------------------------------
// Email template

func renderCancelEmail(t *testing.T, lang string, data CancelEmailData) (subject, body string) {
	t.Helper()
	initI18n(t)
	loc := appI18n.GetLocalizer(lang)
	raw, err := os.ReadFile("../../" + emailTemplate)
	if err != nil {
		t.Fatal(err)
	}
	translate := func(key string, _ ...interface{}) template.HTML {
		msg, err := loc.Localize(&goi18n.LocalizeConfig{MessageID: key, TemplateData: data})
		if err != nil {
			return template.HTML("!!!" + key + "!!!")
		}
		return template.HTML(msg)
	}
	tmpl, err := template.New("cancel").Parse(string(raw))
	if err != nil {
		t.Fatal(err)
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, struct {
		Data interface{}
		T    func(string, ...interface{}) template.HTML
	}{data, translate}); err != nil {
		t.Fatal(err)
	}
	subject, err = loc.Localize(&goi18n.LocalizeConfig{MessageID: emailSubjectKey, TemplateData: data})
	if err != nil {
		t.Fatal(err)
	}
	return subject, buf.String()
}

func TestCancellationEmailUsesTheThemeOfTheOtherTemplatesAndEscapesNames(t *testing.T) {
	hostile := CancelEmailData{
		ItemName: `<script>alert(1)</script>`, HasSlot: true, Slot: "02/07/2026 14:00",
		AdminName: `<img src=x onerror=alert(2)>`, CancelledAt: "01/07/2026 14:00",
	}
	for _, lang := range []string{"fr", "en", "de", "es"} {
		for _, hasSlot := range []bool{true, false} {
			d := hostile
			d.HasSlot = hasSlot
			subject, body := renderCancelEmail(t, lang, d)
			for _, raw := range []string{"<script", "<img src=x"} {
				if strings.Contains(body, raw) {
					t.Fatalf("%s: unescaped %q", lang, raw)
				}
			}
			if !strings.Contains(body, "&lt;img src=x onerror=alert(2)&gt;") {
				t.Fatalf("%s: the admin's name is not shown", lang)
			}
			if strings.Contains(body, "!!!") || strings.Contains(subject, "!!!") {
				t.Fatalf("%s: missing translation", lang)
			}
			if !strings.HasPrefix(subject, "[Transat]") {
				t.Fatalf("%s: subject %q", lang, subject)
			}
			if hasSlot != strings.Contains(body, "02/07/2026 14:00") {
				t.Fatalf("%s: slot row shown = %v", lang, !hasSlot)
			}
			// Same visual identity as the existing templates: same palette, card, logos, footer and fonts.
			for _, want := range []string{
				`background-color:#0D0505`, `background-color:#181010`, `color:#ffe6cc`, `https://transat.dev/logo_header.png`,
				`https://transat.dev/logo_footer.png`, `border-radius:10px`, `font-family:roboto,'helvetica neue',helvetica,arial,sans-serif`,
			} {
				if !strings.Contains(body, want) {
					t.Fatalf("%s: missing %q", lang, want)
				}
			}
		}
	}
	_, body := renderCancelEmail(t, "fr", CancelEmailData{ItemName: "Salle", AdminName: "Ada Lovelace", CancelledAt: "x"})
	if !strings.Contains(body, "Ada Lovelace") || !strings.Contains(body, "Annulée par") {
		t.Fatal("the admin's name is missing from the email")
	}
}
