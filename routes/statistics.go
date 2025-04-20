package routes

import (
	"database/sql"

	"Transat_2.0_Backend/handlers/statistics"
	"Transat_2.0_Backend/services"
	"github.com/gofiber/fiber/v2"
)

// SetupStatisticsRoutes configures routes for statistics
func SetupStatisticsRoutes(router fiber.Router, db *sql.DB, statisticsService *services.StatisticsService) {
	// Initialize Statistics Handler
	statsHandler := statistics.NewStatisticsHandler(db, statisticsService)

	// Create a statistics group (public for easier testing)
	statsGroup := router.Group("/statistics")

	// Endpoints
	statsGroup.Get("/endpoints", statsHandler.GetEndpointStatistics)
	statsGroup.Get("/global", statsHandler.GetGlobalStatistics)
} 