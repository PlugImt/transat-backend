package routes

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"github.com/plugimt/transat-backend/handlers/traq"
	"github.com/plugimt/transat-backend/middlewares"
	"github.com/plugimt/transat-backend/utils"
)

func SetupTraqRoutes(router fiber.Router, db *sql.DB) {
	traqHandler := traq.NewTraqHandler(db)
	traqGroup := router.Group("/traq")

	adminOnly := []fiber.Handler{
		middlewares.JWTMiddleware,
		utils.EnhanceSentryEventWithEmail,
		middlewares.AdminAuthMiddleware(db),
	}

	// Types must be registered before /:id so "types" is not captured as an article ID.
	traqTypesGroup := traqGroup.Group("/types")
	traqTypesGroup.Get("/", traqHandler.GetAllTraqTypes)
	traqTypesGroup.Get("/:id", traqHandler.GetTraqType)
	traqTypesGroup.Post("/", append(adminOnly, traqHandler.CreateTraqType)...)
	traqTypesGroup.Patch("/:id", append(adminOnly, traqHandler.UpdateTraqType)...)
	traqTypesGroup.Delete("/:id", append(adminOnly, traqHandler.DeleteTraqType)...)

	traqGroup.Get("/", traqHandler.GetAllTraqArticles)
	traqGroup.Get("/available", traqHandler.GetAvailableTraqArticles)
	traqGroup.Get("/:id", traqHandler.GetTraqArticle)
	traqGroup.Post("/", append(adminOnly, traqHandler.CreateTraqArticle)...)
	traqGroup.Patch("/:id", append(adminOnly, traqHandler.UpdateTraqArticle)...)
	traqGroup.Delete("/:id", append(adminOnly, traqHandler.DeleteTraqArticle)...)
}
