package photobooth

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/plugimt/transat-backend/internal/database/dbtest"
)

// fakeDrive is an in-memory Drive: folders, files, and upload sessions delivered by the test.
type fakeDrive struct {
	mu       sync.Mutex
	seq      int
	folders  map[string]Folder
	parentOf map[string]string
	files    map[string]*FileInfo
	sessions map[string]FileInfo
	created  []string
}

func newFakeDrive(root string) *fakeDrive {
	return &fakeDrive{
		folders: map[string]Folder{root: {ID: root, Name: "root"}}, parentOf: map[string]string{},
		files: map[string]*FileInfo{}, sessions: map[string]FileInfo{},
	}
}

func (d *fakeDrive) id(prefix string) string {
	d.seq++
	return fmt.Sprintf("%s%012d", prefix, d.seq)
}

func (d *fakeDrive) addFolder(parent, name string) string {
	d.mu.Lock()
	defer d.mu.Unlock()
	id := d.id("folder")
	d.folders[id] = Folder{ID: id, Name: name}
	d.parentOf[id] = parent
	return id
}

func (d *fakeDrive) ListChildFolders(_ context.Context, parent string) ([]Folder, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	var out []Folder
	for id, p := range d.parentOf {
		if p == parent {
			out = append(out, d.folders[id])
		}
	}
	return out, nil
}

func (d *fakeDrive) ListChildFilesNamed(_ context.Context, parent, name string) ([]FileInfo, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	var out []FileInfo
	for _, f := range d.files {
		if f.Parents[0] == parent && f.Name == name && !f.Trashed {
			out = append(out, *f)
		}
	}
	return out, nil
}

func (d *fakeDrive) CreateFolder(_ context.Context, parent, name string) (string, error) {
	id := d.addFolder(parent, name)
	d.mu.Lock()
	d.created = append(d.created, name)
	d.mu.Unlock()
	return id, nil
}

func (d *fakeDrive) StartUpload(_ context.Context, parent, name, mime string, size int64) (string, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if _, ok := d.folders[parent]; !ok {
		return "", ErrNotFound
	}
	u := "https://upload.example/" + d.id("session")
	d.sessions[u] = FileInfo{Name: name, Parents: []string{parent}, Size: size, MimeType: mime}
	return u, nil
}

// deliver plays the booth sending bytes to the upload URL; returns the new Drive file id.
func (d *fakeDrive) deliver(t *testing.T, uploadURL string, content []byte) string {
	t.Helper()
	d.mu.Lock()
	defer d.mu.Unlock()
	s, ok := d.sessions[uploadURL]
	if !ok {
		t.Fatalf("unknown upload session %s", uploadURL)
	}
	delete(d.sessions, uploadURL)
	sum := md5.Sum(content)
	f := s
	f.ID = d.id("file")
	f.Size = int64(len(content))
	f.MD5 = hex.EncodeToString(sum[:])
	d.files[f.ID] = &f
	return f.ID
}

func (d *fakeDrive) GetFile(_ context.Context, id string) (FileInfo, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if f, ok := d.files[id]; ok {
		return *f, nil
	}
	if f, ok := d.folders[id]; ok {
		return FileInfo{ID: f.ID, Name: f.Name}, nil
	}
	return FileInfo{}, ErrNotFound
}

func (d *fakeDrive) Trash(_ context.Context, id string) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	if f, ok := d.files[id]; ok {
		f.Trashed = true
	}
	return nil
}

func (d *fakeDrive) Download(_ context.Context, id string) (io.ReadCloser, error) {
	if _, err := d.GetFile(context.Background(), id); err != nil {
		return nil, err
	}
	return io.NopCloser(strings.NewReader("original-" + id)), nil
}

func (d *fakeDrive) Thumbnail(context.Context, FileInfo, int) (io.ReadCloser, string, error) {
	return nil, "", ErrNotFound
}

func (d *fakeDrive) live() map[string]bool {
	d.mu.Lock()
	defer d.mu.Unlock()
	out := map[string]bool{}
	for _, f := range d.files {
		if !f.Trashed {
			out[f.Name] = true
		}
	}
	return out
}

// waitLive waits for background trashing to settle on the expected set of live file names.
func waitLive(t *testing.T, d *fakeDrive, want ...string) {
	t.Helper()
	deadline := time.Now().Add(3 * time.Second)
	for {
		got := d.live()
		ok := len(got) == len(want)
		for _, w := range want {
			ok = ok && got[w]
		}
		if ok {
			return
		}
		if time.Now().After(deadline) {
			t.Fatalf("live files = %v, want %v", got, want)
		}
		time.Sleep(20 * time.Millisecond)
	}
}

var event = Event{Name: "Soirée Intégration", Date: "2026-10-27", LayoutType: "GRID_2X2"}

type booth struct {
	t     *testing.T
	svc   *Service
	drive *fakeDrive
	sess  string
}

func (b booth) upload(kind string, n int, content string) (StartResult, string) {
	b.t.Helper()
	sum := md5.Sum([]byte(content))
	res, err := b.svc.StartUpload(context.Background(), "booth-1", b.sess, StartRequest{
		Event: event, Kind: kind, Number: n, Size: int64(len(content)), MD5: hex.EncodeToString(sum[:]),
	})
	if err != nil {
		b.t.Fatalf("start %s %d: %v", kind, n, err)
	}
	return res, b.drive.deliver(b.t, res.UploadURL, []byte(content))
}

func (b booth) put(kind string, n int, content string) StartResult {
	b.t.Helper()
	res, driveID := b.upload(kind, n, content)
	if err := b.svc.CompleteUpload(context.Background(), "booth-1", b.sess, res.FileID, driveID); err != nil {
		b.t.Fatalf("complete %s %d: %v", kind, n, err)
	}
	return res
}

func setup(t *testing.T) (*Service, *fakeDrive) {
	db := dbtest.Open(t)
	d := newFakeDrive("rootfolder0001")
	return NewService(db, d, "rootfolder0001", "https://transat.dev/photobooth/"), d
}

func status(err error) int {
	var e *Error
	if errors.As(err, &e) {
		return e.Status
	}
	return 0
}

func TestFolderIsCreatedOnceThenReused(t *testing.T) {
	svc, d := setup(t)
	b := booth{t, svc, d, "session-0001"}
	b.put(KindPhoto, 1, "p1")
	b.put(KindPhoto, 2, "p2")
	other := booth{t, svc, d, "session-0002"}
	other.put(KindPhoto, 1, "q1")
	want := []string{"2026-2027", "2026.10.27 - Soirée Intégration", "Photobooth"}
	if fmt.Sprint(d.created) != fmt.Sprint(want) {
		t.Fatalf("created folders %v, want %v", d.created, want)
	}
}

func TestExistingFoldersAreUsed(t *testing.T) {
	svc, d := setup(t)
	year := d.addFolder("rootfolder0001", "2026 - 2027")
	d.addFolder(year, "2026.10.26 - Soirée Intégration")
	ev := d.addFolder(year, "2026.10.27 - soiree integration")
	pbID := d.addFolder(ev, "PhotoBooth")
	b := booth{t, svc, d, "session-0001"}
	res := b.put(KindPhoto, 1, "p1")
	if len(d.created) != 0 {
		t.Fatalf("created %v although every folder existed", d.created)
	}
	f, _ := d.ListChildFilesNamed(context.Background(), pbID, res.FileName)
	if len(f) != 1 {
		t.Fatal("photo not in the existing Photobooth folder")
	}
}

func TestFullSessionGivesAPermanentLink(t *testing.T) {
	svc, d := setup(t)
	b := booth{t, svc, d, "session-0001"}
	ctx := context.Background()
	b.put(KindPhoto, 1, "p1")
	b.put(KindPhoto, 2, "p2")

	if _, err := svc.Finalize(ctx, "booth-1", b.sess, 2); status(err) != 409 {
		t.Fatalf("finalize without collage: %v", err)
	}
	b.put(KindCollage, 0, "collage")
	if _, err := svc.Finalize(ctx, "booth-1", b.sess, 3); status(err) != 409 {
		t.Fatalf("finalize with a missing photo: %v", err)
	}
	link, err := svc.Finalize(ctx, "booth-1", b.sess, 2)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(link, "https://transat.dev/photobooth/") {
		t.Fatalf("link %s", link)
	}
	again, err := svc.Finalize(ctx, "booth-1", b.sess, 2)
	if err != nil || again != link {
		t.Fatalf("second finalize gave %s (%v), want the same link", again, err)
	}
	token := strings.TrimPrefix(link, "https://transat.dev/photobooth/")

	share, err := svc.GetShare(ctx, token)
	if err != nil {
		t.Fatal(err)
	}
	if share.EventName != event.Name || share.EventDate != event.Date || len(share.Files) != 3 || share.Files[0].Kind != KindCollage ||
		share.Files[1].Number != 1 || share.Files[2].Number != 2 {
		t.Fatalf("share %+v", share)
	}
	f, err := svc.OpenSharedFile(ctx, token, share.Files[1].ID, 0)
	if err != nil {
		t.Fatal(err)
	}
	body, _ := io.ReadAll(f.Body)
	if !bytes.HasPrefix(body, []byte("original-")) || f.Name != share.Files[1].Name {
		t.Fatalf("file %q %s", body, f.Name)
	}

	// A share link reveals nothing of other shoots, and wrong tokens see nothing.
	other := booth{t, svc, d, "session-0002"}
	o := other.put(KindPhoto, 1, "x")
	if _, err := svc.OpenSharedFile(ctx, token, o.FileID, 0); status(err) != 404 {
		t.Fatalf("file of another shoot reachable: %v", err)
	}
	if _, err := svc.GetShare(ctx, strings.Repeat("A", 32)); status(err) != 404 {
		t.Fatalf("unknown token: %v", err)
	}

	// Once shared, the shoot cannot be cancelled or changed.
	if err := svc.Cancel(ctx, "booth-1", b.sess); status(err) != 409 {
		t.Fatalf("cancel after share: %v", err)
	}
	if err := svc.DiscardSlot(ctx, "booth-1", b.sess, KindPhoto, 1); status(err) != 409 {
		t.Fatalf("discard after share: %v", err)
	}
}

func TestRetakeAndCancelRemoveEverythingFromDrive(t *testing.T) {
	svc, d := setup(t)
	ctx := context.Background()
	b := booth{t, svc, d, "session-0001"}
	first := b.put(KindPhoto, 1, "p1")
	// Retake: the booth discards the slot, then uploads again.
	if err := svc.DiscardSlot(ctx, "booth-1", b.sess, KindPhoto, 1); err != nil {
		t.Fatal(err)
	}
	second := b.put(KindPhoto, 1, "p1-again")
	if first.FileName == second.FileName {
		t.Fatal("a retake reuses the previous file name")
	}
	waitLive(t, d, second.FileName)

	// A retake that skips the discard call replaces the slot as well.
	third := b.put(KindPhoto, 1, "p1-third")
	waitLive(t, d, third.FileName)

	// An upload still in flight when the guest cancels is trashed when it completes.
	res, driveID := b.upload(KindPhoto, 2, "p2")
	if err := svc.Cancel(ctx, "booth-1", b.sess); err != nil {
		t.Fatal(err)
	}
	if err := svc.CompleteUpload(ctx, "booth-1", b.sess, res.FileID, driveID); status(err) != 410 {
		t.Fatalf("late complete after cancel: %v", err)
	}
	waitLive(t, d)
	if err := svc.Cancel(ctx, "booth-1", b.sess); err != nil {
		t.Fatalf("second cancel: %v", err)
	}
	if _, err := svc.StartUpload(ctx, "booth-1", b.sess, StartRequest{Event: event, Kind: KindPhoto, Number: 3, Size: 1,
		MD5: strings.Repeat("0", 32)}); status(err) != 409 {
		t.Fatalf("upload after cancel: %v", err)
	}
	if _, err := svc.Finalize(ctx, "booth-1", b.sess, 1); status(err) != 409 {
		t.Fatalf("finalize after cancel: %v", err)
	}
}

func TestOrphanUploadIsFoundAndTrashed(t *testing.T) {
	svc, d := setup(t)
	ctx := context.Background()
	b := booth{t, svc, d, "session-0001"}
	_, _ = b.upload(KindPhoto, 1, "p1") // the booth died before confirming
	if err := svc.Cancel(ctx, "booth-1", b.sess); err != nil {
		t.Fatal(err)
	}
	if n, _ := svc.RunTrash(ctx); n != 0 {
		t.Fatal("orphan handled before its grace period")
	}
	svc.now = func() time.Time { return time.Now().Add(orphanGrace + time.Minute) }
	if _, err := svc.RunTrash(ctx); err != nil {
		t.Fatal(err)
	}
	waitLive(t, d)
}

func TestCompleteChecksTheUploadedContent(t *testing.T) {
	svc, d := setup(t)
	ctx := context.Background()
	b := booth{t, svc, d, "session-0001"}
	sum := md5.Sum([]byte("expected"))
	res, err := svc.StartUpload(ctx, "booth-1", b.sess, StartRequest{Event: event, Kind: KindPhoto, Number: 1, Size: 8,
		MD5: hex.EncodeToString(sum[:])})
	if err != nil {
		t.Fatal(err)
	}
	driveID := d.deliver(t, res.UploadURL, []byte("tampered"))
	if err := svc.CompleteUpload(ctx, "booth-1", b.sess, res.FileID, driveID); status(err) != 400 {
		t.Fatalf("mismatching checksum accepted: %v", err)
	}
	// Another booth cannot complete or touch this session.
	if err := svc.CompleteUpload(ctx, "booth-2", b.sess, res.FileID, driveID); status(err) != 404 {
		t.Fatalf("other booth: %v", err)
	}
	// A file that is not the announced upload is refused.
	unrelated := d.addFolder("rootfolder0001", "elsewhere")
	if err := svc.CompleteUpload(ctx, "booth-1", b.sess, res.FileID, unrelated); status(err) != 400 {
		t.Fatalf("unrelated drive item accepted: %v", err)
	}
}

func TestSessionCannotSwitchEvent(t *testing.T) {
	svc, d := setup(t)
	b := booth{t, svc, d, "session-0001"}
	b.put(KindPhoto, 1, "p1")
	_, err := svc.StartUpload(context.Background(), "booth-1", b.sess, StartRequest{
		Event: Event{Name: "Gala", Date: event.Date}, Kind: KindPhoto, Number: 2, Size: 1, MD5: strings.Repeat("a", 32),
	})
	if status(err) != 409 {
		t.Fatalf("event switch: %v", err)
	}
}
