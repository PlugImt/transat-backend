package main

import (
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"log"
	"strings"
)

func register(c *fiber.Ctx) error {
	var newf Newf

	log.Println("╔======== 👶 Newf Registration 👶 ========╗")

	if err := c.BodyParser(&newf); err != nil {
		log.Println("║ 💥 Failed to parse request body: ", err)
		log.Println("╚=========================================╝")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Failed to parse your data"})
	}
	_, err := checkEmail(newf.Email)
	if err != nil {
		log.Println("║ 💥 Failed to check email: ", err)
		log.Println("╚=========================================╝")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "There is an error with your email"})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newf.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("║ 💥 Failed to hash password: ", err)
		log.Println("╚=========================================╝")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Something went wrong"})
	}

	newf.Password = string(hashedPassword)
	firstName := strings.Split(newf.Email, ".")[0]
	lastName := strings.Split(strings.Split(newf.Email, ".")[1], "@")[0]
	newf.FirstName = strings.ToUpper(string(firstName[0])) + firstName[1:]
	newf.LastName = strings.ToUpper(lastName)

	request := `
		INSERT INTO newf (email, password, first_name, last_name)
		VALUES ($1, $2, $3, $4);
	`
	// Execute the SQL query
	_, err = db.Exec(request, newf.Email, newf.Password, newf.FirstName, newf.LastName)
	if err != nil {
		// Check if the error is a unique violation (PostgreSQL error code for unique violation)
		if strings.Contains(err.Error(), "duplicate key") {
			log.Println("║ 💥 User already exists: ", err)
			log.Println("╚=========================================╝")
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "You already have an account"})
		}

		log.Println("║ 💥 Failed to insert newf: ", err)
		log.Println("╚=========================================╝")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Something went wrong"})
	}

	addRole := `
		INSERT INTO newf_roles (email, id_roles)
		VALUES ($1, (SELECT id_roles FROM roles WHERE name = 'VERIFYING'));
	`
	_, err = db.Exec(addRole, newf.Email)
	if err != nil {
		log.Println("║ 💥 Failed to add role: ", err)
		log.Println("╚=========================================╝")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Something went wrong"})
	}

	// TODO: Send verification email

	log.Println("║ ✅ Newf registered successfully")
	log.Println("║ 📧 Email: ", newf.Email)
	log.Println("╚=========================================╝")

	return c.SendStatus(fiber.StatusCreated)
}

func login(c *fiber.Ctx) error {
	var newf Newf
	var candidate Newf

	log.Println("╔======== 🔐 Login 🔐 ========╗")

	if err := c.BodyParser(&candidate); err != nil {
		log.Println("║ 💥 Failed to parse request body: ", err)
		log.Println("╚=========================================╝")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Failed to parse your data"})
	}

	request := `
		SELECT email, password
		FROM newf
		WHERE email = $1;
	`
	// Execute the SQL query
	row := db.QueryRow(request, candidate.Email)
	err := row.Scan(&newf.Email, &newf.Password)
	if err != nil {
		log.Println("║ 💥 Failed to fetch newf: ", err)
		log.Println("║ 📧 Email: ", newf.Email)
		log.Println("╚=========================================╝")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	err = bcrypt.CompareHashAndPassword([]byte(newf.Password), []byte(candidate.Password))
	if err != nil {
		log.Println("║ 💥 Failed to compare password: ", err)
		log.Println("║ 📧 Email: ", newf.Email)
		log.Println("╚=========================================╝")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	// TODO: Generate JWT token and send it to the user

	log.Println("║ ✅ Login successful")
	log.Println("║ 📧 Email: ", newf.Email)
	log.Println("╚=========================================╝")

	return c.SendStatus(fiber.StatusOK)
}
