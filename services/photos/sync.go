package photos

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/lib/pq"
	"google.golang.org/api/googleapi"
)

// advisoryLockKey serialises syncs across every backend instance, without needing Redis.
const advisoryLockKey int64 = 7_303_001

var ErrSyncRunning = errors.New("photos sync already running")

type Mode string

const (
	ModeFull        Mode = "full"
	ModeIncremental Mode = "incremental"
)

type SyncStats struct {
	Mode       Mode             `json:"mode"`
	Resumed    bool             `json:"resumed,omitempty"`
	Changes    int              `json:"changes"`
	Upserted   int              `json:"upserted"`
	Deleted    int64            `json:"deleted"`
	Rebuilds   int              `json:"structure_rebuilds"`
	Folders    int              `json:"folders"`
	Albums     int              `json:"albums"`
	DriveCalls map[string]int64 `json:"drive_calls"`
	Duration   string           `json:"duration"`
}

type syncState struct {
	ChangesToken    sql.NullString
	PendingToken    sql.NullString
	FullRunID       sql.NullInt64
	FullPhase       sql.NullString
	FullPageToken   sql.NullString
	LastFull        sql.NullTime
	LastIncremental sql.NullTime
	DataVersion     int64
}

// SyncService mirrors the shared Drive into PostgreSQL. It only ever reads from Drive.
type SyncService struct {
	db             *sql.DB
	api            DriveAPI
	driveID        string
	reconcileEvery time.Duration
	now            func() time.Time
	// Optional: announces new albums after each tick. Nil when album notifications are off.
	notifier *AlbumNotifier
}

// NewSyncService: reconcileEvery is how often a full resync corrects drift the changes feed cannot express.
func NewSyncService(db *sql.DB, api DriveAPI, driveID string, reconcileEvery time.Duration) *SyncService {
	return &SyncService{db: db, api: api, driveID: driveID, reconcileEvery: reconcileEvery, now: time.Now}
}

func (s *SyncService) loadState(ctx context.Context) (syncState, error) {
	var st syncState
	err := s.db.QueryRowContext(ctx, `
		SELECT changes_page_token, pending_changes_token, full_run_id, full_phase, full_page_token,
		       last_full_at, last_incremental_at, data_version
		FROM photo_sync_state WHERE id = 1`).
		Scan(&st.ChangesToken, &st.PendingToken, &st.FullRunID, &st.FullPhase, &st.FullPageToken,
			&st.LastFull, &st.LastIncremental, &st.DataVersion)
	return st, err
}

// Tick runs whatever is due: resume an interrupted full sync, a first or reconciling full sync, else an incremental one.
func (s *SyncService) Tick(ctx context.Context) (SyncStats, error) {
	st, err := s.loadState(ctx)
	if err != nil {
		return SyncStats{}, err
	}
	full := !st.ChangesToken.Valid || st.FullRunID.Valid ||
		!st.LastFull.Valid || s.now().Sub(st.LastFull.Time) > s.reconcileEvery
	if full {
		return s.Run(ctx, ModeFull)
	}
	return s.Run(ctx, ModeIncremental)
}

// Run executes one sync of the given mode while holding the cluster-wide advisory lock.
func (s *SyncService) Run(ctx context.Context, mode Mode) (stats SyncStats, err error) {
	conn, err := s.db.Conn(ctx)
	if err != nil {
		return stats, err
	}
	defer conn.Close()
	var locked bool
	if err = conn.QueryRowContext(ctx, `SELECT pg_try_advisory_lock($1)`, advisoryLockKey).Scan(&locked); err != nil {
		return stats, err
	}
	if !locked {
		return stats, ErrSyncRunning
	}
	defer func() {
		_, _ = conn.ExecContext(context.Background(), `SELECT pg_advisory_unlock($1)`, advisoryLockKey)
	}()

	start := s.now()
	before := s.api.Calls()
	stats.Mode = mode
	var runID int64
	switch mode {
	case ModeFull:
		runID, err = s.runFull(ctx, &stats)
	default:
		runID, err = s.runIncremental(ctx, &stats)
	}
	stats.Duration = s.now().Sub(start).Round(time.Millisecond).String()
	stats.DriveCalls = diffCalls(before, s.api.Calls())
	if runID != 0 {
		s.finishRun(runID, stats, err)
	}
	if err != nil {
		log.Printf("📷 Photos %s sync failed: %v", mode, err)
		return stats, err
	}
	if stats.Changes > 0 || mode == ModeFull {
		log.Printf("📷 Photos %s sync done: %+v", mode, stats)
	}
	return stats, nil
}

func diffCalls(before, after map[string]int64) map[string]int64 {
	out := map[string]int64{}
	for k, v := range after {
		if d := v - before[k]; d > 0 {
			out[k] = d
		}
	}
	return out
}

func (s *SyncService) startRun(ctx context.Context, mode Mode) (int64, error) {
	var id int64
	err := s.db.QueryRowContext(ctx, `INSERT INTO photo_sync_runs (mode) VALUES ($1) RETURNING id`, string(mode)).Scan(&id)
	return id, err
}

func (s *SyncService) finishRun(id int64, stats SyncStats, runErr error) {
	b, _ := json.Marshal(stats)
	msg := sql.NullString{}
	if runErr != nil {
		msg = sql.NullString{String: runErr.Error(), Valid: true}
	}
	_, _ = s.db.ExecContext(context.Background(),
		`UPDATE photo_sync_runs SET finished_at = now(), complete = $2, stats = $3, error = $4 WHERE id = $1`,
		id, runErr == nil, string(b), msg)
}

// ---------------------------------------------------------------- full (initial / reconciliation)

// runFull is resumable: structure first (one transaction), then files page by page, each page committed together
// with its checkpoint. A crash or Drive error resumes at the last committed page on the next Tick.
func (s *SyncService) runFull(ctx context.Context, stats *SyncStats) (int64, error) {
	st, err := s.loadState(ctx)
	if err != nil {
		return 0, err
	}
	var runID int64
	resumed := st.FullRunID.Valid && st.FullPhase.String == "files"
	if resumed {
		runID = st.FullRunID.Int64
		stats.Resumed = true
		_, _ = s.db.ExecContext(ctx, `UPDATE photo_sync_runs SET finished_at = NULL, error = NULL WHERE id = $1`, runID)
	} else {
		if st.FullRunID.Valid {
			_, _ = s.db.ExecContext(ctx, `UPDATE photo_sync_runs SET finished_at = now(), error = 'abandoned before the files phase' WHERE id = $1`, st.FullRunID.Int64)
		}
		if runID, err = s.startRun(ctx, ModeFull); err != nil {
			return 0, err
		}
		// Taken before listing so that changes made while we list are replayed by the first incremental run.
		token, err := s.api.StartPageToken(ctx)
		if err != nil {
			return runID, fmt.Errorf("start page token: %w", err)
		}
		if _, err = s.db.ExecContext(ctx, `
			UPDATE photo_sync_state SET full_run_id = $1, full_phase = 'structure', full_page_token = NULL,
			       pending_changes_token = $2 WHERE id = 1`, runID, token); err != nil {
			return runID, err
		}
		snap, err := s.planStructure(ctx)
		if err != nil {
			return runID, fmt.Errorf("structure: %w", err)
		}
		tx, err := s.db.BeginTx(ctx, nil)
		if err != nil {
			return runID, err
		}
		defer func() { _ = tx.Rollback() }()
		if _, err = s.applyStructure(ctx, tx, snap, runID, stats); err != nil {
			return runID, fmt.Errorf("structure: %w", err)
		}
		if _, err = tx.ExecContext(ctx, `UPDATE photo_sync_state SET full_phase = 'files' WHERE id = 1`); err != nil {
			return runID, err
		}
		if err = tx.Commit(); err != nil {
			return runID, err
		}
	}

	idx, err := loadFolderIndex(ctx, s.db)
	if err != nil {
		return runID, err
	}
	token := st.FullPageToken.String
	if !resumed {
		token = ""
	}
	for {
		items, next, err := s.api.ListFilesPage(ctx, token)
		if err != nil {
			if isBadToken(err) && token != "" {
				// The checkpoint expired: start over rather than loop on it.
				_, _ = s.db.ExecContext(ctx, `UPDATE photo_sync_state SET full_run_id = NULL, full_phase = NULL, full_page_token = NULL WHERE id = 1`)
			}
			return runID, fmt.Errorf("list files: %w", err)
		}
		inputs := make([]photoInput, 0, len(items))
		for _, it := range items {
			inputs = append(inputs, photoInput{Item: it, FolderID: s.folderOf(it)})
		}
		tx, err := s.db.BeginTx(ctx, nil)
		if err != nil {
			return runID, err
		}
		if err = upsertPhotos(ctx, tx, buildRows(inputs, idx, s.now()), runID); err == nil {
			_, err = tx.ExecContext(ctx, `UPDATE photo_sync_state SET full_page_token = $1 WHERE id = 1`, nullIfEmpty(next))
		}
		if err != nil {
			_ = tx.Rollback()
			return runID, err
		}
		if err = tx.Commit(); err != nil {
			return runID, err
		}
		stats.Upserted += len(items)
		if next == "" {
			break
		}
		token = next
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return runID, err
	}
	defer func() { _ = tx.Rollback() }()
	if err = sweepPhotos(ctx, tx, runID); err != nil {
		return runID, err
	}
	if err = refreshDerived(ctx, tx); err != nil {
		return runID, err
	}
	if _, err = tx.ExecContext(ctx, `
		UPDATE photo_sync_state SET changes_page_token = pending_changes_token, pending_changes_token = NULL,
		       last_full_at = now(), last_incremental_at = now(),
		       full_run_id = NULL, full_phase = NULL, full_page_token = NULL WHERE id = 1`); err != nil {
		return runID, err
	}
	if err = bumpVersion(ctx, tx); err != nil {
		return runID, err
	}
	stats.Changes = stats.Upserted
	return runID, tx.Commit()
}

// folderOf returns the indexed folder a file lives in ("" for the drive root).
func (s *SyncService) folderOf(it DriveItem) string {
	if it.ParentID == s.driveID {
		return ""
	}
	return it.ParentID
}

// ---------------------------------------------------------------- structure

func (s *SyncService) planStructure(ctx context.Context) (*Snapshot, error) {
	folders, err := s.api.ListFolders(ctx)
	if err != nil {
		return nil, err
	}
	return BuildSnapshot(ctx, folders, s.driveID, s.api)
}

// applyStructure writes albums and folders, re-points photos at their album and reports whether any folder
// changed album or event date, or came back from the trash (photos then need the next reconciliation to converge).
func (s *SyncService) applyStructure(ctx context.Context, tx *sql.Tx, snap *Snapshot, run int64, stats *SyncStats) (datesChanged bool, err error) {
	before, err := loadFolderIndex(ctx, tx)
	if err != nil {
		return false, err
	}
	// A folder restored from the trash brings its files back, but the feed does not necessarily list them again:
	// reconcile with a full sync instead of leaving them hidden.
	ids := make([]string, 0, len(snap.Folders))
	for _, f := range snap.Folders {
		ids = append(ids, f.ID)
	}
	var revived int
	if err = tx.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM photo_folders WHERE deleted_at IS NOT NULL AND drive_folder_id = ANY($1)`,
		pq.Array(ids)).Scan(&revived); err != nil {
		return false, err
	}
	albumIDs, err := upsertAlbums(ctx, tx, snap, run)
	if err != nil {
		return false, err
	}
	if err = upsertFolders(ctx, tx, snap, albumIDs, run); err != nil {
		return false, err
	}
	if err = sweepStructure(ctx, tx, run); err != nil {
		return false, err
	}
	if err = syncAlbumMembership(ctx, tx); err != nil {
		return false, err
	}
	after, err := loadFolderIndex(ctx, tx)
	if err != nil {
		return false, err
	}
	inputs := make([]photoInput, 0, len(snap.Media))
	for _, m := range snap.Media {
		inputs = append(inputs, photoInput{Item: m.DriveItem, FolderID: m.FolderID})
	}
	if err = upsertPhotos(ctx, tx, buildRows(inputs, after, s.now()), run); err != nil {
		return false, err
	}
	stats.Folders, stats.Albums = len(snap.Folders), len(albumIDs)
	stats.Rebuilds++
	return revived > 0 || !before.sameAs(after), nil
}

// ---------------------------------------------------------------- incremental

func (s *SyncService) runIncremental(ctx context.Context, stats *SyncStats) (runID int64, err error) {
	st, err := s.loadState(ctx)
	if err != nil {
		return 0, err
	}
	if !st.ChangesToken.Valid {
		return 0, errors.New("no changes checkpoint: run a full sync first")
	}
	token := st.ChangesToken.String
	ensureRun := func() error {
		if runID != 0 {
			return nil
		}
		runID, err = s.startRun(ctx, ModeIncremental)
		return err
	}

	for {
		page, err := s.api.Changes(ctx, token)
		if err != nil {
			if isBadToken(err) {
				// An expired checkpoint cannot be replayed: force a full resync on the next Tick.
				_, _ = s.db.ExecContext(ctx, `UPDATE photo_sync_state SET changes_page_token = NULL WHERE id = 1`)
			}
			return runID, fmt.Errorf("changes: %w", err)
		}
		next := page.NextPageToken
		if next == "" {
			next = page.NewStartPageToken
		}
		changes := lastChangePerFile(page.Changes)

		if len(changes) == 0 {
			if _, err = s.db.ExecContext(ctx, `UPDATE photo_sync_state SET changes_page_token = $1, last_incremental_at = now() WHERE id = 1`, next); err != nil {
				return runID, err
			}
		} else {
			if err = ensureRun(); err != nil {
				return runID, err
			}
			if err = s.applyChanges(ctx, changes, runID, next, stats); err != nil {
				return runID, err
			}
		}
		if page.NextPageToken == "" {
			return runID, nil
		}
		token = page.NextPageToken
	}
}

func lastChangePerFile(in []Change) []Change {
	pos := map[string]int{}
	out := make([]Change, 0, len(in))
	for _, c := range in {
		if i, ok := pos[c.FileID]; ok {
			out[i] = c
			continue
		}
		pos[c.FileID] = len(out)
		out = append(out, c)
	}
	return out
}

func (s *SyncService) applyChanges(ctx context.Context, changes []Change, run int64, nextToken string, stats *SyncStats) error {
	stats.Changes += len(changes)

	structural, err := s.touchesStructure(ctx, changes)
	if err != nil {
		return err
	}
	var snap *Snapshot
	if structural {
		if snap, err = s.planStructure(ctx); err != nil {
			return fmt.Errorf("structure: %w", err)
		}
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	datesChanged := false
	if structural {
		if datesChanged, err = s.applyStructure(ctx, tx, snap, run, stats); err != nil {
			return err
		}
	}
	idx, err := loadFolderIndex(ctx, tx)
	if err != nil {
		return err
	}

	var inputs []photoInput
	var gone []string
	for _, c := range changes {
		switch {
		case c.Removed || c.Item == nil || c.Item.Trashed:
			gone = append(gone, c.FileID)
		case c.Item.IsFolder() || c.Item.IsShortcut():
			// Handled by the structure rebuild.
		default:
			inputs = append(inputs, photoInput{Item: *c.Item, FolderID: s.folderOf(*c.Item)})
		}
	}
	deleted, err := softDeletePhotos(ctx, tx, gone)
	if err != nil {
		return err
	}
	if err = upsertPhotos(ctx, tx, buildRows(inputs, idx, s.now()), run); err != nil {
		return err
	}
	if err = refreshDerived(ctx, tx); err != nil {
		return err
	}
	stats.Upserted += len(inputs)
	stats.Deleted += deleted

	if _, err = tx.ExecContext(ctx, `UPDATE photo_sync_state SET changes_page_token = $1, last_incremental_at = now() WHERE id = 1`, nextToken); err != nil {
		return err
	}
	if datesChanged {
		// Sort times derive from event dates: converge them with a full resync on the next Tick.
		if _, err = tx.ExecContext(ctx, `UPDATE photo_sync_state SET last_full_at = NULL WHERE id = 1`); err != nil {
			return err
		}
	}
	if err = bumpVersion(ctx, tx); err != nil {
		return err
	}
	return tx.Commit()
}

// touchesStructure tells whether a batch creates, renames, moves or removes a folder or shortcut.
func (s *SyncService) touchesStructure(ctx context.Context, changes []Change) (bool, error) {
	var removed []string
	for _, c := range changes {
		if c.Item != nil && !c.Item.Trashed && (c.Item.IsFolder() || c.Item.IsShortcut()) {
			return true, nil
		}
		if c.Removed || c.Item == nil || c.Item.Trashed {
			removed = append(removed, c.FileID)
		}
	}
	if len(removed) == 0 {
		return false, nil
	}
	var known bool
	err := s.db.QueryRowContext(ctx, `
		SELECT EXISTS (SELECT 1 FROM photo_folders WHERE drive_folder_id = ANY($1) OR shortcut_id = ANY($1))`,
		pq.Array(removed)).Scan(&known)
	return known, err
}

func isBadToken(err error) bool {
	var gerr *googleapi.Error
	return errors.As(err, &gerr) && (gerr.Code == 400 || gerr.Code == 404)
}

func nullIfEmpty(s string) any {
	if s == "" {
		return nil
	}
	return s
}
