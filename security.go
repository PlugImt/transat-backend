package main

import (
	"github.com/gofiber/fiber/v2"
	"log"
)

func verifyAccount(c *fiber.Ctx) error {
	var newf Newf

	log.Println("╔======== 📧 Verify Account 📧 ========╗")

	if err := c.BodyParser(&newf); err != nil {
		log.Println("║ 💥 Failed to parse request body: ", err)
		log.Println("╚=========================================╝")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Failed to parse your data"})
	}

	request := `
		UPDATE newf_roles
		SET id_roles = (SELECT id_roles FROM roles WHERE name = 'NEWF')
		WHERE email = $1 AND id_roles = (SELECT id_roles FROM roles WHERE name = 'VERIFYING')
  			AND (SELECT verification_code FROM newf WHERE email = $1) = $2
  			AND (SELECT verification_code_expiration FROM newf WHERE email = $1) > NOW();;
		;
	`
	// Execute the SQL query
	res, err := db.Exec(request, newf.Email, newf.VerificationCode)
	if err != nil {
		log.Println("║ 💥 Failed to verify account: ", err)
		log.Println("║ 📧 Email: ", newf.Email)
		log.Println("╚=========================================╝")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Something went wrong"})
	}

	rows, err := res.RowsAffected()
	if err != nil {
		log.Println("║ 💥 Failed to get rows affected: ", err)
		log.Println("║ 📧 Email: ", newf.Email)
		log.Println("╚=========================================╝")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Something went wrong"})
	}

	if rows == 0 {
		log.Println("║ 💥 Failed to verify account: ", err)
		log.Println("║ 📧 Email: ", newf.Email)
		log.Println("╚=========================================╝")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid verification code"})
	}

	log.Println("║ ✅ Account verified successfully")
	log.Println("║ 📧 Email: ", newf.Email)
	log.Println("╚=========================================╝")

	return c.SendStatus(fiber.StatusOK)
}
