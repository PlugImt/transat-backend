package routes

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"github.com/plugimt/transat-backend/handlers/schedule"
	"github.com/plugimt/transat-backend/handlers/schedule/service"
	"github.com/plugimt/transat-backend/middlewares"
	"github.com/plugimt/transat-backend/utils"
)

func SetupScheduleRoutes(router fiber.Router, db *sql.DB, icsService *service.IcsService) {
	handler := schedule.NewHandler(db, icsService)

	protected := router.Group("/schedule", middlewares.JWTMiddleware, utils.EnhanceSentryEventWithEmail)
	protected.Get("/me", handler.GetMySchedule)
	protected.Patch("/me", handler.UpdateMySchedule)
	protected.Delete("/me", handler.DeleteMySchedule)

	router.Get("/schedule/inte", handler.GetInteSchedule)
}
