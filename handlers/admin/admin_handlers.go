package admin

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/plugimt/transat-backend/models"
	"github.com/plugimt/transat-backend/utils"
	"golang.org/x/crypto/bcrypt"
)

type AdminHandler struct {
	DB *sql.DB
}

func NewAdminHandler(db *sql.DB) *AdminHandler {
	return &AdminHandler{DB: db}
}

func (h *AdminHandler) GetAllUsers(c *fiber.Ctx) error {
	utils.LogHeader("👥 Get All Users (Admin)")

	query := `
		SELECT 
			n.id_newf,
			n.email,
			n.first_name,
			n.last_name,
			COALESCE(n.phone_number, '') AS phone_number,
			COALESCE(n.profile_picture, '') AS profile_picture,
			n.graduation_year,
			COALESCE(n.formation_name, '') AS formation_name,
			COALESCE(n.campus, '') AS campus,
			COALESCE(l.code, 'fr') AS language,
			n.creation_date,
			n.verification_code,
			n.verification_code_expiration
		FROM newf n
		LEFT JOIN languages l ON n.language = l.id_languages
		ORDER BY n.creation_date DESC;
	`

	rows, err := h.DB.Query(query)
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to fetch users")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch users"})
	}
	defer rows.Close()

	type UserWithRoles struct {
		models.Newf
		Roles []string `json:"roles,omitempty"`
	}

	var users []UserWithRoles
	userRolesMap := make(map[string][]string)

	for rows.Next() {
		var user UserWithRoles
		var verificationCode sql.NullString
		var verificationCodeExpiration sql.NullString
		var graduationYear sql.NullInt32

		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.FirstName,
			&user.LastName,
			&user.PhoneNumber,
			&user.ProfilePicture,
			&graduationYear,
			&user.FormationName,
			&user.Campus,
			&user.Language,
			&user.CreationDate,
			&verificationCode,
			&verificationCodeExpiration,
		)
		if err != nil {
			utils.LogMessage(utils.LevelError, "Failed to scan user row")
			utils.LogLineKeyValue(utils.LevelError, "Error", err)
			continue
		}

		if verificationCode.Valid {
			user.VerificationCode = verificationCode.String
		}
		if verificationCodeExpiration.Valid {
			user.VerificationCodeExpiration = verificationCodeExpiration.String
		}
		if graduationYear.Valid {
			year := int(graduationYear.Int32)
			user.GraduationYear = &year
		}

		users = append(users, user)
	}

	rolesQuery := `
		SELECT nr.email, r.name
		FROM newf_roles nr
		JOIN roles r ON nr.id_roles = r.id_roles
	`
	rolesRows, err := h.DB.Query(rolesQuery)
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to fetch user roles")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
	} else {
		defer rolesRows.Close()
		for rolesRows.Next() {
			var email, roleName string
			if err := rolesRows.Scan(&email, &roleName); err == nil {
				userRolesMap[email] = append(userRolesMap[email], roleName)
			}
		}
	}

	for i := range users {
		if roles, exists := userRolesMap[users[i].Email]; exists {
			users[i].Roles = roles
		}
	}

	utils.LogMessage(utils.LevelInfo, "Successfully fetched all users")
	utils.LogLineKeyValue(utils.LevelInfo, "User Count", len(users))
	utils.LogFooter()

	return c.JSON(users)
}

func (h *AdminHandler) DeleteUser(c *fiber.Ctx) error {
	email := c.Params("email")
	utils.LogHeader("🗑️ Delete User (Admin)")
	utils.LogLineKeyValue(utils.LevelInfo, "Target Email", email)

	query := `DELETE FROM newf WHERE email = $1`
	result, err := h.DB.Exec(query, email)
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to delete user")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete user"})
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		utils.LogMessage(utils.LevelWarn, "User not found for deletion")
		utils.LogFooter()
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	utils.LogMessage(utils.LevelInfo, "User deleted successfully")
	utils.LogFooter()
	return c.SendStatus(fiber.StatusOK)
}

func (h *AdminHandler) CreateUser(c *fiber.Ctx) error {
	utils.LogHeader("👤 Create User (Admin)")

	var req struct {
		Email          string   `json:"email"`
		FirstName      string   `json:"first_name"`
		LastName       string   `json:"last_name"`
		PhoneNumber    string   `json:"phone_number"`
		Campus         string   `json:"campus"`
		FormationName  string   `json:"formation_name"`
		GraduationYear int      `json:"graduation_year"`
		Roles          []string `json:"roles"`
	}

	if err := c.BodyParser(&req); err != nil {
		utils.LogMessage(utils.LevelError, "Failed to parse request body")
		utils.LogFooter()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request format"})
	}

	tx, err := h.DB.Begin()
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to begin transaction")
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database error"})
	}
	defer tx.Rollback()

	defaultPassword := utils.GenerateRandomString(12)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(defaultPassword), bcrypt.DefaultCost)
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to hash password")
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Password generation error"})
	}

	insertQuery := `
		INSERT INTO newf (email, password, first_name, last_name`

	args := []interface{}{req.Email, string(hashedPassword), req.FirstName, req.LastName}
	argCount := 4

	if req.PhoneNumber != "" {
		insertQuery += `, phone_number`
		args = append(args, req.PhoneNumber)
		argCount++
	}

	if req.Campus != "" {
		insertQuery += `, campus`
		args = append(args, req.Campus)
		argCount++
	}

	if req.FormationName != "" {
		insertQuery += `, formation_name`
		args = append(args, req.FormationName)
		argCount++
	}

	if req.GraduationYear != 0 {
		insertQuery += `, graduation_year`
		args = append(args, req.GraduationYear)
		argCount++
	}

	insertQuery += `) VALUES ($1, $2, $3, $4`
	for i := 5; i <= argCount; i++ {
		insertQuery += fmt.Sprintf(`, $%d`, i)
	}
	insertQuery += `) RETURNING id_newf`

	var userID int
	err = tx.QueryRow(insertQuery, args...).Scan(&userID)
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to create user")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create user"})
	}

	for _, roleName := range req.Roles {
		var roleID int
		err := tx.QueryRow("SELECT id_roles FROM roles WHERE name = $1", roleName).Scan(&roleID)
		if err != nil {
			utils.LogMessage(utils.LevelWarn, "Role not found, skipping")
			utils.LogLineKeyValue(utils.LevelWarn, "Role", roleName)
			continue
		}

		_, err = tx.Exec("INSERT INTO newf_roles (email, id_roles) VALUES ($1, $2)", req.Email, roleID)
		if err != nil {
			utils.LogMessage(utils.LevelError, "Failed to assign role")
			utils.LogLineKeyValue(utils.LevelError, "Role", roleName)
		}
	}

	if err := tx.Commit(); err != nil {
		utils.LogMessage(utils.LevelError, "Failed to commit transaction")
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create user"})
	}

	utils.LogMessage(utils.LevelInfo, "User created successfully")
	utils.LogLineKeyValue(utils.LevelInfo, "Email", req.Email)
	utils.LogFooter()

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "User created successfully", "password": defaultPassword})
}

func (h *AdminHandler) UpdateUser(c *fiber.Ctx) error {
	email := c.Params("email")
	utils.LogHeader("✏️ Update User (Admin)")
	utils.LogLineKeyValue(utils.LevelInfo, "Target Email", email)

	var req struct {
		FirstName        *string  `json:"first_name"`
		LastName         *string  `json:"last_name"`
		PhoneNumber      *string  `json:"phone_number"`
		Campus           *string  `json:"campus"`
		FormationName    *string  `json:"formation_name"`
		GraduationYear   *int     `json:"graduation_year"`
		Roles            []string `json:"roles"`
		VerificationCode *string  `json:"verification_code"`
	}

	if err := c.BodyParser(&req); err != nil {
		utils.LogMessage(utils.LevelError, "Failed to parse request body")
		utils.LogFooter()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request format"})
	}

	tx, err := h.DB.Begin()
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to begin transaction")
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database error"})
	}
	defer tx.Rollback()

	// Construire dynamiquement la requête de mise à jour
	var updateFields []string
	var updateValues []interface{}
	paramCount := 1

	if req.FirstName != nil {
		updateFields = append(updateFields, "first_name = $"+fmt.Sprint(paramCount))
		updateValues = append(updateValues, *req.FirstName)
		paramCount++
	}

	if req.LastName != nil {
		updateFields = append(updateFields, "last_name = $"+fmt.Sprint(paramCount))
		updateValues = append(updateValues, *req.LastName)
		paramCount++
	}

	if req.PhoneNumber != nil {
		updateFields = append(updateFields, "phone_number = $"+fmt.Sprint(paramCount))
		updateValues = append(updateValues, *req.PhoneNumber)
		paramCount++
	}

	if req.Campus != nil {
		updateFields = append(updateFields, "campus = $"+fmt.Sprint(paramCount))
		updateValues = append(updateValues, *req.Campus)
		paramCount++
	}

	if req.FormationName != nil {
		updateFields = append(updateFields, "formation_name = $"+fmt.Sprint(paramCount))
		updateValues = append(updateValues, *req.FormationName)
		paramCount++
	}

	if req.GraduationYear != nil {
		if *req.GraduationYear == 0 {
			updateFields = append(updateFields, "graduation_year = NULL")
		} else {
			updateFields = append(updateFields, "graduation_year = $"+fmt.Sprint(paramCount))
			updateValues = append(updateValues, *req.GraduationYear)
			paramCount++
		}
	}

	if req.VerificationCode != nil {
		if *req.VerificationCode == "0" {
			updateFields = append(updateFields, "verification_code = NULL, verification_code_expiration = NULL")
		} else {
			updateFields = append(updateFields, "verification_code = $"+fmt.Sprint(paramCount))
			updateValues = append(updateValues, *req.VerificationCode)
			paramCount++
		}
	}

	// Ajouter l'email comme dernier paramètre
	updateValues = append(updateValues, email)

	// S'il n'y a pas de champs à mettre à jour, retourner une erreur
	if len(updateFields) == 0 {
		utils.LogMessage(utils.LevelWarn, "No fields to update")
		utils.LogFooter()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "No fields to update"})
	}

	updateQuery := fmt.Sprintf(`
		UPDATE newf 
		SET %s
		WHERE email = $%d
	`, strings.Join(updateFields, ", "), paramCount)

	result, err := tx.Exec(updateQuery, updateValues...)
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to update user")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update user"})
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		utils.LogMessage(utils.LevelWarn, "User not found for update")
		utils.LogFooter()
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	// Mettre à jour les rôles seulement si ils sont fournis
	if req.Roles != nil {
		_, err = tx.Exec("DELETE FROM newf_roles WHERE email = $1", email)
		if err != nil {
			utils.LogMessage(utils.LevelError, "Failed to clear existing roles")
			utils.LogLineKeyValue(utils.LevelError, "Error", err)
		}

		for _, roleName := range req.Roles {
			var roleID int
			err := tx.QueryRow("SELECT id_roles FROM roles WHERE name = $1", roleName).Scan(&roleID)
			if err != nil {
				utils.LogMessage(utils.LevelWarn, "Role not found, skipping")
				utils.LogLineKeyValue(utils.LevelWarn, "Role", roleName)
				continue
			}

			_, err = tx.Exec("INSERT INTO newf_roles (email, id_roles) VALUES ($1, $2)", email, roleID)
			if err != nil {
				utils.LogMessage(utils.LevelError, "Failed to assign role")
				utils.LogLineKeyValue(utils.LevelError, "Role", roleName)
			}
		}
	}

	if err := tx.Commit(); err != nil {
		utils.LogMessage(utils.LevelError, "Failed to commit transaction")
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update user"})
	}

	utils.LogMessage(utils.LevelInfo, "User updated successfully")
	utils.LogFooter()

	return c.SendStatus(fiber.StatusOK)
}

func (h *AdminHandler) GetAllEvents(c *fiber.Ctx) error {
	utils.LogHeader("🎉 Get All Events (Admin)")

	query := `
		SELECT 
			e.id_events,
			e.name,
			e.description,
			e.link,
			e.start_date,
			e.end_date,
			e.location,
			e.picture,
			e.creation_date,
			e.creator,
			e.id_club,
			c.name as club_name,
			COALESCE(attendee_count.count, 0) as attendee_count
		FROM events e
		LEFT JOIN clubs c ON e.id_club = c.id_clubs
		LEFT JOIN (
			SELECT id_events, COUNT(*) as count
			FROM events_attendents
			GROUP BY id_events
		) attendee_count ON e.id_events = attendee_count.id_events
		ORDER BY e.creation_date DESC
	`

	rows, err := h.DB.Query(query)
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to fetch events")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch events"})
	}
	defer rows.Close()

	var events []map[string]interface{}
	for rows.Next() {
		var event map[string]interface{} = make(map[string]interface{})
		var id, clubID, attendeeCount int
		var name, location, creator, clubName string
		var description, link, picture sql.NullString
		var startDate, endDate, creationDate sql.NullTime

		err := rows.Scan(
			&id, &name, &description, &link, &startDate, &endDate,
			&location, &picture, &creationDate, &creator, &clubID,
			&clubName, &attendeeCount,
		)
		if err != nil {
			utils.LogMessage(utils.LevelError, "Failed to scan event")
			utils.LogLineKeyValue(utils.LevelError, "Error", err)
			continue
		}

		event["id_events"] = id
		event["name"] = name
		event["description"] = description.String
		event["link"] = link.String
		event["location"] = location
		event["picture"] = picture.String
		event["creator"] = creator
		event["id_club"] = clubID
		event["club_name"] = clubName
		event["attendee_count"] = attendeeCount

		if startDate.Valid {
			event["start_date"] = startDate.Time
		}
		if endDate.Valid {
			event["end_date"] = endDate.Time
		}
		if creationDate.Valid {
			event["creation_date"] = creationDate.Time
		}

		events = append(events, event)
	}

	utils.LogMessage(utils.LevelInfo, "Successfully fetched all events")
	utils.LogLineKeyValue(utils.LevelInfo, "Event Count", len(events))
	utils.LogFooter()

	return c.JSON(events)
}

func (h *AdminHandler) CreateEvent(c *fiber.Ctx) error {
	utils.LogHeader("🎉 Create Event (Admin)")

	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Link        string `json:"link"`
		StartDate   string `json:"start_date"`
		EndDate     string `json:"end_date"`
		Location    string `json:"location"`
		Picture     string `json:"picture"`
		Creator     string `json:"creator"`
		ClubID      int    `json:"id_club"`
	}

	if err := c.BodyParser(&req); err != nil {
		utils.LogMessage(utils.LevelError, "Failed to parse request body")
		utils.LogFooter()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request format"})
	}

	insertQuery := `
		INSERT INTO events (name, description, link, start_date, end_date, location, picture, creator, id_club)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id_events
	`

	var eventID int
	err := h.DB.QueryRow(insertQuery, req.Name, req.Description, req.Link,
		req.StartDate, req.EndDate, req.Location, req.Picture, req.Creator, req.ClubID).Scan(&eventID)
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to create event")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create event"})
	}

	utils.LogMessage(utils.LevelInfo, "Event created successfully")
	utils.LogLineKeyValue(utils.LevelInfo, "Event ID", eventID)
	utils.LogFooter()

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Event created successfully", "id": eventID})
}

func (h *AdminHandler) UpdateEvent(c *fiber.Ctx) error {
	eventID := c.Params("id")
	utils.LogHeader("🎉 Update Event (Admin)")
	utils.LogLineKeyValue(utils.LevelInfo, "Event ID", eventID)

	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Link        string `json:"link"`
		StartDate   string `json:"start_date"`
		EndDate     string `json:"end_date"`
		Location    string `json:"location"`
		Picture     string `json:"picture"`
		Creator     string `json:"creator"`
		ClubID      int    `json:"id_club"`
	}

	if err := c.BodyParser(&req); err != nil {
		utils.LogMessage(utils.LevelError, "Failed to parse request body")
		utils.LogFooter()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request format"})
	}

	updateQuery := `
		UPDATE events 
		SET name = $1, description = $2, link = $3, start_date = $4, end_date = $5, 
		    location = $6, picture = $7, creator = $8, id_club = $9
		WHERE id_events = $10
	`

	result, err := h.DB.Exec(updateQuery, req.Name, req.Description, req.Link,
		req.StartDate, req.EndDate, req.Location, req.Picture, req.Creator, req.ClubID, eventID)
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to update event")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update event"})
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		utils.LogMessage(utils.LevelWarn, "Event not found for update")
		utils.LogFooter()
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Event not found"})
	}

	utils.LogMessage(utils.LevelInfo, "Event updated successfully")
	utils.LogFooter()

	return c.SendStatus(fiber.StatusOK)
}

func (h *AdminHandler) DeleteEvent(c *fiber.Ctx) error {
	eventID := c.Params("id")
	utils.LogHeader("🗑️ Delete Event (Admin)")
	utils.LogLineKeyValue(utils.LevelInfo, "Event ID", eventID)

	tx, err := h.DB.Begin()
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to begin transaction")
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database error"})
	}
	defer tx.Rollback()

	_, err = tx.Exec("DELETE FROM events_attendents WHERE id_events = $1", eventID)
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to delete event attendees")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete event"})
	}

	result, err := tx.Exec("DELETE FROM events WHERE id_events = $1", eventID)
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to delete event")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete event"})
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		utils.LogMessage(utils.LevelWarn, "Event not found for deletion")
		utils.LogFooter()
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Event not found"})
	}

	if err := tx.Commit(); err != nil {
		utils.LogMessage(utils.LevelError, "Failed to commit transaction")
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete event"})
	}

	utils.LogMessage(utils.LevelInfo, "Event deleted successfully")
	utils.LogFooter()
	return c.SendStatus(fiber.StatusOK)
}

func (h *AdminHandler) GetAllClubs(c *fiber.Ctx) error {
	utils.LogHeader("🏛️ Get All Clubs (Admin)")

	query := `
		SELECT 
			c.id_clubs,
			c.name,
			c.picture,
			c.description,
			c.location,
			c.link,
			COALESCE(member_count.count, 0) as member_count,
			respo.email as responsible_email,
			COALESCE(respo.first_name, '') as responsible_first_name,
			COALESCE(respo.last_name, '') as responsible_last_name
		FROM clubs c
		LEFT JOIN (
			SELECT id_clubs, COUNT(*) as count
			FROM clubs_members
			GROUP BY id_clubs
		) member_count ON c.id_clubs = member_count.id_clubs
		LEFT JOIN (
			SELECT DISTINCT c.id_clubs, n.email, n.first_name, n.last_name
			FROM clubs c
			JOIN newf_roles nr ON nr.email IN (
				SELECT email FROM clubs_members WHERE id_clubs = c.id_clubs
			)
			JOIN roles r ON nr.id_roles = r.id_roles
			JOIN newf n ON nr.email = n.email
			WHERE r.name LIKE CONCAT(LOWER(REPLACE(c.name, ' ', '')), '_respo')
		) respo ON c.id_clubs = respo.id_clubs
		ORDER BY c.name
	`

	rows, err := h.DB.Query(query)
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to fetch clubs")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch clubs"})
	}
	defer rows.Close()

	var clubs []map[string]interface{}
	for rows.Next() {
		var club map[string]interface{} = make(map[string]interface{})
		var id, memberCount int
		var name, picture string
		var description, location, link sql.NullString
		var responsibleEmail, responsibleFirstName, responsibleLastName sql.NullString

		err := rows.Scan(
			&id, &name, &picture, &description, &location, &link,
			&memberCount, &responsibleEmail, &responsibleFirstName, &responsibleLastName,
		)
		if err != nil {
			utils.LogMessage(utils.LevelError, "Failed to scan club")
			utils.LogLineKeyValue(utils.LevelError, "Error", err)
			continue
		}

		club["id_clubs"] = id
		club["name"] = name
		club["picture"] = picture
		club["description"] = description.String
		club["location"] = location.String
		club["link"] = link.String
		club["member_count"] = memberCount

		if responsibleEmail.Valid {
			club["responsible"] = map[string]interface{}{
				"email":      responsibleEmail.String,
				"first_name": responsibleFirstName.String,
				"last_name":  responsibleLastName.String,
			}
		}

		clubs = append(clubs, club)
	}

	utils.LogMessage(utils.LevelInfo, "Successfully fetched all clubs")
	utils.LogLineKeyValue(utils.LevelInfo, "Club Count", len(clubs))
	utils.LogFooter()

	return c.JSON(clubs)
}

func (h *AdminHandler) CreateClub(c *fiber.Ctx) error {
	utils.LogHeader("🏛️ Create Club (Admin)")

	var req struct {
		Name        string `json:"name"`
		Picture     string `json:"picture"`
		Description string `json:"description"`
		Location    string `json:"location"`
		Link        string `json:"link"`
	}

	if err := c.BodyParser(&req); err != nil {
		utils.LogMessage(utils.LevelError, "Failed to parse request body")
		utils.LogFooter()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request format"})
	}

	insertQuery := `
		INSERT INTO clubs (name, picture, description, location, link)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id_clubs
	`

	var clubID int
	err := h.DB.QueryRow(insertQuery, req.Name, req.Picture, req.Description, req.Location, req.Link).Scan(&clubID)
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to create club")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create club"})
	}

	utils.LogMessage(utils.LevelInfo, "Club created successfully")
	utils.LogLineKeyValue(utils.LevelInfo, "Club ID", clubID)
	utils.LogFooter()

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Club created successfully", "id": clubID})
}

func (h *AdminHandler) UpdateClub(c *fiber.Ctx) error {
	clubID := c.Params("id")
	utils.LogHeader("🏛️ Update Club (Admin)")
	utils.LogLineKeyValue(utils.LevelInfo, "Club ID", clubID)

	var req struct {
		Name        string `json:"name"`
		Picture     string `json:"picture"`
		Description string `json:"description"`
		Location    string `json:"location"`
		Link        string `json:"link"`
	}

	if err := c.BodyParser(&req); err != nil {
		utils.LogMessage(utils.LevelError, "Failed to parse request body")
		utils.LogFooter()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request format"})
	}

	updateQuery := `
		UPDATE clubs 
		SET name = $1, picture = $2, description = $3, location = $4, link = $5
		WHERE id_clubs = $6
	`

	result, err := h.DB.Exec(updateQuery, req.Name, req.Picture, req.Description, req.Location, req.Link, clubID)
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to update club")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update club"})
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		utils.LogMessage(utils.LevelWarn, "Club not found for update")
		utils.LogFooter()
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Club not found"})
	}

	utils.LogMessage(utils.LevelInfo, "Club updated successfully")
	utils.LogFooter()

	return c.SendStatus(fiber.StatusOK)
}

func (h *AdminHandler) DeleteClub(c *fiber.Ctx) error {
	clubID := c.Params("id")
	utils.LogHeader("🗑️ Delete Club (Admin)")
	utils.LogLineKeyValue(utils.LevelInfo, "Club ID", clubID)

	tx, err := h.DB.Begin()
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to begin transaction")
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database error"})
	}
	defer tx.Rollback()

	_, err = tx.Exec("DELETE FROM events WHERE id_club = $1", clubID)
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to delete club events")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete club"})
	}

	_, err = tx.Exec("DELETE FROM clubs_members WHERE id_clubs = $1", clubID)
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to delete club members")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete club"})
	}

	result, err := tx.Exec("DELETE FROM clubs WHERE id_clubs = $1", clubID)
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to delete club")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete club"})
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		utils.LogMessage(utils.LevelWarn, "Club not found for deletion")
		utils.LogFooter()
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Club not found"})
	}

	if err := tx.Commit(); err != nil {
		utils.LogMessage(utils.LevelError, "Failed to commit transaction")
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete club"})
	}

	utils.LogMessage(utils.LevelInfo, "Club deleted successfully")
	utils.LogFooter()
	return c.SendStatus(fiber.StatusOK)
}

func (h *AdminHandler) GetAllRoles(c *fiber.Ctx) error {
	utils.LogHeader("🏷️ Get All Roles (Admin)")

	query := `SELECT id_roles, name FROM roles ORDER BY name`
	rows, err := h.DB.Query(query)
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to fetch roles")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch roles"})
	}
	defer rows.Close()

	type Role struct {
		ID   int    `json:"id_roles"`
		Name string `json:"name"`
	}

	var roles []Role
	for rows.Next() {
		var role Role
		err := rows.Scan(&role.ID, &role.Name)
		if err != nil {
			utils.LogMessage(utils.LevelError, "Failed to scan role row")
			utils.LogLineKeyValue(utils.LevelError, "Error", err)
			continue
		}
		roles = append(roles, role)
	}

	utils.LogMessage(utils.LevelInfo, "Successfully fetched all roles")
	utils.LogLineKeyValue(utils.LevelInfo, "Role Count", len(roles))
	utils.LogFooter()

	return c.JSON(roles)
}

func (h *AdminHandler) ValidateUser(c *fiber.Ctx) error {
	email := c.Params("email")
	utils.LogHeader("✅ Validate User (Admin)")
	utils.LogLineKeyValue(utils.LevelInfo, "Target Email", email)

	tx, err := h.DB.Begin()
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to begin transaction")
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database error"})
	}
	defer tx.Rollback()

	// Clear verification codes
	_, err = tx.Exec("UPDATE newf SET verification_code = NULL, verification_code_expiration = NULL WHERE email = $1", email)
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to clear verification codes")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to validate user"})
	}

	// Get VERIFYING role ID
	var verifyingRoleID int
	err = tx.QueryRow("SELECT id_roles FROM roles WHERE name = 'VERIFYING'").Scan(&verifyingRoleID)
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to find VERIFYING role")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to validate user"})
	}

	// Get NEWF role ID
	var newfRoleID int
	err = tx.QueryRow("SELECT id_roles FROM roles WHERE name = 'NEWF'").Scan(&newfRoleID)
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to find NEWF role")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to validate user"})
	}

	// Remove VERIFYING role
	_, err = tx.Exec("DELETE FROM newf_roles WHERE email = $1 AND id_roles = $2", email, verifyingRoleID)
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to remove VERIFYING role")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to validate user"})
	}

	// Add NEWF role (if not already present)
	_, err = tx.Exec("INSERT INTO newf_roles (email, id_roles) VALUES ($1, $2) ON CONFLICT (email, id_roles) DO NOTHING", email, newfRoleID)
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to add NEWF role")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to validate user"})
	}

	if err := tx.Commit(); err != nil {
		utils.LogMessage(utils.LevelError, "Failed to commit transaction")
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to validate user"})
	}

	utils.LogMessage(utils.LevelInfo, "User validated successfully")
	utils.LogFooter()

	return c.SendStatus(fiber.StatusOK)
}
