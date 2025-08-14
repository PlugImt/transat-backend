package routes

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"github.com/plugimt/transat-backend/handlers/bassine" // Import the reservation handlers
	"github.com/plugimt/transat-backend/middlewares"
	"github.com/plugimt/transat-backend/utils"
)

func SetupBassineRoutes(router fiber.Router, db *sql.DB) {
	// Initialize Reservation Handler
	bassineHandler := bassine.NewBassineHandler(db)

	bassineGroup := router.Group("/bassine", middlewares.JWTMiddleware, utils.EnhanceSentryEventWithEmail)

	// Root reservation routes
	bassineGroup.Patch("", bassineHandler.IncrementBassine)                // Increments or decrements the bassine count
	bassineGroup.Patch("/", bassineHandler.IncrementBassine)               // Increments or decrements the bassine count
	bassineGroup.Get("", bassineHandler.GetMyBassine)                      // Returns the user's bassine count, rank, and surrounding users
	bassineGroup.Get("/", bassineHandler.GetMyBassine)                     // Returns the user's bassine count, rank, and surrounding users
	bassineGroup.Get("/leaderboard", bassineHandler.GetBassineLeaderboard) // Returns the leaderboard of users
	bassineGroup.Get("/history", bassineHandler.GetBassineGlobalHistory)   // Returns the history of bassine consumption for the logged-in user
	bassineGroup.Get("/history/", bassineHandler.GetBassineGlobalHistory)  // Returns the history of bassine consumption for the logged-in user
	bassineGroup.Get("/history/:email", bassineHandler.GetBassineHistory)  // Returns the history of bassine consumption for a specific user
}
