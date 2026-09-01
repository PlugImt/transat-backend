package middlewares

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/plugimt/transat-backend/utils"
)

func JWTMiddleware(db *sql.DB) fiber.Handler {
	sem := make(chan struct{}, 50)
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")

		utils.LogHeader("📧 JWT Middleware")

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

		email, _ := claims["email"].(string)
		c.Locals("email", email)

		// Note: We only store email in context - roles are checked from database for security

		utils.LogMessage(utils.LevelInfo, "Token is valid")
		utils.LogLineKeyValue(utils.LevelInfo, "Email", email)
		utils.LogFooter()

		// Update last_activity at most once per day; bounded to 50 concurrent writes.
		select {
		case sem <- struct{}{}:
			go func() {
				defer func() { <-sem }()
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()
				if _, err := db.ExecContext(ctx, `UPDATE newf SET last_activity = NOW() WHERE email = $1 AND (last_activity IS NULL OR last_activity < NOW() - INTERVAL '1 day')`, email); err != nil {
					utils.LogMessage(utils.LevelError, "Failed to update last_activity")
					utils.LogLineKeyValue(utils.LevelError, "Error", err)
				}
			}()
		default:
		}

		return c.Next()
	}
}
