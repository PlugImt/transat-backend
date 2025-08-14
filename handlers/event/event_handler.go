package event

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/plugimt/transat-backend/models"
	"github.com/plugimt/transat-backend/utils"
)

type EventHandler struct {
	db *sql.DB
}

func NewEventHandler(db *sql.DB) *EventHandler {
	return &EventHandler{
		db: db,
	}
}

// GetEvent returns all events based on time filter
func (h *EventHandler) GetEvent(c *fiber.Ctx) error {
	utils.LogHeader("🎉 Get All Events")

	timeFilter := c.Query("time", "upcoming") // default to upcoming
	utils.LogLineKeyValue(utils.LevelInfo, "Time Filter", timeFilter)

	var query string
	var args []interface{}

	switch timeFilter {
	case "past":
		query = `
			SELECT 
				e.id_events,
				e.id_club,
				e.name,
				e.start_date,
				e.end_date,
				e.location,
				e.picture,
				COALESCE(attendee_count.count, 0) as attendee_count
			FROM events e
			LEFT JOIN (
				SELECT id_events, COUNT(*) as count
				FROM events_attendents
				GROUP BY id_events
			) attendee_count ON e.id_events = attendee_count.id_events
			WHERE e.start_date < NOW()
			ORDER BY e.start_date DESC
		`
	case "all":
		query = `
			SELECT 
				e.id_events,
				e.id_club,
				e.name,
				e.start_date,
				e.end_date,
				e.location,
				e.picture,
				COALESCE(attendee_count.count, 0) as attendee_count
			FROM events e
			LEFT JOIN (
				SELECT id_events, COUNT(*) as count
				FROM events_attendents
				GROUP BY id_events
			) attendee_count ON e.id_events = attendee_count.id_events
			ORDER BY e.start_date ASC
		`
	default: // upcoming
		query = `
			SELECT 
				e.id_events,
				e.id_club,
				e.name,
				e.start_date,
				e.end_date,
				e.location,
				e.picture,
				COALESCE(attendee_count.count, 0) as attendee_count
			FROM events e
			LEFT JOIN (
				SELECT id_events, COUNT(*) as count
				FROM events_attendents
				GROUP BY id_events
			) attendee_count ON e.id_events = attendee_count.id_events
			WHERE e.start_date >= NOW()
			ORDER BY e.start_date ASC
		`
	}

	rows, err := h.db.Query(query, args...)
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to fetch events")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch events",
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

	var events []map[string]interface{}
	for rows.Next() {
		var id int
		var idClub int
		var name, location, picture string
		var startDate time.Time
		var endDate sql.NullTime
		var attendeeCount int

		err := rows.Scan(&id, &idClub, &name, &startDate, &endDate, &location, &picture, &attendeeCount)
		if err != nil {
			utils.LogMessage(utils.LevelError, "Failed to scan event")
			utils.LogLineKeyValue(utils.LevelError, "Error", err)
			continue
		}

		event := map[string]interface{}{
			"id":             id,
			"id_club":        idClub,
			"name":           name,
			"start_date":     startDate,
			"end_date":       endDate.Time,
			"location":       location,
			"picture":        picture,
			"attendee_count": attendeeCount,
		}

		// Handle null end date
		if !endDate.Valid {
			event["end_date"] = nil
		}

		// Get first 10 attendee profile pictures
		photosQuery := `
			SELECT COALESCE(n.profile_picture,'')
			FROM events_attendents ea
			JOIN newf n ON ea.email = n.email
			WHERE ea.id_events = $1 AND n.profile_picture IS NOT NULL AND n.profile_picture != ''
			ORDER BY n.first_name, n.last_name
			LIMIT 10
		`
		photoRows, err := h.db.Query(photosQuery, id)
		if err == nil {
			defer photoRows.Close()
			var memberPhotos []string
			for photoRows.Next() {
				var photo string
				if err := photoRows.Scan(&photo); err == nil {
					memberPhotos = append(memberPhotos, photo)
				}
			}
			event["member_photos"] = memberPhotos
		}

		events = append(events, event)
	}

	utils.LogMessage(utils.LevelInfo, "Successfully fetched events")
	utils.LogLineKeyValue(utils.LevelInfo, "Count", len(events))
	utils.LogFooter()

	return c.JSON(events)
}

func (h *EventHandler) GetEventByClubID(c *fiber.Ctx) error {
	utils.LogHeader("🎉 Get All Events")

	timeFilter := c.Query("time", "upcoming") // default to upcoming
	utils.LogLineKeyValue(utils.LevelInfo, "Time Filter", timeFilter)

	clubID, err := strconv.Atoi(c.Params("id"))

	var query string
	var args []interface{}

	switch timeFilter {
	case "past":
		query = `
			SELECT 
				e.id_events,
				e.name,
				e.start_date,
				e.end_date,
				e.location,
				e.picture,
				COALESCE(attendee_count.count, 0) as attendee_count
			FROM events e
			LEFT JOIN (
				SELECT id_events, COUNT(*) as count
				FROM events_attendents
				GROUP BY id_events
			) attendee_count ON e.id_events = attendee_count.id_events
			WHERE e.start_date < NOW()
			AND e.id_club = $1
			ORDER BY e.start_date DESC
		`
	case "all":
		query = `
			SELECT 
				e.id_events,
				e.name,
				e.start_date,
				e.end_date,
				e.location,
				e.picture,
				COALESCE(attendee_count.count, 0) as attendee_count
			FROM events e
			LEFT JOIN (
				SELECT id_events, COUNT(*) as count
				FROM events_attendents
				GROUP BY id_events
			) attendee_count ON e.id_events = attendee_count.id_events
			WHERE e.id_club = $1
			ORDER BY e.start_date ASC
		`
	default: // upcoming
		query = `
			SELECT 
				e.id_events,
				e.name,
				e.start_date,
				e.end_date,
				e.location,
				e.picture,
				COALESCE(attendee_count.count, 0) as attendee_count
			FROM events e
			LEFT JOIN (
				SELECT id_events, COUNT(*) as count
				FROM events_attendents
				GROUP BY id_events
			) attendee_count ON e.id_events = attendee_count.id_events
			WHERE e.start_date >= NOW()
			AND e.id_club = $1
			ORDER BY e.start_date ASC
		`
	}

	if err != nil {
		utils.LogMessage(utils.LevelError, "Invalid club ID")
		utils.LogFooter()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid club ID",
		})
	}

	args = append(args, clubID)
	rows, err := h.db.Query(query, args...)
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to fetch events")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch events",
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

	var events []map[string]interface{}
	for rows.Next() {
		var id int
		var name, location, picture string
		var startDate time.Time
		var endDate sql.NullTime
		var attendeeCount int

		err := rows.Scan(&id, &name, &startDate, &endDate, &location, &picture, &attendeeCount)
		if err != nil {
			utils.LogMessage(utils.LevelError, "Failed to scan event")
			utils.LogLineKeyValue(utils.LevelError, "Error", err)
			continue
		}

		event := map[string]interface{}{
			"id":             id,
			"name":           name,
			"start_date":     startDate,
			"end_date":       endDate.Time,
			"location":       location,
			"picture":        picture,
			"attendee_count": attendeeCount,
		}

		// Handle null end date
		if !endDate.Valid {
			event["end_date"] = nil
		}

		// Get first 10 attendee profile pictures
		photosQuery := `
			SELECT COALESCE(n.profile_picture, '')
			FROM events_attendents ea
			JOIN newf n ON ea.email = n.email
			WHERE ea.id_events = $1 AND n.profile_picture IS NOT NULL AND n.profile_picture != ''
			ORDER BY n.first_name, n.last_name
			LIMIT 10
		`
		photoRows, err := h.db.Query(photosQuery, id)
		if err == nil {
			defer photoRows.Close()
			var memberPhotos []string
			for photoRows.Next() {
				var photo string
				if err := photoRows.Scan(&photo); err == nil {
					memberPhotos = append(memberPhotos, photo)
				}
			}
			event["member_photos"] = memberPhotos
		}

		events = append(events, event)
	}

	utils.LogMessage(utils.LevelInfo, "Successfully fetched events")
	utils.LogLineKeyValue(utils.LevelInfo, "Count", len(events))
	utils.LogFooter()

	return c.JSON(events)
}

// GetEventByID returns a specific event by ID with detailed information
func (h *EventHandler) GetEventByID(c *fiber.Ctx) error {
	utils.LogHeader("🎉 Get Event By ID")

	eventID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		utils.LogMessage(utils.LevelWarn, "Invalid event ID")
		utils.LogFooter()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid event ID",
		})
	}

	utils.LogLineKeyValue(utils.LevelInfo, "Event ID", eventID)

	// Get event basic info
	eventQuery := `
		SELECT 
			e.id_events,
			e.name,
			e.description,
			e.link,
			e.start_date,
			e.end_date,
			e.location,
			e.picture,
			e.creator,
			e.id_club
		FROM events e
		WHERE e.id_events = $1
	`

	var event models.Event
	var description, link sql.NullString
	var endDate sql.NullTime

	err = h.db.QueryRow(eventQuery, eventID).Scan(
		&event.ID,
		&event.Name,
		&description,
		&link,
		&event.StartDate,
		&endDate,
		&event.Location,
		&event.Picture,
		&event.Creator,
		&event.ClubID,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			utils.LogMessage(utils.LevelWarn, "Event not found")
			utils.LogFooter()
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Event not found",
			})
		}
		utils.LogMessage(utils.LevelError, "Failed to fetch event")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch event",
		})
	}

	event.Description = description.String
	event.Link = link.String
	if endDate.Valid {
		event.EndDate = &endDate.Time
	}

	// Get creator details
	creatorQuery := `
		SELECT 
			n.email,
			n.first_name,
			n.last_name,
			COALESCE(n.profile_picture,''),
			n.graduation_year
		FROM newf n
		WHERE n.email = $1
	`

	var creator map[string]interface{}
	var creatorEmail, creatorFirstName, creatorLastName string
	var creatorProfilePicture sql.NullString
	var creatorGraduationYear sql.NullInt64

	err = h.db.QueryRow(creatorQuery, event.Creator).Scan(
		&creatorEmail,
		&creatorFirstName,
		&creatorLastName,
		&creatorProfilePicture,
		&creatorGraduationYear,
	)

	if err == nil {
		creator = map[string]interface{}{
			"email":           creatorEmail,
			"first_name":      creatorFirstName,
			"last_name":       creatorLastName,
			"profile_picture": creatorProfilePicture.String,
			"graduation_year": creatorGraduationYear.Int64,
		}
	} else if err != sql.ErrNoRows {
		utils.LogMessage(utils.LevelError, "Failed to fetch creator")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
	}

	// Get club details
	clubQuery := `
		SELECT 
			c.name,
			c.picture
		FROM clubs c
		WHERE c.id_clubs = $1
	`

	var club map[string]interface{}
	var clubName, clubPicture string

	err = h.db.QueryRow(clubQuery, event.ClubID).Scan(&clubName, &clubPicture)
	if err == nil {
		club = map[string]interface{}{
			"name":    clubName,
			"picture": clubPicture,
		}
	} else if err != sql.ErrNoRows {
		utils.LogMessage(utils.LevelError, "Failed to fetch club")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
	}

	// Get total attendee count
	var attendeeCount int
	err = h.db.QueryRow("SELECT COUNT(*) FROM events_attendents WHERE id_events = $1", eventID).Scan(&attendeeCount)
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to get attendee count")
		attendeeCount = 0
	}

	// Get first 15 attendees
	attendeesQuery := `
		SELECT 
			n.email,
			n.first_name,
			n.last_name,
			COALESCE(n.profile_picture,''),
			n.graduation_year
		FROM events_attendents ea
		JOIN newf n ON ea.email = n.email
		WHERE ea.id_events = $1
		ORDER BY n.first_name, n.last_name
		LIMIT 15
	`

	rows, err := h.db.Query(attendeesQuery, eventID)
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to fetch attendees")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
	} else {
		defer func(rows *sql.Rows) {
			err := rows.Close()
			if err != nil {
				utils.LogMessage(utils.LevelError, "Failed to close rows")
				utils.LogLineKeyValue(utils.LevelError, "Error", err)
			} else {
				utils.LogMessage(utils.LevelInfo, "Rows closed successfully")
			}
		}(rows)
	}

	var attendees []map[string]interface{}
	if rows != nil {
		for rows.Next() {
			var email, firstName, lastName string
			var profilePicture sql.NullString
			var graduationYear sql.NullInt64

			if err := rows.Scan(&email, &firstName, &lastName, &profilePicture, &graduationYear); err == nil {
				attendee := map[string]interface{}{
					"email":           email,
					"first_name":      firstName,
					"last_name":       lastName,
					"profile_picture": profilePicture.String,
					"graduation_year": graduationYear.Int64,
				}
				attendees = append(attendees, attendee)
			}
		}
	}

	// Build response
	response := map[string]interface{}{
		"id":             event.ID,
		"name":           event.Name,
		"description":    event.Description,
		"link":           event.Link,
		"start_date":     event.StartDate,
		"end_date":       event.EndDate,
		"location":       event.Location,
		"picture":        event.Picture,
		"attendee_count": attendeeCount,
		"attendees":      attendees,
	}

	if creator != nil {
		response["creator"] = creator
	}

	if club != nil {
		response["club"] = club
	}

	utils.LogMessage(utils.LevelInfo, "Successfully fetched event details")
	utils.LogFooter()

	return c.JSON(response)
}

// GetEventMembers returns all attendees of a specific event with detailed info
func (h *EventHandler) GetEventMembers(c *fiber.Ctx) error {
	utils.LogHeader("👥 Get Event Members")

	eventID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		utils.LogMessage(utils.LevelWarn, "Invalid event ID")
		utils.LogFooter()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid event ID",
		})
	}

	utils.LogLineKeyValue(utils.LevelInfo, "Event ID", eventID)

	// First check if event exists
	var eventExists bool
	err = h.db.QueryRow("SELECT EXISTS(SELECT 1 FROM events WHERE id_events = $1)", eventID).Scan(&eventExists)
	if err != nil || !eventExists {
		utils.LogMessage(utils.LevelWarn, "Event not found")
		utils.LogFooter()
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Event not found",
		})
	}

	query := `
		SELECT 
			ea.email,
			n.first_name,
			n.last_name,
			COALESCE(n.profile_picture,''),
			n.graduation_year
		FROM events_attendents ea
		JOIN newf n ON ea.email = n.email
		WHERE ea.id_events = $1
		ORDER BY n.first_name, n.last_name
	`

	rows, err := h.db.Query(query, eventID)
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to fetch event members")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch event members",
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

	var members []map[string]interface{}
	for rows.Next() {
		var email, firstName, lastName string
		var profilePicture sql.NullString
		var graduationYear sql.NullInt64

		err := rows.Scan(&email, &firstName, &lastName, &profilePicture, &graduationYear)
		if err != nil {
			utils.LogMessage(utils.LevelError, "Failed to scan member")
			continue
		}

		member := map[string]interface{}{
			"email":           email,
			"first_name":      firstName,
			"last_name":       lastName,
			"profile_picture": profilePicture.String,
			"graduation_year": graduationYear.Int64,
		}

		members = append(members, member)
	}

	utils.LogMessage(utils.LevelInfo, "Successfully fetched event members")
	utils.LogLineKeyValue(utils.LevelInfo, "Count", len(members))
	utils.LogFooter()

	return c.JSON(fiber.Map{
		"members": members,
		"count":   len(members),
	})
}

// CreateEvent creates a new event
func (h *EventHandler) CreateEvent(c *fiber.Ctx) error {
	utils.LogHeader("🎉 Create Event")

	userEmail := c.Locals("email").(string)
	utils.LogLineKeyValue(utils.LevelInfo, "Creator", userEmail)

	var req models.CreateEventRequest
	if err := c.BodyParser(&req); err != nil {
		utils.LogMessage(utils.LevelError, "Failed to parse request body")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate required fields
	if req.Name == "" || req.Location == "" || req.ClubID == 0 {
		utils.LogMessage(utils.LevelWarn, "Missing required fields")
		utils.LogFooter()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Name, location, and club ID are required",
		})
	}

	utils.LogLineKeyValue(utils.LevelInfo, "Event Name", req.Name)

	// Check if club exists
	var clubExists bool
	err := h.db.QueryRow("SELECT EXISTS(SELECT 1 FROM clubs WHERE id_clubs = $1)", req.ClubID).Scan(&clubExists)
	if err != nil || !clubExists {
		utils.LogMessage(utils.LevelWarn, "Club not found")
		utils.LogFooter()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Club not found",
		})
	}

	// Start transaction
	tx, err := h.db.Begin()
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to begin transaction")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database transaction error",
		})
	}
	defer tx.Rollback()

	// Create event
	var eventID int
	insertEventQuery := `
		INSERT INTO events (name, description, link, start_date, end_date, location, picture, creator, id_club)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id_events
	`
	err = tx.QueryRow(
		insertEventQuery,
		req.Name,
		req.Description,
		req.Link,
		req.StartDate,
		req.EndDate,
		req.Location,
		req.Picture,
		userEmail,
		req.ClubID,
	).Scan(&eventID)

	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to create event")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create event",
		})
	}

	// Add creator as attendee
	_, err = tx.Exec(
		"INSERT INTO events_attendents (email, id_events) VALUES ($1, $2)",
		userEmail, eventID,
	)
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to add creator as attendee")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to add creator as attendee",
		})
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		utils.LogMessage(utils.LevelError, "Failed to commit transaction")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save event",
		})
	}

	utils.LogMessage(utils.LevelInfo, "Event created successfully")
	utils.LogLineKeyValue(utils.LevelInfo, "Event ID", eventID)
	utils.LogFooter()

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Event created successfully",
		"event": map[string]interface{}{
			"id_event":    eventID,
			"name":        req.Name,
			"description": req.Description,
			"link":        req.Link,
			"start_date":  req.StartDate,
			"end_date":    req.EndDate,
			"location":    req.Location,
			"picture":     req.Picture,
			"creator":     userEmail,
			"id_club":     req.ClubID,
		},
	})
}

// UpdateEvent updates event information (only creator can update)
func (h *EventHandler) UpdateEvent(c *fiber.Ctx) error {
	utils.LogHeader("🎉 Update Event")

	eventID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		utils.LogMessage(utils.LevelWarn, "Invalid event ID")
		utils.LogFooter()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid event ID",
		})
	}

	userEmail := c.Locals("email").(string)
	utils.LogLineKeyValue(utils.LevelInfo, "Event ID", eventID)
	utils.LogLineKeyValue(utils.LevelInfo, "User", userEmail)

	// Check if user is the creator of this event
	var creator string
	err = h.db.QueryRow("SELECT creator FROM events WHERE id_events = $1", eventID).Scan(&creator)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.LogMessage(utils.LevelWarn, "Event not found")
			utils.LogFooter()
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Event not found",
			})
		}
		utils.LogMessage(utils.LevelError, "Failed to get event")
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get event information",
		})
	}

	if creator != userEmail {
		utils.LogMessage(utils.LevelWarn, "User not authorized to update event")
		utils.LogFooter()
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Only event creators can update event information",
		})
	}

	var req models.UpdateEventRequest
	if err := c.BodyParser(&req); err != nil {
		utils.LogMessage(utils.LevelError, "Failed to parse request body")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Build dynamic update query
	var setParts []string
	var args []interface{}
	argIndex := 1

	if req.Name != "" {
		setParts = append(setParts, fmt.Sprintf("name = $%d", argIndex))
		args = append(args, req.Name)
		argIndex++
	}
	if req.Description != "" {
		setParts = append(setParts, fmt.Sprintf("description = $%d", argIndex))
		args = append(args, req.Description)
		argIndex++
	}
	if req.Link != "" {
		setParts = append(setParts, fmt.Sprintf("link = $%d", argIndex))
		args = append(args, req.Link)
		argIndex++
	}
	if req.StartDate != nil {
		setParts = append(setParts, fmt.Sprintf("start_date = $%d", argIndex))
		args = append(args, req.StartDate)
		argIndex++
	}
	if req.EndDate != nil {
		setParts = append(setParts, fmt.Sprintf("end_date = $%d", argIndex))
		args = append(args, req.EndDate)
		argIndex++
	}
	if req.Location != "" {
		setParts = append(setParts, fmt.Sprintf("location = $%d", argIndex))
		args = append(args, req.Location)
		argIndex++
	}
	if req.Picture != "" {
		setParts = append(setParts, fmt.Sprintf("picture = $%d", argIndex))
		args = append(args, req.Picture)
		argIndex++
	}

	if len(setParts) == 0 {
		utils.LogMessage(utils.LevelWarn, "No fields to update")
		utils.LogFooter()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No fields to update",
		})
	}

	query := fmt.Sprintf("UPDATE events SET %s WHERE id_events = $%d", strings.Join(setParts, ", "), argIndex)
	args = append(args, eventID)

	_, err = h.db.Exec(query, args...)
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to update event")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update event",
		})
	}

	utils.LogMessage(utils.LevelInfo, "Event updated successfully")
	utils.LogFooter()

	return c.JSON(fiber.Map{
		"message": "Event updated successfully",
	})
}

func (h *EventHandler) DeleteEvent(c *fiber.Ctx) error {
	utils.LogHeader("🗑️ Delete Event")

	eventID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		utils.LogMessage(utils.LevelWarn, "Invalid event ID")
		utils.LogFooter()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid event ID",
		})
	}

	userEmail := c.Locals("email").(string)
	utils.LogLineKeyValue(utils.LevelInfo, "Event ID", eventID)
	utils.LogLineKeyValue(utils.LevelInfo, "User", userEmail)

	// Check if user is the creator of this event
	var creator string
	err = h.db.QueryRow("SELECT creator FROM events WHERE id_events = $1", eventID).Scan(&creator)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.LogMessage(utils.LevelWarn, "Event not found")
			utils.LogFooter()
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Event not found",
			})
		}
		utils.LogMessage(utils.LevelError, "Failed to get event")
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get event information",
		})
	}

	if creator != userEmail {
		utils.LogMessage(utils.LevelWarn, "User not authorized to delete event")
		utils.LogFooter()
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Only event creators can delete the event",
		})
	}

	// Start transaction
	tx, err := h.db.Begin()
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to begin transaction")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database transaction error",
		})
	}
	defer tx.Rollback()

	// Delete attendees first
	result, err := tx.Exec("DELETE FROM events_attendents WHERE id_events = $1", eventID)
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to delete attendees")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete event attendees",
		})
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to get rows affected")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete event attendees",
		})
	}
	if rowsAffected == 0 {
		utils.LogMessage(utils.LevelWarn, "No attendees found for event")
		utils.LogFooter()
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "No attendees found for this event",
		})
	}
	// Delete the event
	_, err = tx.Exec("DELETE FROM events WHERE id_events = $1", eventID)
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to delete event")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete event",
		})
	}
	// Commit transaction
	if err = tx.Commit(); err != nil {
		utils.LogMessage(utils.LevelError, "Failed to commit transaction")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete event",
		})
	}
	utils.LogMessage(utils.LevelInfo, "Event deleted successfully")
	utils.LogFooter()
	return c.JSON(fiber.Map{
		"message": "Event deleted successfully",
	})
}

// JoinEvent adds the user to an event
func (h *EventHandler) JoinEvent(c *fiber.Ctx) error {
	utils.LogHeader("🤝 Join Event")

	eventID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		utils.LogMessage(utils.LevelWarn, "Invalid event ID")
		utils.LogFooter()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid event ID",
		})
	}

	userEmail := c.Locals("email").(string)
	utils.LogLineKeyValue(utils.LevelInfo, "Event ID", eventID)
	utils.LogLineKeyValue(utils.LevelInfo, "User", userEmail)

	// Check if event exists
	var eventExists bool
	err = h.db.QueryRow("SELECT EXISTS(SELECT 1 FROM events WHERE id_events = $1)", eventID).Scan(&eventExists)
	if err != nil || !eventExists {
		utils.LogMessage(utils.LevelWarn, "Event not found")
		utils.LogFooter()
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Event not found",
		})
	}

	// Check if user is already an attendee
	var isAttendee bool
	err = h.db.QueryRow("SELECT EXISTS(SELECT 1 FROM events_attendents WHERE email = $1 AND id_events = $2)", userEmail, eventID).Scan(&isAttendee)
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to check attendance")
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to check attendance status",
		})
	}

	if isAttendee {
		utils.LogMessage(utils.LevelWarn, "User already an attendee")
		utils.LogFooter()
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "You are already attending this event",
		})
	}

	// Add user to event
	_, err = h.db.Exec("INSERT INTO events_attendents (email, id_events) VALUES ($1, $2)", userEmail, eventID)
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to join event")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to join event",
		})
	}

	utils.LogMessage(utils.LevelInfo, "Successfully joined event")
	utils.LogFooter()

	return c.JSON(fiber.Map{
		"message": "Successfully joined the event",
	})
}

// LeaveEvent removes the user from an event
func (h *EventHandler) LeaveEvent(c *fiber.Ctx) error {
	utils.LogHeader("👋 Leave Event")

	eventID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		utils.LogMessage(utils.LevelWarn, "Invalid event ID")
		utils.LogFooter()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid event ID",
		})
	}

	userEmail := c.Locals("email").(string)
	utils.LogLineKeyValue(utils.LevelInfo, "Event ID", eventID)
	utils.LogLineKeyValue(utils.LevelInfo, "User", userEmail)

	// Check if user is an attendee
	var isAttendee bool
	err = h.db.QueryRow("SELECT EXISTS(SELECT 1 FROM events_attendents WHERE email = $1 AND id_events = $2)", userEmail, eventID).Scan(&isAttendee)
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to check attendance")
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to check attendance status",
		})
	}

	if !isAttendee {
		utils.LogMessage(utils.LevelWarn, "User not an attendee")
		utils.LogFooter()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "You are not attending this event",
		})
	}

	// Check if user is the creator of this event
	var creator string
	err = h.db.QueryRow("SELECT creator FROM events WHERE id_events = $1", eventID).Scan(&creator)
	if err == nil && creator == userEmail {
		utils.LogMessage(utils.LevelWarn, "Creator cannot leave event")
		utils.LogFooter()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Event creators cannot leave the event. Please delete the event instead.",
		})
	}

	// Remove user from event
	result, err := h.db.Exec("DELETE FROM events_attendents WHERE email = $1 AND id_events = $2", userEmail, eventID)
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to leave event")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to leave event",
		})
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		utils.LogMessage(utils.LevelWarn, "No rows affected when leaving event")
		utils.LogFooter()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "You are not attending this event",
		})
	}

	utils.LogMessage(utils.LevelInfo, "Successfully left event")
	utils.LogFooter()

	return c.JSON(fiber.Map{
		"message": "Successfully left the event",
	})
}
