package routes

import (
	"database/sql"

	"github.com/plugimt/transat-backend/handlers/traq" // Import the traq handlers
	"github.com/plugimt/transat-backend/middlewares"   // Import if specific middlewares are needed

	"github.com/gofiber/fiber/v2"
)

// SetupTraqRoutes configures routes related to Traq articles and types.
func SetupTraqRoutes(router fiber.Router, db *sql.DB) {
	// Initialize Traq Handler
	traqHandler := traq.NewTraqHandler(db)

	// Group Traq routes under /traq
	// Add authentication middleware if needed for specific actions
	// Example: traqGroup := router.Group("/traq", middlewares.JWTMiddleware)
	traqGroup := router.Group("/traq")

	// Article routes
	// Consider adding middleware (e.g., JWT, admin check) to POST/PUT/DELETE
	traqGroup.Post("/", middlewares.JWTMiddleware, traqHandler.CreateTraqArticle) // Example: Require JWT
	traqGroup.Get("/", traqHandler.GetAllTraqArticles)                            // Publicly accessible?
	// Uncomment and implement handlers for specific article actions if needed
	// traqGroup.Get("/:id", traqHandler.GetTraqArticle)    // GET /api/traq/:id
	// traqGroup.Put("/:id", middlewares.JWTMiddleware, traqHandler.UpdateTraqArticle) // Example: Require JWT
	// traqGroup.Delete("/:id", middlewares.JWTMiddleware, traqHandler.DeleteTraqArticle) // Example: Require JWT

	// Type routes
	traqTypesGroup := traqGroup.Group("/types") // /api/traq/types

	// Assuming creating types might need auth, getting them might be public
	// Add middlewares as appropriate
	traqTypesGroup.Post("/", middlewares.JWTMiddleware, traqHandler.CreateTraqType) // Example: Require JWT
	traqTypesGroup.Get("/", traqHandler.GetAllTraqTypes)                            // Publicly accessible?
}
