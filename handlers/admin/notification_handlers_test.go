package admin

import (
	"io"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/plugimt/transat-backend/internal/database"
	"github.com/plugimt/transat-backend/services"
)

// Needs a disposable Postgres: TEST_DATABASE_URL=postgres://user@host/db?sslmode=disable go test ./handlers/admin
func TestSendNotificationDryRunCountsRecipients(t *testing.T) {
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("TEST_DATABASE_URL not set")
	}
	db, err := database.Open(dsn, "../../db/migrations")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	const email = "admin.notif@test.fr"
	for _, q := range []string{
		`DELETE FROM newf WHERE email = $1`,
		`INSERT INTO newf (email, password, first_name, last_name) VALUES ($1, 'x', 'T', 'T')`,
		`INSERT INTO user_notification_tokens (email, token) VALUES ($1, 'ExponentPushToken[a]'), ($1, 'ExponentPushToken[b]')`,
	} {
		if _, err := db.Exec(q, email); err != nil {
			t.Fatal(err)
		}
	}
	defer db.Exec(`DELETE FROM newf WHERE email = $1`, email)

	app := fiber.New()
	app.Post("/send", (&AdminHandler{Notifications: services.NewNotificationService(db)}).SendNotification)

	body := `{"title":"t","audience":{"type":"users","emails":["` + email + `"]},"dryRun":true}`
	req := httptest.NewRequest("POST", "/send", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	raw, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 || string(raw) != `{"devices":2,"sent":false,"users":1}` {
		t.Fatalf("dry run = %d %s", resp.StatusCode, raw)
	}

	// Opted out of a category: nobody left to notify, and a real send is refused.
	body = `{"title":"t","audience":{"type":"users","emails":["` + email + `"]},"service":"EVENTS"}`
	req = httptest.NewRequest("POST", "/send", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err = app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 422 {
		t.Fatalf("send without opted-in recipients = %d, want 422", resp.StatusCode)
	}
}

// Invalid requests are rejected before the notification service is touched.
func TestSendNotificationRejectsInvalidRequests(t *testing.T) {
	app := fiber.New()
	app.Post("/send", (&AdminHandler{}).SendNotification)

	cases := map[string]string{
		"not json":         `nope`,
		"missing title":    `{"title":" ","audience":{"type":"all"}}`,
		"unknown audience": `{"title":"t","audience":{"type":"everyone"}}`,
		"club without id":  `{"title":"t","audience":{"type":"club"}}`,
		"unknown campus":   `{"title":"t","audience":{"type":"campus","campus":"MARS"}}`,
		"users empty":      `{"title":"t","audience":{"type":"users"}}`,
		"unknown service":  `{"title":"t","audience":{"type":"all"},"service":"NOPE"}`,
		"bad navigation":   `{"title":"t","audience":{"type":"all"},"navigation":{"type":"event"}}`,
		"title too long":   `{"title":"` + strings.Repeat("a", 101) + `","audience":{"type":"all"}}`,
	}

	for name, body := range cases {
		t.Run(name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/send", strings.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			resp, err := app.Test(req)
			if err != nil {
				t.Fatal(err)
			}
			raw, _ := io.ReadAll(resp.Body)
			if resp.StatusCode != 400 {
				t.Fatalf("status = %d (%s), want 400", resp.StatusCode, raw)
			}
		})
	}
}
