package handlers

import (
	"database/sql"
	"os"
	"strings"

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

func (h *NaolibHandler) GetDepartures(c *fiber.Ctx) error {
	stopPlaceId := c.Query("stopPlaceId")
	if stopPlaceId == "" {
		return c.Status(fiber.StatusBadRequest).SendString("Stop place ID is required")
	}

	// get the quays from the stop place id
	rows, err := h.db.Query("SELECT id FROM NETEX_Quay WHERE site_ref_stopplace_id = $1", stopPlaceId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	defer rows.Close()

	var quays []string
	for rows.Next() {
		var quay string
		err = rows.Scan(&quay)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}
		quays = append(quays, quay)
	}

	siriResponse, err := netex.CallStopMonitoringRequest(quays)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	departures := siriResponse.ServiceDelivery.StopMonitoringDelivery.MonitoredStopVisits

	type DepartureDirection struct {
		Direction  string                      `json:"direction"`
		Departures []models.MonitoredStopVisit `json:"departures"`
	}

	type Departures struct {
		DepartureDirectionAller  DepartureDirection `json:"aller"`
		DepartureDirectionRetour DepartureDirection `json:"retour"`
	}

	departuresMap := make(map[string]Departures)

	for _, departure := range departures {
		lineRef := departure.MonitoredVehicleJourney.LineRef

		lineDepartures, ok := departuresMap[lineRef]
		if !ok {
			lineDepartures = Departures{
				DepartureDirectionAller: DepartureDirection{
					Direction:  "",
					Departures: []models.MonitoredStopVisit{},
				},
				DepartureDirectionRetour: DepartureDirection{
					Direction:  "",
					Departures: []models.MonitoredStopVisit{},
				},
			}
		}

		if departure.MonitoredVehicleJourney.DirectionName == "A" {
			if lineDepartures.DepartureDirectionAller.Direction == "" {
				lineDepartures.DepartureDirectionAller.Direction = departure.MonitoredVehicleJourney.DestinationName
			} else {
				if !strings.Contains(lineDepartures.DepartureDirectionAller.Direction, departure.MonitoredVehicleJourney.DestinationName) {
					lineDepartures.DepartureDirectionAller.Direction += " / " + departure.MonitoredVehicleJourney.DestinationName
				}
			}
			lineDepartures.DepartureDirectionAller.Departures = append(lineDepartures.DepartureDirectionAller.Departures, departure)
		} else {
			// si la destination n'est pas B (valeur par défaut), ça signifie qu'on l'a déjà changé. on doit aller vérifier si c'est la même destination, sinon on rajoute la destination avec "/ <destination>"
			if lineDepartures.DepartureDirectionRetour.Direction == "" {
				lineDepartures.DepartureDirectionRetour.Direction = departure.MonitoredVehicleJourney.DestinationName
			} else {
				if !strings.Contains(lineDepartures.DepartureDirectionRetour.Direction, departure.MonitoredVehicleJourney.DestinationName) {
					lineDepartures.DepartureDirectionRetour.Direction += " / " + departure.MonitoredVehicleJourney.DestinationName
				}
			}

			lineDepartures.DepartureDirectionRetour.Departures = append(lineDepartures.DepartureDirectionRetour.Departures, departure)
		}

		departuresMap[lineRef] = lineDepartures
	}

	return c.JSON(departuresMap)
}
