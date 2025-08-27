package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/plugimt/transat-backend/handlers"
)

func SetupWhatsAppRoutes(app *fiber.App, whatsappHandler *handlers.WhatsAppHandler) {
	// route de test pour les templates/code de vérif (mais je doute que ça fonctionne avant un moment anyway vu la vérif)
	app.Get("/whatsapp/send-message", whatsappHandler.SendMessage)
}
