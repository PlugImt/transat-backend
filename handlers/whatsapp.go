package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/plugimt/transat-backend/services"
	"github.com/plugimt/transat-backend/utils"
)

type WhatsAppHandler struct {
	whatsappService *services.WhatsAppService
}

func NewWhatsAppHandler(whatsappService *services.WhatsAppService) *WhatsAppHandler {
	return &WhatsAppHandler{
		whatsappService: whatsappService,
	}
}

func (h *WhatsAppHandler) SendMessage(c *fiber.Ctx) error {
	code := c.Query("code")
	phoneNumber := c.Query("phoneNumber")
	language := c.Query("language")

	if code == "" || phoneNumber == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Code and phoneNumber parameters are required",
		})
	}

	if language == "" {
		language = "en_US" // pour tester
	}

	err := h.whatsappService.SendVerificationCode(phoneNumber, code, language)
	if err != nil {
		utils.LogMessage(utils.LevelError, "Error sending WhatsApp message")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to send WhatsApp message",
		})
	}

	return c.JSON(fiber.Map{
		"message": "message envoyé correctement",
	})
}
