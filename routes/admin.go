package routes

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"github.com/plugimt/transat-backend/handlers/admin"
	"github.com/plugimt/transat-backend/middlewares"
	"github.com/plugimt/transat-backend/services"
)

// SetupAdminRoutes configures routes for admin functionality
func SetupAdminRoutes(router fiber.Router, db *sql.DB, statisticsService *services.StatisticsService) {
	// Initialize Admin Handler
	adminHandler := admin.NewAdminHandler(db, statisticsService)

	// Create admin group with restrictTo middleware protection
	adminGroup := router.Group("/admin")
	adminGroup.Use(middlewares.RestrictTo(db, "ADMIN"))

	// Dashboard endpoints (restricted to ADMIN role)
	adminGroup.Get("/dashboard/stats", adminHandler.GetDashboardStats)

	// User management endpoints (restricted to ADMIN role)
	adminGroup.Get("/users", adminHandler.GetUsers)

	// Profile endpoint accessible by any authenticated admin
	adminGroup.Get("/profile", adminHandler.GetCurrentUser)

	// Example of routes with multiple allowed roles:
	// Create a group that allows both ADMIN and MODERATOR roles
	moderatedGroup := router.Group("/moderated")
	moderatedGroup.Use(middlewares.RestrictTo(db, "ADMIN", "MODERATOR"))

	// These endpoints would be accessible by users with either ADMIN or MODERATOR role
	// moderatedGroup.Get("/some-endpoint", someHandler.SomeMethod)

	// Example of a route that allows multiple roles including regular users
	userGroup := router.Group("/user-admin")
	userGroup.Use(middlewares.RestrictTo(db, "ADMIN", "MODERATOR", "NEWF"))
	// userGroup.Get("/accessible-endpoint", someHandler.SomeMethod)
}
