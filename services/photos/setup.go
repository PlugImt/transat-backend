package photos

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"os"
	"strings"
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

// EnableAlbumNotifications announces new albums through notify after every tick, unless
// PHOTOS_ALBUM_NOTIFICATIONS is set to "off" (the switch that stops announcements without touching the sync).
func (s *SyncService) EnableAlbumNotifications(notify NotifyFunc) {
	if v := strings.ToLower(os.Getenv("PHOTOS_ALBUM_NOTIFICATIONS")); v == "off" || v == "false" || v == "0" {
		log.Println("Info: Photos album notifications disabled by PHOTOS_ALBUM_NOTIFICATIONS")
		return
	}
	s.notifier = NewAlbumNotifier(s.db, notify, envDuration("PHOTOS_NOTIFY_SETTLE", defaultNotifySettle))
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
		// Whatever the tick did, pending announcements still get their turn; a failure here never reaches the sync.
		s.runNotifier()
	}
	if _, err := c.AddFunc("@every "+every.String(), job); err != nil {
		log.Printf("Photos sync scheduling failed: %v", err)
		return
	}
	go job()
}

func (s *SyncService) runNotifier() {
	if s.notifier == nil {
		return
	}
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Photos album notifier panicked: %v", r)
		}
	}()
	if _, err := s.notifier.Run(context.Background()); err != nil {
		utils.LogMessage(utils.LevelError, "Photos album notifier failed")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
	}
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
