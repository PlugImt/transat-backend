package user

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/getsentry/sentry-go"
	sentryfiber "github.com/getsentry/sentry-go/fiber"
	"github.com/gofiber/fiber/v2"
	"github.com/plugimt/transat-backend/models"
	"github.com/plugimt/transat-backend/services" // For NotificationService
	"github.com/plugimt/transat-backend/utils"    // For logger
)

// UserHandler handles user profile and related actions.
type UserHandler struct {
	DB           *sql.DB
	NotifService *services.NotificationService // Inject if needed for notification handlers
}

// NewUserHandler creates a new UserHandler.
func NewUserHandler(db *sql.DB, notifService *services.NotificationService) *UserHandler {
	return &UserHandler{
		DB:           db,
		NotifService: notifService,
	}
}

// GetNewf retrieves the profile information for the logged-in user.
func (h *UserHandler) GetNewf(c *fiber.Ctx) error {
	email := c.Locals("email").(string) // Assumes email is set by JWTMiddleware
	ctx := c.UserContext()              // Obtenir le context.Context de Fiber

	utils.LogHeader("📧 Get Newf Profile")
	utils.LogLineKeyValue(utils.LevelInfo, "User", email)

	query := `
		SELECT
			n.id_newf,
			n.email,
			n.first_name,
			n.last_name,
			COALESCE(n.profile_picture, '') AS profile_picture,
			COALESCE(n.phone_number, '') AS phone_number,
			n.graduation_year,
			COALESCE(n.formation_name, '') AS formation_name,
			COALESCE(n.campus, '') AS campus,
			-- COALESCE(n.notification_token, '') AS notification_token, -- Maybe don't expose token?
			n.password_updated_date, -- Consider format or omitting
			COALESCE(l.code, 'fr') AS language, -- Get language code
			(SELECT id_newf FROM newf ORDER BY creation_date DESC LIMIT 1) AS total_newf -- Calculate total users separately if needed
		FROM newf n
		LEFT JOIN languages l ON n.language = l.id_languages
		WHERE n.email = $1;
	`

	var newf models.Newf             // Use the full model, but only populate relevant fields for response
	var passwordUpdated sql.NullTime // Use sql.NullTime for potentially null dates

	parentSpan := sentryfiber.GetSpanFromContext(c)
	var querySpan *sentry.Span

	if parentSpan != nil {
		querySpan = parentSpan.StartChild("db.sql.query")
		querySpan.SetTag("db.system", "postgresql")
		querySpan.SetData("db.statement", query)
		querySpan.SetData("db.operation", "SELECT")
		querySpan.SetData("db.table", "newf, languages")
		defer querySpan.Finish()
	}

	var graduationYear sql.NullInt32
	err := h.DB.QueryRowContext(ctx, query, email).Scan(
		&newf.ID,
		&newf.Email,
		&newf.FirstName,
		&newf.LastName,
		&newf.ProfilePicture,
		&newf.PhoneNumber,
		&graduationYear,
		&newf.FormationName,
		&newf.Campus,
		// &newf.NotificationToken, // Omitted
		&passwordUpdated,
		&newf.Language,
		&newf.TotalUsers,
	)

	if querySpan != nil {
		if err != nil {
			if err == sql.ErrNoRows {
				querySpan.Status = sentry.SpanStatusNotFound
			} else {
				querySpan.Status = sentry.SpanStatusInternalError
			}
			querySpan.SetData("error", err.Error())
		} else {
			querySpan.Status = sentry.SpanStatusOK
		}
	}

	if err != nil {
		if err == sql.ErrNoRows {
			utils.LogMessage(utils.LevelWarn, "User profile not found")
			utils.LogFooter()
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User profile not found"})
		}
		utils.LogMessage(utils.LevelError, "Failed to fetch user profile")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve profile"})
	}

	// Set graduation year if valid
	if graduationYear.Valid {
		year := int(graduationYear.Int32)
		newf.GraduationYear = &year
	}

	// Create response map, explicitly adding non-zero/non-empty fields
	response := make(map[string]interface{})
	response["id_newf"] = newf.ID
	response["email"] = newf.Email
	response["first_name"] = newf.FirstName
	response["last_name"] = newf.LastName
	response["language"] = newf.Language
	response["total_newf"] = newf.TotalUsers // Assuming total_newf is calculated correctly

	if newf.ProfilePicture != "" {
		response["profile_picture"] = newf.ProfilePicture
	}
	if newf.PhoneNumber != "" {
		response["phone_number"] = newf.PhoneNumber
	}
	if newf.GraduationYear != nil {
		response["graduation_year"] = *newf.GraduationYear
	}
	if newf.FormationName != "" {
		response["formation_name"] = newf.FormationName
	}
	if newf.Campus != "" {
		response["campus"] = newf.Campus
	}
	if passwordUpdated.Valid {
		response["password_updated_date"] = utils.FormatParis(passwordUpdated.Time, time.RFC3339)
	}

	utils.LogMessage(utils.LevelInfo, "User profile fetched successfully")
	utils.LogFooter()

	return c.Status(fiber.StatusOK).JSON(response)
}

// UpdateNewf updates the profile information for the logged-in user.
func (h *UserHandler) UpdateNewf(c *fiber.Ctx) error {
	email := c.Locals("email").(string)
	var req models.Newf // Use the Newf model to parse incoming update data

	utils.LogHeader("📧 Update Newf Profile")
	utils.LogLineKeyValue(utils.LevelInfo, "User", email)

	if err := c.BodyParser(&req); err != nil {
		utils.LogMessage(utils.LevelError, "Failed to parse request body")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Failed to parse your data"})
	}

	// Build query dynamically based on provided fields
	updateFields := make(map[string]interface{})
	queryArgs := []interface{}{}
	argIndex := 1

	if req.FirstName != "" {
		updateFields["first_name"] = req.FirstName
	}
	if req.LastName != "" {
		updateFields["last_name"] = req.LastName
	}
	if req.PhoneNumber != "" {
		// Add validation for phone number format if needed
		updateFields["phone_number"] = req.PhoneNumber
	}
	if req.GraduationYear != nil && *req.GraduationYear != 0 {
		if *req.GraduationYear < 1900 || *req.GraduationYear > utils.GetYearParis(utils.Now())+5 {
			utils.LogMessage(utils.LevelWarn, "Invalid graduation year")
			utils.LogFooter()
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid graduation year"})
		}
		updateFields["graduation_year"] = *req.GraduationYear
	} else if req.GraduationYear != nil && *req.GraduationYear == 0 {
		// Set to NULL if explicitly set to 0
		updateFields["graduation_year"] = nil
	}
	if req.FormationName != "" {
		req.FormationName = strings.ToUpper(req.FormationName)

		// Validate formation name against known values if needed
		if !utils.IsValidFormationName(req.FormationName) {
			utils.LogMessage(utils.LevelWarn, "Invalid formation name")
			utils.LogFooter()
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid formation name"})
		}
		updateFields["formation_name"] = req.FormationName
	}
	if req.Campus != "" {
		updateFields["campus"] = req.Campus
	}
	if req.ProfilePicture != "" {
		// Potentially validate the picture URL/path format
		updateFields["profile_picture"] = req.ProfilePicture
	}
	if req.NotificationToken != "" { // Save notification token to user_notification_tokens table
		// Insert token into user_notification_tokens table (allows multiple tokens per user)
		// ON CONFLICT DO NOTHING prevents duplicate tokens for the same user
		tokenQuery := `INSERT INTO user_notification_tokens (email, token) VALUES ($1, $2) ON CONFLICT (email, token) DO NOTHING`
		_, err := h.DB.Exec(tokenQuery, email, req.NotificationToken)
		if err != nil {
			utils.LogMessage(utils.LevelError, "Failed to save notification token")
			utils.LogLineKeyValue(utils.LevelError, "Error", err)
			// Log but continue with other fields
		} else {
			utils.LogMessage(utils.LevelInfo, "Notification token saved successfully")
		}
	}
	if req.Language != "" { // Allow updating language preference
		// Language update needs a subquery to get the ID
		langQuery := `UPDATE newf SET language = (SELECT id_languages FROM languages WHERE code = $1 LIMIT 1) WHERE email = $2`
		_, err := h.DB.Exec(langQuery, strings.ToLower(req.Language), email)
		if err != nil {
			utils.LogMessage(utils.LevelError, "Failed to update language preference")
			utils.LogLineKeyValue(utils.LevelError, "Language Code", req.Language)
			utils.LogLineKeyValue(utils.LevelError, "Error", err)
			// Log but maybe don't fail the whole update? Or return partial success?
			// For now, log and continue with other fields.
		} else {
			utils.LogMessage(utils.LevelInfo, "Language preference updated")
			utils.LogLineKeyValue(utils.LevelInfo, "Language Code", req.Language)
		}
	}

	if len(updateFields) == 0 {
		// Check if only language or notification token was updated
		if req.Language != "" || req.NotificationToken != "" {
			utils.LogMessage(utils.LevelInfo, "Only language preference or notification token was updated")
			utils.LogFooter()
			return c.SendStatus(fiber.StatusOK) // Return OK if only language or token updated successfully
		}
		utils.LogMessage(utils.LevelWarn, "No fields provided for update")
		utils.LogFooter()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "No data provided to update"})
	}

	// Construct the SET part of the SQL query
	var setClauses []string
	for column, value := range updateFields {
		setClauses = append(setClauses, fmt.Sprintf("%s = $%d", column, argIndex))
		queryArgs = append(queryArgs, value)
		argIndex++
	}

	// Add email to the end of arguments for the WHERE clause
	queryArgs = append(queryArgs, email)

	// Replace with parameterized query construction
	setClause := strings.Join(setClauses, ", ")
	query := "UPDATE newf SET " + setClause + fmt.Sprintf(" WHERE email = $%d;", argIndex)

	// Execute the update
	result, err := h.DB.Exec(query, queryArgs...)
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to update user profile")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update profile"})
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		// This shouldn't happen if the JWT middleware ensures user exists
		utils.LogMessage(utils.LevelWarn, "Update profile query affected 0 rows")
		utils.LogFooter()
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User profile not found for update"})
	}

	utils.LogMessage(utils.LevelInfo, "User profile updated successfully")
	utils.LogLineKeyValue(utils.LevelInfo, "Fields Updated", updateFields) // Log only keys for brevity/privacy
	utils.LogFooter()

	return c.SendStatus(fiber.StatusOK)
}

// DeleteNewf handles the deletion of the logged-in user's account.
// IMPORTANT: This is a destructive action and needs careful consideration.
// Consider soft delete instead of hard delete.
func (h *UserHandler) DeleteNewf(c *fiber.Ctx) error {
	email := c.Locals("email").(string)

	utils.LogHeader("🚫 Delete Newf Account")
	utils.LogLineKeyValue(utils.LevelWarn, "User requesting deletion", email) // Log clearly this is a deletion

	// --- !! DANGER ZONE: Hard Delete !! ---
	// Consider alternatives:
	// 1. Soft Delete: Add an 'is_deleted' flag or 'deleted_at' timestamp.
	// 2. Anonymization: Remove PII but keep related non-PII data.
	// 3. Confirmation Step: Require password or email confirmation.

	// Hard delete implementation (use with caution):
	query := `DELETE FROM newf WHERE email = $1;`
	result, err := h.DB.Exec(query, email)
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to delete user account")
		utils.LogLineKeyValue(utils.LevelError, "Email", email)
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete account"})
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		// Should not happen if JWT is valid
		utils.LogMessage(utils.LevelError, "Attempted to delete non-existent user?")
		utils.LogLineKeyValue(utils.LevelError, "Email", email)
		utils.LogFooter()
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	// Consider deleting related data (files, posts, reactions, etc.) based on foreign key constraints or explicitly.

	utils.LogMessage(utils.LevelInfo, "User account deleted successfully")
	utils.LogFooter()

	// Maybe clear JWT cookie or advise client to log out?
	return c.SendStatus(fiber.StatusOK) // Or 204 No Content
}
