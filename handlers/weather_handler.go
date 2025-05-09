package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/plugimt/transat-backend/services"
	"github.com/plugimt/transat-backend/utils"
)

type WeatherHandler struct {
	weatherService *services.WeatherService
}

func NewWeatherHandler(weatherService *services.WeatherService) *WeatherHandler {
	return &WeatherHandler{
		weatherService: weatherService,
	}
}

func (h *WeatherHandler) GetWeather(c *fiber.Ctx) error {
	lang := c.Query("language")

	weatherData, err := h.weatherService.GetWeather(lang)
	if err != nil {
		utils.LogMessage(utils.LevelError, "Error getting weather data")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve weather data",
		})
	}
	return c.JSON(weatherData)
}
