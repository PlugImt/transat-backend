package association

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/plugimt/transat-backend/models"
	"github.com/plugimt/transat-backend/utils"
)

type AssociationHandler struct {
	db *sql.DB
}

func NewAssociationHandler(db *sql.DB) *AssociationHandler {
	return &AssociationHandler{
		db: db,
	}
}

// GetAssociation returns all associations with basic info, sorted by member count
func (h *AssociationHandler) GetAssociation(c *fiber.Ctx) error {
	utils.LogHeader("🏛️ Get All Associations")

	query := `
		SELECT 
		    c.id_associations,
			c.name,
			c.description,
			c.picture,
			COALESCE(member_count.count, 0) as member_count
		FROM associations c
		LEFT JOIN (
			SELECT id_associations, COUNT(*) as count
			FROM associations_members
			GROUP BY id_associations
		) member_count ON c.id_associations = member_count.id_associations
		ORDER BY member_count DESC, c.name
	`

	rows, err := h.db.Query(query)
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to fetch associations")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch associations",
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

	var associations []map[string]interface{}
	for rows.Next() {
		var id int
		var name, picture string
		var description sql.NullString
		var memberCount int

		err := rows.Scan(&id, &name, &description, &picture, &memberCount)
		if err != nil {
			utils.LogMessage(utils.LevelError, "Failed to scan association")
			utils.LogLineKeyValue(utils.LevelError, "Error", err)
			continue
		}

		association := map[string]interface{}{
			"id":           id,
			"name":         name,
			"description":  description.String,
			"picture":      picture,
			"member_count": memberCount,
		}

		associations = append(associations, association)
	}

	utils.LogMessage(utils.LevelInfo, "Successfully fetched associations")
	utils.LogLineKeyValue(utils.LevelInfo, "Count", len(associations))
	utils.LogFooter()

	return c.JSON(associations)
}

func (h *AssociationHandler) GetAssociationByID(c *fiber.Ctx) error {
	utils.LogHeader("🏛️ Get Association By ID")

	associationID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		utils.LogMessage(utils.LevelWarn, "Invalid association ID")
		utils.LogFooter()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid association ID",
		})
	}

	utils.LogLineKeyValue(utils.LevelInfo, "Association ID", associationID)

	// Get association basic info
	associationQuery := `
       SELECT 
          c.id_associations,
          c.name,
          c.picture,
          c.description,
          c.location,
          c.link
       FROM associations c
       WHERE c.id_associations = $1
    `

	var association models.Association
	var description, location, link sql.NullString

	err = h.db.QueryRow(associationQuery, associationID).Scan(
		&association.ID,
		&association.Name,
		&association.Picture,
		&description,
		&location,
		&link,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			utils.LogMessage(utils.LevelWarn, "Association not found")
			utils.LogFooter()
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Association not found",
			})
		}
		utils.LogMessage(utils.LevelError, "Failed to fetch association")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch association",
		})
	}

	association.Description = description.String
	association.Location = location.String
	association.Link = link.String

	// Get responsible info (Fetch ALL managers via db.Query)
	roleName := fmt.Sprintf("%s_respo", strings.ToLower(strings.ReplaceAll(association.Name, " ", "")))
	respoQuery := `
       SELECT 
          n.email,
          n.first_name,
          n.last_name,
          COALESCE(n.profile_picture, ''),
          COALESCE(n.graduation_year, 0)
       FROM newf_roles nr
       JOIN roles r ON nr.id_roles = r.id_roles
       JOIN newf n ON nr.email = n.email
       WHERE r.name = $1
    `

	respoRows, err := h.db.Query(respoQuery, roleName)
	var responsibleList []map[string]interface{}

	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to fetch association responsibles")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
	} else {
		defer respoRows.Close()
		for respoRows.Next() {
			var email, firstName, lastName, profilePicture string
			var graduationYear int

			err := respoRows.Scan(&email, &firstName, &lastName, &profilePicture, &graduationYear)
			if err != nil {
				utils.LogMessage(utils.LevelError, "Failed to scan responsible row")
				continue
			}

			responsible := map[string]interface{}{
				"email":           email,
				"first_name":      firstName,
				"last_name":       lastName,
				"profile_picture": profilePicture,
				"graduation_year": graduationYear,
			}
			responsibleList = append(responsibleList, responsible)
		}
	}

	// Get total member count
	var memberCount int
	err = h.db.QueryRow("SELECT COUNT(*) FROM associations_members WHERE id_associations = $1", associationID).Scan(&memberCount)
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to get member count")
		memberCount = 0
	}

	// Get photos of first 5 members
	membersQuery := `
       SELECT COALESCE(n.profile_picture, '')
       FROM associations_members cm
       JOIN newf n ON cm.email = n.email
       WHERE cm.id_associations = $1 AND n.profile_picture IS NOT NULL AND n.profile_picture != ''
       ORDER BY n.first_name, n.last_name
       LIMIT 5
    `

	rows, err := h.db.Query(membersQuery, associationID)
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

	isMemberQuery := `
       SELECT EXISTS(SELECT 1 FROM associations_members WHERE email = $1 AND id_associations = $2)
    `

	var hasJoined bool
	err = h.db.QueryRow(isMemberQuery, c.Locals("email").(string), associationID).Scan(&hasJoined)
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to check if user has joined association")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		hasJoined = false
	}

	// Build response
	response := map[string]interface{}{
		"id":            association.ID,
		"name":          association.Name,
		"picture":       association.Picture,
		"description":   association.Description,
		"location":      association.Location,
		"link":          association.Link,
		"member_count":  memberCount,
		"member_photos": memberPhotos,
		"has_joined":    hasJoined,
		"responsible":   responsibleList,
	}

	utils.LogMessage(utils.LevelInfo, "Successfully fetched association details")
	utils.LogFooter()

	return c.JSON(response)
}

// GetAssociationMembers returns all members of a specific association with detailed info
func (h *AssociationHandler) GetAssociationMembers(c *fiber.Ctx) error {
	utils.LogHeader("👥 Get Association Members")

	associationID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		utils.LogMessage(utils.LevelWarn, "Invalid association ID")
		utils.LogFooter()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid association ID",
		})
	}

	utils.LogLineKeyValue(utils.LevelInfo, "Association ID", associationID)

	// First check if association exists
	var associationExists bool
	err = h.db.QueryRow("SELECT EXISTS(SELECT 1 FROM associations WHERE id_associations = $1)", associationID).Scan(&associationExists)
	if err != nil || !associationExists {
		utils.LogMessage(utils.LevelWarn, "Association not found")
		utils.LogFooter()
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Association not found",
		})
	}

	// Get association name for role checking
	var associationName string
	err = h.db.QueryRow("SELECT name FROM associations WHERE id_associations = $1", associationID).Scan(&associationName)
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to get association name")
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get association information",
		})
	}

	query := `
		SELECT 
    	cm.email,
    	n.first_name,
    	n.last_name,
    	COALESCE(n.profile_picture, ''),
    	n.graduation_year,
    	EXISTS (
    	    SELECT 1
    	    FROM newf_roles nr
    	    JOIN roles r ON nr.id_roles = r.id_roles
    	    WHERE nr.email = cm.email AND r.name = $1
    	) AS is_respo
		FROM associations_members cm
		JOIN newf n ON cm.email = n.email
		WHERE cm.id_associations = $2
		ORDER BY is_respo DESC, n.first_name, n.last_name;
	`

	roleName := fmt.Sprintf("%s_respo", strings.ToLower(strings.ReplaceAll(associationName, " ", "")))

	rows, err := h.db.Query(query, roleName, associationID)
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to fetch association members")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch association members",
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

	utils.LogMessage(utils.LevelInfo, "Successfully fetched association members")
	utils.LogLineKeyValue(utils.LevelInfo, "Count", len(members))
	utils.LogFooter()

	return c.JSON(fiber.Map{
		"members": members,
		"count":   len(members),
	})
}

// CreateAssociation creates a new association and assigns the creator as responsible
func (h *AssociationHandler) CreateAssociation(c *fiber.Ctx) error {
	utils.LogHeader("🏛️ Create Association")

	userEmail := c.Locals("email").(string)
	utils.LogLineKeyValue(utils.LevelInfo, "Creator", userEmail)

	var req models.CreateAssociationRequest
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

	utils.LogLineKeyValue(utils.LevelInfo, "Association Name", req.Name)

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

	// Create association
	var associationID int
	insertAssociationQuery := `
		INSERT INTO associations (name, picture, description, location, link)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id_associations
	`
	err = tx.QueryRow(
		insertAssociationQuery,
		req.Name,
		req.Picture,
		req.Description,
		req.Location,
		req.Link,
	).Scan(&associationID)

	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to create association")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create association",
		})
	}

	// Add creator as member
	_, err = tx.Exec(
		"INSERT INTO associations_members (email, id_associations) VALUES ($1, $2)",
		userEmail, associationID,
	)
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to add creator as member")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to add creator as member",
		})
	}

	// Create role for association responsible
	roleName := fmt.Sprintf("%s_respo", strings.ToLower(strings.ReplaceAll(req.Name, " ", "")))
	roleDescription := fmt.Sprintf("Responsible for the %s association", req.Name)

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
			"error": "Failed to create association role",
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
			"error": "Failed to save association",
		})
	}

	utils.LogMessage(utils.LevelInfo, "Association created successfully")
	utils.LogLineKeyValue(utils.LevelInfo, "Association ID", associationID)
	utils.LogLineKeyValue(utils.LevelInfo, "Role", roleName)
	utils.LogFooter()

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Association created successfully",
		"association": map[string]interface{}{
			"id_associations": associationID,
			"name":            req.Name,
			"picture":         req.Picture,
			"description":     req.Description,
			"location":        req.Location,
			"link":            req.Link,
		},
	})
}

// UpdateAssociation updates association information (only responsibles can update)
func (h *AssociationHandler) UpdateAssociation(c *fiber.Ctx) error {
	utils.LogHeader("🏛️ Update Association")

	associationID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		utils.LogMessage(utils.LevelWarn, "Invalid association ID")
		utils.LogFooter()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid association ID",
		})
	}

	userEmail := c.Locals("email").(string)
	utils.LogLineKeyValue(utils.LevelInfo, "Association ID", associationID)
	utils.LogLineKeyValue(utils.LevelInfo, "User", userEmail)

	// Check if user is responsible for this association
	var associationName string
	err = h.db.QueryRow("SELECT name FROM associations WHERE id_associations = $1", associationID).Scan(&associationName)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.LogMessage(utils.LevelWarn, "Association not found")
			utils.LogFooter()
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Association not found",
			})
		}
		utils.LogMessage(utils.LevelError, "Failed to get association")
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get association information",
		})
	}

	// Check if user has the association responsible role
	roleName := fmt.Sprintf("%s_respo", strings.ToLower(strings.ReplaceAll(associationName, " ", "")))
	var hasRole bool
	err = h.db.QueryRow(`
		SELECT EXISTS(
			SELECT 1 FROM newf_roles nr
			JOIN roles r ON nr.id_roles = r.id_roles
			WHERE nr.email = $1 AND r.name = $2
		)
	`, userEmail, roleName).Scan(&hasRole)

	if err != nil || !hasRole {
		utils.LogMessage(utils.LevelWarn, "User not authorized to update association")
		utils.LogFooter()
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Only association responsibles can update association information",
		})
	}

	var req models.UpdateAssociationRequest
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

	query := fmt.Sprintf("UPDATE associations SET %s WHERE id_associations = $%d", strings.Join(setParts, ", "), argIndex)
	args = append(args, associationID)

	_, err = h.db.Exec(query, args...)
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to update association")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update association",
		})
	}

	utils.LogMessage(utils.LevelInfo, "Association updated successfully")
	utils.LogFooter()

	return c.JSON(fiber.Map{
		"message": "Association updated successfully",
	})
}

// AddAssociationRespo adds a co-responsible without wiping existing ones
func (h *AssociationHandler) AddAssociationRespo(c *fiber.Ctx) error {
	utils.LogHeader("👑 Add Association Responsible")

	associationID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		utils.LogMessage(utils.LevelWarn, "Invalid association ID")
		utils.LogFooter()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid association ID",
		})
	}

	userEmail := c.Locals("email").(string)
	utils.LogLineKeyValue(utils.LevelInfo, "Association ID", associationID)

	var associationName string
	err = h.db.QueryRow("SELECT name FROM associations WHERE id_associations = $1", associationID).Scan(&associationName)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.LogMessage(utils.LevelWarn, "Association not found")
			utils.LogFooter()
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Association not found",
			})
		}
		utils.LogMessage(utils.LevelError, "Failed to get association")
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get association information",
		})
	}

	roleName := fmt.Sprintf("%s_respo", strings.ToLower(strings.ReplaceAll(associationName, " ", "")))

	var hasPermission bool
	permissionQuery := `
       SELECT EXISTS(
          SELECT 1 FROM newf n
          JOIN newf_roles nr ON n.email = nr.email
          JOIN roles r ON nr.id_roles = r.id_roles
          WHERE n.email = $1 AND (r.name = 'ADMIN' OR r.name = $2)
       )`
	err = h.db.QueryRow(permissionQuery, userEmail, roleName).Scan(&hasPermission)
	if err != nil || !hasPermission {
		utils.LogMessage(utils.LevelWarn, "User not authorized to change responsible")
		utils.LogFooter()
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Only admins or current responsibles can add an association responsible",
		})
	}

	var req models.AddRespoAssociationRequest
	if err := c.BodyParser(&req); err != nil {
		utils.LogFooter()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.Email == "" {
		utils.LogFooter()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Email is required",
		})
	}

	var userExists, isMember bool
	err = h.db.QueryRow("SELECT EXISTS(SELECT 1 FROM newf WHERE email = $1)", req.Email).Scan(&userExists)
	if err != nil || !userExists {
		utils.LogFooter()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	err = h.db.QueryRow("SELECT EXISTS(SELECT 1 FROM associations_members WHERE email = $1 AND id_associations = $2)", req.Email, associationID).Scan(&isMember)
	if err != nil || !isMember {
		_, err = h.db.Exec("INSERT INTO associations_members (email, id_associations) VALUES ($1, $2)", req.Email, associationID)
		if err != nil {
			utils.LogFooter()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to add user to association",
			})
		}
	}

	tx, err := h.db.Begin()
	if err != nil {
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database transaction error",
		})
	}
	defer tx.Rollback()

	var roleID int
	err = tx.QueryRow("SELECT id_roles FROM roles WHERE name = $1", roleName).Scan(&roleID)
	if err != nil {
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get role information",
		})
	}

	// Safety check to avoid duplicate assignments
	var alreadyRespo bool
	err = tx.QueryRow("SELECT EXISTS(SELECT 1 FROM newf_roles WHERE email = $1 AND id_roles = $2)", req.Email, roleID).Scan(&alreadyRespo)
	if alreadyRespo {
		utils.LogMessage(utils.LevelWarn, "User is already a responsible")
		utils.LogFooter()
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "User is already a responsible",
		})
	}

	// Redundant destructive DELETE instruction removed here to support multi-responsible mapping.
	_, err = tx.Exec("INSERT INTO newf_roles (email, id_roles) VALUES ($1, $2)", req.Email, roleID)
	if err != nil {
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to add new responsible",
		})
	}

	if err = tx.Commit(); err != nil {
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save changes",
		})
	}

	utils.LogMessage(utils.LevelInfo, "Association responsible added successfully")
	utils.LogFooter()

	return c.JSON(fiber.Map{
		"message":         "Association responsible added successfully",
		"new_responsible": req.Email,
	})
}

// JoinAssociation adds the user to a association
func (h *AssociationHandler) JoinAssociation(c *fiber.Ctx) error {
	utils.LogHeader("🤝 Join Association")

	associationID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		utils.LogMessage(utils.LevelWarn, "Invalid association ID")
		utils.LogFooter()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid association ID",
		})
	}

	userEmail := c.Locals("email").(string)
	utils.LogLineKeyValue(utils.LevelInfo, "Association ID", associationID)
	utils.LogLineKeyValue(utils.LevelInfo, "User", userEmail)

	// Check if association exists
	var associationExists bool
	err = h.db.QueryRow("SELECT EXISTS(SELECT 1 FROM associations WHERE id_associations = $1)", associationID).Scan(&associationExists)
	if err != nil || !associationExists {
		utils.LogMessage(utils.LevelWarn, "Association not found")
		utils.LogFooter()
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Association not found",
		})
	}

	// Check if user is already a member
	var isMember bool
	err = h.db.QueryRow("SELECT EXISTS(SELECT 1 FROM associations_members WHERE email = $1 AND id_associations = $2)", userEmail, associationID).Scan(&isMember)
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
			"error": "You are already a member of this association",
		})
	}

	// Add user to association
	_, err = h.db.Exec("INSERT INTO associations_members (email, id_associations) VALUES ($1, $2)", userEmail, associationID)
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to join association")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to join association",
		})
	}

	utils.LogMessage(utils.LevelInfo, "Successfully joined association")
	utils.LogFooter()

	return c.JSON(fiber.Map{
		"message": "Successfully joined the association",
	})
}

// LeaveAssociation removes the user from a association
func (h *AssociationHandler) LeaveAssociation(c *fiber.Ctx) error {
	utils.LogHeader("👋 Leave Association")

	associationID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		utils.LogMessage(utils.LevelWarn, "Invalid association ID")
		utils.LogFooter()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid association ID",
		})
	}

	userEmail := c.Locals("email").(string)
	utils.LogLineKeyValue(utils.LevelInfo, "Association ID", associationID)
	utils.LogLineKeyValue(utils.LevelInfo, "User", userEmail)

	// Check if user is a member
	var isMember bool
	err = h.db.QueryRow("SELECT EXISTS(SELECT 1 FROM associations_members WHERE email = $1 AND id_associations = $2)", userEmail, associationID).Scan(&isMember)
	if err != nil || !isMember {
		utils.LogMessage(utils.LevelWarn, "User not a member")
		utils.LogFooter()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "You are not a member of this association",
		})
	}

	// Check if user is responsible for this association
	var associationName string
	err = h.db.QueryRow("SELECT name FROM associations WHERE id_associations = $1", associationID).Scan(&associationName)
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to get association name")
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get association information",
		})
	}

	roleName := fmt.Sprintf("%s_respo", strings.ToLower(strings.ReplaceAll(associationName, " ", "")))
	var isRespo bool
	err = h.db.QueryRow(`
       SELECT EXISTS(
          SELECT 1 FROM newf_roles nr
          JOIN roles r ON nr.id_roles = r.id_roles
          WHERE nr.email = $1 AND r.name = $2
       )
    `, userEmail, roleName).Scan(&isRespo)

	tx, err := h.db.Begin()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Transaction failed"})
	}
	defer tx.Rollback()

	if err == nil && isRespo {
		// Check if there are other managers in the association
		var respoCount int
		countQuery := `
          SELECT COUNT(*) 
          FROM newf_roles nr
          JOIN roles r ON nr.id_roles = r.id_roles
          WHERE r.name = $1
       `
		err = tx.QueryRow(countQuery, roleName).Scan(&respoCount)

		// Stop the execution if this user is the single remaining manager
		if respoCount <= 1 {
			utils.LogMessage(utils.LevelWarn, "Last responsible cannot leave association")
			utils.LogFooter()
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "You are the last responsible. Please assign another responsible before leaving.",
			})
		}

		// Multiple responsibles exist: user is safely unlinked from their role
		deleteRoleQuery := `
          DELETE FROM newf_roles 
          WHERE email = $1 AND id_roles = (SELECT id_roles FROM roles WHERE name = $2)
       `
		_, err = tx.Exec(deleteRoleQuery, userEmail, roleName)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to remove responsible role"})
		}
		utils.LogMessage(utils.LevelInfo, "Responsible role removed since other responsibles exist")
	}

	// Remove user from association members table
	_, err = tx.Exec("DELETE FROM associations_members WHERE email = $1 AND id_associations = $2", userEmail, associationID)
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to leave association")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to leave association",
		})
	}

	if err = tx.Commit(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to commit transaction"})
	}

	utils.LogMessage(utils.LevelInfo, "Successfully left association")
	utils.LogFooter()

	return c.JSON(fiber.Map{
		"message": "Successfully left the association",
	})
}
