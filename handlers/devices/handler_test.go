package devices_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	devicesHandler "github.com/plugimt/transat-backend/handlers/devices"
	"github.com/plugimt/transat-backend/internal/database/dbtest"
	"github.com/plugimt/transat-backend/middlewares"
	"github.com/plugimt/transat-backend/routes"
	service "github.com/plugimt/transat-backend/services/devices"
	"github.com/plugimt/transat-backend/utils"
)

const (
	idA = "aaaaaaaa-aaaa-4aaa-8aaa-aaaaaaaaaaaa"
	idB = "bbbbbbbb-bbbb-4bbb-8bbb-bbbbbbbbbbbb"

	tok1 = "ExponentPushToken[token-one-aaaaaaaaaaaa]"
	tok2 = "ExponentPushToken[token-two-bbbbbbbbbbbb]"
)

type api struct {
	t   *testing.T
	db  *sql.DB
	svc *service.Service
	app *fiber.App
}

func newAPI(t *testing.T) *api {
	t.Setenv("JWT_SECRET", "test-secret-test-secret-test-secret")
	db := dbtest.Open(t)
	middlewares.InitJWTMiddleware(db)
	svc := service.NewService(db, service.DefaultConfig())

	// The same input filters as production, in the same order.
	app := fiber.New()
	app.Use(middlewares.ValidationMiddleware(), middlewares.SQLInjectionProtectionMiddleware(), middlewares.XSSProtectionMiddleware())
	routes.SetupDeviceRoutes(app, devicesHandler.NewHandler(svc))

	a := &api{t: t, db: db, svc: svc, app: app}
	for _, email := range []string{"alice@x.test", "bob@x.test"} {
		a.exec(`INSERT INTO newf (email, password, first_name, last_name, password_updated_date)
			VALUES ($1, 'x', 'A', 'B', now() - interval '3 years')`, email)
	}
	return a
}

func (a *api) exec(q string, args ...any) {
	a.t.Helper()
	if _, err := a.db.Exec(q, args...); err != nil {
		a.t.Fatalf("%s: %v", q, err)
	}
}

type reply struct {
	status int
	body   string
}

func (a *api) call(method, path, as string, body any) reply {
	a.t.Helper()
	var buf io.Reader
	if s, ok := body.(string); ok {
		buf = strings.NewReader(s)
	} else if body != nil {
		raw, _ := json.Marshal(body)
		buf = bytes.NewReader(raw)
	}
	req := httptest.NewRequest(method, path, buf)
	req.Header.Set("Content-Type", "application/json")
	if as != "" {
		token := as
		if strings.Contains(as, "@") {
			var err error
			if token, err = utils.GenerateJWT(as, nil, ""); err != nil {
				a.t.Fatal(err)
			}
		}
		req.Header.Set("Authorization", "Bearer "+token)
	}
	res, err := a.app.Test(req, -1)
	if err != nil {
		a.t.Fatal(err)
	}
	defer res.Body.Close()
	raw, _ := io.ReadAll(res.Body)
	return reply{res.StatusCode, string(raw)}
}

func body(token string) map[string]any {
	return map[string]any{"platform": "ios", "push_token": token, "permission_status": "granted", "app_version": "1.2.3"}
}

func (a *api) put(as, id string, b any) reply {
	return a.call(http.MethodPut, "/devices/"+id, as, b)
}

func (a *api) owner(id string) string {
	a.t.Helper()
	var email string
	if err := a.db.QueryRow(`SELECT user_email FROM devices WHERE installation_id = $1`, id).Scan(&email); err != nil {
		a.t.Fatal(err)
	}
	return email
}

type listing struct {
	Devices []map[string]any `json:"devices"`
	// Connected and ActiveSession are the two answers of "is this user here".
	Connected     bool `json:"connected"`
	ActiveSession bool `json:"has_active_session"`
}

func (a *api) list(as string) listing {
	a.t.Helper()
	r := a.call(http.MethodGet, "/devices", as, nil)
	if r.status != 200 {
		a.t.Fatalf("GET /devices: %d %s", r.status, r.body)
	}
	var l listing
	if err := json.Unmarshal([]byte(r.body), &l); err != nil {
		a.t.Fatal(err)
	}
	return l
}

func TestEndpointsRequireAValidSessionExceptSignOut(t *testing.T) {
	a := newAPI(t)
	for _, as := range []string{"", "garbage"} {
		if got := a.put(as, idA, body(tok1)); got.status != http.StatusUnauthorized {
			t.Fatalf("PUT as %q: %d", as, got.status)
		}
		if got := a.call(http.MethodGet, "/devices", as, nil); got.status != http.StatusUnauthorized {
			t.Fatalf("GET as %q: %d", as, got.status)
		}
	}
	var n int
	_ = a.db.QueryRow(`SELECT count(*) FROM devices`).Scan(&n)
	if n != 0 {
		t.Fatalf("%d devices stored without a session", n)
	}
}

func TestTheOwnerIsTheAuthenticatedUser(t *testing.T) {
	a := newAPI(t)
	b := body(tok1)
	b["email"] = "bob@x.test" // ignored: the owner never comes from the body
	b["user_email"] = "bob@x.test"
	got := a.put("alice@x.test", idA, b)
	if got.status != 200 {
		t.Fatalf("PUT: %d %s", got.status, got.body)
	}
	if a.owner(idA) != "alice@x.test" {
		t.Fatalf("owner = %s", a.owner(idA))
	}

	var res map[string]any
	_ = json.Unmarshal([]byte(got.body), &res)
	if res["heartbeat_interval_seconds"] != float64(300) || res["presence_window_seconds"] != float64(600) {
		t.Fatalf("response = %s", got.body)
	}
}

func TestInvalidRegistrationsAreRejected(t *testing.T) {
	a := newAPI(t)
	cases := map[string]struct {
		id   string
		body any
	}{
		"bad installation id": {"abc", body(tok1)},
		"no platform":         {idA, map[string]any{"push_token": tok1}},
		"unknown platform":    {idA, map[string]any{"platform": "web"}},
		"bad token":           {idA, body("not a token")},
		"bad permission":      {idA, map[string]any{"platform": "ios", "permission_status": "yes"}},
		"not json":            {idA, "{nope"},
		"wrong types":         {idA, map[string]any{"platform": 3}},
	}
	for name, c := range cases {
		if got := a.put("alice@x.test", c.id, c.body); got.status != http.StatusBadRequest {
			t.Fatalf("%s: %d %s", name, got.status, got.body)
		}
	}
	var n int
	_ = a.db.QueryRow(`SELECT count(*) FROM devices`).Scan(&n)
	if n != 0 {
		t.Fatalf("%d devices stored from invalid requests", n)
	}
}

// Expo tokens can contain "--", which the SQL-comment filter used to reject: those users never registered.
func TestATokenWithDoubleDashPassesTheInputFilters(t *testing.T) {
	a := newAPI(t)
	if got := a.put("alice@x.test", idA, body("ExponentPushToken[ab--cd_ef]")); got.status != 200 {
		t.Fatalf("PUT: %d %s", got.status, got.body)
	}
	// The filter still protects every other body.
	if got := a.call(http.MethodPost, "/devices/"+idA+"/sign-out", "", `{"x":"1; DROP TABLE devices"}`); got.status != http.StatusNoContent {
		t.Fatalf("sign-out: %d", got.status)
	}
}

func TestListShowsOnlyTheCallersDevicesWithoutTokens(t *testing.T) {
	a := newAPI(t)
	a.put("alice@x.test", idA, body(tok1))
	a.put("bob@x.test", idB, body(tok2))

	l := a.list("alice@x.test")
	if len(l.Devices) != 1 || l.Devices[0]["installation_id"] != idA || !l.Connected || !l.ActiveSession {
		t.Fatalf("listing = %+v", l)
	}
	raw := a.call(http.MethodGet, "/devices", "alice@x.test", nil).body
	if strings.Contains(raw, "ExponentPushToken") || strings.Contains(raw, `"push_token"`) {
		t.Fatalf("the push token leaked: %s", raw)
	}
	if l.Devices[0]["state"] != "connected" || l.Devices[0]["has_push_token"] != true {
		t.Fatalf("device = %v", l.Devices[0])
	}
}

func TestZeroDevices(t *testing.T) {
	a := newAPI(t)
	l := a.list("alice@x.test")
	if l.Devices == nil || len(l.Devices) != 0 || l.Connected || l.ActiveSession {
		t.Fatalf("listing = %+v", l)
	}
}

func TestLogoutThenAccountSwitch(t *testing.T) {
	a := newAPI(t)
	a.put("alice@x.test", idA, body(tok1))

	// Alice logs out: no token at all is needed.
	if got := a.call(http.MethodPost, "/devices/"+idA+"/sign-out", "", nil); got.status != http.StatusNoContent {
		t.Fatalf("sign-out: %d %s", got.status, got.body)
	}
	l := a.list("alice@x.test")
	if l.Connected || l.ActiveSession || l.Devices[0]["state"] != "disabled" || l.Devices[0]["has_push_token"] != false {
		t.Fatalf("alice after logout = %+v", l)
	}
	var legacy int
	_ = a.db.QueryRow(`SELECT count(*) FROM user_notification_tokens`).Scan(&legacy)
	if legacy != 0 {
		t.Fatalf("%d legacy tokens left after logout", legacy)
	}

	// Bob logs in on the same phone.
	if got := a.put("bob@x.test", idA, body(tok1)); got.status != 200 {
		t.Fatalf("bob PUT: %d %s", got.status, got.body)
	}
	if a.owner(idA) != "bob@x.test" {
		t.Fatalf("owner = %s", a.owner(idA))
	}
	if l := a.list("alice@x.test"); l.Connected || len(l.Devices) != 0 {
		t.Fatalf("alice must not see bob's phone: %+v", l)
	}
	if l := a.list("bob@x.test"); !l.Connected || len(l.Devices) != 1 {
		t.Fatalf("bob = %+v", l)
	}
}

func TestSignOutIsLimitedPerInstallationNotPerAddress(t *testing.T) {
	a := newAPI(t)
	// Behind the proxy every user shares an address, and logins already count against the account limiter:
	// a burst of sign-outs from different phones must still all go through.
	for i := 0; i < 80; i++ {
		id := fmt.Sprintf("dddddddd-dddd-4ddd-8ddd-%012d", i)
		if got := a.call(http.MethodPost, "/devices/"+id+"/sign-out", "", nil); got.status != http.StatusNoContent {
			t.Fatalf("sign-out %d: %d %s", i, got.status, got.body)
		}
	}
	// One installation cannot be hammered.
	const idC = "cccccccc-cccc-4ccc-8ccc-cccccccccccc"
	limited := 0
	for i := 0; i < 40; i++ {
		if got := a.call(http.MethodPost, "/devices/"+idC+"/sign-out", "", nil); got.status == http.StatusTooManyRequests {
			limited++
		}
	}
	if limited == 0 {
		t.Fatal("an installation can be signed out without limit")
	}
}

func TestSignOutIsIdempotentAndHidesWhatExists(t *testing.T) {
	a := newAPI(t)
	a.put("alice@x.test", idA, body(tok1))
	for i := 0; i < 2; i++ {
		if got := a.call(http.MethodPost, "/devices/"+idA+"/sign-out", "", nil); got.status != http.StatusNoContent {
			t.Fatalf("known id, call %d: %d", i, got.status)
		}
	}
	if got := a.call(http.MethodPost, "/devices/"+idB+"/sign-out", "", nil); got.status != http.StatusNoContent {
		t.Fatalf("unknown id: %d", got.status)
	}
	if got := a.call(http.MethodPost, "/devices/nope/sign-out", "", nil); got.status != http.StatusBadRequest {
		t.Fatalf("malformed id: %d", got.status)
	}
}

func TestARevokedSessionCannotHeartbeatNorBeActive(t *testing.T) {
	a := newAPI(t)
	a.put("alice@x.test", idA, body(tok1))
	token, _ := utils.GenerateJWT("alice@x.test", nil, "")

	// The password changes: the token issued before it is refused, so no heartbeat can keep the device alive.
	a.exec(`UPDATE newf SET password_updated_date = now() + interval '3 hours' WHERE email = 'alice@x.test'`)
	if got := a.put(token, idA, body(tok1)); got.status != http.StatusUnauthorized {
		t.Fatalf("heartbeat with a revoked token: %d %s", got.status, got.body)
	}
	devices, err := a.svc.List(t.Context(), "alice@x.test")
	if err != nil || len(devices) != 1 || devices[0].State != service.StateRevoked {
		t.Fatalf("devices = %+v, %v", devices, err)
	}
	if service.AnyActiveSession(devices) || service.AnyConnected(devices) {
		t.Fatal("a revoked device is neither connected nor an active session")
	}
}

func TestADatabaseFailureIsNotAnAcceptedRegistration(t *testing.T) {
	a := newAPI(t)
	// A service on an empty schema: the authentication works, the write cannot.
	broken := service.NewService(dbtest.Open(t), service.DefaultConfig())
	app := fiber.New()
	routes.SetupDeviceRoutes(app, devicesHandler.NewHandler(broken))
	a.app = app

	if got := a.put("alice@x.test", idA, body(tok1)); got.status != http.StatusInternalServerError {
		t.Fatalf("PUT = %d %s, want 500 so the app retries", got.status, got.body)
	}
}
