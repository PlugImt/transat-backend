package services

import (
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"sort"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/plugimt/transat-backend/internal/database/dbtest"
	"github.com/plugimt/transat-backend/models"
)

type expoMock struct {
	mu       sync.Mutex
	requests [][]map[string]any
	// Per request, the HTTP status to answer; empty means 200 with a ticket per message.
	statuses []int
	// Tokens answered with DeviceNotRegistered.
	dead   map[string]bool
	server *httptest.Server
}

func newExpoMock(t *testing.T) *expoMock {
	m := &expoMock{dead: map[string]bool{}}
	m.server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		raw, _ := io.ReadAll(r.Body)
		var messages []map[string]any
		_ = json.Unmarshal(raw, &messages)

		m.mu.Lock()
		call := len(m.requests)
		m.requests = append(m.requests, messages)
		status := http.StatusOK
		if call < len(m.statuses) {
			status = m.statuses[call]
		}
		m.mu.Unlock()

		if status != http.StatusOK {
			w.WriteHeader(status)
			_, _ = w.Write([]byte(`{"errors":[{"message":"nope"}]}`))
			return
		}
		tickets := make([]map[string]any, len(messages))
		for i, msg := range messages {
			token, _ := msg["to"].(string)
			if m.dead[token] {
				tickets[i] = map[string]any{"status": "error", "message": `"` + token + `" is not a registered push notification recipient`, "details": map[string]any{"error": "DeviceNotRegistered"}}
			} else {
				tickets[i] = map[string]any{"status": "ok", "id": "x"}
			}
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"data": tickets})
	}))
	t.Cleanup(m.server.Close)
	return m
}

// tokens lists every device addressed, across requests.
func (m *expoMock) tokens() []string {
	m.mu.Lock()
	defer m.mu.Unlock()
	var out []string
	for _, req := range m.requests {
		for _, msg := range req {
			out = append(out, msg["to"].(string))
		}
	}
	sort.Strings(out)
	return out
}

func (m *expoMock) messages() []map[string]any {
	m.mu.Lock()
	defer m.mu.Unlock()
	var out []map[string]any
	for _, req := range m.requests {
		out = append(out, req...)
	}
	return out
}

type fixture struct {
	t  *testing.T
	db *sql.DB
	ns *NotificationService
	ex *expoMock
}

func newFixture(t *testing.T) *fixture {
	db := dbtest.Open(t)
	ex := newExpoMock(t)
	ns := NewNotificationService(db)
	ns.expoPushURL = ex.server.URL
	ns.retryDelay = time.Millisecond
	return &fixture{t: t, db: db, ns: ns, ex: ex}
}

func (f *fixture) exec(q string, args ...any) {
	f.t.Helper()
	if _, err := f.db.Exec(q, args...); err != nil {
		f.t.Fatalf("%s: %v", q, err)
	}
}

// user creates an account (language 1 = fr, 2 = en) with its tokens. Photos is on only when optIn is set.
func (f *fixture) user(email string, lang int, optIn bool, tokens ...string) {
	f.exec(`INSERT INTO newf (email, password, first_name, last_name, language) VALUES ($1, 'x', 'A', 'B', $2)`, email, lang)
	if err := SubscribeToDefaultCategories(f.db, email); err != nil {
		f.t.Fatal(err)
	}
	if optIn {
		if err := f.ns.SetPreference(email, models.NotificationPhotos, true); err != nil {
			f.t.Fatal(err)
		}
	}
	for _, token := range tokens {
		f.exec(`INSERT INTO user_notification_tokens (email, token) VALUES ($1, $2)`, email, token)
	}
}

func tok(n string) string { return "ExponentPushToken[" + n + "]" }

func TestPhotosNotificationsAreOffByDefault(t *testing.T) {
	f := newFixture(t)
	f.user("new@imt.fr", 1, false, tok("a"))

	prefs, err := f.ns.GetPreferences("new@imt.fr")
	if err != nil {
		t.Fatal(err)
	}
	seen := map[models.NotificationCategory]bool{}
	for _, p := range prefs {
		seen[p.Service] = p.Enabled
	}
	if on, known := seen[models.NotificationPhotos]; !known || on {
		t.Fatalf("photos preference = %v (known=%v), want known and off", on, known)
	}
	for _, c := range []models.NotificationCategory{models.NotificationRestaurant, models.NotificationEvents, models.NotificationReservations} {
		if !seen[c] {
			t.Fatalf("%s is off for a new user: registration must keep subscribing the other categories", c)
		}
	}

	delivery, err := f.ns.NotifyNewPhotoAlbum(1, "Gala", nil)
	if err != nil || delivery.Devices != 0 || len(f.ex.tokens()) != 0 {
		t.Fatalf("a user who never opted in was notified: %+v, %v, %v", delivery, err, f.ex.tokens())
	}
}

func TestSwitchingPhotosOnOrOffLeavesOtherPreferencesAlone(t *testing.T) {
	f := newFixture(t)
	f.user("a@imt.fr", 1, false)
	before, _ := f.ns.GetPreferences("a@imt.fr")

	if err := f.ns.SetPreference("a@imt.fr", models.NotificationPhotos, true); err != nil {
		t.Fatal(err)
	}
	on, _ := f.ns.GetPreferences("a@imt.fr")
	for i := range before {
		want := before[i].Enabled
		if before[i].Service == models.NotificationPhotos {
			want = true
		}
		if on[i].Enabled != want {
			t.Fatalf("%s = %v after switching photos on", on[i].Service, on[i].Enabled)
		}
	}
	// Idempotent, and off again restores the starting state.
	if err := f.ns.SetPreference("a@imt.fr", models.NotificationPhotos, true); err != nil {
		t.Fatal(err)
	}
	if err := f.ns.SetPreference("a@imt.fr", models.NotificationPhotos, false); err != nil {
		t.Fatal(err)
	}
	off, _ := f.ns.GetPreferences("a@imt.fr")
	for i := range before {
		if off[i] != before[i] {
			t.Fatalf("%+v differs from the starting %+v", off[i], before[i])
		}
	}
}

func TestOnlyOptedInUsersWithADeviceAreNotified(t *testing.T) {
	f := newFixture(t)
	f.user("in1@imt.fr", 1, true, tok("phone1"), tok("tablet1")) // two devices
	f.user("in2@imt.fr", 2, true, tok("phone2"))
	f.user("in3@imt.fr", 1, true) // opted in, no device
	f.user("out@imt.fr", 1, false, tok("phone3"))
	// A phone shared after a logout: its token is also registered to an account that declined.
	f.user("in4@imt.fr", 1, true, tok("shared"))
	f.user("out2@imt.fr", 1, false, tok("shared"))
	// Two opted-in accounts on one phone: addressed once.
	f.user("in5@imt.fr", 1, true, tok("both"))
	f.user("in6@imt.fr", 1, true, tok("both"))
	// Junk in the token table must not reach Expo.
	f.user("in7@imt.fr", 1, true, "not-a-token")

	date := time.Date(2026, 3, 12, 0, 0, 0, 0, time.UTC)
	delivery, err := f.ns.NotifyNewPhotoAlbum(42, "Soirée de gala", &date)
	if err != nil {
		t.Fatal(err)
	}

	want := []string{tok("both"), tok("phone1"), tok("phone2"), tok("tablet1")}
	got := f.ex.tokens()
	if strings.Join(got, ",") != strings.Join(want, ",") {
		t.Fatalf("addressed %v, want %v", got, want)
	}
	if delivery.Accepted != 4 || delivery.Skipped != 1 {
		t.Fatalf("delivery = %+v", delivery)
	}
	var leftover int
	_ = f.db.QueryRow(`SELECT count(*) FROM user_notification_tokens WHERE token = 'not-a-token'`).Scan(&leftover)
	if leftover != 0 {
		t.Fatal("a malformed token was kept")
	}
}

func TestNotificationContentIsLocalizedAndPrivate(t *testing.T) {
	f := newFixture(t)
	f.user("fr@imt.fr", 1, true, tok("fr"))
	f.user("en@imt.fr", 2, true, tok("en"))

	date := time.Date(2026, 3, 12, 0, 0, 0, 0, time.UTC)
	if _, err := f.ns.NotifyNewPhotoAlbum(42, "Soirée de gala", &date); err != nil {
		t.Fatal(err)
	}
	if _, err := f.ns.NotifyNewPhotoAlbum(43, "Sans date", nil); err != nil {
		t.Fatal(err)
	}

	byToken := map[string]map[string]any{}
	var undated map[string]any
	for _, m := range f.ex.messages() {
		if m["data"].(map[string]any)["navigation"].(map[string]any)["id"] == "43" {
			undated = m
			continue
		}
		byToken[m["to"].(string)] = m
	}
	fr, en := byToken[tok("fr")], byToken[tok("en")]
	if fr["title"] != "Nouvel album photo" || fr["body"] != "Soirée de gala · 12/03/2026" {
		t.Fatalf("fr = %v / %v", fr["title"], fr["body"])
	}
	if en["title"] != "New photo album" || en["body"] != "Soirée de gala · Mar 12, 2026" {
		t.Fatalf("en = %v / %v", en["title"], en["body"])
	}
	if undated["body"] != "Sans date" {
		t.Fatalf("undated body = %v", undated["body"])
	}

	nav := fr["data"].(map[string]any)["navigation"].(map[string]any)
	params := nav["params"].(map[string]any)
	if nav["type"] != "photo_album" || nav["id"] != "42" || params["title"] != "Soirée de gala" {
		t.Fatalf("navigation = %v", nav)
	}
	raw, _ := json.Marshal(fr)
	for _, private := range []string{"@imt.fr", "fr@", "password"} {
		if strings.Contains(string(raw), private) {
			t.Fatalf("payload leaks %q: %s", private, raw)
		}
	}
}

func TestUnregisteredDevicesAreRemovedAndDoNotFailTheRest(t *testing.T) {
	f := newFixture(t)
	f.user("a@imt.fr", 1, true, tok("alive"), tok("dead"))
	f.ex.dead[tok("dead")] = true

	delivery, err := f.ns.NotifyNewPhotoAlbum(1, "Gala", nil)
	if err != nil {
		t.Fatalf("one dead device failed the announcement: %v", err)
	}
	if delivery.Accepted != 1 || delivery.Failed != 1 {
		t.Fatalf("delivery = %+v", delivery)
	}
	var remaining []string
	rows, _ := f.db.Query(`SELECT token FROM user_notification_tokens ORDER BY token`)
	for rows.Next() {
		var s string
		_ = rows.Scan(&s)
		remaining = append(remaining, s)
	}
	rows.Close()
	if len(remaining) != 1 || remaining[0] != tok("alive") {
		t.Fatalf("tokens left = %v", remaining)
	}
}

func TestTransientExpoErrorsAreRetried(t *testing.T) {
	f := newFixture(t)
	f.user("a@imt.fr", 1, true, tok("a"))
	f.ex.statuses = []int{http.StatusInternalServerError, http.StatusTooManyRequests}

	delivery, err := f.ns.NotifyNewPhotoAlbum(1, "Gala", nil)
	if err != nil || delivery.Accepted != 1 {
		t.Fatalf("delivery = %+v, err = %v", delivery, err)
	}
	if len(f.ex.requests) != 3 {
		t.Fatalf("%d requests, want 3 (two transient failures then success)", len(f.ex.requests))
	}
}

func TestTotalDeliveryFailureIsReportedSoItCanBeRetried(t *testing.T) {
	f := newFixture(t)
	f.user("a@imt.fr", 1, true, tok("a"))
	f.ex.statuses = []int{http.StatusBadRequest}

	delivery, err := f.ns.NotifyNewPhotoAlbum(1, "Gala", nil)
	if err == nil || delivery.Accepted != 0 {
		t.Fatalf("err = %v, delivery = %+v", err, delivery)
	}
	if strings.Contains(err.Error(), "ExponentPushToken") {
		t.Fatalf("error exposes a push token: %v", err)
	}
	if len(f.ex.requests) != 1 {
		t.Fatalf("a client error was retried %d times", len(f.ex.requests))
	}
}

func TestLargeAudiencesAreSentInBatches(t *testing.T) {
	f := newFixture(t)
	for i := 0; i < expoBatchSize+30; i++ {
		f.user(strings.Repeat("u", 1)+string(rune('a'+i%26))+string(rune('a'+i/26))+"@imt.fr", 1, true, tok(string(rune('a'+i%26))+string(rune('a'+i/26))))
	}
	delivery, err := f.ns.NotifyNewPhotoAlbum(1, "Gala", nil)
	if err != nil || delivery.Accepted != expoBatchSize+30 {
		t.Fatalf("delivery = %+v, err = %v", delivery, err)
	}
	if len(f.ex.requests) != 2 || len(f.ex.requests[0]) != expoBatchSize {
		t.Fatalf("%d requests, first of %d", len(f.ex.requests), len(f.ex.requests[0]))
	}
}
