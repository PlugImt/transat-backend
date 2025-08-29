package middlewares

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"github.com/plugimt/transat-backend/utils"
)

func AdminAuthMiddleware(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		email := c.Locals("email").(string)
		
		utils.LogHeader("🔐 Admin Auth Check")
		utils.LogLineKeyValue(utils.LevelInfo, "User", email)

		query := `
			SELECT EXISTS(
				SELECT 1 FROM newf_roles nr
				JOIN roles r ON nr.id_roles = r.id_roles
				WHERE nr.email = $1 AND r.name = 'ADMIN'
			)
		`
		
		var isAdmin bool
		err := db.QueryRow(query, email).Scan(&isAdmin)
		if err != nil {
			utils.LogMessage(utils.LevelError, "Failed to check admin status")
			utils.LogLineKeyValue(utils.LevelError, "Error", err)
			utils.LogFooter()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Authentication error"})
		}

		if !isAdmin {
			utils.LogMessage(utils.LevelWarn, "Non-admin user attempted admin access")
			utils.LogFooter()
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Admin access required"})
		}

		utils.LogMessage(utils.LevelInfo, "Admin access granted")
		utils.LogFooter()
		return c.Next()
	}
}