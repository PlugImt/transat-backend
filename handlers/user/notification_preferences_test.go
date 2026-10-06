package user

import (
	"encoding/json"
	"io"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/plugimt/transat-backend/internal/database"
	"github.com/plugimt/transat-backend/services"
)

// Needs a disposable Postgres: TEST_DATABASE_URL=postgres://user@host/db?sslmode=disable go test ./handlers/user
func TestNotificationPreferencesEndpoints(t *testing.T) {
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("TEST_DATABASE_URL not set")
	}
	db, err := database.Open(dsn, "../../db/migrations")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	const email = "endpoint.prefs@test.fr"
	for _, q := range []string{
		`DELETE FROM newf WHERE email = $1`,
		`INSERT INTO newf (email, password, first_name, last_name) VALUES ($1, 'x', 'T', 'T')`,
		`INSERT INTO notifications (email, id_services) SELECT $1, id_services FROM services`,
	} {
		if _, err := db.Exec(q, email); err != nil {
			t.Fatal(err)
		}
	}
	defer db.Exec(`DELETE FROM newf WHERE email = $1`, email)

	h := NewUserHandler(db, services.NewNotificationService(db))
	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("email", email)
		return c.Next()
	})
	app.Get("/prefs", h.GetNotificationPreferences)
	app.Put("/prefs", h.SetNotificationPreference)

	call := func(method, body string) (int, string) {
		t.Helper()
		req := httptest.NewRequest(method, "/prefs", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		if err != nil {
			t.Fatal(err)
		}
		raw, _ := io.ReadAll(resp.Body)
		return resp.StatusCode, string(raw)
	}

	if status, body := call("PUT", `{"service":"RESTAURANT","enabled":false}`); status != 200 || body != `{"service":"RESTAURANT","enabled":false}` {
		t.Fatalf("PUT = %d %s", status, body)
	}

	status, body := call("GET", "")
	if status != 200 {
		t.Fatalf("GET = %d %s", status, body)
	}
	var got struct {
		Preferences []struct {
			Service string `json:"service"`
			Enabled bool   `json:"enabled"`
		} `json:"preferences"`
	}
	if err := json.Unmarshal([]byte(body), &got); err != nil {
		t.Fatal(err)
	}
	state := map[string]bool{}
	for _, p := range got.Preferences {
		state[p.Service] = p.Enabled
	}
	if len(state) != 5 || state["RESTAURANT"] || !state["EVENTS"] || !state["RESERVATIONS"] {
		t.Fatalf("unexpected preferences: %s", body)
	}

	for _, bad := range []string{`{"service":"NOPE","enabled":true}`, `{"service":"EVENTS"}`, `not json`} {
		if status, _ := call("PUT", bad); status != 400 {
			t.Fatalf("PUT %s = %d, want 400", bad, status)
		}
	}
}
