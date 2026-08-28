package routes

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/plugimt/transat-backend/handlers/washingmachine"
)

// SetupWashingMachineRoutes configures the routes for washing machine data
func SetupWashingMachineRoutes(router fiber.Router) {
	handler := washingmachine.NewWashingMachineHandler()

	washingMachines := router.Group("/washingmachines")
	washingMachines.Get("", handler.GetWashingMachines())

	if os.Getenv("ENV") != "production" {
		washingMachines.Get("/test", handler.GetWashingMachinesTest())
	}
}
