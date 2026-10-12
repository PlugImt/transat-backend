package push

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/plugimt/transat-backend/internal/database/dbtest"
	"github.com/plugimt/transat-backend/models"
	"github.com/plugimt/transat-backend/services/devices"
)

// fakeProvider records what it is given and answers as the test scripts it.
type fakeProvider struct {
	mu    sync.Mutex
	calls [][]Outbound
	// send answers a request; the default accepts every message with a ticket. call counts from 1.
	send func(call int, msgs []Outbound) ([]SendResult, error)
	// receipts are what Receipts returns; receiptErr fails the whole request.
	receipts       map[string]ReceiptResult
	receiptErr     error
	receiptCalls   int
	receiptTickets []string
}

func (p *fakeProvider) Name() string  { return "fake" }
func (p *fakeProvider) MaxBatch() int { return 100 }

func (p *fakeProvider) Send(_ context.Context, msgs []Outbound) ([]SendResult, error) {
	p.mu.Lock()
	p.calls = append(p.calls, append([]Outbound(nil), msgs...))
	call := len(p.calls)
	send := p.send
	p.mu.Unlock()
	if send != nil {
		return send(call, msgs)
	}
	return acceptAll(call, msgs), nil
}

func acceptAll(call int, msgs []Outbound) []SendResult {
	out := make([]SendResult, len(msgs))
	for i := range msgs {
		out[i] = SendResult{Kind: ResultAccepted, TicketID: fmt.Sprintf("tk-%d-%d", call, i)}
	}
	return out
}

func (p *fakeProvider) Receipts(_ context.Context, ids []string) (map[string]ReceiptResult, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.receiptCalls++
	p.receiptTickets = append(p.receiptTickets, ids...)
	if p.receiptErr != nil {
		return nil, p.receiptErr
	}
	out := map[string]ReceiptResult{}
	for _, id := range ids {
		if r, ok := p.receipts[id]; ok {
			out[id] = r
		}
	}
	return out, nil
}

// sent returns every message handed to the provider, across calls.
func (p *fakeProvider) sent() []Outbound {
	p.mu.Lock()
	defer p.mu.Unlock()
	var out []Outbound
	for _, c := range p.calls {
		out = append(out, c...)
	}
	return out
}

func (p *fakeProvider) callCount() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	return len(p.calls)
}

type fx struct {
	t   *testing.T
	db  *sql.DB
	p   *Pipeline
	fp  *fakeProvider
	dev *devices.Service
	ctx context.Context
}

func newFx(t *testing.T, tweak ...func(*Config)) *fx {
	db := dbtest.Open(t)
	cfg := DefaultConfig()
	cfg.MaxPerSecond = 0
	cfg.MaxAttempts = 3
	cfg.RetryBase = 2 * time.Second
	cfg.RetryMax = 20 * time.Second
	cfg.ReceiptDelay = time.Hour
	for _, f := range tweak {
		f(&cfg)
	}
	fp := &fakeProvider{}
	return &fx{t: t, db: db, fp: fp, ctx: context.Background(),
		p:   New(db, fp, cfg, devices.DefaultConfig()),
		dev: devices.NewService(db, devices.DefaultConfig())}
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

// user creates an account (language 1 = fr, 2 = en) subscribed to the given categories.
func (f *fx) user(email string, lang int, categories ...string) {
	f.exec(`INSERT INTO newf (email, password, first_name, last_name, language, password_updated_date)
		VALUES ($1, 'x', 'A', 'B', $2, now() - interval '3 years')`, email, lang)
	for _, c := range categories {
		f.exec(`INSERT INTO notifications (email, id_services) SELECT $1, id_services FROM services WHERE name = $2`, email, c)
	}
}

func uuid(n int) string { return fmt.Sprintf("00000000-0000-4000-8000-%012d", n) }
func tok(n int) string  { return fmt.Sprintf("ExponentPushToken[device-%d]", n) }

// device registers device n for a user and returns its token ("" registers it without one).
func (f *fx) device(n int, email, platform string, withToken bool) string {
	f.t.Helper()
	token := ""
	if withToken {
		token = tok(n)
	}
	if _, err := f.dev.Register(f.ctx, devices.Registration{
		InstallationID: uuid(n), Email: email, Platform: platform, PushToken: token, PermissionStatus: "granted",
	}); err != nil {
		f.t.Fatal(err)
	}
	return token
}

func (f *fx) message(category models.NotificationCategory, key string) Message {
	return Message{
		Kind: "test", Category: category, DedupeKey: key,
		Audience: Audience{NotificationAudience: models.NotificationAudience{Type: models.AudienceAll}},
		Content: Content{
			"fr": {Title: "Titre", Body: "Corps"},
			"en": {Title: "Title", Body: "Body"},
		},
		Navigation: models.RestaurantNavigation(),
	}
}

func (f *fx) enqueue(m Message) int64 {
	f.t.Helper()
	res, err := f.p.Enqueue(f.ctx, m)
	if err != nil {
		f.t.Fatal(err)
	}
	return res.ID
}

func (f *fx) run() Stats {
	f.t.Helper()
	st, err := f.p.RunOnce(f.ctx)
	if err != nil {
		f.t.Fatal(err)
	}
	return st
}

// statuses counts the deliveries of a message by status.
func (f *fx) statuses(id int64) map[string]int {
	f.t.Helper()
	out, err := f.p.groupCount(f.ctx, `SELECT status, count(*) FROM push_deliveries WHERE message_id = $1 GROUP BY status`, id)
	if err != nil {
		f.t.Fatal(err)
	}
	return out
}

func (f *fx) deliveryOf(id int64, device int) (status, code, skip string, attempts int) {
	f.t.Helper()
	var c, s sql.NullString
	err := f.db.QueryRow(`SELECT status, error_code, skip_reason, attempts FROM push_deliveries WHERE message_id = $1 AND device_id = $2`,
		id, uuid(device)).Scan(&status, &c, &s, &attempts)
	if err != nil {
		f.t.Fatalf("delivery of device %d: %v", device, err)
	}
	return status, c.String, s.String, attempts
}

func (f *fx) dueNow(id int64) {
	f.exec(`UPDATE push_deliveries SET next_attempt_at = now() WHERE message_id = $1 AND status = 'pending'`, id)
}

func (f *fx) wantStatuses(id int64, want map[string]int) {
	f.t.Helper()
	got := f.statuses(id)
	if fmt.Sprint(got) != fmt.Sprint(want) {
		f.t.Fatalf("statuses = %v, want %v", got, want)
	}
}

func tokensOf(msgs []Outbound) string {
	var s []string
	for _, m := range msgs {
		s = append(s, m.Token)
	}
	return strings.Join(s, " ")
}

func TestMultiDeviceDeliveryUsesEachUsersLanguage(t *testing.T) {
	f := newFx(t)
	f.user("alice@x.test", 1, "RESTAURANT")
	f.user("bob@x.test", 2, "RESTAURANT")
	a1 := f.device(1, "alice@x.test", "ios", true)
	a2 := f.device(2, "alice@x.test", "android", true)
	b1 := f.device(3, "bob@x.test", "ios", true)

	id := f.enqueue(f.message(models.NotificationRestaurant, "k1"))
	st := f.run()

	if st.FannedOut != 1 || st.Accepted != 3 {
		t.Fatalf("stats = %+v", st)
	}
	if f.fp.callCount() != 1 {
		t.Fatalf("%d provider calls for 3 messages", f.fp.callCount())
	}
	titles := map[string]string{}
	platforms := map[string]string{}
	tags := map[string]string{}
	for _, m := range f.fp.sent() {
		titles[m.Token] = m.Title
		platforms[m.Token] = m.Platform
		tags[m.Token] = m.DedupeTag
	}
	if titles[a1] != "Titre" || titles[a2] != "Titre" || titles[b1] != "Title" {
		t.Fatalf("titles = %v", titles)
	}
	if platforms[a1] != "ios" || platforms[a2] != "android" {
		t.Fatalf("platforms = %v", platforms)
	}
	if tags[a1] != fmt.Sprintf("m%d", id) {
		t.Fatalf("dedupe tag = %q", tags[a1])
	}
	f.wantStatuses(id, map[string]int{"accepted": 3})

	sum, err := f.p.Summary(f.ctx, id)
	if err != nil || *sum.Users != 2 || *sum.EligibleDevices != 3 || sum.Attempts != 3 {
		t.Fatalf("summary = %+v, %v", sum, err)
	}
	// The log never holds a token.
	for _, table := range []string{"push_deliveries", "push_messages"} {
		if n := f.count(`SELECT count(*) FROM ` + table + ` t WHERE t::text LIKE '%PushToken%'`); n != 0 {
			t.Fatalf("%s holds a push token", table)
		}
	}
}

func TestPreferencesAndTheSessionPolicyDecideWhoIsReached(t *testing.T) {
	f := newFx(t)
	for _, u := range []string{"ok", "optout", "expired", "revoked", "signedout", "notoken", "denied", "nodevice"} {
		cats := []string{"RESTAURANT"}
		if u == "optout" {
			cats = nil
		}
		f.user(u+"@x.test", 1, cats...)
	}
	good := f.device(1, "ok@x.test", "ios", true)
	f.device(2, "optout@x.test", "ios", true)
	f.device(3, "expired@x.test", "ios", true)
	f.exec(`UPDATE devices SET last_seen_at = now() - interval '90 days' WHERE installation_id = $1`, uuid(3))
	f.device(4, "revoked@x.test", "ios", true)
	f.exec(`UPDATE newf SET password_updated_date = now() WHERE email = 'revoked@x.test'`)
	f.device(5, "signedout@x.test", "ios", true)
	if _, err := f.dev.SignOut(f.ctx, uuid(5)); err != nil {
		t.Fatal(err)
	}
	f.device(6, "notoken@x.test", "ios", false)
	f.device(7, "denied@x.test", "ios", true)
	f.exec(`UPDATE devices SET permission_status = 'denied' WHERE installation_id = $1`, uuid(7))

	id := f.enqueue(f.message(models.NotificationRestaurant, "k"))
	f.run()

	if got := tokensOf(f.fp.sent()); got != good {
		t.Fatalf("sent to %q, want only %q", got, good)
	}
	sum, _ := f.p.Summary(f.ctx, id)
	want := map[string]int{"expired": 1, "revoked": 1, "disabled": 1, "no_token": 1, "permission_denied": 1, "no_device": 1}
	if fmt.Sprint(sum.Filtered) != fmt.Sprint(want) {
		t.Fatalf("filtered = %v, want %v", sum.Filtered, want)
	}
	if *sum.Users != 7 || *sum.EligibleDevices != 1 {
		t.Fatalf("users %d, eligible %d (the opted-out user is not in the audience)", *sum.Users, *sum.EligibleDevices)
	}
}

func TestEveryUsersDevicesAreReachedNotJustOne(t *testing.T) {
	f := newFx(t)
	f.user("alice@x.test", 1, "RESTAURANT")
	f.device(1, "alice@x.test", "ios", true)
	f.device(2, "alice@x.test", "android", true)
	f.device(3, "alice@x.test", "android", true)
	// One of them is a background session, one has not been seen for a month: both still hold a valid session.
	f.exec(`UPDATE devices SET last_seen_at = now() - interval '30 days' WHERE installation_id = $1`, uuid(3))
	f.enqueue(f.message(models.NotificationRestaurant, "k"))
	f.run()
	if n := len(f.fp.sent()); n != 3 {
		t.Fatalf("%d devices reached, want 3", n)
	}
}

func TestEligibilityIsCheckedAgainWhenTheMessageIsSent(t *testing.T) {
	f := newFx(t)
	for _, u := range []string{"signout", "switch", "optout", "refresh", "stay"} {
		f.user(u+"@x.test", 1, "RESTAURANT")
	}
	f.user("other@x.test", 1, "RESTAURANT")
	f.device(1, "signout@x.test", "ios", true)
	f.device(2, "switch@x.test", "ios", true)
	f.device(3, "optout@x.test", "ios", true)
	f.device(4, "refresh@x.test", "ios", true)
	stay := f.device(5, "stay@x.test", "ios", true)

	id := f.enqueue(f.message(models.NotificationRestaurant, "k"))
	if _, err := f.p.fanOutOne(f.ctx); err != nil {
		t.Fatal(err)
	}
	if n := f.count(`SELECT count(*) FROM push_deliveries WHERE message_id = $1`, id); n != 5 {
		t.Fatalf("%d deliveries after fan-out", n)
	}

	// Between fan-out and sending: one user signs out, another account takes the phone over, one switches the
	// category off, one phone gets a refreshed token.
	if _, err := f.dev.SignOut(f.ctx, uuid(1)); err != nil {
		t.Fatal(err)
	}
	if _, err := f.dev.Register(f.ctx, devices.Registration{InstallationID: uuid(2), Email: "other@x.test", Platform: "ios", PushToken: tok(2), PermissionStatus: "granted"}); err != nil {
		t.Fatal(err)
	}
	f.exec(`DELETE FROM notifications WHERE email = 'optout@x.test'`)
	newToken := "ExponentPushToken[refreshed]"
	if _, err := f.dev.Register(f.ctx, devices.Registration{InstallationID: uuid(4), Email: "refresh@x.test", Platform: "ios", PushToken: newToken, PermissionStatus: "granted"}); err != nil {
		t.Fatal(err)
	}

	f.run()
	got := map[string]bool{}
	for _, m := range f.fp.sent() {
		got[m.Token] = true
	}
	if len(got) != 2 || !got[stay] || !got[newToken] {
		t.Fatalf("sent to %v, want the unchanged device and the refreshed token", got)
	}
	for device, reason := range map[int]string{1: "disabled", 2: "owner_changed", 3: "preference_off"} {
		if status, _, skip, _ := f.deliveryOf(id, device); status != "skipped" || skip != reason {
			t.Fatalf("device %d: %s/%s, want skipped/%s", device, status, skip, reason)
		}
	}
}

func TestAMessageThatExpiredIsNotSent(t *testing.T) {
	f := newFx(t)
	f.user("alice@x.test", 1, "RESTAURANT")
	f.device(1, "alice@x.test", "ios", true)

	late := f.enqueue(f.message(models.NotificationRestaurant, "late"))
	f.exec(`UPDATE push_messages SET expires_at = now() - interval '1 minute' WHERE id = $1`, late)
	f.run()
	var status string
	_ = f.db.QueryRow(`SELECT status FROM push_messages WHERE id = $1`, late).Scan(&status)
	if status != "cancelled" || len(f.fp.sent()) != 0 {
		t.Fatalf("message %s, sent %d", status, len(f.fp.sent()))
	}

	// Fanned out in time, but the deadline passes before the delivery is attempted.
	waiting := f.enqueue(f.message(models.NotificationRestaurant, "waiting"))
	if _, err := f.p.fanOutOne(f.ctx); err != nil {
		t.Fatal(err)
	}
	f.exec(`UPDATE push_messages SET expires_at = now() - interval '1 second' WHERE id = $1`, waiting)
	f.run()
	f.wantStatuses(waiting, map[string]int{"expired": 1})
	if len(f.fp.sent()) != 0 {
		t.Fatal("an expired message was sent")
	}
}

func TestTheProviderIsGivenTheRemainingLifetimeAsTTL(t *testing.T) {
	f := newFx(t)
	f.user("alice@x.test", 1, "RESTAURANT")
	f.device(1, "alice@x.test", "ios", true)
	m := f.message(models.NotificationRestaurant, "k")
	m.ExpiresAt = time.Now().Add(10 * time.Minute)
	f.enqueue(m)
	f.run()
	ttl := f.fp.sent()[0].TTL
	if ttl < 9*time.Minute || ttl > 10*time.Minute {
		t.Fatalf("ttl = %v", ttl)
	}
}

func TestAnInvalidTokenDoesNotBlockTheOthersAndIsDeactivated(t *testing.T) {
	f := newFx(t)
	f.user("alice@x.test", 1, "RESTAURANT")
	f.user("bob@x.test", 1, "RESTAURANT")
	f.user("carol@x.test", 1, "RESTAURANT")
	f.device(1, "alice@x.test", "ios", true)
	dead := f.device(2, "bob@x.test", "android", true)
	f.device(3, "carol@x.test", "ios", true)
	f.fp.send = func(call int, msgs []Outbound) ([]SendResult, error) {
		res := acceptAll(call, msgs)
		for i, m := range msgs {
			if m.Token == dead {
				res[i] = SendResult{Kind: ResultInvalidDevice, Code: "DeviceNotRegistered", Detail: "gone"}
			}
		}
		return res, nil
	}
	id := f.enqueue(f.message(models.NotificationRestaurant, "k"))
	st := f.run()

	f.wantStatuses(id, map[string]int{"accepted": 2, "failed": 1})
	if st.Invalidated != 1 {
		t.Fatalf("stats = %+v", st)
	}
	var status, reason string
	var token sql.NullString
	_ = f.db.QueryRow(`SELECT status, COALESCE(disabled_reason, ''), push_token FROM devices WHERE installation_id = $1`, uuid(2)).Scan(&status, &reason, &token)
	if status != "disabled" || reason != "invalid_token" || token.Valid {
		t.Fatalf("device = %s/%s token=%v", status, reason, token)
	}
	if n := f.count(`SELECT count(*) FROM user_notification_tokens WHERE token = $1`, dead); n != 0 {
		t.Fatal("the dead token stays in the legacy table")
	}

	// The next message does not even try the dead device.
	f.fp.send = nil
	before := len(f.fp.sent())
	f.enqueue(f.message(models.NotificationRestaurant, "k2"))
	f.run()
	if sent := f.fp.sent()[before:]; len(sent) != 2 {
		t.Fatalf("second message went to %d devices", len(sent))
	}
}

func TestADeviceThatRefreshedItsTokenIsNotDisabledByNewsAboutTheOldOne(t *testing.T) {
	f := newFx(t)
	f.user("alice@x.test", 1, "RESTAURANT")
	f.device(1, "alice@x.test", "ios", true)
	f.fp.send = func(call int, msgs []Outbound) ([]SendResult, error) {
		// The app refreshes its token while the request is in flight.
		_, _ = f.dev.Register(f.ctx, devices.Registration{InstallationID: uuid(1), Email: "alice@x.test", Platform: "ios", PushToken: "ExponentPushToken[fresh]", PermissionStatus: "granted"})
		return []SendResult{{Kind: ResultInvalidDevice, Code: "DeviceNotRegistered"}}, nil
	}
	f.enqueue(f.message(models.NotificationRestaurant, "k"))
	st := f.run()
	if st.Invalidated != 0 {
		t.Fatal("the refreshed device was disabled")
	}
	var status string
	_ = f.db.QueryRow(`SELECT status FROM devices WHERE installation_id = $1`, uuid(1)).Scan(&status)
	if status != "active" {
		t.Fatalf("device %s", status)
	}
}

func TestTransientErrorsAreRetriedWithABoundedBackoff(t *testing.T) {
	f := newFx(t)
	f.user("alice@x.test", 1, "RESTAURANT")
	f.device(1, "alice@x.test", "ios", true)
	failing := true
	f.fp.send = func(call int, msgs []Outbound) ([]SendResult, error) {
		if failing {
			return nil, &RequestError{Code: "network", Retryable: true}
		}
		return acceptAll(call, msgs), nil
	}
	id := f.enqueue(f.message(models.NotificationRestaurant, "k"))

	f.run()
	status, code, _, attempts := f.deliveryOf(id, 1)
	if status != "pending" || code != "network" || attempts != 1 {
		t.Fatalf("after the first failure: %s/%s attempt %d", status, code, attempts)
	}
	var wait float64
	_ = f.db.QueryRow(`SELECT EXTRACT(EPOCH FROM next_attempt_at - now()) FROM push_deliveries WHERE message_id = $1`, id).Scan(&wait)
	if wait < 0.9 || wait > 2.1 {
		t.Fatalf("first retry in %.1fs, want between base/2 and base (1s to 2s)", wait)
	}

	// Not due yet: running again does not call the provider.
	f.run()
	if f.fp.callCount() != 1 {
		t.Fatalf("%d calls before the retry was due", f.fp.callCount())
	}

	f.dueNow(id)
	f.run()
	if _, _, _, attempts := f.deliveryOf(id, 1); attempts != 2 {
		t.Fatalf("attempts = %d", attempts)
	}
	_ = f.db.QueryRow(`SELECT EXTRACT(EPOCH FROM next_attempt_at - now()) FROM push_deliveries WHERE message_id = $1`, id).Scan(&wait)
	if wait < 1.9 || wait > 4.1 {
		t.Fatalf("second retry in %.1fs, want 2s to 4s (the delay doubles)", wait)
	}

	failing = false
	f.dueNow(id)
	f.run()
	status, code, _, attempts = f.deliveryOf(id, 1)
	if status != "accepted" || attempts != 3 || code != "" {
		t.Fatalf("after recovery: %s/%s attempt %d", status, code, attempts)
	}
}

func TestRetriesStopAtTheMaximumNumberOfAttempts(t *testing.T) {
	f := newFx(t)
	f.user("alice@x.test", 1, "RESTAURANT")
	f.device(1, "alice@x.test", "ios", true)
	f.fp.send = func(int, []Outbound) ([]SendResult, error) {
		return nil, &RequestError{Code: "TOO_MANY_REQUESTS", Retryable: true}
	}
	id := f.enqueue(f.message(models.NotificationRestaurant, "k"))
	for i := 0; i < 6; i++ {
		f.run()
		f.dueNow(id)
	}
	status, code, _, attempts := f.deliveryOf(id, 1)
	if status != "failed" || code != "retries_exhausted:TOO_MANY_REQUESTS" || attempts != 3 {
		t.Fatalf("%s/%s after %d attempts", status, code, attempts)
	}
	if f.fp.callCount() != 3 {
		t.Fatalf("%d provider calls, want exactly MaxAttempts", f.fp.callCount())
	}
}

func TestBackoffGrowsIsCappedAndJittered(t *testing.T) {
	base, max := 30*time.Second, 15*time.Minute
	prevHigh := time.Duration(0)
	for attempt := 1; attempt <= 12; attempt++ {
		low, high := Backoff(attempt, base, max, 0), Backoff(attempt, base, max, 0.999999)
		if low > high || high > max {
			t.Fatalf("attempt %d: %v..%v", attempt, low, high)
		}
		if high < prevHigh {
			t.Fatalf("attempt %d: the delay shrank", attempt)
		}
		prevHigh = high
	}
	if Backoff(1, base, max, 0) != 15*time.Second || Backoff(1, base, max, 1) != 30*time.Second {
		t.Fatal("the first delay is not between base/2 and base")
	}
	if Backoff(30, base, max, 1) != max {
		t.Fatal("the delay is not capped")
	}
	if Backoff(0, base, max, 0) != Backoff(1, base, max, 0) {
		t.Fatal("attempt 0 must behave like attempt 1")
	}
}

func TestARefusedBatchIsSplitToFindTheBadMessage(t *testing.T) {
	f := newFx(t)
	for i := 1; i <= 7; i++ {
		email := fmt.Sprintf("u%d@x.test", i)
		f.user(email, 1, "RESTAURANT")
		f.device(i, email, "ios", true)
	}
	poison := tok(5)
	f.fp.send = func(call int, msgs []Outbound) ([]SendResult, error) {
		for _, m := range msgs {
			if m.Token == poison {
				return nil, &RequestError{Code: "PUSH_TOO_MANY_EXPERIENCE_IDS", Status: 400, Isolate: true}
			}
		}
		return acceptAll(call, msgs), nil
	}
	id := f.enqueue(f.message(models.NotificationRestaurant, "k"))
	f.run()
	f.wantStatuses(id, map[string]int{"accepted": 6, "failed": 1})
	if status, code, _, _ := f.deliveryOf(id, 5); status != "failed" || code != "PUSH_TOO_MANY_EXPERIENCE_IDS" {
		t.Fatalf("poisoned device: %s/%s", status, code)
	}
}

func TestTicketLevelOutcomesAreRecordedPerDevice(t *testing.T) {
	f := newFx(t)
	for i := 1; i <= 3; i++ {
		email := fmt.Sprintf("u%d@x.test", i)
		f.user(email, 1, "RESTAURANT")
		f.device(i, email, "ios", true)
	}
	f.fp.send = func(call int, msgs []Outbound) ([]SendResult, error) {
		return []SendResult{
			{Kind: ResultRetry, Code: "MessageRateExceeded", Detail: "slow down " + tok(1)},
			{Kind: ResultRejected, Code: "MessageTooBig"},
			{Kind: ResultAccepted},
		}, nil
	}
	id := f.enqueue(f.message(models.NotificationRestaurant, "k"))
	f.run()
	for device, want := range map[int]string{1: "pending", 2: "failed", 3: "unconfirmed"} {
		if status, _, _, _ := f.deliveryOf(id, device); status != want {
			t.Fatalf("device %d: %s, want %s", device, status, want)
		}
	}
	var detail string
	_ = f.db.QueryRow(`SELECT error_detail FROM push_deliveries WHERE message_id = $1 AND device_id = $2`, id, uuid(1)).Scan(&detail)
	if strings.Contains(detail, "PushToken") {
		t.Fatalf("a token reached the log: %q", detail)
	}
}

func TestDuplicateMessagesAreQueuedOnce(t *testing.T) {
	f := newFx(t)
	f.user("alice@x.test", 1, "RESTAURANT")
	f.device(1, "alice@x.test", "ios", true)
	first, err := f.p.Enqueue(f.ctx, f.message(models.NotificationRestaurant, "restaurant_menu:2026-10-12"))
	if err != nil || first.Duplicate {
		t.Fatalf("first = %+v, %v", first, err)
	}
	second, err := f.p.Enqueue(f.ctx, f.message(models.NotificationRestaurant, "restaurant_menu:2026-10-12"))
	if err != nil || !second.Duplicate || second.ID != first.ID {
		t.Fatalf("second = %+v, %v", second, err)
	}
	f.run()
	f.run()
	if len(f.fp.sent()) != 1 {
		t.Fatalf("%d notifications sent for a message queued twice", len(f.fp.sent()))
	}
	// Without a key nothing is deduplicated: two admin messages are two messages.
	a, _ := f.p.Enqueue(f.ctx, f.message("", ""))
	b, _ := f.p.Enqueue(f.ctx, f.message("", ""))
	if a.ID == b.ID {
		t.Fatal("messages without a key were merged")
	}
}

func TestFanningOutTwiceCreatesNoSecondDelivery(t *testing.T) {
	f := newFx(t)
	f.user("alice@x.test", 1, "RESTAURANT")
	f.device(1, "alice@x.test", "ios", true)
	id := f.enqueue(f.message(models.NotificationRestaurant, "k"))
	f.run()
	// A crashed worker would have left the message pending: simulate it and fan out again.
	f.exec(`UPDATE push_messages SET status = 'pending' WHERE id = $1`, id)
	f.run()
	if n := f.count(`SELECT count(*) FROM push_deliveries WHERE message_id = $1`, id); n != 1 {
		t.Fatalf("%d deliveries", n)
	}
	if len(f.fp.sent()) != 1 {
		t.Fatalf("%d notifications sent", len(f.fp.sent()))
	}
}

func TestConcurrentWorkersSendEveryDeviceExactlyOnce(t *testing.T) {
	f := newFx(t, func(c *Config) { c.ClaimSize = 7 })
	const users = 60
	for i := 1; i <= users; i++ {
		email := fmt.Sprintf("u%d@x.test", i)
		f.user(email, 1, "RESTAURANT")
		f.device(i, email, "ios", true)
	}
	f.enqueue(f.message(models.NotificationRestaurant, "k"))

	var wg sync.WaitGroup
	for w := 0; w < 6; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < 4; i++ {
				if _, err := f.p.RunOnce(f.ctx); err != nil {
					t.Error(err)
				}
			}
		}()
	}
	wg.Wait()

	seen := map[string]int{}
	for _, m := range f.fp.sent() {
		seen[m.Token]++
	}
	if len(seen) != users {
		t.Fatalf("%d devices reached, want %d", len(seen), users)
	}
	for token, n := range seen {
		if n != 1 {
			t.Fatalf("%s sent %d times", token, n)
		}
	}
	if n := f.count(`SELECT count(*) FROM push_deliveries WHERE status = 'accepted' AND attempts = 1`); n != users {
		t.Fatalf("%d deliveries accepted at the first attempt", n)
	}
}

func TestBatchesRespectTheProvidersLimit(t *testing.T) {
	f := newFx(t)
	for i := 1; i <= 250; i++ {
		email := fmt.Sprintf("u%d@x.test", i)
		f.user(email, 1, "RESTAURANT")
		f.device(i, email, "ios", true)
	}
	f.enqueue(f.message(models.NotificationRestaurant, "k"))
	f.run()
	var sizes []int
	for _, c := range f.fp.calls {
		sizes = append(sizes, len(c))
	}
	if fmt.Sprint(sizes) != "[100 100 50]" {
		t.Fatalf("batch sizes = %v", sizes)
	}
}

func TestAWorkerThatDiedHasItsClaimsTakenBack(t *testing.T) {
	f := newFx(t)
	f.user("alice@x.test", 1, "RESTAURANT")
	f.user("bob@x.test", 1, "RESTAURANT")
	f.device(1, "alice@x.test", "ios", true)
	f.device(2, "bob@x.test", "ios", true)
	id := f.enqueue(f.message(models.NotificationRestaurant, "k"))
	if _, err := f.p.fanOutOne(f.ctx); err != nil {
		t.Fatal(err)
	}
	// Alice's delivery was claimed once and its worker died; Bob's has used all its attempts that way.
	f.exec(`UPDATE push_deliveries SET status = 'sending', attempts = 1, lease_until = now() - interval '1 minute' WHERE device_id = $1`, uuid(1))
	f.exec(`UPDATE push_deliveries SET status = 'sending', attempts = 3, lease_until = now() - interval '1 minute' WHERE device_id = $1`, uuid(2))
	f.run()

	if status, _, _, attempts := f.deliveryOf(id, 1); status != "accepted" || attempts != 2 {
		t.Fatalf("recovered delivery: %s after %d attempts", status, attempts)
	}
	if status, code, _, _ := f.deliveryOf(id, 2); status != "failed" || code != "interrupted" {
		t.Fatalf("a delivery that kept killing its worker: %s/%s", status, code)
	}
	if len(f.fp.sent()) != 1 {
		t.Fatalf("%d sent", len(f.fp.sent()))
	}
}

func TestAWorkerWhoseLeaseWasTakenOverCannotOverwriteTheResult(t *testing.T) {
	f := newFx(t)
	f.user("alice@x.test", 1, "RESTAURANT")
	f.device(1, "alice@x.test", "ios", true)
	id := f.enqueue(f.message(models.NotificationRestaurant, "k"))
	f.fp.send = func(call int, msgs []Outbound) ([]SendResult, error) {
		// Another worker finished this delivery while the request was in flight.
		f.exec(`UPDATE push_deliveries SET status = 'failed', error_code = 'other_worker', lease_until = NULL WHERE message_id = $1`, id)
		return acceptAll(call, msgs), nil
	}
	f.run()
	if status, code, _, _ := f.deliveryOf(id, 1); status != "failed" || code != "other_worker" {
		t.Fatalf("a stale worker overwrote the result: %s/%s", status, code)
	}
}

func TestReceiptsAreReadLaterAndEachOutcomeIsRecorded(t *testing.T) {
	f := newFx(t)
	for i := 1; i <= 5; i++ {
		email := fmt.Sprintf("u%d@x.test", i)
		f.user(email, 1, "RESTAURANT")
		f.device(i, email, "ios", true)
	}
	f.fp.send = func(call int, msgs []Outbound) ([]SendResult, error) {
		out := make([]SendResult, len(msgs))
		for i, m := range msgs {
			out[i] = SendResult{Kind: ResultAccepted, TicketID: "t-" + m.Token}
		}
		return out, nil
	}
	id := f.enqueue(f.message(models.NotificationRestaurant, "k"))
	f.run()
	f.wantStatuses(id, map[string]int{"accepted": 5})

	// Not due yet: no receipt request, and the message is not completed.
	f.run()
	if f.fp.receiptCalls != 0 {
		t.Fatal("receipts were asked for before their delay")
	}

	f.fp.receipts = map[string]ReceiptResult{
		"t-" + tok(1): {Kind: ReceiptHandedOff},
		"t-" + tok(2): {Kind: ReceiptInvalidDevice, Code: "DeviceNotRegistered"},
		"t-" + tok(3): {Kind: ReceiptRetry, Code: "MessageRateExceeded"},
		"t-" + tok(4): {Kind: ReceiptFailed, Code: "InvalidCredentials"},
		// device 5: no receipt yet
	}
	f.exec(`UPDATE push_deliveries SET receipt_due_at = now() - interval '1 second' WHERE message_id = $1`, id)
	st := f.run()
	if st.Receipts != 5 || st.Invalidated != 1 {
		t.Fatalf("stats = %+v", st)
	}
	f.wantStatuses(id, map[string]int{"accepted": 1, "failed": 2, "handed_off": 1, "pending": 1})
	if status, code, _, _ := f.deliveryOf(id, 2); status != "failed" || code != "DeviceNotRegistered" {
		t.Fatalf("device 2: %s/%s", status, code)
	}
	var deviceStatus string
	_ = f.db.QueryRow(`SELECT status FROM devices WHERE installation_id = $1`, uuid(2)).Scan(&deviceStatus)
	if deviceStatus != "disabled" {
		t.Fatalf("the device behind a confirmed dead token is %s", deviceStatus)
	}
	var due float64
	_ = f.db.QueryRow(`SELECT EXTRACT(EPOCH FROM receipt_due_at - now()) FROM push_deliveries WHERE device_id = $1`, uuid(5)).Scan(&due)
	if due < 200 || due > 310 {
		t.Fatalf("a missing receipt is asked again in %.0fs, want about 5 minutes", due)
	}

	// 24 hours later the provider no longer has it: unconfirmed, not "delivered" and not "failed".
	f.exec(`UPDATE push_deliveries SET accepted_at = now() - interval '25 hours', receipt_due_at = now() - interval '1 second' WHERE device_id = $1`, uuid(5))
	f.run()
	if status, code, _, _ := f.deliveryOf(id, 5); status != "unconfirmed" || code != "receipt_unavailable" {
		t.Fatalf("device 5: %s/%s", status, code)
	}

	// The rate-limited one goes through the normal retry path and is sent again.
	f.dueNow(id)
	f.run()
	if status, _, _, attempts := f.deliveryOf(id, 3); status != "accepted" || attempts != 2 {
		t.Fatalf("device 3: %s after %d attempts", status, attempts)
	}
}

func TestAFailingReceiptRequestIsAskedAgainWithoutLosingTickets(t *testing.T) {
	f := newFx(t)
	f.user("alice@x.test", 1, "RESTAURANT")
	f.device(1, "alice@x.test", "ios", true)
	id := f.enqueue(f.message(models.NotificationRestaurant, "k"))
	f.run()
	f.exec(`UPDATE push_deliveries SET receipt_due_at = now() - interval '1 second' WHERE message_id = $1`, id)
	f.fp.receiptErr = &RequestError{Code: "network", Retryable: true}
	if _, err := f.p.RunOnce(f.ctx); err == nil {
		t.Fatal("the failure was swallowed")
	}
	if status, _, _, _ := f.deliveryOf(id, 1); status != "accepted" {
		t.Fatalf("status = %s", status)
	}
	if n := f.count(`SELECT count(*) FROM push_deliveries WHERE lease_until IS NOT NULL`); n != 0 {
		t.Fatal("the lease was kept")
	}
}

func TestAMessageCompletesWhenEveryDeliveryIsSettled(t *testing.T) {
	f := newFx(t)
	f.user("alice@x.test", 1, "RESTAURANT")
	f.device(1, "alice@x.test", "ios", true)
	id := f.enqueue(f.message(models.NotificationRestaurant, "k"))
	f.run()
	if sum, _ := f.p.Summary(f.ctx, id); sum.Status != "fanned_out" {
		t.Fatalf("status = %s while a receipt is pending", sum.Status)
	}
	f.fp.receipts = map[string]ReceiptResult{"tk-1-0": {Kind: ReceiptHandedOff}}
	f.exec(`UPDATE push_deliveries SET receipt_due_at = now() - interval '1 second'`)
	st := f.run()
	sum, _ := f.p.Summary(f.ctx, id)
	if sum.Status != "completed" || st.Completed != 1 || sum.CompletedAt == nil {
		t.Fatalf("summary = %+v, stats = %+v", sum, st)
	}
	if sum.Deliveries["handed_off"] != 1 {
		t.Fatalf("deliveries = %v", sum.Deliveries)
	}
	// An audience nobody qualifies for completes at once.
	empty := f.enqueue(f.message(models.NotificationPhotos, "photos"))
	f.run()
	if sum, _ := f.p.Summary(f.ctx, empty); sum.Status != "completed" || *sum.EligibleDevices != 0 {
		t.Fatalf("empty message: %+v", sum)
	}
}

func TestOldFinishedMessagesArePurged(t *testing.T) {
	f := newFx(t, func(c *Config) { c.Retention = time.Hour })
	f.user("alice@x.test", 1, "RESTAURANT")
	f.device(1, "alice@x.test", "ios", true)
	old := f.enqueue(f.message(models.NotificationRestaurant, "old"))
	f.run()
	f.exec(`UPDATE push_messages SET status = 'completed', completed_at = now() - interval '2 hours' WHERE id = $1`, old)
	fresh := f.enqueue(f.message(models.NotificationRestaurant, "fresh"))
	f.run()
	if err := f.p.purge(f.ctx); err != nil {
		t.Fatal(err)
	}
	if n := f.count(`SELECT count(*) FROM push_messages WHERE id = $1`, old); n != 0 {
		t.Fatal("the old message was kept")
	}
	if n := f.count(`SELECT count(*) FROM push_deliveries WHERE message_id = $1`, old); n != 0 {
		t.Fatal("its deliveries were kept")
	}
	if n := f.count(`SELECT count(*) FROM push_messages WHERE id = $1`, fresh); n != 1 {
		t.Fatal("a recent message was purged")
	}
}

func TestLegacyTokensAreDevicesWhoseSessionIsTheAccountsLastActivity(t *testing.T) {
	f := newFx(t)
	f.user("alice@x.test", 1, "RESTAURANT")
	f.user("bob@x.test", 1, "RESTAURANT")
	if err := devices.RegisterLegacyToken(f.ctx, f.db, "alice@x.test", tok(1)); err != nil {
		t.Fatal(err)
	}
	if err := devices.RegisterLegacyToken(f.ctx, f.db, "alice@x.test", tok(1)); err != nil {
		t.Fatal(err)
	}
	if n := f.count(`SELECT count(*) FROM devices WHERE source = 'legacy'`); n != 1 {
		t.Fatalf("%d legacy devices for one token", n)
	}

	// An old version never says it is still used: only the account's activity keeps its device eligible.
	f.exec(`UPDATE devices SET last_seen_at = now() - interval '90 days'`)
	f.exec(`UPDATE newf SET last_activity = now() - interval '2 days' WHERE email = 'alice@x.test'`)
	f.enqueue(f.message(models.NotificationRestaurant, "active"))
	f.run()
	if len(f.fp.sent()) != 1 {
		t.Fatalf("a legacy device of an active account got %d notifications", len(f.fp.sent()))
	}

	f.exec(`UPDATE newf SET last_activity = now() - interval '90 days' WHERE email = 'alice@x.test'`)
	f.enqueue(f.message(models.NotificationRestaurant, "inactive"))
	f.run()
	if len(f.fp.sent()) != 1 {
		t.Fatal("a legacy device of a long inactive account was reached")
	}

	// The phone is handed to another account through the old version: the legacy device follows it.
	if err := devices.RegisterLegacyToken(f.ctx, f.db, "bob@x.test", tok(1)); err != nil {
		t.Fatal(err)
	}
	var owner string
	_ = f.db.QueryRow(`SELECT user_email FROM devices WHERE push_token = $1`, tok(1)).Scan(&owner)
	if owner != "bob@x.test" {
		t.Fatalf("owner = %s", owner)
	}
}

func TestLegacyRegistrationNeverTakesAnInstallationsTokenAway(t *testing.T) {
	f := newFx(t)
	f.user("alice@x.test", 1, "RESTAURANT")
	f.user("bob@x.test", 1, "RESTAURANT")
	f.device(1, "alice@x.test", "ios", true)
	if err := devices.RegisterLegacyToken(f.ctx, f.db, "bob@x.test", tok(1)); err != nil {
		t.Fatal(err)
	}
	var owner, source string
	_ = f.db.QueryRow(`SELECT user_email, source FROM devices WHERE push_token = $1`, tok(1)).Scan(&owner, &source)
	if owner != "alice@x.test" || source != "client" {
		t.Fatalf("%s/%s", owner, source)
	}
	if err := devices.RegisterLegacyToken(f.ctx, f.db, "bob@x.test", "not a token"); err != devices.ErrInvalidPushToken {
		t.Fatalf("error = %v", err)
	}
}

func TestAudiences(t *testing.T) {
	f := newFx(t)
	f.exec(`INSERT INTO clubs (name, picture) VALUES ('Club', 'p')`)
	var club int
	_ = f.db.QueryRow(`SELECT id_clubs FROM clubs`).Scan(&club)
	for i, c := range []struct{ email, campus string }{{"a@x.test", "NANTES"}, {"b@x.test", "BREST"}, {"c@x.test", "NANTES"}} {
		f.user(c.email, 1, "EVENT_REMINDERS")
		f.exec(`UPDATE newf SET campus = $2, formation_name = 'FISE', graduation_year = 2027 WHERE email = $1`, c.email, c.campus)
		f.device(i+1, c.email, "ios", true)
	}
	f.exec(`INSERT INTO clubs_members (email, id_clubs) VALUES ('a@x.test', $1), ('b@x.test', $1)`, club)

	emails := func(aud Audience, category models.NotificationCategory) string {
		res, err := f.p.Resolve(f.ctx, aud, category)
		if err != nil {
			t.Fatal(err)
		}
		var out []string
		for _, tg := range res.Targets {
			out = append(out, tg.Email)
		}
		return strings.Join(out, ",")
	}
	all := models.NotificationAudience{Type: models.AudienceAll}
	cases := map[string]struct {
		aud  Audience
		want string
	}{
		"all":      {Audience{NotificationAudience: all}, "a@x.test,b@x.test,c@x.test"},
		"club":     {Audience{NotificationAudience: models.NotificationAudience{Type: models.AudienceClub, ClubID: club}}, "a@x.test,b@x.test"},
		"campus":   {Audience{NotificationAudience: models.NotificationAudience{Type: models.AudienceCampus, Campus: "nantes"}}, "a@x.test,c@x.test"},
		"users":    {Audience{NotificationAudience: models.NotificationAudience{Type: models.AudienceUsers, Emails: []string{"c@x.test"}}}, "c@x.test"},
		"cohort":   {Audience{NotificationAudience: models.NotificationAudience{Type: models.AudienceCohort, Formation: "fise", GraduationYear: 2027}}, "a@x.test,b@x.test,c@x.test"},
		"interest": {Audience{NotificationAudience: all, InterestedClubID: club}, "a@x.test,b@x.test"},
	}
	for name, c := range cases {
		if got := emails(c.aud, "EVENT_REMINDERS"); got != c.want {
			t.Fatalf("%s: %q, want %q", name, got, c.want)
		}
	}
	if got := emails(Audience{NotificationAudience: all}, models.NotificationRestaurant); got != "" {
		t.Fatalf("a category nobody subscribed to reached %q", got)
	}
	if err := (Audience{NotificationAudience: models.NotificationAudience{Type: models.AudienceClub}}).Validate(); err == nil {
		t.Fatal("a club audience without a club is valid")
	}
}

func TestEnqueueRefusesInvalidMessages(t *testing.T) {
	f := newFx(t)
	base := f.message(models.NotificationRestaurant, "k")

	noFrench := base
	noFrench.Content = Content{"en": {Title: "Title"}}
	if _, err := f.p.Enqueue(f.ctx, noFrench); err != ErrNoContent {
		t.Fatalf("no default language: %v", err)
	}
	badAudience := base
	badAudience.Audience = Audience{NotificationAudience: models.NotificationAudience{Type: "everyone"}}
	if _, err := f.p.Enqueue(f.ctx, badAudience); err == nil {
		t.Fatal("unknown audience accepted")
	}
	badCategory := base
	badCategory.Category = "NOPE"
	if _, err := f.p.Enqueue(f.ctx, badCategory); err == nil {
		t.Fatal("unknown category accepted")
	}
	badNav := base
	badNav.Navigation = &models.NavigationTarget{Type: "nowhere"}
	if _, err := f.p.Enqueue(f.ctx, badNav); err == nil {
		t.Fatal("unknown destination accepted")
	}
	noKind := base
	noKind.Kind = ""
	if _, err := f.p.Enqueue(f.ctx, noKind); err == nil {
		t.Fatal("a message without a kind was accepted")
	}
	if n := f.count(`SELECT count(*) FROM push_messages`); n != 0 {
		t.Fatalf("%d invalid messages were stored", n)
	}
}

func TestOversizedTextIsShrunkToFitThePayloadLimit(t *testing.T) {
	f := newFx(t)
	f.user("alice@x.test", 1, "RESTAURANT")
	f.device(1, "alice@x.test", "ios", true)
	m := f.message(models.NotificationRestaurant, "big")
	m.Content = Content{"fr": {Title: "Titre", Body: strings.Repeat("é", 6000)}}
	f.enqueue(m)
	f.run()
	body := f.fp.sent()[0].Body
	if len(body) == 0 || payloadSize(Text{Title: "Titre", Body: body}, f.fp.sent()[0].Data) > maxPayloadBytes {
		t.Fatalf("body of %d bytes does not fit", len(body))
	}
	if !strings.HasSuffix(body, "…") {
		t.Fatal("a shortened text should say so")
	}

	// Navigation data that alone exceeds the limit is refused rather than sent truncated.
	huge := f.message(models.NotificationRestaurant, "huge")
	huge.Navigation = &models.NavigationTarget{Version: 1, Type: models.NavigationPhotoAlbum, ID: "1", Params: map[string]string{"title": strings.Repeat("x", 5000)}}
	if _, err := f.p.Enqueue(f.ctx, huge); err != ErrPayloadTooLarge {
		t.Fatalf("error = %v", err)
	}
}

func TestAnImageIsCarriedOnlyByTheMessagesThatHaveOne(t *testing.T) {
	f := newFx(t)
	f.user("alice@x.test", 1, "RESTAURANT")
	f.device(1, "alice@x.test", "ios", true)
	f.device(2, "alice@x.test", "android", true)

	plain := f.message(models.NotificationRestaurant, "plain")
	f.enqueue(plain)
	f.run()
	for _, out := range f.fp.sent() {
		if out.ImageURL != "" {
			t.Fatalf("a text notification got an image: %+v", out)
		}
		if out.Data["type"] != "test" {
			t.Fatalf("the notification type is missing: %v", out.Data)
		}
	}

	withImage := f.message(models.NotificationRestaurant, "image")
	withImage.ImageURL = "https://cdn.x.test/menu.jpg"
	id := f.enqueue(withImage)
	f.run()
	sent := f.fp.sent()[2:]
	if len(sent) != 2 {
		t.Fatalf("%d messages sent for the image notification", len(sent))
	}
	for _, out := range sent {
		if out.ImageURL != "https://cdn.x.test/menu.jpg" {
			t.Fatalf("image = %q", out.ImageURL)
		}
	}
	if got := f.count(`SELECT count(*) FROM push_messages WHERE id = $1 AND image_url = $2`, id, "https://cdn.x.test/menu.jpg"); got != 1 {
		t.Fatal("the image link was not stored with the message")
	}
}

func TestImagesThatCannotBeShownAreRefusedBeforeQueueing(t *testing.T) {
	f := newFx(t)
	for name, link := range map[string]string{
		"http":     "http://cdn.x.test/a.jpg",
		"script":   "javascript:alert(1)",
		"relative": "/a.jpg",
		"too long": "https://cdn.x.test/" + strings.Repeat("a", models.MaxNotificationImageURLLength),
	} {
		m := f.message(models.NotificationRestaurant, name)
		m.ImageURL = link
		if _, err := f.p.Enqueue(f.ctx, m); err == nil {
			t.Fatalf("%s: the image was accepted", name)
		}
	}
	if n := f.count(`SELECT count(*) FROM push_messages`); n != 0 {
		t.Fatalf("%d messages with a bad image were stored", n)
	}
}

func TestTheImageLinkTakesItsShareOfThePayloadLimit(t *testing.T) {
	fit := func(image string) string {
		f := newFx(t)
		f.user("alice@x.test", 1, "RESTAURANT")
		f.device(1, "alice@x.test", "ios", true)
		m := f.message(models.NotificationRestaurant, "k")
		m.Content = Content{"fr": {Title: "Titre", Body: strings.Repeat("a", 7000)}}
		m.ImageURL = image
		f.enqueue(m)
		f.run()
		return f.fp.sent()[0].Body
	}
	link := "https://cdn.x.test/" + strings.Repeat("p", 400)
	if without, with := fit(""), fit(link); len(with) >= len(without) {
		t.Fatalf("the text did not make room for the image: %d bytes with, %d without", len(with), len(without))
	}
}

func TestSummariesCarryCountsNeverTokensOrText(t *testing.T) {
	f := newFx(t)
	f.user("alice@x.test", 1, "RESTAURANT")
	f.device(1, "alice@x.test", "ios", true)
	f.fp.send = func(int, []Outbound) ([]SendResult, error) {
		return []SendResult{{Kind: ResultRejected, Code: "MessageTooBig", Detail: "too big " + tok(1)}}, nil
	}
	id := f.enqueue(f.message(models.NotificationRestaurant, "k"))
	f.run()
	sum, err := f.p.Summary(f.ctx, id)
	if err != nil {
		t.Fatal(err)
	}
	if sum.Errors["MessageTooBig"] != 1 || sum.Deliveries["failed"] != 1 {
		t.Fatalf("summary = %+v", sum)
	}
	recent, err := f.p.Recent(f.ctx, 10)
	if err != nil || len(recent) != 1 || recent[0].ID != id {
		t.Fatalf("recent = %+v, %v", recent, err)
	}
	if _, err := f.p.Summary(f.ctx, 9999); err != ErrMessageNotFound {
		t.Fatalf("error = %v", err)
	}
}

func TestPacingKeepsTheRateUnderTheLimit(t *testing.T) {
	f := newFx(t, func(c *Config) { c.MaxPerSecond = 100; c.BatchSize = 10 })
	for i := 1; i <= 30; i++ {
		email := fmt.Sprintf("u%d@x.test", i)
		f.user(email, 1, "RESTAURANT")
		f.device(i, email, "ios", true)
	}
	f.enqueue(f.message(models.NotificationRestaurant, "k"))
	start := time.Now()
	f.run()
	// 30 messages at 100 per second take at least 0.3 s.
	if elapsed := time.Since(start); elapsed < 250*time.Millisecond {
		t.Fatalf("30 messages went out in %v, faster than 100 per second", elapsed)
	}
}

func TestConfigFromEnv(t *testing.T) {
	t.Setenv("PUSH_MAX_ATTEMPTS", "7")
	t.Setenv("PUSH_RETRY_BASE", "10s")
	t.Setenv("PUSH_RETRY_MAX", "5m")
	t.Setenv("PUSH_BATCH_SIZE", "50")
	t.Setenv("PUSH_DEFAULT_TTL", "garbage")
	cfg := ConfigFromEnv()
	if cfg.MaxAttempts != 7 || cfg.RetryBase != 10*time.Second || cfg.RetryMax != 5*time.Minute || cfg.BatchSize != 50 {
		t.Fatalf("config = %+v", cfg)
	}
	if cfg.DefaultTTL != DefaultConfig().DefaultTTL {
		t.Fatal("an invalid value must fall back to the default")
	}
	t.Setenv("PUSH_RETRY_MAX", "1s")
	if cfg := ConfigFromEnv(); cfg.RetryBase != DefaultConfig().RetryBase || cfg.RetryMax != DefaultConfig().RetryMax {
		t.Fatalf("an inconsistent backoff must fall back to the defaults: %+v", cfg)
	}
}
