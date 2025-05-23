package handlers

import (
	"database/sql"
	"os"

	"github.com/gofiber/fiber/v2"
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

	stopPlaces := netexData.DataObjects.GeneralFrame.Members.StopPlaces
	quays := netexData.DataObjects.GeneralFrame.Members.Quays

	// on crée une transaction et on supprime toutes les données existantes, avant de les insérer
	tx, err := h.db.Begin()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	defer tx.Rollback()

	_, err = tx.Exec("DELETE FROM NETEX_StopPlace")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	_, err = tx.Exec("DELETE FROM NETEX_Quay")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	_, err = tx.Exec("DELETE FROM NETEX_StopPlace_QuayRef")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	for _, stopPlace := range stopPlaces {
		_, err := tx.Exec("INSERT INTO NETEX_StopPlace (id, modification, name, longitude, latitude, transport_mode, other_transport_modes, stop_place_type, weighting) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)", stopPlace.ID, stopPlace.Modification, stopPlace.Name, stopPlace.Centroid.Location.Longitude, stopPlace.Centroid.Location.Latitude, stopPlace.TransportMode, stopPlace.OtherTransportModes, stopPlace.StopPlaceType, stopPlace.Weighting)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}
	}

	for _, quay := range quays {
		_, err := tx.Exec("INSERT INTO NETEX_Quay (id, name, longitude, latitude, site_ref_stopplace_id, transport_mode) VALUES ($1, $2, $3, $4, $5, $6)", quay.ID, quay.Name, quay.Centroid.Location.Longitude, quay.Centroid.Location.Latitude, quay.SiteRef.Ref, quay.TransportMode)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}
	}

	// now, populate the NETEX_StopPlace_QuayRef link table
	for _, stopPlace := range stopPlaces {
		for _, quayRef := range stopPlace.QuayRefs {
			_, err := tx.Exec("INSERT INTO NETEX_StopPlace_QuayRef (stop_place_id, quay_id, quay_ref_version) VALUES ($1, $2, $3)", stopPlace.ID, quayRef.Ref, quayRef.Version)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
			}
		}
	}

	err = tx.Commit()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.JSON(map[string]interface{}{
		"success": true,
	})
}
