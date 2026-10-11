package photos

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"os"
	"time"

	"github.com/plugimt/transat-backend/utils"
	"github.com/robfig/cron/v3"
)

const (
	defaultSyncInterval      = 10 * time.Minute
	defaultReconcileInterval = 7 * 24 * time.Hour
)

// FromEnv builds the sync service from GOOGLE_DRIVE_CREDENTIALS_JSON and PHOTOS_DRIVE_ID.
// It returns nil, nil when Drive is not configured so the rest of the API keeps working.
func FromEnv(ctx context.Context, db *sql.DB) (*SyncService, error) {
	creds, driveID := os.Getenv("GOOGLE_DRIVE_CREDENTIALS_JSON"), os.Getenv("PHOTOS_DRIVE_ID")
	if creds == "" && driveID == "" {
		return nil, nil
	}
	client, err := NewDriveClient(ctx, creds, driveID)
	if err != nil {
		return nil, err
	}
	return NewSyncService(db, client, driveID, envDuration("PHOTOS_RECONCILE_INTERVAL", defaultReconcileInterval)), nil
}

// Schedule polls the Drive changes feed every PHOTOS_SYNC_INTERVAL and runs one tick right away.
// Overlapping ticks (also across instances) are harmless: the advisory lock turns them into no-ops.
func (s *SyncService) Schedule(c *cron.Cron) {
	every := envDuration("PHOTOS_SYNC_INTERVAL", defaultSyncInterval)
	job := func() {
		if _, err := s.Tick(context.Background()); err != nil && !errors.Is(err, ErrSyncRunning) {
			utils.LogMessage(utils.LevelError, "Photos sync tick failed")
			utils.LogLineKeyValue(utils.LevelError, "Error", err)
		}
	}
	if _, err := c.AddFunc("@every "+every.String(), job); err != nil {
		log.Printf("Photos sync scheduling failed: %v", err)
		return
	}
	go job()
}

func envDuration(name string, def time.Duration) time.Duration {
	raw := os.Getenv(name)
	if raw == "" {
		return def
	}
	d, err := time.ParseDuration(raw)
	if err != nil || d < time.Minute {
		log.Printf("Warning: %s=%q is not a duration of at least 1m, using %s", name, raw, def)
		return def
	}
	return d
}
