package routes

import (
	"database/sql"
	"log"

	"github.com/plugimt/transat-backend/handlers/files" // Import the files handlers
	"github.com/plugimt/transat-backend/middlewares"
	"github.com/plugimt/transat-backend/services"

	"github.com/gofiber/fiber/v2"
)

// SetupFileRoutes configures routes for file operations.
func SetupFileRoutes(router fiber.Router, db *sql.DB, r2Service *services.R2Service) {
	// Initialize File Handler
	fileHandler, err := files.NewFileHandler(db, r2Service)
	if err != nil {
		// Log the error and potentially panic or exit if file handling is critical
		log.Fatalf("💥 Failed to initialize File Handler: %v", err)
	}

	// Public route to serve files (no auth needed usually)
	// Matches /api/data/some_unique_filename.ext
	router.Get("/data/:filename", fileHandler.ServeFile)

	// Authenticated routes for managing files
	// Placed directly under /api, protected by JWT middleware
	router.Post("/upload", middlewares.JWTMiddleware, fileHandler.UploadFile)
	router.Get("/files", middlewares.JWTMiddleware, fileHandler.ListUserFiles)           // List user's own files
	router.Delete("/files/:filename", middlewares.JWTMiddleware, fileHandler.DeleteFile) // Delete user's own file (by original name)

	// Admin route (example - secure appropriately!)
	// Requires admin check within handler or a specific admin middleware
	router.Get("/all-files", middlewares.JWTMiddleware, fileHandler.ListAllFiles) // List all files on server
}
