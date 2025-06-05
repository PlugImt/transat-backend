package routes

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"github.com/plugimt/transat-backend/handlers/planning"
	"github.com/plugimt/transat-backend/middlewares"
)

// SetupPlanningRoutes configures planning-related routes.
func SetupPlanningRoutes(router fiber.Router, db *sql.DB) {
	planningHandler := planning.NewPlanningHandler(db)

	planningGroup := router.Group("/planning", middlewares.JWTMiddleware)

	// GET /planning/users - List users with only NewfID, FirstName, LastName, PassID
	planningGroup.Get("/users", planningHandler.GetUsersWithPassID)

	// PATCH /planning/users/:id_newf/passid - Update user's PassID if null
	planningGroup.Patch("/users/:id_newf/passid", planningHandler.UpdateUserPassIDIfNull)

	// POST /planning/courses - Add a new course
	planningGroup.Post("/courses", planningHandler.CreateCourse)

	// GET /planning/users/:email/courses?start=YYYY-MM-DD&end=YYYY-MM-DD - Get user's courses between two dates
	planningGroup.Get("/users/:email/courses", planningHandler.GetUserCoursesBetweenDates)

	// GET /planning/users/:email/courses/today - Get user's courses for the current day
	planningGroup.Get("/users/:email/courses/today", planningHandler.GetUserCoursesToday)
}
