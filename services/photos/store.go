package photos

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/lib/pq"
)

const (
	batchSize      = 1000
	driveFolderURL = "https://drive.google.com/drive/folders/"
	driveFileURL   = "https://drive.google.com/file/d/%s/view"
)

type queryer interface {
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}

type folderEntry struct {
	AlbumID   int64
	EventDate *time.Time
}

// folderIndex maps a Drive folder ID to the album it belongs to; it is how a file finds its album.
type folderIndex map[string]folderEntry

func loadFolderIndex(ctx context.Context, q queryer) (folderIndex, error) {
	rows, err := q.QueryContext(ctx, `
		SELECT f.drive_folder_id, COALESCE(f.album_id, 0), a.event_date
		FROM photo_folders f LEFT JOIN photo_albums a ON a.id_photo_albums = f.album_id
		WHERE f.deleted_at IS NULL`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	idx := folderIndex{}
	for rows.Next() {
		var id string
		var e folderEntry
		var d sql.NullTime
		if err := rows.Scan(&id, &e.AlbumID, &d); err != nil {
			return nil, err
		}
		if d.Valid {
			t := d.Time
			e.EventDate = &t
		}
		idx[id] = e
	}
	return idx, rows.Err()
}

func (a folderIndex) sameAs(b folderIndex) bool {
	for id, ea := range a {
		eb, ok := b[id]
		if !ok {
			continue // folders that appear or vanish are handled through their own file changes
		}
		if ea.AlbumID != eb.AlbumID || (ea.EventDate == nil) != (eb.EventDate == nil) ||
			(ea.EventDate != nil && !ea.EventDate.Equal(*eb.EventDate)) {
			return false
		}
	}
	return true
}

func upsertAlbums(ctx context.Context, tx *sql.Tx, snap *Snapshot, run int64) (map[string]int64, error) {
	ids := map[string]int64{}
	for _, f := range snap.Folders {
		r := snap.Results[f.ID]
		if r.Kind != KindEvent {
			continue
		}
		visibility := "needs_review"
		if r.Publishable() {
			visibility = "published"
		}
		var date any
		if r.EventDate != nil {
			date = r.EventDate.Format("2006-01-02")
		}
		var id int64
		err := tx.QueryRowContext(ctx, `
			INSERT INTO photo_albums (root_folder_id, title, raw_folder_name, event_date, event_date_source, academic_year,
			                          category, drive_url, visibility, confidence, rule, last_seen_run)
			VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)
			ON CONFLICT (root_folder_id) DO UPDATE SET
			  title = EXCLUDED.title, raw_folder_name = EXCLUDED.raw_folder_name, event_date = EXCLUDED.event_date,
			  event_date_source = EXCLUDED.event_date_source, academic_year = EXCLUDED.academic_year,
			  category = EXCLUDED.category, drive_url = EXCLUDED.drive_url, visibility = EXCLUDED.visibility,
			  confidence = EXCLUDED.confidence, rule = EXCLUDED.rule, last_seen_run = EXCLUDED.last_seen_run,
			  parent_album_id = NULL, deleted_at = NULL
			RETURNING id_photo_albums`,
			f.ID, r.Title, f.EffectiveName, date, string(r.EventDateSource), r.AcademicYear, r.Category,
			driveFolderURL+f.ID, visibility, r.Confidence, string(r.Rule), run).Scan(&id)
		if err != nil {
			return nil, err
		}
		ids[f.ID] = id
	}
	for root, id := range ids {
		if parent := snap.Results[root].ParentAlbumRootID; parent != "" {
			if _, err := tx.ExecContext(ctx, `UPDATE photo_albums SET parent_album_id = $1 WHERE id_photo_albums = $2`, ids[parent], id); err != nil {
				return nil, err
			}
		}
	}
	return ids, nil
}

func upsertFolders(ctx context.Context, tx *sql.Tx, snap *Snapshot, albumIDs map[string]int64, run int64) error {
	for _, f := range snap.Folders {
		r := snap.Results[f.ID]
		var album any
		if id, ok := albumIDs[r.AlbumRootID]; ok && r.AlbumRootID != "" {
			album = id
		}
		webURL := f.WebURL
		if webURL == "" {
			webURL = driveFolderURL + f.ID
		}
		_, err := tx.ExecContext(ctx, `
			INSERT INTO photo_folders (drive_folder_id, parent_drive_id, name, effective_name, logical_path, shortcut_id, web_url,
			    drive_modified_at, kind, rule, confidence, status, signals, classifier_version, academic_year, category,
			    section_path, album_id, last_seen_run)
			VALUES ($1, NULLIF($2,''), $3, $4, $5, NULLIF($6,''), $7, NULLIF($8,'')::timestamptz, $9, $10, $11, $12, $13, $14,
			        $15, $16, $17, $18, $19)
			ON CONFLICT (drive_folder_id) DO UPDATE SET
			  parent_drive_id = EXCLUDED.parent_drive_id, name = EXCLUDED.name, effective_name = EXCLUDED.effective_name,
			  logical_path = EXCLUDED.logical_path, shortcut_id = EXCLUDED.shortcut_id, web_url = EXCLUDED.web_url,
			  drive_modified_at = EXCLUDED.drive_modified_at, kind = EXCLUDED.kind, rule = EXCLUDED.rule,
			  confidence = EXCLUDED.confidence, status = EXCLUDED.status, signals = EXCLUDED.signals,
			  classifier_version = EXCLUDED.classifier_version, academic_year = EXCLUDED.academic_year,
			  category = EXCLUDED.category, section_path = EXCLUDED.section_path, album_id = EXCLUDED.album_id,
			  last_seen_run = EXCLUDED.last_seen_run, deleted_at = NULL`,
			f.ID, f.ParentID, f.Name, f.EffectiveName, f.Path, f.ShortcutID, webURL, timeText(f.ModifiedTime),
			string(r.Kind), string(r.Rule), r.Confidence, string(r.Status), pq.Array(nonNil(r.Flags)), ClassifierVersion,
			r.AcademicYear, r.Category, r.SectionPath, album, run)
		if err != nil {
			return err
		}
	}
	return nil
}

// sweepStructure soft-deletes folders and albums absent from the latest structure listing, and their photos.
func sweepStructure(ctx context.Context, tx *sql.Tx, run int64) error {
	for _, table := range []string{"photo_folders", "photo_albums"} {
		if _, err := tx.ExecContext(ctx,
			`UPDATE `+table+` SET deleted_at = now() WHERE last_seen_run <> $1 AND deleted_at IS NULL`, run); err != nil {
			return err
		}
	}
	_, err := tx.ExecContext(ctx, `
		UPDATE photos SET deleted_at = now()
		WHERE deleted_at IS NULL AND folder_id IN (SELECT drive_folder_id FROM photo_folders WHERE deleted_at IS NOT NULL)`)
	return err
}

// syncAlbumMembership re-points photos at their folder's album after folders moved or were renamed.
func syncAlbumMembership(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `
		UPDATE photos p SET album_id = f.album_id
		FROM photo_folders f
		WHERE p.folder_id = f.drive_folder_id AND p.album_id IS DISTINCT FROM f.album_id`)
	return err
}

type photoInput struct {
	Item     DriveItem
	FolderID string // "" when the file sits at the drive root or in an unknown folder
}

type photoRow struct {
	id, folder, name, mime, kind, md5, takenSrc, sortSrc, url string
	album, size, width, height                                int64
	taken, created, modified, sortAt                          string
	suspect, hasThumb                                         bool
}

func buildRows(inputs []photoInput, idx folderIndex, now time.Time) []photoRow {
	rows := make([]photoRow, 0, len(inputs))
	for _, in := range inputs {
		m := in.Item
		entry, known := idx[in.FolderID]
		folder := ""
		if known {
			folder = in.FolderID
		}
		t := ResolveTimes(PhotoTimes{ExifTime: m.ExifTime, Name: m.Name, EventDate: entry.EventDate,
			DriveCreated: m.CreatedTime, DriveModified: m.ModifiedTime}, now)
		row := photoRow{
			id: m.ID, folder: folder, name: m.Name, mime: m.MimeType, kind: string(KindOfMime(m.MimeType)), md5: m.MD5,
			album: entry.AlbumID, size: m.Size, width: int64(m.Width), height: int64(m.Height),
			takenSrc: string(t.TakenSource), suspect: t.TakenSuspect, sortSrc: string(t.SortAtSource),
			sortAt: t.SortAt.UTC().Format(time.RFC3339Nano), created: timeText(m.CreatedTime), modified: timeText(m.ModifiedTime),
			url: m.WebViewLink, hasThumb: m.HasThumbnail,
		}
		if t.TakenAt != nil {
			row.taken = t.TakenAt.UTC().Format(time.RFC3339Nano)
		}
		if row.url == "" {
			row.url = fmt.Sprintf(driveFileURL, m.ID)
		}
		rows = append(rows, row)
	}
	return rows
}

func upsertPhotos(ctx context.Context, tx *sql.Tx, rows []photoRow, run int64) error {
	for start := 0; start < len(rows); start += batchSize {
		chunk := rows[start:min(start+batchSize, len(rows))]
		var c struct {
			ids, folders, names, mimes, kinds, md5s, takenSrc, sortSrc, urls, taken, created, modified, sortAt []string
			albums, sizes, widths, heights                                                                     []int64
			suspects, thumbs                                                                                   []bool
		}
		for _, r := range chunk {
			c.ids, c.folders, c.names, c.mimes, c.kinds, c.md5s = append(c.ids, r.id), append(c.folders, r.folder), append(c.names, r.name), append(c.mimes, r.mime), append(c.kinds, r.kind), append(c.md5s, r.md5)
			c.takenSrc, c.sortSrc, c.urls = append(c.takenSrc, r.takenSrc), append(c.sortSrc, r.sortSrc), append(c.urls, r.url)
			c.taken, c.created, c.modified, c.sortAt = append(c.taken, r.taken), append(c.created, r.created), append(c.modified, r.modified), append(c.sortAt, r.sortAt)
			c.albums, c.sizes, c.widths, c.heights = append(c.albums, r.album), append(c.sizes, r.size), append(c.widths, r.width), append(c.heights, r.height)
			c.suspects, c.thumbs = append(c.suspects, r.suspect), append(c.thumbs, r.hasThumb)
		}
		_, err := tx.ExecContext(ctx, `
			INSERT INTO photos (drive_file_id, folder_id, album_id, name, mime_type, media_type, size_bytes, md5, width, height,
			    taken_at, taken_at_source, taken_at_suspect, drive_created_at, drive_modified_at, sort_at, sort_at_source,
			    web_url, has_thumbnail, last_seen_run)
			SELECT t.id, NULLIF(t.folder,''), NULLIF(t.album,0), t.name, t.mime, t.kind, t.size, t.md5, t.w, t.h,
			    NULLIF(t.taken,'')::timestamptz, t.tsrc, t.susp, NULLIF(t.cr,'')::timestamptz, NULLIF(t.md,'')::timestamptz,
			    t.sort_at::timestamptz, t.ssrc, t.url, t.thumb, $19
			FROM unnest($1::text[], $2::text[], $3::bigint[], $4::text[], $5::text[], $6::text[], $7::bigint[], $8::text[],
			            $9::bigint[], $10::bigint[], $11::text[], $12::text[], $13::bool[], $14::text[], $15::text[],
			            $16::text[], $17::text[], $18::text[], $20::bool[])
			  AS t(id, folder, album, name, mime, kind, size, md5, w, h, taken, tsrc, susp, cr, md, sort_at, ssrc, url, thumb)
			ON CONFLICT (drive_file_id) DO UPDATE SET
			  folder_id = EXCLUDED.folder_id, album_id = EXCLUDED.album_id, name = EXCLUDED.name, mime_type = EXCLUDED.mime_type,
			  media_type = EXCLUDED.media_type, size_bytes = EXCLUDED.size_bytes, md5 = EXCLUDED.md5, width = EXCLUDED.width,
			  height = EXCLUDED.height, taken_at = EXCLUDED.taken_at, taken_at_source = EXCLUDED.taken_at_source,
			  taken_at_suspect = EXCLUDED.taken_at_suspect, drive_created_at = EXCLUDED.drive_created_at,
			  drive_modified_at = EXCLUDED.drive_modified_at, sort_at = EXCLUDED.sort_at, sort_at_source = EXCLUDED.sort_at_source,
			  web_url = EXCLUDED.web_url, has_thumbnail = EXCLUDED.has_thumbnail,
			  last_seen_run = EXCLUDED.last_seen_run, deleted_at = NULL`,
			pq.Array(c.ids), pq.Array(c.folders), pq.Array(c.albums), pq.Array(c.names), pq.Array(c.mimes), pq.Array(c.kinds),
			pq.Array(c.sizes), pq.Array(c.md5s), pq.Array(c.widths), pq.Array(c.heights), pq.Array(c.taken), pq.Array(c.takenSrc),
			pq.Array(c.suspects), pq.Array(c.created), pq.Array(c.modified), pq.Array(c.sortAt), pq.Array(c.sortSrc), pq.Array(c.urls),
			run, pq.Array(c.thumbs))
		if err != nil {
			return err
		}
	}
	return nil
}

func softDeletePhotos(ctx context.Context, tx *sql.Tx, ids []string) (int64, error) {
	if len(ids) == 0 {
		return 0, nil
	}
	res, err := tx.ExecContext(ctx,
		`UPDATE photos SET deleted_at = now() WHERE drive_file_id = ANY($1) AND deleted_at IS NULL`, pq.Array(ids))
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func sweepPhotos(ctx context.Context, tx *sql.Tx, run int64) error {
	_, err := tx.ExecContext(ctx, `UPDATE photos SET deleted_at = now() WHERE last_seen_run <> $1 AND deleted_at IS NULL`, run)
	return err
}

// refreshDerived recomputes exact-duplicate flags and album statistics from the indexed rows.
func refreshDerived(ctx context.Context, tx *sql.Tx) error {
	// One canonical file per md5: originals before "Copie de …", then oldest, then lowest ID. A copy is also
	// hidden inside an album when its original sits in the same album.
	if _, err := tx.ExecContext(ctx, `
		WITH ranked AS (
		  SELECT drive_file_id, album_id,
		         FIRST_VALUE(drive_file_id) OVER w AS canon,
		         FIRST_VALUE(album_id) OVER w AS canon_album,
		         ROW_NUMBER() OVER w AS rn
		  FROM photos
		  WHERE deleted_at IS NULL AND md5 <> '' AND media_type <> 'other'
		  WINDOW w AS (PARTITION BY md5 ORDER BY (name LIKE 'Copie de %' OR name LIKE 'Copy of %'),
		                                         drive_created_at NULLS LAST, drive_file_id)
		), target AS (
		  SELECT drive_file_id,
		         CASE WHEN rn = 1 THEN NULL ELSE canon END AS dup_of,
		         (rn > 1 AND album_id IS NOT NULL AND canon_album IS NOT NULL AND album_id = canon_album) AS album_dup
		  FROM ranked
		)
		UPDATE photos p SET duplicate_of = t.dup_of, album_duplicate = t.album_dup
		FROM target t
		WHERE p.drive_file_id = t.drive_file_id
		  AND (p.duplicate_of IS DISTINCT FROM t.dup_of OR p.album_duplicate <> t.album_dup)`); err != nil {
		return err
	}
	_, err := tx.ExecContext(ctx, `
		UPDATE photo_albums a SET
		  photo_count = COALESCE(s.n, 0), first_taken_at = s.first, last_taken_at = s.last, cover_photo_id = s.cover,
		  sort_date = COALESCE(a.event_date, (s.first AT TIME ZONE 'Europe/Paris')::date, DATE '1970-01-01')
		FROM photo_albums b
		LEFT JOIN (
		  SELECT album_id, COUNT(*) AS n,
		         MIN(taken_at) FILTER (WHERE NOT taken_at_suspect) AS first,
		         MAX(taken_at) FILTER (WHERE NOT taken_at_suspect) AS last,
		         (ARRAY_AGG(drive_file_id ORDER BY (media_type = 'image' AND has_thumbnail) DESC, sort_at, drive_file_id))[1] AS cover
		  FROM photos
		  WHERE deleted_at IS NULL AND media_type <> 'other' AND album_id IS NOT NULL
		    AND (duplicate_of IS NULL OR NOT album_duplicate)
		  GROUP BY album_id
		) s ON s.album_id = b.id_photo_albums
		WHERE a.id_photo_albums = b.id_photo_albums`)
	return err
}

func bumpVersion(ctx context.Context, q queryer) error {
	_, err := q.ExecContext(ctx, `UPDATE photo_sync_state SET data_version = data_version + 1 WHERE id = 1`)
	return err
}

func timeText(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.UTC().Format(time.RFC3339Nano)
}

func nonNil(s []string) []string {
	if s == nil {
		return []string{}
	}
	return s
}
