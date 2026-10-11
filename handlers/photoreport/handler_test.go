package photoreport_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	reportHandler "github.com/plugimt/transat-backend/handlers/photoreport"
	"github.com/plugimt/transat-backend/internal/database/dbtest"
	"github.com/plugimt/transat-backend/middlewares"
	"github.com/plugimt/transat-backend/models"
	"github.com/plugimt/transat-backend/routes"
	service "github.com/plugimt/transat-backend/services/photoreport"
	"github.com/plugimt/transat-backend/utils"
)

type mailer struct {
	mu   sync.Mutex
	to   []string
	fail bool
}

func (m *mailer) SendEmail(mail models.Email, _ interface{}) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.fail {
		return errors.New("provider down")
	}
	m.to = append(m.to, mail.Recipient)
	return nil
}

func (m *mailer) count() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.to)
}

type api struct {
	t   *testing.T
	db  *sql.DB
	app *fiber.App
	m   *mailer
}

const photoID = "PHOTO-0000000001"

func newAPI(t *testing.T) *api {
	t.Setenv("JWT_SECRET", "test-secret-test-secret-test-secret")
	db := dbtest.Open(t)
	middlewares.InitJWTMiddleware(db)
	m := &mailer{}
	app := fiber.New()
	routes.SetupPhotosRoutes(app, db, nil, reportHandler.NewHandler(service.NewService(db, m, "")))
	a := &api{t: t, db: db, app: app, m: m}

	a.user("owner@x.test", "NEWF", "clubphoto_respo")
	a.user("admin@x.test", "NEWF", "ADMIN")
	a.user("student@x.test", "NEWF")
	a.user("teacher@x.test", "ACADEMICS")
	a.exec(`INSERT INTO clubs (name, picture) VALUES ('Club Photo', 'p')`)
	a.exec(`INSERT INTO photo_albums (root_folder_id, title, raw_folder_name, drive_url, visibility, confidence, rule, last_seen_run, club_id)
		VALUES ('r', 'Gala', 'Gala', 'u', 'published', 1, 'R', 1, (SELECT id_clubs FROM clubs))`)
	a.exec(`INSERT INTO photos (drive_file_id, album_id, name, mime_type, media_type, sort_at, sort_at_source, web_url, last_seen_run)
		VALUES ($1, (SELECT id_photo_albums FROM photo_albums), 'a.jpg', 'image/jpeg', 'image', now(), 'x', 'https://drive.example/f', 1)`, photoID)
	return a
}

func (a *api) exec(q string, args ...any) {
	a.t.Helper()
	if _, err := a.db.Exec(q, args...); err != nil {
		a.t.Fatalf("%s: %v", q, err)
	}
}

func (a *api) user(email string, roles ...string) {
	a.exec(`INSERT INTO newf (email, password, first_name, last_name) VALUES ($1, 'x', $2, 'Doe')`, email, strings.Split(email, "@")[0])
	for _, r := range roles {
		a.exec(`INSERT INTO roles (name) VALUES ($1) ON CONFLICT (name) DO NOTHING`, r)
		a.exec(`INSERT INTO newf_roles (email, id_roles) SELECT $1, id_roles FROM roles WHERE name = $2`, email, r)
	}
}

type reply struct {
	status int
	header http.Header
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
	return reply{res.StatusCode, res.Header, string(raw)}
}

func (a *api) report(as string, reason, comment string) reply {
	return a.call(http.MethodPost, "/photos/reports", as, map[string]any{"photo_id": photoID, "reason": reason, "comment": comment})
}

func (a *api) waitMails(n int) {
	a.t.Helper()
	deadline := time.Now().Add(5 * time.Second)
	for a.m.count() < n && time.Now().Before(deadline) {
		time.Sleep(20 * time.Millisecond)
	}
	if got := a.m.count(); got != n {
		a.t.Fatalf("%d emails, want %d", got, n)
	}
}

func TestSubmitRequiresAStudentSession(t *testing.T) {
	a := newAPI(t)
	for name, c := range map[string]struct {
		as   string
		want int
	}{
		"no token":      {"", http.StatusUnauthorized},
		"garbage token": {"not-a-token", http.StatusUnauthorized},
		"academic":      {"teacher@x.test", http.StatusForbidden},
	} {
		if got := a.report(c.as, "other", ""); got.status != c.want {
			t.Fatalf("%s: status %d, want %d (%s)", name, got.status, c.want, got.body)
		}
	}
	var n int
	_ = a.db.QueryRow(`SELECT count(*) FROM photo_reports`).Scan(&n)
	if n != 0 {
		t.Fatalf("%d reports stored without a valid student", n)
	}
}

func TestSubmitRejectsInvalidInput(t *testing.T) {
	a := newAPI(t)
	post := func(body any) reply { return a.call(http.MethodPost, "/photos/reports", "student@x.test", body) }
	cases := map[string]struct {
		body any
		want int
	}{
		"unknown reason": {map[string]any{"photo_id": photoID, "reason": "spam"}, 400},
		"empty reason":   {map[string]any{"photo_id": photoID}, 400},
		"short photo id": {map[string]any{"photo_id": "x", "reason": "other"}, 400},
		"missing photo":  {map[string]any{"reason": "other"}, 400},
		"unknown photo":  {map[string]any{"photo_id": "UNKNOWN-0000000", "reason": "other"}, 404},
		"comment length": {map[string]any{"photo_id": photoID, "reason": "other", "comment": strings.Repeat("a", 1001)}, 400},
		"not json":       {"{nope", 400},
		"wrong types":    {map[string]any{"photo_id": 5, "reason": []int{1}}, 400},
	}
	for name, c := range cases {
		if got := post(c.body); got.status != c.want {
			t.Fatalf("%s: status %d, want %d (%s)", name, got.status, c.want, got.body)
		}
	}
}

func TestEveryReasonIsAcceptedAndMailed(t *testing.T) {
	for _, reason := range []string{"appears_in_photo", "inappropriate", "poor_quality", "other"} {
		t.Run(reason, func(t *testing.T) {
			a := newAPI(t)
			got := a.report("student@x.test", reason, "")
			if got.status != http.StatusCreated {
				t.Fatalf("status %d: %s", got.status, got.body)
			}
			a.waitMails(2) // the club owner and the administrator
			if strings.Contains(got.body, "@") {
				t.Fatalf("the answer must not expose any address: %s", got.body)
			}
		})
	}
}

func TestReporterAndRecipientsCannotBeChosenByTheClient(t *testing.T) {
	a := newAPI(t)
	got := a.call(http.MethodPost, "/photos/reports", "student@x.test", map[string]any{
		"photo_id": photoID, "reason": "inappropriate", "comment": "x",
		"reporter_email": "admin@x.test", "reporter_id": 1, "recipients": []string{"victim@evil.test"},
		"club_id": 999, "to": "victim@evil.test", "cc": []string{"victim@evil.test"},
	})
	if got.status != http.StatusCreated {
		t.Fatalf("status %d: %s", got.status, got.body)
	}
	a.waitMails(2)
	var reporter, club string
	_ = a.db.QueryRow(`SELECT reporter_email, club_name FROM photo_reports`).Scan(&reporter, &club)
	if reporter != "student@x.test" || club != "Club Photo" {
		t.Fatalf("reporter %q club %q", reporter, club)
	}
	for _, to := range a.m.to {
		if to != "owner@x.test" && to != "admin@x.test" {
			t.Fatalf("mail sent to %s", to)
		}
	}
}

func TestDuplicateRateLimitAndEmailFailure(t *testing.T) {
	a := newAPI(t)
	a.m.fail = true // the email service is down: the user still gets a confirmation and nothing is lost

	first := a.report("student@x.test", "other", "hello")
	if first.status != http.StatusCreated {
		t.Fatalf("first: %d %s", first.status, first.body)
	}
	again := a.report("student@x.test", "other", "hello")
	if again.status != http.StatusOK || !strings.Contains(again.body, `"duplicate":true`) {
		t.Fatalf("double tap: %d %s", again.status, again.body)
	}
	time.Sleep(200 * time.Millisecond) // background delivery attempt
	var reports, failed int
	_ = a.db.QueryRow(`SELECT count(*) FROM photo_reports`).Scan(&reports)
	_ = a.db.QueryRow(`SELECT count(*) FROM photo_report_deliveries WHERE status = 'failed' AND last_error IS NOT NULL`).Scan(&failed)
	if reports != 1 || failed != 2 {
		t.Fatalf("%d reports and %d failed deliveries; want 1 stored report with 2 failed deliveries to retry", reports, failed)
	}

	var last reply
	for i := 0; i < 8; i++ {
		last = a.report("student@x.test", "other", fmt.Sprint("different ", i))
	}
	if last.status != http.StatusTooManyRequests || last.header.Get("Retry-After") == "" {
		t.Fatalf("flood: %d %v", last.status, last.header)
	}
	if got := a.report("admin@x.test", "other", "someone else"); got.status != http.StatusCreated {
		t.Fatalf("another user is not limited: %d", got.status)
	}
}

func TestAdminReportsAreAdministratorOnly(t *testing.T) {
	a := newAPI(t)
	a.report("student@x.test", "appears_in_photo", "c")
	a.waitMails(2)

	for _, path := range []string{"/admin/photos/reports", "/admin/photos/reports/1"} {
		for _, who := range []string{"", "student@x.test", "owner@x.test", "teacher@x.test"} {
			if got := a.call(http.MethodGet, path, who, nil); got.status != 401 && got.status != 403 {
				t.Fatalf("%s as %q: status %d", path, who, got.status)
			}
		}
	}
	if got := a.call(http.MethodPost, "/admin/photos/reports/1/retry", "student@x.test", nil); got.status != 403 {
		t.Fatalf("retry as student: %d", got.status)
	}
	if got := a.call(http.MethodPut, "/admin/photos/albums/1/club", "owner@x.test", map[string]any{"club_id": nil}); got.status != 403 {
		t.Fatalf("club link as club owner: %d", got.status)
	}

	list := a.call(http.MethodGet, "/admin/photos/reports", "admin@x.test", nil)
	var parsed struct {
		Total   int `json:"total"`
		Reports []struct {
			Reason   string `json:"reason"`
			PhotoURL string `json:"photo_url"`
			Reporter struct {
				Email string `json:"email"`
			} `json:"reporter"`
			DeliveryStatus string `json:"delivery_status"`
		} `json:"reports"`
	}
	if err := json.Unmarshal([]byte(list.body), &parsed); err != nil || list.status != 200 {
		t.Fatalf("list: %d %s", list.status, list.body)
	}
	if parsed.Total != 1 || parsed.Reports[0].Reporter.Email != "student@x.test" || parsed.Reports[0].PhotoURL != "https://drive.example/f" ||
		parsed.Reports[0].DeliveryStatus != "sent" {
		t.Fatalf("list content: %s", list.body)
	}

	for path, want := range map[string]int{
		"/admin/photos/reports?reason=nope":     400,
		"/admin/photos/reports?delivery=nope":   400,
		"/admin/photos/reports?reason=other":    200,
		"/admin/photos/reports/abc":             400,
		"/admin/photos/reports/999":             404,
		"/admin/photos/reports/1":               200,
		"/admin/photos/reports?delivery=failed": 200,
	} {
		if got := a.call(http.MethodGet, path, "admin@x.test", nil); got.status != want {
			t.Fatalf("%s: %d, want %d", path, got.status, want)
		}
	}
	if got := a.call(http.MethodPost, "/admin/photos/reports/999/retry", "admin@x.test", nil); got.status != 404 {
		t.Fatalf("retry unknown: %d", got.status)
	}
	if got := a.call(http.MethodPost, "/admin/photos/reports/1/retry", "admin@x.test", nil); got.status != 200 {
		t.Fatalf("retry: %d", got.status)
	}
	if got := a.call(http.MethodPut, "/admin/photos/albums/1/club", "admin@x.test", map[string]any{"club_id": nil}); got.status != 204 {
		t.Fatalf("unlink: %d", got.status)
	}
	if got := a.call(http.MethodPut, "/admin/photos/albums/1/club", "admin@x.test", map[string]any{"club_id": 9999}); got.status != 400 {
		t.Fatalf("unknown club: %d", got.status)
	}
	if got := a.call(http.MethodPut, "/admin/photos/albums/999/club", "admin@x.test", map[string]any{"club_id": nil}); got.status != 404 {
		t.Fatalf("unknown album: %d", got.status)
	}

	// With the album unlinked from its club, the owner no longer hears about its photos.
	before := a.m.count()
	a.report("admin@x.test", "poor_quality", "after unlink")
	a.waitMails(before + 1)
}
