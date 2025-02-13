package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"log"
)

func verifyAccount(c *fiber.Ctx) error {
	var newf Newf

	log.Println("╔======== 📧 Verify Account 📧 ========╗")

	if err := c.BodyParser(&newf); err != nil {
		log.Println("║ 💥 Failed to parse request body: ", err)
		log.Println("╚=======================================╝")
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
		log.Println("║ 💥 Failed to prepare statement: ", err)
		log.Println("╚=======================================╝")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Something went wrong"})
	}
	defer stmt.Close()

	res, err := stmt.Exec(newf.Email, newf.VerificationCode)
	if err != nil {
		log.Println("║ 💥 Failed to verify account: ", err)
		log.Println("║ 📧 Email: ", newf.Email)
		log.Println("╚=======================================╝")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Something went wrong"})
	}

	rows, err := res.RowsAffected()
	if err != nil {
		log.Println("║ 💥 Failed to get rows affected: ", err)
		log.Println("║ 📧 Email: ", newf.Email)
		log.Println("╚=======================================╝")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Something went wrong"})
	}

	if rows == 0 {
		log.Println("║ 💥 Invalid verification code or expired: ", newf.Email)
		log.Println("╚=======================================╝")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid verification code or expired"})
	}

	// TODO: Send welcome email
	token, err := generateJWT(newf)
	if err != nil {
		log.Println("║ 💥 Failed to generate JWT: ", err)
		log.Println("╚=========================================╝")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Something went wrong"})
	}

	log.Println("║ ✅ Account verified successfully")
	log.Println("║ 📧 Email: ", newf.Email)
	log.Println("╚=======================================╝")

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"token": token})
}

func jwtMiddleware(c *fiber.Ctx) error {
	tokenString := c.Get("Authorization")

	log.Println("╔======== 📧 JWT Middleware 📧 ========╗")

	if tokenString == "" {
		log.Println("║ 💥 Missing token")
		log.Println("╚=======================================╝")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing token"})
	}

	token, err := validateJWT(tokenString)
	if err != nil {
		log.Println("║ 💥 Invalid token: ", err)
		log.Println("╚=======================================╝")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
	}

	// Extract claims (optional)
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		log.Println("║ 💥 Invalid claims")
		log.Println("╚=======================================╝")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid claims"})
	}

	log.Println("║ ✅ Token is valid")
	log.Println("║ 📧 Email: ", claims["email"])
	log.Println("╚=======================================╝")

	return c.Next()
}
