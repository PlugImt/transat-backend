package club

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/plugimt/transat-backend/models"
	"github.com/plugimt/transat-backend/utils"
)

type ClubHandler struct {
	db *sql.DB
}

func NewclubHandler(db *sql.DB) *ClubHandler {
	return &ClubHandler{
		db: db,
	}
}

// GetClub returns all clubs with basic info, sorted by member count
func (h *ClubHandler) GetClub(c *fiber.Ctx) error {
	utils.LogHeader("🏛️ Get All Clubs")

	query := `
		SELECT 
		    c.id_clubs,
			c.name,
			c.description,
			c.picture,
			COALESCE(member_count.count, 0) as member_count
		FROM clubs c
		LEFT JOIN (
			SELECT id_clubs, COUNT(*) as count
			FROM clubs_members
			GROUP BY id_clubs
		) member_count ON c.id_clubs = member_count.id_clubs
		ORDER BY member_count DESC, c.name
	`

	rows, err := h.db.Query(query)
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to fetch clubs")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch clubs",
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

	var clubs []map[string]interface{}
	for rows.Next() {
		var id int
		var name, picture string
		var description sql.NullString
		var memberCount int

		err := rows.Scan(&id, &name, &description, &picture, &memberCount)
		if err != nil {
			utils.LogMessage(utils.LevelError, "Failed to scan club")
			utils.LogLineKeyValue(utils.LevelError, "Error", err)
			continue
		}

		club := map[string]interface{}{
			"id":           id,
			"name":         name,
			"description":  description.String,
			"picture":      picture,
			"member_count": memberCount,
		}

		clubs = append(clubs, club)
	}

	utils.LogMessage(utils.LevelInfo, "Successfully fetched clubs")
	utils.LogLineKeyValue(utils.LevelInfo, "Count", len(clubs))
	utils.LogFooter()

	return c.JSON(clubs)
}

// GetClubByID returns a specific club by ID with detailed information
func (h *ClubHandler) GetClubByID(c *fiber.Ctx) error {
	utils.LogHeader("🏛️ Get Club By ID")

	clubID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		utils.LogMessage(utils.LevelWarn, "Invalid club ID")
		utils.LogFooter()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid club ID",
		})
	}

	utils.LogLineKeyValue(utils.LevelInfo, "Club ID", clubID)

	// Get club basic info
	clubQuery := `
		SELECT 
			c.id_clubs,
			c.name,
			c.picture,
			c.description,
			c.location,
			c.link
		FROM clubs c
		WHERE c.id_clubs = $1
	`

	var club models.Club
	var description, location, link sql.NullString

	err = h.db.QueryRow(clubQuery, clubID).Scan(
		&club.ID,
		&club.Name,
		&club.Picture,
		&description,
		&location,
		&link,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			utils.LogMessage(utils.LevelWarn, "Club not found")
			utils.LogFooter()
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Club not found",
			})
		}
		utils.LogMessage(utils.LevelError, "Failed to fetch club")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch club",
		})
	}

	club.Description = description.String
	club.Location = location.String
	club.Link = link.String

	// Get responsible info
	roleName := fmt.Sprintf("%s_respo", strings.ToLower(strings.ReplaceAll(club.Name, " ", "")))
	respoQuery := `
		SELECT 
			n.email,
			n.first_name,
			n.last_name,
			n.profile_picture,
			n.graduation_year
		FROM newf_roles nr
		JOIN roles r ON nr.id_roles = r.id_roles
		JOIN newf n ON nr.email = n.email
		WHERE r.name = $1
		LIMIT 1
	`

	var responsible map[string]interface{}
	var respoEmail, respoFirstName, respoLastName string
	var respoProfilePicture sql.NullString
	var respoGraduationYear sql.NullInt64

	err = h.db.QueryRow(respoQuery, roleName).Scan(
		&respoEmail,
		&respoFirstName,
		&respoLastName,
		&respoProfilePicture,
		&respoGraduationYear,
	)

	if err == nil {
		responsible = map[string]interface{}{
			"email":           respoEmail,
			"first_name":      respoFirstName,
			"last_name":       respoLastName,
			"profile_picture": respoProfilePicture.String,
			"graduation_year": respoGraduationYear.Int64,
		}
	} else if err != sql.ErrNoRows {
		utils.LogMessage(utils.LevelError, "Failed to fetch responsible")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
	}

	// Get total member count
	var memberCount int
	err = h.db.QueryRow("SELECT COUNT(*) FROM clubs_members WHERE id_clubs = $1", clubID).Scan(&memberCount)
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to get member count")
		memberCount = 0
	}

	// Get photos of first 5 members
	membersQuery := `
		SELECT n.profile_picture
		FROM clubs_members cm
		JOIN newf n ON cm.email = n.email
		WHERE cm.id_clubs = $1 AND n.profile_picture IS NOT NULL AND n.profile_picture != ''
		ORDER BY n.first_name, n.last_name
		LIMIT 5
	`

	rows, err := h.db.Query(membersQuery, clubID)
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to fetch member photos")
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

	var memberPhotos []string
	if rows != nil {
		for rows.Next() {
			var photo string
			if err := rows.Scan(&photo); err == nil {
				memberPhotos = append(memberPhotos, photo)
			}
		}
	}

	// Build response
	response := map[string]interface{}{
		"id":            club.ID,
		"name":          club.Name,
		"picture":       club.Picture,
		"description":   club.Description,
		"location":      club.Location,
		"link":          club.Link,
		"member_count":  memberCount,
		"member_photos": memberPhotos,
	}

	if responsible != nil {
		response["responsible"] = responsible
	}

	utils.LogMessage(utils.LevelInfo, "Successfully fetched club details")
	utils.LogFooter()

	return c.JSON(response)
}

// GetClubMembers returns all members of a specific club with detailed info
func (h *ClubHandler) GetClubMembers(c *fiber.Ctx) error {
	utils.LogHeader("👥 Get Club Members")

	clubID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		utils.LogMessage(utils.LevelWarn, "Invalid club ID")
		utils.LogFooter()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid club ID",
		})
	}

	utils.LogLineKeyValue(utils.LevelInfo, "Club ID", clubID)

	// First check if club exists
	var clubExists bool
	err = h.db.QueryRow("SELECT EXISTS(SELECT 1 FROM clubs WHERE id_clubs = $1)", clubID).Scan(&clubExists)
	if err != nil || !clubExists {
		utils.LogMessage(utils.LevelWarn, "Club not found")
		utils.LogFooter()
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Club not found",
		})
	}

	// Get club name for role checking
	var clubName string
	err = h.db.QueryRow("SELECT name FROM clubs WHERE id_clubs = $1", clubID).Scan(&clubName)
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to get club name")
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get club information",
		})
	}

	query := `
		SELECT 
    	cm.email,
    	n.first_name,
    	n.last_name,
    	n.profile_picture,
    	n.graduation_year,
    	EXISTS (
    	    SELECT 1
    	    FROM newf_roles nr
    	    JOIN roles r ON nr.id_roles = r.id_roles
    	    WHERE nr.email = cm.email AND r.name = $1
    	) AS is_respo
		FROM clubs_members cm
		JOIN newf n ON cm.email = n.email
		WHERE cm.id_clubs = $2
		ORDER BY is_respo DESC, n.first_name, n.last_name;
	`

	roleName := fmt.Sprintf("%s_respo", strings.ToLower(strings.ReplaceAll(clubName, " ", "")))

	rows, err := h.db.Query(query, roleName, clubID)
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to fetch club members")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch club members",
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
		var isRespo bool

		err := rows.Scan(&email, &firstName, &lastName, &profilePicture, &graduationYear, &isRespo)
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
			"is_respo":        isRespo,
		}

		members = append(members, member)
	}

	utils.LogMessage(utils.LevelInfo, "Successfully fetched club members")
	utils.LogLineKeyValue(utils.LevelInfo, "Count", len(members))
	utils.LogFooter()

	return c.JSON(fiber.Map{
		"members": members,
		"count":   len(members),
	})
}

// CreateClub creates a new club and assigns the creator as responsible
func (h *ClubHandler) CreateClub(c *fiber.Ctx) error {
	utils.LogHeader("🏛️ Create Club")

	userEmail := c.Locals("email").(string)
	utils.LogLineKeyValue(utils.LevelInfo, "Creator", userEmail)

	var req models.CreateClubRequest
	if err := c.BodyParser(&req); err != nil {
		utils.LogMessage(utils.LevelError, "Failed to parse request body")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate required fields
	if req.Name == "" || req.Picture == "" {
		utils.LogMessage(utils.LevelWarn, "Missing required fields")
		utils.LogFooter()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Name and picture are required",
		})
	}

	utils.LogLineKeyValue(utils.LevelInfo, "Club Name", req.Name)

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

	// Create club
	var clubID int
	insertClubQuery := `
		INSERT INTO clubs (name, picture, description, location, link)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id_clubs
	`
	err = tx.QueryRow(
		insertClubQuery,
		req.Name,
		req.Picture,
		req.Description,
		req.Location,
		req.Link,
	).Scan(&clubID)

	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to create club")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create club",
		})
	}

	// Add creator as member
	_, err = tx.Exec(
		"INSERT INTO clubs_members (email, id_clubs) VALUES ($1, $2)",
		userEmail, clubID,
	)
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to add creator as member")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to add creator as member",
		})
	}

	// Create role for club responsible
	roleName := fmt.Sprintf("%s_respo", strings.ToLower(strings.ReplaceAll(req.Name, " ", "")))
	roleDescription := fmt.Sprintf("Responsible for the %s club", req.Name)

	var roleID int
	_, err = tx.Exec(
		"INSERT INTO roles (name, description) VALUES ($1, $2)",
		roleName, roleDescription,
	)
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to create role")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create club role",
		})
	}

	// Get the role ID
	err = tx.QueryRow("SELECT id_roles FROM roles WHERE name = $1", roleName).Scan(&roleID)
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to get role ID")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get role information",
		})
	}

	// Assign creator as responsible
	_, err = tx.Exec(
		"INSERT INTO newf_roles (email, id_roles) VALUES ($1, $2)",
		userEmail, roleID,
	)
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to assign creator as responsible")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to assign creator as responsible",
		})
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		utils.LogMessage(utils.LevelError, "Failed to commit transaction")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save club",
		})
	}

	utils.LogMessage(utils.LevelInfo, "Club created successfully")
	utils.LogLineKeyValue(utils.LevelInfo, "Club ID", clubID)
	utils.LogLineKeyValue(utils.LevelInfo, "Role", roleName)
	utils.LogFooter()

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Club created successfully",
		"club": map[string]interface{}{
			"id_clubs":    clubID,
			"name":        req.Name,
			"picture":     req.Picture,
			"description": req.Description,
			"location":    req.Location,
			"link":        req.Link,
		},
	})
}

// UpdateClub updates club information (only responsibles can update)
func (h *ClubHandler) UpdateClub(c *fiber.Ctx) error {
	utils.LogHeader("🏛️ Update Club")

	clubID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		utils.LogMessage(utils.LevelWarn, "Invalid club ID")
		utils.LogFooter()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid club ID",
		})
	}

	userEmail := c.Locals("email").(string)
	utils.LogLineKeyValue(utils.LevelInfo, "Club ID", clubID)
	utils.LogLineKeyValue(utils.LevelInfo, "User", userEmail)

	// Check if user is responsible for this club
	var clubName string
	err = h.db.QueryRow("SELECT name FROM clubs WHERE id_clubs = $1", clubID).Scan(&clubName)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.LogMessage(utils.LevelWarn, "Club not found")
			utils.LogFooter()
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Club not found",
			})
		}
		utils.LogMessage(utils.LevelError, "Failed to get club")
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get club information",
		})
	}

	// Check if user has the club responsible role
	roleName := fmt.Sprintf("%s_respo", strings.ToLower(strings.ReplaceAll(clubName, " ", "")))
	var hasRole bool
	err = h.db.QueryRow(`
		SELECT EXISTS(
			SELECT 1 FROM newf_roles nr
			JOIN roles r ON nr.id_roles = r.id_roles
			WHERE nr.email = $1 AND r.name = $2
		)
	`, userEmail, roleName).Scan(&hasRole)

	if err != nil || !hasRole {
		utils.LogMessage(utils.LevelWarn, "User not authorized to update club")
		utils.LogFooter()
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Only club responsibles can update club information",
		})
	}

	var req models.UpdateClubRequest
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
	if req.Picture != "" {
		setParts = append(setParts, fmt.Sprintf("picture = $%d", argIndex))
		args = append(args, req.Picture)
		argIndex++
	}
	if req.Description != "" {
		setParts = append(setParts, fmt.Sprintf("description = $%d", argIndex))
		args = append(args, req.Description)
		argIndex++
	}
	if req.Location != "" {
		setParts = append(setParts, fmt.Sprintf("location = $%d", argIndex))
		args = append(args, req.Location)
		argIndex++
	}
	if req.Link != "" {
		setParts = append(setParts, fmt.Sprintf("link = $%d", argIndex))
		args = append(args, req.Link)
		argIndex++
	}

	if len(setParts) == 0 {
		utils.LogMessage(utils.LevelWarn, "No fields to update")
		utils.LogFooter()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No fields to update",
		})
	}

	query := fmt.Sprintf("UPDATE clubs SET %s WHERE id_clubs = $%d", strings.Join(setParts, ", "), argIndex)
	args = append(args, clubID)

	_, err = h.db.Exec(query, args...)
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to update club")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update club",
		})
	}

	utils.LogMessage(utils.LevelInfo, "Club updated successfully")
	utils.LogFooter()

	return c.JSON(fiber.Map{
		"message": "Club updated successfully",
	})
}

// AddClubRespo assigns a new responsible to a club (removes old one)
func (h *ClubHandler) AddClubRespo(c *fiber.Ctx) error {
	utils.LogHeader("👑 Add Club Responsible")

	clubID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		utils.LogMessage(utils.LevelWarn, "Invalid club ID")
		utils.LogFooter()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid club ID",
		})
	}

	userEmail := c.Locals("email").(string)
	utils.LogLineKeyValue(utils.LevelInfo, "Club ID", clubID)
	utils.LogLineKeyValue(utils.LevelInfo, "Current User", userEmail)

	// Get club name first
	var clubName string
	err = h.db.QueryRow("SELECT name FROM clubs WHERE id_clubs = $1", clubID).Scan(&clubName)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.LogMessage(utils.LevelWarn, "Club not found")
			utils.LogFooter()
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Club not found",
			})
		}
		utils.LogMessage(utils.LevelError, "Failed to get club")
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get club information",
		})
	}

	// Check permissions directly from database - more secure!
	roleName := fmt.Sprintf("%s_respo", strings.ToLower(strings.ReplaceAll(clubName, " ", "")))
	utils.LogLineKeyValue(utils.LevelInfo, "Expected Role Name", roleName)

	// Check if user has permission (ADMIN or current responsible) in one efficient query
	var hasPermission bool
	permissionQuery := `
		SELECT EXISTS(
			SELECT 1 FROM newf n
			JOIN newf_roles nr ON n.email = nr.email
			JOIN roles r ON nr.id_roles = r.id_roles
			WHERE n.email = $1 AND (r.name = 'ADMIN' OR r.name = $2)
		)`
	err = h.db.QueryRow(permissionQuery, userEmail, roleName).Scan(&hasPermission)
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to check permissions")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to check permissions",
		})
	}

	utils.LogLineKeyValue(utils.LevelInfo, "Has Permission (ADMIN or Respo)", hasPermission)

	if !hasPermission {
		utils.LogMessage(utils.LevelWarn, "User not authorized to change responsible")
		utils.LogLineKeyValue(utils.LevelWarn, "Has Permission", hasPermission)
		utils.LogLineKeyValue(utils.LevelWarn, "Expected Role", roleName)
		utils.LogFooter()
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Only admins or current responsibles can change club responsible",
		})
	}

	// Debug: Log the raw body content
	bodyBytes := c.Body()
	utils.LogLineKeyValue(utils.LevelInfo, "Raw Body", string(bodyBytes))
	utils.LogLineKeyValue(utils.LevelInfo, "Content-Type", c.Get("Content-Type"))

	var req models.AddRespoRequest
	if err := c.BodyParser(&req); err != nil {
		utils.LogMessage(utils.LevelError, "Failed to parse request body")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogLineKeyValue(utils.LevelError, "Raw Body", string(bodyBytes))
		utils.LogLineKeyValue(utils.LevelError, "Content-Type", c.Get("Content-Type"))
		utils.LogFooter()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	utils.LogLineKeyValue(utils.LevelInfo, "Parsed Request", req)

	if req.Email == "" {
		utils.LogMessage(utils.LevelWarn, "Email is required")
		utils.LogFooter()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Email is required",
		})
	}

	utils.LogLineKeyValue(utils.LevelInfo, "New Responsible", req.Email)

	// Check if new responsible exists and is a club member
	var userExists, isMember bool
	err = h.db.QueryRow("SELECT EXISTS(SELECT 1 FROM newf WHERE email = $1)", req.Email).Scan(&userExists)
	if err != nil || !userExists {
		utils.LogMessage(utils.LevelWarn, "User not found")
		utils.LogFooter()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	err = h.db.QueryRow("SELECT EXISTS(SELECT 1 FROM clubs_members WHERE email = $1 AND id_clubs = $2)", req.Email, clubID).Scan(&isMember)
	if err != nil || !isMember {
		utils.LogMessage(utils.LevelWarn, "User is not a club member")
		utils.LogFooter()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "User must be a club member to become responsible",
		})
	}

	// Start transaction
	tx, err := h.db.Begin()
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to begin transaction")
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database transaction error",
		})
	}
	defer tx.Rollback()

	// Get role ID
	var roleID int
	err = tx.QueryRow("SELECT id_roles FROM roles WHERE name = $1", roleName).Scan(&roleID)
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to get role ID")
		utils.LogLineKeyValue(utils.LevelError, "Role Name", roleName)
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get role information",
		})
	}

	// Remove old responsible (if any)
	_, err = tx.Exec("DELETE FROM newf_roles WHERE id_roles = $1", roleID)
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to remove old responsible")
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to remove old responsible",
		})
	}

	// Add new responsible
	_, err = tx.Exec("INSERT INTO newf_roles (email, id_roles) VALUES ($1, $2)", req.Email, roleID)
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to add new responsible")
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to add new responsible",
		})
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		utils.LogMessage(utils.LevelError, "Failed to commit transaction")
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save changes",
		})
	}

	utils.LogMessage(utils.LevelInfo, "Club responsible updated successfully")
	utils.LogFooter()

	return c.JSON(fiber.Map{
		"message":         "Club responsible updated successfully",
		"new_responsible": req.Email,
	})
}

// JoinClub adds the user to a club
func (h *ClubHandler) JoinClub(c *fiber.Ctx) error {
	utils.LogHeader("🤝 Join Club")

	clubID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		utils.LogMessage(utils.LevelWarn, "Invalid club ID")
		utils.LogFooter()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid club ID",
		})
	}

	userEmail := c.Locals("email").(string)
	utils.LogLineKeyValue(utils.LevelInfo, "Club ID", clubID)
	utils.LogLineKeyValue(utils.LevelInfo, "User", userEmail)

	// Check if club exists
	var clubExists bool
	err = h.db.QueryRow("SELECT EXISTS(SELECT 1 FROM clubs WHERE id_clubs = $1)", clubID).Scan(&clubExists)
	if err != nil || !clubExists {
		utils.LogMessage(utils.LevelWarn, "Club not found")
		utils.LogFooter()
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Club not found",
		})
	}

	// Check if user is already a member
	var isMember bool
	err = h.db.QueryRow("SELECT EXISTS(SELECT 1 FROM clubs_members WHERE email = $1 AND id_clubs = $2)", userEmail, clubID).Scan(&isMember)
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to check membership")
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to check membership status",
		})
	}

	if isMember {
		utils.LogMessage(utils.LevelWarn, "User already a member")
		utils.LogFooter()
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "You are already a member of this club",
		})
	}

	// Add user to club
	_, err = h.db.Exec("INSERT INTO clubs_members (email, id_clubs) VALUES ($1, $2)", userEmail, clubID)
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to join club")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to join club",
		})
	}

	utils.LogMessage(utils.LevelInfo, "Successfully joined club")
	utils.LogFooter()

	return c.JSON(fiber.Map{
		"message": "Successfully joined the club",
	})
}

// LeaveClub removes the user from a club
func (h *ClubHandler) LeaveClub(c *fiber.Ctx) error {
	utils.LogHeader("👋 Leave Club")

	clubID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		utils.LogMessage(utils.LevelWarn, "Invalid club ID")
		utils.LogFooter()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid club ID",
		})
	}

	userEmail := c.Locals("email").(string)
	utils.LogLineKeyValue(utils.LevelInfo, "Club ID", clubID)
	utils.LogLineKeyValue(utils.LevelInfo, "User", userEmail)

	// Check if user is a member
	var isMember bool
	err = h.db.QueryRow("SELECT EXISTS(SELECT 1 FROM clubs_members WHERE email = $1 AND id_clubs = $2)", userEmail, clubID).Scan(&isMember)
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to check membership")
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to check membership status",
		})
	}

	if !isMember {
		utils.LogMessage(utils.LevelWarn, "User not a member")
		utils.LogFooter()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "You are not a member of this club",
		})
	}

	// Check if user is responsible for this club
	var clubName string
	err = h.db.QueryRow("SELECT name FROM clubs WHERE id_clubs = $1", clubID).Scan(&clubName)
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to get club name")
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get club information",
		})
	}

	roleName := fmt.Sprintf("%s_respo", strings.ToLower(strings.ReplaceAll(clubName, " ", "")))
	var isRespo bool
	err = h.db.QueryRow(`
		SELECT EXISTS(
			SELECT 1 FROM newf_roles nr
			JOIN roles r ON nr.id_roles = r.id_roles
			WHERE nr.email = $1 AND r.name = $2
		)
	`, userEmail, roleName).Scan(&isRespo)

	if err == nil && isRespo {
		utils.LogMessage(utils.LevelWarn, "Responsible cannot leave club")
		utils.LogFooter()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Club responsibles cannot leave the club. Please assign a new responsible first.",
		})
	}

	// Remove user from club
	result, err := h.db.Exec("DELETE FROM clubs_members WHERE email = $1 AND id_clubs = $2", userEmail, clubID)
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to leave club")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to leave club",
		})
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		utils.LogMessage(utils.LevelWarn, "No rows affected when leaving club")
		utils.LogFooter()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "You are not a member of this club",
		})
	}

	utils.LogMessage(utils.LevelInfo, "Successfully left club")
	utils.LogFooter()

	return c.JSON(fiber.Map{
		"message": "Successfully left the club",
	})
}
