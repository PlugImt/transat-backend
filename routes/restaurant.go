package routes

import (
	"github.com/gofiber/fiber/v2"
	restaurantHandler "github.com/plugimt/transat-backend/handlers/restaurant"
	"github.com/plugimt/transat-backend/middlewares"
	"github.com/plugimt/transat-backend/utils"
)

func SetupRestaurantRoutes(router fiber.Router, h *restaurantHandler.RestaurantHandler) {
	// Public routes (no authentication required)
	router.Get("/restaurant", h.GetRestaurantMenu)
	router.Get("/restaurant/:id", h.GetDishDetails)

	// Protected routes (require JWT authentication)
	router.Post("/restaurant/:id", middlewares.JWTMiddleware, utils.EnhanceSentryEventWithEmail, h.PostDishReview)
}
