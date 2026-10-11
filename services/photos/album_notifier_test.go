package photos

import (
	"context"
	"database/sql"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/plugimt/transat-backend/internal/database/dbtest"
)

type recorder struct {
	mu    sync.Mutex
	calls []AlbumNotice
	err   error
	panic bool
}

func (r *recorder) notify(_ context.Context, n AlbumNotice) (Delivery, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.calls = append(r.calls, n)
	if r.panic {
		panic("boom")
	}
	if r.err != nil {
		return Delivery{}, r.err
	}
	return Delivery{Users: 2, Devices: 3}, nil
}

func (r *recorder) count() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return len(r.calls)
}

type clock struct{ t time.Time }

func (c *clock) now() time.Time          { return c.t }
func (c *clock) advance(d time.Duration) { c.t = c.t.Add(d) }

type env struct {
	t     *testing.T
	db    *sql.DB
	rec   *recorder
	clock *clock
	n     *AlbumNotifier
}

func setup(t *testing.T) *env {
	db := dbtest.Open(t)
	e := &env{t: t, db: db, rec: &recorder{}, clock: &clock{t: time.Date(2026, 3, 13, 20, 0, 0, 0, time.UTC)}}
	e.n = e.notifier()
	return e
}

func (e *env) notifier() *AlbumNotifier {
	n := NewAlbumNotifier(e.db, e.rec.notify, 30*time.Minute)
	n.now = e.clock.now
	return n
}

func (e *env) exec(q string, args ...any) {
	e.t.Helper()
	if _, err := e.db.Exec(q, args...); err != nil {
		e.t.Fatalf("%s: %v", q, err)
	}
}

func (e *env) fullSyncDone() {
	e.exec(`UPDATE photo_sync_state SET last_full_at = now(), full_run_id = NULL WHERE id = 1`)
}

func (e *env) album(root, title, visibility string, count int, parent any) {
	e.exec(`INSERT INTO photo_albums (root_folder_id, parent_album_id, title, raw_folder_name, event_date, drive_url, visibility, confidence, rule, last_seen_run, photo_count)
		VALUES ($1, $2, $3, $3, '2026-03-12', 'u', $4, 1, 'R', 1, $5)`, root, parent, title, visibility, count)
}

func (e *env) run() NotifyStats {
	e.t.Helper()
	stats, err := e.n.Run(context.Background())
	if err != nil {
		e.t.Fatalf("run: %v", err)
	}
	return stats
}

func (e *env) status(root string) string {
	var s sql.NullString
	_ = e.db.QueryRow(`SELECT status FROM photo_album_notifications WHERE root_folder_id = $1`, root).Scan(&s)
	return s.String
}

// settled runs once to register the album, waits out the settle delay and runs again.
func (e *env) settled() NotifyStats {
	e.run()
	e.clock.advance(31 * time.Minute)
	return e.run()
}

func TestNothingIsAnnouncedBeforeTheFirstCompleteSync(t *testing.T) {
	e := setup(t)
	e.album("OLD1", "Soirée A", "published", 10, nil)

	e.run()
	e.clock.advance(2 * time.Hour)
	e.run()
	if e.rec.count() != 0 {
		t.Fatalf("%d notifications before the first sync completed", e.rec.count())
	}
	var rows int
	_ = e.db.QueryRow(`SELECT count(*) FROM photo_album_notifications`).Scan(&rows)
	if rows != 0 {
		t.Fatalf("%d ledger rows written before the baseline", rows)
	}
}

func TestExistingAlbumsAreRecordedNotAnnounced(t *testing.T) {
	e := setup(t)
	e.album("OLD1", "Soirée A", "published", 10, nil)
	e.album("OLD2", "Tournoi B", "published", 40, nil)
	e.fullSyncDone()

	stats := e.run()
	if stats.Baselined != 2 || e.status("OLD1") != "baseline" || e.status("OLD2") != "baseline" {
		t.Fatalf("baseline = %+v", stats)
	}
	for i := 0; i < 5; i++ {
		e.clock.advance(time.Hour)
		e.run()
	}
	if e.rec.count() != 0 {
		t.Fatalf("%d notifications for albums that already existed", e.rec.count())
	}
}

func TestNewAlbumIsAnnouncedOnceItIsSettled(t *testing.T) {
	e := setup(t)
	e.fullSyncDone()
	e.run() // baseline of an empty library

	e.album("NEW1", "Gala", "published", 30, nil)
	e.run()
	if e.rec.count() != 0 {
		t.Fatal("announced before the settle delay")
	}
	e.clock.advance(31 * time.Minute)
	if stats := e.run(); stats.Sent != 1 {
		t.Fatalf("stats = %+v", stats)
	}
	got := e.rec.calls[0]
	if got.Title != "Gala" || got.AlbumID == 0 || got.EventDate == nil || got.EventDate.Format("2006-01-02") != "2026-03-12" {
		t.Fatalf("notice = %+v", got)
	}
	for i := 0; i < 4; i++ {
		e.clock.advance(time.Hour)
		e.run()
	}
	if e.rec.count() != 1 || e.status("NEW1") != "sent" {
		t.Fatalf("calls = %d, status = %s", e.rec.count(), e.status("NEW1"))
	}
}

func TestAlbumStillBeingUploadedIsNotAnnouncedHalfEmpty(t *testing.T) {
	e := setup(t)
	e.fullSyncDone()
	e.run()

	e.album("UP1", "Festival", "published", 5, nil)
	e.run()
	for _, count := range []int{40, 120, 300} {
		e.clock.advance(20 * time.Minute)
		e.exec(`UPDATE photo_albums SET photo_count = $1 WHERE root_folder_id = 'UP1'`, count)
		e.run()
		if e.rec.count() != 0 {
			t.Fatalf("announced while photos were still arriving (%d)", count)
		}
	}
	e.clock.advance(31 * time.Minute)
	e.run()
	if e.rec.count() != 1 {
		t.Fatalf("calls = %d once the count was stable", e.rec.count())
	}
}

func TestAlbumDiscoveredBeforeItsPhotosWaitsForThem(t *testing.T) {
	e := setup(t)
	e.fullSyncDone()
	e.run()

	e.album("EMPTY1", "Concert", "published", 0, nil)
	e.clock.advance(3 * time.Hour)
	e.run()
	if e.rec.count() != 0 || e.status("EMPTY1") != "" {
		t.Fatal("an album without photos is not available yet")
	}
	e.exec(`UPDATE photo_albums SET photo_count = 12 WHERE root_folder_id = 'EMPTY1'`)
	e.settled()
	if e.rec.count() != 1 {
		t.Fatalf("calls = %d", e.rec.count())
	}
}

func TestNestedHiddenAndUnpublishedAlbumsAreNotAnnounced(t *testing.T) {
	e := setup(t)
	e.fullSyncDone()
	e.run()

	e.album("PARENT", "Tournoi", "published", 50, nil)
	var parentID int
	_ = e.db.QueryRow(`SELECT id_photo_albums FROM photo_albums WHERE root_folder_id = 'PARENT'`).Scan(&parentID)
	e.album("CHILD", "Finale", "published", 20, parentID)
	e.album("REVIEW", "Dossier à vérifier", "needs_review", 20, nil)
	e.album("HIDDEN", "Caché", "hidden", 20, nil)
	e.exec(`INSERT INTO photo_albums (root_folder_id, title, raw_folder_name, drive_url, visibility, confidence, rule, last_seen_run, photo_count, deleted_at)
		VALUES ('GONE', 'Supprimé', 'Supprimé', 'u', 'published', 1, 'R', 1, 9, now())`)

	e.settled()
	if e.rec.count() != 1 || e.rec.calls[0].Title != "Tournoi" {
		t.Fatalf("announced %+v, want only the top-level published album", e.rec.calls)
	}

	// A reviewed album that becomes published is new to the users, and is announced once.
	e.exec(`UPDATE photo_albums SET visibility = 'published' WHERE root_folder_id = 'REVIEW'`)
	e.settled()
	if e.rec.count() != 2 || e.rec.calls[1].Title != "Dossier à vérifier" {
		t.Fatalf("calls = %+v", e.rec.calls)
	}
}

func TestRenamingOrUpdatingAnAlbumDoesNotAnnounceItAgain(t *testing.T) {
	e := setup(t)
	e.album("OLD1", "Ancien", "published", 10, nil)
	e.fullSyncDone()
	e.run()
	e.album("NEW1", "Nouveau", "published", 10, nil)
	e.settled()
	if e.rec.count() != 1 {
		t.Fatalf("calls = %d", e.rec.count())
	}

	e.exec(`UPDATE photo_albums SET title = 'Renommé', event_date = '2026-04-01', category = 'Sport' WHERE root_folder_id IN ('OLD1', 'NEW1')`)
	for i := 0; i < 3; i++ {
		e.clock.advance(time.Hour)
		e.run()
	}
	if e.rec.count() != 1 {
		t.Fatalf("a rename produced %d notifications", e.rec.count())
	}
}

func TestReindexedAlbumKeepsItsHistory(t *testing.T) {
	e := setup(t)
	e.fullSyncDone()
	e.run()
	e.album("NEW1", "Gala", "published", 10, nil)
	e.settled()

	// The index is rebuilt: the album row is recreated with another id but the same Drive folder.
	e.exec(`DELETE FROM photo_albums WHERE root_folder_id = 'NEW1'`)
	e.album("NEW1", "Gala", "published", 10, nil)
	e.clock.advance(time.Hour)
	e.run()
	e.clock.advance(time.Hour)
	e.run()
	if e.rec.count() != 1 {
		t.Fatalf("a rebuilt index produced %d notifications", e.rec.count())
	}
}

func TestConcurrentRunsAnnounceOnce(t *testing.T) {
	e := setup(t)
	e.fullSyncDone()
	e.run()
	e.album("NEW1", "Gala", "published", 10, nil)
	e.run()
	e.clock.advance(31 * time.Minute)

	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			n := e.notifier() // as if several instances ran the job at the same moment
			if _, err := n.Run(context.Background()); err != nil {
				t.Errorf("run: %v", err)
			}
		}()
	}
	wg.Wait()
	if e.rec.count() != 1 {
		t.Fatalf("%d concurrent runs announced the album %d times", 8, e.rec.count())
	}
}

func TestConcurrentFirstRunsBaselineOnce(t *testing.T) {
	e := setup(t)
	e.album("OLD1", "Ancien", "published", 10, nil)
	e.fullSyncDone()

	var wg sync.WaitGroup
	for i := 0; i < 6; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if _, err := e.notifier().Run(context.Background()); err != nil {
				t.Errorf("run: %v", err)
			}
		}()
	}
	wg.Wait()
	e.clock.advance(2 * time.Hour)
	e.run()
	if e.rec.count() != 0 || e.status("OLD1") != "baseline" {
		t.Fatalf("calls = %d, status = %s", e.rec.count(), e.status("OLD1"))
	}
}

func TestFailedDeliveryIsRetriedThenGivenUp(t *testing.T) {
	e := setup(t)
	e.fullSyncDone()
	e.run()
	e.album("NEW1", "Gala", "published", 10, nil)
	e.rec.err = errors.New("expo unavailable")

	e.settled() // attempt 1
	if e.rec.count() != 1 || e.status("NEW1") != "pending" {
		t.Fatalf("after attempt 1: calls = %d, status = %s", e.rec.count(), e.status("NEW1"))
	}
	e.clock.advance(time.Minute)
	e.run()
	if e.rec.count() != 1 {
		t.Fatal("retried before its delay")
	}
	e.clock.advance(15 * time.Minute)
	e.run() // attempt 2
	e.clock.advance(25 * time.Minute)
	e.run() // attempt 3
	if e.rec.count() != 3 || e.status("NEW1") != "failed" {
		t.Fatalf("after attempt 3: calls = %d, status = %s", e.rec.count(), e.status("NEW1"))
	}
	e.clock.advance(24 * time.Hour)
	e.run()
	if e.rec.count() != 3 {
		t.Fatal("kept retrying after giving up")
	}
}

func TestDeliveryThatRecoversIsSentOnce(t *testing.T) {
	e := setup(t)
	e.fullSyncDone()
	e.run()
	e.album("NEW1", "Gala", "published", 10, nil)
	e.rec.err = errors.New("expo unavailable")
	e.settled()
	e.rec.err = nil
	e.clock.advance(15 * time.Minute)
	e.run()
	e.clock.advance(time.Hour)
	e.run()
	if e.rec.count() != 2 || e.status("NEW1") != "sent" {
		t.Fatalf("calls = %d, status = %s", e.rec.count(), e.status("NEW1"))
	}
}

func TestAPanicInDeliveryNeverBreaksTheRun(t *testing.T) {
	e := setup(t)
	e.fullSyncDone()
	e.run()
	e.album("NEW1", "Gala", "published", 10, nil)
	e.rec.panic = true
	stats := e.settled()
	if stats.Retrying != 1 || e.status("NEW1") != "pending" {
		t.Fatalf("stats = %+v, status = %s", stats, e.status("NEW1"))
	}
}

func TestAClaimLeftBehindByACrashIsRecovered(t *testing.T) {
	e := setup(t)
	e.fullSyncDone()
	e.run()
	e.album("NEW1", "Gala", "published", 10, nil)
	e.run()
	e.clock.advance(31 * time.Minute)
	// Another process claimed it and died before recording the result.
	e.exec(`UPDATE photo_album_notifications SET status = 'sending', attempts = 1, claimed_at = $1 WHERE root_folder_id = 'NEW1'`, e.clock.t.Add(-time.Minute))
	e.run()
	if e.rec.count() != 0 {
		t.Fatal("took over a claim that is still held")
	}
	e.clock.advance(20 * time.Minute)
	e.run()
	if e.rec.count() != 1 || e.status("NEW1") != "sent" {
		t.Fatalf("calls = %d, status = %s", e.rec.count(), e.status("NEW1"))
	}
}
