package main

import (
	"database/sql"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"log"
	"strings"
	"time"
)

func register(c *fiber.Ctx) error {
	var newf Newf

	log.Println("╔======== 👶 Newf Registration 👶 ========╗")

	if err := c.BodyParser(&newf); err != nil {
		log.Println("║ 💥 Failed to parse request body: ", err)
		log.Println("╚=========================================╝")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Failed to parse your data"})
	}
	emailValid, err := checkEmail(newf.Email)
	if err != nil {
		log.Println("║ 💥 Failed to check email: ", err)
		log.Println("╚=========================================╝")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "There is an error with your email"})
	}
	if !emailValid {
		log.Println("║ 💥 Invalid email: ", newf.Email)
		log.Println("╚=========================================╝")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid email"})
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

	stmt, err := db.Prepare(request)
	if err != nil {
		log.Println("║ 💥 Failed to prepare statement: ", err)
		log.Println("╚=========================================╝")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Something went wrong"})
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			log.Println("║ 💥 Failed to close statement: ", err)
			log.Println("╚=========================================╝")
			return
		}
	}(stmt)

	_, err = stmt.Exec(newf.Email, newf.Password, newf.FirstName, newf.LastName)
	if err != nil {
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

	stmt, err = db.Prepare(addRole)
	if err != nil {
		log.Println("║ 💥 Failed to prepare statement: ", err)
		log.Println("╚=========================================╝")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Something went wrong"})
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			log.Println("║ 💥 Failed to close statement: ", err)
			log.Println("╚=========================================╝")
			return
		}
	}(stmt)

	_, err = stmt.Exec(newf.Email)
	if err != nil {
		log.Println("║ 💥 Failed to add role: ", err)
		log.Println("╚=========================================╝")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Something went wrong"})
	}

	verifCode, err := getVerificationCode(newf)
	if err != nil {
		log.Println("║ 💥 Failed to get verification code: ", err)
		log.Println("╚=========================================╝")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Something went wrong"})
	}

	errEmail := sendEmail(Email{
		Recipient: newf.Email,
		Subject:   fmt.Sprintf("🔐 Ton code de vérification : %s", verifCode.VerificationCode),
		Template:  "email_templates/email_template_verif_code.html",
		Sender: EmailSender{
			Name:  "Transat Team",
			Email: "admin@destimt.fr",
		},
	}, verifCode)
	if errEmail != nil {
		log.Println("║ 💥 Failed to send verification email: ", errEmail)
		log.Println("╚=========================================╝")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Something went wrong"})
	}

	log.Println("║ ✅ Newf registered successfully")
	log.Println("║ 📧 Email: ", newf.Email)
	log.Println("╚=========================================╝")

	return c.SendStatus(fiber.StatusCreated)
}

func login(c *fiber.Ctx) error {
	var newf Newf
	var candidate Newf

	log.Println("╔============== 🔐 Login 🔐 ============╗")

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
	stmt, err := db.Prepare(request)
	if err != nil {
		log.Println("║ 💥 Failed to prepare statement: ", err)
		log.Println("╚=========================================╝")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Something went wrong"})
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			log.Println("║ 💥 Failed to close statement: ", err)
			log.Println("╚=========================================╝")
			return
		}
	}(stmt)

	row := stmt.QueryRow(candidate.Email)
	err = row.Scan(&newf.Email, &newf.Password)
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

	if !isValidated(candidate) {
		log.Println("║ 💥 Newf is not validated")
		log.Println("╚=========================================╝")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Validate your account first"})
	}

	token, err := generateJWT(newf)
	if err != nil {
		log.Println("║ 💥 Failed to generate JWT: ", err)
		log.Println("╚=========================================╝")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Something went wrong"})
	}

	errEmail := sendEmail(Email{
		Recipient: newf.Email,
		Subject:   "🔐 Nouvelle connexion à votre compte",
		Template:  "email_templates/email_template_new_signin.html",
		Sender: EmailSender{
			Name:  "Transat Team",
			Email: "admin@destimt.fr",
		},
	}, struct {
		FirstName string
		Date      string
		Time      string
	}{
		FirstName: strings.ToUpper(strings.Split(newf.Email, ".")[0])[0:1] + strings.Split(newf.Email, ".")[0][1:],
		Date:      time.Now().Format("02/01/2006"),
		Time:      time.Now().Format("15:04:05"),
	})
	if errEmail != nil {
		log.Println("║ 💥 Failed to send login email: ", errEmail)
		log.Println("╚=========================================╝")
	}

	log.Println("║ ✅ Login successful")
	log.Println("║ 📧 Email: ", newf.Email)
	log.Println("╚=========================================╝")

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"token": token})
}

func verificationCode(c *fiber.Ctx) error {
	var newf Newf

	log.Println("╔======== 📧 Verification Code 📧 ========╗")

	if err := c.BodyParser(&newf); err != nil {
		log.Println("║ 💥 Failed to parse request body: ", err)
		log.Println("╚=========================================╝")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Failed to parse your data"})
	}

	request := `
		UPDATE newf
		SET 
		    verification_code = $1,
		    verification_code_expiration = CURRENT_TIMESTAMP + INTERVAL '10 minutes'
		WHERE email = $2;
	`

	stmt, err := db.Prepare(request)
	if err != nil {
		log.Println("║ 💥 Failed to prepare statement: ", err)
		log.Println("╚=========================================╝")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Something went wrong"})
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			log.Println("║ 💥 Failed to close statement: ", err)
			log.Println("╚=========================================╝")
			return
		}
	}(stmt)

	code := generate2FACode(6)

	_, err = stmt.Exec(code, newf.Email)
	if err != nil {
		log.Println("║ 💥 Failed to send the code: ", err)
		log.Println("║ 📧 Email: ", newf.Email)
		log.Println("╚=========================================╝")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Something went wrong"})
	}

	verifCode, err := getVerificationCode(newf)
	if err != nil {
		log.Println("║ 💥 Failed to get verification code: ", err)
		log.Println("╚=========================================╝")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Something went wrong"})
	}

	errEmail := sendEmail(Email{
		Recipient: newf.Email,
		Subject:   fmt.Sprintf("🔐 Ton code de vérification : %s", verifCode.VerificationCode),
		Template:  "email_templates/email_template_verif_code.html",
		Sender: EmailSender{
			Name:  "Transat Team",
			Email: "admin@destimt.fr",
		},
	}, verifCode)
	if errEmail != nil {
		log.Println("║ 💥 Failed to send verification email: ", errEmail)
		log.Println("╚=========================================╝")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Something went wrong"})
	}

	log.Println("║ ✅ Verification code sent")
	log.Println("║ 📧 Email: ", newf.Email)
	log.Println("║ 📧 Verification Code: ", code)
	log.Println("╚=========================================╝")

	return c.SendStatus(fiber.StatusOK)
}

func changePassword(c *fiber.Ctx) error {
	var newf Newf

	log.Println("╔======== 🔐 Change Password 🔐 ========╗")

	if err := c.BodyParser(&newf); err != nil {
		log.Println("║ 💥 Failed to parse request body: ", err)
		log.Println("╚=========================================╝")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Failed to parse your data"})
	}

	if newf.NewPassword == "" || newf.NewPasswordConfirmation == "" {
		log.Println("║ 💥 New passwords are missing")
		log.Println("╚=========================================╝")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "New passwords are missing"})
	}

	if newf.NewPassword != newf.NewPasswordConfirmation {
		log.Println("║ 💥 Passwords do not match")
		log.Println("╚=========================================╝")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Passwords do not match"})
	}

	if newf.VerificationCode != "" {
		request := `
			UPDATE newf
			SET password = $1,
			    password_updated_date = CURRENT_TIMESTAMP
			WHERE email = $2
			AND verification_code = $3
			AND verification_code_expiration > NOW();
		`
		stmt, err := db.Prepare(request)
		if err != nil {
			log.Println("║ 💥 Failed to prepare statement: ", err)
			log.Println("╚=========================================╝")
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Something went wrong"})
		}
		defer func(stmt *sql.Stmt) {
			err := stmt.Close()
			if err != nil {
				log.Println("║ 💥 Failed to close statement: ", err)
				log.Println("╚=========================================╝")
				return
			}
		}(stmt)

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newf.NewPassword), bcrypt.DefaultCost)
		if err != nil {
			log.Println("║ 💥 Failed to hash password: ", err)
			log.Println("╚=========================================╝")
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Something went wrong"})
		}

		result, err := stmt.Exec(string(hashedPassword), newf.Email, newf.VerificationCode)
		if err != nil {
			log.Println("║ 💥 Failed to change password: ", err)
			log.Println("║ 📧 Email: ", newf.Email)
			log.Println("╚=========================================╝")
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Something went wrong"})
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			log.Println("║ 💥 Failed to get affected rows: ", err)
			log.Println("╚=========================================╝")
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Something went wrong"})
		}

		if rowsAffected == 0 {
			log.Println("║ 💥 No rows were updated")
			log.Println("╚=========================================╝")
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid verification code or expired"})
		}
	}

	if newf.Password != "" {
		requestOldPassword := `
			SELECT password
			FROM newf
			WHERE email = $1;
		`
		stmt, err := db.Prepare(requestOldPassword)
		if err != nil {
			log.Println("║ 💥 Failed to prepare statement: ", err)
			log.Println("╚=========================================╝")
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Something went wrong"})
		}
		defer func(stmt *sql.Stmt) {
			err := stmt.Close()
			if err != nil {
				log.Println("║ 💥 Failed to close statement: ", err)
				log.Println("╚=========================================╝")
				return
			}
		}(stmt)

		var oldPassword string
		err = stmt.QueryRow(newf.Email).Scan(&oldPassword)
		if err != nil {
			log.Println("║ 💥 Failed to get old password: ", err)
			log.Println("╚=========================================╝")
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Something went wrong"})
		}

		err = bcrypt.CompareHashAndPassword([]byte(oldPassword), []byte(newf.Password))
		if err != nil {
			log.Println("║ 💥 Failed to compare password: ", err)
			log.Println("╚=========================================╝")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
		}

		request := `
			UPDATE newf
			SET password = $1,
			    password_updated_date = CURRENT_TIMESTAMP
			WHERE email = $2;
		`
		stmt, err = db.Prepare(request)
		if err != nil {
			log.Println("║ 💥 Failed to prepare statement: ", err)
			log.Println("╚=========================================╝")
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Something went wrong"})
		}
		defer func(stmt *sql.Stmt) {
			err := stmt.Close()
			if err != nil {
				log.Println("║ 💥 Failed to close statement: ", err)
				log.Println("╚=========================================╝")
				return
			}
		}(stmt)

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newf.NewPassword), bcrypt.DefaultCost)
		if err != nil {
			log.Println("║ 💥 Failed to hash password: ", err)
			log.Println("╚=========================================╝")
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Something went wrong"})
		}

		result, err := stmt.Exec(string(hashedPassword), newf.Email)
		if err != nil {
			log.Println("║ 💥 Failed to change password: ", err)
			log.Println("║ 📧 Email: ", newf.Email)
			log.Println("╚=========================================╝")
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Something went wrong"})
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			log.Println("║ 💥 Failed to get affected rows: ", err)
			log.Println("╚=========================================╝")
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Something went wrong"})
		}

		if rowsAffected == 0 {
			log.Println("║ 💥 No rows were updated")
			log.Println("╚=========================================╝")
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid verification code or expired"})
		}
	}

	log.Println("║ ✅ Password changed successfully")
	log.Println("║ 📧 Email: ", newf.Email)
	log.Println("╚=========================================╝")

	return c.SendStatus(fiber.StatusOK)
}

func deleteNewf(c *fiber.Ctx) error {
	email := c.Params("email")

	log.Println("╔======== 🚫 Delete Newf 🚫 ========╗")

	request := `
		DELETE FROM newf
		WHERE email = $1;
	`

	stmt, err := db.Prepare(request)
	if err != nil {
		log.Println("║ 💥 Failed to prepare statement: ", err)
		log.Println("╚=========================================╝")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Something went wrong"})
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			log.Println("║ 💥 Failed to close statement: ", err)
			log.Println("╚=========================================╝")
			return
		}
	}(stmt)

	_, err = stmt.Exec(email)
	if err != nil {
		log.Println("║ 💥 Failed to delete newf: ", err)
		log.Println("║ 📧 Email: ", email)
		log.Println("╚=========================================╝")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Something went wrong"})
	}

	log.Println("║ ✅ Newf deleted successfully")
	log.Println("║ 📧 Email: ", email)
	log.Println("╚=========================================╝")

	return c.SendStatus(fiber.StatusOK)
}
