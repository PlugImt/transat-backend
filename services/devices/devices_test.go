package devices

import (
	"context"
	"database/sql"
	"fmt"
	"sort"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/plugimt/transat-backend/internal/database/dbtest"
)

const (
	idA = "aaaaaaaa-aaaa-4aaa-8aaa-aaaaaaaaaaaa"
	idB = "bbbbbbbb-bbbb-4bbb-8bbb-bbbbbbbbbbbb"
	idC = "cccccccc-cccc-4ccc-8ccc-cccccccccccc"

	tok1 = "ExponentPushToken[token-one-aaaaaaaaaaaa]"
	tok2 = "ExponentPushToken[token-two-bbbbbbbbbbbb]"
	tok3 = "ExponentPushToken[token-three-cccccccccc]"
)

type fx struct {
	t   *testing.T
	db  *sql.DB
	svc *Service
	ctx context.Context
}

func newFx(t *testing.T) *fx {
	db := dbtest.Open(t)
	f := &fx{t: t, db: db, svc: NewService(db, DefaultConfig()), ctx: context.Background()}
	f.user("alice@x.test")
	f.user("bob@x.test")
	return f
}

func (f *fx) exec(q string, args ...any) {
	f.t.Helper()
	if _, err := f.db.Exec(q, args...); err != nil {
		f.t.Fatalf("%s: %v", q, err)
	}
}

func (f *fx) user(email string) {
	// Accounts get their password date at creation: make them old so moving a device into the past is possible.
	f.exec(`INSERT INTO newf (email, password, first_name, last_name, password_updated_date)
		VALUES ($1, 'x', 'A', 'B', now() - interval '3 years')`, email)
}

func (f *fx) register(id, email, token string) bool {
	f.t.Helper()
	written, err := f.svc.Register(f.ctx, Registration{
		InstallationID: id, Email: email, Platform: "ios", PushToken: token, PermissionStatus: "granted", AppVersion: "1.0.0",
	})
	if err != nil {
		f.t.Fatalf("register %s for %s: %v", id, email, err)
	}
	return written
}

// seen moves the last proof of a device into the past.
func (f *fx) seen(id string, ago time.Duration) {
	f.t.Helper()
	f.exec(`UPDATE devices SET last_seen_at = now() - make_interval(secs => $2) WHERE installation_id = $1`, id, ago.Seconds())
}

func (f *fx) device(email, id string) (Device, bool) {
	f.t.Helper()
	devices, err := f.svc.List(f.ctx, email)
	if err != nil {
		f.t.Fatal(err)
	}
	for _, d := range devices {
		if d.InstallationID == id {
			return d, true
		}
	}
	return Device{}, false
}

func (f *fx) state(email, id string) State {
	f.t.Helper()
	d, ok := f.device(email, id)
	if !ok {
		f.t.Fatalf("%s has no device %s", email, id)
	}
	return d.State
}

func (f *fx) token(id string) (string, bool) {
	f.t.Helper()
	var tok sql.NullString
	if err := f.db.QueryRow(`SELECT push_token FROM devices WHERE installation_id = $1`, id).Scan(&tok); err != nil {
		f.t.Fatal(err)
	}
	return tok.String, tok.Valid
}

func (f *fx) count(q string, args ...any) int {
	f.t.Helper()
	var n int
	if err := f.db.QueryRow(q, args...).Scan(&n); err != nil {
		f.t.Fatal(err)
	}
	return n
}

// legacy returns the tokens the current dispatch would read for a user.
func (f *fx) legacy(email string) []string {
	f.t.Helper()
	rows, err := f.db.Query(`SELECT token FROM user_notification_tokens WHERE email = $1 ORDER BY token`, email)
	if err != nil {
		f.t.Fatal(err)
	}
	defer rows.Close()
	var out []string
	for rows.Next() {
		var s string
		_ = rows.Scan(&s)
		out = append(out, s)
	}
	return out
}

func same(a, b []string) bool {
	sort.Strings(a)
	sort.Strings(b)
	return strings.Join(a, ",") == strings.Join(b, ",")
}

func TestUserWithoutDevicesIsNotConnected(t *testing.T) {
	f := newFx(t)
	devices, err := f.svc.List(f.ctx, "alice@x.test")
	if err != nil || devices == nil || len(devices) != 0 {
		t.Fatalf("List = %v, %v; want an empty, non-nil list", devices, err)
	}
	if ok, err := f.svc.IsConnected(f.ctx, "alice@x.test"); err != nil || ok {
		t.Fatalf("IsConnected = %v, %v", ok, err)
	}
	if ok, err := f.svc.HasActiveSession(f.ctx, "alice@x.test"); err != nil || ok {
		t.Fatalf("HasActiveSession = %v, %v", ok, err)
	}
}

func TestManyDevicesPerUserKeepTheirOwnTokens(t *testing.T) {
	f := newFx(t)
	f.register(idA, "alice@x.test", tok1)
	f.register(idB, "alice@x.test", tok2)
	f.register(idC, "alice@x.test", "") // a third installation without token

	if n := f.count(`SELECT count(*) FROM devices WHERE user_email = 'alice@x.test'`); n != 3 {
		t.Fatalf("%d devices, want 3", n)
	}
	if got, _ := f.token(idA); got != tok1 {
		t.Fatalf("device A token = %q: another device overwrote it", got)
	}
	if got, _ := f.token(idB); got != tok2 {
		t.Fatalf("device B token = %q", got)
	}
	if _, has := f.token(idC); has {
		t.Fatal("device C must have no token")
	}
	if got := f.legacy("alice@x.test"); !same(got, []string{tok1, tok2}) {
		t.Fatalf("legacy tokens = %v", got)
	}
	if ok, _ := f.svc.IsConnected(f.ctx, "alice@x.test"); !ok {
		t.Fatal("alice should be connected")
	}
}

func TestDuplicateRegistrationIsOneRow(t *testing.T) {
	f := newFx(t)
	if !f.register(idA, "alice@x.test", tok1) {
		t.Fatal("the first registration must write")
	}
	var created, updated time.Time
	_ = f.db.QueryRow(`SELECT created_at, updated_at FROM devices WHERE installation_id = $1`, idA).Scan(&created, &updated)

	if f.register(idA, "alice@x.test", tok1) {
		t.Fatal("an identical registration right after the first must not write")
	}
	if n := f.count(`SELECT count(*) FROM devices`); n != 1 {
		t.Fatalf("%d rows after a duplicate registration", n)
	}
	if n := f.count(`SELECT count(*) FROM user_notification_tokens WHERE email = 'alice@x.test'`); n != 1 {
		t.Fatalf("%d legacy rows", n)
	}
	var created2, updated2 time.Time
	_ = f.db.QueryRow(`SELECT created_at, updated_at FROM devices WHERE installation_id = $1`, idA).Scan(&created2, &updated2)
	if !created.Equal(created2) || !updated.Equal(updated2) {
		t.Fatal("a duplicate registration changed created_at or updated_at")
	}
}

func TestHeartbeatRefreshesLastSeenOnlyAfterTheThrottle(t *testing.T) {
	f := newFx(t)
	f.register(idA, "alice@x.test", tok1)

	f.seen(idA, 20*time.Second)
	if f.register(idA, "alice@x.test", tok1) {
		t.Fatal("a heartbeat 20s after the last one must be skipped")
	}
	if f.state("alice@x.test", idA) != StateConnected {
		t.Fatal("still connected")
	}

	f.seen(idA, 3*time.Minute)
	var updatedBefore time.Time
	_ = f.db.QueryRow(`SELECT updated_at FROM devices WHERE installation_id = $1`, idA).Scan(&updatedBefore)
	if !f.register(idA, "alice@x.test", tok1) {
		t.Fatal("a heartbeat 3 minutes after the last one must write")
	}
	d, _ := f.device("alice@x.test", idA)
	if time.Since(d.LastSeenAt) > 10*time.Second {
		t.Fatalf("last_seen_at not refreshed: %v", d.LastSeenAt)
	}
	if !d.UpdatedAt.Equal(updatedBefore) {
		t.Fatal("a plain heartbeat must not touch updated_at")
	}
}

func TestTokenRefreshReplacesTheTokenOfTheSameInstallation(t *testing.T) {
	f := newFx(t)
	f.register(idA, "alice@x.test", tok1)
	f.register(idA, "alice@x.test", tok2)

	if n := f.count(`SELECT count(*) FROM devices`); n != 1 {
		t.Fatalf("%d rows after a token refresh", n)
	}
	if got, _ := f.token(idA); got != tok2 {
		t.Fatalf("token = %q, want the refreshed one", got)
	}
	if got := f.legacy("alice@x.test"); !same(got, []string{tok2}) {
		t.Fatalf("legacy tokens = %v, the old token must be gone", got)
	}
}

func TestMissingTokenAndRefusedPermissionAreHandled(t *testing.T) {
	f := newFx(t)
	f.register(idA, "alice@x.test", tok1)

	_, err := f.svc.Register(f.ctx, Registration{
		InstallationID: idA, Email: "alice@x.test", Platform: "android", PermissionStatus: "denied",
	})
	if err != nil {
		t.Fatal(err)
	}
	d, _ := f.device("alice@x.test", idA)
	if d.HasPushToken || d.PermissionStatus != "denied" || d.Platform != "android" {
		t.Fatalf("device = %+v", d)
	}
	if d.State != StateConnected {
		t.Fatalf("a device without permission is still connected, got %s", d.State)
	}
	if got := f.legacy("alice@x.test"); len(got) != 0 {
		t.Fatalf("legacy tokens = %v: a device that cannot show notifications must not keep a token", got)
	}

	// The permission is granted later: the same installation gets its token back.
	f.register(idA, "alice@x.test", tok3)
	if got, _ := f.token(idA); got != tok3 {
		t.Fatalf("token = %q", got)
	}
}

func TestATokenNeverBelongsToTwoInstallations(t *testing.T) {
	f := newFx(t)
	f.register(idA, "alice@x.test", tok1)
	// The app was reinstalled: a new installation id arrives with the token the old one still holds.
	f.register(idB, "alice@x.test", tok1)

	if got, _ := f.token(idB); got != tok1 {
		t.Fatalf("new installation token = %q", got)
	}
	old, _ := f.device("alice@x.test", idA)
	if old.HasPushToken || old.State != StateDisabled || old.DisabledReason != "token_replaced" {
		t.Fatalf("old installation = %+v", old)
	}
	if n := f.count(`SELECT count(*) FROM devices WHERE push_token = $1`, tok1); n != 1 {
		t.Fatalf("%d installations hold the token", n)
	}
	if got := f.legacy("alice@x.test"); !same(got, []string{tok1}) {
		t.Fatalf("legacy tokens = %v", got)
	}

	// The stale installation checks in again with a token of its own and comes back.
	f.register(idA, "alice@x.test", tok2)
	if f.state("alice@x.test", idA) != StateConnected {
		t.Fatal("the installation must be reactivated by its next registration")
	}
}

func TestAccountSwitchMovesTheDeviceToTheNewAccount(t *testing.T) {
	f := newFx(t)
	f.register(idA, "alice@x.test", tok1)

	// Bob signs in on the same phone: same installation id, the token may be the same or refreshed.
	f.register(idA, "bob@x.test", tok1)

	if _, ok := f.device("alice@x.test", idA); ok {
		t.Fatal("alice must no longer own the device")
	}
	if d, ok := f.device("bob@x.test", idA); !ok || d.State != StateConnected {
		t.Fatalf("bob's device = %+v, %v", d, ok)
	}
	if got := f.legacy("alice@x.test"); len(got) != 0 {
		t.Fatalf("alice still has legacy tokens %v: she would keep receiving bob's phone", got)
	}
	if got := f.legacy("bob@x.test"); !same(got, []string{tok1}) {
		t.Fatalf("bob legacy tokens = %v", got)
	}
	if ok, _ := f.svc.IsConnected(f.ctx, "alice@x.test"); ok {
		t.Fatal("alice is no longer connected")
	}

	// Back to alice with a refreshed token.
	f.register(idA, "alice@x.test", tok2)
	if got := f.legacy("bob@x.test"); len(got) != 0 {
		t.Fatalf("bob legacy tokens = %v after alice came back", got)
	}
	if got := f.legacy("alice@x.test"); !same(got, []string{tok2}) {
		t.Fatalf("alice legacy tokens = %v", got)
	}
}

func TestLegacyTokenOfAPreviousAccountIsRemovedWhenTheTokenIsRegistered(t *testing.T) {
	f := newFx(t)
	// Written by an old app version: no installation row, the token sits under alice.
	f.exec(`INSERT INTO user_notification_tokens (email, token) VALUES ('alice@x.test', $1)`, tok1)
	f.register(idA, "bob@x.test", tok1)

	if got := f.legacy("alice@x.test"); len(got) != 0 {
		t.Fatalf("alice legacy tokens = %v", got)
	}
}

func TestSignOutDisablesTheDeviceAndItsToken(t *testing.T) {
	f := newFx(t)
	f.register(idA, "alice@x.test", tok1)
	f.register(idB, "alice@x.test", tok2)

	ok, err := f.svc.SignOut(f.ctx, strings.ToUpper(idA))
	if err != nil || !ok {
		t.Fatalf("SignOut = %v, %v", ok, err)
	}
	d, _ := f.device("alice@x.test", idA)
	if d.State != StateDisabled || d.DisabledReason != "signed_out" || d.HasPushToken {
		t.Fatalf("device = %+v", d)
	}
	if got := f.legacy("alice@x.test"); !same(got, []string{tok2}) {
		t.Fatalf("legacy tokens = %v, only the other phone must remain", got)
	}
	if f.state("alice@x.test", idB) != StateConnected {
		t.Fatal("signing out one device must not touch the others")
	}

	if ok, err := f.svc.SignOut(f.ctx, idA); err != nil || ok {
		t.Fatalf("second SignOut = %v, %v; want a no-op", ok, err)
	}
	if ok, err := f.svc.SignOut(f.ctx, idC); err != nil || ok {
		t.Fatalf("SignOut of an unknown installation = %v, %v", ok, err)
	}
	if _, err := f.svc.SignOut(f.ctx, "not-a-uuid"); err != ErrInvalidInstallationID {
		t.Fatalf("SignOut(invalid) error = %v", err)
	}

	f.exec(`DELETE FROM devices WHERE installation_id = $1`, idB)
	if ok, _ := f.svc.IsConnected(f.ctx, "alice@x.test"); ok {
		t.Fatal("a signed-out device must not make the user connected")
	}
}

func TestNextLoginReactivatesASignedOutDevice(t *testing.T) {
	f := newFx(t)
	f.register(idA, "alice@x.test", tok1)
	if _, err := f.svc.SignOut(f.ctx, idA); err != nil {
		t.Fatal(err)
	}
	f.register(idA, "bob@x.test", tok1)

	d, ok := f.device("bob@x.test", idA)
	if !ok || d.State != StateConnected || d.DisabledReason != "" || !d.HasPushToken {
		t.Fatalf("device = %+v, %v", d, ok)
	}
}

func TestStateFollowsTheLastProofOfSession(t *testing.T) {
	cfg := DefaultConfig()
	cases := map[string]struct {
		ago  time.Duration
		want State
	}{
		"just now":                 {0, StateConnected},
		"inside the presence":      {cfg.PresenceWindow - time.Minute, StateConnected},
		"app in the background":    {cfg.PresenceWindow + time.Minute, StateActive},
		"a week ago":               {7 * 24 * time.Hour, StateActive},
		"inside the session":       {cfg.SessionTimeout - time.Hour, StateActive},
		"past the session timeout": {cfg.SessionTimeout + time.Hour, StateExpired},
		"a year without any proof": {365 * 24 * time.Hour, StateExpired},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			f := newFx(t)
			f.register(idA, "alice@x.test", tok1)
			f.seen(idA, c.ago)
			if got := f.state("alice@x.test", idA); got != c.want {
				t.Fatalf("state = %s, want %s", got, c.want)
			}
		})
	}
}

func TestExpiredSessionsAreNeitherConnectedNorActive(t *testing.T) {
	f := newFx(t)
	f.register(idA, "alice@x.test", tok1)
	f.register(idB, "alice@x.test", tok2)
	f.seen(idA, DefaultConfig().SessionTimeout+time.Hour)
	f.seen(idB, DefaultConfig().SessionTimeout+2*time.Hour)

	if ok, _ := f.svc.IsConnected(f.ctx, "alice@x.test"); ok {
		t.Fatal("expired devices must not be connected")
	}
	if ok, _ := f.svc.HasActiveSession(f.ctx, "alice@x.test"); ok {
		t.Fatal("expired devices must not be an active session")
	}
	// A registered token is not a session.
	if _, has := f.token(idA); !has {
		t.Fatal("the test needs devices that still hold a token")
	}
}

func TestConcurrentSessionsOfOneUser(t *testing.T) {
	f := newFx(t)
	f.register(idA, "alice@x.test", tok1) // open now
	f.register(idB, "alice@x.test", tok2) // in the background
	f.register(idC, "alice@x.test", tok3) // expired
	f.seen(idB, 2*time.Hour)
	f.seen(idC, DefaultConfig().SessionTimeout+time.Hour)

	want := map[string]State{idA: StateConnected, idB: StateActive, idC: StateExpired}
	for id, st := range want {
		if got := f.state("alice@x.test", id); got != st {
			t.Fatalf("device %s = %s, want %s", id, got, st)
		}
	}
	if ok, _ := f.svc.IsConnected(f.ctx, "alice@x.test"); !ok {
		t.Fatal("alice has an open device")
	}

	f.seen(idA, 2*time.Hour)
	if ok, _ := f.svc.IsConnected(f.ctx, "alice@x.test"); ok {
		t.Fatal("nothing is open any more")
	}
	if ok, _ := f.svc.HasActiveSession(f.ctx, "alice@x.test"); !ok {
		t.Fatal("two devices still hold a valid session")
	}
}

func TestPasswordChangeRevokesEverySessionUntilTheNextLogin(t *testing.T) {
	f := newFx(t)
	f.register(idA, "alice@x.test", tok1)
	f.register(idB, "alice@x.test", tok2)

	f.exec(`UPDATE newf SET password_updated_date = now() WHERE email = 'alice@x.test'`)
	for _, id := range []string{idA, idB} {
		if got := f.state("alice@x.test", id); got != StateRevoked {
			t.Fatalf("device %s = %s, want revoked", id, got)
		}
	}
	if ok, _ := f.svc.HasActiveSession(f.ctx, "alice@x.test"); ok {
		t.Fatal("revoked sessions are not active")
	}

	// The device that changed the password holds a new token and checks in at once: no waiting for the throttle.
	f.register(idA, "alice@x.test", tok1)
	if f.state("alice@x.test", idA) != StateConnected {
		t.Fatal("the device that logged in after the change is connected")
	}
	if f.state("alice@x.test", idB) != StateRevoked {
		t.Fatal("the other device stays revoked")
	}
}

func TestDeletingTheAccountDeletesItsDevices(t *testing.T) {
	f := newFx(t)
	f.register(idA, "alice@x.test", tok1)
	f.exec(`DELETE FROM newf WHERE email = 'alice@x.test'`)
	if n := f.count(`SELECT count(*) FROM devices`); n != 0 {
		t.Fatalf("%d devices left", n)
	}
}

func TestRegistrationIsValidated(t *testing.T) {
	f := newFx(t)
	base := Registration{InstallationID: idA, Email: "alice@x.test", Platform: "ios", PushToken: tok1, PermissionStatus: "granted"}
	cases := map[string]struct {
		mutate func(*Registration)
		want   error
	}{
		"empty installation id": {func(r *Registration) { r.InstallationID = "" }, ErrInvalidInstallationID},
		"not a uuid":            {func(r *Registration) { r.InstallationID = "abc" }, ErrInvalidInstallationID},
		"unknown platform":      {func(r *Registration) { r.Platform = "windows" }, ErrInvalidPlatform},
		"missing platform":      {func(r *Registration) { r.Platform = "" }, ErrInvalidPlatform},
		"not an expo token":     {func(r *Registration) { r.PushToken = "fcm:APA91b" }, ErrInvalidPushToken},
		"token with quotes":     {func(r *Registration) { r.PushToken = "ExponentPushToken[a'b]" }, ErrInvalidPushToken},
		"empty token brackets":  {func(r *Registration) { r.PushToken = "ExponentPushToken[]" }, ErrInvalidPushToken},
		"unknown permission":    {func(r *Registration) { r.PermissionStatus = "maybe" }, ErrInvalidPermission},
		"version too long":      {func(r *Registration) { r.AppVersion = strings.Repeat("1", 33) }, ErrInvalidAppVersion},
	}
	for name, c := range cases {
		r := base
		c.mutate(&r)
		if _, err := f.svc.Register(f.ctx, r); err != c.want {
			t.Fatalf("%s: error = %v, want %v", name, err, c.want)
		}
	}
	if n := f.count(`SELECT count(*) FROM devices`); n != 0 {
		t.Fatalf("%d devices stored from invalid registrations", n)
	}

	// Valid variations: a token with dashes (including "--") and upper-case ids.
	for _, token := range []string{"ExpoPushToken[abc_DEF-123]", "ExponentPushToken[ab--cd]"} {
		r := base
		r.InstallationID = strings.ToUpper(idB)
		r.PushToken = token
		if _, err := f.svc.Register(f.ctx, r); err != nil {
			t.Fatalf("token %q: %v", token, err)
		}
	}
}

func TestUnknownUserCannotRegister(t *testing.T) {
	f := newFx(t)
	_, err := f.svc.Register(f.ctx, Registration{InstallationID: idA, Email: "ghost@x.test", Platform: "ios"})
	if err == nil {
		t.Fatal("a registration for an account that does not exist must fail")
	}
}

func TestConcurrentHeartbeatsOfOneInstallation(t *testing.T) {
	f := newFx(t)
	var wg sync.WaitGroup
	errs := make(chan error, 40)
	for i := 0; i < 40; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := f.svc.Register(f.ctx, Registration{
				InstallationID: idA, Email: "alice@x.test", Platform: "ios", PushToken: tok1, PermissionStatus: "granted",
			})
			errs <- err
		}()
	}
	wg.Wait()
	close(errs)
	for err := range errs {
		if err != nil {
			t.Fatal(err)
		}
	}
	if n := f.count(`SELECT count(*) FROM devices`); n != 1 {
		t.Fatalf("%d rows", n)
	}
	if n := f.count(`SELECT count(*) FROM user_notification_tokens`); n != 1 {
		t.Fatalf("%d legacy rows", n)
	}
}

func TestConcurrentClaimsOfOneTokenLeaveOneHolder(t *testing.T) {
	f := newFx(t)
	for round := 0; round < 15; round++ {
		token := fmt.Sprintf("ExponentPushToken[race-%d-xxxxxxxxxx]", round)
		var wg sync.WaitGroup
		errs := make(chan error, 3)
		for _, id := range []string{idA, idB, idC} {
			wg.Add(1)
			go func() {
				defer wg.Done()
				_, err := f.svc.Register(f.ctx, Registration{
					InstallationID: id, Email: "alice@x.test", Platform: "android", PushToken: token, PermissionStatus: "granted",
				})
				errs <- err
			}()
		}
		wg.Wait()
		close(errs)
		for err := range errs {
			if err != nil {
				t.Fatalf("round %d: %v", round, err)
			}
		}
		if n := f.count(`SELECT count(*) FROM devices WHERE push_token = $1`, token); n != 1 {
			t.Fatalf("round %d: %d installations hold the token", round, n)
		}
	}
}

func TestConfigValidation(t *testing.T) {
	if err := DefaultConfig().Validate(); err != nil {
		t.Fatal(err)
	}
	bad := map[string]func(*Config){
		"zero timeout":              func(c *Config) { c.SessionTimeout = 0 },
		"presence beyond session":   func(c *Config) { c.PresenceWindow = c.SessionTimeout },
		"heartbeat too slow":        func(c *Config) { c.HeartbeatInterval = c.PresenceWindow },
		"throttle beyond heartbeat": func(c *Config) { c.WriteThrottle = c.HeartbeatInterval },
	}
	for name, mutate := range bad {
		cfg := DefaultConfig()
		mutate(&cfg)
		if cfg.Validate() == nil {
			t.Fatalf("%s: accepted", name)
		}
	}
}

func TestConfigFromEnv(t *testing.T) {
	t.Setenv("DEVICE_SESSION_TIMEOUT", "720h")
	t.Setenv("DEVICE_PRESENCE_WINDOW", "4m")
	t.Setenv("DEVICE_HEARTBEAT_INTERVAL", "2m")
	t.Setenv("DEVICE_WRITE_THROTTLE", "30s")
	cfg := ConfigFromEnv()
	if cfg.SessionTimeout != 720*time.Hour || cfg.PresenceWindow != 4*time.Minute ||
		cfg.HeartbeatInterval != 2*time.Minute || cfg.WriteThrottle != 30*time.Second {
		t.Fatalf("config = %+v", cfg)
	}

	// An inconsistent set falls back to the defaults as a whole.
	t.Setenv("DEVICE_HEARTBEAT_INTERVAL", "9m")
	if got := ConfigFromEnv(); got != DefaultConfig() {
		t.Fatalf("config = %+v, want the defaults", got)
	}
	t.Setenv("DEVICE_HEARTBEAT_INTERVAL", "garbage")
	t.Setenv("DEVICE_SESSION_TIMEOUT", "-5h")
	t.Setenv("DEVICE_PRESENCE_WINDOW", "")
	want := DefaultConfig()
	want.WriteThrottle = 30 * time.Second
	if got := ConfigFromEnv(); got != want {
		t.Fatalf("config = %+v, want %+v", got, want)
	}
}

func TestTimeoutsAreConfigurable(t *testing.T) {
	f := newFx(t)
	f.svc = NewService(f.db, Config{
		SessionTimeout: 2 * time.Hour, PresenceWindow: 2 * time.Minute, HeartbeatInterval: time.Minute, WriteThrottle: 10 * time.Second,
	})
	f.register(idA, "alice@x.test", tok1)
	f.seen(idA, 3*time.Minute)
	if got := f.state("alice@x.test", idA); got != StateActive {
		t.Fatalf("3 minutes with a 2 minute window = %s", got)
	}
	f.seen(idA, 3*time.Hour)
	if got := f.state("alice@x.test", idA); got != StateExpired {
		t.Fatalf("3 hours with a 2 hour timeout = %s", got)
	}
}
