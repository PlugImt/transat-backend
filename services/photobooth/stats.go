package photobooth

import (
	"cmp"
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"regexp"
	"slices"
	"time"

	"github.com/lib/pq"
)

// Audience of the public pages: who opened the share link of a print (its QR code) and the event's public album.
// The site's pages report each opening with a random id the visitor's browser keeps; only its hash is stored.

const (
	VisitShare = "photobooth_share"
	VisitAlbum = "photo_album"
)

var visitorRe = regexp.MustCompile(`^[A-Za-z0-9-]{16,64}$`)

// ValidVisitor tells whether a visitor id is one the site's pages make.
func ValidVisitor(v string) bool { return visitorRe.MatchString(v) }

func visitorHash(v string) string {
	sum := sha256.Sum256([]byte(v))
	return hex.EncodeToString(sum[:])
}

// RecordShareVisit counts an opening of a share link (a scan of a print's QR code).
func (s *Service) RecordShareVisit(ctx context.Context, token, visitor string) error {
	if !shareTokenRe.MatchString(token) || !ValidVisitor(visitor) {
		return errNotFound
	}
	var id int64
	err := s.db.QueryRowContext(ctx, `SELECT id FROM photobooth_shoots WHERE share_token = $1 AND status = 'finalized'`, token).Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		return errNotFound
	}
	if err != nil {
		return err
	}
	return s.recordVisit(ctx, VisitShare, id, visitor)
}

// RecordAlbumVisit counts an opening of a public album page.
func (s *Service) RecordAlbumVisit(ctx context.Context, albumID int64, visitor string) error {
	if albumID <= 0 || !ValidVisitor(visitor) {
		return errNotFound
	}
	var ok bool
	err := s.db.QueryRowContext(ctx, `SELECT TRUE FROM photo_albums WHERE id_photo_albums = $1 AND deleted_at IS NULL`, albumID).Scan(&ok)
	if errors.Is(err, sql.ErrNoRows) {
		return errNotFound
	}
	if err != nil {
		return err
	}
	return s.recordVisit(ctx, VisitAlbum, albumID, visitor)
}

func (s *Service) recordVisit(ctx context.Context, kind string, target int64, visitor string) error {
	_, err := s.db.ExecContext(ctx, `INSERT INTO public_visits (kind, target_id, visitor) VALUES ($1, $2, $3)
		ON CONFLICT (kind, target_id, visitor) DO UPDATE SET views = public_visits.views + 1, last_seen = now()`,
		kind, target, visitorHash(visitor))
	return err
}

// EventStats is the audience of one evening: the booth's sessions, the people who scanned their QR codes, and the
// people who opened the event's public album.
type EventStats struct {
	Date    string `json:"date"`
	Name    string `json:"name"`
	EventID *int64 `json:"event_id,omitempty"` // the event prepared in the admin panel, when there is one
	// Sessions shared by the booth (each print carries its QR code).
	Sessions int `json:"sessions"`
	// Sessions whose QR code was opened at least once.
	ScannedSessions int `json:"scanned_sessions"`
	// Different people who opened at least one of the evening's QR codes, and how many times they did.
	Scanners  int `json:"scanners"`
	ScanViews int `json:"scan_views"`
	// The public album of the evening's Drive folder, once the gallery knows it.
	AlbumID       *int64 `json:"album_id,omitempty"`
	AlbumTitle    string `json:"album_title,omitempty"`
	AlbumVisitors int    `json:"album_visitors"`
	AlbumViews    int    `json:"album_views"`
}

// StatsTotals sums every evening; Scanners and AlbumVisitors count each person once.
type StatsTotals struct {
	Events          int `json:"events"`
	Sessions        int `json:"sessions"`
	ScannedSessions int `json:"scanned_sessions"`
	Scanners        int `json:"scanners"`
	ScanViews       int `json:"scan_views"`
	AlbumVisitors   int `json:"album_visitors"`
	AlbumViews      int `json:"album_views"`
	// Over the evenings with at least one session.
	AvgScannersPerEvent      float64 `json:"avg_scanners_per_event"`
	AvgAlbumVisitorsPerEvent float64 `json:"avg_album_visitors_per_event"`
	// Share of the sessions whose QR code was opened (0..1), and people per session.
	ScanRate              float64 `json:"scan_rate"`
	AvgScannersPerSession float64 `json:"avg_scanners_per_session"`
}

type Stats struct {
	Events []EventStats `json:"events"`
	Totals StatsTotals  `json:"totals"`
	// When the first opening was recorded: earlier evenings have sessions but no audience.
	TrackingSince *time.Time `json:"tracking_since,omitempty"`
}

type eventKey struct{ date, name string }

// Stats gathers the audience of every evening the booth worked, and of the panel's past events, newest first.
func (s *Service) Stats(ctx context.Context) (Stats, error) {
	byKey := map[eventKey]*EventStats{}
	scanners := map[eventKey]map[string]bool{}
	get := func(date time.Time, name string) *EventStats {
		k := eventKey{date.Format("2006-01-02"), NameKey(name)}
		if e := byKey[k]; e != nil {
			return e
		}
		e := &EventStats{Date: k.date, Name: name}
		byKey[k] = e
		scanners[k] = map[string]bool{}
		return e
	}

	// Panel events already started, even without a session yet.
	today := s.now().Format("2006-01-02")
	if err := s.eachRow(ctx, `SELECT id, name, event_date FROM photobooth_events WHERE event_date <= $1`, []any{today},
		func(r *sql.Rows) error {
			var id int64
			var name string
			var date time.Time
			if err := r.Scan(&id, &name, &date); err != nil {
				return err
			}
			e := get(date, name)
			e.Name, e.EventID = name, &id
			return nil
		}); err != nil {
		return Stats{}, err
	}

	// Shared sessions and the openings of their links.
	shootEvent := map[int64]eventKey{}
	if err := s.eachRow(ctx, `SELECT id, event_name, event_date FROM photobooth_shoots WHERE status = 'finalized'`, nil,
		func(r *sql.Rows) error {
			var id int64
			var name string
			var date time.Time
			if err := r.Scan(&id, &name, &date); err != nil {
				return err
			}
			get(date, name).Sessions++
			shootEvent[id] = eventKey{date.Format("2006-01-02"), NameKey(name)}
			return nil
		}); err != nil {
		return Stats{}, err
	}
	allScanners := map[string]bool{}
	scanned := map[int64]bool{}
	if err := s.eachRow(ctx, `SELECT target_id, visitor, views FROM public_visits WHERE kind = $1`, []any{VisitShare},
		func(r *sql.Rows) error {
			var shoot int64
			var visitor string
			var views int
			if err := r.Scan(&shoot, &visitor, &views); err != nil {
				return err
			}
			k, ok := shootEvent[shoot]
			if !ok {
				return nil
			}
			e := byKey[k]
			e.ScanViews += views
			if !scanned[shoot] {
				scanned[shoot] = true
				e.ScannedSessions++
			}
			scanners[k][visitor] = true
			allScanners[visitor] = true
			return nil
		}); err != nil {
		return Stats{}, err
	}

	// Each evening's public album: the gallery's album of the Drive folder the booth filed its pictures in.
	if err := s.eachRow(ctx, `SELECT ef.event_date, ef.event_key, a.id_photo_albums, a.title,
			COALESCE(v.visitors, 0), COALESCE(v.views, 0)
		FROM photobooth_event_folders ef
		JOIN photo_albums a ON a.root_folder_id = ef.event_folder_id AND a.deleted_at IS NULL
		LEFT JOIN (SELECT target_id, count(*) AS visitors, sum(views) AS views FROM public_visits
			WHERE kind = $1 GROUP BY target_id) v ON v.target_id = a.id_photo_albums`, []any{VisitAlbum},
		func(r *sql.Rows) error {
			var date time.Time
			var key, title string
			var id int64
			var visitors, views int
			if err := r.Scan(&date, &key, &id, &title, &visitors, &views); err != nil {
				return err
			}
			e := byKey[eventKey{date.Format("2006-01-02"), key}]
			if e == nil || e.AlbumID != nil {
				return nil
			}
			e.AlbumID, e.AlbumTitle, e.AlbumVisitors, e.AlbumViews = &id, title, visitors, views
			return nil
		}); err != nil {
		return Stats{}, err
	}

	out := Stats{Events: make([]EventStats, 0, len(byKey))}
	albumVisitors := map[int64]bool{}
	for k, e := range byKey {
		e.Scanners = len(scanners[k])
		out.Events = append(out.Events, *e)
		t := &out.Totals
		t.Events++
		t.Sessions += e.Sessions
		t.ScannedSessions += e.ScannedSessions
		t.ScanViews += e.ScanViews
		t.AlbumViews += e.AlbumViews
		if e.AlbumID != nil {
			albumVisitors[*e.AlbumID] = true
		}
	}
	slices.SortFunc(out.Events, func(a, b EventStats) int {
		return cmp.Or(cmp.Compare(b.Date, a.Date), cmp.Compare(a.Name, b.Name))
	})

	t := &out.Totals
	t.Scanners = len(allScanners)
	if len(albumVisitors) > 0 {
		ids := make([]int64, 0, len(albumVisitors))
		for id := range albumVisitors {
			ids = append(ids, id)
		}
		if err := s.db.QueryRowContext(ctx, `SELECT count(DISTINCT visitor) FROM public_visits
			WHERE kind = $1 AND target_id = ANY($2)`, VisitAlbum, pq.Array(ids)).Scan(&t.AlbumVisitors); err != nil {
			return Stats{}, err
		}
	}
	var withSessions, scannersSum, albumSum int
	for _, e := range out.Events {
		if e.Sessions > 0 {
			withSessions++
			scannersSum += e.Scanners
			albumSum += e.AlbumVisitors
		}
	}
	if withSessions > 0 {
		t.AvgScannersPerEvent = float64(scannersSum) / float64(withSessions)
		t.AvgAlbumVisitorsPerEvent = float64(albumSum) / float64(withSessions)
	}
	if t.Sessions > 0 {
		t.ScanRate = float64(t.ScannedSessions) / float64(t.Sessions)
		t.AvgScannersPerSession = float64(scannersSum) / float64(t.Sessions)
	}

	var since sql.NullTime
	if err := s.db.QueryRowContext(ctx, `SELECT min(first_seen) FROM public_visits`).Scan(&since); err != nil {
		return Stats{}, err
	}
	if since.Valid {
		out.TrackingSince = &since.Time
	}
	return out, nil
}

func (s *Service) eachRow(ctx context.Context, query string, args []any, fn func(*sql.Rows) error) error {
	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		if err := fn(rows); err != nil {
			return err
		}
	}
	return rows.Err()
}
