package auth

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/plugimt/transat-backend/models"
	"github.com/plugimt/transat-backend/services" // For NotificationService
	"github.com/plugimt/transat-backend/utils"    // For EmailService and helpers
	"golang.org/x/crypto/bcrypt"
)

// AuthHandler handles authentication related requests.
type AuthHandler struct {
	DB           *sql.DB
	JwtSecret    []byte
	NotifService *services.NotificationService
	EmailService *services.EmailService
}

// NewAuthHandler creates a new AuthHandler.
func NewAuthHandler(db *sql.DB, jwtSecret []byte, notifService *services.NotificationService, emailService *services.EmailService) *AuthHandler {
	if emailService == nil {
		log.Println("Warning: EmailService is nil in NewAuthHandler")
	}
	return &AuthHandler{
		DB:           db,
		JwtSecret:    jwtSecret,
		NotifService: notifService, // Store the notification service if needed later
		EmailService: emailService,
	}
}

// Register handles new user registration.
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var newf models.Newf

	utils.LogHeader("👶 Newf Registration") // Using logger from utils

	if err := c.BodyParser(&newf); err != nil {
		utils.LogMessage(utils.LevelError, "Failed to parse request body")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Failed to parse your data"})
	}

	// Validate required fields
	if newf.Email == "" || newf.Password == "" {
		utils.LogMessage(utils.LevelWarn, "Registration request missing email or password")
		utils.LogFooter()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Email and password are required"})
	}

	emailValid, err := utils.CheckEmail(newf.Email) // Use helper from utils
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to check email validity")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "There is an error checking your email"})
	}
	if !emailValid {
		utils.LogMessage(utils.LevelError, "Invalid email format")
		utils.LogLineKeyValue(utils.LevelError, "Email", newf.Email)
		utils.LogFooter()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid email format"})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newf.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to hash password")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Something went wrong hashing password"})
	}

	newf.Password = string(hashedPassword)
	// Auto-generate names from email
	parts := strings.Split(newf.Email, "@")
	nameParts := strings.Split(parts[0], ".")
	if len(nameParts) >= 1 {
		firstName := nameParts[0]
		if len(firstName) > 0 {
			newf.FirstName = strings.ToUpper(string(firstName[0])) + strings.ToLower(firstName[1:])
		}
	}
	if len(nameParts) >= 2 {
		lastName := nameParts[1]
		if len(lastName) > 0 {
			newf.LastName = strings.ToUpper(lastName)
		}
	}

	// Set default language if not provided
	if newf.Language == "" {
		newf.Language = "fr" // Default to French
	} else {
		newf.Language = strings.ToLower(newf.Language)
	}

	// Use transaction for multi-step DB operations
	tx, err := h.DB.Begin()
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to begin database transaction")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database transaction error"})
	}
	// Defer rollback in case of errors
	commitOrRollback := func(tx *sql.Tx, err error) {
		if err != nil {
			rbErr := tx.Rollback()
			if rbErr != nil {
				utils.LogMessage(utils.LevelError, "Failed to rollback transaction")
				utils.LogLineKeyValue(utils.LevelError, "Rollback Error", rbErr)
			}
		} else {
			commitErr := tx.Commit()
			if commitErr != nil {
				utils.LogMessage(utils.LevelError, "Failed to commit transaction")
				utils.LogLineKeyValue(utils.LevelError, "Commit Error", commitErr)
				// Set err so subsequent steps know commit failed? Or handle differently?
				err = commitErr // Propagate commit error if needed
			}
		}
	}

	// Insert user
	insertUserQuery := `
		INSERT INTO newf (email, password, first_name, last_name, language)
		VALUES ($1, $2, $3, $4, COALESCE ((SELECT id_languages FROM languages WHERE code = $5 LIMIT 1), 1));
	`
	_, err = tx.Exec(insertUserQuery, strings.ToLower(newf.Email), newf.Password, newf.FirstName, newf.LastName, newf.Language)
	if err != nil {
		// Handle specific errors like duplicate email
		if strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "unique constraint") {
			utils.LogMessage(utils.LevelWarn, "User already exists")
			utils.LogLineKeyValue(utils.LevelWarn, "Email", newf.Email)
			utils.LogFooter()
			commitOrRollback(tx, err) // Rollback
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "You already have an account"})
		}
		utils.LogMessage(utils.LevelError, "Failed to insert newf")
		utils.LogLineKeyValue(utils.LevelError, "Email", newf.Email)
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		commitOrRollback(tx, err) // Rollback
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Something went wrong creating account"})
	}

	// Add initial 'VERIFYING' role
	addRoleQuery := `
		INSERT INTO newf_roles (email, id_roles)
		VALUES ($1, (SELECT id_roles FROM roles WHERE name = 'VERIFYING'));
	`
	_, err = tx.Exec(addRoleQuery, newf.Email)
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to add initial role")
		utils.LogLineKeyValue(utils.LevelError, "Email", newf.Email)
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		commitOrRollback(tx, err) // Rollback
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Something went wrong setting up account roles"})
	}

	// Generate and store verification code (within the same transaction)
	code := utils.Generate2FACode(6) // Use helper
	setVerificationCodeQuery := `
		UPDATE newf
		SET verification_code = $1, verification_code_expiration = NOW() + INTERVAL '10 minutes'
		WHERE email = $2;
	`
	_, err = tx.Exec(setVerificationCodeQuery, code, newf.Email)
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to set verification code")
		utils.LogLineKeyValue(utils.LevelError, "Email", newf.Email)
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		commitOrRollback(tx, err) // Rollback
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Something went wrong setting up verification"})
	}

	// Subscribe user to all available notification services
	subscribeToAllServicesQuery := `
		INSERT INTO notifications (email, id_services)
		SELECT $1, id_services FROM services;
	`
	_, err = tx.Exec(subscribeToAllServicesQuery, newf.Email)
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to subscribe user to all notification services")
		utils.LogLineKeyValue(utils.LevelError, "Email", newf.Email)
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		// Don't fail registration if subscription fails, just log the error
		// We'll continue with the commit
		utils.LogMessage(utils.LevelWarn, "Registration will proceed despite subscription error")
		err = nil // Reset error so commit proceeds
	} else {
		utils.LogMessage(utils.LevelInfo, "User subscribed to all notification services")
	}

	// Commit the transaction now
	commitErr := tx.Commit()
	if commitErr != nil {
		utils.LogMessage(utils.LevelError, "Failed to commit transaction")
		utils.LogLineKeyValue(utils.LevelError, "Commit Error", commitErr)
		utils.LogFooter()
		// Don't rollback here, commit failed after successful ops
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database error saving registration"})
	}

	// Send verification email *after* successful commit
	if h.EmailService != nil {
		verifCodeData, err := utils.GetVerificationCodeData(h.DB, newf.Email) // Use helper
		if err != nil {
			utils.LogMessage(utils.LevelError, "Failed to get formatted verification code after registration")
			utils.LogLineKeyValue(utils.LevelError, "Email", newf.Email)
			utils.LogLineKeyValue(utils.LevelError, "Error", err)
			// Don't fail the whole request, but log it. Email will lack expiration.
			verifCodeData.VerificationCode = code            // Ensure code is sent
			verifCodeData.VerificationCodeExpiration = "???" // Indicate missing expiration
		}

		// Capture necessary variables for the goroutine
		recipient := newf.Email
		template := "email_templates/email_template_verif_code.html"
		subjectKey := "email_verification.subject"
		lang := newf.Language
		data := verifCodeData

		// Run email sending in a separate goroutine
		go func(recipient, template, subjectKey, lang string, data models.VerificationCodeData) {
			emailErr := h.EmailService.SendEmail(models.Email{
				Recipient:  recipient,
				Template:   template,
				SubjectKey: subjectKey,
				Language:   lang,
			}, data) // Pass struct containing Code and Expiration

			if emailErr != nil {
				utils.LogMessage(utils.LevelError, "Failed to send verification email (asynchronously)")
				utils.LogLineKeyValue(utils.LevelError, "Email", recipient)
				utils.LogLineKeyValue(utils.LevelError, "Error", emailErr)
			}
		}(recipient, template, subjectKey, lang, data)

	} else {
		utils.LogMessage(utils.LevelWarn, "EmailService is not initialized, cannot send verification email.")
	}

	utils.LogMessage(utils.LevelInfo, "Newf registered successfully")
	utils.LogLineKeyValue(utils.LevelInfo, "Email", newf.Email)
	utils.LogFooter()

	return c.SendStatus(fiber.StatusCreated)
}

// Login handles user login.
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var candidate models.Newf // Data from request body

	utils.LogHeader("🔐 Login")

	if err := c.BodyParser(&candidate); err != nil {
		utils.LogMessage(utils.LevelError, "Failed to parse request body")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Failed to parse your data"})
	}

	var storedNewf models.Newf // Data fetched from DB
	var role string
	var languageCode string

	// Fetch user, their role, and language in one query
	query := `
		SELECT n.email, n.password, COALESCE(r.name, 'UNKNOWN') as role, COALESCE(l.code, 'fr') as lang_code
		FROM newf n
		LEFT JOIN newf_roles nr ON n.email = nr.email
		LEFT JOIN roles r ON nr.id_roles = r.id_roles
		LEFT JOIN languages l ON n.language = l.id_languages
		WHERE n.email = $1;
	`
	err := h.DB.QueryRow(query, strings.ToLower(candidate.Email)).Scan(
		&storedNewf.Email,
		&storedNewf.Password,
		&role,
		&languageCode,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			utils.LogMessage(utils.LevelWarn, "Login attempt failed: User not found")
			utils.LogLineKeyValue(utils.LevelWarn, "Email", candidate.Email)
		} else {
			utils.LogMessage(utils.LevelError, "Failed to fetch user during login")
			utils.LogLineKeyValue(utils.LevelError, "Email", candidate.Email)
			utils.LogLineKeyValue(utils.LevelError, "Error", err)
		}
		utils.LogFooter()
		// Return generic error for security
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	// Compare hashed password
	err = bcrypt.CompareHashAndPassword([]byte(storedNewf.Password), []byte(candidate.Password))
	if err != nil {
		utils.LogMessage(utils.LevelWarn, "Login attempt failed: Invalid password")
		utils.LogLineKeyValue(utils.LevelWarn, "Email", candidate.Email)
		utils.LogFooter()
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	// Check if account is validated
	if role == "VERIFYING" {
		utils.LogMessage(utils.LevelWarn, "Login attempt failed: Account not verified")
		utils.LogLineKeyValue(utils.LevelWarn, "Email", candidate.Email)
		utils.LogFooter()
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Validate your account first"})
	}
	if role == "UNKNOWN" {
		// This shouldn't happen if DB constraints are set up correctly
		utils.LogMessage(utils.LevelError, "Login failed: User role is unknown")
		utils.LogLineKeyValue(utils.LevelError, "Email", candidate.Email)
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Account configuration error"})
	}

	// Generate JWT
	token, err := utils.GenerateJWT(storedNewf.Email, role, "") // Use device fingerprint or empty string if not available
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to generate JWT")
		utils.LogLineKeyValue(utils.LevelError, "Email", storedNewf.Email)
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Something went wrong during login"})
	}

	// Send notification email about new sign-in
	if h.EmailService != nil {
		// Capture data for the goroutine
		recipient := storedNewf.Email
		lang := languageCode
		template := "email_templates/email_template_new_signin.html"
		subjectKey := "email_new_signin.subject"

		firstName := ""
		emailParts := strings.Split(storedNewf.Email, "@")
		nameParts := strings.Split(emailParts[0], ".")
		if len(nameParts) > 0 && len(nameParts[0]) > 0 {
			fName := nameParts[0]
			firstName = strings.ToUpper(string(fName[0])) + strings.ToLower(fName[1:])
		}
		emailData := struct {
			FirstName string
			Date      string
			Time      string
		}{
			FirstName: firstName,
			Date:      time.Now().Format("02/01/2006"), // Adjust format as needed
			Time:      time.Now().Format("15:04:05"),   // Adjust format as needed
		}

		// Run email sending in a separate goroutine
		go func(recipient, template, subjectKey, lang string, data interface{}) {
			emailErr := h.EmailService.SendEmail(models.Email{
				Recipient:  recipient,
				Template:   template,
				SubjectKey: subjectKey,
				Language:   lang,
			}, data)
			if emailErr != nil {
				// Log error but don't fail the login
				utils.LogMessage(utils.LevelError, "Failed to send new sign-in notification email asynchronously")
				utils.LogLineKeyValue(utils.LevelError, "Email", recipient)
				utils.LogLineKeyValue(utils.LevelError, "Error", emailErr)
			}
		}(recipient, template, subjectKey, lang, emailData)

	} else {
		utils.LogMessage(utils.LevelWarn, "EmailService is not initialized, skipping new sign-in email.")
	}

	utils.LogMessage(utils.LevelInfo, "Login successful")
	utils.LogLineKeyValue(utils.LevelInfo, "Email", storedNewf.Email)
	utils.LogFooter()

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"token": token})
}

// VerifyAccount handles account verification using a code.
func (h *AuthHandler) VerifyAccount(c *fiber.Ctx) error {
	var req models.VerificationRequest // Assuming a model like { Email string, VerificationCode string }

	utils.LogHeader("📧 Verify Account")

	if err := c.BodyParser(&req); err != nil {
		utils.LogMessage(utils.LevelError, "Failed to parse request body")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Failed to parse your data"})
	}

	if req.Email == "" || req.VerificationCode == "" {
		utils.LogMessage(utils.LevelWarn, "Verify account request missing email or code")
		utils.LogFooter()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Email and verification code are required"})
	}

	// Start transaction
	tx, err := h.DB.Begin()
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to begin transaction")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database error"})
	}
	// Use defer for rollback
	defer func() {
		if err != nil { // Check the named return error
			rbErr := tx.Rollback()
			if rbErr != nil {
				utils.LogMessage(utils.LevelError, "Failed to rollback transaction")
				utils.LogLineKeyValue(utils.LevelError, "Rollback Error", rbErr)
			}
		}
	}()

	// Check code and expiration, and update role if valid
	query := `
		UPDATE newf_roles
		SET id_roles = (SELECT id_roles FROM roles WHERE name = 'NEWF')
		WHERE email = $1
		  AND id_roles = (SELECT id_roles FROM roles WHERE name = 'VERIFYING')
		  AND EXISTS (
			  SELECT 1 FROM newf
			  WHERE email = $1
				AND verification_code = $2
				AND verification_code_expiration > NOW()
		  );
	`
	var result sql.Result
	result, err = tx.Exec(query, strings.ToLower(req.Email), req.VerificationCode)
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to execute verification update")
		utils.LogLineKeyValue(utils.LevelError, "Email", req.Email)
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Something went wrong during verification"})
	}

	var rowsAffected int64
	rowsAffected, err = result.RowsAffected()
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to get rows affected after verification update")
		utils.LogLineKeyValue(utils.LevelError, "Email", req.Email)
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Something went wrong confirming verification"})
	}

	if rowsAffected == 0 {
		// Check if the user/code was ever valid or just expired/wrong code
		var exists int
		// Use the main DB connection for check, no need for tx
		checkQuery := `SELECT COUNT(*) FROM newf WHERE email = $1 AND verification_code = $2`
		errCheck := h.DB.QueryRow(checkQuery, strings.ToLower(req.Email), req.VerificationCode).Scan(&exists)

		// Set err for defer rollback before returning
		err = fmt.Errorf("verification failed")

		if errCheck == nil && exists > 0 {
			utils.LogMessage(utils.LevelWarn, "Verification failed: Code expired")
			utils.LogLineKeyValue(utils.LevelWarn, "Email", req.Email)
			utils.LogFooter()
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Verification code has expired"})
		} else {
			utils.LogMessage(utils.LevelWarn, "Verification failed: Invalid code or email")
			utils.LogLineKeyValue(utils.LevelWarn, "Email", req.Email)
			utils.LogFooter()
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid verification code"})
		}
	}

	// Verification successful, clear the code (optional)
	clearCodeQuery := `UPDATE newf SET verification_code = NULL, verification_code_expiration = NULL WHERE email = $1`
	_, err = tx.Exec(clearCodeQuery, strings.ToLower(req.Email))
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to clear verification code after success")
		utils.LogLineKeyValue(utils.LevelError, "Email", req.Email)
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		// Log but don't fail the request or rollback transaction
		err = nil // Reset error so commit proceeds
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to commit verification transaction")
		utils.LogLineKeyValue(utils.LevelError, "Email", req.Email)
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		// Don't need to rollback here, commit failed
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Something went wrong saving verification"})
	}

	// Generate JWT for immediate login
	var token string
	token, err = utils.GenerateJWT(req.Email, "NEWF", "") // Empty fingerprint for now
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to generate JWT after verification")
		utils.LogLineKeyValue(utils.LevelError, "Email", req.Email)
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		// Don't fail request, account is verified, but advise user to log in manually.
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Account verified successfully. Please log in."})
	}

	// Send welcome email
	if h.EmailService != nil {
		// Capture data for the goroutine
		recipient := req.Email
		languageCode, langErr := utils.GetLanguageCode(h.DB, strings.ToLower(req.Email)) // Use helper
		if langErr != nil {
			utils.LogMessage(utils.LevelWarn, "Failed to get language for welcome email, using default 'fr'")
			languageCode = "fr" // Default
		}
		template := "email_templates/email_template_welcome.html"
		subjectKey := "email_welcome.subject"
		lang := languageCode

		firstName := ""
		emailParts := strings.Split(req.Email, "@")
		nameParts := strings.Split(emailParts[0], ".")
		if len(nameParts) > 0 && len(nameParts[0]) > 0 {
			fName := nameParts[0]
			firstName = strings.ToUpper(string(fName[0])) + strings.ToLower(fName[1:])
		}
		emailData := struct { // Data for welcome email
			FirstName string
		}{
			FirstName: firstName,
		}

		// Run email sending in a separate goroutine
		go func(recipient, template, subjectKey, lang string, data interface{}) {
			emailErr := h.EmailService.SendEmail(models.Email{
				Recipient:  recipient,
				Template:   template,
				SubjectKey: subjectKey,
				Language:   lang,
			}, data)
			if emailErr != nil {
				utils.LogMessage(utils.LevelError, "Failed to send welcome email asynchronously")
				utils.LogLineKeyValue(utils.LevelError, "Email", recipient)
				utils.LogLineKeyValue(utils.LevelError, "Error", emailErr)
				// Log but don't fail
			}
		}(recipient, template, subjectKey, lang, emailData)
	} else {
		utils.LogMessage(utils.LevelWarn, "EmailService is not initialized, skipping welcome email.")
	}

	utils.LogMessage(utils.LevelInfo, "Account verified successfully")
	utils.LogLineKeyValue(utils.LevelInfo, "Email", req.Email)
	utils.LogFooter()

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"token": token}) // Return token for seamless experience
}

// RequestVerificationCode handles requests to resend a verification code.
func (h *AuthHandler) RequestVerificationCode(c *fiber.Ctx) error {
	var req models.EmailRequest // Assuming { Email string }

	utils.LogHeader("📧 Resend Verification Code")

	if err := c.BodyParser(&req); err != nil {
		utils.LogMessage(utils.LevelError, "Failed to parse request body")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Failed to parse your data"})
	}

	if req.Email == "" {
		utils.LogMessage(utils.LevelWarn, "Resend code request missing email")
		utils.LogFooter()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Email is required"})
	}

	// Check if user exists and potentially if they are already verified
	var isVerified bool
	var userExists bool
	checkQuery := `
		SELECT EXISTS (SELECT 1 FROM newf WHERE email = $1),
		       EXISTS (
				SELECT 1 FROM newf_roles nr JOIN roles r ON nr.id_roles = r.id_roles
				WHERE nr.email = $1 AND r.name != 'VERIFYING'
			);
	`
	err := h.DB.QueryRow(checkQuery, strings.ToLower(req.Email)).Scan(&userExists, &isVerified)
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed check user status for verification code request")
		utils.LogLineKeyValue(utils.LevelError, "Email", req.Email)
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database error checking user status"})
	}

	if !userExists {
		// User doesn't exist at all
		utils.LogMessage(utils.LevelWarn, "Verification code request for non-existent user")
		utils.LogLineKeyValue(utils.LevelWarn, "Email", req.Email)
		utils.LogFooter()
		// Return generic message for security
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "If an account with this email exists, a verification code has been sent."}) // Ambiguous for security
	}
	// If user exists but is already verified, we might still allow sending code for password reset.
	// The logic here assumes this endpoint can be used for both initial verification and password reset request.
	// if isVerified {
	// 	utils.LogMessage(utils.LevelInfo, "Verification code request for already verified user (potential password reset)")
	//     utils.LogLineKeyValue(utils.LevelInfo, "Email", req.Email)
	// }

	// Generate and update code
	code := utils.Generate2FACode(6)
	updateQuery := `
		UPDATE newf
		SET verification_code = $1,
			verification_code_expiration = NOW() + INTERVAL '10 minutes'
		WHERE email = $2;
	`
	result, err := h.DB.Exec(updateQuery, code, req.Email)
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to update verification code")
		utils.LogLineKeyValue(utils.LevelError, "Email", req.Email)
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update verification code"})
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		// Should not happen if user existence check passed, but good practice to check
		utils.LogMessage(utils.LevelError, "Failed to update verification code (no rows affected)")
		utils.LogLineKeyValue(utils.LevelError, "Email", req.Email)
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update verification code for user"})
	}

	// Get formatted code data for email
	verifCodeData, codeErr := utils.GetVerificationCodeData(h.DB, req.Email)
	if codeErr != nil {
		utils.LogMessage(utils.LevelError, "Failed to get formatted verification code")
		utils.LogLineKeyValue(utils.LevelError, "Email", req.Email)
		utils.LogLineKeyValue(utils.LevelError, "Error", codeErr)
		// Proceed without formatted data? Send raw code?
		// For now, log and continue, email might lack expiration time.
		verifCodeData.VerificationCode = code // Ensure code is sent at least
		verifCodeData.VerificationCodeExpiration = "???"
	}

	// Get user language
	languageCode, langErr := utils.GetLanguageCode(h.DB, req.Email)
	if langErr != nil {
		utils.LogMessage(utils.LevelWarn, "Failed to get language for verification email, using default 'fr'")
		languageCode = "fr" // Default
	}

	// Send email
	if h.EmailService != nil {
		// Capture data for the goroutine
		recipient := req.Email
		template := "email_templates/email_template_verif_code.html"
		subjectKey := "email_verification.subject" // Or a different key for password reset if applicable
		lang := languageCode
		data := verifCodeData

		// Run email sending in a separate goroutine
		go func(recipient, template, subjectKey, lang string, data models.VerificationCodeData) {
			emailErr := h.EmailService.SendEmail(models.Email{
				Recipient:  recipient,
				Template:   template,
				SubjectKey: subjectKey,
				Language:   lang,
			}, data)

			if emailErr != nil {
				utils.LogMessage(utils.LevelError, "Failed to send verification code email asynchronously")
				utils.LogLineKeyValue(utils.LevelError, "Email", recipient)
				utils.LogLineKeyValue(utils.LevelError, "Error", emailErr)
				// Don't expose email failure to client if code was updated successfully in DB
			}
		}(recipient, template, subjectKey, lang, data)
	} else {
		utils.LogMessage(utils.LevelWarn, "EmailService is not initialized, cannot send verification code email.")
		// Return an error? Or just log? Depends on requirements.
		// If email is critical, maybe return an error.
		// return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Email service unavailable"})
	}

	utils.LogMessage(utils.LevelInfo, "Verification code sent request processed")
	utils.LogLineKeyValue(utils.LevelInfo, "Email", req.Email)
	// utils.LogLineKeyValue(utils.LevelInfo, "Code", code) // Avoid logging code in production
	utils.LogFooter()

	// Return generic success message for security (don't confirm if email exists)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "If an account with this email exists, a verification code has been sent."})
}

// ChangePassword handles password change requests (both verified and unverified users).
func (h *AuthHandler) ChangePassword(c *fiber.Ctx) error {
	var req models.ChangePasswordRequest // Needs fields like Email, OldPassword (optional), NewPassword, NewPasswordConfirmation, VerificationCode (optional)

	utils.LogHeader("🔐 Change Password")

	if err := c.BodyParser(&req); err != nil {
		utils.LogMessage(utils.LevelError, "Failed to parse request body")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Failed to parse your data"})
	}

	// Basic validation
	if req.Email == "" {
		utils.LogMessage(utils.LevelWarn, "Password change request missing email")
		utils.LogFooter()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Email is required"})
	}
	if req.NewPassword == "" || req.NewPasswordConfirmation == "" {
		utils.LogMessage(utils.LevelWarn, "Password change request missing new password fields")
		utils.LogLineKeyValue(utils.LevelWarn, "Email", req.Email)
		utils.LogFooter()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "New password fields are required"})
	}
	if req.NewPassword != req.NewPasswordConfirmation {
		utils.LogMessage(utils.LevelWarn, "Password change request passwords do not match")
		utils.LogLineKeyValue(utils.LevelWarn, "Email", req.Email)
		utils.LogFooter()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Passwords do not match"})
	}
	// Add password strength validation if needed

	// Hash the new password
	hashedNewPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to hash new password")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Something went wrong hashing password"})
	}

	// Determine update logic based on provided fields (Verification Code vs Old Password)
	var rowsAffected int64
	var updateErr error

	if req.VerificationCode != "" {
		// --- Change password using verification code (e.g., forgot password flow) ---
		utils.LogMessage(utils.LevelInfo, "Attempting password change using verification code")
		utils.LogLineKeyValue(utils.LevelInfo, "Email", req.Email)

		query := `
			UPDATE newf
			SET password = $1,
				password_updated_date = NOW(),
				verification_code = NULL,          -- Clear code after use
				verification_code_expiration = NULL
			WHERE email = $2
			  AND verification_code = $3
			  AND verification_code_expiration > NOW();
		`
		result, err := h.DB.Exec(query, string(hashedNewPassword), strings.ToLower(req.Email), req.VerificationCode)
		if err != nil {
			updateErr = err // Store error to log later
		} else {
			rowsAffected, _ = result.RowsAffected()
		}

		if updateErr != nil || rowsAffected == 0 {
			if updateErr != nil {
				utils.LogMessage(utils.LevelError, "DB error during password change with code")
				utils.LogLineKeyValue(utils.LevelError, "Error", updateErr)
			} else {
				utils.LogMessage(utils.LevelWarn, "Password change with code failed (invalid/expired code or email)")
			}
			utils.LogLineKeyValue(utils.LevelWarn, "Email", req.Email)
			utils.LogFooter()
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid or expired verification code"})
		}

	} else if req.OldPassword != "" {
		// --- Change password using old password (logged-in user flow) ---
		utils.LogMessage(utils.LevelInfo, "Attempting password change using old password")
		utils.LogLineKeyValue(utils.LevelInfo, "Email", req.Email)

		// 1. Fetch current password hash
		var currentPasswordHash string
		fetchQuery := `SELECT password FROM newf WHERE email = $1;`
		err := h.DB.QueryRow(fetchQuery, req.Email).Scan(&currentPasswordHash)
		if err != nil {
			if err == sql.ErrNoRows {
				utils.LogMessage(utils.LevelWarn, "Password change failed: User not found")
			} else {
				utils.LogMessage(utils.LevelError, "Failed to fetch user for password change")
				utils.LogLineKeyValue(utils.LevelError, "Error", err)
			}
			utils.LogLineKeyValue(utils.LevelWarn, "Email", req.Email)
			utils.LogFooter()
			// Return generic error for security if user not found
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
		}

		// 2. Compare old password
		err = bcrypt.CompareHashAndPassword([]byte(currentPasswordHash), []byte(req.OldPassword))
		if err != nil {
			utils.LogMessage(utils.LevelWarn, "Password change failed: Incorrect old password")
			utils.LogLineKeyValue(utils.LevelWarn, "Email", req.Email)
			utils.LogFooter()
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
		}

		// 3. Update password
		updateQuery := `
			UPDATE newf
			SET password = $1,
				password_updated_date = NOW()
			WHERE email = $2;
		`
		result, err := h.DB.Exec(updateQuery, string(hashedNewPassword), req.Email)
		if err != nil {
			updateErr = err // Store error
		} else {
			rowsAffected, _ = result.RowsAffected() // Should always be 1 if user exists
		}

		if updateErr != nil || rowsAffected == 0 {
			if updateErr != nil {
				utils.LogMessage(utils.LevelError, "DB error during password change with old password")
				utils.LogLineKeyValue(utils.LevelError, "Error", updateErr)
			} else {
				utils.LogMessage(utils.LevelError, "Password change with old password failed (no rows affected, user should exist)")
			}
			utils.LogLineKeyValue(utils.LevelError, "Email", req.Email)
			utils.LogFooter()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update password"})
		}

	} else {
		// --- Invalid request: Neither verification code nor old password provided ---
		utils.LogMessage(utils.LevelWarn, "Password change request missing verification code or old password")
		utils.LogLineKeyValue(utils.LevelWarn, "Email", req.Email)
		utils.LogFooter()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Either verification code or old password is required"})
	}

	// Password successfully changed
	utils.LogMessage(utils.LevelInfo, "Password changed successfully")
	utils.LogLineKeyValue(utils.LevelInfo, "Email", req.Email)
	utils.LogFooter()

	// Optionally send a confirmation email
	if h.EmailService != nil {
		// Capture data for the goroutine
		recipient := req.Email
		languageCode, langErr := utils.GetLanguageCode(h.DB, req.Email)
		if langErr != nil {
			utils.LogMessage(utils.LevelWarn, "Failed to get language for password change confirmation email, using default 'fr'")
			languageCode = "fr" // Default
		}
		template := "email_templates/email_template_password_changed.html" // Create this template
		subjectKey := "email_password_changed.subject"                     // Define this translation key
		lang := languageCode

		firstName := ""
		emailParts := strings.Split(req.Email, "@")
		nameParts := strings.Split(emailParts[0], ".")
		if len(nameParts) > 0 && len(nameParts[0]) > 0 {
			fName := nameParts[0]
			firstName = strings.ToUpper(string(fName[0])) + strings.ToLower(fName[1:])
		}
		emailData := struct { // Data for the template
			FirstName string
			Date      string
			Time      string
		}{
			FirstName: firstName,
			Date:      time.Now().Format("02/01/2006"),
			Time:      time.Now().Format("15:04:05"),
		}

		// Run email sending in a separate goroutine
		go func(recipient, template, subjectKey, lang string, data interface{}) {
			emailErr := h.EmailService.SendEmail(models.Email{
				Recipient:  recipient,
				Template:   template,
				SubjectKey: subjectKey,
				Language:   lang,
			}, data)
			if emailErr != nil {
				utils.LogMessage(utils.LevelError, "Failed to send password change confirmation email asynchronously")
				utils.LogLineKeyValue(utils.LevelError, "Email", recipient)
				utils.LogLineKeyValue(utils.LevelError, "Error", emailErr)
			}
		}(recipient, template, subjectKey, lang, emailData)
	} else {
		utils.LogMessage(utils.LevelWarn, "EmailService not initialized, skipping password change confirmation email.")
	}

	return c.SendStatus(fiber.StatusOK)
}
