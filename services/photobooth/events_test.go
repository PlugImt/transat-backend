package photobooth

import (
	"context"
	"slices"
	"strings"
	"testing"
	"time"
)

func testCatalog() map[string][]CatalogLayout {
	out := map[string][]CatalogLayout{}
	for _, aspect := range []string{"RATIO_4_3", "RATIO_3_2"} {
		for i, l := range testLayouts() {
			out[aspect] = append(out[aspect], CatalogLayout{FrameLayout: l, EnabledByDefault: i == 0})
		}
	}
	return out
}

func eventSettings(name, date string, layouts ...string) EventSettings {
	s := DefaultEventSettings()
	s.Name, s.Date, s.EnabledLayouts = name, date, layouts
	return s
}

func TestBoothPublishesItsLayoutCatalog(t *testing.T) {
	svc, _ := setup(t)
	ctx := context.Background()
	empty, err := svc.LayoutCatalog(ctx)
	if err != nil || len(empty.Aspects) != 0 || empty.UpdatedAt != nil {
		t.Fatalf("empty catalog %+v %v", empty, err)
	}
	catalog := testCatalog()
	catalog["RATIO_4_3"][0].HasFrame = true // booth-local state, never kept
	if err := svc.PublishLayoutCatalog(ctx, "booth-1", catalog); err != nil {
		t.Fatal(err)
	}
	got, err := svc.LayoutCatalog(ctx)
	if err != nil || len(got.Aspects) != 2 || got.Device != "booth-1" || got.UpdatedAt == nil {
		t.Fatalf("catalog %+v %v", got, err)
	}
	first := got.Aspects["RATIO_4_3"][0]
	if first.Type != "PORTRAIT_DUO" || !first.EnabledByDefault || first.HasFrame || len(first.Photos) != 2 || first.TitleZone == nil {
		t.Fatalf("first layout %+v", first)
	}

	// A newer app replaces the catalog as a whole.
	if err := svc.PublishLayoutCatalog(ctx, "booth-2", map[string][]CatalogLayout{"RATIO_16_9": testCatalog()["RATIO_4_3"][1:]}); err != nil {
		t.Fatal(err)
	}
	got, _ = svc.LayoutCatalog(ctx)
	if len(got.Aspects) != 1 || len(got.Aspects["RATIO_16_9"]) != 1 || got.Device != "booth-2" {
		t.Fatalf("replaced catalog %+v", got)
	}

	for name, bad := range map[string]map[string][]CatalogLayout{
		"empty":         {},
		"unknown ratio": {"RATIO_1_1": testCatalog()["RATIO_4_3"]},
		"bad geometry":  {"RATIO_4_3": {{FrameLayout: FrameLayout{Type: "X", Title: "X", Width: 1200, Height: 1800, Photos: []FrameZone{{X: 0, Y: 0, W: 2000, H: 10, Shot: 1}}}}}},
	} {
		if err := svc.PublishLayoutCatalog(ctx, "booth-1", bad); status(err) != 400 {
			t.Fatalf("%s: %v", name, err)
		}
	}
}

func TestEventsAreManagedFromThePanel(t *testing.T) {
	svc, _ := setup(t)
	ctx := context.Background()
	if err := svc.PublishLayoutCatalog(ctx, "booth-1", testCatalog()); err != nil {
		t.Fatal(err)
	}

	e, err := svc.CreateEvent(ctx, "respo@imt-atlantique.net", eventSettings("  Soirée   Égypte ", "2026-11-20", "PORTRAIT_DUO"))
	if err != nil {
		t.Fatal(err)
	}
	if e.Name != "Soirée Égypte" || e.Date != "2026-11-20" || e.PhotoAspect != "RATIO_4_3" || e.CountdownSeconds != 5 ||
		!e.AllowRetakes || e.CreatedBy != "respo@imt-atlantique.net" {
		t.Fatalf("created %+v", e)
	}
	if len(e.Layouts) != 1 || !e.Layouts[0].Enabled || e.Layouts[0].Frame != nil || !slices.Equal(e.EnabledLayouts, []string{"PORTRAIT_DUO"}) {
		t.Fatalf("layouts %+v", e.Layouts)
	}

	// The frame is checked against the catalog's geometry for this layout.
	if _, err := svc.SetEventFrame(ctx, e.ID, "PORTRAIT_DUO", pngOf(t, 1800, 1200)); status(err) != 400 {
		t.Fatalf("landscape frame on a portrait layout: %v", err)
	}
	if _, err := svc.SetEventFrame(ctx, e.ID, "GRID_2X2", pngOf(t, 1800, 1200)); status(err) != 400 {
		t.Fatalf("layout unknown to the catalog: %v", err)
	}
	if _, err := svc.SetEventFrame(ctx, 999999, "PORTRAIT_DUO", pngOf(t, 1200, 1800)); status(err) != 404 {
		t.Fatalf("unknown event: %v", err)
	}
	png := pngOf(t, 1200, 1800)
	frame, err := svc.SetEventFrame(ctx, e.ID, "PORTRAIT_DUO", png)
	if err != nil || frame.Width != 1200 || len(frame.SHA256) != 64 || frame.Size != len(png) {
		t.Fatalf("frame %+v %v", frame, err)
	}
	// A frame on a layout that is switched off is kept, and the layout stays off.
	if _, err := svc.SetEventFrame(ctx, e.ID, "LANDSCAPE_BANNER", pngOf(t, 1800, 1200)); err != nil {
		t.Fatal(err)
	}
	got, _ := svc.Event(ctx, e.ID)
	if len(got.Layouts) != 2 || got.Layouts[0].Type != "LANDSCAPE_BANNER" || got.Layouts[0].Enabled || got.Layouts[0].Frame == nil {
		t.Fatalf("after frames %+v", got.Layouts)
	}
	if b, err := svc.EventFrameImage(ctx, e.ID, "PORTRAIT_DUO"); err != nil || string(b) != string(png) {
		t.Fatalf("image %d bytes, %v", len(b), err)
	}

	// Switching layouts keeps every frame.
	in := eventSettings("Soirée Égypte", "2026-11-21", "LANDSCAPE_BANNER")
	in.PhotoAspect, in.CountdownSeconds, in.AllowRetakes = "RATIO_3_2", 8, false
	got, err = svc.UpdateEvent(ctx, e.ID, in)
	if err != nil || got.Date != "2026-11-21" || got.PhotoAspect != "RATIO_3_2" || got.CountdownSeconds != 8 || got.AllowRetakes {
		t.Fatalf("updated %+v %v", got, err)
	}
	if !slices.Equal(got.EnabledLayouts, []string{"LANDSCAPE_BANNER"}) || got.Layouts[1].Frame == nil || got.Layouts[1].Enabled {
		t.Fatalf("updated layouts %+v", got.Layouts)
	}

	if err := svc.DeleteEventFrame(ctx, e.ID, "PORTRAIT_DUO"); err != nil {
		t.Fatal(err)
	}
	if err := svc.DeleteEventFrame(ctx, e.ID, "PORTRAIT_DUO"); status(err) != 404 {
		t.Fatalf("second delete: %v", err)
	}
	if _, err := svc.EventFrameImage(ctx, e.ID, "PORTRAIT_DUO"); status(err) != 404 {
		t.Fatalf("deleted image: %v", err)
	}

	if err := svc.DeleteEvent(ctx, e.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.Event(ctx, e.ID); status(err) != 404 {
		t.Fatalf("deleted event: %v", err)
	}
	if err := svc.DeleteEvent(ctx, e.ID); status(err) != 404 {
		t.Fatalf("second delete: %v", err)
	}
}

func TestEventSettingsAreValidated(t *testing.T) {
	svc, _ := setup(t)
	ctx := context.Background()
	if err := svc.PublishLayoutCatalog(ctx, "booth-1", testCatalog()); err != nil {
		t.Fatal(err)
	}
	ok := eventSettings("Gala", "2027-02-14", "PORTRAIT_DUO")
	for name, mutate := range map[string]func(*EventSettings){
		"no name":          func(s *EventSettings) { s.Name = "  " },
		"bad date":         func(s *EventSettings) { s.Date = "14/02/2027" },
		"unknown ratio":    func(s *EventSettings) { s.PhotoAspect = "RATIO_1_1" },
		"countdown":        func(s *EventSettings) { s.CountdownSeconds = 0 },
		"prints":           func(s *EventSettings) { s.MaxPrintsPerSession = 21 },
		"no layout":        func(s *EventSettings) { s.EnabledLayouts = nil },
		"repeated layout":  func(s *EventSettings) { s.EnabledLayouts = []string{"PORTRAIT_DUO", "PORTRAIT_DUO"} },
		"layout not known": func(s *EventSettings) { s.EnabledLayouts = []string{"GRID_2X2"} },
		"ratio not known":  func(s *EventSettings) { s.PhotoAspect = "RATIO_16_9" },
	} {
		in := ok
		in.EnabledLayouts = slices.Clone(ok.EnabledLayouts)
		mutate(&in)
		if _, err := svc.CreateEvent(ctx, "a@b", in); status(err) != 400 {
			t.Errorf("%s: %v", name, err)
		}
	}
	if _, err := svc.UpdateEvent(ctx, 999999, ok); status(err) != 404 {
		t.Fatalf("update of an unknown event: %v", err)
	}
	if _, err := svc.CreateEvent(ctx, "a@b", ok); err != nil {
		t.Fatal(err)
	}
}

func TestBoothsSeeRecentAndUpcomingEvents(t *testing.T) {
	svc, _ := setup(t)
	ctx := context.Background()
	svc.now = func() time.Time { return time.Date(2026, 11, 13, 20, 0, 0, 0, time.UTC) }
	for _, d := range []string{"2026-12-01", "2026-10-01", "2026-11-13", "2026-10-20"} {
		if _, err := svc.CreateEvent(ctx, "a@b", eventSettings("Soirée "+d, d, "PORTRAIT_DUO")); err != nil {
			t.Fatal(err)
		}
	}
	list, err := svc.BoothEvents(ctx)
	if err != nil {
		t.Fatal(err)
	}
	var dates []string
	for _, e := range list {
		dates = append(dates, e.Date)
	}
	if strings.Join(dates, ",") != "2026-10-20,2026-11-13,2026-12-01" {
		t.Fatalf("booth sees %v", dates)
	}
	all, _ := svc.Events(ctx)
	if len(all) != 4 || all[0].Date != "2026-12-01" {
		t.Fatalf("panel sees %d events, first %s", len(all), all[0].Date)
	}
}

func TestFramesSentByQRCodeAreFiledInTheEvent(t *testing.T) {
	svc, _ := setup(t)
	ctx := context.Background()
	if err := svc.PublishLayoutCatalog(ctx, "booth-1", testCatalog()); err != nil {
		t.Fatal(err)
	}
	e, err := svc.CreateEvent(ctx, "a@b", eventSettings("Gala", "2027-02-14", "LANDSCAPE_BANNER"))
	if err != nil {
		t.Fatal(err)
	}
	session, err := svc.CreateFrameSession(ctx, "booth-1", "Gala", e.ID, testLayouts())
	if err != nil {
		t.Fatal(err)
	}
	token := session.URL[strings.LastIndex(session.URL, "/")+1:]
	u, err := svc.UploadFrame(ctx, token, "PORTRAIT_DUO", pngOf(t, 600, 900))
	if err != nil {
		t.Fatal(err)
	}
	got, _ := svc.Event(ctx, e.ID)
	var duo *EventLayout
	for i := range got.Layouts {
		if got.Layouts[i].Type == "PORTRAIT_DUO" {
			duo = &got.Layouts[i]
		}
	}
	if duo == nil || !duo.Enabled || duo.Frame == nil || duo.Frame.SHA256 != u.SHA256 {
		t.Fatalf("event layouts after a QR upload %+v", got.Layouts)
	}

	// A link for an event the panel does not know (deleted, or local to the booth) still works.
	orphan, err := svc.CreateFrameSession(ctx, "booth-1", "Local", 424242, testLayouts())
	if err != nil {
		t.Fatal(err)
	}
	if _, err := svc.UploadFrame(ctx, orphan.URL[strings.LastIndex(orphan.URL, "/")+1:], "PORTRAIT_DUO", pngOf(t, 600, 900)); err != nil {
		t.Fatalf("upload to a link without event: %v", err)
	}
}
