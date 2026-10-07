package routes

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"github.com/plugimt/transat-backend/handlers/profile"
	"github.com/plugimt/transat-backend/middlewares"
	"github.com/plugimt/transat-backend/services"
	"github.com/plugimt/transat-backend/utils"
)

// SetupProfileRoutes exposes public user profiles and lets users customize their own.
func SetupProfileRoutes(router fiber.Router, db *sql.DB, r2Service *services.R2Service) {
	handler := profile.NewProfileHandler(db, r2Service)

	group := router.Group("/users", middlewares.JWTMiddleware, middlewares.NewfAuthMiddleware(db), utils.EnhanceSentryEventWithEmail)

	// "/me" must be declared before "/:id".
	group.Get("/me", handler.GetMyProfile)
	group.Patch("/me/profile", handler.UpdateMyProfile)
	group.Get("/:id", handler.GetProfile)
}
