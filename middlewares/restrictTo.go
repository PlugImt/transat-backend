package middlewares

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/plugimt/transat-backend/utils"
)

// RestrictTo creates a middleware that restricts access to users with specific roles
// Usage: app.Use(middlewares.RestrictTo(db, "ADMIN", "MODERATOR"))
func RestrictTo(db *sql.DB, allowedRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		utils.LogHeader("🔐 Role Authorization")

		// Validate that at least one role is specified
		if len(allowedRoles) == 0 {
			utils.LogMessage(utils.LevelError, "No roles specified for restrictTo middleware")
			utils.LogFooter()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Server configuration error"})
		}

		// First check if user is authenticated (JWT validation)
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			utils.LogMessage(utils.LevelError, "Missing token")
			utils.LogFooter()
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing token"})
		}

		tokenString := authHeader
		if strings.HasPrefix(authHeader, "Bearer ") {
			tokenString = authHeader[7:]
		}

		token, err := utils.ValidateJWT(tokenString)
		if err != nil {
			utils.LogMessage(utils.LevelError, "Invalid token")
			utils.LogLineKeyValue(utils.LevelError, "Error", err)
			utils.LogFooter()
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			utils.LogMessage(utils.LevelError, "Invalid claims")
			utils.LogFooter()
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid claims"})
		}

		email, ok := claims["email"].(string)
		if !ok {
			utils.LogMessage(utils.LevelError, "Invalid email in claims")
			utils.LogFooter()
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token claims"})
		}

		// Build query to check if user has any of the allowed roles
		// Create placeholders for the IN clause
		placeholders := make([]string, len(allowedRoles))
		args := make([]interface{}, len(allowedRoles)+1)
		args[0] = strings.ToLower(email)

		for i, role := range allowedRoles {
			placeholders[i] = fmt.Sprintf("$%d", i+2) // Start from $2 since $1 is email
			args[i+1] = strings.ToUpper(role)
		}

		query := `
			SELECT r.name
			FROM newf_roles nr
			JOIN roles r ON nr.id_roles = r.id_roles
			WHERE nr.email = $1 AND r.name IN (` + strings.Join(placeholders, ",") + `)
			LIMIT 1
		`

		var userRole string
		err = db.QueryRow(query, args...).Scan(&userRole)
		if err != nil {
			if err == sql.ErrNoRows {
				utils.LogMessage(utils.LevelWarn, "Access denied: User does not have required role")
				utils.LogLineKeyValue(utils.LevelWarn, "Email", email)
				utils.LogLineKeyValue(utils.LevelWarn, "Required Roles", strings.Join(allowedRoles, ", "))
				utils.LogFooter()
				return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
					"error": "Access denied: Insufficient permissions",
				})
			}
			utils.LogMessage(utils.LevelError, "Failed to check user roles")
			utils.LogLineKeyValue(utils.LevelError, "Email", email)
			utils.LogLineKeyValue(utils.LevelError, "Error", err)
			utils.LogFooter()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Authorization check failed"})
		}

		// Set user info in context for use in handlers
		c.Locals("email", email)
		c.Locals("role", userRole)
		c.Locals("allowedRoles", allowedRoles)

		utils.LogMessage(utils.LevelInfo, "Access granted")
		utils.LogLineKeyValue(utils.LevelInfo, "Email", email)
		utils.LogLineKeyValue(utils.LevelInfo, "User Role", userRole)
		utils.LogLineKeyValue(utils.LevelInfo, "Required Roles", strings.Join(allowedRoles, ", "))
		utils.LogFooter()

		return c.Next()
	}
}

// GetUserRoles returns all roles for a given user (utility function)
func GetUserRoles(db *sql.DB, email string) ([]string, error) {
	query := `
		SELECT r.name
		FROM newf_roles nr
		JOIN roles r ON nr.id_roles = r.id_roles
		WHERE nr.email = $1
	`

	rows, err := db.Query(query, strings.ToLower(email))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []string
	for rows.Next() {
		var role string
		if err := rows.Scan(&role); err != nil {
			continue // Skip invalid rows
		}
		roles = append(roles, role)
	}

	return roles, nil
}
