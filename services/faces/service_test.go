package faces

import (
	"context"
	"database/sql"
	"errors"
	"sort"
	"sync"
	"testing"

	"github.com/plugimt/transat-backend/handlers/photos/repository"
	"github.com/plugimt/transat-backend/internal/database/dbtest"
	"github.com/plugimt/transat-backend/services/accounts"
)

// fakeML answers by image content: the images are the photo ids (or "selfie-…"), each mapped to its faces.
type fakeML struct {
	mu    sync.Mutex
	faces map[string][]Face
	fail  map[string]error
	calls int
}

func (f *fakeML) Model() string { return "buffalo_l" }

func (f *fakeML) Detect(_ context.Context, img []byte) (Detection, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.calls++
	if err := f.fail[string(img)]; err != nil {
		return Detection{}, err
	}
	return Detection{ImageWidth: 2000, ImageHeight: 1500, Faces: f.faces[string(img)]}, nil
}

type fakeImages struct{ missing map[string]bool }

func (f fakeImages) Image(_ context.Context, id string, _ int) ([]byte, error) {
	if f.missing[id] {
		return nil, ErrImageUnavailable
	}
	return []byte(id), nil
}

// Three people as unit vectors; a slightly different shot of the same person stays well under 0.5.
var (
	alice    = Embedding{1, 0, 0, 0}
	aliceAlt = mustNorm(Embedding{0.9, 0.3, 0, 0})
	bob      = Embedding{0, 1, 0, 0}
	stranger = Embedding{0, 0, 1, 0}
)

func mustNorm(e Embedding) Embedding {
	n, err := normalize(e)
	if err != nil {
		panic(err)
	}
	return n
}

func on(e Embedding, x float32) Face {
	return Face{X1: x, Y1: 0.2, X2: x + 0.2, Y2: 0.5, Score: 0.95, Embedding: e}
}

type env struct {
	t   *testing.T
	db  *sql.DB
	ml  *fakeML
	svc *Service
}

func setup(t *testing.T) env {
	db := dbtest.Open(t)
	ml := &fakeML{faces: map[string][]Face{}, fail: map[string]error{}}
	cfg := DefaultConfig()
	cfg.TickBudget = 0
	cfg.ScanBatch = 2 // several batches per tick
	svc := NewService(db, ml, fakeImages{missing: map[string]bool{"GONE-PHOTO-1": true}}, cfg)
	return env{t, db, ml, svc}
}

func (e env) exec(q string, args ...any) {
	e.t.Helper()
	if _, err := e.db.Exec(q, args...); err != nil {
		e.t.Fatal(err)
	}
}

func (e env) user(email string) {
	e.exec(`INSERT INTO newf (email, password, first_name, last_name) VALUES ($1, 'x', 'Jo', 'Doe')`, email)
}

func (e env) photo(id string, minutes int, faces ...Face) {
	e.exec(`INSERT INTO photos (drive_file_id, name, mime_type, media_type, md5, sort_at, sort_at_source, web_url, last_seen_run)
		VALUES ($1, $1, 'image/jpeg', 'image', $1, now() - make_interval(mins => $2), 'x', 'u', 1)`, id, minutes)
	e.ml.faces[id] = faces
}

func (e env) tick() TickStats {
	e.t.Helper()
	st, err := e.svc.Tick(context.Background())
	if err != nil {
		e.t.Fatal(err)
	}
	return st
}

func (e env) enrol(email string, selfie Embedding) {
	e.t.Helper()
	ctx := context.Background()
	if err := e.svc.Consent(ctx, email, ConsentVersion); err != nil {
		e.t.Fatal(err)
	}
	e.ml.faces["selfie-"+email] = []Face{{X1: 0.3, Y1: 0.2, X2: 0.7, Y2: 0.8, Score: 0.99, Embedding: selfie}}
	if err := e.svc.AddReference(ctx, email, []byte("selfie-"+email)); err != nil {
		e.t.Fatal(err)
	}
}

func (e env) album(email string) []string {
	e.t.Helper()
	rows, err := repository.NewPhotosRepository(e.db).MyPhotos(context.Background(), email, nil, 100)
	if err != nil {
		e.t.Fatal(err)
	}
	ids := []string{}
	for _, r := range rows {
		ids = append(ids, r.ID)
	}
	sort.Strings(ids)
	return ids
}

func (e env) count(q string, args ...any) int {
	e.t.Helper()
	var n int
	if err := e.db.QueryRow(q, args...).Scan(&n); err != nil {
		e.t.Fatal(err)
	}
	return n
}

func equal(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestAlbumsFollowConsentPastAndFuturePhotos(t *testing.T) {
	e := setup(t)
	ctx := context.Background()
	for _, u := range []string{"alice@imt.fr", "bob@imt.fr", "carol@imt.fr"} {
		e.user(u)
	}
	e.photo("PAST-ALICE-1", 300, on(aliceAlt, 0.1), on(stranger, 0.5))
	e.photo("PAST-GROUP-1", 200, on(alice, 0.1), on(bob, 0.4))
	e.photo("PAST-BOB-01", 100, on(bob, 0.2))
	e.photo("PAST-EMPTY-1", 50)
	e.photo("GONE-PHOTO-1", 40, on(alice, 0.2))

	// Before anyone opts in, the library is scanned but nobody is matched.
	st := e.tick()
	if st.Scanned != 4 || st.Failed != 1 || st.Faces != 5 || st.Matches != 0 {
		t.Fatalf("first tick: %+v", st)
	}
	if n := e.count(`SELECT count(*) FROM face_matches`); n != 0 {
		t.Fatalf("%d matches without consent", n)
	}

	// Alice opts in: her past photos come up; Bob, who never consented, gets nothing.
	e.enrol("alice@imt.fr", alice)
	e.tick()
	if got := e.album("alice@imt.fr"); !equal(got, []string{"PAST-ALICE-1", "PAST-GROUP-1"}) {
		t.Fatalf("alice past album: %v", got)
	}
	if got := e.album("bob@imt.fr"); len(got) != 0 {
		t.Fatalf("bob has no consent but got %v", got)
	}
	status, err := e.svc.Status(ctx, "alice@imt.fr")
	if err != nil {
		t.Fatal(err)
	}
	if !status.Enabled || status.References != 1 || status.PhotoCount != 2 || status.Matching || status.LibraryScanned != 5 || status.LibraryTotal != 5 {
		t.Fatalf("status: %+v", status)
	}

	// A future photo is matched on the next tick, without a new scan of the past.
	calls := e.ml.calls
	e.photo("NEW-ALICE-01", 1, on(alice, 0.3))
	e.tick()
	if e.ml.calls != calls+1 {
		t.Fatalf("expected one model call for the new photo, got %d", e.ml.calls-calls)
	}
	if got := e.album("alice@imt.fr"); !equal(got, []string{"NEW-ALICE-01", "PAST-ALICE-1", "PAST-GROUP-1"}) {
		t.Fatalf("alice album after new photo: %v", got)
	}

	// "Not me" removes a photo and a later tick does not bring it back.
	if err := e.svc.Reject(ctx, "alice@imt.fr", "PAST-GROUP-1"); err != nil {
		t.Fatal(err)
	}
	if err := e.svc.Reject(ctx, "alice@imt.fr", "PAST-BOB-01"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("rejecting a photo not in the album: %v", err)
	}
	e.exec(`UPDATE face_references SET matched_through = 0`)
	e.tick()
	if got := e.album("alice@imt.fr"); !equal(got, []string{"NEW-ALICE-01", "PAST-ALICE-1"}) {
		t.Fatalf("rejected photo came back: %v", got)
	}

	// A photo leaving the gallery takes its faces (and so the matches) with it.
	e.exec(`UPDATE photos SET deleted_at = now() WHERE drive_file_id = 'NEW-ALICE-01'`)
	st = e.tick()
	if st.Purged != 1 {
		t.Fatalf("purge: %+v", st)
	}
	if n := e.count(`SELECT count(*) FROM photo_faces WHERE drive_file_id = 'NEW-ALICE-01'`); n != 0 {
		t.Fatalf("faces of a deleted photo kept: %d", n)
	}
	if got := e.album("alice@imt.fr"); !equal(got, []string{"PAST-ALICE-1"}) {
		t.Fatalf("album after deletion: %v", got)
	}

	// Withdrawing erases every face data of Alice, and nothing comes back.
	if err := e.svc.Withdraw(ctx, "alice@imt.fr"); err != nil {
		t.Fatal(err)
	}
	e.tick()
	for _, table := range []string{"face_consents", "face_references", "face_matches"} {
		if n := e.count(`SELECT count(*) FROM `+table+` WHERE email = $1`, "alice@imt.fr"); n != 0 {
			t.Fatalf("%s still has %d rows after withdrawal", table, n)
		}
	}
	if err := e.svc.AddReference(ctx, "alice@imt.fr", []byte("selfie-alice@imt.fr")); !errors.Is(err, ErrNoConsent) {
		t.Fatalf("reference without consent: %v", err)
	}
}

func TestConsentAndReferenceRules(t *testing.T) {
	e := setup(t)
	ctx := context.Background()
	e.user("dan@imt.fr")

	if err := e.svc.Consent(ctx, "dan@imt.fr", "2020-old"); !errors.Is(err, ErrConsentVersion) {
		t.Fatalf("old consent version: %v", err)
	}
	if err := e.svc.AddReference(ctx, "dan@imt.fr", []byte("x")); !errors.Is(err, ErrNoConsent) {
		t.Fatalf("reference before consent: %v", err)
	}
	if err := e.svc.Consent(ctx, "dan@imt.fr", ConsentVersion); err != nil {
		t.Fatal(err)
	}
	e.ml.faces["group-selfie"] = []Face{on(alice, 0.1), on(bob, 0.5)}
	if err := e.svc.AddReference(ctx, "dan@imt.fr", []byte("group-selfie")); !errors.Is(err, ErrManyFaces) {
		t.Fatalf("group selfie: %v", err)
	}
	if err := e.svc.AddReference(ctx, "dan@imt.fr", []byte("landscape")); !errors.Is(err, ErrNoFace) {
		t.Fatalf("no face: %v", err)
	}
	e.ml.fail["corrupt"] = ErrUnreadableImage
	if err := e.svc.AddReference(ctx, "dan@imt.fr", []byte("corrupt")); !errors.Is(err, ErrNoFace) {
		t.Fatalf("unreadable: %v", err)
	}
	e.ml.faces["selfie"] = []Face{{X1: 0.3, Y1: 0.2, X2: 0.7, Y2: 0.8, Score: 0.99, Embedding: alice}}
	for i := 0; i < MaxReferences; i++ {
		if err := e.svc.AddReference(ctx, "dan@imt.fr", []byte("selfie")); err != nil {
			t.Fatal(err)
		}
	}
	if err := e.svc.AddReference(ctx, "dan@imt.fr", []byte("selfie")); !errors.Is(err, ErrTooManyRefs) {
		t.Fatalf("over the limit: %v", err)
	}

	// A new consent text: existing references stay, new ones need the new consent.
	e.exec(`UPDATE face_consents SET consent_version = 'previous'`)
	e.exec(`DELETE FROM face_references WHERE id = (SELECT min(id) FROM face_references)`)
	if err := e.svc.AddReference(ctx, "dan@imt.fr", []byte("selfie")); !errors.Is(err, ErrConsentVersion) {
		t.Fatalf("outdated consent: %v", err)
	}
	st, err := e.svc.Status(ctx, "dan@imt.fr")
	if err != nil {
		t.Fatal(err)
	}
	if st.ConsentVersion == nil || *st.ConsentVersion != "previous" || st.CurrentConsentVersion != ConsentVersion {
		t.Fatalf("status must expose the outdated version: %+v", st)
	}

	// Without the ML service the feature is off.
	off := NewService(e.db, nil, nil, DefaultConfig())
	if err := off.Consent(ctx, "dan@imt.fr", ConsentVersion); !errors.Is(err, ErrUnavailable) {
		t.Fatalf("unavailable: %v", err)
	}
	if st, err := off.Status(ctx, "dan@imt.fr"); err != nil || st.Available {
		t.Fatalf("unavailable status: %+v %v", st, err)
	}
}

func TestFailedScansAreRetriedThenGivenUp(t *testing.T) {
	e := setup(t)
	e.photo("FLAKY-PHOTO1", 10, on(stranger, 0.1))
	e.ml.fail["FLAKY-PHOTO1"] = errors.New("model overloaded")
	if st := e.tick(); st.Failed != 1 {
		t.Fatalf("first: %+v", st)
	}
	// Not retried within the hour.
	if st := e.tick(); st.Failed != 0 || st.Scanned != 0 {
		t.Fatalf("retried too soon: %+v", st)
	}
	for i := 0; i < maxScanAttempts; i++ {
		e.exec(`UPDATE face_photo_scans SET scanned_at = now() - interval '2 hours'`)
		e.tick()
	}
	if n := e.count(`SELECT attempts FROM face_photo_scans WHERE drive_file_id = 'FLAKY-PHOTO1'`); n != maxScanAttempts {
		t.Fatalf("attempts: %d", n)
	}
	// Once the model recovers and the file changes, it is scanned again from scratch.
	delete(e.ml.fail, "FLAKY-PHOTO1")
	e.exec(`UPDATE photos SET md5 = 'new-content' WHERE drive_file_id = 'FLAKY-PHOTO1'`)
	if st := e.tick(); st.Scanned != 1 || st.Faces != 1 {
		t.Fatalf("after change: %+v", st)
	}
}

func TestAccountDeletionErasesFaceData(t *testing.T) {
	e := setup(t)
	e.user("erin@imt.fr")
	e.photo("ERIN-PHOTO-1", 10, on(alice, 0.1))
	e.enrol("erin@imt.fr", alice)
	e.tick()
	if got := e.album("erin@imt.fr"); len(got) != 1 {
		t.Fatalf("album before deletion: %v", got)
	}
	if err := accounts.NewService(e.db, nil).Delete(context.Background(), "erin@imt.fr"); err != nil {
		t.Fatal(err)
	}
	for _, table := range []string{"face_consents", "face_references", "face_matches"} {
		if n := e.count(`SELECT count(*) FROM ` + table); n != 0 {
			t.Fatalf("%s keeps %d rows after the account was deleted", table, n)
		}
	}
}

func TestOnlyOneWorkerRunsAtATime(t *testing.T) {
	e := setup(t)
	conn, err := e.db.Conn(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	if _, err := conn.ExecContext(context.Background(), `SELECT pg_advisory_lock($1, hashtext(current_schema()))`, advisoryLockKey); err != nil {
		t.Fatal(err)
	}
	if _, err := e.svc.Tick(context.Background()); !errors.Is(err, ErrWorkerRunning) {
		t.Fatalf("second worker: %v", err)
	}
	_, _ = conn.ExecContext(context.Background(), unlockSQL, advisoryLockKey)
}
