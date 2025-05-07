package utils

import (
	"crypto/rand"
	"database/sql"
	"fmt"
	"log" // Use log instead of custom logger for basic helpers?
	"math/big"
	"regexp"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/plugimt/transat-backend/models" // Need this for GetVerificationCodeData
)

// CheckEmail validates if the email format is correct for imt-atlantique.net.
func CheckEmail(email string) (bool, error) {
	// Consider making the domain configurable
	regex := `^[a-z0-9.\-]+@imt-atlantique\.net$`
	matched, err := regexp.MatchString(regex, email)
	if err != nil {
		log.Printf("Error matching email regex: %v", err)
		return false, fmt.Errorf("regex error: %w", err)
	}
	return matched, nil
}

// Generate2FACode generates a random numeric string of the specified length.
func Generate2FACode(digits int) string {
	if digits <= 0 {
		digits = 6
	}

	// Calcul de la borne supérieure (ex : 10^6 pour 6 chiffres)
	max := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(digits)), nil)

	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "000000"
	}

	// Formatage avec zéros en tête
	format := "%0" + strconv.Itoa(digits) + "s"
	return fmt.Sprintf(format, n.Text(10))
}

// GenerateJWT creates a new JWT token for a user with their email and role.
func GenerateJWT(email string, role string, jwtSecret []byte) (string, error) {
	if len(jwtSecret) == 0 {
		log.Println("Error generating JWT: JWT secret is empty")
		return "", fmt.Errorf("JWT secret is not configured")
	}

	expirationTime := time.Now().Add(365 * 24 * time.Hour) // Token valid for 1 year

	claims := jwt.MapClaims{
		"email": email,
		"role":  role,
		"exp":   expirationTime.Unix(),
		"iat":   time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		log.Printf("Error signing JWT for %s: %v", email, err)
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

// IsValidated checks if a user role is not 'VERIFYING'.
// This function might be less useful now as role is checked directly in handlers.
// func IsValidated(db *sql.DB, email string) bool {
// 	query := `
// 		SELECT EXISTS (
// 			SELECT 1 FROM newf_roles nr JOIN roles r ON nr.id_roles = r.id_roles
// 			WHERE nr.email = $1 AND r.name != 'VERIFYING'
// 		);
// 	`
// 	var validated bool
// 	err := db.QueryRow(query, email).Scan(&validated)
// 	if err != nil {
// 		log.Printf("Error checking if user %s is validated: %v", email, err)
// 		return false
// 	}
// 	return validated
// }

// GetVerificationCodeData fetches the verification code and its formatted expiration time.
func GetVerificationCodeData(db *sql.DB, email string) (models.VerificationCodeData, error) {
	query := `
		SELECT verification_code, verification_code_expiration
		FROM newf
		WHERE email = $1;
	`
	var data models.VerificationCodeData
	var expiration time.Time // Scan into time.Time first

	err := db.QueryRow(query, email).Scan(&data.VerificationCode, &expiration)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("No verification code data found for user %s", email)
			return data, fmt.Errorf("user not found or no code set")
		}
		log.Printf("Error fetching verification code data for %s: %v", email, err)
		return data, fmt.Errorf("database error fetching code data: %w", err)
	}

	// Format the expiration time (e.g., 15h04)
	// Consider making the format configurable or using a standard like RFC3339
	if !expiration.IsZero() {
		data.VerificationCodeExpiration = expiration.Format("15h04") // Example format
	} else {
		data.VerificationCodeExpiration = "N/A" // Or empty string
	}

	return data, nil
}

// GetLanguageCode fetches the language code (e.g., 'fr', 'en') for a user.
func GetLanguageCode(db *sql.DB, email string) (string, error) {
	query := `
		SELECT COALESCE(l.code, 'fr') -- Default to 'fr' if language or mapping is missing
		FROM newf n
		LEFT JOIN languages l ON n.language = l.id_languages
		WHERE n.email = $1;
	`
	var languageCode string
	err := db.QueryRow(query, email).Scan(&languageCode)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("User %s not found when fetching language code", email)
			// Should we return default 'fr' or an error?
			return "fr", fmt.Errorf("user not found")
		}
		log.Printf("Error fetching language code for %s: %v", email, err)
		return "fr", fmt.Errorf("database error fetching language code: %w", err) // Return default and error
	}
	return languageCode, nil
}
