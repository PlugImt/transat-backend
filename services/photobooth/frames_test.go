package photobooth

import (
	"bytes"
	"context"
	"image"
	"image/png"
	"strings"
	"testing"
	"time"
)

func pngOf(t *testing.T, w, h int) []byte {
	t.Helper()
	var b bytes.Buffer
	if err := png.Encode(&b, image.NewNRGBA(image.Rect(0, 0, w, h))); err != nil {
		t.Fatal(err)
	}
	return b.Bytes()
}

func testLayouts() []FrameLayout {
	return []FrameLayout{
		{Type: "PORTRAIT_DUO", Title: "Portrait Duo", Width: 1200, Height: 1800,
			TitleZone: &FrameZone{X: 0, Y: 1440, W: 1200, H: 360},
			Photos:    []FrameZone{{X: 150, Y: 48, W: 900, H: 675, Shot: 1}, {X: 150, Y: 753, W: 900, H: 675, Shot: 2}}},
		{Type: "LANDSCAPE_BANNER", Title: "Bannière", Width: 1800, Height: 1200,
			Photos: []FrameZone{{X: 276, Y: 48, W: 1248, H: 936, Shot: 1}}},
	}
}

func TestFramesTravelFromTheLinkToTheBooth(t *testing.T) {
	svc, _ := setup(t)
	ctx := context.Background()
	session, err := svc.CreateFrameSession(ctx, "booth-1", "Soirée Égypte", 0, testLayouts())
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(session.URL, "https://transat.dev/photobooth/frames/") {
		t.Fatalf("url %s", session.URL)
	}
	token := strings.TrimPrefix(session.URL, "https://transat.dev/photobooth/frames/")

	page, err := svc.PublicFrameSession(ctx, token)
	if err != nil || page.EventName != "Soirée Égypte" || len(page.Layouts) != 2 || len(page.Uploaded) != 0 {
		t.Fatalf("page %+v %v", page, err)
	}

	// Same shape at a smaller size is fine (the booth scales it); another shape or format is not.
	if _, err := svc.UploadFrame(ctx, token, "PORTRAIT_DUO", pngOf(t, 600, 900)); err != nil {
		t.Fatalf("half-size frame: %v", err)
	}
	if _, err := svc.UploadFrame(ctx, token, "PORTRAIT_DUO", pngOf(t, 1800, 1200)); status(err) != 400 {
		t.Fatalf("landscape frame for a portrait layout: %v", err)
	}
	if _, err := svc.UploadFrame(ctx, token, "PORTRAIT_DUO", []byte("\xff\xd8\xff\xe0 not a png")); status(err) != 400 {
		t.Fatalf("jpeg accepted: %v", err)
	}
	if _, err := svc.UploadFrame(ctx, token, "PORTRAIT_DUO", pngOf(t, 300, 450)); status(err) != 400 {
		t.Fatalf("tiny frame accepted: %v", err)
	}
	if _, err := svc.UploadFrame(ctx, token, "GRID_2X2", pngOf(t, 1800, 1200)); status(err) != 404 {
		t.Fatalf("layout not in the session: %v", err)
	}
	final := pngOf(t, 1200, 1800)
	if _, err := svc.UploadFrame(ctx, token, "PORTRAIT_DUO", final); err != nil {
		t.Fatalf("exact frame: %v", err)
	}

	uploads, _, err := svc.FrameUploads(ctx, "booth-1", session.ID)
	if err != nil || len(uploads) != 1 || uploads[0].Width != 1200 || uploads[0].LayoutType != "PORTRAIT_DUO" {
		t.Fatalf("uploads %+v %v", uploads, err)
	}
	got, err := svc.FrameImage(ctx, "booth-1", session.ID, "PORTRAIT_DUO")
	if err != nil || !bytes.Equal(got, final) {
		t.Fatalf("image: %v (the latest upload replaces the previous one)", err)
	}

	// Another booth sees nothing of it.
	if _, _, err := svc.FrameUploads(ctx, "booth-2", session.ID); status(err) != 404 {
		t.Fatalf("other booth: %v", err)
	}
	if _, err := svc.FrameImage(ctx, "booth-2", session.ID, "PORTRAIT_DUO"); status(err) != 404 {
		t.Fatalf("other booth image: %v", err)
	}

	// Expired, the link stops accepting; an hour later it is gone.
	svc.now = func() time.Time { return time.Now().Add(FrameSessionTTL + time.Minute) }
	if _, err := svc.UploadFrame(ctx, token, "LANDSCAPE_BANNER", pngOf(t, 1800, 1200)); status(err) != 404 {
		t.Fatalf("upload after expiry: %v", err)
	}
	svc.now = func() time.Time { return time.Now().Add(FrameSessionTTL + 2*time.Hour) }
	if n, err := svc.PurgeFrameSessions(ctx); err != nil || n != 1 {
		t.Fatalf("purge: %d %v", n, err)
	}
}

func TestFrameSessionsRejectBadGeometry(t *testing.T) {
	svc, _ := setup(t)
	ctx := context.Background()
	bad := testLayouts()
	bad[0].Photos[0].W = 5000
	if _, err := svc.CreateFrameSession(ctx, "booth-1", "Gala", 0, bad); status(err) != 400 {
		t.Fatalf("zone outside the canvas: %v", err)
	}
	dup := append(testLayouts(), testLayouts()[0])
	if _, err := svc.CreateFrameSession(ctx, "booth-1", "Gala", 0, dup); status(err) != 400 {
		t.Fatalf("repeated layout: %v", err)
	}
	if _, err := svc.CreateFrameSession(ctx, "booth-1", "  ", 0, testLayouts()); status(err) != 400 {
		t.Fatalf("no event name: %v", err)
	}
}
