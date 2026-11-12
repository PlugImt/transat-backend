package faces

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

const defaultModel = "buffalo_l"

// FromEnv builds the service from:
//   - FACES_ML_URL: the machine-learning container on the private network (e.g. http://face-ml.railway.internal:3003).
//     Unset, the feature is off (the app hides it) and nothing is processed;
//   - FACES_SCOPE_FOLDERS: comma-separated Drive folder IDs; only photos under them are scanned (all when unset);
//   - FACES_MODEL (buffalo_l), FACES_MAX_DISTANCE (0.5), FACES_MIN_SCORE (0.7), FACES_SCAN_CONCURRENCY (2);
//   - GOOGLE_DRIVE_CREDENTIALS_JSON, the photo sync's read-only account, to read the photos.
//
// It never returns nil: an unavailable service still answers status requests.
func FromEnv(ctx context.Context, db *sql.DB) *Service {
	cfg := DefaultConfig()
	cfg.MaxDistance = envFloat("FACES_MAX_DISTANCE", cfg.MaxDistance)
	cfg.MinFaceScore = envFloat("FACES_MIN_SCORE", cfg.MinFaceScore)
	if n, err := strconv.Atoi(os.Getenv("FACES_SCAN_CONCURRENCY")); err == nil && n > 0 && n <= 16 {
		cfg.ScanConcurrency = n
	}

	for _, id := range strings.Split(os.Getenv("FACES_SCOPE_FOLDERS"), ",") {
		if id = strings.TrimSpace(id); id != "" {
			cfg.ScopeFolders = append(cfg.ScopeFolders, id)
		}
	}
	if len(cfg.ScopeFolders) > 0 {
		log.Printf("Info: ℹ️ Face recognition limited to Drive folders %v", cfg.ScopeFolders)
	}

	mlURL := os.Getenv("FACES_ML_URL")
	if mlURL == "" {
		log.Println("Info: ℹ️ Face recognition disabled: FACES_ML_URL is not set")
		return NewService(db, nil, nil, cfg)
	}
	model := os.Getenv("FACES_MODEL")
	if model == "" {
		model = defaultModel
	}
	ml := NewMLClient(mlURL, model, cfg.MinFaceScore)

	var images ImageSource
	if creds := os.Getenv("GOOGLE_DRIVE_CREDENTIALS_JSON"); creds != "" {
		d, err := NewDriveImages(ctx, creds)
		if err != nil {
			log.Printf("Warning: face recognition cannot read Drive, photos will not be scanned: %v", err)
		} else {
			images = d
		}
	} else {
		log.Println("Info: ℹ️ Face recognition will not scan photos: GOOGLE_DRIVE_CREDENTIALS_JSON is not set")
	}
	return NewService(db, ml, images, cfg)
}

// Start runs the worker every interval and whenever Wake is called, until ctx ends.
func (s *Service) Start(ctx context.Context, interval time.Duration) {
	if !s.Available() {
		return
	}
	go func() {
		t := time.NewTicker(interval)
		defer t.Stop()
		for {
			s.runOnce(ctx)
			select {
			case <-ctx.Done():
				return
			case <-t.C:
			case <-s.wake:
			}
		}
	}()
}

func (s *Service) runOnce(ctx context.Context) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Faces worker panicked: %v", r)
		}
	}()
	stats, err := s.Tick(ctx)
	switch {
	case errors.Is(err, ErrWorkerRunning):
	case err != nil:
		log.Printf("Faces worker failed: %v (%+v)", err, stats)
	case stats.Scanned+stats.Failed+stats.Matches > 0 || stats.Purged > 0:
		log.Printf("🙂 Faces worker: %+v", stats)
	}
}

func envFloat(name string, def float32) float32 {
	raw := os.Getenv(name)
	if raw == "" {
		return def
	}
	v, err := strconv.ParseFloat(raw, 32)
	if err != nil || v <= 0 || v >= 2 {
		log.Printf("Warning: %s=%q is invalid, using %v", name, raw, def)
		return def
	}
	return float32(v)
}
