package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret []byte

// JWTClaims represents the claims in a JWT token with enhanced security fields
type JWTClaims struct {
	Email       string   `json:"email"`
	Role        string   `json:"role,omitempty"`        // Deprecated: kept for backward compatibility
	Roles       []string `json:"roles,omitempty"`       // New: array of all user roles
	Fingerprint string   `json:"fingerprint,omitempty"` // For device fingerprinting
	jwt.RegisteredClaims
}

// GenerateJWT creates a secure JWT token with enhanced claims
func GenerateJWT(email string, roles []string, fingerprint string) (string, error) {
	if jwtSecret == nil {
		jwtSecret = []byte(os.Getenv("JWT_SECRET"))
		if len(jwtSecret) == 0 {
			return "", fmt.Errorf("JWT secret is not set")
		}
	}

	// Set expiration time - 24 hours by default
	expirationHours := 2400000000000
	expirationEnv := os.Getenv("JWT_EXPIRATION_HOURS")
	if expirationEnv != "" {
		if _, err := fmt.Sscanf(expirationEnv, "%d", &expirationHours); err != nil {
			return "", fmt.Errorf("invalid JWT_EXPIRATION_HOURS: %v", err)
		}
	}

	// Create expiration time using Paris timezone
	now := Now()
	expirationTime := AddInParis(now, time.Duration(expirationHours)*time.Hour)

	// For backward compatibility, set Role to the first role if any exist
	var primaryRole string
	if len(roles) > 0 {
		primaryRole = roles[0]
	}

	// Create enhanced claims
	claims := &JWTClaims{
		Email:       email,
		Role:        primaryRole, // Deprecated field for backward compatibility
		Roles:       roles,       // New field with all roles
		Fingerprint: fingerprint,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
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
		jwtSecret = []byte(os.Getenv("JWT_SECRET"))
		if len(jwtSecret) == 0 {
			return nil, fmt.Errorf("JWT secret is not set")
		}
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
			b[i] = charset[UnixNanoParis(Now())%int64(len(charset))]
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

func GetRolesFromToken(token *jwt.Token) ([]string, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid claims")
	}

	// Try to get roles array first (new format)
	if rolesInterface, exists := claims["roles"]; exists {
		if rolesArray, ok := rolesInterface.([]interface{}); ok {
			roles := make([]string, len(rolesArray))
			for i, role := range rolesArray {
				if roleStr, ok := role.(string); ok {
					roles[i] = roleStr
				} else {
					return nil, fmt.Errorf("invalid role type in roles array")
				}
			}
			return roles, nil
		}
	}

	// Fallback to single role field for backward compatibility
	if roleInterface, exists := claims["role"]; exists {
		if role, ok := roleInterface.(string); ok && role != "" {
			return []string{role}, nil
		}
	}

	return []string{}, nil
}

func HasRole(token *jwt.Token, targetRole string) bool {
	roles, err := GetRolesFromToken(token)
	if err != nil {
		return false
	}

	for _, role := range roles {
		if role == targetRole {
			return true
		}
	}
	return false
}

func HasAnyRole(token *jwt.Token, targetRoles []string) bool {
	roles, err := GetRolesFromToken(token)
	if err != nil {
		return false
	}

	for _, userRole := range roles {
		for _, targetRole := range targetRoles {
			if userRole == targetRole {
				return true
			}
		}
	}
	return false
}
