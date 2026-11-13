package photos_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/gofiber/fiber/v2"
	reportHandler "github.com/plugimt/transat-backend/handlers/photoreport"
	"github.com/plugimt/transat-backend/internal/database/dbtest"
	"github.com/plugimt/transat-backend/middlewares"
	"github.com/plugimt/transat-backend/models"
	"github.com/plugimt/transat-backend/routes"
	"github.com/plugimt/transat-backend/services/photoreport"
	"github.com/plugimt/transat-backend/services/photos"
	"github.com/plugimt/transat-backend/utils"
)

// fakeMedia records what was read from "Drive" and answers with the file id as the body.
type fakeMedia struct {
	mu    sync.Mutex
	reads []string
}

func (m *fakeMedia) record(s string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.reads = append(m.reads, s)
}

func (m *fakeMedia) Thumbnail(_ context.Context, id string, size int) (photos.Media, error) {
	m.record(fmt.Sprintf("thumb:%s:%d", id, size))
	return photos.Media{Body: io.NopCloser(strings.NewReader("thumb-" + id)), ContentType: "image/jpeg", Size: -1}, nil
}

func (m *fakeMedia) Original(_ context.Context, id string) (photos.Media, error) {
	m.record("original:" + id)
	body := "original-" + id
	return photos.Media{Body: io.NopCloser(strings.NewReader(body)), ContentType: "application/octet-stream", Size: int64(len(body))}, nil
}

func (m *fakeMedia) count() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.reads)
}

type mailer struct{}

func (mailer) SendEmail(models.Email, interface{}) error { return nil }

type env struct {
	t     *testing.T
	db    *sql.DB
	app   *fiber.App
	media *fakeMedia
}

func newEnv(t *testing.T) *env {
	t.Setenv("JWT_SECRET", "test-secret-test-secret-test-secret")
	t.Setenv("PUBLIC_SITE_URL", "https://transat.test/")
	db := dbtest.Open(t)
	middlewares.InitJWTMiddleware(db)
	e := &env{t: t, db: db, app: fiber.New(), media: &fakeMedia{}}
	routes.SetupPhotosRoutes(e.app, db, nil, e.media, reportHandler.NewHandler(photoreport.NewService(db, mailer{}, "")))
	return e
}

func (e *env) exec(q string, args ...any) {
	e.t.Helper()
	if _, err := e.db.Exec(q, args...); err != nil {
		e.t.Fatalf("%s: %v", q, err)
	}
}

func (e *env) user(email string, roles ...string) {
	e.exec(`INSERT INTO newf (email, password, first_name, last_name) VALUES ($1, 'x', 'A', 'B')`, email)
	for _, r := range roles {
		e.exec(`INSERT INTO roles (name) VALUES ($1) ON CONFLICT (name) DO NOTHING`, r)
		e.exec(`INSERT INTO newf_roles (email, id_roles) SELECT $1, id_roles FROM roles WHERE name = $2`, email, r)
	}
}

// album creates an album with the given state; its photos are added with photo.
func (e *env) album(title, visibility string, hidden bool) int64 {
	e.t.Helper()
	var id int64
	if err := e.db.QueryRow(`INSERT INTO photo_albums (root_folder_id, title, raw_folder_name, event_date, drive_url, visibility,
		confidence, rule, last_seen_run, hidden, photo_count, sort_date)
		VALUES ('root-'||$1, $1, $1, '2026-03-12', 'https://drive.google.com/drive/folders/root-'||$1, $2, 1, 'R', 1, $3, 1, '2026-03-12')
		RETURNING id_photo_albums`, title, visibility, hidden).Scan(&id); err != nil {
		e.t.Fatal(err)
	}
	return id
}

func (e *env) photo(id string, album int64, name, media string, minutes int) {
	e.exec(`INSERT INTO photos (drive_file_id, album_id, name, mime_type, media_type, width, height, sort_at, sort_at_source, web_url, last_seen_run)
		VALUES ($1, $2, $3, CASE WHEN $4 = 'image' THEN 'image/jpeg' ELSE 'video/mp4' END, $4, 4000, 3000,
		        '2026-03-12 20:00+00'::timestamptz + make_interval(mins => $5), 'x', 'https://drive.google.com/file/d/'||$1, 1)`,
		id, album, name, media, minutes)
}

func (e *env) setCover(album int64, photo string) {
	e.exec(`UPDATE photo_albums SET cover_photo_id = $2 WHERE id_photo_albums = $1`, album, photo)
}

type reply struct {
	status int
	header http.Header
	body   string
}

func (e *env) call(method, path, as string, body string) reply {
	e.t.Helper()
	var r io.Reader
	if body != "" {
		r = strings.NewReader(body)
	}
	req := httptest.NewRequest(method, path, r)
	req.Header.Set("Content-Type", "application/json")
	if as != "" {
		token, err := utils.GenerateJWT(as, nil, "")
		if err != nil {
			e.t.Fatal(err)
		}
		req.Header.Set("Authorization", "Bearer "+token)
	}
	res, err := e.app.Test(req, -1)
	if err != nil {
		e.t.Fatal(err)
	}
	defer res.Body.Close()
	raw, _ := io.ReadAll(res.Body)
	return reply{res.StatusCode, res.Header, string(raw)}
}

func (e *env) get(path string) reply { return e.call(http.MethodGet, path, "", "") }

func decode[T any](t *testing.T, r reply) T {
	t.Helper()
	var v T
	if r.status != http.StatusOK {
		t.Fatalf("status %d: %s", r.status, r.body)
	}
	if err := json.Unmarshal([]byte(r.body), &v); err != nil {
		t.Fatalf("%v: %s", err, r.body)
	}
	return v
}

type publicAlbum struct {
	ID         int64   `json:"id"`
	Title      string  `json:"title"`
	Slug       string  `json:"slug"`
	PhotoCount int     `json:"photo_count"`
	CoverID    *string `json:"cover_id"`
	ShareURL   string  `json:"share_url"`
}

type publicPage struct {
	Photos []struct {
		ID     string `json:"id"`
		Width  int    `json:"width"`
		Height int    `json:"height"`
	} `json:"photos"`
	NextCursor *string `json:"next_cursor"`
}

func TestPublicAlbumServesOnlyShownAlbumsWithoutDriveLinks(t *testing.T) {
	e := newEnv(t)
	gala := e.album("Soirée de Gala", "published", false)
	e.photo("GALA-PHOTO-01", gala, "IMG_1.jpg", "image", 1)
	e.photo("GALA-PHOTO-02", gala, "IMG_2.jpg", "image", 2)
	e.photo("GALA-VIDEO-01", gala, "VID_1.mp4", "video", 3)
	e.photo("GALA-COVER-01", gala, "cover.jpg", "image", 4)
	e.setCover(gala, "GALA-COVER-01")
	hidden := e.album("Caché", "published", true)
	e.photo("HIDDEN-PHOTO1", hidden, "IMG_3.jpg", "image", 1)
	review := e.album("À revoir", "needs_review", false)
	e.photo("REVIEW-PHOTO1", review, "IMG_4.jpg", "image", 1)

	// The share link's whole segment is accepted: the id decides, the slug is decoration.
	a := decode[publicAlbum](t, e.get(fmt.Sprintf("/public/albums/%d-n-importe-quoi", gala)))
	if a.ID != gala || a.Title != "Soirée de Gala" || a.Slug != "soiree-de-gala" {
		t.Fatalf("album = %+v", a)
	}
	if want := fmt.Sprintf("https://transat.test/albums/%d-soiree-de-gala", gala); a.ShareURL != want {
		t.Fatalf("share url = %q, want %q", a.ShareURL, want)
	}
	if a.CoverID == nil || *a.CoverID != "GALA-COVER-01" {
		t.Fatalf("cover = %v", a.CoverID)
	}
	raw := e.get(fmt.Sprintf("/public/albums/%d", gala))
	if strings.Contains(raw.body, "drive.google.com") {
		t.Fatalf("public album leaks a Drive link: %s", raw.body)
	}
	if !strings.HasPrefix(raw.header.Get("Cache-Control"), "public") {
		t.Fatalf("cache-control = %q", raw.header.Get("Cache-Control"))
	}

	// Images only: no video, no cover file.
	page := decode[publicPage](t, e.get(fmt.Sprintf("/public/albums/%d/photos", gala)))
	if len(page.Photos) != 2 || page.Photos[0].ID != "GALA-PHOTO-02" || page.Photos[1].ID != "GALA-PHOTO-01" || page.NextCursor != nil {
		t.Fatalf("photos = %+v", page)
	}
	if page.Photos[0].Width != 4000 || page.Photos[0].Height != 3000 {
		t.Fatalf("dimensions = %+v", page.Photos[0])
	}

	for _, id := range []int64{hidden, review, 99999} {
		for _, path := range []string{"", "/photos", "/files", "/cover"} {
			if r := e.get(fmt.Sprintf("/public/albums/%d%s", id, path)); r.status != http.StatusNotFound {
				t.Errorf("album %d%s: status %d, want 404", id, path, r.status)
			}
		}
	}
	for _, path := range []string{"/public/albums/abc", "/public/albums/-1", "/public/albums/0/photos"} {
		if r := e.get(path); r.status != http.StatusNotFound {
			t.Errorf("%s: status %d, want 404", path, r.status)
		}
	}
}

func TestPublicPhotosArePaginatedByCursor(t *testing.T) {
	e := newEnv(t)
	album := e.album("Rallye", "published", false)
	for i := range 5 {
		e.photo(fmt.Sprintf("RALLYE-PHOTO-%d", i), album, "x.jpg", "image", i)
	}
	seen := map[string]bool{}
	path := fmt.Sprintf("/public/albums/%d/photos?limit=2", album)
	for pages := 0; ; pages++ {
		if pages > 5 {
			t.Fatal("pagination does not end")
		}
		p := decode[publicPage](t, e.get(path))
		for _, ph := range p.Photos {
			if seen[ph.ID] {
				t.Fatalf("%s listed twice", ph.ID)
			}
			seen[ph.ID] = true
		}
		if p.NextCursor == nil {
			break
		}
		path = fmt.Sprintf("/public/albums/%d/photos?limit=2&cursor=%s", album, *p.NextCursor)
	}
	if len(seen) != 5 {
		t.Fatalf("saw %d photos, want 5", len(seen))
	}
}

func TestPublicImagesAreStreamedOnlyForTheirOwnShownAlbum(t *testing.T) {
	e := newEnv(t)
	gala := e.album("Gala", "published", false)
	e.photo("GALA-PHOTO-01", gala, "IMG 1 é.jpg", "image", 1)
	e.photo("GALA-VIDEO-01", gala, "VID_1.mp4", "video", 2)
	e.photo("GALA-COVER-01", gala, "Cover.jpg", "image", 3)
	e.setCover(gala, "GALA-COVER-01")
	other := e.album("Rallye", "published", false)
	e.photo("RALLYE-PHOTO1", other, "IMG_2.jpg", "image", 1)
	hidden := e.album("Caché", "published", true)
	e.photo("HIDDEN-PHOTO1", hidden, "IMG_3.jpg", "image", 1)
	e.setCover(hidden, "HIDDEN-PHOTO1")

	grid := e.get(fmt.Sprintf("/public/albums/%d/photos/GALA-PHOTO-01/grid", gala))
	if grid.status != http.StatusOK || grid.body != "thumb-GALA-PHOTO-01" || grid.header.Get("Content-Type") != "image/jpeg" {
		t.Fatalf("grid: %d %q %v", grid.status, grid.body, grid.header)
	}
	preview := e.get(fmt.Sprintf("/public/albums/%d-gala/photos/GALA-PHOTO-01/preview", gala))
	if preview.status != http.StatusOK {
		t.Fatalf("preview: %d", preview.status)
	}
	original := e.get(fmt.Sprintf("/public/albums/%d/photos/GALA-PHOTO-01/original?download=1", gala))
	if original.status != http.StatusOK || original.body != "original-GALA-PHOTO-01" {
		t.Fatalf("original: %d %q", original.status, original.body)
	}
	// The type recorded at indexing wins over Drive's generic one; the file keeps its name.
	if ct := original.header.Get("Content-Type"); ct != "image/jpeg" {
		t.Fatalf("original content-type = %q", ct)
	}
	if cd := original.header.Get("Content-Disposition"); !strings.HasPrefix(cd, `attachment; filename="IMG 1 _.jpg"`) || !strings.Contains(cd, "IMG%201%20%C3%A9.jpg") {
		t.Fatalf("content-disposition = %q", cd)
	}
	if cover := e.get(fmt.Sprintf("/public/albums/%d/cover", gala)); cover.status != http.StatusOK || cover.body != "thumb-GALA-COVER-01" {
		t.Fatalf("cover: %d %q", cover.status, cover.body)
	}
	before := e.media.count()

	denied := []string{
		fmt.Sprintf("/public/albums/%d/photos/RALLYE-PHOTO1/original", gala),   // another album's photo
		fmt.Sprintf("/public/albums/%d/photos/HIDDEN-PHOTO1/original", hidden), // hidden album
		fmt.Sprintf("/public/albums/%d/photos/HIDDEN-PHOTO1/grid", hidden),
		fmt.Sprintf("/public/albums/%d/cover", hidden),
		fmt.Sprintf("/public/albums/%d/photos/GALA-VIDEO-01/original", gala), // videos are not served
		fmt.Sprintf("/public/albums/%d/photos/GALA-PHOTO-01/full", gala),     // unknown variant
		fmt.Sprintf("/public/albums/%d/photos/short/grid", gala),             // not a Drive id
	}
	for _, path := range denied {
		if r := e.get(path); r.status != http.StatusNotFound {
			t.Errorf("%s: status %d, want 404", path, r.status)
		}
	}
	if e.media.count() != before {
		t.Fatalf("Drive was read for a refused request: %v", e.media.reads[before:])
	}

	// Without Drive credentials the metadata still answers, images do not.
	plain := fiber.New()
	routes.SetupPhotosRoutes(plain, e.db, nil, nil, reportHandler.NewHandler(photoreport.NewService(e.db, mailer{}, "")))
	res, err := plain.Test(httptest.NewRequest(http.MethodGet, fmt.Sprintf("/public/albums/%d/photos/GALA-PHOTO-01/grid", gala), nil), -1)
	if err != nil || res.StatusCode != http.StatusServiceUnavailable {
		t.Fatalf("without media: %v %v", res.StatusCode, err)
	}
}

func TestPublicAlbumFilesListTheImagesForTheArchive(t *testing.T) {
	e := newEnv(t)
	gala := e.album("Gala", "published", false)
	e.photo("GALA-PHOTO-02", gala, "b.jpg", "image", 2)
	e.photo("GALA-PHOTO-01", gala, "a.jpg", "image", 1)
	e.photo("GALA-VIDEO-01", gala, "v.mp4", "video", 3)
	e.photo("GALA-COVER-01", gala, "cover.png", "image", 4)
	files := decode[struct {
		Slug  string `json:"slug"`
		Files []struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"files"`
	}](t, e.get(fmt.Sprintf("/public/albums/%d/files", gala)))
	if files.Slug != "gala" || len(files.Files) != 2 || files.Files[0].ID != "GALA-PHOTO-01" || files.Files[1].Name != "b.jpg" {
		t.Fatalf("files = %+v", files)
	}
}

func TestHidingAnAlbumRemovesItEverywhereAndSurvivesTheSync(t *testing.T) {
	e := newEnv(t)
	e.user("admin@x.test", "ADMIN")
	e.user("club@x.test", photoreport.PhotoClubRole)
	e.user("student@x.test", "NEWF")
	gala := e.album("Gala", "published", false)
	e.photo("GALA-PHOTO-01", gala, "a.jpg", "image", 1)
	e.album("Rallye", "needs_review", false)

	hide := func(as string, hidden bool) reply {
		return e.call(http.MethodPut, fmt.Sprintf("/admin/photos/albums/%d/hidden", gala), as, fmt.Sprintf(`{"hidden": %v}`, hidden))
	}
	if r := hide("student@x.test", true); r.status != http.StatusForbidden {
		t.Fatalf("a student hid an album: %d", r.status)
	}
	if r := hide("", true); r.status != http.StatusUnauthorized {
		t.Fatalf("anonymous: %d", r.status)
	}
	if r := e.call(http.MethodPut, fmt.Sprintf("/admin/photos/albums/%d/hidden", gala), "admin@x.test", `{}`); r.status != http.StatusBadRequest {
		t.Fatalf("missing flag: %d", r.status)
	}
	if r := e.call(http.MethodPut, "/admin/photos/albums/99999/hidden", "admin@x.test", `{"hidden": true}`); r.status != http.StatusNotFound {
		t.Fatalf("unknown album: %d", r.status)
	}

	var version0 int64
	_ = e.db.QueryRow(`SELECT data_version FROM photo_sync_state WHERE id = 1`).Scan(&version0)

	// Fill the public cache before hiding: hiding must invalidate it.
	if r := e.get(fmt.Sprintf("/public/albums/%d", gala)); r.status != http.StatusOK {
		t.Fatalf("before hiding: %d", r.status)
	}

	type adminAlbum struct {
		ID        int64   `json:"id"`
		Hidden    bool    `json:"hidden"`
		HiddenBy  *string `json:"hidden_by"`
		Shown     bool    `json:"shown"`
		ShareURL  string  `json:"share_url"`
		CoverPath *string `json:"cover_path"`
	}
	got := decode[adminAlbum](t, hide("club@x.test", true))
	if !got.Hidden || got.Shown || got.HiddenBy == nil || *got.HiddenBy != "club@x.test" {
		t.Fatalf("after hiding: %+v", got)
	}
	// The panel still shows the cover of a hidden album, to its own users only.
	e.setCover(gala, "GALA-PHOTO-01")
	if r := e.call(http.MethodGet, fmt.Sprintf("/admin/photos/albums/%d/cover", gala), "club@x.test", ""); r.status != http.StatusOK || r.body != "thumb-GALA-PHOTO-01" {
		t.Fatalf("admin cover: %d %q", r.status, r.body)
	}
	if r := e.call(http.MethodGet, fmt.Sprintf("/admin/photos/albums/%d/cover", gala), "student@x.test", ""); r.status != http.StatusForbidden {
		t.Fatalf("a student read an admin cover: %d", r.status)
	}
	var version1 int64
	_ = e.db.QueryRow(`SELECT data_version FROM photo_sync_state WHERE id = 1`).Scan(&version1)
	if version1 <= version0 {
		t.Fatal("hiding did not invalidate the gallery caches")
	}

	if r := e.get(fmt.Sprintf("/public/albums/%d", gala)); r.status != http.StatusNotFound {
		t.Fatalf("hidden album still public: %d", r.status)
	}
	student := func(path string) string {
		r := e.call(http.MethodGet, path, "student@x.test", "")
		if r.status != http.StatusOK && r.status != http.StatusNotFound {
			t.Fatalf("%s: %d %s", path, r.status, r.body)
		}
		return fmt.Sprint(r.status, r.body)
	}
	if body := student("/photos/albums"); strings.Contains(body, "Gala") {
		t.Fatalf("hidden album listed in the app: %s", body)
	}
	if body := student("/photos"); strings.Contains(body, "GALA-PHOTO-01") {
		t.Fatalf("photo of a hidden album in the app gallery: %s", body)
	}
	if body := student(fmt.Sprintf("/photos/albums/%d", gala)); !strings.HasPrefix(body, "404") {
		t.Fatalf("hidden album opened in the app: %s", body)
	}
	if body := student("/photos/GALA-PHOTO-01"); !strings.HasPrefix(body, "404") {
		t.Fatalf("photo of a hidden album opened in the app: %s", body)
	}

	// The sync rewrites the album row on every change; the administrator's choice stays.
	e.exec(`INSERT INTO photo_albums (root_folder_id, title, raw_folder_name, drive_url, visibility, confidence, rule, last_seen_run)
		VALUES ('root-Gala', 'Gala renamed', 'Gala', 'u', 'published', 1, 'R', 2)
		ON CONFLICT (root_folder_id) DO UPDATE SET title = EXCLUDED.title, visibility = EXCLUDED.visibility, last_seen_run = EXCLUDED.last_seen_run`)
	var stillHidden bool
	_ = e.db.QueryRow(`SELECT hidden FROM photo_albums WHERE id_photo_albums = $1`, gala).Scan(&stillHidden)
	if !stillHidden {
		t.Fatal("the sync unhid the album")
	}

	// Listing: counts by status, filters and search.
	type list struct {
		Albums []adminAlbum   `json:"albums"`
		Counts map[string]int `json:"counts"`
		Total  int            `json:"total"`
	}
	all := decode[list](t, e.call(http.MethodGet, "/admin/photos/albums", "admin@x.test", ""))
	if len(all.Albums) != 2 || all.Counts["all"] != 2 || all.Counts["hidden"] != 1 || all.Counts["review"] != 1 || all.Counts["shown"] != 0 {
		t.Fatalf("list = %+v", all)
	}
	onlyHidden := decode[list](t, e.call(http.MethodGet, "/admin/photos/albums?status=hidden", "club@x.test", ""))
	if len(onlyHidden.Albums) != 1 || onlyHidden.Albums[0].ID != gala || onlyHidden.Total != 1 {
		t.Fatalf("hidden filter = %+v", onlyHidden)
	}
	search := decode[list](t, e.call(http.MethodGet, "/admin/photos/albums?q=ralLYE", "admin@x.test", ""))
	if len(search.Albums) != 1 || search.Albums[0].ID == gala {
		t.Fatalf("search = %+v", search)
	}
	if r := e.call(http.MethodGet, "/admin/photos/albums?q=%25", "admin@x.test", ""); len(decode[list](t, r).Albums) != 0 {
		t.Fatal("a % in the search matched everything")
	}
	if r := e.call(http.MethodGet, "/admin/photos/albums", "student@x.test", ""); r.status != http.StatusForbidden {
		t.Fatalf("a student listed albums: %d", r.status)
	}

	// Showing it again brings everything back.
	got = decode[adminAlbum](t, hide("admin@x.test", false))
	if got.Hidden || !got.Shown || got.HiddenBy != nil {
		t.Fatalf("after showing: %+v", got)
	}
	if r := e.get(fmt.Sprintf("/public/albums/%d", gala)); r.status != http.StatusOK {
		t.Fatalf("shown album not public: %d", r.status)
	}
	if body := student("/photos"); !strings.Contains(body, "GALA-PHOTO-01") {
		t.Fatalf("photo missing from the gallery once shown again: %s", body)
	}
}

func TestAppAlbumResponsesCarryTheTransatLink(t *testing.T) {
	e := newEnv(t)
	e.user("student@x.test", "NEWF")
	gala := e.album("Soirée de Gala", "published", false)
	e.photo("GALA-PHOTO-01", gala, "a.jpg", "image", 1)
	want := fmt.Sprintf("https://transat.test/albums/%d-soiree-de-gala", gala)
	for _, path := range []string{"/photos/albums", fmt.Sprintf("/photos/albums/%d", gala), "/photos/GALA-PHOTO-01"} {
		r := e.call(http.MethodGet, path, "student@x.test", "")
		if r.status != http.StatusOK || !strings.Contains(r.body, `"share_url":"`+want+`"`) {
			t.Errorf("%s: %d %s", path, r.status, r.body)
		}
	}
}

func TestPhotoManagersOpenThePhotoSectionsOfThePanel(t *testing.T) {
	e := newEnv(t)
	routes.SetupAccessRoutes(e.app, e.db)
	e.user("manager@x.test", "NEWF", photoreport.PhotoManagerRole)
	e.user("club@x.test", "NEWF", photoreport.PhotoClubRole)
	e.user("student@x.test", "NEWF")
	gala := e.album("Gala", "published", false)

	type access struct {
		Admin    bool     `json:"admin"`
		Sections []string `json:"sections"`
	}
	for _, who := range []string{"manager@x.test", "club@x.test"} {
		got := decode[access](t, e.call(http.MethodGet, "/admin/access", who, ""))
		want := []string{"dashboard", "photo-reports", "photobooth", "photo-albums"}
		if got.Admin || strings.Join(got.Sections, ",") != strings.Join(want, ",") {
			t.Fatalf("%s: access = %+v", who, got)
		}
	}
	if got := decode[access](t, e.call(http.MethodGet, "/admin/access", "student@x.test", "")); len(got.Sections) != 0 {
		t.Fatalf("student: access = %+v", got)
	}

	if r := e.call(http.MethodGet, "/admin/photos/albums", "manager@x.test", ""); r.status != http.StatusOK {
		t.Fatalf("albums as photo manager: %d %s", r.status, r.body)
	}
	if r := e.call(http.MethodPut, fmt.Sprintf("/admin/photos/albums/%d/hidden", gala), "manager@x.test", `{"hidden": true}`); r.status != http.StatusOK {
		t.Fatalf("hiding as photo manager: %d %s", r.status, r.body)
	}
}
