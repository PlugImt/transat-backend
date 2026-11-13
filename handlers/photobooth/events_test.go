package photobooth_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	photoboothHandler "github.com/plugimt/transat-backend/handlers/photobooth"
	"github.com/plugimt/transat-backend/internal/database/dbtest"
	"github.com/plugimt/transat-backend/routes"
	pb "github.com/plugimt/transat-backend/services/photobooth"
)

var secret = bytes.Repeat([]byte("k"), 32)

// boothRequest builds a request signed with booth-1's key, as the tablet does.
func boothRequest(t *testing.T, app *fiber.App, method, path string, body []byte, nonce string) (int, []byte) {
	t.Helper()
	ts := strconv.FormatInt(time.Now().Unix(), 10)
	req := httptest.NewRequest(method, path, bytes.NewReader(body))
	req.Header.Set(pb.HeaderKey, "booth-1")
	req.Header.Set(pb.HeaderTimestamp, ts)
	req.Header.Set(pb.HeaderNonce, nonce)
	req.Header.Set(pb.HeaderSignature, pb.Sign(secret, pb.Canonical(method, path, ts, nonce, body)))
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	res, err := app.Test(req, -1)
	if err != nil {
		t.Fatal(err)
	}
	out, _ := io.ReadAll(res.Body)
	return res.StatusCode, out
}

func TestBoothPublishesLayoutsAndFetchesEvents(t *testing.T) {
	db := dbtest.Open(t)
	svc := pb.NewService(db, nil, "root", "https://transat.dev/photobooth")
	h := photoboothHandler.NewHandler(svc, pb.NewVerifier(map[string][]byte{"booth-1": secret}))
	app := fiber.New()
	routes.SetupPhotoboothRoutes(app, db, h)

	catalog := []byte(`{"aspects":{"RATIO_4_3":[{"type":"PORTRAIT_DUO","title":"Portrait Duo","width":1200,"height":1800,
		"title_zone":{"x":0,"y":1440,"w":1200,"h":360},"photos":[{"x":150,"y":48,"w":900,"h":675,"shot":1}],"enabled_by_default":true}]}}`)
	if code, body := boothRequest(t, app, "PUT", "/photobooth/booth/layouts", catalog, "nonce-0000000000000001"); code != 200 {
		t.Fatalf("publish: %d %s", code, body)
	}
	if _, err := svc.CreateEvent(t.Context(), "a@b", pb.EventSettings{Name: "Gala", Date: time.Now().Format("2006-01-02"),
		PhotoAspect: "RATIO_4_3", CountdownSeconds: 5, MaxPrintsPerSession: 5, EnabledLayouts: []string{"PORTRAIT_DUO"}}); err != nil {
		t.Fatal(err)
	}

	code, body := boothRequest(t, app, "GET", "/photobooth/booth/events", nil, "nonce-0000000000000002")
	if code != 200 {
		t.Fatalf("events: %d %s", code, body)
	}
	var got struct {
		Events []pb.BoothEvent `json:"events"`
	}
	if err := json.Unmarshal(body, &got); err != nil || len(got.Events) != 1 || got.Events[0].Name != "Gala" ||
		len(got.Events[0].Layouts) != 1 || !got.Events[0].Layouts[0].Enabled {
		t.Fatalf("events %s (%v)", body, err)
	}

	// Unsigned, or replayed: refused.
	if res, _ := app.Test(httptest.NewRequest("GET", "/photobooth/booth/events", nil), -1); res.StatusCode != 401 {
		t.Fatalf("unsigned: %d", res.StatusCode)
	}
	if code, _ := boothRequest(t, app, "GET", "/photobooth/booth/events", nil, "nonce-0000000000000002"); code != 401 {
		t.Fatalf("replayed nonce: %d", code)
	}
	// The admin routes need a user token.
	if res, _ := app.Test(httptest.NewRequest("GET", "/admin/photobooth/events", nil), -1); res.StatusCode != 401 {
		t.Fatalf("admin route without token: %d", res.StatusCode)
	}
}
