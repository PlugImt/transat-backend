package middlewares

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"github.com/plugimt/transat-backend/utils"
)

// NewfAuthMiddleware rejects requests from users who don't have the NEWF or ADMIN role (aka les staff)
// Must be used after JWTMiddleware.
func NewfAuthMiddleware(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		email, ok := c.Locals("email").(string)
		if !ok || email == "" {
			utils.LogHeader("🎓 Newf Auth Check")
			utils.LogMessage(utils.LevelError, "Missing or invalid email in context")
			utils.LogFooter()
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Authentication required"})
		}

		utils.LogHeader("🎓 Newf Auth Check")
		utils.LogLineKeyValue(utils.LevelInfo, "User", email)

		query := `
			SELECT EXISTS(
				SELECT 1 FROM newf_roles nr
				JOIN roles r ON nr.id_roles = r.id_roles
				WHERE nr.email = $1 AND r.name IN ('NEWF', 'ADMIN')
			)
		`

		var hasAccess bool
		err := db.QueryRow(query, email).Scan(&hasAccess)
		if err != nil {
			utils.LogMessage(utils.LevelError, "Failed to check newf status")
			utils.LogLineKeyValue(utils.LevelError, "Error", err)
			utils.LogFooter()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Authentication error"})
		}

		if !hasAccess {
			utils.LogMessage(utils.LevelWarn, "Staff user attempted student-only access")
			utils.LogFooter()
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "This feature is only available to students"})
		}

		utils.LogMessage(utils.LevelInfo, "Newf access granted")
		utils.LogFooter()
		return c.Next()
	}
}
