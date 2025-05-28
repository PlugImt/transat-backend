package utils

import (
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret []byte
var jwtExpirationHours time.Duration = 24 * time.Hour // Default expiration

// InitJWT initializes the JWT secret and expiration time.
func InitJWT(secret []byte, expirationHoursEnv string) error {
	if len(secret) == 0 {
		return fmt.Errorf("JWT secret cannot be empty")
	}
	jwtSecret = secret

	if expirationHoursEnv != "" {
		expHours, err := strconv.Atoi(expirationHoursEnv)
		if err != nil {
			// Log and continue with default if parsing fails
			LogMessage(LevelWarn, fmt.Sprintf("Invalid JWT_EXPIRATION_HOURS: '%s'. Using default %v. Error: %v", expirationHoursEnv, jwtExpirationHours, err))
		} else {
			jwtExpirationHours = time.Duration(expHours) * time.Hour
		}
	} // else, default is used
	LogMessage(LevelInfo, fmt.Sprintf("JWT Initialized. Expiration: %v", jwtExpirationHours))
	return nil
}

// JWTClaims represents the claims in a JWT token with enhanced security fields
type JWTClaims struct {
	Email       string `json:"email"`
	Role        string `json:"role,omitempty"`
	Fingerprint string `json:"fingerprint,omitempty"` // For device fingerprinting
	jwt.RegisteredClaims
}

// GenerateJWT creates a secure JWT token with enhanced claims
func GenerateJWT(email, role string, fingerprint string) (string, error) {
	if jwtSecret == nil {
		// This should not happen if InitJWT was called
		return "", fmt.Errorf("JWT secret is not initialized. Call utils.InitJWT during startup.")
	}

	// Create expiration time using the configured duration
	expirationTime := time.Now().Add(jwtExpirationHours)

	// Create enhanced claims
	claims := &JWTClaims{
		Email:       email,
		Role:        role,
		Fingerprint: fingerprint,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "transat-backend",
			Subject:   email,
			ID:        GenerateRandomString(16), // Unique JWT ID
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate signed token
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateJWT validates a JWT token with enhanced security checks
func ValidateJWT(tokenString string) (*jwt.Token, error) {
	if jwtSecret == nil {
		// This should not happen if InitJWT was called
		return nil, fmt.Errorf("JWT secret is not initialized. Call utils.InitJWT during startup.")
	}

	// Parse token with custom claims
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		// Validate signing method
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

	// Get claims and perform additional validation checks
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid claims")
	}

	// Check issuer
	if iss, ok := claims["iss"].(string); !ok || iss != "transat-backend" {
		return nil, fmt.Errorf("invalid issuer")
	}

	return token, nil
}

// GenerateRandomString generates a cryptographically secure random string
func GenerateRandomString(length int) string {
	// Implementation using crypto/rand
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)

	// Use crypto/rand to generate secure random bytes
	_, err := SecureRandom(b)
	if err != nil {
		// Fallback to less secure method in case of error
		for i := range b {
			b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
			time.Sleep(1 * time.Nanosecond)
		}
		return string(b)
	}

	// Map random bytes to charset
	for i := 0; i < length; i++ {
		b[i] = charset[int(b[i])%len(charset)]
	}

	return string(b)
}

// ValidateTokenFingerprint checks if the token's fingerprint matches the provided one
func ValidateTokenFingerprint(token *jwt.Token, fingerprint string) bool {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return false
	}

	tokenFingerprint, ok := claims["fingerprint"].(string)
	if !ok {
		return false
	}

	return tokenFingerprint == fingerprint
}
