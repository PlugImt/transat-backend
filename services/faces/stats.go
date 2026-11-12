package faces

import (
	"context"
	"database/sql"
	"time"
)

// AdminStats is the face recognition overview for the admin panel. Counts only: no user is named.
type AdminStats struct {
	Available bool `json:"available"`
	// Users who opted in, and those of them who sent at least one selfie.
	UsersEnabled    int `json:"users_enabled"`
	UsersWithSelfie int `json:"users_with_selfie"`
	// Users whose album has at least one photo, and the average album size among users with a selfie.
	UsersWithPhotos     int     `json:"users_with_photos"`
	AveragePhotosByUser float64 `json:"average_photos_by_user"`
	RejectedMatches     int     `json:"rejected_matches"`
	// Library: eligible photos (within FACES_SCOPE_FOLDERS when set), analysed, failed for good, still to do.
	PhotosTotal    int        `json:"photos_total"`
	PhotosAnalysed int        `json:"photos_analysed"`
	PhotosFailed   int        `json:"photos_failed"`
	PhotosPending  int        `json:"photos_pending"`
	FacesFound     int        `json:"faces_found"`
	LastScanAt     *time.Time `json:"last_scan_at"`
	Scoped         bool       `json:"scoped"` // the library is limited to some folders (dev)
}

func (s *Service) AdminStats(ctx context.Context) (AdminStats, error) {
	st := AdminStats{Available: s.Available(), Scoped: len(s.cfg.ScopeFolders) > 0}
	if err := s.db.QueryRowContext(ctx, `
		SELECT (SELECT count(*) FROM face_consents),
		       (SELECT count(DISTINCT email) FROM face_references),
		       (SELECT count(*) FROM face_matches WHERE rejected_at IS NOT NULL),
		       (SELECT count(*) FROM photo_faces f JOIN photos p ON p.drive_file_id = f.drive_file_id
		        WHERE `+inScope(1)+`)`, s.scopeParam()).
		Scan(&st.UsersEnabled, &st.UsersWithSelfie, &st.RejectedMatches, &st.FacesFound); err != nil {
		return st, err
	}
	if err := s.db.QueryRowContext(ctx, `
		WITH albums AS (
			SELECT r.email, count(DISTINCT p.drive_file_id) AS photos
			FROM (SELECT DISTINCT email FROM face_references) r
			LEFT JOIN face_matches m ON m.email = r.email AND m.rejected_at IS NULL
			LEFT JOIN photos p ON p.drive_file_id = m.drive_file_id AND p.deleted_at IS NULL
			GROUP BY r.email)
		SELECT count(*) FILTER (WHERE photos > 0), COALESCE(avg(photos), 0) FROM albums`).
		Scan(&st.UsersWithPhotos, &st.AveragePhotosByUser); err != nil {
		return st, err
	}

	model := defaultModel
	if s.ml != nil {
		model = s.ml.Model()
	}
	var last sql.NullTime
	if err := s.db.QueryRowContext(ctx, `
		SELECT count(*),
		       count(*) FILTER (WHERE sc.model = $1 AND sc.md5 = p.md5 AND sc.status IN ('done', 'skipped')),
		       count(*) FILTER (WHERE sc.model = $1 AND sc.md5 = p.md5 AND sc.status = 'failed' AND sc.attempts >= $3),
		       max(sc.scanned_at)
		FROM photos p
		LEFT JOIN face_photo_scans sc ON sc.drive_file_id = p.drive_file_id
		WHERE `+eligiblePhotos+` AND `+inScope(2), model, s.scopeParam(), maxScanAttempts).
		Scan(&st.PhotosTotal, &st.PhotosAnalysed, &st.PhotosFailed, &last); err != nil {
		return st, err
	}
	st.PhotosPending = st.PhotosTotal - st.PhotosAnalysed - st.PhotosFailed
	if last.Valid {
		st.LastScanAt = &last.Time
	}
	return st, nil
}
