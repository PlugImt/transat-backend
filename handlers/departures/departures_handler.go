package departures

import (
	"github.com/gofiber/fiber/v2"
	"github.com/plugimt/transat-backend/services"
	"github.com/plugimt/transat-backend/utils"
)

type DeparturesHandler struct {
	gtfsService *services.GTFSService
}

func NewDeparturesHandler(gtfsService *services.GTFSService) *DeparturesHandler {
	return &DeparturesHandler{gtfsService: gtfsService}
}

func (h *DeparturesHandler) GetChantrerieDepartures(c *fiber.Ctx) error {
	if !h.gtfsService.IsReady() {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"error": "Bus departure data is not available yet",
		})
	}

	departures, err := h.gtfsService.GetChantrerieDepartures()
	if err != nil {
		utils.LogMessage(utils.LevelError, "Error getting Chantrerie bus departures")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve bus departures",
		})
	}

	return c.JSON(departures)
}
