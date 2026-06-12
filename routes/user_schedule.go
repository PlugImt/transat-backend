package routes

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	user_schedule "github.com/plugimt/transat-backend/handlers/user_schedule"
	"github.com/plugimt/transat-backend/middlewares"
	"github.com/plugimt/transat-backend/utils"
)

func SetupUserScheduleRoutes(router fiber.Router, db *sql.DB) {
	handler := user_schedule.NewUserScheduleHandler(db)

	scheduleGroup := router.Group("/schedule", middlewares.JWTMiddleware, utils.EnhanceSentryEventWithEmail)

	scheduleGroup.Get("/me", handler.GetMySchedule)
	scheduleGroup.Patch("/me", handler.UpdateMySchedule)
	scheduleGroup.Delete("/me", handler.DeleteMySchedule)
}
