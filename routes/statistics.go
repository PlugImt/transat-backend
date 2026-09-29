package routes

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"github.com/plugimt/transat-backend/handlers/statistics"
	"github.com/plugimt/transat-backend/middlewares"
	"github.com/plugimt/transat-backend/services"
)

// SetupStatisticsRoutes configures routes for statistics
func SetupStatisticsRoutes(router fiber.Router, db *sql.DB, statisticsService *services.StatisticsService) {
	// Initialize Statistics Handler
	statsHandler := statistics.NewStatisticsHandler(db, statisticsService)

	// Toutes les routes de statistiques nécessitent un JWT admin
	statsGroup := router.Group("/statistics", middlewares.JWTMiddleware, middlewares.AdminAuthMiddleware(db))

	statsGroup.Get("/endpoints", statsHandler.GetEndpointStatistics)
	statsGroup.Get("/global", statsHandler.GetGlobalStatistics)
	statsGroup.Get("/top-users", statsHandler.GetTopUserStatistics)
	statsGroup.Get("/dashboard", statsHandler.GetDashboardStatistics)
}
