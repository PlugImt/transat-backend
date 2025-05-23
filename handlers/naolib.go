package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/plugimt/transat-backend/services"
	"github.com/plugimt/transat-backend/utils"
)

type NaolibHandler struct {
	service *services.NaolibService
}

func NewNaolibHandler(service *services.NaolibService) *NaolibHandler {
	return &NaolibHandler{
		service: service,
	}
}

func (h *NaolibHandler) GetNextDeparturesChantrerie(c *fiber.Ctx) error {
	departuresC6, err := h.service.GetNextDepartures("CTRE2")
	if err != nil {
		utils.LogMessage(utils.LevelError, "Error fetching departures: "+err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch departure information",
		})
	}
	departures75, err := h.service.GetNextDepartures("CTRE4")
	if err != nil {
		utils.LogMessage(utils.LevelError, "Error fetching departures: "+err.Error())
	}

	departures := append(departuresC6, departures75...)
	return c.JSON(fiber.Map{
		"success": true,
		"data":    departures,
	})
}
