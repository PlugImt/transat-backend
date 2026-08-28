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
func SetupAuthRoutes(router fiber.Router, db *sql.DB, jwtSecret []byte, notifService *services.NotificationService, emailService *services.EmailService, discordService *services.DiscordService) {
	// Initialize Auth Handler with dependencies
	// Note: jwtSecret is now passed to the handler constructor
	authHandler := auth.NewAuthHandler(db, jwtSecret, notifService, emailService, discordService)

	// Group routes for auth operations and apply multiple security layers
	authGroup := router.Group("/auth")

	authGroup.Post("/register", middlewares.LoginRegisterLimiter, authHandler.Register)
	authGroup.Post("/login", middlewares.LoginRegisterLimiter, authHandler.Login)
	authGroup.Post("/verify-account", middlewares.AccessRateLimiter, authHandler.VerifyAccount)
	
	// Apply slightly less restrictive rate limiter to these endpoints
	authGroup.Post("/verification-code", middlewares.AccountRateLimiter, authHandler.RequestVerificationCode)
	authGroup.Patch("/change-password", middlewares.AccountRateLimiter, authHandler.ChangePassword)
}
