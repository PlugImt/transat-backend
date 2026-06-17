package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/plugimt/transat-backend/handlers/departures"
	"github.com/plugimt/transat-backend/services"
)

func SetupBusDepartureRoutes(router fiber.Router, gtfsService *services.GTFSService) {
	handler := departures.NewDeparturesHandler(gtfsService)

	router.Get("/departures", handler.GetChantrerieDepartures)
}
