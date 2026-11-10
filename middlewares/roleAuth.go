package middlewares

import (
	"context"
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"github.com/lib/pq"
	"github.com/plugimt/transat-backend/utils"
)

// AdminRole is the role that opens every page of the admin panel.
const AdminRole = "ADMIN"

// UserRoles returns the names of the roles held by an account.
func UserRoles(ctx context.Context, db *sql.DB, email string) ([]string, error) {
	rows, err := db.QueryContext(ctx, `
		SELECT r.name FROM newf_roles nr JOIN roles r ON r.id_roles = nr.id_roles WHERE nr.email = $1 ORDER BY r.name`, email)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	roles := []string{}
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		roles = append(roles, name)
	}
	return roles, rows.Err()
}

// AdminOrRolesMiddleware lets administrators in, plus the holders of any of the given roles. It runs after
// JWTMiddleware, which puts the verified email in the context.
func AdminOrRolesMiddleware(db *sql.DB, roles ...string) fiber.Handler {
	allowed := append([]string{AdminRole}, roles...)
	return func(c *fiber.Ctx) error {
		email, ok := c.Locals("email").(string)
		if !ok || email == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Authentication required"})
		}
		var granted bool
		err := db.QueryRowContext(c.UserContext(), `
			SELECT EXISTS (SELECT 1 FROM newf_roles nr JOIN roles r ON r.id_roles = nr.id_roles
			               WHERE nr.email = $1 AND r.name = ANY($2))`, email, pq.Array(allowed)).Scan(&granted)
		if err != nil {
			utils.LogMessage(utils.LevelError, "Failed to check panel access")
			utils.LogLineKeyValue(utils.LevelError, "Error", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Authentication error"})
		}
		if !granted {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Access denied"})
		}
		return c.Next()
	}
}
