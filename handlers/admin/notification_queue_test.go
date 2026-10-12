package admin_test

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

	"github.com/gofiber/fiber/v2"
	"github.com/plugimt/transat-backend/internal/database/dbtest"
	"github.com/plugimt/transat-backend/middlewares"
	"github.com/plugimt/transat-backend/routes"
	"github.com/plugimt/transat-backend/services"
	"github.com/plugimt/transat-backend/services/devices"
	"github.com/plugimt/transat-backend/services/push"
	"github.com/plugimt/transat-backend/utils"
)

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
func (p *provider) count() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	return len(p.sent)
}

type api struct {
	t        *testing.T
	db       *sql.DB
	app      *fiber.App
	pipeline *push.Pipeline
	provider *provider
}

func newAPI(t *testing.T, withPipeline bool) *api {
	t.Setenv("JWT_SECRET", "test-secret-test-secret-test-secret")
	db := dbtest.Open(t)
	middlewares.InitJWTMiddleware(db)
	ns := services.NewNotificationService(db)
	a := &api{t: t, db: db, provider: &provider{}}
	if withPipeline {
		cfg := push.DefaultConfig()
		cfg.MaxPerSecond = 0
		a.pipeline = push.New(db, a.provider, cfg, devices.DefaultConfig())
		ns.SetPipeline(a.pipeline)
	}
	a.app = fiber.New()
	routes.SetupAdminRoutes(a.app, db, ns)

	a.exec(`INSERT INTO roles (name) VALUES ('ADMIN') ON CONFLICT (name) DO NOTHING`)
	a.user("admin@x.test", 1)
	a.exec(`INSERT INTO newf_roles (email, id_roles) SELECT 'admin@x.test', id_roles FROM roles WHERE name = 'ADMIN'`)
	a.user("student@x.test", 1)
	return a
}

func (a *api) exec(q string, args ...any) {
	a.t.Helper()
	if _, err := a.db.Exec(q, args...); err != nil {
		a.t.Fatalf("%s: %v", q, err)
	}
}

func (a *api) user(email string, lang int, categories ...string) {
	a.exec(`INSERT INTO newf (email, password, first_name, last_name, language, password_updated_date)
		VALUES ($1, 'x', 'A', 'B', $2, now() - interval '3 years')`, email, lang)
	for _, c := range categories {
		a.exec(`INSERT INTO notifications (email, id_services) SELECT $1, id_services FROM services WHERE name = $2`, email, c)
	}
}

func (a *api) device(n int, email string, withToken bool) {
	a.t.Helper()
	token := ""
	if withToken {
		token = fmt.Sprintf("ExponentPushToken[device-%d]", n)
	}
	_, err := devices.NewService(a.db, devices.DefaultConfig()).Register(context.Background(), devices.Registration{
		InstallationID: fmt.Sprintf("00000000-0000-4000-8000-%012d", n), Email: email, Platform: "android",
		PushToken: token, PermissionStatus: "granted",
	})
	if err != nil {
		a.t.Fatal(err)
	}
}

type reply struct {
	status int
	body   map[string]any
	raw    string
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
	if as != "" {
		token, err := utils.GenerateJWT(as, nil, "")
		if err != nil {
			a.t.Fatal(err)
		}
		req.Header.Set("Authorization", "Bearer "+token)
	}
	res, err := a.app.Test(req, -1)
	if err != nil {
		a.t.Fatal(err)
	}
	defer res.Body.Close()
	raw, _ := io.ReadAll(res.Body)
	out := reply{status: res.StatusCode, raw: string(raw)}
	_ = json.Unmarshal(raw, &out.body)
	return out
}

func (a *api) send(body map[string]any) reply {
	if _, ok := body["title"]; !ok {
		body["title"] = "Hello"
	}
	if _, ok := body["audience"]; !ok {
		body["audience"] = map[string]any{"type": "all"}
	}
	return a.call(http.MethodPost, "/admin/notifications/send", "admin@x.test", body)
}

func TestSendingQueuesTheMessageAndAnswersWithoutSendingAnything(t *testing.T) {
	a := newAPI(t, true)
	a.user("a@x.test", 1, "RESTAURANT")
	a.user("b@x.test", 2, "RESTAURANT")
	a.user("c@x.test", 1, "RESTAURANT")
	a.device(1, "a@x.test", true)
	a.device(2, "a@x.test", true)
	a.device(3, "b@x.test", true)
	a.device(4, "c@x.test", false)

	got := a.send(map[string]any{"service": "RESTAURANT", "message": "Body", "navigation": map[string]any{"type": "restaurant"}})
	if got.status != http.StatusAccepted {
		t.Fatalf("status %d: %s", got.status, got.raw)
	}
	if got.body["queued"] != true || got.body["sent"] != false || got.body["devices"] != float64(3) || got.body["users"] != float64(3) {
		t.Fatalf("body = %s", got.raw)
	}
	if without := got.body["withoutDevice"].([]any); len(without) != 1 || without[0] != "c@x.test" {
		t.Fatalf("withoutDevice = %v", without)
	}
	if a.provider.count() != 0 {
		t.Fatal("the request sent a push itself")
	}

	if _, err := a.pipeline.RunOnce(context.Background()); err != nil {
		t.Fatal(err)
	}
	if a.provider.count() != 3 {
		t.Fatalf("%d pushes after the worker ran", a.provider.count())
	}

	id := int64(got.body["id"].(float64))
	detail := a.call(http.MethodGet, fmt.Sprintf("/admin/notifications/messages/%d", id), "admin@x.test", nil)
	if detail.status != 200 || detail.body["kind"] != "admin" || detail.body["createdBy"] != "admin@x.test" ||
		detail.body["deliveries"].(map[string]any)["accepted"] != float64(3) {
		t.Fatalf("detail = %s", detail.raw)
	}
	if strings.Contains(detail.raw, "PushToken") || strings.Contains(detail.raw, "Hello") {
		t.Fatalf("the diagnostics leak a token or the text: %s", detail.raw)
	}
	list := a.call(http.MethodGet, "/admin/notifications/messages?limit=5", "admin@x.test", nil)
	if list.status != 200 || len(list.body["messages"].([]any)) != 1 {
		t.Fatalf("list = %s", list.raw)
	}
}

func TestDryRunCountsOnlyDevicesThatQualify(t *testing.T) {
	a := newAPI(t, true)
	a.user("ok@x.test", 1)
	a.user("expired@x.test", 1)
	a.user("nodevice@x.test", 1)
	a.device(1, "ok@x.test", true)
	a.device(2, "expired@x.test", true)
	a.exec(`UPDATE devices SET last_seen_at = now() - interval '90 days' WHERE installation_id = $1`, "00000000-0000-4000-8000-000000000002")

	got := a.send(map[string]any{"dryRun": true, "audience": map[string]any{"type": "users", "emails": []string{"ok@x.test", "expired@x.test", "nodevice@x.test"}}})
	if got.status != 200 || got.body["devices"] != float64(1) || got.body["users"] != float64(3) {
		t.Fatalf("body = %s", got.raw)
	}
	filtered := got.body["filtered"].(map[string]any)
	if filtered["expired"] != float64(1) || filtered["no_device"] != float64(1) {
		t.Fatalf("filtered = %v", filtered)
	}
	if n := countRows(t, a.db, "push_messages"); n != 0 {
		t.Fatalf("a dry run queued %d messages", n)
	}
}

func TestNothingIsQueuedWhenNoDeviceQualifies(t *testing.T) {
	a := newAPI(t, true)
	a.user("a@x.test", 1)
	a.device(1, "a@x.test", false)
	got := a.send(map[string]any{})
	if got.status != http.StatusUnprocessableEntity {
		t.Fatalf("status %d: %s", got.status, got.raw)
	}
	if n := countRows(t, a.db, "push_messages"); n != 0 {
		t.Fatalf("%d messages queued", n)
	}
}

func TestAnImageAndTheNotificationTypeReachTheProvider(t *testing.T) {
	a := newAPI(t, true)
	a.user("a@x.test", 1, "EVENTS")
	a.device(1, "a@x.test", true)

	got := a.send(map[string]any{"service": "EVENTS", "message": "Body", "imageUrl": "https://cdn.x.test/poster.jpg",
		"navigation": map[string]any{"type": "event", "id": "7"}})
	if got.status != http.StatusAccepted {
		t.Fatalf("status %d: %s", got.status, got.raw)
	}
	if _, err := a.pipeline.RunOnce(context.Background()); err != nil {
		t.Fatal(err)
	}
	if a.provider.count() != 1 {
		t.Fatalf("%d pushes", a.provider.count())
	}
	sent := a.provider.sent[0]
	if sent.ImageURL != "https://cdn.x.test/poster.jpg" || sent.Data["type"] != "admin" {
		t.Fatalf("outbound = %+v", sent)
	}

	// Without an image the notification stays a plain text one.
	a.send(map[string]any{"service": "EVENTS", "message": "Again"})
	if _, err := a.pipeline.RunOnce(context.Background()); err != nil {
		t.Fatal(err)
	}
	if a.provider.count() != 2 || a.provider.sent[1].ImageURL != "" {
		t.Fatalf("outbound = %+v", a.provider.sent)
	}
}

func TestInvalidRequestsAreRefusedBeforeAnythingIsQueued(t *testing.T) {
	a := newAPI(t, true)
	a.user("a@x.test", 1)
	a.device(1, "a@x.test", true)
	for name, body := range map[string]map[string]any{
		"no title":         {"title": ""},
		"unknown audience": {"audience": map[string]any{"type": "nobody"}},
		"bad navigation":   {"navigation": map[string]any{"type": "nowhere"}},
		"unknown service":  {"service": "NOPE"},
		"http image":       {"imageUrl": "http://cdn.x.test/a.jpg"},
		"script image":     {"imageUrl": "javascript:alert(1)"},
		"relative image":   {"imageUrl": "/a.jpg"},
	} {
		if got := a.send(body); got.status != http.StatusBadRequest {
			t.Fatalf("%s: %d %s", name, got.status, got.raw)
		}
	}
	if n := countRows(t, a.db, "push_messages"); n != 0 {
		t.Fatalf("%d messages queued", n)
	}
}

func TestOnlyAdministratorsUseTheNotificationEndpoints(t *testing.T) {
	a := newAPI(t, true)
	for _, path := range []string{"/admin/notifications/messages", "/admin/notifications/messages/1"} {
		if got := a.call(http.MethodGet, path, "student@x.test", nil); got.status != http.StatusForbidden {
			t.Fatalf("%s as a student: %d", path, got.status)
		}
		if got := a.call(http.MethodGet, path, "", nil); got.status != http.StatusUnauthorized {
			t.Fatalf("%s without a token: %d", path, got.status)
		}
	}
	if got := a.call(http.MethodGet, "/admin/notifications/messages/9999", "admin@x.test", nil); got.status != http.StatusNotFound {
		t.Fatalf("unknown message: %d", got.status)
	}
	if got := a.call(http.MethodGet, "/admin/notifications/messages/abc", "admin@x.test", nil); got.status != http.StatusBadRequest {
		t.Fatalf("bad id: %d", got.status)
	}
}

func TestWithoutThePipelineTheDiagnosticsAreNotAvailable(t *testing.T) {
	a := newAPI(t, false)
	if got := a.call(http.MethodGet, "/admin/notifications/messages", "admin@x.test", nil); got.status != http.StatusNotFound {
		t.Fatalf("status %d", got.status)
	}
}

func countRows(t *testing.T, db *sql.DB, table string) int {
	t.Helper()
	var n int
	if err := db.QueryRow(`SELECT count(*) FROM ` + table).Scan(&n); err != nil {
		t.Fatal(err)
	}
	return n
}
