package handlers

import (
	"database/sql"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/plugimt/transat-backend/models"
	"github.com/plugimt/transat-backend/services"
	"github.com/plugimt/transat-backend/services/netex"
	"github.com/plugimt/transat-backend/utils"
)

type NaolibHandler struct {
	service *services.NaolibService
	db      *sql.DB
}

func NewNaolibHandler(service *services.NaolibService, db *sql.DB) *NaolibHandler {
	return &NaolibHandler{
		service: service,
		db:      db,
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

func (h *NaolibHandler) ImportNetexData(c *fiber.Ctx) error {
	var body struct {
		Url string `json:"url"`
	}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	url := body.Url

	if url == "" {
		return c.Status(fiber.StatusBadRequest).SendString("URL is required")
	}

	fileName, err := netex.DownloadAndExtractIfNeeded(url)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	defer os.Remove(fileName)

	netexData, err := netex.DecodeNetexData(fileName)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	err = netex.SaveNetexToDatabase(netexData, h.db)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.JSON(map[string]any{
		"success": true,
	})
}

func (h *NaolibHandler) SearchStopPlace(c *fiber.Ctx) error {
	query := c.Query("query")
	if query == "" {
		return c.Status(fiber.StatusBadRequest).SendString("Query is required")
	}

	rows, err := h.db.Query("SELECT id, name FROM NETEX_StopPlace WHERE name ILIKE $1", "%"+query+"%")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	defer rows.Close()

	var stopPlaces []models.StopPlace
	for rows.Next() {
		var stopPlace models.StopPlace
		err = rows.Scan(&stopPlace.ID, &stopPlace.Name)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}
		stopPlaces = append(stopPlaces, stopPlace)
	}

	type StopPlaceResponse struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}

	stopPlacesArray := make([]StopPlaceResponse, len(stopPlaces))
	for i, stopPlace := range stopPlaces {
		stopPlacesArray[i] = StopPlaceResponse{
			ID:   stopPlace.ID,
			Name: stopPlace.Name,
		}
	}

	return c.JSON(stopPlacesArray)
}

func (h *NaolibHandler) GenerateNetexRequest(c *fiber.Ctx) error {
	stops := []string{"CTRE2", "CTRE4"}
	request, err := netex.GenerateStopMonitoringRequest(stops)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error() + "\n")
	}

	return c.SendString(request + "\n")
}
