package photobooth

import (
	"context"
	"database/sql"
	"time"
)

// The booth's pictures seen from the gallery's albums: an album is an event folder of the Drive, and the booth files
// its pictures in a sub-folder of the event folder ("photobooth_temp" while they await approval, then "Photobooth").
// The admin panel shows them in the album, approved or not; the public pages never do before approval (the gallery
// does not index the pending folder).

// AlbumBoothFolder is one photobooth folder of an album, awaiting approval or approved, with its pictures.
type AlbumBoothFolder struct {
	FolderID  string           `json:"folder_id"`
	Pending   bool             `json:"pending"`
	EventName string           `json:"event_name"`
	EventDate string           `json:"event_date"`
	DriveURL  string           `json:"drive_url"`
	Pictures  []PendingPicture `json:"pictures"`
}

// AlbumBoothCount is how many booth pictures an album holds.
type AlbumBoothCount struct {
	AlbumID  int64 `json:"album_id"`
	Pending  int   `json:"pending"`
	Approved int   `json:"approved"`
}

type albumPicture struct {
	album  int64
	folder string
	event  string
	date   string
	pic    PendingPicture
}

// albumPictures lists the booth's online pictures filed in the event folder of an album (every album when albumID
// is 0), in upload order. A shoot belongs to the album whose event folder its event resolved to.
func (s *Service) albumPictures(ctx context.Context, albumID int64) ([]albumPicture, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT a.id_photo_albums, ef.event_key, sh.drive_folder_id, sh.event_name, to_char(sh.event_date, 'YYYY-MM-DD'),
		       f.id, f.shoot_id, f.kind, f.number, f.file_name, f.width, f.height, COALESCE(f.uploaded_at, f.created_at)
		FROM photobooth_files f
		JOIN photobooth_shoots sh ON sh.id = f.shoot_id AND sh.status <> 'cancelled'
		JOIN photobooth_event_folders ef ON ef.event_date = sh.event_date
		JOIN photo_albums a ON a.root_folder_id = ef.event_folder_id AND a.deleted_at IS NULL
		WHERE f.status = 'uploaded' AND ($1 = 0 OR a.id_photo_albums = $1)
		ORDER BY COALESCE(f.uploaded_at, f.created_at), f.id`, albumID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []albumPicture
	for rows.Next() {
		var p albumPicture
		var key string
		var w, h sql.NullInt64
		var at time.Time
		if err := rows.Scan(&p.album, &key, &p.folder, &p.event, &p.date, &p.pic.ID, &p.pic.ShootID, &p.pic.Kind,
			&p.pic.Number, &p.pic.Name, &w, &h, &at); err != nil {
			return nil, err
		}
		// Several events may share a day: only the folder of this shoot's own event counts.
		if NameKey(p.event) != key {
			continue
		}
		if w.Valid && h.Valid {
			wi, hi := int(w.Int64), int(h.Int64)
			p.pic.Width, p.pic.Height = &wi, &hi
		}
		p.pic.UploadedAt = at
		out = append(out, p)
	}
	return out, rows.Err()
}

// pendingFolders tells, for each folder, whether it still awaits approval (one Drive lookup per folder).
func (s *Service) pendingFolders(ctx context.Context, pics []albumPicture) (map[string]bool, error) {
	pending := map[string]bool{}
	for _, p := range pics {
		if _, seen := pending[p.folder]; seen {
			continue
		}
		ok, err := s.isPendingFolder(ctx, p.folder)
		if err != nil {
			return nil, err
		}
		pending[p.folder] = ok
	}
	return pending, nil
}

// AlbumPictures returns the booth's pictures of an album, folder by folder: those awaiting approval first.
func (s *Service) AlbumPictures(ctx context.Context, albumID int64) ([]AlbumBoothFolder, error) {
	if albumID <= 0 {
		return nil, errNotFound
	}
	pics, err := s.albumPictures(ctx, albumID)
	if err != nil {
		return nil, err
	}
	pending, err := s.pendingFolders(ctx, pics)
	if err != nil {
		return nil, err
	}
	out := []AlbumBoothFolder{}
	index := map[string]int{}
	for _, p := range pics {
		i, ok := index[p.folder]
		if !ok {
			i = len(out)
			index[p.folder] = i
			out = append(out, AlbumBoothFolder{FolderID: p.folder, Pending: pending[p.folder], EventName: p.event,
				EventDate: p.date, DriveURL: driveFolderURL + p.folder, Pictures: []PendingPicture{}})
		}
		out[i].Pictures = append(out[i].Pictures, p.pic)
	}
	// Awaiting approval first, then the approved ones, each in upload order.
	pendingFirst := make([]AlbumBoothFolder, 0, len(out))
	for _, want := range []bool{true, false} {
		for _, f := range out {
			if f.Pending == want {
				pendingFirst = append(pendingFirst, f)
			}
		}
	}
	return pendingFirst, nil
}

// AlbumCounts returns, for every album holding booth pictures, how many await approval and how many are approved.
func (s *Service) AlbumCounts(ctx context.Context) ([]AlbumBoothCount, error) {
	pics, err := s.albumPictures(ctx, 0)
	if err != nil {
		return nil, err
	}
	pending, err := s.pendingFolders(ctx, pics)
	if err != nil {
		return nil, err
	}
	out := []AlbumBoothCount{}
	index := map[int64]int{}
	for _, p := range pics {
		i, ok := index[p.album]
		if !ok {
			i = len(out)
			index[p.album] = i
			out = append(out, AlbumBoothCount{AlbumID: p.album})
		}
		if pending[p.folder] {
			out[i].Pending++
		} else {
			out[i].Approved++
		}
	}
	return out, nil
}

// linkPending fills the panel event and the gallery album of each pending folder.
func (s *Service) linkPending(ctx context.Context, folders []PendingFolder) error {
	if len(folders) == 0 {
		return nil
	}
	events := map[eventKey]int64{}
	if err := s.eachRow(ctx, `SELECT id, name, to_char(event_date, 'YYYY-MM-DD') FROM photobooth_events ORDER BY id`, nil,
		func(r *sql.Rows) error {
			var id int64
			var name, date string
			if err := r.Scan(&id, &name, &date); err != nil {
				return err
			}
			k := eventKey{date, NameKey(name)}
			if _, taken := events[k]; !taken {
				events[k] = id
			}
			return nil
		}); err != nil {
		return err
	}
	albums := map[string]int64{}
	if err := s.eachRow(ctx, `SELECT ef.photobooth_folder, a.id_photo_albums FROM photobooth_event_folders ef
		JOIN photo_albums a ON a.root_folder_id = ef.event_folder_id AND a.deleted_at IS NULL`, nil,
		func(r *sql.Rows) error {
			var folder string
			var id int64
			if err := r.Scan(&folder, &id); err != nil {
				return err
			}
			albums[folder] = id
			return nil
		}); err != nil {
		return err
	}
	for i := range folders {
		f := &folders[i]
		if id, ok := events[eventKey{f.EventDate, f.eventKey}]; ok {
			f.EventID = &id
		}
		if id, ok := albums[f.FolderID]; ok {
			f.AlbumID = &id
		}
	}
	return nil
}
