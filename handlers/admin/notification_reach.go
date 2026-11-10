package admin

import (
	"context"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/lib/pq"
	"github.com/plugimt/transat-backend/utils"
)

// unreachablePeople names the users who cannot be notified: the panel shows people, not addresses.
func (h *AdminHandler) unreachablePeople(ctx context.Context, emails []string) []string {
	names := []string{}
	if len(emails) == 0 {
		return names
	}
	rows, err := h.DB.QueryContext(ctx, `
		SELECT COALESCE(NULLIF(trim(first_name || ' ' || last_name), ''), email)
		FROM newf WHERE email = ANY($1) ORDER BY last_name, first_name`, pq.Array(emails))
	if err != nil {
		return names
	}
	defer rows.Close()
	for rows.Next() {
		var n string
		if rows.Scan(&n) == nil {
			names = append(names, n)
		}
	}
	return names
}

// GetAppVersions handles GET /admin/notifications/app-versions: how many devices and people run each version of
// the app, newest first, to pick who to ask to update.
func (h *AdminHandler) GetAppVersions(c *fiber.Ctx) error {
	rows, err := h.DB.QueryContext(c.UserContext(), `
		SELECT COALESCE(NULLIF(app_version, ''), 'inconnue'), platform, count(*), count(DISTINCT user_email),
		       max(last_seen_at)::text
		FROM devices
		WHERE status = 'active'
		GROUP BY NULLIF(app_version, ''), platform
		ORDER BY app_version_key(NULLIF(app_version, '')) DESC, platform`)
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to load app versions")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to load app versions"})
	}
	defer rows.Close()
	type row struct {
		Version  string  `json:"version"`
		Platform string  `json:"platform"`
		Devices  int     `json:"devices"`
		Users    int     `json:"users"`
		LastSeen *string `json:"lastSeen"`
	}
	out := []row{}
	for rows.Next() {
		var r row
		var seen *string
		if err := rows.Scan(&r.Version, &r.Platform, &r.Devices, &r.Users, &seen); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to load app versions"})
		}
		r.Platform = strings.ToLower(r.Platform)
		r.LastSeen = seen
		out = append(out, r)
	}
	return c.JSON(fiber.Map{"versions": out})
}
