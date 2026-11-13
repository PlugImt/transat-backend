package repository

import (
	"context"
	"database/sql"
	"strings"
	"time"
)

// PhotoFile is what streaming one file of a public album needs.
type PhotoFile struct {
	ID           string
	Name         string
	MimeType     string
	MediaType    string
	SizeBytes    int64
	HasThumbnail bool
}

// AlbumPhotoFile returns a file of a shown album: one of its photos or videos, or its cover. Anything else, a
// file of a hidden or unpublished album included, is sql.ErrNoRows.
func (r *PhotosRepository) AlbumPhotoFile(ctx context.Context, albumID int64, photoID string) (PhotoFile, error) {
	var f PhotoFile
	err := r.db.QueryRowContext(ctx, `
		SELECT p.drive_file_id, p.name, p.mime_type, p.media_type, p.size_bytes, p.has_thumbnail
		FROM photos p JOIN photo_albums a ON a.id_photo_albums = p.album_id
		WHERE p.drive_file_id = $1 AND p.album_id = $2 AND `+shownAlbum+`
		  AND p.deleted_at IS NULL AND p.media_type <> 'other'
		  AND (NOT p.is_cover OR p.drive_file_id = a.cover_photo_id)`, photoID, albumID).
		Scan(&f.ID, &f.Name, &f.MimeType, &f.MediaType, &f.SizeBytes, &f.HasThumbnail)
	return f, err
}

// AlbumImagesPage is AlbumPhotos restricted to images: what the public album page shows.
func (r *PhotosRepository) AlbumImagesPage(ctx context.Context, albumID int64, after *PhotoCursor, limit int) ([]PhotoRow, error) {
	return r.photos(ctx, `SELECT `+photoColumns+`
		WHERE p.album_id = $4 AND `+shownAlbum+` AND p.deleted_at IS NULL AND p.media_type = 'image' AND NOT p.is_cover
		  AND (p.duplicate_of IS NULL OR NOT p.album_duplicate)
		  AND ($1::timestamptz IS NULL OR (p.sort_at, p.drive_file_id) < ($1, $2))
		ORDER BY p.sort_at DESC, p.drive_file_id DESC LIMIT $3`, after, limit, albumID)
}

// AlbumImages lists the images of a shown album, oldest first, for the whole-album download: the photos the page
// shows, without videos, covers or copies.
func (r *PhotosRepository) AlbumImages(ctx context.Context, albumID int64) ([]PhotoFile, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT p.drive_file_id, p.name, p.mime_type, p.media_type, p.size_bytes, p.has_thumbnail
		FROM photos p JOIN photo_albums a ON a.id_photo_albums = p.album_id
		WHERE p.album_id = $1 AND `+shownAlbum+` AND p.deleted_at IS NULL AND p.media_type = 'image' AND NOT p.is_cover
		  AND (p.duplicate_of IS NULL OR NOT p.album_duplicate)
		ORDER BY p.sort_at, p.drive_file_id`, albumID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []PhotoFile{}
	for rows.Next() {
		var f PhotoFile
		if err := rows.Scan(&f.ID, &f.Name, &f.MimeType, &f.MediaType, &f.SizeBytes, &f.HasThumbnail); err != nil {
			return nil, err
		}
		out = append(out, f)
	}
	return out, rows.Err()
}

// ---- administration

// AdminAlbumRow is an album as the admin panel lists it, whatever its state.
type AdminAlbumRow struct {
	ID           int64
	ParentID     sql.NullInt64
	Title        string
	EventDate    sql.NullTime
	AcademicYear string
	Category     string
	DriveURL     string
	Visibility   string
	PhotoCount   int
	CoverID      sql.NullString
	Hidden       bool
	HiddenAt     sql.NullTime
	HiddenBy     sql.NullString
	FirstTakenAt sql.NullTime
}

// AdminAlbumFilter selects albums in the admin panel. Status is "", "shown", "hidden" or "review".
type AdminAlbumFilter struct {
	Query  string
	Status string
	Limit  int
	Offset int
}

type AdminAlbumCounts struct {
	All    int `json:"all"`
	Shown  int `json:"shown"`
	Hidden int `json:"hidden"`
	Review int `json:"review"`
}

const adminAlbumColumns = `a.id_photo_albums, a.parent_album_id, a.title, a.event_date, a.academic_year, a.category, a.drive_url,
	a.visibility, a.photo_count, a.cover_photo_id, a.hidden, a.hidden_at, a.hidden_by, a.first_taken_at
	FROM photo_albums a`

func scanAdminAlbum(s interface{ Scan(...any) error }) (AdminAlbumRow, error) {
	var a AdminAlbumRow
	err := s.Scan(&a.ID, &a.ParentID, &a.Title, &a.EventDate, &a.AcademicYear, &a.Category, &a.DriveURL,
		&a.Visibility, &a.PhotoCount, &a.CoverID, &a.Hidden, &a.HiddenAt, &a.HiddenBy, &a.FirstTakenAt)
	return a, err
}

// statusCondition is the WHERE clause of a status filter, on photo_albums aliased "a".
func statusCondition(status string) string {
	switch status {
	case "shown":
		return shownAlbum
	case "hidden":
		return "a.hidden"
	case "review":
		return "a.visibility <> 'published' AND NOT a.hidden"
	}
	return "TRUE"
}

// AdminAlbums lists the albums that are not deleted, newest event first, with the count of each status for the
// same search.
func (r *PhotosRepository) AdminAlbums(ctx context.Context, f AdminAlbumFilter) ([]AdminAlbumRow, AdminAlbumCounts, error) {
	var counts AdminAlbumCounts
	search := "%" + escapeLike(strings.TrimSpace(f.Query)) + "%"
	const base = `a.deleted_at IS NULL AND ($1 = '%%' OR a.title ILIKE $1 OR a.raw_folder_name ILIKE $1 OR a.category ILIKE $1)`
	if err := r.db.QueryRowContext(ctx, `
		SELECT count(*),
		       count(*) FILTER (WHERE `+shownAlbum+`),
		       count(*) FILTER (WHERE a.hidden),
		       count(*) FILTER (WHERE a.visibility <> 'published' AND NOT a.hidden)
		FROM photo_albums a WHERE `+base, search).Scan(&counts.All, &counts.Shown, &counts.Hidden, &counts.Review); err != nil {
		return nil, counts, err
	}

	rows, err := r.db.QueryContext(ctx, `SELECT `+adminAlbumColumns+`
		WHERE `+base+` AND `+statusCondition(f.Status)+`
		ORDER BY a.sort_date DESC, a.id_photo_albums DESC LIMIT $2 OFFSET $3`, search, f.Limit, f.Offset)
	if err != nil {
		return nil, counts, err
	}
	defer rows.Close()
	out := []AdminAlbumRow{}
	for rows.Next() {
		a, err := scanAdminAlbum(rows)
		if err != nil {
			return nil, counts, err
		}
		out = append(out, a)
	}
	return out, counts, rows.Err()
}

// AdminAlbum returns one album that is not deleted, or sql.ErrNoRows.
func (r *PhotosRepository) AdminAlbum(ctx context.Context, id int64) (AdminAlbumRow, error) {
	return scanAdminAlbum(r.db.QueryRowContext(ctx, `SELECT `+adminAlbumColumns+`
		WHERE a.id_photo_albums = $1 AND a.deleted_at IS NULL`, id))
}

// SetAlbumHidden hides or shows an album, remembering who did it, and invalidates every cached gallery response.
func (r *PhotosRepository) SetAlbumHidden(ctx context.Context, id int64, hidden bool, by string, now time.Time) (AdminAlbumRow, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return AdminAlbumRow{}, err
	}
	defer func() { _ = tx.Rollback() }()

	res, err := tx.ExecContext(ctx, `
		UPDATE photo_albums SET hidden = $2,
		       hidden_at = CASE WHEN $2 THEN $3::timestamptz END,
		       hidden_by = CASE WHEN $2 THEN NULLIF($4, '') END
		WHERE id_photo_albums = $1 AND deleted_at IS NULL AND hidden <> $2`, id, hidden, now, by)
	if err != nil {
		return AdminAlbumRow{}, err
	}
	if n, _ := res.RowsAffected(); n > 0 {
		if _, err := tx.ExecContext(ctx, `UPDATE photo_sync_state SET data_version = data_version + 1 WHERE id = 1`); err != nil {
			return AdminAlbumRow{}, err
		}
	}
	a, err := scanAdminAlbum(tx.QueryRowContext(ctx, `SELECT `+adminAlbumColumns+`
		WHERE a.id_photo_albums = $1 AND a.deleted_at IS NULL`, id))
	if err != nil {
		return AdminAlbumRow{}, err
	}
	return a, tx.Commit()
}

func escapeLike(s string) string {
	return strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`).Replace(s)
}
