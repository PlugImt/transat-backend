package routes

import (
	"database/sql"
	"github.com/plugimt/transat-backend/services/photoreport"

	"github.com/gofiber/fiber/v2"
	"github.com/plugimt/transat-backend/handlers/statistics"
	"github.com/plugimt/transat-backend/middlewares"
	"github.com/plugimt/transat-backend/services"
	"github.com/plugimt/transat-backend/services/faces"
)

// SetupStatisticsRoutes configures routes for statistics
func SetupStatisticsRoutes(router fiber.Router, db *sql.DB, statisticsService *services.StatisticsService, faceService *faces.Service) {
	// Initialize Statistics Handler
	statsHandler := statistics.NewStatisticsHandler(db, statisticsService)

	// Toutes les routes de statistiques nécessitent un JWT admin
	// The dashboard is also open to the Club Photo's owners, who review photo reports.
	statsGroup := router.Group("/statistics", middlewares.JWTMiddleware, middlewares.AdminOrRolesMiddleware(db, photoreport.PhotoManagerRoles...))

	statsGroup.Get("/endpoints", statsHandler.GetEndpointStatistics)
	statsGroup.Get("/global", statsHandler.GetGlobalStatistics)
	statsGroup.Get("/top-users", statsHandler.GetTopUserStatistics)
	statsGroup.Get("/dashboard", statsHandler.GetDashboardStatistics)
	statsGroup.Get("/active-users", statsHandler.GetActiveUsersOverTime)
	statsGroup.Get("/activity-heatmap", statsHandler.GetActivityHeatmap)
	// Platforms and app versions in use, face recognition usage.
	statsGroup.Get("/app", statistics.NewAppStatsHandler(db, faceService).Get)
}
