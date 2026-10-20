package admin_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
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
)

type rsvMailer struct {
	mu   sync.Mutex
	sent []rsvMail
}

type rsvMail struct {
	mail models.Email
	data reservations.CancelEmailData
}

func (m *rsvMailer) SendEmail(mail models.Email, data interface{}) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.sent = append(m.sent, rsvMail{mail, data.(reservations.CancelEmailData)})
	return nil
}

type rsvAPI struct {
	*api
	mail    *rsvMailer
	club    int
	devices int
}

// rawFirst is the first element of a reply that is a JSON array.
func (r reply) rawFirst() map[string]any {
	var list []map[string]any
	if err := json.Unmarshal([]byte(r.raw), &list); err != nil || len(list) == 0 {
		return nil
	}
	return list[0]
}

func newReservationsAPI(t *testing.T) *rsvAPI {
	t.Setenv("JWT_SECRET", "test-secret-test-secret-test-secret")
	db := dbtest.Open(t)
	middlewares.InitJWTMiddleware(db)
	ns := services.NewNotificationService(db)
	a := &api{t: t, db: db, provider: &provider{}}
	cfg := push.DefaultConfig()
	cfg.MaxPerSecond = 0
	a.pipeline = push.New(db, a.provider, cfg, devices.DefaultConfig())
	ns.SetPipeline(a.pipeline)

	mail := &rsvMailer{}
	svc := reservations.New(db, ns, mail)
	svc.SetBackground(func(f func()) { f() })
	a.app = fiber.New()
	adminGroup := routes.SetupAdminRoutes(a.app, db, ns)
	routes.SetupAdminReservationRoutes(adminGroup, db, nil, svc)
	routes.SetupReservationRoutes(a.app, db, nil, svc)

	r := &rsvAPI{api: a, mail: mail}
	a.exec(`INSERT INTO roles (name) VALUES ('ADMIN'), ('NEWF') ON CONFLICT (name) DO NOTHING`)
	if err := db.QueryRow(`INSERT INTO clubs (name, picture) VALUES ('Club', 'p') RETURNING id_clubs`).Scan(&r.club); err != nil {
		t.Fatal(err)
	}
	return r
}

// person creates a student with a phone that receives reservation notifications.
func (r *rsvAPI) person(email, first, last string, roles ...string) {
	r.t.Helper()
	r.user(email, 1, "RESERVATIONS")
	r.exec(`UPDATE newf SET first_name = $2, last_name = $3 WHERE email = $1`, email, first, last)
	roles = append(roles, "NEWF")
	for _, role := range roles {
		r.exec(`INSERT INTO newf_roles (email, id_roles) SELECT $1, id_roles FROM roles WHERE name = $2`, email, role)
	}
	r.devices++
	r.device(1000+r.devices, email, true)
}

func (r *rsvAPI) item(name string, slot bool) int {
	r.t.Helper()
	var id int
	if err := r.db.QueryRow(`INSERT INTO reservation_element (name, slot, id_clubs) VALUES ($1, $2, $3) RETURNING id_reservation_element`,
		name, slot, r.club).Scan(&id); err != nil {
		r.t.Fatal(err)
	}
	return id
}

func (r *rsvAPI) deliver() []push.Outbound {
	r.t.Helper()
	for i := 0; i < 3; i++ {
		if _, err := r.pipeline.RunOnce(context.Background()); err != nil {
			r.t.Fatal(err)
		}
	}
	r.provider.mu.Lock()
	defer r.provider.mu.Unlock()
	out := r.provider.sent
	r.provider.sent = nil
	return out
}

func rsvSlot() string {
	return time.Now().UTC().Truncate(time.Hour).Add(48 * time.Hour).Format(reservations.APITimeLayout)
}

func TestTheReservationAdministrationRoutesAreForAdministratorsOnly(t *testing.T) {
	a := newReservationsAPI(t)
	a.person("admin@x.test", "Ada", "Lovelace", "ADMIN")
	a.person("student@x.test", "Sam", "Student")
	item := a.item("Salle", true)

	calls := []struct{ method, path string }{
		{http.MethodGet, "/admin/reservations"},
		{http.MethodGet, "/admin/reservations/stats"},
		{http.MethodGet, "/admin/reservations/items"},
		{http.MethodPost, "/admin/reservations/cancel"},
		{http.MethodPost, fmt.Sprintf("/admin/reservations/items/%d/return", item)},
		{http.MethodGet, fmt.Sprintf("/admin/reservations/waitlist?item_id=%d", item)},
		{http.MethodPost, "/admin/reservations/waitlist/advance"},
		{http.MethodDelete, "/admin/reservations/waitlist/1"},
		{http.MethodPost, "/admin/reservations/items"},
		{http.MethodPost, "/admin/reservations/categories"},
	}
	for _, c := range calls {
		if r := a.call(c.method, c.path, "student@x.test", map[string]any{}); r.status != http.StatusForbidden {
			t.Fatalf("a student got %d on %s %s", r.status, c.method, c.path)
		}
		if r := a.call(c.method, c.path, "", nil); r.status != http.StatusUnauthorized {
			t.Fatalf("no token: %d on %s %s", r.status, c.method, c.path)
		}
	}
	for _, path := range []string{"/admin/reservations", "/admin/reservations/stats", "/admin/reservations/items"} {
		if r := a.call(http.MethodGet, path, "admin@x.test", nil); r.status != 200 {
			t.Fatalf("admin: %d on %s", r.status, path)
		}
	}
}

func TestAnAdministratorCancelsFromThePanelAndTheQueueMoves(t *testing.T) {
	a := newReservationsAPI(t)
	a.person("admin@x.test", "Ada", "Lovelace", "ADMIN")
	a.person("bob@x.test", "Bob", "Martin")
	a.person("q1@x.test", "Quinn", "One")
	room := a.item("Salle Ampère", true)
	start := rsvSlot()
	path := fmt.Sprintf("/reservation/item/%d", room)

	if r := a.call(http.MethodPatch, path, "bob@x.test", map[string]string{"start_date": start}); r.status != 200 {
		t.Fatalf("booking: %d %v", r.status, r.body)
	}
	if r := a.call(http.MethodPost, path+"/waitlist", "q1@x.test", map[string]string{"start_date": start}); r.status != 201 {
		t.Fatalf("waitlist: %d %v", r.status, r.body)
	}

	// The panel sees the reservation and its queue.
	items := a.call(http.MethodGet, "/admin/reservations?status=upcoming&q=martin", "admin@x.test", nil)
	list, _ := items.body["items"].([]any)
	if items.status != 200 || len(list) != 1 || list[0].(map[string]any)["cancellable"] != true {
		t.Fatalf("upcoming: %d %v", items.status, items.body)
	}
	queue := a.call(http.MethodGet, fmt.Sprintf("/admin/reservations/waitlist?item_id=%d&start_date=%s", room, url.QueryEscape(start)), "admin@x.test", nil)
	if entries, _ := queue.body["entries"].([]any); queue.status != 200 || len(entries) != 1 || queue.body["free"] != false {
		t.Fatalf("queue: %d %v", queue.status, queue.body)
	}

	// Cancelling: the user and the first of the queue are told.
	r := a.call(http.MethodPost, "/admin/reservations/cancel", "admin@x.test", map[string]any{"item_id": room, "start_date": start})
	if r.status != 200 || r.body["cancelled_by_admin"] != true || r.body["waitlist_notified"] != true || r.body["owner"] != "bob@x.test" {
		t.Fatalf("cancel: %d %v", r.status, r.body)
	}
	sent := a.deliver()
	byUser := map[string]push.Outbound{}
	for _, m := range sent {
		byUser[m.Token] = m
	}
	bob, quinn := byUser[tokenFor(a.api, "bob@x.test")], byUser[tokenFor(a.api, "q1@x.test")]
	if len(sent) != 2 || !strings.Contains(bob.Body, "Ada Lovelace") || !strings.Contains(bob.Body, "Salle Ampère") {
		t.Fatalf("pushes = %+v", sent)
	}
	if quinn.ActionCategory != reservations.WaitlistActionCategory {
		t.Fatalf("the first of the queue got %+v", quinn)
	}
	if len(a.mail.sent) != 1 || a.mail.sent[0].mail.Recipient != "bob@x.test" || a.mail.sent[0].data.AdminName != "Ada Lovelace" {
		t.Fatalf("emails = %+v", a.mail.sent)
	}

	// The panel keeps the history, and the queue now shows who holds the turn.
	cancelled := a.call(http.MethodGet, "/admin/reservations?status=cancelled", "admin@x.test", nil)
	lines, _ := cancelled.body["items"].([]any)
	if len(lines) != 1 || lines[0].(map[string]any)["cancel_kind"] != "cancelled_by_admin" || lines[0].(map[string]any)["cancelled_by"] != "admin@x.test" {
		t.Fatalf("cancelled: %v", cancelled.body)
	}
	queue = a.call(http.MethodGet, fmt.Sprintf("/admin/reservations/waitlist?item_id=%d&start_date=%s", room, url.QueryEscape(start)), "admin@x.test", nil)
	if queue.body["free"] != true || queue.body["awaiting_response"] != true {
		t.Fatalf("queue after release: %v", queue.body)
	}

	// The turn can be passed on by hand, and entries removed.
	if r := a.call(http.MethodPost, "/admin/reservations/waitlist/advance", "admin@x.test", map[string]any{"item_id": room, "start_date": start}); r.status != 200 || r.body["next_notified"] != false {
		t.Fatalf("advance: %d %v", r.status, r.body)
	}
	if r := a.call(http.MethodPost, "/admin/reservations/cancel", "admin@x.test", map[string]any{"item_id": room, "start_date": start}); r.status != 404 {
		t.Fatalf("cancelling twice: %d", r.status)
	}
	if r := a.call(http.MethodPost, "/admin/reservations/cancel", "admin@x.test", map[string]any{"item_id": room}); r.status != 400 {
		t.Fatalf("slot item without slot: %d", r.status)
	}
	if r := a.call(http.MethodDelete, "/admin/reservations/waitlist/424242", "admin@x.test", nil); r.status != 404 {
		t.Fatalf("unknown entry: %d", r.status)
	}
}

func tokenFor(a *api, email string) string {
	var token string
	if err := a.db.QueryRow(`SELECT push_token FROM devices WHERE user_email = $1`, email).Scan(&token); err != nil {
		a.t.Fatal(err)
	}
	return token
}

func TestAnAdministratorMarksALoanAsReturned(t *testing.T) {
	a := newReservationsAPI(t)
	a.person("admin@x.test", "Ada", "Lovelace", "ADMIN")
	a.person("bob@x.test", "Bob", "Martin")
	cable := a.item("Câble", false)
	a.exec(`INSERT INTO reservation (email, id_reservation_element, start_date) VALUES ('bob@x.test', $1, now() AT TIME ZONE 'UTC' - interval '9 days')`, cable)

	late := a.call(http.MethodGet, "/admin/reservations?status=active&overdue_days=7", "admin@x.test", nil)
	lines, _ := late.body["items"].([]any)
	if len(lines) != 1 || lines[0].(map[string]any)["overdue_level"] != float64(7) {
		t.Fatalf("late loans: %v", late.body)
	}
	catalogue := a.call(http.MethodGet, "/admin/reservations/items", "admin@x.test", nil)
	if first := catalogue.rawFirst(); first["availability"] != "borrowed" {
		t.Fatalf("catalogue = %v", first)
	}

	r := a.call(http.MethodPost, fmt.Sprintf("/admin/reservations/items/%d/return", cable), "admin@x.test", nil)
	if r.status != 200 || r.body["owner"] != "bob@x.test" {
		t.Fatalf("return: %d %v", r.status, r.body)
	}
	if len(a.deliver()) != 0 || len(a.mail.sent) != 0 {
		t.Fatal("somebody was notified of a return")
	}
	if got := a.call(http.MethodGet, "/admin/reservations/items", "admin@x.test", nil).rawFirst(); got["availability"] != "available" {
		t.Fatalf("catalogue after = %v", got)
	}
	again := a.call(http.MethodPost, fmt.Sprintf("/admin/reservations/items/%d/return", cable), "admin@x.test", nil)
	if again.status != 404 {
		t.Fatalf("returning twice: %d", again.status)
	}
	room := a.item("Salle", true)
	if r := a.call(http.MethodPost, fmt.Sprintf("/admin/reservations/items/%d/return", room), "admin@x.test", nil); r.status != 400 {
		t.Fatalf("a slot item: %d", r.status)
	}
	history := a.call(http.MethodGet, "/admin/reservations?status=past&kind=loan", "admin@x.test", nil)
	if past, _ := history.body["items"].([]any); len(past) != 1 || past[0].(map[string]any)["cancel_kind"] != "returned_by_admin" {
		t.Fatalf("history: %v", history.body)
	}
}

func TestAnItemUnderMaintenanceDisappearsFromTheAppAndCannotBeBooked(t *testing.T) {
	a := newReservationsAPI(t)
	a.person("admin@x.test", "Ada", "Lovelace", "ADMIN")
	a.person("bob@x.test", "Bob", "Martin")

	created := a.call(http.MethodPost, "/admin/reservations/items", "admin@x.test",
		map[string]any{"name": "Vidéoprojecteur", "slot": false, "id_club_parent": a.club, "status": "maintenance"})
	if created.status != 201 {
		t.Fatalf("create: %d %v", created.status, created.body)
	}
	item := int(created.body["item"].(map[string]any)["id"].(float64))
	if r := a.call(http.MethodPost, "/admin/reservations/items", "admin@x.test", map[string]any{"name": "X", "id_club_parent": a.club, "status": "broken"}); r.status != 400 {
		t.Fatalf("an unknown status: %d", r.status)
	}

	catalog := func() int {
		r := a.call(http.MethodGet, fmt.Sprintf("/reservation/club/%d", a.club), "bob@x.test", nil)
		items, _ := r.body["items"].([]any)
		return len(items)
	}
	if n := catalog(); n != 0 {
		t.Fatalf("the app lists %d items in maintenance", n)
	}
	now := time.Now().UTC().Format(reservations.APITimeLayout)
	if r := a.call(http.MethodPatch, fmt.Sprintf("/reservation/item/%d", item), "bob@x.test", map[string]string{"start_date": now}); r.status != http.StatusLocked {
		t.Fatalf("booking an item in maintenance: %d %v", r.status, r.body)
	}

	// Back in service.
	if r := a.call(http.MethodPatch, fmt.Sprintf("/admin/reservations/items/%d", item), "admin@x.test", map[string]any{"status": "active", "status_note": ""}); r.status != 200 {
		t.Fatalf("activate: %d", r.status)
	}
	if n := catalog(); n != 1 {
		t.Fatalf("the app lists %d items", n)
	}
	if r := a.call(http.MethodPatch, fmt.Sprintf("/reservation/item/%d", item), "bob@x.test", map[string]string{"start_date": now}); r.status != 200 {
		t.Fatalf("booking: %d %v", r.status, r.body)
	}

	// With a loan in progress the type cannot change, nor can the item be deleted.
	if r := a.call(http.MethodPatch, fmt.Sprintf("/admin/reservations/items/%d", item), "admin@x.test", map[string]any{"slot": true}); r.status != http.StatusConflict {
		t.Fatalf("changing the type under a loan: %d", r.status)
	}
	if r := a.call(http.MethodDelete, fmt.Sprintf("/admin/reservations/items/%d", item), "admin@x.test", nil); r.status != 400 {
		t.Fatalf("deleting a borrowed item: %d", r.status)
	}
	if r := a.call(http.MethodPatch, fmt.Sprintf("/admin/reservations/items/%d", item), "admin@x.test", map[string]any{"status": "nope"}); r.status != 400 {
		t.Fatalf("invalid status: %d", r.status)
	}
	// An upcoming slot blocks the deletion of a slot item too.
	room := a.item("Salle", true)
	if r := a.call(http.MethodPatch, fmt.Sprintf("/reservation/item/%d", room), "bob@x.test", map[string]string{"start_date": rsvSlot()}); r.status != 200 {
		t.Fatalf("slot booking: %d", r.status)
	}
	if r := a.call(http.MethodDelete, fmt.Sprintf("/admin/reservations/items/%d", room), "admin@x.test", nil); r.status != 400 {
		t.Fatalf("deleting an item with an upcoming slot: %d", r.status)
	}
}

func TestTheWholeCascadeFromAnAdminCancellationDownTheQueue(t *testing.T) {
	a := newReservationsAPI(t)
	for _, p := range []struct{ email, first string }{{"owner@x.test", "Olivia"}, {"q1@x.test", "Quinn"}, {"q2@x.test", "Quincy"}, {"q3@x.test", "Quentin"}} {
		a.person(p.email, p.first, "X")
	}
	a.person("admin@x.test", "Ada", "Lovelace", "ADMIN")
	room := a.item("Salle", true)
	start := rsvSlot()
	item := fmt.Sprintf("/reservation/item/%d", room)
	body := map[string]string{"start_date": start}

	a.call(http.MethodPatch, item, "owner@x.test", body)
	for _, u := range []string{"q1", "q2", "q3"} {
		if r := a.call(http.MethodPost, item+"/waitlist", u+"@x.test", body); r.status != 201 {
			t.Fatalf("%s joins: %d", u, r.status)
		}
	}
	tokens := func(sent []push.Outbound) string {
		var out []string
		for _, m := range sent {
			out = append(out, m.Token)
		}
		return strings.Join(out, ",")
	}
	only := func(sent []push.Outbound, who string) {
		t.Helper()
		if tokens(sent) != tokenFor(a.api, who+"@x.test") {
			t.Fatalf("pushes went to %q, want only %s", tokens(sent), who)
		}
		if sent[0].ActionCategory != reservations.WaitlistActionCategory {
			t.Fatalf("no Réserver / Annuler buttons: %+v", sent[0])
		}
	}

	// Cancelled from the panel: the owner is told, and only the first of the queue.
	a.call(http.MethodPost, "/admin/reservations/cancel", "admin@x.test", map[string]any{"item_id": room, "start_date": start})
	sent := a.deliver()
	if len(sent) != 2 {
		t.Fatalf("%d pushes after the cancellation", len(sent))
	}
	var waitlistPush []push.Outbound
	for _, m := range sent {
		if m.ActionCategory != "" {
			waitlistPush = append(waitlistPush, m)
		}
	}
	only(waitlistPush, "q1")

	// "Annuler" passes it on, in order, to the next person only.
	if r := a.call(http.MethodPost, item+"/waitlist/decline", "q1@x.test", body); r.status != 200 || r.body["next_notified"] != true {
		t.Fatalf("q1 declines: %d %v", r.status, r.body)
	}
	only(a.deliver(), "q2")

	// "Réserver" books it, takes q2 out of the queue and tells nobody else.
	if r := a.call(http.MethodPost, item+"/waitlist/reserve", "q2@x.test", body); r.status != 200 {
		t.Fatalf("q2 reserves: %d %v", r.status, r.body)
	}
	if got := a.deliver(); len(got) != 0 {
		t.Fatalf("booking notified %q", tokens(got))
	}
	queue := a.call(http.MethodGet, fmt.Sprintf("/admin/reservations/waitlist?item_id=%d&start_date=%s", room, url.QueryEscape(start)), "admin@x.test", nil)
	entries, _ := queue.body["entries"].([]any)
	if len(entries) != 1 || entries[0].(map[string]any)["email"] != "q3@x.test" || entries[0].(map[string]any)["position"] != float64(1) || queue.body["awaiting_response"] != false {
		t.Fatalf("queue = %v", queue.body)
	}

	// When q2 gives it up, the cascade resumes with q3.
	a.call(http.MethodDelete, item, "q2@x.test", body)
	only(a.deliver(), "q3")
}
