package faces

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	photosHandler "github.com/plugimt/transat-backend/handlers/photos"
	"github.com/plugimt/transat-backend/internal/database/dbtest"
	svc "github.com/plugimt/transat-backend/services/faces"
)

// The fake model sees one face, the same for every image: the selfie and the photo below.
func fakeModel(t *testing.T) *httptest.Server {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"facial-recognition":[{"boundingBox":{"x1":200,"y1":100,"x2":600,"y2":600},
			"embedding":"[1,0,0,0]","score":0.99}],"imageWidth":1000,"imageHeight":1000}`))
	}))
	t.Cleanup(srv.Close)
	return srv
}

type images struct{}

func (images) Image(_ context.Context, id string, _ int) ([]byte, error) { return []byte(id), nil }

func TestFaceRoutes(t *testing.T) {
	db := dbtest.Open(t)
	for _, u := range []string{"alice@imt.fr", "bob@imt.fr"} {
		if _, err := db.Exec(`INSERT INTO newf (email, password, first_name, last_name) VALUES ($1, 'x', 'Jo', 'Doe')`, u); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := db.Exec(`INSERT INTO photos (drive_file_id, name, mime_type, media_type, md5, sort_at, sort_at_source, web_url, last_seen_run)
		VALUES ('GALA-PHOTO-1', 'a.jpg', 'image/jpeg', 'image', 'm', now(), 'x', 'u', 1)`); err != nil {
		t.Fatal(err)
	}
	cfg := svc.DefaultConfig()
	cfg.TickBudget = 0
	s := svc.NewService(db, svc.NewMLClient(fakeModel(t).URL, "buffalo_l", 0.7), images{}, cfg)

	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error { c.Locals("email", c.Get("X-Test-User")); return c.Next() })
	h := NewHandler(s)
	app.Get("/faces", h.GetStatus)
	app.Post("/faces/consent", h.PostConsent)
	app.Post("/faces/references", h.PostReference)
	app.Delete("/faces", h.Delete)
	app.Post("/faces/photos/:id/reject", h.RejectPhoto)
	app.Get("/photos/me", photosHandler.NewPhotosHandler(db, nil).GetMyPhotos)

	call := func(user, method, path string, body io.Reader, contentType string) (int, map[string]any) {
		t.Helper()
		req := httptest.NewRequest(method, path, body)
		req.Header.Set("X-Test-User", user)
		if contentType != "" {
			req.Header.Set("Content-Type", contentType)
		}
		res, err := app.Test(req, -1)
		if err != nil {
			t.Fatal(err)
		}
		out := map[string]any{}
		_ = json.NewDecoder(res.Body).Decode(&out)
		return res.StatusCode, out
	}
	selfie := func(user string) (int, map[string]any) {
		var b bytes.Buffer
		w := multipart.NewWriter(&b)
		part, _ := w.CreateFormFile("photo", "selfie.jpg")
		_, _ = part.Write([]byte("selfie"))
		_ = w.Close()
		return call(user, "POST", "/faces/references", &b, w.FormDataContentType())
	}
	jsonBody := func(v string) io.Reader { return strings.NewReader(v) }

	if code, st := call("alice@imt.fr", "GET", "/faces", nil, ""); code != 200 || st["available"] != true || st["enabled"] != false {
		t.Fatalf("status before consent: %d %v", code, st)
	}
	if code, body := selfie("alice@imt.fr"); code != 403 || body["code"] != "consent_required" {
		t.Fatalf("selfie before consent: %d %v", code, body)
	}
	if code, _ := call("alice@imt.fr", "POST", "/faces/consent", jsonBody(`{"version":"`+svc.ConsentVersion+`","accepted":false}`), "application/json"); code != 400 {
		t.Fatalf("consent not accepted: %d", code)
	}
	if code, body := call("alice@imt.fr", "POST", "/faces/consent", jsonBody(`{"version":"old","accepted":true}`), "application/json"); code != 409 || body["code"] != "consent_outdated" {
		t.Fatalf("old consent text: %d %v", code, body)
	}
	if code, st := call("alice@imt.fr", "POST", "/faces/consent", jsonBody(`{"version":"`+svc.ConsentVersion+`","accepted":true}`), "application/json"); code != 200 || st["enabled"] != true {
		t.Fatalf("consent: %d %v", code, st)
	}
	if code, st := selfie("alice@imt.fr"); code != 201 || st["references"] != float64(1) {
		t.Fatalf("selfie: %d %v", code, st)
	}
	if _, err := s.Tick(context.Background()); err != nil {
		t.Fatal(err)
	}

	// Alice sees the photo; Bob, who has the same face for the fake model but never consented, does not.
	if code, page := call("alice@imt.fr", "GET", "/photos/me", nil, ""); code != 200 || len(page["photos"].([]any)) != 1 {
		t.Fatalf("alice album: %d %v", code, page)
	}
	if code, page := call("bob@imt.fr", "GET", "/photos/me", nil, ""); code != 200 || len(page["photos"].([]any)) != 0 {
		t.Fatalf("bob album: %d %v", code, page)
	}
	if code, _ := call("bob@imt.fr", "POST", "/faces/photos/GALA-PHOTO-1/reject", nil, ""); code != 404 {
		t.Fatalf("bob rejecting alice's match: %d", code)
	}
	if code, _ := call("alice@imt.fr", "POST", "/faces/photos/GALA-PHOTO-1/reject", nil, ""); code != 204 {
		t.Fatalf("reject: %d", code)
	}
	if _, page := call("alice@imt.fr", "GET", "/photos/me", nil, ""); len(page["photos"].([]any)) != 0 {
		t.Fatalf("rejected photo still listed: %v", page)
	}
	if code, _ := call("alice@imt.fr", "DELETE", "/faces", nil, ""); code != 204 {
		t.Fatalf("withdraw: %d", code)
	}
	if _, st := call("alice@imt.fr", "GET", "/faces", nil, ""); st["enabled"] != false || st["references"] != float64(0) {
		t.Fatalf("status after withdrawal: %v", st)
	}
}
