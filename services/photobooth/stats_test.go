package photobooth

import (
	"context"
	"testing"
	"time"
)

func TestStatsCountPeopleWhoScannedAndOpenedTheAlbum(t *testing.T) {
	svc, _ := setup(t)
	ctx := context.Background()
	svc.now = func() time.Time { return time.Date(2026, 3, 30, 12, 0, 0, 0, time.UTC) }
	exec := func(q string, args ...any) {
		t.Helper()
		if _, err := svc.db.ExecContext(ctx, q, args...); err != nil {
			t.Fatal(err)
		}
	}
	token := func(i int) string { return "token-" + string(rune('a'+i)) + "-0123456789abcdefghijklmnopqrstu" }
	// Gala: 3 shared sessions (names differing by case and accents), one cancelled. Bal: 1 session.
	for i, ev := range []struct{ name, date, status string }{
		{"Gala", "2026-03-27", "finalized"}, {"gala", "2026-03-27", "finalized"}, {"Galà ", "2026-03-27", "finalized"},
		{"Gala", "2026-03-27", "cancelled"}, {"Bal", "2026-03-20", "finalized"},
	} {
		exec(`INSERT INTO photobooth_shoots (device_key_id, client_session_id, event_name, event_date, drive_folder_id, status, share_token)
			VALUES ('booth-1', $1, $2, $3, 'f', $4, $5)`, "s"+token(i), ev.name, ev.date, ev.status, token(i))
	}
	if _, err := svc.CreateEvent(ctx, "a@b", eventSettings("Gala", "2026-03-27", "PORTRAIT_DUO")); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.CreateEvent(ctx, "a@b", eventSettings("Futur", "2026-04-27", "PORTRAIT_DUO")); err != nil {
		t.Fatal(err)
	}
	exec(`INSERT INTO photobooth_event_folders (event_date, event_key, event_folder_id, photobooth_folder) VALUES ('2026-03-27', 'gala', 'drive-gala', 'p')`)
	var album int64
	if err := svc.db.QueryRowContext(ctx, `INSERT INTO photo_albums (root_folder_id, title, raw_folder_name, drive_url, visibility, confidence, rule, last_seen_run)
		VALUES ('drive-gala', 'Gala 2026', '2026.03.27 - Gala', 'u', 'published', 1, 'r', 1) RETURNING id_photo_albums`).Scan(&album); err != nil {
		t.Fatal(err)
	}

	alice, bob, carol := "visitor-alice-0000000001", "visitor-bob-000000000002", "visitor-carol-0000000003"
	visits := []struct {
		token, visitor string
	}{{token(0), alice}, {token(0), alice}, {token(0), bob}, {token(1), alice}, {token(4), carol}}
	for _, v := range visits {
		if err := svc.RecordShareVisit(ctx, v.token, v.visitor); err != nil {
			t.Fatal(err)
		}
	}
	for _, v := range []string{alice, bob, bob} {
		if err := svc.RecordAlbumVisit(ctx, album, v); err != nil {
			t.Fatal(err)
		}
	}
	// Not counted: a cancelled shoot, an unknown album, a malformed visitor.
	if status(svc.RecordShareVisit(ctx, token(3), alice)) != 404 || status(svc.RecordAlbumVisit(ctx, album+1, alice)) != 404 ||
		status(svc.RecordShareVisit(ctx, token(0), "short")) != 404 {
		t.Fatal("bad visits recorded")
	}

	st, err := svc.Stats(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(st.Events) != 2 || st.TrackingSince == nil {
		t.Fatalf("events %+v", st)
	}
	gala, bal := st.Events[0], st.Events[1]
	if gala.Name != "Gala" || gala.EventID == nil || gala.Sessions != 3 || gala.ScannedSessions != 2 || gala.Scanners != 2 ||
		gala.ScanViews != 4 || gala.AlbumID == nil || *gala.AlbumID != album || gala.AlbumVisitors != 2 || gala.AlbumViews != 3 {
		t.Fatalf("gala %+v", gala)
	}
	if bal.Name != "Bal" || bal.EventID != nil || bal.Sessions != 1 || bal.Scanners != 1 || bal.AlbumID != nil {
		t.Fatalf("bal %+v", bal)
	}
	tot := st.Totals
	if tot.Events != 2 || tot.Sessions != 4 || tot.ScannedSessions != 3 || tot.Scanners != 3 || tot.ScanViews != 5 ||
		tot.AlbumVisitors != 2 || tot.AvgScannersPerEvent != 1.5 || tot.ScanRate != 0.75 || tot.AvgScannersPerSession != 0.75 {
		t.Fatalf("totals %+v", tot)
	}
}
