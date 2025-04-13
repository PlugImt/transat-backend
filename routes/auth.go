package routes

import (
	"database/sql"

	"Transat_2.0_Backend/handlers/auth" // Import the auth handlers
	"Transat_2.0_Backend/middlewares"
	"Transat_2.0_Backend/services" // For NotificationService
	"Transat_2.0_Backend/utils"    // For EmailService

	"github.com/gofiber/fiber/v2"
)

// SetupAuthRoutes configures authentication related routes.
func SetupAuthRoutes(router fiber.Router, db *sql.DB, jwtSecret []byte, notifService *services.NotificationService, emailService *utils.EmailService) {
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
