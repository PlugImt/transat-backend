package routes

import (
	"database/sql"

	"github.com/plugimt/transat-backend/handlers/user" // Import the user handlers
	"github.com/plugimt/transat-backend/middlewares"
	"github.com/plugimt/transat-backend/services" // Import NotificationService
	"github.com/plugimt/transat-backend/utils"

	"github.com/gofiber/fiber/v2"
)

// SetupUserRoutes configures user profile and related routes.
func SetupUserRoutes(router fiber.Router, db *sql.DB, notifService *services.NotificationService) {
	// Initialize User Handler with dependencies
	userHandler := user.NewUserHandler(db, notifService)

	// Group routes that require JWT authentication
	// Changed group name from "/newf" to "/user" for clarity
	userGroup := router.Group("/newf", middlewares.JWTMiddleware, utils.EnhanceSentryEventWithEmail)

	// Profile routes
	userGroup.Get("/me", userHandler.GetNewf)       // GET /api/user/me
	userGroup.Patch("/me", userHandler.UpdateNewf)  // PATCH /api/user/me
	userGroup.Delete("/me", userHandler.DeleteNewf) // DELETE /api/user/me (Use with caution!)

	// Notification preferences: one idempotent read and write per category.
	userGroup.Get("/notifications/preferences", userHandler.GetNotificationPreferences) // GET /api/newf/notifications/preferences
	userGroup.Put("/notifications/preferences", userHandler.SetNotificationPreference)  // PUT /api/newf/notifications/preferences {"service","enabled"}
}
