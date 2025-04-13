package routes

import (
	// "database/sql" // No longer needed here

	restaurantHandler "Transat_2.0_Backend/handlers/restaurant" // Import the handler package
	// Remove unused service/utils imports if handler is passed directly
	// "Transat_2.0_Backend/services"
	// "Transat_2.0_Backend/utils"

	"github.com/gofiber/fiber/v2"
)

// SetupRestaurantRoutes configures the routes related to the restaurant menu.
// It now accepts an initialized RestaurantHandler instance.
func SetupRestaurantRoutes(router fiber.Router, h *restaurantHandler.RestaurantHandler) {
	// Use the passed-in handler instance
	// h := restaurantHandler.NewRestaurantHandler(db, transService, notifService) // Removed: Handler is now injected

	// Define the route for getting the menu
	// GET is used because the client might send a language preference in the query params
	router.Get("/restaurant", h.GetRestaurantMenu)
}
