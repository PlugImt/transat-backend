package photos

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"
)

// AlbumNotice is what an announcement needs: nothing about the users, only the album.
type AlbumNotice struct {
	AlbumID   int64
	Title     string
	EventDate *time.Time
}

// Delivery is the outcome of one announcement, in counts.
type Delivery struct {
	Users   int
	Devices int
}

// NotifyFunc announces a new album. It returns an error only when nothing could be delivered, which is the
// signal to try again later; a partial delivery is final, so devices that got it are never notified twice.
type NotifyFunc func(ctx context.Context, notice AlbumNotice) (Delivery, error)

const (
	defaultNotifySettle = 30 * time.Minute
	notifyLease         = 15 * time.Minute
	notifyMaxAttempts   = 3
	notifyRetryBase     = 10 * time.Minute
)

// available is when an album shows up in the app: published and holding at least one photo.
const available = `a.visibility = 'published' AND a.deleted_at IS NULL AND a.photo_count > 0`

// AlbumNotifier announces new event albums once, using the albums the sync already indexed: it never talks to
// Google Drive. Its state lives in PostgreSQL (photo_album_notifications), so a restart, a retry or a second
// instance can neither lose nor repeat an announcement.
type AlbumNotifier struct {
	db     *sql.DB
	notify NotifyFunc
	// An album is announced once its photo count has stopped changing for this long, so an event still being
	// uploaded is not announced half-empty.
	settle time.Duration
	now    func() time.Time
}

func NewAlbumNotifier(db *sql.DB, notify NotifyFunc, settle time.Duration) *AlbumNotifier {
	if settle <= 0 {
		settle = defaultNotifySettle
	}
	return &AlbumNotifier{db: db, notify: notify, settle: settle, now: time.Now}
}

type NotifyStats struct {
	// Baselined is how many albums were recorded without announcing them (first run only).
	Baselined int
	Sent      int
	Retrying  int
	Failed    int
}

// Run observes the indexed albums and announces the new ones that are ready. Call it after every sync tick.
func (n *AlbumNotifier) Run(ctx context.Context) (NotifyStats, error) {
	now := n.now()
	var stats NotifyStats

	ready, baselined, err := n.ensureBaseline(ctx, now)
	if err != nil {
		return stats, err
	}
	stats.Baselined = baselined
	// Until the first complete sync has recorded what already exists, nothing may be announced.
	if !ready {
		return stats, nil
	}

	if err := n.observe(ctx, now); err != nil {
		return stats, err
	}
	claimed, err := n.claim(ctx, now)
	if err != nil {
		return stats, err
	}

	for _, album := range claimed {
		outcome, err := n.announce(ctx, album)
		switch {
		case err == nil:
			stats.Sent++
		case outcome == "failed":
			stats.Failed++
		default:
			stats.Retrying++
		}
	}
	return stats, nil
}

// ensureBaseline records every album that exists when notifications start as already announced. It happens once,
// after the first complete sync, so deploying the feature or indexing the Drive for the first time never sends a
// notification per existing album.
func (n *AlbumNotifier) ensureBaseline(ctx context.Context, now time.Time) (ready bool, count int, err error) {
	tx, err := n.db.BeginTx(ctx, nil)
	if err != nil {
		return false, 0, err
	}
	defer func() { _ = tx.Rollback() }()

	var done sql.NullTime
	if err = tx.QueryRowContext(ctx, `SELECT baselined_at FROM photo_notify_state WHERE id = 1 FOR UPDATE`).Scan(&done); err != nil {
		return false, 0, fmt.Errorf("read notification state: %w", err)
	}
	if done.Valid {
		return true, 0, nil
	}

	var synced bool
	if err = tx.QueryRowContext(ctx, `
		SELECT EXISTS (SELECT 1 FROM photo_sync_state WHERE id = 1 AND last_full_at IS NOT NULL AND full_run_id IS NULL)`).
		Scan(&synced); err != nil {
		return false, 0, err
	}
	if !synced {
		return false, 0, nil
	}

	res, err := tx.ExecContext(ctx, `
		INSERT INTO photo_album_notifications (root_folder_id, album_id, status, last_count, first_seen_at, count_changed_at)
		SELECT a.root_folder_id, a.id_photo_albums, 'baseline', a.photo_count, $1, $1
		FROM photo_albums a WHERE `+available+`
		ON CONFLICT (root_folder_id) DO NOTHING`, now)
	if err != nil {
		return false, 0, fmt.Errorf("record existing albums: %w", err)
	}
	inserted, _ := res.RowsAffected()
	if _, err = tx.ExecContext(ctx, `UPDATE photo_notify_state SET baselined_at = $1 WHERE id = 1`, now); err != nil {
		return false, 0, err
	}
	if err = tx.Commit(); err != nil {
		return false, 0, err
	}
	log.Printf("Photos: %d existing albums recorded without notification", inserted)
	return true, int(inserted), nil
}

// observe registers albums seen for the first time and restarts the settle delay of those still growing.
// Only top-level albums count: events nested inside another event are not announced on their own.
func (n *AlbumNotifier) observe(ctx context.Context, now time.Time) error {
	if _, err := n.db.ExecContext(ctx, `
		INSERT INTO photo_album_notifications (root_folder_id, album_id, status, last_count, first_seen_at, count_changed_at)
		SELECT a.root_folder_id, a.id_photo_albums, 'pending', a.photo_count, $1, $1
		FROM photo_albums a WHERE `+available+` AND a.parent_album_id IS NULL
		ON CONFLICT (root_folder_id) DO NOTHING`, now); err != nil {
		return fmt.Errorf("register new albums: %w", err)
	}
	if _, err := n.db.ExecContext(ctx, `
		UPDATE photo_album_notifications n
		SET last_count = a.photo_count, count_changed_at = $1, album_id = a.id_photo_albums
		FROM photo_albums a
		WHERE a.root_folder_id = n.root_folder_id AND n.status = 'pending' AND n.last_count <> a.photo_count`, now); err != nil {
		return fmt.Errorf("track album growth: %w", err)
	}
	// A claim whose owner vanished (crash, deploy) is retried; after the last attempt it is given up.
	if _, err := n.db.ExecContext(ctx, `
		UPDATE photo_album_notifications
		SET status = 'failed', last_error = 'claim expired after the last attempt'
		WHERE status = 'sending' AND claimed_at <= $1 AND attempts >= $2`, now.Add(-notifyLease), notifyMaxAttempts); err != nil {
		return fmt.Errorf("close expired claims: %w", err)
	}
	return nil
}

type claimedAlbum struct {
	rootFolderID string
	notice       AlbumNotice
}

// claim atomically takes the albums that are ready. Only one process can take a given album: the status change
// is a single UPDATE, so instances racing for it see it already claimed.
func (n *AlbumNotifier) claim(ctx context.Context, now time.Time) ([]claimedAlbum, error) {
	rows, err := n.db.QueryContext(ctx, `
		UPDATE photo_album_notifications n
		SET status = 'sending', claimed_at = $1, attempts = n.attempts + 1
		FROM photo_albums a
		WHERE a.root_folder_id = n.root_folder_id AND `+available+` AND a.parent_album_id IS NULL
		  AND n.attempts < $4
		  AND ((n.status = 'pending' AND n.count_changed_at <= $2 AND (n.retry_at IS NULL OR n.retry_at <= $1))
		       OR (n.status = 'sending' AND n.claimed_at <= $3))
		RETURNING n.root_folder_id, a.id_photo_albums, a.title, a.event_date`,
		now, now.Add(-n.settle), now.Add(-notifyLease), notifyMaxAttempts)
	if err != nil {
		return nil, fmt.Errorf("claim albums: %w", err)
	}
	defer rows.Close()

	var out []claimedAlbum
	for rows.Next() {
		var c claimedAlbum
		var date sql.NullTime
		if err := rows.Scan(&c.rootFolderID, &c.notice.AlbumID, &c.notice.Title, &date); err != nil {
			return nil, err
		}
		if date.Valid {
			d := date.Time
			c.notice.EventDate = &d
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

// announce sends one claimed album and records the result. A failure never escapes to the sync.
func (n *AlbumNotifier) announce(ctx context.Context, album claimedAlbum) (outcome string, err error) {
	var delivery Delivery
	func() {
		defer func() {
			if r := recover(); r != nil {
				err = fmt.Errorf("panic while notifying: %v", r)
			}
		}()
		delivery, err = n.notify(ctx, album.notice)
	}()

	now := n.now()
	if err == nil {
		_, dbErr := n.db.ExecContext(ctx, `
			UPDATE photo_album_notifications
			SET status = 'sent', notified_at = $2, devices = $3, last_error = NULL, claimed_at = NULL
			WHERE root_folder_id = $1 AND status = 'sending'`, album.rootFolderID, now, delivery.Devices)
		if dbErr != nil {
			log.Printf("Photos: album %d announced to %d devices but not recorded: %v", album.notice.AlbumID, delivery.Devices, dbErr)
		} else {
			log.Printf("Photos: album %d announced to %d devices (%d users)", album.notice.AlbumID, delivery.Devices, delivery.Users)
		}
		return "sent", nil
	}

	message := err.Error()
	if len(message) > 300 {
		message = message[:300]
	}
	var attempts int
	row := n.db.QueryRowContext(ctx, `
		UPDATE photo_album_notifications
		SET status = CASE WHEN attempts >= $3 THEN 'failed' ELSE 'pending' END,
		    retry_at = $2::timestamptz + (attempts * $4::bigint) * interval '1 microsecond',
		    last_error = $5, claimed_at = NULL
		WHERE root_folder_id = $1 AND status = 'sending'
		RETURNING attempts`,
		album.rootFolderID, now, notifyMaxAttempts, notifyRetryBase.Microseconds(), message)
	if scanErr := row.Scan(&attempts); scanErr != nil && scanErr != sql.ErrNoRows {
		log.Printf("Photos: could not record the failure of album %d: %v", album.notice.AlbumID, scanErr)
	}
	if attempts >= notifyMaxAttempts {
		log.Printf("Photos: album %d notification abandoned after %d attempts: %s", album.notice.AlbumID, attempts, message)
		return "failed", err
	}
	log.Printf("Photos: album %d notification failed (attempt %d), will retry: %s", album.notice.AlbumID, attempts, message)
	return "retry", err
}
