package bus_departures

import (
	"github.com/gofiber/fiber/v2"
	"github.com/plugimt/transat-backend/services"
	"github.com/plugimt/transat-backend/utils"
)

type BusDeparturesHandler struct {
	gtfsService *services.GTFSService
}

func NewBusDeparturesHandler(gtfsService *services.GTFSService) *BusDeparturesHandler {
	return &BusDeparturesHandler{gtfsService: gtfsService}
}

func (h *BusDeparturesHandler) GetChantrerieDepartures(c *fiber.Ctx) error {
	if !h.gtfsService.IsReady() {
		c.Set(fiber.HeaderCacheControl, "no-store")
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"error": "Bus departure data is not available yet",
		})
	}

	departures, err := h.gtfsService.GetChantrerieDepartures()
	if err != nil {
		utils.LogMessage(utils.LevelError, "Error getting Chantrerie bus departures")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		c.Set(fiber.HeaderCacheControl, "no-store")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve bus departures",
		})
	}

	c.Set(fiber.HeaderCacheControl, "public, max-age=20")
	return c.JSON(departures)
}
