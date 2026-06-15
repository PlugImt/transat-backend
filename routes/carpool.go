package routes

import (
	"github.com/plugimt/transat-backend/handlers/carpool"
	"github.com/plugimt/transat-backend/middlewares"
	"github.com/plugimt/transat-backend/utils"

	"github.com/gofiber/fiber/v2"
)

func SetupCarpoolRoutes(router fiber.Router, carpoolHandler *carpool.CarpoolHandler) {

	carpoolGroup := router.Group("/carpool", middlewares.JWTMiddleware, utils.EnhanceSentryEventWithEmail)

	carpoolGroup.Get("", carpoolHandler.GetCarpools)
	carpoolGroup.Get("/", carpoolHandler.GetCarpools)

	carpoolGroup.Get("/:id", carpoolHandler.GetCarpoolByID)

	carpoolGroup.Post("", carpoolHandler.CreateCarpool)
	carpoolGroup.Post("/", carpoolHandler.CreateCarpool)

	carpoolGroup.Patch("/:id/status", carpoolHandler.UpdateCarpoolStatus)
	carpoolGroup.Delete("/:id", carpoolHandler.DeleteCarpool)
}
