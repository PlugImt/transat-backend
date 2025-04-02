package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
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
		INSERT INTO newf (email, password, first_name, last_name, language)
		VALUES ($1, $2, $3, $4, COALESCE ((SELECT id_languages FROM languages WHERE name = $5 OR code = $5 LIMIT 1), 1));
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

	if newf.Language == "" {
		newf.Language = "fr"
	} else {
		newf.Language = strings.ToLower(newf.Language)
	}

	fmt.Println(newf)

	_, err = stmt.Exec(newf.Email, newf.Password, newf.FirstName, newf.LastName, newf.Language)
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
			Email: os.Getenv("EMAIL_SENDER"),
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
			Email: os.Getenv("EMAIL_SENDER"),
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
			Email: os.Getenv("EMAIL_SENDER"),
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
	// TODO: change this to use the JWT token
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

func getNewf(c *fiber.Ctx) error {
	email := c.Locals("email").(string)

	log.Println("╔======== 📧 Get Newf 📧 ========╗")

	request := `
			SELECT id_newf,
       			email,
       			first_name,
       			last_name,
       			COALESCE(profile_picture, '')                              as profile_picture,
       			COALESCE(phone_number, '')                                 as phone_number,
       			COALESCE(graduation_year, 0)                               as graduation_year,
       			COALESCE(campus, '')                                       as campus,
       			(SELECT COUNT(*)
        				FROM newf)                                                as total_newf,
       			password_updated_date,
       			(SELECT code FROM languages WHERE language = id_languages) as language
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

	var newf Newf
	err = stmt.QueryRow(email).Scan(&newf.ID, &newf.Email, &newf.FirstName, &newf.LastName, &newf.ProfilePicture, &newf.PhoneNumber, &newf.GraduationYear, &newf.Campus, &newf.TotalUsers, &newf.PasswordUpdatedDate, &newf.Language)
	if err != nil {
		log.Println("║ 💥 Failed to get newf: ", err)
		log.Println("║ 📧 Email: ", email)
		log.Println("╚=========================================╝")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Something went wrong"})
	}

	log.Println("║ ✅ Newf fetched successfully")
	log.Println("║ 📧 Email: ", email)
	log.Println("╚=========================================╝")

	response := make(map[string]interface{})

	if newf.ID != 0 {
		response["id_newf"] = newf.ID
	}
	if newf.Email != "" {
		response["email"] = newf.Email
	}
	if newf.FirstName != "" {
		response["first_name"] = newf.FirstName
	}
	if newf.LastName != "" {
		response["last_name"] = newf.LastName
	}
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
	if newf.TotalUsers != 0 {
		response["total_newf"] = newf.TotalUsers
	}
	if newf.PasswordUpdatedDate != "" {
		response["password_updated_date"] = newf.PasswordUpdatedDate
	}
	if newf.Language != "" {
		response["language"] = newf.Language
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

func updateNewf(c *fiber.Ctx) error {
	var newf Newf

	log.Println("╔======== 📧 Update Newf 📧 ========╗")

	if err := c.BodyParser(&newf); err != nil {
		log.Println("║ 💥 Failed to parse request body: ", err)
		log.Println("╚=========================================╝")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Failed to parse your data"})
	}

	email := c.Locals("email").(string)

	updateFields := make(map[string]interface{})

	if newf.FirstName != "" {
		updateFields["first_name"] = newf.FirstName
	}
	if newf.LastName != "" {
		updateFields["last_name"] = newf.LastName
	}
	if newf.PhoneNumber != "" {
		updateFields["phone_number"] = newf.PhoneNumber
	}
	if newf.GraduationYear != 0 {
		updateFields["graduation_year"] = newf.GraduationYear
	}
	if newf.Campus != "" {
		updateFields["campus"] = newf.Campus
	}
	if newf.NotificationToken != "" {
		updateFields["notification_token"] = newf.NotificationToken
	}
	if newf.ProfilePicture != "" {
		updateFields["profile_picture"] = newf.ProfilePicture
	}
	if newf.Language != "" {
		request := `
			UPDATE newf 
			SET language = (SELECT id_languages FROM languages WHERE name = $1 OR code = $1 LIMIT 1)
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

		_, err = stmt.Exec(newf.Language, email)
		if err != nil {
			log.Println("║ 💥 Failed to update language: ", err)
			log.Println("║ 📧 Email: ", email)
			log.Println("╚=========================================╝")
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Something went wrong"})
		}

	}

	if len(updateFields) == 0 && newf.Language == "" {
		log.Println("║ ⚠️ No fields provided for update")
		log.Println("╚=========================================╝")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "No data provided to update"})
	}
	if len(updateFields) > 0 {
		var queryParts []string
		var values []interface{}
		i := 1

		for column, value := range updateFields {
			queryParts = append(queryParts, fmt.Sprintf("%s = $%d", column, i))
			values = append(values, value)
			i++
		}

		values = append(values, email)

		query := fmt.Sprintf("UPDATE newf SET %s WHERE email = $%d;", strings.Join(queryParts, ", "), i)

		fmt.Println(query, values)

		stmt, err := db.Prepare(query)
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

		_, err = stmt.Exec(values...)
		if err != nil {
			log.Println("║ 💥 Failed to update newf: ", err)
			log.Println("║ 📧 Email: ", email)
			log.Println("╚=========================================╝")
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Something went wrong"})
		}
	}

	log.Println("║ ✅ Newf updated successfully")
	log.Println("║ 🗽 Update Fields: ", updateFields)
	log.Println("║ 📧 Email: ", email)
	log.Println("╚=========================================╝")

	return c.SendStatus(fiber.StatusOK)
}

func addNotification(c *fiber.Ctx) error {
	type NotificationService struct {
		Service string `json:"service"`
	}
	var notificationService NotificationService

	log.Println("╔======== 📞 Add Notification 📞 ========╗")

	if err := c.BodyParser(&notificationService); err != nil {
		log.Println("║ 💥 Failed to parse request body: ", err)
		log.Println("╚=========================================╝")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Failed to parse your data"})
	}

	email := c.Locals("email").(string)

	request := `
		INSERT INTO notifications (email, id_services)
		VALUES ($1, (SELECT id_services FROM services WHERE name = $2));
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

	_, err = stmt.Exec(email, notificationService.Service)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			log.Println("║ 💥 Notification already exists, unsubscribing instead")

			request := `
				DELETE FROM notifications		
				WHERE email = $1	
				AND id_services = (SELECT id_services FROM services WHERE name = $2);
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

			_, err = stmt.Exec(email, notificationService.Service)
			if err != nil {
				log.Println("║ 💥 Failed to delete notification: ", err)
				log.Println("║ 📧 Email: ", email)
				log.Println("╚=========================================╝")
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Something went wrong"})
			}

			log.Println("║ ✅ Notification deleted successfully")
			log.Println("║ 📧 Email: ", email)
			log.Println("║ 📞 Service: ", notificationService.Service)
			log.Println("╚=========================================╝")

			return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": false})
		}

		log.Println("║ 💥 Failed to add notification: ", err)
		log.Println("║ 📧 Email: ", email)
		log.Println("║ 📞 Service: ", notificationService.Service)
		log.Println("╚=========================================╝")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Something went wrong"})
	}

	log.Println("║ ✅ Notification added successfully")
	log.Println("║ 📧 Email: ", email)
	log.Println("║ 📞 Service: ", notificationService.Service)
	log.Println("╚=========================================╝")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": true})
}

func getNotification(c *fiber.Ctx) error {
	email := c.Locals("email").(string)
	log.Println("╔======== 📞 Get Notification 📞 ========╗")

	// Check if body is empty or contains data
	body := c.Body()
	if len(body) == 0 {
		// No body provided, return all subscribed services
		request := `
			SELECT s.name 
			FROM public.notifications n 
			JOIN public.services s ON s.id_services = n.id_services
			WHERE n.email = $1;
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

		rows, err := stmt.Query(email)
		if err != nil {
			log.Println("║ 💥 Failed to get notifications: ", err)
			log.Println("║ 📧 Email: ", email)
			log.Println("╚=========================================╝")
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Something went wrong"})
		}
		defer rows.Close()

		var services []string
		for rows.Next() {
			var service string
			if err := rows.Scan(&service); err != nil {
				log.Println("║ 💥 Failed to scan service: ", err)
				log.Println("╚=========================================╝")
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Something went wrong"})
			}
			services = append(services, service)
		}

		log.Println("║ ✅ All subscribed services fetched successfully")
		log.Println("║ 📧 Email: ", email)
		log.Println("╚=========================================╝")

		return c.Status(fiber.StatusOK).JSON(fiber.Map{"services": services})
	} else {
		// Body contains data, check if it's a single service or array of services
		var payload struct {
			Services []string `json:"services"`
			Service  string   `json:"service"`
		}

		if err := c.BodyParser(&payload); err != nil {
			log.Println("║ 💥 Failed to parse request body: ", err)
			log.Println("╚=========================================╝")
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Failed to parse your data"})
		}

		// Handle single service case (backward compatibility)
		if payload.Service != "" && len(payload.Services) == 0 {
			request := `
				SELECT COUNT(*) FROM public.notifications n join public.services s on s.id_services = n.id_services
				WHERE n.email = $1 AND s.name = $2;
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

			var count int
			err = stmt.QueryRow(email, payload.Service).Scan(&count)
			if err != nil {
				log.Println("║ 💥 Failed to get notification: ", err)
				log.Println("║ 📧 Email: ", email)
				log.Println("║ 📞 Service: ", payload.Service)
				log.Println("╚=========================================╝")
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Something went wrong"})
			}

			log.Println("║ ✅ Notification status fetched successfully")
			log.Println("║ 📧 Email: ", email)
			log.Println("║ 📞 Service: ", payload.Service)
			log.Println("╚=========================================╝")

			return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": count > 0})
		}

		// Handle array of services
		if len(payload.Services) > 0 {
			// Generate placeholders for the SQL query
			placeholders := make([]string, len(payload.Services))
			args := make([]interface{}, len(payload.Services)+1)
			args[0] = email

			for i, service := range payload.Services {
				placeholders[i] = fmt.Sprintf("$%d", i+2)
				args[i+1] = service
			}

			query := fmt.Sprintf(`
				SELECT s.name 
				FROM public.notifications n 
				JOIN public.services s ON s.id_services = n.id_services
				WHERE n.email = $1 AND s.name IN (%s);
			`, strings.Join(placeholders, ", "))

			stmt, err := db.Prepare(query)
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

			rows, err := stmt.Query(args...)
			if err != nil {
				log.Println("║ 💥 Failed to query subscribed services: ", err)
				log.Println("╚=========================================╝")
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Something went wrong"})
			}
			defer func(rows *sql.Rows) {
				err := rows.Close()
				if err != nil {
					log.Println("║ 💥 Failed to close rows: ", err)
					log.Println("╚=========================================╝")
					return
				}
			}(rows)

			subscribedServices := []string{}
			for rows.Next() {
				var service string
				if err := rows.Scan(&service); err != nil {
					log.Println("║ 💥 Failed to scan service: ", err)
					log.Println("╚=========================================╝")
					return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Something went wrong"})
				}
				subscribedServices = append(subscribedServices, service)
			}

			log.Println("║ ✅ Subscribed services fetched successfully")
			log.Println("║ 📧 Email: ", email)
			log.Println("╚=========================================╝")

			return c.Status(fiber.StatusOK).JSON(fiber.Map{"services": subscribedServices})
		}

		// No recognized format in the request
		log.Println("║ 💥 Invalid request format")
		log.Println("╚=========================================╝")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request format"})
	}
}
