package statistics

import (
	"github.com/gofiber/fiber/v2"
	"github.com/plugimt/transat-backend/utils"
)

// GetTopUserStatistics returns statistics about the top 10 users by request count
func (h *StatisticsHandler) GetTopUserStatistics(c *fiber.Ctx) error {
	utils.LogHeader("📊 Get Top Users Statistics")

	// Get top user statistics from service
	stats, err := h.statisticsService.GetTopUserStatistics()
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to get top user statistics")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve top user statistics",
		})
	}

	utils.LogMessage(utils.LevelInfo, "Successfully retrieved top users statistics")
	utils.LogLineKeyValue(utils.LevelInfo, "Users Count", len(stats))
	utils.LogFooter()

	// Return json response
	return c.JSON(fiber.Map{
		"success":    true,
		"statistics": stats,
	})
}
