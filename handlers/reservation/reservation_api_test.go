package reservation_test

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/plugimt/transat-backend/internal/database/dbtest"
	"github.com/plugimt/transat-backend/middlewares"
	"github.com/plugimt/transat-backend/models"
	"github.com/plugimt/transat-backend/routes"
	"github.com/plugimt/transat-backend/services"
	"github.com/plugimt/transat-backend/services/devices"
	"github.com/plugimt/transat-backend/services/push"
	"github.com/plugimt/transat-backend/services/reservations"
	"github.com/plugimt/transat-backend/utils"
)

// provider records what the delivery pipeline hands to the push provider.
type provider struct {
	mu   sync.Mutex
	sent []push.Outbound
}

func (p *provider) Name() string  { return "test" }
func (p *provider) MaxBatch() int { return 100 }
func (p *provider) Send(_ context.Context, msgs []push.Outbound) ([]push.SendResult, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.sent = append(p.sent, msgs...)
	out := make([]push.SendResult, len(msgs))
	for i := range out {
		out[i] = push.SendResult{Kind: push.ResultAccepted, TicketID: fmt.Sprintf("t%d", len(p.sent)+i)}
	}
	return out, nil
}
func (p *provider) Receipts(context.Context, []string) (map[string]push.ReceiptResult, error) {
	return nil, nil
}

// take returns what was sent since the last call.
func (p *provider) take() []push.Outbound {
	p.mu.Lock()
	defer p.mu.Unlock()
	out := p.sent
	p.sent = nil
	return out
}

type mailer struct {
	mu   sync.Mutex
	sent []sentMail
}

type sentMail struct {
	mail models.Email
	data reservations.CancelEmailData
}

func (m *mailer) SendEmail(mail models.Email, data interface{}) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.sent = append(m.sent, sentMail{mail, data.(reservations.CancelEmailData)})
	return nil
}

type api struct {
	t        *testing.T
	db       *sql.DB
	app      *fiber.App
	svc      *reservations.Service
	pipeline *push.Pipeline
	provider *provider
	mail     *mailer
	club     int
	device   int
	tokens   map[string]string
}

func newAPI(t *testing.T) *api {
	t.Setenv("JWT_SECRET", "test-secret-test-secret-test-secret")
	db := dbtest.Open(t)
	middlewares.InitJWTMiddleware(db)
	a := &api{t: t, db: db, provider: &provider{}, mail: &mailer{}, tokens: map[string]string{}}

	ns := services.NewNotificationService(db)
	cfg := push.DefaultConfig()
	cfg.MaxPerSecond = 0
	a.pipeline = push.New(db, a.provider, cfg, devices.DefaultConfig())
	ns.SetPipeline(a.pipeline)

	a.svc = reservations.New(db, ns, a.mail)
	a.svc.SetBackground(func(f func()) { f() })
	a.app = fiber.New()
	routes.SetupReservationRoutes(a.app, db, nil, a.svc)

	a.exec(`INSERT INTO roles (name) VALUES ('ADMIN'), ('NEWF') ON CONFLICT (name) DO NOTHING`)
	if err := db.QueryRow(`INSERT INTO clubs (name, picture) VALUES ('Club', 'p') RETURNING id_clubs`).Scan(&a.club); err != nil {
		t.Fatal(err)
	}
	return a
}

func (a *api) exec(q string, args ...any) {
	a.t.Helper()
	if _, err := a.db.Exec(q, args...); err != nil {
		a.t.Fatalf("%s: %v", q, err)
	}
}

// user creates a student (language 1 = fr) with a phone that receives reservation notifications.
func (a *api) user(email, first, last string, roles ...string) {
	a.t.Helper()
	a.exec(`INSERT INTO newf (email, password, first_name, last_name, language, password_updated_date)
		VALUES ($1, 'x', $2, $3, 1, now() - interval '3 years')`, email, first, last)
	a.exec(`INSERT INTO newf_roles (email, id_roles) SELECT $1, id_roles FROM roles WHERE name = 'NEWF'`, email)
	for _, r := range roles {
		a.exec(`INSERT INTO newf_roles (email, id_roles) SELECT $1, id_roles FROM roles WHERE name = $2`, email, r)
	}
	a.exec(`INSERT INTO notifications (email, id_services) SELECT $1, id_services FROM services WHERE name = 'RESERVATIONS'`, email)
	a.device++
	a.tokens[email] = fmt.Sprintf("ExponentPushToken[device-%d]", a.device)
	_, err := devices.NewService(a.db, devices.DefaultConfig()).Register(context.Background(), devices.Registration{
		InstallationID: fmt.Sprintf("00000000-0000-4000-8000-%012d", a.device), Email: email, Platform: "android",
		PushToken: a.tokens[email], PermissionStatus: "granted",
	})
	if err != nil {
		a.t.Fatal(err)
	}
}

func (a *api) item(name string, slot bool) int {
	a.t.Helper()
	var id int
	if err := a.db.QueryRow(`INSERT INTO reservation_element (name, slot, id_clubs) VALUES ($1, $2, $3) RETURNING id_reservation_element`,
		name, slot, a.club).Scan(&id); err != nil {
		a.t.Fatal(err)
	}
	return id
}

type reply struct {
	status int
	body   map[string]any
}

func (a *api) call(method, path, as string, body any) reply {
	a.t.Helper()
	var buf io.Reader
	if body != nil {
		raw, _ := json.Marshal(body)
		buf = bytes.NewReader(raw)
	}
	req := httptest.NewRequest(method, path, buf)
	req.Header.Set("Content-Type", "application/json")
	token, err := utils.GenerateJWT(as, nil, "")
	if err != nil {
		a.t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	res, err := a.app.Test(req, -1)
	if err != nil {
		a.t.Fatal(err)
	}
	defer res.Body.Close()
	raw, _ := io.ReadAll(res.Body)
	out := reply{status: res.StatusCode}
	_ = json.Unmarshal(raw, &out.body)
	return out
}

// deliver runs the push worker and returns what reached the phones.
func (a *api) deliver() []push.Outbound {
	a.t.Helper()
	for i := 0; i < 3; i++ {
		if _, err := a.pipeline.RunOnce(context.Background()); err != nil {
			a.t.Fatal(err)
		}
	}
	return a.provider.take()
}

// holder returns who holds a slot item (any reservation) or a loan item (an open loan).
func (a *api) holder(item int, slot bool) string {
	a.t.Helper()
	query := `SELECT email FROM reservation WHERE id_reservation_element = $1 AND end_date IS NULL`
	if slot {
		query = `SELECT email FROM reservation WHERE id_reservation_element = $1`
	}
	var email sql.NullString
	if err := a.db.QueryRow(query+` LIMIT 1`, item).Scan(&email); err != nil && err != sql.ErrNoRows {
		a.t.Fatal(err)
	}
	return email.String
}

func (a *api) rows(item int) int {
	var n int
	if err := a.db.QueryRow(`SELECT count(*) FROM reservation WHERE id_reservation_element = $1`, item).Scan(&n); err != nil {
		a.t.Fatal(err)
	}
	return n
}

// slot is a plain hour two days from now, formatted the way the API expects.
func slot() (time.Time, string) {
	s := time.Now().UTC().Truncate(time.Hour).Add(48 * time.Hour)
	return s, s.Format(reservations.APITimeLayout)
}

func (a *api) tokenOf(email string) string { return a.tokens[email] }

func recipients(sent []push.Outbound) string {
	var out []string
	for _, m := range sent {
		out = append(out, m.Token)
	}
	return strings.Join(out, ",")
}

func TestAnAdminCancelsAThirdPartysSlotAndAUserCannot(t *testing.T) {
	a := newAPI(t)
	a.user("admin@x.test", "Ada", "Lovelace", "ADMIN")
	a.user("bob@x.test", "Bob", "Martin")
	a.user("carol@x.test", "Carol", "Dupont")
	room := a.item("Salle Ampère", true)
	_, start := slot()
	path := fmt.Sprintf("/reservation/item/%d", room)

	if r := a.call(http.MethodPatch, path, "bob@x.test", map[string]string{"start_date": start}); r.status != 200 {
		t.Fatalf("booking: %d %v", r.status, r.body)
	}

	// A student cannot cancel somebody else's reservation.
	r := a.call(http.MethodDelete, path, "carol@x.test", map[string]string{"start_date": start})
	if r.status != 404 || a.rows(room) != 1 {
		t.Fatalf("carol: %d %v, %d rows", r.status, r.body, a.rows(room))
	}
	if len(a.deliver()) != 0 || len(a.mail.sent) != 0 {
		t.Fatal("a refused cancellation notified somebody")
	}

	// An administrator can.
	r = a.call(http.MethodDelete, path, "admin@x.test", map[string]string{"start_date": start})
	if r.status != 200 || r.body["cancelled_by_admin"] != true || a.rows(room) != 0 {
		t.Fatalf("admin: %d %v, %d rows", r.status, r.body, a.rows(room))
	}

	// The user got the push notification, naming the item, the slot and the administrator.
	sent := a.deliver()
	if len(sent) != 1 || sent[0].Token != a.tokenOf("bob@x.test") {
		t.Fatalf("pushes = %v", recipients(sent))
	}
	if sent[0].Title != "Réservation annulée" || !strings.Contains(sent[0].Body, "Ada Lovelace") || !strings.Contains(sent[0].Body, "Salle Ampère") {
		t.Fatalf("push = %q / %q", sent[0].Title, sent[0].Body)
	}
	nav, _ := sent[0].Data["navigation"].(*models.NavigationTarget)
	if sent[0].Data["type"] != "reservation_cancelled" || nav == nil || nav.Type != models.NavigationReservation || nav.ID != fmt.Sprint(room) {
		t.Fatalf("data = %v", sent[0].Data)
	}

	// And the email, with the same information.
	if len(a.mail.sent) != 1 {
		t.Fatalf("%d emails", len(a.mail.sent))
	}
	m := a.mail.sent[0]
	if m.mail.Recipient != "bob@x.test" || m.data.AdminName != "Ada Lovelace" || m.data.ItemName != "Salle Ampère" || !m.data.HasSlot {
		t.Fatalf("email = %+v", m)
	}

	// The owner can still cancel their own slot, and that tells nobody.
	if r := a.call(http.MethodPatch, path, "bob@x.test", map[string]string{"start_date": start}); r.status != 200 {
		t.Fatalf("booking again: %d", r.status)
	}
	if r := a.call(http.MethodDelete, path, "bob@x.test", map[string]string{"start_date": start}); r.status != 200 || r.body["cancelled_by_admin"] != false {
		t.Fatalf("owner cancel: %d %v", r.status, r.body)
	}
	if len(a.deliver()) != 0 || len(a.mail.sent) != 1 {
		t.Fatal("cancelling one's own reservation notified somebody")
	}
}

func TestAnAdminEndsALoanWithoutBodyAndAStudentCannot(t *testing.T) {
	a := newAPI(t)
	a.user("admin@x.test", "Ada", "Lovelace", "ADMIN")
	a.user("bob@x.test", "Bob", "Martin")
	a.user("carol@x.test", "Carol", "Dupont")
	cable := a.item("Câble HDMI", false)
	path := fmt.Sprintf("/reservation/item/%d", cable)

	now := time.Now().UTC().Add(-time.Minute).Format(reservations.APITimeLayout)
	if r := a.call(http.MethodPatch, path, "bob@x.test", map[string]string{"start_date": now}); r.status != 200 {
		t.Fatalf("loan: %d %v", r.status, r.body)
	}
	if r := a.call(http.MethodDelete, path, "carol@x.test", nil); r.status != 400 {
		t.Fatalf("carol: %d %v", r.status, r.body)
	}
	if a.holder(cable, false) != "bob@x.test" {
		t.Fatal("the loan was ended by a student")
	}
	r := a.call(http.MethodDelete, path, "admin@x.test", nil)
	if r.status != 200 || r.body["cancelled_by_admin"] != true {
		t.Fatalf("admin: %d %v", r.status, r.body)
	}
	if a.holder(cable, false) != "" {
		t.Fatal("the loan is still open")
	}
	sent := a.deliver()
	if len(sent) != 1 || sent[0].Token != a.tokenOf("bob@x.test") || !strings.Contains(sent[0].Body, "Ada Lovelace") || !strings.Contains(sent[0].Body, "Câble HDMI") {
		t.Fatalf("pushes = %+v", sent)
	}
	if len(a.mail.sent) != 1 || a.mail.sent[0].data.HasSlot {
		t.Fatalf("emails = %+v", a.mail.sent)
	}
}

func TestTheWaitingListEndToEnd(t *testing.T) {
	a := newAPI(t)
	for _, u := range []string{"owner", "q1", "q2", "fast"} {
		a.user(u+"@x.test", u, "X")
	}
	room := a.item("Salle", true)
	_, start := slot()
	item := fmt.Sprintf("/reservation/item/%d", room)
	wait := item + "/waitlist"
	body := map[string]string{"start_date": start}

	// Nothing to wait for while the slot is free.
	if r := a.call(http.MethodPost, wait, "q1@x.test", body); r.status != 409 || r.body["code"] != "still_available" {
		t.Fatalf("free slot: %d %v", r.status, r.body)
	}
	if r := a.call(http.MethodPatch, item, "owner@x.test", body); r.status != 200 {
		t.Fatalf("booking: %d", r.status)
	}

	// Interest in a reserved slot.
	if r := a.call(http.MethodPost, wait, "q1@x.test", body); r.status != 201 || r.body["position"] != float64(1) {
		t.Fatalf("q1: %d %v", r.status, r.body)
	}
	if r := a.call(http.MethodPost, wait, "q2@x.test", body); r.status != 201 || r.body["position"] != float64(2) {
		t.Fatalf("q2: %d %v", r.status, r.body)
	}
	if r := a.call(http.MethodPost, wait, "q2@x.test", body); r.status != 409 || r.body["code"] != "already_queued" {
		t.Fatalf("q2 again: %d %v", r.status, r.body)
	}
	if r := a.call(http.MethodPost, wait, "owner@x.test", body); r.status != 409 || r.body["code"] != "own_reservation" {
		t.Fatalf("owner: %d %v", r.status, r.body)
	}
	if r := a.call(http.MethodPost, wait, "q1@x.test", nil); r.status != 400 || r.body["code"] != "start_required" {
		t.Fatalf("no slot: %d %v", r.status, r.body)
	}
	if r := a.call(http.MethodGet, wait+"?start_date="+strings.ReplaceAll(start, " ", "%20"), "q2@x.test", nil); r.status != 200 ||
		r.body["count"] != float64(2) || r.body["position"] != float64(2) || r.body["joined"] != true {
		t.Fatalf("status: %d %v", r.status, r.body)
	}
	if r := a.call(http.MethodGet, "/reservation/waitlist", "q1@x.test", nil); r.status != 200 {
		t.Fatalf("my waitlist: %d", r.status)
	} else if list, _ := r.body["waitlist"].([]any); len(list) != 1 {
		t.Fatalf("my waitlist = %v", r.body)
	}
	if len(a.deliver()) != 0 {
		t.Fatal("joining a queue notified somebody")
	}

	// The slot is released: only the first of the queue is told, with the actions.
	if r := a.call(http.MethodDelete, item, "owner@x.test", body); r.status != 200 || r.body["waitlist_notified"] != true {
		t.Fatalf("release: %d %v", r.status, r.body)
	}
	sent := a.deliver()
	if len(sent) != 1 || sent[0].Token != a.tokenOf("q1@x.test") {
		t.Fatalf("pushes = %q", recipients(sent))
	}
	if sent[0].ActionCategory != reservations.WaitlistActionCategory || !strings.Contains(sent[0].Title, "Salle") {
		t.Fatalf("push = %+v", sent[0])
	}
	nav, _ := sent[0].Data["navigation"].(*models.NavigationTarget)
	if nav == nil || nav.ID != fmt.Sprint(room) || nav.Params["waitlist"] != "1" || nav.Params["start"] != start {
		t.Fatalf("navigation = %+v", nav)
	}

	// Booking is not blocked: somebody who is not in the queue takes the slot first.
	if r := a.call(http.MethodPatch, item, "fast@x.test", body); r.status != 200 {
		t.Fatalf("direct booking was blocked: %d %v", r.status, r.body)
	}
	// "Réserver" now fails with a clear error, and q1 keeps their place.
	r := a.call(http.MethodPost, wait+"/reserve", "q1@x.test", body)
	if r.status != 409 || r.body["code"] != "not_available" || !strings.Contains(fmt.Sprint(r.body["error"]), "reserved by someone else") {
		t.Fatalf("reserve when taken: %d %v", r.status, r.body)
	}

	// The slot is free again: q1 (still first) is told again, and this time books it.
	if r := a.call(http.MethodDelete, item, "fast@x.test", body); r.status != 200 {
		t.Fatalf("release 2: %d", r.status)
	}
	if sent := a.deliver(); len(sent) != 1 || sent[0].Token != a.tokenOf("q1@x.test") {
		t.Fatalf("second round: %q", recipients(sent))
	}
	r = a.call(http.MethodPost, wait+"/reserve", "q1@x.test", body)
	if r.status != 200 {
		t.Fatalf("reserve: %d %v", r.status, r.body)
	}
	if a.holder(room, true) != "q1@x.test" {
		t.Fatal("the slot was not booked")
	}
	if r := a.call(http.MethodGet, wait+"?start_date="+strings.ReplaceAll(start, " ", "%20"), "q1@x.test", nil); r.body["joined"] != false || r.body["count"] != float64(1) {
		t.Fatalf("after reserving: %v", r.body)
	}
	if len(a.deliver()) != 0 {
		t.Fatal("booking from the queue notified somebody")
	}
}

func TestDecliningPassesTheRelayOverHTTP(t *testing.T) {
	a := newAPI(t)
	for _, u := range []string{"owner", "q1", "q2"} {
		a.user(u+"@x.test", u, "X")
	}
	cable := a.item("Câble", false)
	item := fmt.Sprintf("/reservation/item/%d", cable)
	wait := item + "/waitlist"

	loan := time.Now().UTC().Add(-time.Hour).Format(reservations.APITimeLayout)
	if r := a.call(http.MethodPatch, item, "owner@x.test", map[string]string{"start_date": loan}); r.status != 200 {
		t.Fatalf("loan: %d", r.status)
	}
	for _, u := range []string{"q1", "q2"} {
		if r := a.call(http.MethodPost, wait, u+"@x.test", nil); r.status != 201 {
			t.Fatalf("%s: %d %v", u, r.status, r.body)
		}
	}

	// The owner returns the item through the existing endpoint.
	end := time.Now().UTC().Format(reservations.APITimeLayout)
	if r := a.call(http.MethodPatch, item, "owner@x.test", map[string]string{"end_date": end}); r.status != 200 {
		t.Fatalf("return: %d %v", r.status, r.body)
	}
	if sent := a.deliver(); len(sent) != 1 || sent[0].Token != a.tokenOf("q1@x.test") {
		t.Fatalf("first: %q", recipients(sent))
	}

	// q2 cannot decline for q1; q1 declines and q2 is told at once.
	r := a.call(http.MethodPost, wait+"/decline", "q1@x.test", nil)
	if r.status != 200 || r.body["next_notified"] != true {
		t.Fatalf("decline: %d %v", r.status, r.body)
	}
	if sent := a.deliver(); len(sent) != 1 || sent[0].Token != a.tokenOf("q2@x.test") || sent[0].ActionCategory == "" {
		t.Fatalf("second: %q", recipients(sent))
	}
	r = a.call(http.MethodPost, wait+"/decline", "q2@x.test", nil)
	if r.status != 200 || r.body["next_notified"] != false {
		t.Fatalf("last decline: %d %v", r.status, r.body)
	}
	if r := a.call(http.MethodPost, wait+"/decline", "q2@x.test", nil); r.status != 404 || r.body["code"] != "not_queued" {
		t.Fatalf("declining twice: %d %v", r.status, r.body)
	}
	if len(a.deliver()) != 0 {
		t.Fatal("a notification went out to an empty queue")
	}
}

func TestRemindersThroughTheRealPipeline(t *testing.T) {
	a := newAPI(t)
	a.user("bob@x.test", "Bob", "Martin")
	room := a.item("Salle", true)
	cable := a.item("Câble", false)
	s, start := slot()
	if r := a.call(http.MethodPatch, fmt.Sprintf("/reservation/item/%d", room), "bob@x.test", map[string]string{"start_date": start}); r.status != 200 {
		t.Fatalf("booking: %d", r.status)
	}
	loanStart := time.Now().UTC().Add(-25 * time.Hour)
	if r := a.call(http.MethodPatch, fmt.Sprintf("/reservation/item/%d", cable), "bob@x.test", map[string]string{"start_date": loanStart.Format(reservations.APITimeLayout)}); r.status != 200 {
		t.Fatalf("loan: %d", r.status)
	}

	// One day into the loan, and 40 hours before the slot: only the loan reminder is due.
	sent, err := a.svc.SendDueReminders(context.Background())
	if err != nil || sent != 1 {
		t.Fatalf("sent %d, %v", sent, err)
	}
	got := a.deliver()
	if len(got) != 1 || !strings.Contains(got[0].Body, "Câble") || got[0].Data["type"] != "reservation_reminder" {
		t.Fatalf("pushes = %+v", got)
	}

	// An hour before the slot: the slot reminder. Not twice.
	a.svc.SetClock(func() time.Time { return s.Add(-time.Hour) })
	if sent, err := a.svc.SendDueReminders(context.Background()); err != nil || sent != 1 {
		t.Fatalf("slot reminder: %d, %v", sent, err)
	}
	if sent, _ := a.svc.SendDueReminders(context.Background()); sent != 0 {
		t.Fatalf("second run sent %d", sent)
	}
	got = a.deliver()
	if len(got) != 1 || got[0].Title != "Réservation dans 1h" || !strings.Contains(got[0].Body, "Salle") {
		t.Fatalf("pushes = %+v", got)
	}

	// Returning the item stops the reminders: the 5 day one never comes.
	end := time.Now().UTC().Format(reservations.APITimeLayout)
	if r := a.call(http.MethodPatch, fmt.Sprintf("/reservation/item/%d", cable), "bob@x.test", map[string]string{"end_date": end}); r.status != 200 {
		t.Fatalf("return: %d", r.status)
	}
	a.svc.SetClock(func() time.Time { return loanStart.Add(5*24*time.Hour + time.Hour) })
	if sent, _ := a.svc.SendDueReminders(context.Background()); sent != 0 {
		t.Fatalf("a reminder went out for a returned item: %d", sent)
	}
}
