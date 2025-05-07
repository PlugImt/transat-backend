package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/plugimt/transat-backend/handlers"
)

func SetupWeatherRoutes(api fiber.Router, weatherHandler *handlers.WeatherHandler) {
	weatherGroup := api.Group("/weather")
	weatherGroup.Get("/", weatherHandler.GetWeather)
}
