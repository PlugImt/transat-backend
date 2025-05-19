package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/plugimt/transat-backend/handlers/washingmachine"
)

// SetupWashingMachineRoutes configures the routes for washing machine data
func SetupWashingMachineRoutes(router fiber.Router) {
	// Create a new handler
	handler := washingmachine.NewWashingMachineHandler()

	// Create a group for washing machine routes
	washingMachines := router.Group("/washingmachines")

	// Define routes
	router.Get("/washingmachines", handler.GetWashingMachines())
	washingMachines.Get("", handler.GetWashingMachines())
	washingMachines.Get("/", handler.GetWashingMachines())
}
