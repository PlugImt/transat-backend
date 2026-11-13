package routes

import (
	"database/sql"
	"slices"

	"github.com/gofiber/fiber/v2"
	"github.com/plugimt/transat-backend/middlewares"
	"github.com/plugimt/transat-backend/services/photoreport"
)

// photoSections are what the photo managers open: reports, the photobooth's review and the albums.
var photoSections = []string{"dashboard", "photo-reports", "photobooth", "photo-albums"}

// Panel sections a non-administrator may open, by role.
var sectionsByRole = map[string][]string{
	photoreport.PhotoManagerRole: photoSections,
	photoreport.PhotoClubRole:    photoSections,
}

// SetupAccessRoutes exposes what the signed-in account may open in the admin panel. It must be declared before
// the admin-only "/admin" group.
func SetupAccessRoutes(router fiber.Router, db *sql.DB) {
	router.Get("/admin/access", middlewares.JWTMiddleware, func(c *fiber.Ctx) error {
		email, _ := c.Locals("email").(string)
		roles, err := middlewares.UserRoles(c.UserContext(), db, email)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not load access"})
		}
		admin := slices.Contains(roles, middlewares.AdminRole)
		sections := []string{}
		for _, r := range roles {
			for _, s := range sectionsByRole[r] {
				if !slices.Contains(sections, s) {
					sections = append(sections, s)
				}
			}
		}
		if admin {
			sections = []string{"*"}
		}
		return c.JSON(fiber.Map{"email": email, "roles": roles, "admin": admin, "sections": sections})
	})
}
