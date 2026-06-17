package routes

import (
	"github.com/plugimt/transat-backend/handlers/planning_sport"
	"github.com/plugimt/transat-backend/middlewares"
	"github.com/plugimt/transat-backend/utils"

	"github.com/gofiber/fiber/v2"
)

func SetupPlanningSportRoutes(router fiber.Router, planningSportHandler *planning_sport.PlanningSportHandler) {

	planningSportGroup := router.Group("/planning-sport", middlewares.JWTMiddleware, utils.EnhanceSentryEventWithEmail)

	planningSportGroup.Get("", planningSportHandler.GetPlanning)
	planningSportGroup.Get("/", planningSportHandler.GetPlanning)
}
