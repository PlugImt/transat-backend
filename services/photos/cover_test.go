package photos

import (
	"context"
	"database/sql"
	"testing"

	"github.com/plugimt/transat-backend/handlers/photos/repository"
	"github.com/plugimt/transat-backend/internal/database/dbtest"
)

func addPhoto(t *testing.T, db *sql.DB, id, name string, album int64, minutes int) {
	t.Helper()
	if _, err := db.Exec(`INSERT INTO photos (drive_file_id, album_id, name, mime_type, media_type, md5, sort_at, sort_at_source, web_url, last_seen_run)
		VALUES ($1, $2, $3, 'image/jpeg', 'image', $1, now() - make_interval(mins => $4), 'x', 'u', 1)`, id, album, name, minutes); err != nil {
		t.Fatal(err)
	}
}

func refresh(t *testing.T, db *sql.DB) {
	t.Helper()
	tx, err := db.Begin()
	if err != nil {
		t.Fatal(err)
	}
	if err := refreshDerived(context.Background(), tx); err != nil {
		t.Fatal(err)
	}
	if err := tx.Commit(); err != nil {
		t.Fatal(err)
	}
}

func TestACoverFileIsTheAlbumCoverAndNotOneOfItsPhotos(t *testing.T) {
	db := dbtest.Open(t)
	ctx := context.Background()
	var withCover, without int64
	for _, a := range []struct {
		title string
		id    *int64
	}{{"Gala", &withCover}, {"Rallye", &without}} {
		if err := db.QueryRow(`INSERT INTO photo_albums (root_folder_id, title, raw_folder_name, drive_url, visibility, confidence, rule, last_seen_run)
			VALUES ('root-'||$1, $1, $1, 'u', 'published', 1, 'R', 1) RETURNING id_photo_albums`, a.title).Scan(a.id); err != nil {
			t.Fatal(err)
		}
	}
	addPhoto(t, db, "GALA-PHOTO-1", "IMG_1.jpg", withCover, 30)
	addPhoto(t, db, "GALA-PHOTO-2", "IMG_2.jpg", withCover, 20)
	addPhoto(t, db, "GALA-COVER-0", "Cover.JPG", withCover, 10)
	addPhoto(t, db, "RALLYE-PHOTO", "IMG_9.jpg", without, 5)
	addPhoto(t, db, "NOT-A-COVER0", "discover.jpg", without, 6)
	refresh(t, db)

	repo := repository.NewPhotosRepository(db)
	gala, err := repo.Album(ctx, withCover)
	if err != nil {
		t.Fatal(err)
	}
	if gala.CoverID.String != "GALA-COVER-0" || gala.PhotoCount != 2 {
		t.Fatalf("album with a cover file: cover %q, %d photos", gala.CoverID.String, gala.PhotoCount)
	}
	rallye, _ := repo.Album(ctx, without)
	if rallye.CoverID.String == "" || rallye.PhotoCount != 2 {
		t.Fatalf("album without a cover file keeps one of its photos: %q, %d photos", rallye.CoverID.String, rallye.PhotoCount)
	}

	albumPhotos, err := repo.AlbumPhotos(ctx, withCover, nil, 50)
	if err != nil {
		t.Fatal(err)
	}
	gallery, err := repo.Gallery(ctx, nil, 50)
	if err != nil {
		t.Fatal(err)
	}
	for _, list := range [][]repository.PhotoRow{albumPhotos, gallery} {
		for _, p := range list {
			if p.ID == "GALA-COVER-0" {
				t.Fatal("the cover file is listed among the photos")
			}
		}
	}
	if len(albumPhotos) != 2 || len(gallery) != 4 {
		t.Fatalf("album lists %d photos, gallery %d", len(albumPhotos), len(gallery))
	}
}
