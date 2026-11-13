package photobooth

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/plugimt/transat-backend/services/photos"
)

// The booth's pictures wait in the event's "photobooth_temp" folder, which the gallery never indexes. At the
// end of the evening an administrator or the Club Photo reviews them in the admin panel (and may remove some),
// then approves the folder: it is renamed "Photobooth", or merged into the event's existing "Photobooth"
// folder, and the gallery indexes it. Renaming the folder by hand on Drive publishes it just the same.

// PendingFolder is an event's photobooth folder awaiting approval.
type PendingFolder struct {
	FolderID   string     `json:"folder_id"`
	EventName  string     `json:"event_name"`
	EventDate  string     `json:"event_date"` // yyyy-mm-dd
	Pictures   int        `json:"pictures"`   // photos and collages still online
	Shoots     int        `json:"shoots"`
	OpenShoots int        `json:"open_shoots"` // guests still in front of the booth
	FirstAt    *time.Time `json:"first_upload_at"`
	LastAt     *time.Time `json:"last_upload_at"`
	DriveURL   string     `json:"drive_url"`
}

// PendingPicture is one picture of a pending folder.
type PendingPicture struct {
	ID         int64     `json:"id"`
	ShootID    int64     `json:"shoot_id"`
	Kind       string    `json:"kind"`
	Number     int       `json:"number"`
	Name       string    `json:"name"`
	Width      *int      `json:"width"`
	Height     *int      `json:"height"`
	UploadedAt time.Time `json:"uploaded_at"`
}

// ApproveResult tells where the approved pictures now are.
type ApproveResult struct {
	FolderID string `json:"folder_id"`
	Merged   bool   `json:"merged"` // moved into the event's existing "Photobooth" folder
	Moved    int    `json:"moved"`
}

const driveFolderURL = "https://drive.google.com/drive/folders/"

// ListPending returns the events whose photobooth pictures await approval, most recent first. A folder renamed
// by hand on Drive is no longer pending and drops out of the list.
func (s *Service) ListPending(ctx context.Context) ([]PendingFolder, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT ef.photobooth_folder, to_char(ef.event_date, 'YYYY-MM-DD'),
		       COALESCE((SELECT sh.event_name FROM photobooth_shoots sh WHERE sh.drive_folder_id = ef.photobooth_folder
		                 ORDER BY sh.id DESC LIMIT 1), ef.event_key),
		       (SELECT COUNT(*) FROM photobooth_files f JOIN photobooth_shoots sh ON sh.id = f.shoot_id
		         WHERE sh.drive_folder_id = ef.photobooth_folder AND f.status = 'uploaded'),
		       (SELECT COUNT(*) FROM photobooth_shoots sh WHERE sh.drive_folder_id = ef.photobooth_folder AND sh.status <> 'cancelled'),
		       (SELECT COUNT(*) FROM photobooth_shoots sh WHERE sh.drive_folder_id = ef.photobooth_folder AND sh.status = 'open'
		           AND sh.created_at > now() - interval '1 hour'),
		       (SELECT MIN(f.uploaded_at) FROM photobooth_files f JOIN photobooth_shoots sh ON sh.id = f.shoot_id
		         WHERE sh.drive_folder_id = ef.photobooth_folder AND f.status = 'uploaded'),
		       (SELECT MAX(f.uploaded_at) FROM photobooth_files f JOIN photobooth_shoots sh ON sh.id = f.shoot_id
		         WHERE sh.drive_folder_id = ef.photobooth_folder AND f.status = 'uploaded')
		FROM photobooth_event_folders ef
		ORDER BY ef.event_date DESC, ef.updated_at DESC`)
	if err != nil {
		return nil, err
	}
	var candidates []PendingFolder
	for rows.Next() {
		var p PendingFolder
		var first, last sql.NullTime
		if err := rows.Scan(&p.FolderID, &p.EventDate, &p.EventName, &p.Pictures, &p.Shoots, &p.OpenShoots, &first, &last); err != nil {
			rows.Close()
			return nil, err
		}
		if first.Valid {
			p.FirstAt = &first.Time
		}
		if last.Valid {
			p.LastAt = &last.Time
		}
		p.DriveURL = driveFolderURL + p.FolderID
		if p.Pictures > 0 {
			candidates = append(candidates, p)
		}
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return nil, err
	}

	out := []PendingFolder{}
	for _, p := range candidates {
		pending, err := s.isPendingFolder(ctx, p.FolderID)
		if err != nil {
			return nil, err
		}
		if pending {
			out = append(out, p)
		}
	}
	return out, nil
}

// PendingPictures lists the pictures of a pending folder in upload order.
func (s *Service) PendingPictures(ctx context.Context, folderID string) ([]PendingPicture, error) {
	if err := s.requirePending(ctx, folderID); err != nil {
		return nil, err
	}
	rows, err := s.db.QueryContext(ctx, `
		SELECT f.id, f.shoot_id, f.kind, f.number, f.file_name, f.width, f.height, COALESCE(f.uploaded_at, f.created_at)
		FROM photobooth_files f JOIN photobooth_shoots sh ON sh.id = f.shoot_id
		WHERE sh.drive_folder_id = $1 AND f.status = 'uploaded'
		ORDER BY COALESCE(f.uploaded_at, f.created_at), f.id`, folderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []PendingPicture{}
	for rows.Next() {
		var p PendingPicture
		var w, h sql.NullInt64
		if err := rows.Scan(&p.ID, &p.ShootID, &p.Kind, &p.Number, &p.Name, &w, &h, &p.UploadedAt); err != nil {
			return nil, err
		}
		if w.Valid && h.Valid {
			wi, hi := int(w.Int64), int(h.Int64)
			p.Width, p.Height = &wi, &hi
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

// OpenPendingPicture streams a picture under review (a preview when previewSize > 0).
func (s *Service) OpenPendingPicture(ctx context.Context, id int64, previewSize int) (SharedFile, error) {
	driveID, name, size, _, err := s.pictureForReview(ctx, id)
	if err != nil {
		return SharedFile{}, err
	}
	return s.openDriveFile(ctx, driveID, name, size, previewSize)
}

// DeletePendingPicture removes a picture before approval: it leaves Drive and the guests' share page.
func (s *Service) DeletePendingPicture(ctx context.Context, id int64, by string) error {
	driveID, name, _, folderID, err := s.pictureForReview(ctx, id)
	if err != nil {
		return err
	}
	if err := s.requirePending(ctx, folderID); err != nil {
		return err
	}
	res, err := s.db.ExecContext(ctx, `UPDATE photobooth_files SET status = 'trash' WHERE id = $1 AND status = 'uploaded'`, id)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return errNotFound
	}
	log.Printf("Photobooth: %s removed %s before approval", by, name)
	s.trashNow([]string{driveID})
	return nil
}

// Approve publishes a pending folder: renamed "Photobooth", or merged into the event's existing one. The
// gallery then indexes the pictures on its next sync.
func (s *Service) Approve(ctx context.Context, folderID, by string) (ApproveResult, error) {
	s.folderMu.Lock()
	defer s.folderMu.Unlock()

	var eventFolder string
	err := s.db.QueryRowContext(ctx, `SELECT event_folder_id FROM photobooth_event_folders WHERE photobooth_folder = $1`,
		folderID).Scan(&eventFolder)
	if errors.Is(err, sql.ErrNoRows) {
		return ApproveResult{}, errNotFound
	}
	if err != nil {
		return ApproveResult{}, err
	}
	if err := s.requirePending(ctx, folderID); err != nil {
		return ApproveResult{}, err
	}

	siblings, err := s.drive.ListChildFolders(ctx, eventFolder)
	if err != nil {
		return ApproveResult{}, err
	}
	target := ""
	for _, f := range siblings {
		if f.ID != folderID && NameKey(f.Name) == NameKey(PhotoboothFolderName) {
			target = f.ID
			break
		}
	}

	if target == "" {
		if err := s.drive.Rename(ctx, folderID, PhotoboothFolderName); err != nil {
			return ApproveResult{}, err
		}
		log.Printf("Photobooth: %s approved folder %s (renamed %q)", by, folderID, PhotoboothFolderName)
		return ApproveResult{FolderID: folderID}, nil
	}

	// A "Photobooth" folder exists already (an earlier approval of the same event): merge into it.
	kids, err := s.drive.ListChildren(ctx, folderID)
	if err != nil {
		return ApproveResult{}, err
	}
	for _, k := range kids {
		if err := s.drive.Move(ctx, k.ID, folderID, target); err != nil {
			return ApproveResult{}, err
		}
	}
	if _, err := s.db.ExecContext(ctx, `UPDATE photobooth_shoots SET drive_folder_id = $1 WHERE drive_folder_id = $2`,
		target, folderID); err != nil {
		return ApproveResult{}, err
	}
	if err := s.drive.Trash(ctx, folderID); err != nil {
		log.Printf("Photobooth: approved folder %s emptied but not trashed: %v", folderID, err)
	}
	log.Printf("Photobooth: %s approved folder %s (%d items merged into %s)", by, folderID, len(kids), target)
	return ApproveResult{FolderID: target, Merged: true, Moved: len(kids)}, nil
}

func (s *Service) pictureForReview(ctx context.Context, id int64) (driveID, name string, size int64, folderID string, err error) {
	err = s.db.QueryRowContext(ctx, `SELECT f.drive_file_id, f.file_name, f.size_bytes, sh.drive_folder_id
		FROM photobooth_files f JOIN photobooth_shoots sh ON sh.id = f.shoot_id
		WHERE f.id = $1 AND f.status = 'uploaded' AND f.drive_file_id IS NOT NULL`, id).Scan(&driveID, &name, &size, &folderID)
	if errors.Is(err, sql.ErrNoRows) {
		err = errNotFound
	}
	return
}

// requirePending accepts only a folder still awaiting approval.
func (s *Service) requirePending(ctx context.Context, folderID string) error {
	pending, err := s.isPendingFolder(ctx, folderID)
	if err != nil {
		return err
	}
	if !pending {
		return conflict("these pictures have already been approved")
	}
	return nil
}

func (s *Service) isPendingFolder(ctx context.Context, folderID string) (bool, error) {
	info, err := s.drive.GetFile(ctx, folderID)
	if errors.Is(err, ErrNotFound) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return !info.Trashed && photos.IsPendingPhotoboothFolder(info.Name), nil
}
