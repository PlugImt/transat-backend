package routes

import (
	"database/sql"

	"github.com/plugimt/transat-backend/handlers/auth" // Import the auth handlers
	"github.com/plugimt/transat-backend/middlewares"
	"github.com/plugimt/transat-backend/services" // For NotificationService

	// For EmailService
	"github.com/gofiber/fiber/v2"
)

// SetupAuthRoutes configures authentication related routes.
func SetupAuthRoutes(router fiber.Router, db *sql.DB, jwtSecret []byte, notifService *services.NotificationService, emailService *services.EmailService) {
	// Initialize Auth Handler with dependencies
	// Note: jwtSecret is now passed to the handler constructor
	authHandler := auth.NewAuthHandler(db, jwtSecret, notifService, emailService)

	// Group routes with rate limiting for login/register
	authGroup := router.Group("/auth", middlewares.LoginRegisterLimiter) // Apply limiter

	authGroup.Post("/register", authHandler.Register)
	authGroup.Post("/login", authHandler.Login)
	authGroup.Post("/verify-account", authHandler.VerifyAccount)
	authGroup.Post("/verification-code", authHandler.RequestVerificationCode)
	authGroup.Patch("/change-password", authHandler.ChangePassword)
}
