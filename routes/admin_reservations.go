package routes

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"github.com/plugimt/transat-backend/handlers/admin"
	"github.com/plugimt/transat-backend/handlers/reservation"
	"github.com/plugimt/transat-backend/services"
	"github.com/plugimt/transat-backend/services/reservations"
)

// SetupAdminReservationRoutes adds the reservation administration to the admin group, so every route is behind the
// JWT and the ADMIN role checks of that group.
func SetupAdminReservationRoutes(adminGroup fiber.Router, db *sql.DB, discord *services.DiscordService, service *reservations.Service) {
	handler := admin.NewReservationsAdminHandler(db, service, discord)
	creation := reservation.NewReservationHandler(db, discord, service)

	adminGroup.Get("/reservations", handler.GetReservations)
	adminGroup.Get("/reservations/stats", handler.GetStats)
	adminGroup.Get("/reservations/items", handler.GetItems)
	adminGroup.Post("/reservations/cancel", handler.CancelReservation)
	adminGroup.Post("/reservations/items/:id/return", handler.MarkReturned)

	adminGroup.Get("/reservations/waitlist", handler.GetWaitlist)
	adminGroup.Post("/reservations/waitlist/advance", handler.AdvanceWaitlist)
	adminGroup.Delete("/reservations/waitlist/:id", handler.RemoveWaitlistEntry)

	// Creation reuses the validation of the app's handlers, behind the administrator check.
	adminGroup.Post("/reservations/items", creation.CreateReservationItem)
	adminGroup.Post("/reservations/categories", creation.CreateReservationCategory)
}
