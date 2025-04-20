package routes

import (
	"Transat_2.0_Backend/handlers/washingmachine"
	"github.com/gofiber/fiber/v2"
)

// SetupWashingMachineRoutes configures the routes for washing machine data
func SetupWashingMachineRoutes(router fiber.Router) {
	// Create a new handler
	handler := washingmachine.NewWashingMachineHandler()
	
	// Create a group for washing machine routes
	washingMachines := router.Group("/washingmachines")
	
	// Define routes
	washingMachines.Get("/", handler.GetWashingMachines())
} 