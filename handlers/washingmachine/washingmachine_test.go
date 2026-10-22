package washingmachine

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/plugimt/transat-backend/i18n"
)

const validBody = `{"machine_info_status":{"machine_list":[
	{"machine_id":"1","nom_type":"LAVE LINGE","selecteur_machine":"11","status":1,"time_before_off":0},
	{"machine_id":"2","nom_type":"SECHE LINGE","selecteur_machine":"15","status":0,"time_before_off":300}]}}`

// newHandler points a handler at a fake upstream and returns it with a counter of the calls the upstream received.
func newHandler(t *testing.T, upstream http.HandlerFunc) (*fiber.App, *int32) {
	t.Helper()
	if err := i18n.Init(); err != nil {
		t.Fatal(err)
	}
	var calls int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&calls, 1)
		upstream(w, r)
	}))
	t.Cleanup(server.Close)

	h := NewWashingMachineHandler()
	h.url = server.URL
	h.client.Timeout = 300 * time.Millisecond

	app := fiber.New()
	app.Get("/washingmachines", h.GetWashingMachines())
	return app, &calls
}

func get(t *testing.T, app *fiber.App, lang string) (int, map[string]any, http.Header) {
	t.Helper()
	req := httptest.NewRequest(http.MethodGet, "/washingmachines", nil)
	if lang != "" {
		req.Header.Set("Accept-Language", lang)
	}
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)
	var body map[string]any
	if err := json.Unmarshal(raw, &body); err != nil {
		t.Fatalf("response is not JSON: %s", raw)
	}
	return resp.StatusCode, body, resp.Header
}

func TestHealthyUpstreamIsServedAndReusedBriefly(t *testing.T) {
	app, calls := newHandler(t, func(w http.ResponseWriter, _ *http.Request) { _, _ = io.WriteString(w, validBody) })

	for range 2 {
		status, body, _ := get(t, app, "")
		if status != 200 || body["success"] != true {
			t.Fatalf("status %d, body %v", status, body)
		}
	}
	if *calls != 1 {
		t.Fatalf("upstream called %d times, want 1 (the second answer comes from the short cache)", *calls)
	}
}

func TestBrokenUpstreamGivesAClearMessage(t *testing.T) {
	cases := map[string]http.HandlerFunc{
		"error status": func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusServiceUnavailable) },
		"html page": func(w http.ResponseWriter, _ *http.Request) {
			_, _ = io.WriteString(w, "<html><body>Application paused</body></html>")
		},
		"empty list": func(w http.ResponseWriter, _ *http.Request) {
			_, _ = io.WriteString(w, `{"machine_info_status":{"machine_list":[]}}`)
		},
		"timeout": func(w http.ResponseWriter, _ *http.Request) { time.Sleep(time.Second) },
	}
	for name, upstream := range cases {
		t.Run(name, func(t *testing.T) {
			app, _ := newHandler(t, upstream)
			status, body, header := get(t, app, "fr-FR,fr;q=0.9")

			if status != http.StatusServiceUnavailable {
				t.Fatalf("status %d, want 503", status)
			}
			// The app displays "error" as is.
			if msg, _ := body["error"].(string); msg == "" || msg[:3] != "Les" {
				t.Fatalf("error %q is not the French message", body["error"])
			}
			if body["success"] != false || body["code"] != "washing_machines_unavailable" || header.Get("Retry-After") == "" {
				t.Fatalf("body %v, headers %v", body, header)
			}
		})
	}
}

func TestMessageFollowsTheLanguageOfTheApp(t *testing.T) {
	app, _ := newHandler(t, func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusBadGateway) })

	_, english, _ := get(t, app, "en-GB")
	_, unknown, _ := get(t, app, "")
	if msg, _ := english["error"].(string); msg[:7] != "Washing" {
		t.Fatalf("english message: %q", msg)
	}
	if msg, _ := unknown["error"].(string); msg[:3] != "Les" {
		t.Fatalf("without a language the French default applies: %q", msg)
	}
}

func TestUpstreamIsLeftAloneAfterAFailure(t *testing.T) {
	app, calls := newHandler(t, func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusServiceUnavailable) })

	get(t, app, "")
	get(t, app, "")
	if *calls != 1 {
		t.Fatalf("upstream called %d times while down, want 1", *calls)
	}
}
