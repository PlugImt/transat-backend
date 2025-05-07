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
			COALESCE(n.graduation_year, 0) AS graduation_year,
			COALESCE(n.campus, '') AS campus,
			-- COALESCE(n.notification_token, '') AS notification_token, -- Maybe don't expose token?
			n.password_updated_date, -- Consider format or omitting
			COALESCE(l.code, 'fr') AS language, -- Get language code
			(SELECT COUNT(*) FROM newf) AS total_newf -- Calculate total users separately if needed
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

	err := h.DB.QueryRowContext(ctx, query, email).Scan(
		&newf.ID,
		&newf.Email,
		&newf.FirstName,
		&newf.LastName,
		&newf.ProfilePicture,
		&newf.PhoneNumber,
		&newf.GraduationYear,
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
	if newf.GraduationYear != 0 {
		response["graduation_year"] = newf.GraduationYear
	}
	if newf.Campus != "" {
		response["campus"] = newf.Campus
	}
	if passwordUpdated.Valid {
		// Format the date/time as needed, e.g., RFC3339 or simpler date
		response["password_updated_date"] = passwordUpdated.Time.Format(time.RFC3339)
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
	if req.GraduationYear != 0 { // Assuming 0 is not a valid year
		// Add validation for year range if needed
		updateFields["graduation_year"] = req.GraduationYear
	}
	if req.Campus != "" {
		updateFields["campus"] = req.Campus
	}
	if req.ProfilePicture != "" {
		// Potentially validate the picture URL/path format
		updateFields["profile_picture"] = req.ProfilePicture
	}
	if req.NotificationToken != "" { // Allow updating push token
		updateFields["notification_token"] = req.NotificationToken
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
		// Check if only language was updated
		if req.Language != "" {
			utils.LogMessage(utils.LevelInfo, "Only language preference was updated")
			utils.LogFooter()
			return c.SendStatus(fiber.StatusOK) // Return OK if only language updated successfully
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

	// Final query
	query := fmt.Sprintf("UPDATE newf SET %s WHERE email = $%d;", strings.Join(setClauses, ", "), argIndex)

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

// --- Notification Preference Handlers ---

// AddOrRemoveNotificationSubscription adds or removes a subscription to a service notification.
func (h *UserHandler) AddOrRemoveNotificationSubscription(c *fiber.Ctx) error {
	email := c.Locals("email").(string)
	var req struct {
		Service string `json:"service"`
	}

	utils.LogHeader("📞 Add/Remove Notification Subscription")
	utils.LogLineKeyValue(utils.LevelInfo, "User", email)

	if err := c.BodyParser(&req); err != nil {
		utils.LogMessage(utils.LevelError, "Failed to parse request body")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request format (expected {\"service\":\"service_name\"})"})
	}

	if req.Service == "" {
		utils.LogMessage(utils.LevelWarn, "Missing service name in request")
		utils.LogFooter()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Service name is required"})
	}
	utils.LogLineKeyValue(utils.LevelInfo, "Service", req.Service)

	// Attempt to insert the notification preference
	insertQuery := `
		INSERT INTO notifications (email, id_services)
		SELECT $1, id_services FROM services WHERE name = $2
		ON CONFLICT (email, id_services) DO NOTHING; -- Try to insert, do nothing if it exists
	`
	result, err := h.DB.Exec(insertQuery, email, req.Service)
	if err != nil {
		// Check if it's an error because the service name doesn't exist
		// This might require a separate check or adjusting the query/error handling
		utils.LogMessage(utils.LevelError, "Error during notification insert attempt")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database error managing subscription"})
	}

	rowsAffected, _ := result.RowsAffected()

	subscribed := false
	if rowsAffected > 0 {
		// Inserted successfully, user is now subscribed
		subscribed = true
		utils.LogMessage(utils.LevelInfo, "User subscribed successfully")
	} else {
		// Insert failed (likely due to ON CONFLICT), meaning user was already subscribed. We delete it.
		utils.LogMessage(utils.LevelInfo, "User already subscribed, attempting to unsubscribe")
		deleteQuery := `
			DELETE FROM notifications
			WHERE email = $1
			  AND id_services = (SELECT id_services FROM services WHERE name = $2);
		`
		delResult, delErr := h.DB.Exec(deleteQuery, email, req.Service)
		if delErr != nil {
			utils.LogMessage(utils.LevelError, "Failed to delete existing subscription")
			utils.LogLineKeyValue(utils.LevelError, "Error", delErr)
			utils.LogFooter()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to unsubscribe"})
		}
		delRowsAffected, _ := delResult.RowsAffected()
		if delRowsAffected > 0 {
			subscribed = false // Unsubscribed successfully
			utils.LogMessage(utils.LevelInfo, "User unsubscribed successfully")
		} else {
			// This state should be rare (insert failed, delete failed) - maybe service name invalid?
			utils.LogMessage(utils.LevelWarn, "Unsubscribe affected 0 rows after insert conflict")
			// Check if service exists
			var serviceExists bool
			checkService := `SELECT EXISTS(SELECT 1 FROM services WHERE name = $1)`
			err := h.DB.QueryRow(checkService, req.Service).Scan(&serviceExists)
			if err != nil {
				utils.LogMessage(utils.LevelError, "Failed to check if service exists")
				utils.LogLineKeyValue(utils.LevelError, "Error", err)
			}
			if !serviceExists {
				utils.LogFooter()
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": fmt.Sprintf("Service '%s' not found", req.Service)})
			}
			utils.LogFooter()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not update subscription status"})
		}
	}

	utils.LogFooter()
	// Return the final state
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"subscribed": subscribed})
}

// GetNotificationSubscriptions retrieves the notification subscription status for the user.
// Can get all subscriptions or check specific ones.
func (h *UserHandler) GetNotificationSubscriptions(c *fiber.Ctx) error {
	email := c.Locals("email").(string)
	utils.LogHeader("📞 Get Notification Subscriptions")
	utils.LogLineKeyValue(utils.LevelInfo, "User", email)

	// Check if request body exists to determine mode (get all vs check specific)
	// Using Peek allows checking without consuming the body if it's needed later (though BodyParser consumes it)
	if len(c.Body()) == 0 {
		// --- Get All Subscribed Services ---
		utils.LogMessage(utils.LevelInfo, "Requesting all subscribed services")
		query := `
			SELECT s.name
			FROM notifications n
			JOIN services s ON n.id_services = s.id_services
			WHERE n.email = $1;
		`
		rows, err := h.DB.Query(query, email)
		if err != nil {
			utils.LogMessage(utils.LevelError, "Failed to query all subscriptions")
			utils.LogLineKeyValue(utils.LevelError, "Error", err)
			utils.LogFooter()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database error fetching subscriptions"})
		}
		defer rows.Close()

		var servicesList []string
		for rows.Next() {
			var serviceName string
			if err := rows.Scan(&serviceName); err != nil {
				utils.LogMessage(utils.LevelError, "Failed to scan service name")
				utils.LogLineKeyValue(utils.LevelError, "Error", err)
				// Continue scanning other rows
			} else {
				servicesList = append(servicesList, serviceName)
			}
		}
		if err := rows.Err(); err != nil {
			utils.LogMessage(utils.LevelError, "Error during row iteration for subscriptions")
			utils.LogLineKeyValue(utils.LevelError, "Error", err)
			// Return partial list or error?
		}

		utils.LogMessage(utils.LevelInfo, "Successfully fetched all subscriptions")
		utils.LogLineKeyValue(utils.LevelInfo, "Services", servicesList)
		utils.LogFooter()
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"services": servicesList})

	} else {
		// --- Check Specific Service(s) ---
		// Expects {"services": ["service1", "service2"]} or {"service": "service1"} (legacy)
		var payload struct {
			Services []string `json:"services"`
			Service  string   `json:"service"` // For legacy single check
		}

		if err := c.BodyParser(&payload); err != nil {
			utils.LogMessage(utils.LevelError, "Failed to parse request body for specific services")
			utils.LogLineKeyValue(utils.LevelError, "Error", err)
			utils.LogFooter()
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request format"})
		}

		servicesToCheck := payload.Services
		if len(servicesToCheck) == 0 && payload.Service != "" {
			utils.LogMessage(utils.LevelInfo, "Requesting single service status (legacy)")
			servicesToCheck = []string{payload.Service}
		}

		if len(servicesToCheck) == 0 {
			utils.LogMessage(utils.LevelWarn, "No services provided in request body")
			utils.LogFooter()
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "No services specified for checking"})
		}

		utils.LogMessage(utils.LevelInfo, "Requesting status for specific services")
		utils.LogLineKeyValue(utils.LevelInfo, "Services", servicesToCheck)

		// Create placeholders for the IN clause: $2, $3, $4...
		placeholders := make([]string, len(servicesToCheck))
		args := make([]interface{}, len(servicesToCheck)+1)
		args[0] = email // $1 is the email
		for i, service := range servicesToCheck {
			placeholders[i] = fmt.Sprintf("$%d", i+2)
			args[i+1] = service
		}

		query := fmt.Sprintf(`
			SELECT s.name
			FROM notifications n
			JOIN services s ON n.id_services = s.id_services
			WHERE n.email = $1 AND s.name IN (%s);
		`, strings.Join(placeholders, ", "))

		rows, err := h.DB.Query(query, args...)
		if err != nil {
			utils.LogMessage(utils.LevelError, "Failed to query specific subscriptions")
			utils.LogLineKeyValue(utils.LevelError, "Error", err)
			utils.LogFooter()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database error checking subscriptions"})
		}
		defer rows.Close()

		subscribedServices := []string{}
		for rows.Next() {
			var serviceName string
			if err := rows.Scan(&serviceName); err != nil {
				utils.LogMessage(utils.LevelError, "Failed to scan subscribed service name")
				utils.LogLineKeyValue(utils.LevelError, "Error", err)
			} else {
				subscribedServices = append(subscribedServices, serviceName)
			}
		}
		if err := rows.Err(); err != nil {
			utils.LogMessage(utils.LevelError, "Error during row iteration for specific subscriptions")
			utils.LogLineKeyValue(utils.LevelError, "Error", err)
		}

		// If the original request was for a single service (legacy), return boolean
		if payload.Service != "" && len(payload.Services) == 0 {
			isSubscribed := false
			for _, subbedService := range subscribedServices {
				if subbedService == payload.Service {
					isSubscribed = true
					break
				}
			}
			utils.LogMessage(utils.LevelInfo, "Successfully checked single subscription status")
			utils.LogFooter()
			return c.Status(fiber.StatusOK).JSON(fiber.Map{"subscribed": isSubscribed})
		}

		// Otherwise, return the list of subscribed services from the requested list
		utils.LogMessage(utils.LevelInfo, "Successfully checked specific subscriptions")
		utils.LogFooter()
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"services": subscribedServices})
	}
}

// SendNotification is the handler to trigger sending a notification.
// It might require admin privileges. It uses the NotificationService.
func (h *UserHandler) SendNotification(c *fiber.Ctx) error {
	// This handler acts as the entry point for the /send-notification route.
	// It needs to:
	// 1. Check user permissions (is the logged-in user allowed to send?). Usually requires admin role.
	// 2. Call the NotificationService's SendNotification method to handle the actual logic.

	// --- Permission Check Placeholder ---
	// role := c.Locals("role").(string) // Get role from JWT
	// if role != "ADMIN" { // Example role check
	//    utils.LogMessage(utils.LevelWarn, "Unauthorized attempt to send notification")
	//    utils.LogLineKeyValue(utils.LevelWarn, "User", c.Locals("email"))
	// 	  utils.LogFooter()
	// 	  return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Permission denied"})
	// }
	// --- End Permission Check ---

	if h.NotifService == nil {
		utils.LogMessage(utils.LevelError, "NotificationService is not available in UserHandler")
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Notification service not configured"})
	}

	// Delegate the core logic to the NotificationService's method
	// Note: The SendNotification method in NotificationService is currently a placeholder.
	// It needs to be implemented fully to handle parsing, target resolution, etc.
	return h.NotifService.SendNotification(c)
}
