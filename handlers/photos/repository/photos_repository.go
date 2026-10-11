package repository

import (
	"context"
	"database/sql"
	"time"
)

type PhotosRepository struct {
	db *sql.DB
}

func NewPhotosRepository(db *sql.DB) *PhotosRepository {
	return &PhotosRepository{db: db}
}

// PhotoCursor is the keyset position of a photo listing. The gallery orders by (SortAt, AlbumKey, ID) so photos of
// one event stay together when several events share the same timestamp (event-dated photos all sit at noon).
type PhotoCursor struct {
	SortAt   time.Time
	AlbumKey int64
	ID       string
}

// AlbumCursor is the keyset position of the album listing.
type AlbumCursor struct {
	SortDate time.Time
	ID       int64
}

// PhotoRow carries only what a grid cell needs.
type PhotoRow struct {
	ID           string
	MediaType    string
	Width        int
	Height       int
	SortAt       time.Time
	AlbumID      sql.NullInt64
	AlbumTitle   sql.NullString // the event name; set only when the album is published, like AlbumID
	AlbumKey     int64          // raw album id (0 for none): the gallery's tie-break, not shown
	HasThumbnail bool
	DriveURL     string
}

// PhotoDetail is the full metadata of one file.
type PhotoDetail struct {
	PhotoRow
	Name        string
	MimeType    string
	SizeBytes   int64
	TakenAt     sql.NullTime
	TakenSource string
	FolderURL   string
	FolderPath  string
	Section     string
	AlbumTitle  sql.NullString
	AlbumURL    sql.NullString
}

type AlbumRow struct {
	ID           int64
	ParentID     sql.NullInt64
	Title        string
	EventDate    sql.NullTime
	SortDate     time.Time
	AcademicYear string
	Category     string
	DriveURL     string
	PhotoCount   int
	CoverID      sql.NullString
	FirstTakenAt sql.NullTime
	LastTakenAt  sql.NullTime
}

// visibleAlbum yields the album id only when that album is published: photos of albums under review stay in the
// gallery but never expose them.
const photoColumns = `
	p.drive_file_id, p.media_type, p.width, p.height, p.sort_at,
	CASE WHEN a.visibility = 'published' AND a.deleted_at IS NULL THEN a.id_photo_albums END,
	CASE WHEN a.visibility = 'published' AND a.deleted_at IS NULL THEN a.title END,
	COALESCE(p.album_id, 0),
	p.has_thumbnail, p.web_url
	FROM photos p
	LEFT JOIN photo_albums a ON a.id_photo_albums = p.album_id`

func scanPhoto(s interface{ Scan(...any) error }, p *PhotoRow) error {
	return s.Scan(&p.ID, &p.MediaType, &p.Width, &p.Height, &p.SortAt, &p.AlbumID, &p.AlbumTitle, &p.AlbumKey, &p.HasThumbnail, &p.DriveURL)
}

func albumKeyOf(c *PhotoCursor) int64 {
	if c == nil {
		return 0
	}
	return c.AlbumKey
}

// Gallery lists every visible photo and video, newest first, events kept together within one timestamp. Exact
// duplicates appear once. Uses idx_photos_gallery_event.
func (r *PhotosRepository) Gallery(ctx context.Context, after *PhotoCursor, limit int) ([]PhotoRow, error) {
	return r.photos(ctx, `SELECT `+photoColumns+`
		WHERE p.deleted_at IS NULL AND p.duplicate_of IS NULL AND p.media_type <> 'other'
		  AND ($1::timestamptz IS NULL OR (p.sort_at, COALESCE(p.album_id, 0), p.drive_file_id) < ($1, $4::bigint, $2))
		ORDER BY p.sort_at DESC, COALESCE(p.album_id, 0) DESC, p.drive_file_id DESC LIMIT $3`,
		after, limit, albumKeyOf(after))
}

// GalleryNewer returns up to limit photos strictly newer than the cursor, oldest first (the closest ones). Uses
// idx_photos_gallery_event backwards.
func (r *PhotosRepository) GalleryNewer(ctx context.Context, than *PhotoCursor, limit int) ([]PhotoRow, error) {
	return r.photos(ctx, `SELECT `+photoColumns+`
		WHERE p.deleted_at IS NULL AND p.duplicate_of IS NULL AND p.media_type <> 'other'
		  AND (p.sort_at, COALESCE(p.album_id, 0), p.drive_file_id) > ($1, $4::bigint, $2)
		ORDER BY p.sort_at, COALESCE(p.album_id, 0), p.drive_file_id LIMIT $3`, than, limit, albumKeyOf(than))
}

// TimelineDay is one calendar day (Paris time, the zone EXIF times are read in) that holds gallery photos.
type TimelineDay struct {
	Day    time.Time
	Count  int
	Latest time.Time
}

// Timeline lists the days that really hold photos, newest first, with the same filter as Gallery.
func (r *PhotosRepository) Timeline(ctx context.Context) ([]TimelineDay, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT (p.sort_at AT TIME ZONE 'Europe/Paris')::date AS day, COUNT(*), MAX(p.sort_at)
		FROM photos p
		WHERE p.deleted_at IS NULL AND p.duplicate_of IS NULL AND p.media_type <> 'other'
		GROUP BY 1 ORDER BY 1 DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := []TimelineDay{}
	for rows.Next() {
		var d TimelineDay
		if err := rows.Scan(&d.Day, &d.Count, &d.Latest); err != nil {
			return nil, err
		}
		out = append(out, d)
	}
	return out, rows.Err()
}

// AlbumPhotos lists an album newest first; a copy is hidden only when its original sits in the same album.
// Uses idx_photos_album backwards.
func (r *PhotosRepository) AlbumPhotos(ctx context.Context, albumID int64, after *PhotoCursor, limit int) ([]PhotoRow, error) {
	return r.photos(ctx, `SELECT `+photoColumns+`
		WHERE p.album_id = $4 AND p.deleted_at IS NULL AND p.media_type <> 'other'
		  AND (p.duplicate_of IS NULL OR NOT p.album_duplicate)
		  AND ($1::timestamptz IS NULL OR (p.sort_at, p.drive_file_id) < ($1, $2))
		ORDER BY p.sort_at DESC, p.drive_file_id DESC LIMIT $3`, after, limit, albumID)
}

func (r *PhotosRepository) photos(ctx context.Context, q string, after *PhotoCursor, limit int, extra ...any) ([]PhotoRow, error) {
	var ts any
	id := ""
	if after != nil {
		ts, id = after.SortAt, after.ID
	}
	rows, err := r.db.QueryContext(ctx, q, append([]any{ts, id, limit}, extra...)...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := []PhotoRow{}
	for rows.Next() {
		var p PhotoRow
		if err := scanPhoto(rows, &p); err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

// Photo returns one file with its album and folder links, or sql.ErrNoRows.
func (r *PhotosRepository) Photo(ctx context.Context, id string) (PhotoDetail, error) {
	var d PhotoDetail
	err := r.db.QueryRowContext(ctx, `
		SELECT p.drive_file_id, p.media_type, p.width, p.height, p.sort_at,
		       CASE WHEN a.visibility = 'published' AND a.deleted_at IS NULL THEN a.id_photo_albums END,
		       p.has_thumbnail, p.web_url,
		       p.name, p.mime_type, p.size_bytes, p.taken_at, p.taken_at_source,
		       COALESCE(f.web_url, ''), COALESCE(f.logical_path, ''), COALESCE(f.section_path, ''),
		       CASE WHEN a.visibility = 'published' AND a.deleted_at IS NULL THEN a.title END,
		       CASE WHEN a.visibility = 'published' AND a.deleted_at IS NULL THEN a.drive_url END
		FROM photos p
		LEFT JOIN photo_albums a ON a.id_photo_albums = p.album_id
		LEFT JOIN photo_folders f ON f.drive_folder_id = p.folder_id
		WHERE p.drive_file_id = $1 AND p.deleted_at IS NULL`, id).
		Scan(&d.ID, &d.MediaType, &d.Width, &d.Height, &d.SortAt, &d.AlbumID, &d.HasThumbnail, &d.DriveURL,
			&d.Name, &d.MimeType, &d.SizeBytes, &d.TakenAt, &d.TakenSource,
			&d.FolderURL, &d.FolderPath, &d.Section, &d.AlbumTitle, &d.AlbumURL)
	return d, err
}

const albumColumns = `id_photo_albums, parent_album_id, title, event_date, sort_date, academic_year, category, drive_url,
	photo_count, cover_photo_id, first_taken_at, last_taken_at FROM photo_albums`

// Albums lists published, non-empty albums by event date, newest first. Uses idx_photo_albums_list.
func (r *PhotosRepository) Albums(ctx context.Context, after *AlbumCursor, limit int) ([]AlbumRow, error) {
	var ts any
	var id int64
	if after != nil {
		ts, id = after.SortDate, after.ID
	}
	rows, err := r.db.QueryContext(ctx, `SELECT `+albumColumns+`
		WHERE visibility = 'published' AND deleted_at IS NULL AND photo_count > 0
		  AND ($1::date IS NULL OR (sort_date, id_photo_albums) < ($1, $2))
		ORDER BY sort_date DESC, id_photo_albums DESC LIMIT $3`, ts, id, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := []AlbumRow{}
	for rows.Next() {
		a, err := scanAlbum(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, a)
	}
	return out, rows.Err()
}

func (r *PhotosRepository) Album(ctx context.Context, id int64) (AlbumRow, error) {
	return scanAlbum(r.db.QueryRowContext(ctx, `SELECT `+albumColumns+`
		WHERE id_photo_albums = $1 AND visibility = 'published' AND deleted_at IS NULL`, id))
}

func scanAlbum(s interface{ Scan(...any) error }) (AlbumRow, error) {
	var a AlbumRow
	err := s.Scan(&a.ID, &a.ParentID, &a.Title, &a.EventDate, &a.SortDate, &a.AcademicYear, &a.Category, &a.DriveURL,
		&a.PhotoCount, &a.CoverID, &a.FirstTakenAt, &a.LastTakenAt)
	return a, err
}

// DataVersion changes whenever indexed data changes; response caches are keyed by it.
func (r *PhotosRepository) DataVersion(ctx context.Context) (int64, error) {
	var v int64
	err := r.db.QueryRowContext(ctx, `SELECT data_version FROM photo_sync_state WHERE id = 1`).Scan(&v)
	return v, err
}

type SyncStatus struct {
	LastFullAt        sql.NullTime
	LastIncrementalAt sql.NullTime
	FullInProgress    bool
	DataVersion       int64
	Runs              []SyncRun
}

type SyncRun struct {
	ID         int64
	Mode       string
	StartedAt  time.Time
	FinishedAt sql.NullTime
	Complete   bool
	Stats      string
	Error      sql.NullString
}

func (r *PhotosRepository) SyncStatus(ctx context.Context) (SyncStatus, error) {
	var s SyncStatus
	err := r.db.QueryRowContext(ctx, `
		SELECT last_full_at, last_incremental_at, full_run_id IS NOT NULL, data_version FROM photo_sync_state WHERE id = 1`).
		Scan(&s.LastFullAt, &s.LastIncrementalAt, &s.FullInProgress, &s.DataVersion)
	if err != nil {
		return s, err
	}
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, mode, started_at, finished_at, complete, stats::text, error
		FROM photo_sync_runs ORDER BY id DESC LIMIT 10`)
	if err != nil {
		return s, err
	}
	defer rows.Close()
	for rows.Next() {
		var run SyncRun
		if err := rows.Scan(&run.ID, &run.Mode, &run.StartedAt, &run.FinishedAt, &run.Complete, &run.Stats, &run.Error); err != nil {
			return s, err
		}
		s.Runs = append(s.Runs, run)
	}
	return s, rows.Err()
}
