package photobooth

import (
	"context"
	"testing"
)

func TestBoothPicturesAreFoundInTheEventAndItsAlbum(t *testing.T) {
	svc, d := setup(t)
	ctx := context.Background()
	b := booth{t, svc, d, "session-0001"}
	b.put(KindPhoto, 1, "p1")
	b.put(KindCollage, 0, "c0")
	ev, err := svc.CreateEvent(ctx, "a@b", eventSettings("soiree integration", event.Date, "PORTRAIT_DUO"))
	if err != nil {
		t.Fatal(err)
	}

	// No album yet (the gallery has not indexed the event folder): the folder is linked to the panel's event only.
	p := pendingFolderOf(t, svc)
	if p.EventID == nil || *p.EventID != ev.ID || p.AlbumID != nil {
		t.Fatalf("pending folder links %+v", p)
	}

	var eventFolder string
	if err := svc.db.QueryRowContext(ctx, `SELECT event_folder_id FROM photobooth_event_folders`).Scan(&eventFolder); err != nil {
		t.Fatal(err)
	}
	var album int64
	if err := svc.db.QueryRowContext(ctx, `INSERT INTO photo_albums (root_folder_id, title, raw_folder_name, drive_url, visibility, confidence, rule, last_seen_run)
		VALUES ($1, 'Soirée Intégration', '2026.10.27 - Soirée Intégration', 'u', 'published', 1, 'r', 1) RETURNING id_photo_albums`, eventFolder).Scan(&album); err != nil {
		t.Fatal(err)
	}
	if p = pendingFolderOf(t, svc); p.AlbumID == nil || *p.AlbumID != album {
		t.Fatalf("pending folder album %+v", p)
	}

	folders, err := svc.AlbumPictures(ctx, album)
	if err != nil || len(folders) != 1 || !folders[0].Pending || len(folders[0].Pictures) != 2 || folders[0].FolderID != p.FolderID {
		t.Fatalf("album pictures %+v (%v)", folders, err)
	}
	if counts, err := svc.AlbumCounts(ctx); err != nil || len(counts) != 1 || counts[0] != (AlbumBoothCount{AlbumID: album, Pending: 2}) {
		t.Fatalf("counts %+v (%v)", counts, err)
	}

	// Approved, then more pictures later in the evening: both folders show, the pending one first.
	if _, err := svc.Approve(ctx, p.FolderID, "a@b"); err != nil {
		t.Fatal(err)
	}
	later := booth{t, svc, d, "session-0002"}
	later.put(KindPhoto, 1, "late")
	folders, err = svc.AlbumPictures(ctx, album)
	if err != nil || len(folders) != 2 || !folders[0].Pending || len(folders[0].Pictures) != 1 || folders[1].Pending ||
		len(folders[1].Pictures) != 2 {
		t.Fatalf("album pictures after approval %+v (%v)", folders, err)
	}
	if counts, _ := svc.AlbumCounts(ctx); len(counts) != 1 || counts[0] != (AlbumBoothCount{AlbumID: album, Pending: 1, Approved: 2}) {
		t.Fatalf("counts after approval %+v", counts)
	}
	if other, err := svc.AlbumPictures(ctx, album+1); err != nil || len(other) != 0 {
		t.Fatalf("another album %+v (%v)", other, err)
	}
}
