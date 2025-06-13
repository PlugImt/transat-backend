package routes

import (
	"github.com/gofiber/fiber/v2"
)

func SetupSystemRoutes(app *fiber.App) {
	// @Summary		Renvoie OK si l'API est en ligne
	// @Tags			System
	// @Success		200			{object}	map[string]string	{"status":"ok"}
	// @Failure		500			{string}	string	"Erreur serveur"
	// @Router			/api/health [get]
	app.Get("/health", func(c *fiber.Ctx) error {
		// TODO: add db check one day
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status": "ok",
		})
	})

	// @Summary		Vérifier le statut de l'API
	// @Description	Retourne le statut de fonctionnement de l'API
	// @Tags			System
	// @Produce		plain
	// @Success		200	{string}	string	"API is up and running"
	// @Router			/status [get]
	app.Get("/status", func(c *fiber.Ctx) error {
		return c.SendString("API is up and running")
	})
}
