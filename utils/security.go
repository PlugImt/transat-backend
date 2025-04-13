package utils

import (
	"fmt"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret []byte

func ValidateJWT(tokenString string) (*jwt.Token, error) {
	if jwtSecret == nil {
		jwtSecret = []byte(os.Getenv("JWT_SECRET"))
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
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
