package photos

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/plugimt/transat-backend/internal/database/dbtest"
)

// fakeDriveAPI is an in-memory shared drive whose change feed is fed by the test.
type fakeDriveAPI struct {
	mu      sync.Mutex
	items   map[string]DriveItem
	pending []Change
}

func (f *fakeDriveAPI) put(it DriveItem, announce bool) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.items == nil {
		f.items = map[string]DriveItem{}
	}
	f.items[it.ID] = it
	if announce {
		c := it
		f.pending = append(f.pending, Change{FileID: it.ID, Item: &c})
	}
}

func (f *fakeDriveAPI) StartPageToken(context.Context) (string, error) { return "t0", nil }

func (f *fakeDriveAPI) ListFolders(context.Context) ([]DriveItem, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	var out []DriveItem
	for _, it := range f.items {
		if it.IsFolder() && !it.Trashed {
			out = append(out, it)
		}
	}
	return out, nil
}

func (f *fakeDriveAPI) ListFilesPage(context.Context, string) ([]DriveItem, string, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	var out []DriveItem
	for _, it := range f.items {
		if !it.IsFolder() && !it.Trashed {
			out = append(out, it)
		}
	}
	return out, "", nil
}

func (f *fakeDriveAPI) Changes(context.Context, string) (ChangePage, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	page := ChangePage{Changes: f.pending, NewStartPageToken: "t1"}
	f.pending = nil
	return page, nil
}

func (f *fakeDriveAPI) GetItem(_ context.Context, id string) (DriveItem, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.items[id], nil
}

func (f *fakeDriveAPI) ListChildren(_ context.Context, folderID string) ([]DriveItem, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	var out []DriveItem
	for _, it := range f.items {
		if it.ParentID == folderID && !it.Trashed {
			out = append(out, it)
		}
	}
	return out, nil
}

func (f *fakeDriveAPI) Calls() map[string]int64 { return map[string]int64{} }

const testDrive = "drive0000001"

func folder(id, parent, name string) DriveItem {
	return DriveItem{ID: id, ParentID: parent, Name: name, MimeType: MimeFolder, ModifiedTime: time.Now()}
}

func photo(id, parent, name string) DriveItem {
	t := time.Date(2026, 10, 27, 22, 0, 0, 0, time.UTC)
	return DriveItem{ID: id, ParentID: parent, Name: name, MimeType: "image/jpeg", Size: 10, CreatedTime: t, ModifiedTime: t}
}

func visible(t *testing.T, s *SyncService, id string) bool {
	t.Helper()
	var n int
	if err := s.db.QueryRow(`SELECT COUNT(*) FROM photos WHERE drive_file_id = $1 AND deleted_at IS NULL`, id).Scan(&n); err != nil {
		t.Fatal(err)
	}
	return n == 1
}

func TestPendingPhotoboothFolderIsIndexedOnlyOnceRenamed(t *testing.T) {
	db := dbtest.Open(t)
	api := &fakeDriveAPI{}
	api.put(folder("year00000001", testDrive, "2026-2027"), false)
	api.put(folder("event0000001", "year00000001", "2026.10.27 - Soirée Intégration"), false)
	api.put(folder("pending00001", "event0000001", "photobooth_temp"), false)
	api.put(folder("session00001", "pending00001", "Session 1"), false)
	api.put(photo("eventphoto01", "event0000001", "dj.jpg"), false)
	api.put(photo("boothphoto01", "pending00001", "booth-1.jpg"), false)
	api.put(photo("boothphoto02", "session00001", "booth-2.jpg"), false)

	s := NewSyncService(db, api, testDrive, time.Hour)
	ctx := context.Background()
	if _, err := s.Run(ctx, ModeFull); err != nil {
		t.Fatal(err)
	}
	if !visible(t, s, "eventphoto01") {
		t.Fatal("the photographers' photo is not indexed")
	}
	if visible(t, s, "boothphoto01") || visible(t, s, "boothphoto02") {
		t.Fatal("photobooth pictures awaiting approval were indexed")
	}

	// A new picture arrives while the evening goes on: still hidden.
	api.put(photo("boothphoto03", "pending00001", "booth-3.jpg"), true)
	if _, err := s.Run(ctx, ModeIncremental); err != nil {
		t.Fatal(err)
	}
	if visible(t, s, "boothphoto03") {
		t.Fatal("a picture uploaded into photobooth_temp was indexed")
	}

	// Approved (by the panel or by hand on Drive): the feed only reports the folder, its pictures follow.
	api.put(folder("pending00001", "event0000001", "Photobooth"), true)
	if _, err := s.Run(ctx, ModeIncremental); err != nil {
		t.Fatal(err)
	}
	for _, id := range []string{"boothphoto01", "boothphoto03"} {
		if !visible(t, s, id) {
			t.Fatalf("%s not indexed after the folder was approved", id)
		}
	}
	var sameAlbum bool
	if err := db.QueryRow(`SELECT (SELECT album_id FROM photos WHERE drive_file_id = 'boothphoto01') =
		(SELECT album_id FROM photos WHERE drive_file_id = 'eventphoto01')`).Scan(&sameAlbum); err != nil {
		t.Fatal(err)
	}
	if !sameAlbum {
		t.Fatal("approved photobooth pictures are not in the event's album")
	}

	// Renamed back to photobooth_temp: hidden again.
	api.put(folder("pending00001", "event0000001", "Photobooth_Temp"), true)
	if _, err := s.Run(ctx, ModeIncremental); err != nil {
		t.Fatal(err)
	}
	if visible(t, s, "boothphoto01") {
		t.Fatal("pictures stayed visible after their folder went back to photobooth_temp")
	}
}

func TestIsPendingPhotoboothFolder(t *testing.T) {
	for name, want := range map[string]bool{
		"photobooth_temp": true, "Photobooth_Temp": true, " photobooth temp ": true, "photobooth-temp": true,
		"Photobooth": false, "photobooth_temporaire": false, "temp": false,
	} {
		if got := IsPendingPhotoboothFolder(name); got != want {
			t.Errorf("IsPendingPhotoboothFolder(%q) = %v, want %v", name, got, want)
		}
	}
}
