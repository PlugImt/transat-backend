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
// @Summary		Obtenir les statistiques par endpoint
// @Description	Récupère les statistiques d'utilisation détaillées pour tous les endpoints de l'API
// @Tags			Statistics
// @Produce		json
// @Security		BearerAuth
// @Success		200	{object}	map[string]interface{}	"Statistiques par endpoint récupérées avec succès"
// @Failure		401	{object}	models.ErrorResponse	"Non autorisé"
// @Failure		500	{object}	models.ErrorResponse	"Erreur serveur"
// @Router			/statistics/endpoints [get]
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
// @Summary		Obtenir les statistiques globales
// @Description	Récupère les statistiques globales d'utilisation de l'API (tous endpoints confondus)
// @Tags			Statistics
// @Produce		json
// @Security		BearerAuth
// @Success		200	{object}	map[string]interface{}	"Statistiques globales récupérées avec succès"
// @Failure		401	{object}	models.ErrorResponse	"Non autorisé"
// @Failure		500	{object}	models.ErrorResponse	"Erreur serveur"
// @Router			/statistics/global [get]
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
