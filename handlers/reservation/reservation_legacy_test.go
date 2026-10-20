package reservation_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/plugimt/transat-backend/services/reservations"
)

// expoCapture stands in for Expo: it answers every message with an ok ticket and keeps what was sent.
type expoCapture struct {
	mu       sync.Mutex
	messages []map[string]any
}

func (e *expoCapture) RoundTrip(req *http.Request) (*http.Response, error) {
	raw, _ := io.ReadAll(req.Body)
	var batch []map[string]any
	_ = json.Unmarshal(raw, &batch)
	tickets := make([]map[string]any, len(batch))
	for i := range batch {
		tickets[i] = map[string]any{"status": "ok", "id": "x"}
	}
	e.mu.Lock()
	e.messages = append(e.messages, batch...)
	e.mu.Unlock()
	out, _ := json.Marshal(map[string]any{"data": tickets})
	return &http.Response{StatusCode: 200, Header: http.Header{"Content-Type": {"application/json"}}, Body: io.NopCloser(bytes.NewReader(out)), Request: req}, nil
}

func (e *expoCapture) take() []map[string]any {
	e.mu.Lock()
	defer e.mu.Unlock()
	out := e.messages
	e.messages = nil
	return out
}

func (e *expoCapture) to(token string) map[string]any {
	for _, m := range e.take() {
		if m["to"] == token {
			return m
		}
	}
	return nil
}

// With PUSH_PIPELINE=legacy the reminders, the administrator's cancellation and the waiting list still reach the
// phones, straight through Expo, with the action buttons.
func TestWithoutThePushPipelineRemindersCancellationsAndTheWaitingListStillReachThePhones(t *testing.T) {
	expo := &expoCapture{}
	previous := http.DefaultTransport
	http.DefaultTransport = expo
	t.Cleanup(func() { http.DefaultTransport = previous })

	a := newAPIWith(t, true)
	if a.svc == nil || a.pipeline == nil {
		t.Fatal("setup")
	}
	for _, u := range []string{"owner", "q1", "q2"} {
		a.user(u+"@x.test", u, "X")
	}
	a.user("admin@x.test", "Ada", "Lovelace", "ADMIN")
	room := a.item("Salle Ampère", true)
	cable := a.item("Câble HDMI", false)
	_, start := slot()
	path := fmt.Sprintf("/reservation/item/%d", room)
	wait := path + "/waitlist"
	body := map[string]string{"start_date": start}

	// A reminder: the loan started 30 hours ago, so its day-one reminder is due.
	ago := time.Now().UTC().Add(-30 * time.Hour).Format(reservations.APITimeLayout)
	if r := a.call(http.MethodPatch, fmt.Sprintf("/reservation/item/%d", cable), "owner@x.test", map[string]string{"start_date": ago}); r.status != 200 {
		t.Fatalf("loan: %d %v", r.status, r.body)
	}
	if sent, err := a.svc.SendDueReminders(t.Context()); err != nil || sent != 1 {
		t.Fatalf("reminders: %d, %v", sent, err)
	}
	m := expo.to(a.tokenOf("owner@x.test"))
	if data, _ := m["data"].(map[string]any); m == nil || data["type"] != "reservation_reminder" || !strings.Contains(fmt.Sprint(m["body"]), "Câble HDMI") {
		t.Fatalf("reminder = %v", m)
	}
	if sent, _ := a.svc.SendDueReminders(t.Context()); sent != 0 {
		t.Fatal("the reminder was sent twice")
	}

	// A slot with a queue; the administrator cancels it: the owner is told by name, the first of the queue gets the buttons.
	if r := a.call(http.MethodPatch, path, "owner@x.test", body); r.status != 200 {
		t.Fatalf("booking: %d", r.status)
	}
	for _, u := range []string{"q1", "q2"} {
		if r := a.call(http.MethodPost, wait, u+"@x.test", body); r.status != 201 {
			t.Fatalf("%s joins: %d %v", u, r.status, r.body)
		}
	}
	parsed, _ := time.Parse(reservations.APITimeLayout, start)
	if _, err := a.svc.Cancel(t.Context(), reservations.CancelRequest{ItemID: room, Start: &parsed, ActorEmail: "admin@x.test", AsAdmin: true}); err != nil {
		t.Fatal(err)
	}
	sent := expo.take()
	byToken := map[any]map[string]any{}
	for _, m := range sent {
		byToken[m["to"]] = m
	}
	owner, q1 := byToken[a.tokenOf("owner@x.test")], byToken[a.tokenOf("q1@x.test")]
	if len(sent) != 2 || owner == nil || !strings.Contains(fmt.Sprint(owner["body"]), "Ada Lovelace") {
		t.Fatalf("pushes = %v", sent)
	}
	if q1 == nil || q1["categoryId"] != reservations.WaitlistActionCategory {
		t.Fatalf("the first of the queue got %v", q1)
	}
	data, _ := q1["data"].(map[string]any)
	if nav, _ := data["navigation"].(map[string]any); nav == nil || nav["type"] != "reservation" {
		t.Fatalf("data = %v: the in-app screen would not open", data)
	}

	// "Annuler": q2 is told straight away.
	if r := a.call(http.MethodDelete, wait, "q1@x.test", body); r.status != 200 || r.body["next_notified"] != true {
		t.Fatalf("decline: %d %v", r.status, r.body)
	}
	if m := expo.to(a.tokenOf("q2@x.test")); m == nil || m["categoryId"] != reservations.WaitlistActionCategory {
		t.Fatalf("q2 was not told: %v", m)
	}
}
