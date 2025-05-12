package routes

import (
	"database/sql"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/plugimt/transat-backend/handlers/support"
	"github.com/plugimt/transat-backend/middlewares"
	"github.com/plugimt/transat-backend/services"
)

// SetupSupportRoutes configures routes for support request operations
func SetupSupportRoutes(router fiber.Router, db *sql.DB, emailService *services.EmailService) {
	// Initialize Support Handler
	supportHandler := support.NewSupportHandler(db, emailService)
	if supportHandler == nil {
		log.Fatalf("💥 Failed to initialize Support Handler")
	}

	// Base route is "/api/support" - these routes are all protected by JWT
	router.Get("/support", middlewares.JWTMiddleware, supportHandler.GetSupportRequests)
	router.Post("/support", middlewares.JWTMiddleware, supportHandler.CreateSupportRequest)
	
	// Admin route for updating support requests
	router.Patch("/support/:id", middlewares.JWTMiddleware, supportHandler.UpdateSupportStatus)
} 