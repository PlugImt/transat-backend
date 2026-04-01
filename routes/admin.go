package routes

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"github.com/plugimt/transat-backend/handlers/admin"
	"github.com/plugimt/transat-backend/middlewares"
	"github.com/plugimt/transat-backend/utils"
)

func SetupAdminRoutes(router fiber.Router, db *sql.DB) {
	adminHandler := admin.NewAdminHandler(db)

	adminGroup := router.Group("/admin",
		middlewares.JWTMiddleware,
		utils.EnhanceSentryEventWithEmail,
		middlewares.AdminAuthMiddleware(db),
	)

	adminGroup.Get("/users", adminHandler.GetAllUsers)
	adminGroup.Post("/users", adminHandler.CreateUser)
	adminGroup.Patch("/users/:email", adminHandler.UpdateUser)
	adminGroup.Delete("/users/:email", adminHandler.DeleteUser)
	adminGroup.Post("/users/:email/validate", adminHandler.ValidateUser)
	adminGroup.Get("/roles", adminHandler.GetAllRoles)

	adminGroup.Get("/events", adminHandler.GetAllEvents)
	adminGroup.Post("/events", adminHandler.CreateEvent)
	adminGroup.Patch("/events/:id", adminHandler.UpdateEvent)
	adminGroup.Delete("/events/:id", adminHandler.DeleteEvent)

	adminGroup.Get("/clubs", adminHandler.GetAllClubs)
	adminGroup.Post("/clubs", adminHandler.CreateClub)
	adminGroup.Patch("/clubs/:id", adminHandler.UpdateClub)
	adminGroup.Delete("/clubs/:id", adminHandler.DeleteClub)

	// Restaurant menu management
	adminGroup.Get("/menu", adminHandler.GetAllMenuItems)
	adminGroup.Delete("/menu/:id", adminHandler.DeleteMenuItem)
	adminGroup.Get("/menu/:id/reviews", adminHandler.GetMenuItemReviews)
	adminGroup.Delete("/menu/:id/reviews/:email", adminHandler.DeleteMenuItemReview)

	// Reviews management
	adminGroup.Get("/reviews", adminHandler.GetAllReviews)

	// Bassine management
	adminGroup.Get("/bassine/scores", adminHandler.GetBassineScores)
	adminGroup.Post("/bassine/update-score", adminHandler.UpdateBassineScore)
	adminGroup.Get("/bassine/history/:email", adminHandler.GetBassineHistory)

	// Reservation item management
	adminGroup.Get("/clubs/:id/reservation-items", adminHandler.GetReservationItemsForClubAdmin)
	adminGroup.Patch("/reservation-items/:id/messages", adminHandler.UpdateReservationItemMessages)
	adminGroup.Get("/reservations/tree", adminHandler.GetReservationTree)
	adminGroup.Delete("/reservations/categories/:id", adminHandler.DeleteReservationCategory)
	adminGroup.Patch("/reservations/categories/:id", adminHandler.UpdateReservationCategory)
	adminGroup.Delete("/reservations/items/:id", adminHandler.DeleteReservationItemAdmin)
	adminGroup.Patch("/reservations/items/:id", adminHandler.UpdateReservationItemAdmin)

}
