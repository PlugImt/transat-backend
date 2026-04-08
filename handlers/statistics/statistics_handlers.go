package statistics

import (
	"database/sql"

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

func (h *StatisticsHandler) GetDashboardStatistics(c *fiber.Ctx) error {
	utils.LogHeader("📊 Get Dashboard Statistics")

	type DashboardStats struct {
		TotalUsers      int                      `json:"totalUsers"`
		UnverifiedUsers int                      `json:"unverifiedUsers"`
		TotalEvents     int                      `json:"totalEvents"`
		TotalClubs      int                      `json:"totalClubs"`
		UserGrowth      []map[string]interface{} `json:"userGrowth"`
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

	utils.LogMessage(utils.LevelInfo, "Successfully fetched dashboard statistics")
	utils.LogFooter()

	return c.JSON(stats)
}
