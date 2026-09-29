package statistics

import (
	"database/sql"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/plugimt/transat-backend/services"
	"github.com/plugimt/transat-backend/utils"
)

// StatisticsHandler handles statistics related operations
type StatisticsHandler struct {
	db                *sql.DB
	statisticsService *services.StatisticsService
}

// NewStatisticsHandler creates a new instance of StatisticsHandler
func NewStatisticsHandler(db *sql.DB, statisticsService *services.StatisticsService) *StatisticsHandler {
	return &StatisticsHandler{
		db:                db,
		statisticsService: statisticsService,
	}
}

// GetEndpointStatistics returns statistics for all endpoints
func (h *StatisticsHandler) GetEndpointStatistics(c *fiber.Ctx) error {
	utils.LogHeader("📊 Get Endpoint Statistics")

	stats, err := h.statisticsService.GetEndpointStatistics()
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to get endpoint statistics")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve endpoint statistics",
		})
	}

	utils.LogMessage(utils.LevelInfo, "Successfully retrieved endpoint statistics")
	utils.LogLineKeyValue(utils.LevelInfo, "Endpoint Count", len(stats))
	utils.LogFooter()

	return c.JSON(fiber.Map{
		"success":    true,
		"statistics": stats,
	})
}

// GetGlobalStatistics returns global statistics across all endpoints
func (h *StatisticsHandler) GetGlobalStatistics(c *fiber.Ctx) error {
	utils.LogHeader("📊 Get Global Statistics")

	stats, err := h.statisticsService.GetGlobalStatistics()
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to get global statistics")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve global statistics",
		})
	}

	utils.LogMessage(utils.LevelInfo, "Successfully retrieved global statistics")
	utils.LogLineKeyValue(utils.LevelInfo, "Total Requests", stats.TotalRequestCount)
	utils.LogLineKeyValue(utils.LevelInfo, "Avg Duration", stats.GlobalAvgDurationMs)
	utils.LogFooter()

	return c.JSON(fiber.Map{
		"success":    true,
		"statistics": stats,
	})
}

// GetActiveUsersOverTime returns the number of distinct active users bucketed by the
// requested granularity (day, week, month, quarter, year) so the dashboard chart can
// switch between views.
func (h *StatisticsHandler) GetActiveUsersOverTime(c *fiber.Ctx) error {
	utils.LogHeader("📊 Get Active Users Over Time")

	period := c.Query("period", "day")

	data, err := h.statisticsService.GetActiveUsersOverTime(period)
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to get active users over time")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid period, expected one of: day, week, month, quarter, year",
		})
	}

	utils.LogMessage(utils.LevelInfo, "Successfully retrieved active users over time")
	utils.LogLineKeyValue(utils.LevelInfo, "Period", period)
	utils.LogFooter()

	return c.JSON(fiber.Map{
		"success": true,
		"period":  period,
		"data":    data,
	})
}

// GetActivityHeatmap returns the number of distinct active users for each day-of-week/hour
// slot within the requested range (week, month, year, all), to visualize when users are
// typically active during the week.
func (h *StatisticsHandler) GetActivityHeatmap(c *fiber.Ctx) error {
	utils.LogHeader("📊 Get Activity Heatmap")

	rangeParam := c.Query("range", "all")

	data, err := h.statisticsService.GetActivityHeatmap(rangeParam)
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to get activity heatmap")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid range, expected one of: week, month, year, all",
		})
	}

	utils.LogMessage(utils.LevelInfo, "Successfully retrieved activity heatmap")
	utils.LogLineKeyValue(utils.LevelInfo, "Range", rangeParam)
	utils.LogFooter()

	return c.JSON(fiber.Map{
		"success": true,
		"range":   rangeParam,
		"data":    data,
	})
}

func (h *StatisticsHandler) GetDashboardStatistics(c *fiber.Ctx) error {
	utils.LogHeader("📊 Get Dashboard Statistics")

	type DashboardStats struct {
		TotalUsers      int                      `json:"totalUsers"`
		UnverifiedUsers int                      `json:"unverifiedUsers"`
		TotalEvents     int                      `json:"totalEvents"`
		TotalClubs      int                      `json:"totalClubs"`
		UserGrowth      []map[string]interface{} `json:"userGrowth"`
		ActiveUsers     struct {
			DAU int `json:"dau"`
			WAU int `json:"wau"`
			MAU int `json:"mau"`
			YAU int `json:"yau"`
		} `json:"activeUsers"`
		DailyActiveUsers []map[string]interface{} `json:"dailyActiveUsers"`
	}

	var stats DashboardStats

	userCountQuery := `SELECT COUNT(*) FROM newf`
	err := h.db.QueryRow(userCountQuery).Scan(&stats.TotalUsers)
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to get total users count")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
	}

	unverifiedQuery := `
		SELECT COUNT(DISTINCT nr.email)
		FROM newf_roles nr
		JOIN roles r ON nr.id_roles = r.id_roles
		WHERE r.name = 'VERIFYING'`
	err = h.db.QueryRow(unverifiedQuery).Scan(&stats.UnverifiedUsers)
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to get unverified users count")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
	}

	eventsQuery := `SELECT COUNT(*) FROM events`
	err = h.db.QueryRow(eventsQuery).Scan(&stats.TotalEvents)
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to get events count")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
	}

	clubsQuery := `SELECT COUNT(*) FROM clubs`
	err = h.db.QueryRow(clubsQuery).Scan(&stats.TotalClubs)
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to get clubs count")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
	}

	growthQuery := `
		WITH dates AS (
			SELECT generate_series(
				(SELECT MIN(DATE(creation_date)) FROM newf),
				(SELECT MAX(DATE(creation_date)) FROM newf),
				interval '1 day'
			)::date AS date
		)
		SELECT
			d.date,
			COUNT(n.*) as count,
			SUM(COUNT(n.*)) OVER (ORDER BY d.date ASC) as cumulativeCount
		FROM dates d
		LEFT JOIN newf n
			ON DATE(n.creation_date) = d.date
		GROUP BY d.date
		ORDER BY d.date ASC;
	`
	rows, err := h.db.Query(growthQuery)
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to get user growth data")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
	} else {
		defer rows.Close()
		for rows.Next() {
			var date string
			var count int
			var cumulativeCount int
			if err := rows.Scan(&date, &count, &cumulativeCount); err == nil {
				stats.UserGrowth = append(stats.UserGrowth, map[string]interface{}{
					"date":            date,
					"count":           count,
					"cumulativeCount": cumulativeCount,
				})
			}
		}
	}

	// Utilisateurs actifs (fenêtres glissantes), basés sur les requêtes authentifiées
	activeQuery := `
		SELECT
			COUNT(DISTINCT user_email) FILTER (WHERE request_received >= NOW() - INTERVAL '1 day'),
			COUNT(DISTINCT user_email) FILTER (WHERE request_received >= NOW() - INTERVAL '7 days'),
			COUNT(DISTINCT user_email) FILTER (WHERE request_received >= NOW() - INTERVAL '30 days'),
			COUNT(DISTINCT user_email)
		FROM request_statistics
		WHERE user_email IS NOT NULL
			AND request_received >= NOW() - INTERVAL '365 days'`
	err = h.db.QueryRow(activeQuery).Scan(
		&stats.ActiveUsers.DAU, &stats.ActiveUsers.WAU, &stats.ActiveUsers.MAU, &stats.ActiveUsers.YAU,
	)
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to get active users counts")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
	}

	// DAU par jour sur les 90 derniers jours (jours sans activité inclus)
	dauQuery := `
		SELECT d.date, COUNT(DISTINCT r.user_email)
		FROM generate_series(CURRENT_DATE - INTERVAL '89 days', CURRENT_DATE, INTERVAL '1 day') AS d(date)
		LEFT JOIN request_statistics r
			ON DATE(r.request_received) = d.date::date AND r.user_email IS NOT NULL
		GROUP BY d.date
		ORDER BY d.date ASC`
	dauRows, err := h.db.Query(dauQuery)
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to get daily active users")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
	} else {
		defer dauRows.Close()
		for dauRows.Next() {
			var date time.Time
			var count int
			if err := dauRows.Scan(&date, &count); err == nil {
				stats.DailyActiveUsers = append(stats.DailyActiveUsers, map[string]interface{}{
					"date":  date.Format("2006-01-02"),
					"count": count,
				})
			}
		}
	}

	utils.LogMessage(utils.LevelInfo, "Successfully fetched dashboard statistics")
	utils.LogFooter()

	return c.JSON(stats)
}
