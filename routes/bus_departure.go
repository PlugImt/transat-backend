package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/plugimt/transat-backend/handlers/bus_departures"
	"github.com/plugimt/transat-backend/services"
)

func SetupBusDepartureRoutes(router fiber.Router, gtfsService *services.GTFSService) {
	handler := bus_departures.NewBusDeparturesHandler(gtfsService)

	router.Get("/departures", handler.GetChantrerieDepartures)
}
