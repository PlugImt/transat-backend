package main

import (
	"database/sql"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"math/rand"
	"regexp"
	"time"
)

func checkEmail(email string) (bool, error) {
	regex := `^[a-z0-9.\-]+@imt-atlantique.net$`
	return regexp.MatchString(regex, email)
}

func generate2FACode(digits int) string {
	return fmt.Sprintf("%0"+fmt.Sprint(digits)+"d", rand.Intn(1000000))
}

func generateJWT(newf Newf) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)

	var role string

	request := `
		SELECT name
		FROM roles join public.newf_roles nr on roles.id_roles = nr.id_roles
		WHERE nr.email = $1;
	`

	err := db.QueryRow(request, newf.Email).Scan(&role)
	if err != nil {
		return "", err
	}

	claims := jwt.MapClaims{
		"email": newf.Email,
		"role":  role,
		"exp":   expirationTime.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func validateJWT(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Ensure the signing method is correct
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return token, nil
}

func isValidated(newf Newf) bool {
	request := `
		SELECT name
		FROM newf_roles join public.roles r on r.id_roles = newf_roles.id_roles
		WHERE email = $1;
	`

	var role string

	err := db.QueryRow(request, newf.Email).Scan(&role)
	if err != nil {
		return false
	}

	return role != "VERIFYING"
}

func getVerificationCode(newf Newf) (verificationCodeData, error) {
	request := `
		SELECT verification_code, verification_code_expiration
		FROM newf
		WHERE email = $1;
	`

	var data verificationCodeData

	stmt, err := db.Prepare(request)
	if err != nil {
		return data, err
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			fmt.Println("Error closing statement:", err)
		}
	}(stmt)

	err = stmt.QueryRow(newf.Email).Scan(&data.VerificationCode, &data.VerificationCodeExpiration)
	if err != nil {
		return data, err
	}

	// Parse the date
	t, err := time.Parse(time.RFC3339Nano, data.VerificationCodeExpiration)
	if err != nil {
		fmt.Println("Error parsing date:", err)
		return data, err
	}

	data.VerificationCodeExpiration = t.Format("15h04")

	return data, nil
}
