package routes

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"github.com/plugimt/transat-backend/handlers/game"
	"github.com/plugimt/transat-backend/middlewares"
)

// SetupGameRoutes configures routes for game features
func SetupGameRoutes(router fiber.Router, db *sql.DB) {
    bassineHandler := game.NewBassineHandler(db)

    gameGroup := router.Group("/game", middlewares.JWTMiddleware)

    gameGroup.Get("/bassine", bassineHandler.GetAllScores)
    gameGroup.Get("/bassine/:email", bassineHandler.GetUserScore)
    gameGroup.Put("/bassine/:email", bassineHandler.SetUserScore)
}


