package routes

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	restaurantHandler "github.com/plugimt/transat-backend/handlers/restaurant"
	"github.com/plugimt/transat-backend/middlewares"
	"github.com/plugimt/transat-backend/utils"
)

func SetupRestaurantRoutes(router fiber.Router, db *sql.DB, h *restaurantHandler.RestaurantHandler) {
	// Public routes (no authentication required)

	restaurant := router.Group("/restaurant", middlewares.JWTMiddleware(db), utils.EnhanceSentryEventWithEmail)

	restaurant.Get("/", h.GetRestaurantMenu)
	restaurant.Get("", h.GetRestaurantMenu)

	restaurant.Get("/test", h.GetRestaurantTestMenu)

	restaurant.Get("/:id", h.GetDishDetails)

	restaurant.Post("/:id", h.PostDishReview)
}
