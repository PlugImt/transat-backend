package admin_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"sort"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	devicesHandler "github.com/plugimt/transat-backend/handlers/devices"
	"github.com/plugimt/transat-backend/internal/database/dbtest"
	"github.com/plugimt/transat-backend/middlewares"
	"github.com/plugimt/transat-backend/routes"
	"github.com/plugimt/transat-backend/services"
	"github.com/plugimt/transat-backend/services/devices"
	"github.com/plugimt/transat-backend/services/push"
)

// lifecycle drives the whole server side of a notification through HTTP: apps register their devices, an
// administrator sends, the worker dispatches to a stand-in for the Expo push service that records every request.
type lifecycle struct {
	*api
	mu     sync.Mutex
	sent   []map[string]any
	answer func(call int, msgs []map[string]any) (int, string)
	calls  int
}

func newLifecycle(t *testing.T) *lifecycle {
	t.Setenv("JWT_SECRET", "test-secret-test-secret-test-secret")
	db := dbtest.Open(t)
	middlewares.InitJWTMiddleware(db)

	l := &lifecycle{api: &api{t: t, db: db}}
	expo := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		raw, _ := io.ReadAll(r.Body)
		var msgs []map[string]any
		_ = json.Unmarshal(raw, &msgs)
		l.mu.Lock()
		l.calls++
		call := l.calls
		l.sent = append(l.sent, msgs...)
		answer := l.answer
		l.mu.Unlock()

		status, body := http.StatusOK, ""
		if answer != nil {
			status, body = answer(call, msgs)
		}
		if body == "" {
			tickets := make([]string, len(msgs))
			for i := range msgs {
				tickets[i] = fmt.Sprintf(`{"status":"ok","id":"t-%d-%d"}`, call, i)
			}
			body = `{"data":[` + strings.Join(tickets, ",") + `]}`
		}
		w.WriteHeader(status)
		_, _ = w.Write([]byte(body))
	}))
	t.Cleanup(expo.Close)

	cfg := push.DefaultConfig()
	cfg.MaxPerSecond = 0
	cfg.RetryBase, cfg.RetryMax = time.Hour, 2*time.Hour
	provider := push.NewExpoProvider(push.ExpoConfig{SendURL: expo.URL, ReceiptsURL: expo.URL})
	l.pipeline = push.New(db, provider, cfg, devices.DefaultConfig())

	ns := services.NewNotificationService(db)
	ns.SetPipeline(l.pipeline)
	l.app = fiber.New()
	routes.SetupAdminRoutes(l.app, db, ns)
	routes.SetupDeviceRoutes(l.app, devicesHandler.NewHandler(devices.NewService(db, devices.DefaultConfig())))

	l.exec(`INSERT INTO roles (name) VALUES ('ADMIN') ON CONFLICT (name) DO NOTHING`)
	l.user("admin@x.test", 1)
	l.exec(`INSERT INTO newf_roles (email, id_roles) SELECT 'admin@x.test', id_roles FROM roles WHERE name = 'ADMIN'`)
	return l
}

func installation(n int) string { return fmt.Sprintf("00000000-0000-4000-8000-%012d", n) }
func token(name string) string  { return "ExponentPushToken[" + name + "]" }

// register is what the app does at sign-in and on every heartbeat.
func (l *lifecycle) register(n int, email, platform, pushToken string) {
	l.t.Helper()
	body := map[string]any{"platform": platform, "permission_status": "granted", "app_version": "2.4.0"}
	if pushToken != "" {
		body["push_token"] = pushToken
	}
	if got := l.call(http.MethodPut, "/devices/"+installation(n), email, body); got.status != http.StatusOK {
		l.t.Fatalf("register: %d %s", got.status, got.raw)
	}
}

func (l *lifecycle) signOut(n int) {
	l.t.Helper()
	if got := l.call(http.MethodPost, "/devices/"+installation(n)+"/sign-out", "", nil); got.status != http.StatusNoContent {
		l.t.Fatalf("sign-out: %d %s", got.status, got.raw)
	}
}

// notify sends a Restaurant notification from the admin API, runs the worker once and returns what Expo got.
func (l *lifecycle) notify(extra map[string]any) []map[string]any {
	l.t.Helper()
	l.mu.Lock()
	before := len(l.sent)
	l.mu.Unlock()
	body := map[string]any{"service": "RESTAURANT", "message": "Menu", "navigation": map[string]any{"type": "restaurant"}}
	for k, v := range extra {
		body[k] = v
	}
	if got := l.send(body); got.status != http.StatusAccepted {
		l.t.Fatalf("send: %d %s", got.status, got.raw)
	}
	if _, err := l.pipeline.RunOnce(context.Background()); err != nil {
		l.t.Fatal(err)
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	return append([]map[string]any(nil), l.sent[before:]...)
}

func tokensOf(msgs []map[string]any) []string {
	out := make([]string, len(msgs))
	for i, m := range msgs {
		out[i], _ = m["to"].(string)
	}
	sort.Strings(out)
	return out
}

func sameTokens(t *testing.T, got []map[string]any, want ...string) {
	t.Helper()
	sort.Strings(want)
	if g := tokensOf(got); fmt.Sprint(g) != fmt.Sprint(want) {
		t.Fatalf("Expo was asked to push to %v, want %v", g, want)
	}
}

func (l *lifecycle) deviceState(email string, n int) string {
	l.t.Helper()
	got := l.call(http.MethodGet, "/devices", email, nil)
	for _, d := range got.body["devices"].([]any) {
		if m := d.(map[string]any); m["installation_id"] == installation(n) {
			return m["state"].(string)
		}
	}
	return "none"
}

func TestTheWholeLifecycleFromSignInToAnInvalidToken(t *testing.T) {
	l := newLifecycle(t)
	l.user("alice@x.test", 1, "RESTAURANT")
	l.user("bob@x.test", 1, "RESTAURANT")
	l.user("carol@x.test", 1) // opted out of the category

	// Several devices of one user register independently; another user registers one; a third has none.
	l.register(1, "alice@x.test", "ios", token("alice-phone"))
	l.register(2, "alice@x.test", "android", token("alice-tablet"))
	l.register(3, "bob@x.test", "ios", token("bob-phone"))
	l.register(4, "carol@x.test", "ios", token("carol-phone"))
	if got := l.call(http.MethodGet, "/devices", "alice@x.test", nil); len(got.body["devices"].([]any)) != 2 {
		t.Fatalf("alice sees %s", got.raw)
	}

	// 1. Every device of every opted-in user is reached once, with the destination and the image.
	first := l.notify(map[string]any{"imageUrl": "https://cdn.x.test/menu.jpg"})
	sameTokens(t, first, token("alice-phone"), token("alice-tablet"), token("bob-phone"))
	for _, m := range first {
		data := m["data"].(map[string]any)
		if data["type"] != "admin" || data["navigation"].(map[string]any)["type"] != "restaurant" {
			t.Fatalf("data = %v", data)
		}
		if m["richContent"].(map[string]any)["image"] != "https://cdn.x.test/menu.jpg" {
			t.Fatalf("image missing: %v", m)
		}
		if (m["to"] == token("alice-phone") || m["to"] == token("bob-phone")) != (m["mutableContent"] == true) {
			t.Fatalf("mutableContent must be set on iOS only: %v", m)
		}
	}

	// 2. A token refresh moves the installation to the new token: the old one is never pushed to again.
	l.register(1, "alice@x.test", "ios", token("alice-phone-2"))
	sameTokens(t, l.notify(nil), token("alice-phone-2"), token("alice-tablet"), token("bob-phone"))

	// 3. Signing out one device leaves the others alone.
	l.signOut(2)
	sameTokens(t, l.notify(nil), token("alice-phone-2"), token("bob-phone"))
	if state := l.deviceState("alice@x.test", 2); state != "disabled" {
		t.Fatalf("signed-out device is %s", state)
	}

	// 4. Account switch: bob signs in on alice's phone. Alice no longer gets pushes there, bob gets them once.
	l.register(1, "bob@x.test", "ios", token("alice-phone-2"))
	sameTokens(t, l.notify(nil), token("alice-phone-2"), token("bob-phone"))
	if got := l.call(http.MethodGet, "/devices", "alice@x.test", nil); len(got.body["devices"].([]any)) != 1 {
		t.Fatalf("alice still sees the phone that was handed over: %s", got.raw)
	}

	// 5. The provider reports a dead token: that device is deactivated, the others are untouched and later
	// notifications skip it.
	l.mu.Lock()
	l.answer = func(_ int, msgs []map[string]any) (int, string) {
		tickets := make([]string, len(msgs))
		for i, m := range msgs {
			if m["to"] == token("bob-phone") {
				tickets[i] = `{"status":"error","message":"not registered","details":{"error":"DeviceNotRegistered"}}`
			} else {
				tickets[i] = fmt.Sprintf(`{"status":"ok","id":"ok-%d"}`, i)
			}
		}
		return http.StatusOK, `{"data":[` + strings.Join(tickets, ",") + `]}`
	}
	l.mu.Unlock()
	l.notify(nil)
	l.mu.Lock()
	l.answer = nil
	l.mu.Unlock()
	if state := l.deviceState("bob@x.test", 3); state != "disabled" {
		t.Fatalf("a dead token left the device %s", state)
	}
	sameTokens(t, l.notify(nil), token("alice-phone-2"))

	// 6. The app registers the device again with a working token: it is reachable again.
	l.register(3, "bob@x.test", "ios", token("bob-phone-2"))
	sameTokens(t, l.notify(nil), token("alice-phone-2"), token("bob-phone-2"))
}

func TestATransientProviderFailureIsRetriedWithoutDuplicates(t *testing.T) {
	l := newLifecycle(t)
	l.user("alice@x.test", 1, "RESTAURANT")
	l.register(1, "alice@x.test", "ios", token("alice-phone"))
	l.register(2, "alice@x.test", "android", token("alice-tablet"))

	// Expo is down for the first request, then recovers.
	l.mu.Lock()
	l.answer = func(call int, _ []map[string]any) (int, string) {
		if call == 1 {
			return http.StatusServiceUnavailable, "down"
		}
		return http.StatusOK, ""
	}
	l.mu.Unlock()

	first := l.notify(nil)
	if len(first) != 2 {
		t.Fatalf("first attempt sent %d", len(first))
	}
	if n := countRows(t, l.db, "push_deliveries WHERE status = 'accepted'"); n != 0 {
		t.Fatalf("%d deliveries were accepted while Expo was down", n)
	}

	l.exec(`UPDATE push_deliveries SET next_attempt_at = now()`)
	if _, err := l.pipeline.RunOnce(context.Background()); err != nil {
		t.Fatal(err)
	}
	if n := countRows(t, l.db, "push_deliveries WHERE status = 'accepted'"); n != 2 {
		t.Fatalf("%d deliveries accepted after the retry, want one per device", n)
	}

	// Nothing is sent again by later runs, and the retried push carries the same collapse key as the first one,
	// so a duplicate that did get through would replace the notification on the phone instead of stacking.
	l.exec(`UPDATE push_deliveries SET next_attempt_at = now()`)
	if _, err := l.pipeline.RunOnce(context.Background()); err != nil {
		t.Fatal(err)
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	if len(l.sent) != 4 {
		t.Fatalf("%d provider messages in all, want 2 failed + 2 retried", len(l.sent))
	}
	keys := map[string]bool{}
	for _, m := range l.sent {
		key, _ := m["collapseId"].(string)
		if key == "" {
			key, _ = m["tag"].(string)
		}
		keys[key] = true
	}
	if len(keys) != 1 || keys[""] {
		t.Fatalf("collapse keys = %v", keys)
	}
}

func TestAnExpiredSessionIsNotPushedToAndAHeartbeatRestoresIt(t *testing.T) {
	l := newLifecycle(t)
	l.user("alice@x.test", 1, "RESTAURANT")
	l.register(1, "alice@x.test", "ios", token("alice-phone"))
	l.exec(`UPDATE devices SET last_seen_at = now() - interval '90 days'`)
	got := l.send(map[string]any{"service": "RESTAURANT", "message": "Menu"})
	if got.status != http.StatusUnprocessableEntity {
		t.Fatalf("an audience with no eligible device: %d %s", got.status, got.raw)
	}
	l.mu.Lock()
	pushed := len(l.sent)
	l.mu.Unlock()
	if pushed != 0 {
		t.Fatalf("an expired session was pushed to")
	}
	l.register(1, "alice@x.test", "ios", token("alice-phone"))
	sameTokens(t, l.notify(nil), token("alice-phone"))
}
