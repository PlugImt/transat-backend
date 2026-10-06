package user

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/plugimt/transat-backend/models"
	"github.com/plugimt/transat-backend/services"
	"github.com/plugimt/transat-backend/utils"
)

// GetNotificationPreferences returns every notification category with the user's choice.
func (h *UserHandler) GetNotificationPreferences(c *fiber.Ctx) error {
	email := c.Locals("email").(string)

	prefs, err := h.NotifService.GetPreferences(email)
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to get notification preferences")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get notification preferences"})
	}

	return c.JSON(fiber.Map{"preferences": prefs})
}

// SetNotificationPreference switches one notification category on or off for the user.
func (h *UserHandler) SetNotificationPreference(c *fiber.Ctx) error {
	email := c.Locals("email").(string)

	var req struct {
		Service models.NotificationCategory `json:"service"`
		Enabled *bool                       `json:"enabled"`
	}
	if err := c.BodyParser(&req); err != nil || req.Enabled == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Expected {\"service\": \"...\", \"enabled\": true|false}"})
	}
	if !req.Service.IsValid() {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Unknown notification service"})
	}

	if err := h.NotifService.SetPreference(email, req.Service, *req.Enabled); err != nil {
		if errors.Is(err, services.ErrUnknownNotificationCategory) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Notification service not available"})
		}
		utils.LogMessage(utils.LevelError, "Failed to save notification preference")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to save notification preference"})
	}

	return c.JSON(models.NotificationPreference{Service: req.Service, Enabled: *req.Enabled})
}
