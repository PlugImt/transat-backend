package carpool

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/plugimt/transat-backend/models"
	"github.com/plugimt/transat-backend/utils"
)

type CarpoolHandler struct {
	db *sql.DB
}

func NewCarpoolHandler(db *sql.DB) *CarpoolHandler {
	return &CarpoolHandler{
		db: db,
	}
}

// GetCarpools returns all active carpools (OPEN or FULL)
func (h *CarpoolHandler) GetCarpools(c *fiber.Ctx) error {
	utils.LogHeader("🚗 Get All Carpools")

	tripTypeFilter := c.Query("type", "all")
	utils.LogLineKeyValue(utils.LevelInfo, "Trip Type Filter", tripTypeFilter)

	var query string
	var args []interface{}
	argIndex := 1

	baseQuery := `
		SELECT 
			c.id_carpools,
			c.creator_email,
			c.trip_type,
			c.departure_place,
			c.destination,
			c.departure_time,
			c.status,
			n.first_name,
			n.last_name,
			COALESCE(n.profile_picture, '')
		FROM carpools c
		JOIN newf n ON c.creator_email = n.email
		WHERE c.status IN ('OPEN', 'FULL') AND c.departure_time >= NOW()
	`

	if tripTypeFilter != "all" {
		query = baseQuery + fmt.Sprintf(" AND c.trip_type = $%d ORDER BY c.departure_time ASC", argIndex)
		args = append(args, tripTypeFilter)
	} else {
		query = baseQuery + " ORDER BY c.departure_time ASC"
	}

	rows, err := h.db.Query(query, args...)
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to fetch carpools")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch carpools",
		})
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			utils.LogMessage(utils.LevelError, "Failed to close rows")
			utils.LogLineKeyValue(utils.LevelError, "Error", err)
		} else {
			utils.LogMessage(utils.LevelInfo, "Rows closed successfully")
		}
	}(rows)

	var carpools []map[string]interface{}
	for rows.Next() {
		var id int
		var creatorEmail, tripType, departurePlace, destination, status, firstName, lastName, profilePicture string
		var departureTime time.Time

		err := rows.Scan(
			&id, &creatorEmail, &tripType, &departurePlace, &destination,
			&departureTime, &status, &firstName, &lastName, &profilePicture,
		)
		if err != nil {
			utils.LogMessage(utils.LevelError, "Failed to scan carpool")
			utils.LogLineKeyValue(utils.LevelError, "Error", err)
			continue
		}

		carpool := map[string]interface{}{
			"id":              id,
			"trip_type":       tripType,
			"departure_place": departurePlace,
			"destination":     destination,
			"departure_time":  departureTime,
			"status":          status,
			"creator": map[string]interface{}{
				"email":           creatorEmail,
				"first_name":      firstName,
				"last_name":       lastName,
				"profile_picture": profilePicture,
			},
		}

		carpools = append(carpools, carpool)
	}

	utils.LogMessage(utils.LevelInfo, "Successfully fetched carpools")
	utils.LogLineKeyValue(utils.LevelInfo, "Count", len(carpools))
	utils.LogFooter()

	return c.JSON(carpools)
}

// GetCarpoolByID returns detailed information for a specific carpool
func (h *CarpoolHandler) GetCarpoolByID(c *fiber.Ctx) error {
	utils.LogHeader("🚗 Get Carpool By ID")

	carpoolID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		utils.LogMessage(utils.LevelWarn, "Invalid carpool ID")
		utils.LogFooter()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid carpool ID",
		})
	}

	utils.LogLineKeyValue(utils.LevelInfo, "Carpool ID", carpoolID)

	query := `
		SELECT 
			c.id_carpools,
			c.creator_email,
			c.trip_type,
			c.departure_place,
			c.destination,
			c.departure_time,
			c.contact_details,
			COALESCE(c.description, ''),
			c.status,
			c.creation_date,
			n.first_name,
			n.last_name,
			COALESCE(n.profile_picture, ''),
			n.graduation_year
		FROM carpools c
		JOIN newf n ON c.creator_email = n.email
		WHERE c.id_carpools = $1
	`

	var id int
	var creatorEmail, tripType, departurePlace, destination, contactDetails, description, status, firstName, lastName, profilePicture string
	var departureTime, creationDate time.Time
	var graduationYear sql.NullInt64

	err = h.db.QueryRow(query, carpoolID).Scan(
		&id, &creatorEmail, &tripType, &departurePlace, &destination,
		&departureTime, &contactDetails, &description, &status, &creationDate,
		&firstName, &lastName, &profilePicture, &graduationYear,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			utils.LogMessage(utils.LevelWarn, "Carpool not found")
			utils.LogFooter()
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Carpool not found",
			})
		}
		utils.LogMessage(utils.LevelError, "Failed to fetch carpool")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch carpool",
		})
	}

	response := map[string]interface{}{
		"id":              id,
		"trip_type":       tripType,
		"departure_place": departurePlace,
		"destination":     destination,
		"departure_time":  departureTime,
		"contact_details": contactDetails,
		"description":     description,
		"status":          status,
		"creation_date":   creationDate,
		"creator": map[string]interface{}{
			"email":           creatorEmail,
			"first_name":      firstName,
			"last_name":       lastName,
			"profile_picture": profilePicture,
			"graduation_year": graduationYear.Int64,
		},
	}

	utils.LogMessage(utils.LevelInfo, "Successfully fetched carpool details")
	utils.LogFooter()

	return c.JSON(response)
}

// CreateCarpool creates a new carpool offer
func (h *CarpoolHandler) CreateCarpool(c *fiber.Ctx) error {
	utils.LogHeader("🚗 Create Carpool")

	userEmail := c.Locals("email").(string)
	utils.LogLineKeyValue(utils.LevelInfo, "Creator", userEmail)

	var req models.CreateCarpoolRequest
	if err := c.BodyParser(&req); err != nil {
		utils.LogMessage(utils.LevelError, "Failed to parse request body")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate basic requirements
	if req.TripType == "" || req.DeparturePlace == "" || req.Destination == "" || req.ContactDetails == "" {
		utils.LogMessage(utils.LevelWarn, "Missing required fields")
		utils.LogFooter()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Trip type, departure place, destination, and contact details are required",
		})
	}

	var carpoolID int
	insertQuery := `
		INSERT INTO carpools (creator_email, trip_type, departure_place, destination, departure_time, contact_details, description)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id_carpools
	`
	err := h.db.QueryRow(
		insertQuery,
		userEmail,
		req.TripType,
		req.DeparturePlace,
		req.Destination,
		req.DepartureTime,
		req.ContactDetails,
		req.Description,
	).Scan(&carpoolID)

	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to create carpool")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create carpool",
		})
	}

	utils.LogMessage(utils.LevelInfo, "Carpool created successfully")
	utils.LogLineKeyValue(utils.LevelInfo, "Carpool ID", carpoolID)
	utils.LogFooter()

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":    "Carpool created successfully",
		"carpool_id": carpoolID,
	})
}

// UpdateCarpoolStatus allows the creator to mark a carpool as FULL or ARCHIVED
func (h *CarpoolHandler) UpdateCarpoolStatus(c *fiber.Ctx) error {
	utils.LogHeader("🔄 Update Carpool Status")

	carpoolID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		utils.LogMessage(utils.LevelWarn, "Invalid carpool ID")
		utils.LogFooter()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid carpool ID"})
	}

	userEmail := c.Locals("email").(string)
	utils.LogLineKeyValue(utils.LevelInfo, "Carpool ID", carpoolID)
	utils.LogLineKeyValue(utils.LevelInfo, "User", userEmail)

	// Verify ownership
	var creator string
	err = h.db.QueryRow("SELECT creator_email FROM carpools WHERE id_carpools = $1", carpoolID).Scan(&creator)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.LogMessage(utils.LevelWarn, "Carpool not found")
			utils.LogFooter()
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Carpool not found"})
		}
		utils.LogMessage(utils.LevelError, "Failed to get carpool")
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get carpool information"})
	}

	if creator != userEmail {
		utils.LogMessage(utils.LevelWarn, "User not authorized to update carpool")
		utils.LogFooter()
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Only creators can update the status"})
	}

	// Parse new status
	var req struct {
		Status string `json:"status"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if req.Status != "OPEN" && req.Status != "FULL" && req.Status != "ARCHIVED" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Status must be OPEN, FULL, or ARCHIVED"})
	}

	query := "UPDATE carpools SET status = $1, updated_date = CURRENT_TIMESTAMP WHERE id_carpools = $2"
	if _, err := h.db.Exec(query, req.Status, carpoolID); err != nil {
		utils.LogMessage(utils.LevelError, "Failed to update carpool status")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update carpool status"})
	}

	utils.LogMessage(utils.LevelInfo, "Carpool status updated successfully")
	utils.LogLineKeyValue(utils.LevelInfo, "New Status", req.Status)
	utils.LogFooter()
	return c.JSON(fiber.Map{"message": "Carpool status updated successfully"})
}

// DeleteCarpool allows the creator to delete their carpool offer
func (h *CarpoolHandler) DeleteCarpool(c *fiber.Ctx) error {
	utils.LogHeader("🗑️ Delete Carpool")

	carpoolID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		utils.LogMessage(utils.LevelWarn, "Invalid carpool ID")
		utils.LogFooter()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid carpool ID",
		})
	}

	userEmail := c.Locals("email").(string)

	// Verify ownership
	var creator string
	err = h.db.QueryRow("SELECT creator_email FROM carpools WHERE id_carpools = $1", carpoolID).Scan(&creator)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Carpool not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get carpool information"})
	}

	if creator != userEmail {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Only carpool creators can delete the offer"})
	}

	_, err = h.db.Exec("DELETE FROM carpools WHERE id_carpools = $1", carpoolID)
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to delete carpool")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete carpool",
		})
	}

	utils.LogMessage(utils.LevelInfo, "Carpool deleted successfully")
	utils.LogFooter()
	return c.JSON(fiber.Map{
		"message": "Carpool deleted successfully",
	})
}
