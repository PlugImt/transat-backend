package admin

import (
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/plugimt/transat-backend/models"
	"github.com/plugimt/transat-backend/utils"
)

const (
	maxNotificationTitle   = 100
	maxNotificationMessage = 500
)

type sendNotificationRequest struct {
	Title    string                      `json:"title"`
	Message  string                      `json:"message"`
	Audience models.NotificationAudience `json:"audience"`
	// Service limits the audience to users who did not opt out of this category.
	Service models.NotificationCategory `json:"service,omitempty"`
	// Navigation is the optional screen opened on tap.
	Navigation *models.NavigationTarget `json:"navigation,omitempty"`
	// DryRun only counts the recipients.
	DryRun bool `json:"dryRun,omitempty"`
}

// SendNotification sends a custom push notification to an audience, or only counts it when dryRun is set.
func (h *AdminHandler) SendNotification(c *fiber.Ctx) error {
	var req sendNotificationRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request format"})
	}

	req.Title = strings.TrimSpace(req.Title)
	req.Message = strings.TrimSpace(req.Message)
	if err := validateNotificationRequest(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	recipients, err := h.Notifications.ResolveAudience(req.Audience, req.Service)
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to resolve notification audience")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to resolve audience"})
	}

	result := fiber.Map{"users": recipients.Users, "devices": len(recipients.Tokens), "sent": false}
	if req.DryRun {
		return c.JSON(result)
	}
	if len(recipients.Tokens) == 0 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": "No recipients with notifications enabled"})
	}

	err = h.Notifications.SendPushNotification(models.NotificationPayload{
		NotificationTokens: recipients.Tokens,
		Title:              req.Title,
		Message:            req.Message,
		Navigation:         req.Navigation,
	})
	utils.LogMessage(utils.LevelInfo, "Admin notification sent")
	utils.LogLineKeyValue(utils.LevelInfo, "Admin", c.Locals("email"))
	utils.LogLineKeyValue(utils.LevelInfo, "Devices", len(recipients.Tokens))
	if err != nil {
		utils.LogLineKeyValue(utils.LevelWarn, "Delivery errors", err)
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"error": err.Error(), "users": recipients.Users, "devices": len(recipients.Tokens)})
	}

	result["sent"] = true
	return c.JSON(result)
}

func validateNotificationRequest(req sendNotificationRequest) error {
	switch {
	case req.Title == "":
		return errors.New("Title is required")
	case len([]rune(req.Title)) > maxNotificationTitle:
		return errors.New("Title is too long")
	case len([]rune(req.Message)) > maxNotificationMessage:
		return errors.New("Message is too long")
	case req.Service != "" && !req.Service.IsValid():
		return errors.New("Unknown notification service")
	}
	if err := req.Audience.Validate(); err != nil {
		return err
	}
	if req.Navigation != nil {
		return req.Navigation.Validate()
	}
	return nil
}
