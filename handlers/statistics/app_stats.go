package statistics

import (
	"context"
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"github.com/plugimt/transat-backend/services/faces"
	"github.com/plugimt/transat-backend/utils"
)

// activeWindow is how recent a device's last heartbeat must be to count as in use.
const activeWindow = "30 days"

type PlatformCount struct {
	Platform string `json:"platform"`
	Devices  int    `json:"devices"`
	Users    int    `json:"users"`
}

type VersionCount struct {
	Platform string `json:"platform"`
	Version  string `json:"version"` // "" when the app did not report one
	Devices  int    `json:"devices"`
	Users    int    `json:"users"`
}

type DeviceStats struct {
	ActiveWindowDays int             `json:"active_window_days"`
	Devices          int             `json:"devices"`
	Users            int             `json:"users"`
	Platforms        []PlatformCount `json:"platforms"`
	Versions         []VersionCount  `json:"versions"`
}

// AppStatsHandler serves GET /statistics/app: platforms and app versions in use, and face recognition usage.
type AppStatsHandler struct {
	db    *sql.DB
	faces *faces.Service
}

func NewAppStatsHandler(db *sql.DB, f *faces.Service) *AppStatsHandler {
	return &AppStatsHandler{db: db, faces: f}
}

func (h *AppStatsHandler) Get(c *fiber.Ctx) error {
	ctx := c.UserContext()
	devices, err := deviceStats(ctx, h.db)
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to compute device statistics")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to compute statistics"})
	}
	faceStats, err := h.faces.AdminStats(ctx)
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to compute face recognition statistics")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to compute statistics"})
	}
	c.Set(fiber.HeaderCacheControl, "no-store")
	return c.JSON(fiber.Map{"devices": devices, "faces": faceStats})
}

// deviceStats counts the installations seen in the last 30 days (active and signed in), by platform and version.
func deviceStats(ctx context.Context, db *sql.DB) (DeviceStats, error) {
	st := DeviceStats{ActiveWindowDays: 30, Platforms: []PlatformCount{}, Versions: []VersionCount{}}
	const recent = `status = 'active' AND last_seen_at > now() - interval '` + activeWindow + `'`

	if err := db.QueryRowContext(ctx, `SELECT count(*), count(DISTINCT user_email) FROM devices WHERE `+recent).
		Scan(&st.Devices, &st.Users); err != nil {
		return st, err
	}
	rows, err := db.QueryContext(ctx, `
		SELECT platform, count(*), count(DISTINCT user_email) FROM devices WHERE `+recent+`
		GROUP BY platform ORDER BY count(*) DESC`)
	if err != nil {
		return st, err
	}
	for rows.Next() {
		var p PlatformCount
		if err := rows.Scan(&p.Platform, &p.Devices, &p.Users); err != nil {
			rows.Close()
			return st, err
		}
		st.Platforms = append(st.Platforms, p)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return st, err
	}

	rows, err = db.QueryContext(ctx, `
		SELECT platform, COALESCE(NULLIF(trim(app_version), ''), ''), count(*), count(DISTINCT user_email)
		FROM devices WHERE `+recent+`
		GROUP BY 1, 2 ORDER BY app_version_key(NULLIF(COALESCE(NULLIF(trim(app_version), ''), ''), '')) DESC NULLS LAST, 1`)
	if err != nil {
		return st, err
	}
	defer rows.Close()
	for rows.Next() {
		var v VersionCount
		if err := rows.Scan(&v.Platform, &v.Version, &v.Devices, &v.Users); err != nil {
			return st, err
		}
		st.Versions = append(st.Versions, v)
	}
	return st, rows.Err()
}
