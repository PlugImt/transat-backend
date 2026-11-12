package faces

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/lib/pq"
)

// advisoryLockKey serialises the worker across backend instances. The lock is keyed by schema too, so that
// isolated test schemas sharing one database do not block each other.
const advisoryLockKey int32 = 7_303_052

const (
	lockSQL   = `SELECT pg_try_advisory_lock($1, hashtext(current_schema()))`
	unlockSQL = `SELECT pg_advisory_unlock($1, hashtext(current_schema()))`
)

// maxScanAttempts is how many times a photo the model or Drive failed on is retried, an hour apart.
const maxScanAttempts = 3

var ErrWorkerRunning = errors.New("face worker already running")

// TickStats summarises one worker run.
type TickStats struct {
	Purged  int64
	Scanned int
	Faces   int
	Failed  int
	Matches int
}

// Tick runs the whole pipeline once, under a cluster-wide advisory lock:
//  1. forget the faces of photos that left the gallery (deleted, duplicates, now a cover);
//  2. match the faces already known against new references, so a new user's album fills in first;
//  3. scan photos not yet processed, newest first, within the tick's time budget, matching after each batch.
func (s *Service) Tick(ctx context.Context) (TickStats, error) {
	var stats TickStats
	if !s.Available() {
		return stats, ErrUnavailable
	}
	conn, err := s.db.Conn(ctx)
	if err != nil {
		return stats, err
	}
	defer conn.Close()
	var locked bool
	if err := conn.QueryRowContext(ctx, lockSQL, advisoryLockKey).Scan(&locked); err != nil {
		return stats, err
	}
	if !locked {
		return stats, ErrWorkerRunning
	}
	defer func() {
		_, _ = conn.ExecContext(context.Background(), unlockSQL, advisoryLockKey)
	}()

	if s.cfg.TickBudget > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, s.cfg.TickBudget)
		defer cancel()
	}

	if stats.Purged, err = s.purge(ctx); err != nil {
		return stats, fmt.Errorf("purge: %w", err)
	}
	if err := s.match(ctx, &stats); err != nil {
		return stats, fmt.Errorf("match: %w", err)
	}
	if s.images == nil {
		return stats, nil
	}
	for ctx.Err() == nil {
		n, err := s.scanBatch(ctx, &stats)
		if err != nil {
			return stats, fmt.Errorf("scan: %w", err)
		}
		if err := s.match(ctx, &stats); err != nil {
			return stats, fmt.Errorf("match: %w", err)
		}
		if n < s.cfg.ScanBatch {
			break
		}
	}
	return stats, nil
}

// purge deletes the faces (and so the matches) of photos that are no longer shown, and the scans of deleted photos.
func (s *Service) purge(ctx context.Context) (int64, error) {
	res, err := s.db.ExecContext(ctx, `
		DELETE FROM photo_faces f USING photos p
		WHERE p.drive_file_id = f.drive_file_id
		  AND (p.deleted_at IS NOT NULL OR p.duplicate_of IS NOT NULL OR p.media_type <> 'image' OR p.is_cover)`)
	if err != nil {
		return 0, err
	}
	n, _ := res.RowsAffected()
	_, err = s.db.ExecContext(ctx, `
		DELETE FROM face_photo_scans sc USING photos p
		WHERE p.drive_file_id = sc.drive_file_id AND p.deleted_at IS NOT NULL`)
	return n, err
}

// scanBatch processes up to ScanBatch photos: never scanned, scanned with another model or on other file
// content, or failed less than maxScanAttempts times and not in the last hour.
func (s *Service) scanBatch(ctx context.Context, stats *TickStats) (int, error) {
	model := s.ml.Model()
	rows, err := s.db.QueryContext(ctx, `
		SELECT p.drive_file_id, p.md5, COALESCE(sc.attempts, 0), sc.status = 'failed' AND sc.model = $1 AND sc.md5 = p.md5
		FROM photos p
		LEFT JOIN face_photo_scans sc ON sc.drive_file_id = p.drive_file_id
		WHERE `+eligiblePhotos+`
		  AND (sc.drive_file_id IS NULL OR sc.model <> $1 OR sc.md5 <> p.md5
		       OR (sc.status = 'failed' AND sc.attempts < $3 AND sc.scanned_at < now() - interval '1 hour'))
		ORDER BY sc.drive_file_id IS NULL DESC, p.sort_at DESC, p.drive_file_id
		LIMIT $2`, model, s.cfg.ScanBatch, maxScanAttempts)
	if err != nil {
		return 0, err
	}
	type job struct {
		id, md5  string
		attempts int
	}
	var jobs []job
	for rows.Next() {
		var j job
		var retry sql.NullBool
		if err := rows.Scan(&j.id, &j.md5, &j.attempts, &retry); err != nil {
			rows.Close()
			return 0, err
		}
		if !retry.Bool {
			j.attempts = 0 // a new model or new content starts over
		}
		jobs = append(jobs, j)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return 0, err
	}

	var mu sync.Mutex
	var firstErr error
	sem := make(chan struct{}, max(1, s.cfg.ScanConcurrency))
	var wg sync.WaitGroup
	for _, j := range jobs {
		if ctx.Err() != nil {
			break
		}
		sem <- struct{}{}
		wg.Add(1)
		go func(j job) {
			defer wg.Done()
			defer func() { <-sem }()
			faces, err := s.scanOne(ctx, j.id, j.md5, j.attempts)
			mu.Lock()
			defer mu.Unlock()
			switch {
			case err != nil && ctx.Err() != nil:
				// The budget ran out: the photo is picked up again next time.
			case err != nil:
				if firstErr == nil {
					firstErr = err
				}
			case faces < 0:
				stats.Failed++
			default:
				stats.Scanned++
				stats.Faces += faces
			}
		}(j)
	}
	wg.Wait()
	return len(jobs), firstErr
}

// scanOne detects the faces of one photo and replaces what was stored for it. It returns -1 when the photo
// failed (recorded for a later retry); an error only when the database itself fails.
func (s *Service) scanOne(ctx context.Context, id, md5 string, attempts int) (int, error) {
	model := s.ml.Model()
	img, err := s.images.Image(ctx, id, s.cfg.ScanSize)
	var det Detection
	if err == nil {
		det, err = s.ml.Detect(ctx, img)
	}
	if err != nil {
		if ctx.Err() != nil {
			return 0, ctx.Err()
		}
		status := "failed"
		if errors.Is(err, ErrImageUnavailable) || errors.Is(err, ErrUnreadableImage) {
			status = "skipped"
		}
		log.Printf("Faces: photo %s %s: %v", id, status, err)
		_, dbErr := s.db.ExecContext(ctx, `
			INSERT INTO face_photo_scans (drive_file_id, model, md5, status, attempts, face_count, last_error, scanned_at)
			VALUES ($1, $2, $3, $4, $5, 0, $6, now())
			ON CONFLICT (drive_file_id) DO UPDATE SET model = $2, md5 = $3, status = $4, attempts = $5, face_count = 0,
				last_error = $6, scanned_at = now()`, id, model, md5, status, attempts+1, truncate(err.Error(), 500))
		return -1, dbErr
	}

	kept := det.Faces[:0]
	for _, f := range det.Faces {
		if f.Score < s.cfg.MinFaceScore {
			continue
		}
		if f.Width()*float32(det.ImageWidth) < float32(s.cfg.MinFacePixels) ||
			f.Height()*float32(det.ImageHeight) < float32(s.cfg.MinFacePixels) {
			continue
		}
		kept = append(kept, f)
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer func() { _ = tx.Rollback() }()
	// The photo may have been removed from the index meanwhile: then there is nothing to store.
	var exists bool
	if err := tx.QueryRowContext(ctx, `SELECT EXISTS (SELECT 1 FROM photos WHERE drive_file_id = $1 FOR SHARE)`, id).Scan(&exists); err != nil {
		return 0, err
	}
	if !exists {
		return 0, nil
	}
	if _, err := tx.ExecContext(ctx, `DELETE FROM photo_faces WHERE drive_file_id = $1`, id); err != nil {
		return 0, err
	}
	for _, f := range kept {
		if _, err := tx.ExecContext(ctx, `
			INSERT INTO photo_faces (drive_file_id, model, box_x1, box_y1, box_x2, box_y2, score, embedding)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`, id, model, f.X1, f.Y1, f.X2, f.Y2, f.Score, encode(f.Embedding)); err != nil {
			return 0, err
		}
	}
	if _, err := tx.ExecContext(ctx, `
		INSERT INTO face_photo_scans (drive_file_id, model, md5, status, attempts, face_count, last_error, scanned_at)
		VALUES ($1, $2, $3, 'done', $4, $5, NULL, now())
		ON CONFLICT (drive_file_id) DO UPDATE SET model = $2, md5 = $3, status = 'done', attempts = $4, face_count = $5,
			last_error = NULL, scanned_at = now()`, id, model, md5, attempts+1, len(kept)); err != nil {
		return 0, err
	}
	return len(kept), tx.Commit()
}

func truncate(s string, n int) string {
	if len(s) > n {
		return s[:n]
	}
	return s
}

type reference struct {
	id      int64
	email   string
	through int64
	emb     Embedding
}

// matchChunk bounds the faces held in memory while matching (2 KiB each).
const matchChunk = 5000

// match compares every face newer than a reference's cursor with that reference and records the close ones.
// Faces are only inserted by the worker, under the same lock, so the id cursor never skips a face.
func (s *Service) match(ctx context.Context, stats *TickStats) error {
	model := s.ml.Model()
	var maxID int64
	if err := s.db.QueryRowContext(ctx, `SELECT COALESCE(max(id), 0) FROM photo_faces`).Scan(&maxID); err != nil {
		return err
	}
	rows, err := s.db.QueryContext(ctx, `
		SELECT r.id, r.email, r.matched_through, r.embedding FROM face_references r
		JOIN face_consents c ON c.email = r.email
		WHERE r.model = $1 AND r.matched_through < $2 ORDER BY r.matched_through`, model, maxID)
	if err != nil {
		return err
	}
	var refs []reference
	for rows.Next() {
		var r reference
		var raw []byte
		if err := rows.Scan(&r.id, &r.email, &r.through, &raw); err != nil {
			rows.Close()
			return err
		}
		if r.emb, err = decode(raw); err != nil {
			rows.Close()
			return err
		}
		refs = append(refs, r)
	}
	rows.Close()
	if err := rows.Err(); err != nil || len(refs) == 0 {
		return err
	}

	cursor := refs[0].through // the smallest, refs are sorted
	for cursor < maxID {
		if err := ctx.Err(); err != nil {
			return nil // progress so far is saved; the rest waits for the next tick
		}
		next, n, err := s.matchFrom(ctx, refs, cursor, maxID, model)
		if err != nil {
			return err
		}
		stats.Matches += n
		cursor = next
		for i := range refs {
			refs[i].through = max(refs[i].through, cursor)
		}
	}
	return nil
}

// matchFrom handles one chunk of faces after cursor and advances every reference past it, in one transaction.
func (s *Service) matchFrom(ctx context.Context, refs []reference, cursor, maxID int64, model string) (int64, int, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, drive_file_id, embedding FROM photo_faces
		WHERE id > $1 AND id <= $2 AND model = $3 ORDER BY id LIMIT $4`, cursor, maxID, model, matchChunk)
	if err != nil {
		return 0, 0, err
	}
	type hit struct {
		email, file string
		face        int64
		dist        float32
	}
	var hits []hit
	last := maxID // an empty chunk (other models only) means everything up to maxID is done
	count := 0
	for rows.Next() {
		var id int64
		var file string
		var raw []byte
		if err := rows.Scan(&id, &file, &raw); err != nil {
			rows.Close()
			return 0, 0, err
		}
		count++
		last = id
		emb, err := decode(raw)
		if err != nil {
			continue
		}
		for _, r := range refs {
			if id <= r.through {
				continue
			}
			if d := distance(r.emb, emb); d <= s.cfg.MaxDistance {
				hits = append(hits, hit{r.email, file, id, d})
			}
		}
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return 0, 0, err
	}
	if count < matchChunk {
		last = maxID
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, 0, err
	}
	defer func() { _ = tx.Rollback() }()
	added := 0
	for _, h := range hits {
		// Through face_consents: a user who withdrew meanwhile gets nothing back.
		res, err := tx.ExecContext(ctx, `
			INSERT INTO face_matches (email, face_id, drive_file_id, distance)
			SELECT c.email, $2, $3, $4 FROM face_consents c WHERE c.email = $1
			ON CONFLICT (email, face_id) DO UPDATE SET distance = LEAST(face_matches.distance, EXCLUDED.distance)`,
			h.email, h.face, h.file, h.dist)
		if err != nil {
			return 0, 0, err
		}
		if n, _ := res.RowsAffected(); n > 0 {
			added++
		}
	}
	ids := make([]int64, len(refs))
	for i, r := range refs {
		ids[i] = r.id
	}
	if _, err := tx.ExecContext(ctx, `
		UPDATE face_references SET matched_through = $2 WHERE id = ANY($1) AND matched_through < $2`,
		pq.Array(ids), last); err != nil {
		return 0, 0, err
	}
	return last, added, tx.Commit()
}
