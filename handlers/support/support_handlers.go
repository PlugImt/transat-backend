package support

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/plugimt/transat-backend/models"
	"github.com/plugimt/transat-backend/services"
	"github.com/plugimt/transat-backend/utils"
)

// SupportHandler handles support request operations
type SupportHandler struct {
	DB           *sql.DB
	EmailService *services.EmailService
}

// NewSupportHandler creates a new SupportHandler
func NewSupportHandler(db *sql.DB, emailService *services.EmailService) *SupportHandler {
	return &SupportHandler{
		DB:           db,
		EmailService: emailService,
	}
}

// GetSupportRequests returns all support requests for the authenticated user
func (h *SupportHandler) GetSupportRequests(c *fiber.Ctx) error {
	utils.LogHeader("📋 Get Support Requests")

	// Get user email from JWT (set by middleware)
	email, ok := c.Locals("email").(string)
	if !ok || email == "" {
		utils.LogMessage(utils.LevelWarn, "User email not found in token")
		utils.LogFooter()
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User email not found in token",
		})
	}
	utils.LogLineKeyValue(utils.LevelInfo, "User", email)

	// Check if user is admin to see all requests
	isAdmin := false
	err := h.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM newf_roles WHERE email = $1 AND id_roles = 1)", email).Scan(&isAdmin)
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to check admin status")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to check user permissions",
		})
	}

	// Build the query based on user role
	query := `
		SELECT 
			id_support, email, subject, message, 
			creation_date, status, response
		FROM support
	`
	var args []interface{}
	if !isAdmin {
		// Regular users can only see their own requests
		query += " WHERE email = $1"
		args = append(args, email)
	}
	query += " ORDER BY creation_date DESC"

	// Execute the query
	rows, err := h.DB.Query(query, args...)
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to retrieve support requests")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve support requests",
		})
	}
	defer rows.Close()

	// Process results
	var requests []models.SupportRequest
	for rows.Next() {
		var req models.SupportRequest
		var response sql.NullString

		err := rows.Scan(
			&req.ID,
			&req.Email,
			&req.Subject,
			&req.Message,
			&req.CreatedAt,
			&req.Status,
			&response,
		)

		if err != nil {
			utils.LogMessage(utils.LevelError, "Error scanning support request")
			utils.LogLineKeyValue(utils.LevelError, "Error", err)
			continue
		}

		if response.Valid {
			req.Response = response.String
		}

		// Get any attached media
		req.ImageURL = h.getImageURL(req.ID)

		requests = append(requests, req)
	}

	if err = rows.Err(); err != nil {
		utils.LogMessage(utils.LevelError, "Error iterating support requests")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
	}

	utils.LogMessage(utils.LevelInfo, "Support requests retrieved successfully")
	utils.LogLineKeyValue(utils.LevelInfo, "Count", len(requests))
	utils.LogFooter()

	return c.JSON(models.SupportRequestResponse{
		Requests: requests,
	})
}

// CreateSupportRequest creates a new support request
func (h *SupportHandler) CreateSupportRequest(c *fiber.Ctx) error {
	utils.LogHeader("📝 Create Support Request")

	// Get user email from JWT (set by middleware)
	email, ok := c.Locals("email").(string)
	if !ok || email == "" {
		utils.LogMessage(utils.LevelWarn, "User email not found in token")
		utils.LogFooter()
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User email not found in token",
		})
	}
	utils.LogLineKeyValue(utils.LevelInfo, "User", email)

	// Parse input
	var input models.SupportRequestInput
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

	// Insert support request
	var requestID int
	err = tx.QueryRow(
		`INSERT INTO support (email, subject, message, status, creation_date) 
         VALUES ($1, $2, $3, 'pending', CURRENT_TIMESTAMP) 
         RETURNING id_support, creation_date`,
		email, input.Subject, input.Message,
	).Scan(&requestID, &time.Time{})

	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to create support request")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create support request",
		})
	}

	// Link file if provided
	imageURL := ""
	if input.FileID > 0 {
		// Check if file exists and belongs to user
		var count int
		err = tx.QueryRow(
			"SELECT COUNT(*) FROM files WHERE id_files = $1 AND email = $2",
			input.FileID, email,
		).Scan(&count)

		if err != nil {
			utils.LogMessage(utils.LevelError, "Failed to verify file ownership")
			utils.LogLineKeyValue(utils.LevelError, "Error", err)
			utils.LogFooter()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to verify file",
			})
		}

		if count == 0 {
			utils.LogMessage(utils.LevelWarn, "File does not exist or user does not own it")
			utils.LogFooter()
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid file ID or you don't own this file",
			})
		}

		// Link file to support request
		_, err = tx.Exec(
			"INSERT INTO support_medias (id_support, id_files) VALUES ($1, $2)",
			requestID, input.FileID,
		)

		if err != nil {
			utils.LogMessage(utils.LevelError, "Failed to link file to support request")
			utils.LogLineKeyValue(utils.LevelError, "Error", err)
			utils.LogFooter()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to link file to support request",
			})
		}

		// Get file URL
		var filename string
		err = tx.QueryRow(
			"SELECT path FROM files WHERE id_files = $1",
			input.FileID,
		).Scan(&filename)

		if err == nil {
			imageURL = fmt.Sprintf("/api/data/%s", filename)
		}
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

	// Send email notification (async)
	go h.sendSupportNotification(email, input.Subject)

	utils.LogMessage(utils.LevelInfo, "Support request created successfully")
	utils.LogLineKeyValue(utils.LevelInfo, "ID", requestID)
	utils.LogFooter()

	// Return request
	request := models.SupportRequest{
		ID:        requestID,
		Email:     email,
		Subject:   input.Subject,
		Message:   input.Message,
		Status:    "pending",
		ImageURL:  imageURL,
		CreatedAt: time.Now(),
	}

	return c.Status(fiber.StatusCreated).JSON(models.SupportRequestCreateResponse{
		Request: request,
	})
}

// Helper method to get image URL for a support request
func (h *SupportHandler) getImageURL(supportID int) string {
	var path string
	err := h.DB.QueryRow(`
		SELECT f.path 
		FROM support_medias sm
		JOIN files f ON sm.id_files = f.id_files
		WHERE sm.id_support = $1
		LIMIT 1
	`, supportID).Scan(&path)

	if err != nil {
		return ""
	}

	return fmt.Sprintf("/api/data/%s", path)
}

// Helper method to send notification about new support request
func (h *SupportHandler) sendSupportNotification(email string, subject string) {
	// Only proceed if email service is available
	if h.EmailService == nil {
		utils.LogMessage(utils.LevelWarn, "EmailService is not initialized, skipping support notification email.")
		return
	}

	// Get admin emails
	rows, err := h.DB.Query(`
		SELECT n.email, COALESCE(l.code, 'fr') as lang_code
		FROM newf_roles nr
		JOIN newf n ON nr.email = n.email
		LEFT JOIN languages l ON n.language = l.id_languages
		WHERE nr.id_roles = 1
	`)
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to get admin emails for notification")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		return
	}
	defer rows.Close()

	// Send notification email to each admin
	for rows.Next() {
		var adminEmail, langCode string
		if err := rows.Scan(&adminEmail, &langCode); err != nil {
			utils.LogMessage(utils.LevelError, "Error scanning admin email row")
			utils.LogLineKeyValue(utils.LevelError, "Error", err)
			continue
		}

		// Extract first name from admin email
		firstName := ""
		emailParts := strings.Split(adminEmail, "@")
		nameParts := strings.Split(emailParts[0], ".")
		if len(nameParts) > 0 && len(nameParts[0]) > 0 {
			fName := nameParts[0]
			firstName = strings.ToUpper(string(fName[0])) + strings.ToLower(fName[1:])
		}

		// Prepare email data
		emailData := struct {
			AdminName  string
			UserEmail  string
			Subject    string
			Date       string
			Time       string
			SupportURL string
		}{
			AdminName:  firstName,
			UserEmail:  email,
			Subject:    subject,
			Date:       time.Now().Format("02/01/2006"),
			Time:       time.Now().Format("15:04:05"),
			SupportURL: "https://your-app-url.com/admin/support", // Update with your admin URL
		}

		// Capture data for the goroutine
		recipient := adminEmail
		template := "email_templates/support_request.html"
		subjectKey := "email_support_request.subject" // Define this in your translation files
		lang := langCode

		// Run email sending in a separate goroutine
		go func(recipient, template, subjectKey, lang string, data interface{}) {
			emailErr := h.EmailService.SendEmail(models.Email{
				Recipient:  recipient,
				Template:   template,
				SubjectKey: subjectKey,
				Language:   lang,
			}, data)
			if emailErr != nil {
				utils.LogMessage(utils.LevelError, "Failed to send support notification email asynchronously")
				utils.LogLineKeyValue(utils.LevelError, "Email", recipient)
				utils.LogLineKeyValue(utils.LevelError, "Error", emailErr)
			}
		}(recipient, template, subjectKey, lang, emailData)
	}
}
