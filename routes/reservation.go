package routes

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"github.com/plugimt/transat-backend/handlers/reservation" // Import the reservation handlers
	"github.com/plugimt/transat-backend/middlewares"
	"github.com/plugimt/transat-backend/services"
	"github.com/plugimt/transat-backend/services/reservations"
)

func SetupReservationRoutes(router fiber.Router, db *sql.DB, discordService *services.DiscordService, reservationService *reservations.Service) {
	// Initialize Reservation Handler
	reservationHandler := reservation.NewReservationHandler(db, discordService, reservationService)

	reservationGroup := router.Group("/reservation", middlewares.JWTMiddleware, middlewares.NewfAuthMiddleware(db))

	// Root reservation routes
	reservationGroup.Get("", reservationHandler.GetReservationItems)      // Returns root categories and items TODO: Add filders && add reservations overview
	reservationGroup.Get("/", reservationHandler.GetReservationItems)     // Returns root categories and items
	reservationGroup.Get("/me", reservationHandler.GetMyReservations)     // Returns logged in user's reservations
	reservationGroup.Get("/search", reservationHandler.SearchReservation) // Search items and categories
	reservationGroup.Get("/waitlist", reservationHandler.GetMyWaitlist)   // Returns the waiting lists the logged in user is in

	// Category routes
	reservationGroup.Post("/category", reservationHandler.CreateReservationCategory) // Creates a new category
	reservationGroup.Get("/category/:id", reservationHandler.GetReservationItems)    // Returns all the items in a specific category by ID
	reservationGroup.Get("/club/:id", reservationHandler.GetReservationItemsForClub) // Returns all the items in a specific category by ID

	// Item routes - these need to come after category routes to avoid conflicts
	reservationGroup.Get("/item/:id", reservationHandler.GetItemDetails)           // Returns the details of a specific item by ID
	reservationGroup.Patch("/item/:id", reservationHandler.UpdateReservationItem)  // Reserve or remove the reservation of an item by ID
	reservationGroup.Post("/item", reservationHandler.CreateReservationItem)       // Creates a new item with a parent category (if applicable)
	reservationGroup.Post("/item/", reservationHandler.CreateReservationItem)      // Creates a new item with a parent category (if applicable)
	reservationGroup.Delete("/item/:id", reservationHandler.DeleteReservationItem) // Deletes a reservation for an item (an admin may delete anyone's)

	// Waiting list: interest in an item or slot that is already reserved. start_date (query or body) names the slot.
	reservationGroup.Get("/item/:id/waitlist", reservationHandler.GetItemWaitlist)              // Queue length and the caller's place
	reservationGroup.Post("/item/:id/waitlist", reservationHandler.JoinWaitlist)                // Join the queue
	reservationGroup.Delete("/item/:id/waitlist", reservationHandler.LeaveWaitlist)             // Leave the queue
	reservationGroup.Post("/item/:id/waitlist/reserve", reservationHandler.ReserveFromWaitlist) // "Réserver" action of the notification
	reservationGroup.Post("/item/:id/waitlist/decline", reservationHandler.LeaveWaitlist)       // "Annuler" action of the notification
}
