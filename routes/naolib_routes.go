package routes

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"github.com/plugimt/transat-backend/handlers"
	"github.com/plugimt/transat-backend/services"
)

func SetupNaolibRoutes(router fiber.Router, naolibService *services.NaolibService, db *sql.DB) {
	// Groupe de routes pour Naolib
	naolib := router.Group("/naolib")
	handler := handlers.NewNaolibHandler(naolibService, db)

	// Route pour obtenir les prochains départs
	naolib.Get("/departures/chantrerie", handler.GetNextDeparturesChantrerie)

	// TODO: protéger ces routes !
	naolib.Post("/import/netex/stops", handler.ImportNetexStops)
	naolib.Post("/import/netex/offer", handler.ImportNetexOffer)

	naolib.Get("/search", handler.SearchStopPlace)

	naolib.Get("/generate-request", handler.GenerateNetexRequest)

	naolib.Get("/get-departures", handler.GetDepartures)
}
