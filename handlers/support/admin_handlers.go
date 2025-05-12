package support

import (
	"database/sql"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/plugimt/transat-backend/models"
	"github.com/plugimt/transat-backend/utils"
)

// UpdateSupportStatus updates the status and response for a support request (admin only)
func (h *SupportHandler) UpdateSupportStatus(c *fiber.Ctx) error {
	utils.LogHeader("✏️ Update Support Request")

	// Get user email from JWT (set by middleware)
	email, ok := c.Locals("email").(string)
	if !ok || email == "" {
		utils.LogMessage(utils.LevelWarn, "User email not found in token")
		utils.LogFooter()
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User email not found in token",
		})
	}
	utils.LogLineKeyValue(utils.LevelInfo, "Admin User", email)

	// Verify admin permissions
	isAdmin := false
	err := h.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM newf_roles WHERE email = $1 AND id_roles = 1)", email).Scan(&isAdmin)
	if err != nil || !isAdmin {
		utils.LogMessage(utils.LevelWarn, "Unauthorized access attempt to update support request")
		utils.LogLineKeyValue(utils.LevelWarn, "User", email)
		utils.LogFooter()
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Admin permission required",
		})
	}

	// Get support request ID from URL params
	supportID, err := c.ParamsInt("id")
	if err != nil {
		utils.LogMessage(utils.LevelWarn, "Invalid support request ID")
		utils.LogLineKeyValue(utils.LevelWarn, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid support request ID",
		})
	}
	utils.LogLineKeyValue(utils.LevelInfo, "Support ID", supportID)

	// Parse input
	var input struct {
		Status   string `json:"status" validate:"required,oneof=pending in_progress resolved"`
		Response string `json:"response"`
	}

	if err := c.BodyParser(&input); err != nil {
		utils.LogMessage(utils.LevelWarn, "Failed to parse request body")
		utils.LogLineKeyValue(utils.LevelWarn, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request format",
		})
	}

	// Validate input
	//if err := utils.Validator.Struct(input); err != nil {
	//	utils.LogMessage(utils.LevelWarn, "Validation failed")
	//	utils.LogLineKeyValue(utils.LevelWarn, "Error", err)
	//	utils.LogFooter()
	//	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
	//		"error": "Validation failed: " + err.Error(),
	//	})
	//}

	// Begin transaction
	tx, err := h.DB.Begin()
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to begin transaction")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database error",
		})
	}
	defer tx.Rollback()

	// Update support request
	result, err := tx.Exec(
		`UPDATE support 
         SET status = $1, response = $2
         WHERE id_support = $3`,
		input.Status, input.Response, supportID,
	)

	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to update support request")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update support request",
		})
	}

	// Check if request exists
	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		utils.LogMessage(utils.LevelWarn, "Support request not found")
		utils.LogFooter()
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Support request not found",
		})
	}

	// Get requester email for notification
	var userEmail string
	err = tx.QueryRow("SELECT email FROM support WHERE id_support = $1", supportID).Scan(&userEmail)
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to get requester email")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to process notification",
		})
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		utils.LogMessage(utils.LevelError, "Failed to commit transaction")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database error",
		})
	}

	// Send status update notification (async)
	if userEmail != "" {
		go h.sendStatusUpdateNotification(userEmail, supportID, input.Status)
	}

	utils.LogMessage(utils.LevelInfo, "Support request updated successfully")
	utils.LogFooter()

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Support request updated successfully",
	})
}

// Helper method to send status update notification to the requester
func (h *SupportHandler) sendStatusUpdateNotification(email string, supportID int, status string) {
	// Only proceed if email service is available
	if h.EmailService == nil {
		utils.LogMessage(utils.LevelWarn, "EmailService is not initialized, skipping support status update email.")
		return
	}

	// Get request details
	var subject string
	var response sql.NullString
	var languageCode string

	err := h.DB.QueryRow(`
		SELECT s.subject, s.response, COALESCE(l.code, 'fr') as lang_code
		FROM support s
		LEFT JOIN newf n ON s.email = n.email
		LEFT JOIN languages l ON n.language = l.id_languages
		WHERE s.id_support = $1
	`, supportID).Scan(&subject, &response, &languageCode)

	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to get support request details for notification")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		return
	}

	// Format status for email
	statusDisplay := "Updated"
	switch status {
	case "in_progress":
		statusDisplay = "In Progress"
	case "resolved":
		statusDisplay = "Resolved"
	}

	// Extract first name from user email
	firstName := ""
	emailParts := strings.Split(email, "@")
	nameParts := strings.Split(emailParts[0], ".")
	if len(nameParts) > 0 && len(nameParts[0]) > 0 {
		fName := nameParts[0]
		firstName = strings.ToUpper(string(fName[0])) + strings.ToLower(fName[1:])
	}

	// Prepare email data
	emailData := struct {
		FirstName  string
		Subject    string
		Status     string
		Response   string
		Date       string
		Time       string
		SupportURL string
	}{
		FirstName:  firstName,
		Subject:    subject,
		Status:     statusDisplay,
		Response:   response.String,
		Date:       time.Now().Format("02/01/2006"),
		Time:       time.Now().Format("15:04:05"),
		SupportURL: "https://your-app-url.com/support", // Update with your app's support URL
	}

	// Capture data for the goroutine
	recipient := email
	template := "email_templates/support_update.html"
	subjectKey := "email_support_update.subject" // Define this in your translation files
	lang := languageCode

	// Run email sending in a separate goroutine
	go func(recipient, template, subjectKey, lang string, data interface{}) {
		emailErr := h.EmailService.SendEmail(models.Email{
			Recipient:  recipient,
			Template:   template,
			SubjectKey: subjectKey,
			Language:   lang,
		}, data)
		if emailErr != nil {
			utils.LogMessage(utils.LevelError, "Failed to send support status update email asynchronously")
			utils.LogLineKeyValue(utils.LevelError, "Email", recipient)
			utils.LogLineKeyValue(utils.LevelError, "Error", emailErr)
		}
	}(recipient, template, subjectKey, lang, emailData)
}
