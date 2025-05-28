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

	// Notification Preferences routes (within the authenticated user group)
	// Combine add/remove into one endpoint toggling state
	userGroup.Post("/notifications/subscriptions", userHandler.AddOrRemoveNotificationSubscription) // POST /api/user/notifications/subscriptions expects {"service": "..."}
	// Endpoint to get subscription status (all or specific)
	// Use GET with optional body
	userGroup.Get("/notifications/subscriptions", userHandler.GetNotificationSubscriptions) // GET /api/user/notifications/subscriptions (body optional for specific check)

	// Route to *trigger* sending a notification (might require admin)
	// This was previously POST /newf/send-notification
	// Kept under /user for now, but permissions must be checked in the handler.
	// Consider moving to /admin or /notifications group later.
	userGroup.Post("/send-notification", userHandler.SendNotification) // POST /api/user/send-notification (Requires Permissions!)
}
