package handlers

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/plugimt/transat-backend/models"
	"github.com/plugimt/transat-backend/services"
	"github.com/plugimt/transat-backend/services/naolib/netex"
	"github.com/plugimt/transat-backend/services/naolib/siri"
)

const (
	ChantrerieStopPlaceId = "FR_NAOLIB:StopPlace:244"
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
	departures, err := h.service.GetDepartures(ChantrerieStopPlaceId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.JSON(departures)
}

func (h *NaolibHandler) ImportNetexOffer(c *fiber.Ctx) error {
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

	fileName, err := netex.DownloadAndExtractIfNeededOffer(url)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	defer os.Remove(fileName)

	netexData, err := netex.DecodeNetexOfferData(fileName)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	fmt.Println(netexData)

	err = netex.SaveNetexOfferToDatabase(netexData, h.db)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.JSON(map[string]any{
		"success": true,
	})
}

func (h *NaolibHandler) ImportNetexStops(c *fiber.Ctx) error {
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

	netexData, err := netex.DecodeNetexStopsData(fileName)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	err = netex.SaveNetexStopsToDatabase(netexData, h.db)
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
	request, err := siri.GenerateStopMonitoringRequest(stops)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error() + "\n")
	}

	return c.SendString(request + "\n")
}

func (h *NaolibHandler) GetDepartures(c *fiber.Ctx) error {
	stopPlaceId := c.Query("stopPlaceId")
	if stopPlaceId == "" {
		return c.Status(fiber.StatusBadRequest).SendString("Stop place ID is required")
	}

	departuresMap, err := h.service.GetDepartures(stopPlaceId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.JSON(departuresMap)
}
