package main

import (
	"os"
	"strings"

	"Transat_2.0_Backend/models"
	"Transat_2.0_Backend/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func verifyAccount(c *fiber.Ctx) error {
	var newf models.Newf

	utils.LogHeader("📧 Verify Account")

	if err := c.BodyParser(&newf); err != nil {
		utils.LogMessage(utils.LevelError, "Failed to parse request body")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Failed to parse your data"})
	}

	request := `
		UPDATE newf_roles
		SET id_roles = (SELECT id_roles FROM roles WHERE name = 'NEWF')
		WHERE email = $1 
		AND id_roles = (SELECT id_roles FROM roles WHERE name = 'VERIFYING')
		AND (SELECT verification_code FROM newf WHERE email = $1) = $2
		AND (SELECT verification_code_expiration FROM newf WHERE email = $1) > NOW();
	`

	stmt, err := db.Prepare(request)
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to prepare statement")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Something went wrong"})
	}
	defer stmt.Close()

	res, err := stmt.Exec(newf.Email, newf.VerificationCode)
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to verify account")
		utils.LogLineKeyValue(utils.LevelError, "Email", newf.Email)
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Something went wrong"})
	}

	rows, err := res.RowsAffected()
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to get rows affected")
		utils.LogLineKeyValue(utils.LevelError, "Email", newf.Email)
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Something went wrong"})
	}

	if rows == 0 {
		utils.LogMessage(utils.LevelError, "Invalid verification code or expired")
		utils.LogLineKeyValue(utils.LevelError, "Email", newf.Email)
		utils.LogFooter()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid verification code or expired"})
	}

	token, err := generateJWT(newf)
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to generate JWT")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Something went wrong"})
	}

	// get language from newf table
	language, err := GetLanguage(newf.Email)
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to get language")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Something went wrong"})
	}

	errEmail := sendEmail(models.Email{
		Recipient: newf.Email,
		Template:  "email_templates/email_template_welcome.html",
		Sender: models.EmailSender{
			Name:  "Transat Team",
			Email: os.Getenv("EMAIL_SENDER"),
		},
		Language: language,
	}, struct {
		FirstName string
	}{
		FirstName: strings.ToUpper(strings.Split(newf.Email, ".")[0])[0:1] + strings.Split(newf.Email, ".")[0][1:],
	})
	if errEmail != nil {
		utils.LogMessage(utils.LevelError, "Failed to send welcome email")
		utils.LogLineKeyValue(utils.LevelError, "Error", errEmail)
	}

	utils.LogMessage(utils.LevelInfo, "Account verified successfully")
	utils.LogLineKeyValue(utils.LevelInfo, "Email", newf.Email)
	utils.LogFooter()

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"token": token})
}

func jwtMiddleware(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")

	utils.LogHeader("📧 JWT Middleware")

	if authHeader == "" {
		utils.LogMessage(utils.LevelError, "Missing token")
		utils.LogFooter()
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing token"})
	}

	tokenString := authHeader
	if strings.HasPrefix(authHeader, "Bearer ") {
		tokenString = authHeader[7:]
	}

	token, err := validateJWT(tokenString)
	if err != nil {
		utils.LogMessage(utils.LevelError, "Invalid token")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		utils.LogMessage(utils.LevelError, "Invalid claims")
		utils.LogFooter()
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid claims"})
	}

	c.Locals("email", claims["email"])

	utils.LogMessage(utils.LevelInfo, "Token is valid")
	utils.LogLineKeyValue(utils.LevelInfo, "Email", claims["email"])
	utils.LogFooter()

	return c.Next()
}
