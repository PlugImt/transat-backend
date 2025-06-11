package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/plugimt/transat-backend/handlers"
	"github.com/plugimt/transat-backend/services"
)

func SetupNaolibRoutes(router fiber.Router, naolibService *services.NaolibService) {
	// Groupe de routes pour Naolib
	naolib := router.Group("/naolib")
	handler := handlers.NewNaolibHandler(naolibService)

	// Route pour obtenir les prochains départs
	naolib.Get("/departures/chantrerie", handler.GetNextDeparturesChantrerie)
}
