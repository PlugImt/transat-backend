package admin

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/plugimt/transat-backend/services/push"
	"github.com/plugimt/transat-backend/utils"
)

// queueNotification is SendNotification for the delivery pipeline: it counts the devices that qualify now and
// queues the message, without sending anything from the request. The pipeline sends later and checks every device
// again at that moment, so the counts can only be an estimate and "sent" stays false.
func (h *AdminHandler) queueNotification(c *fiber.Ctx, req sendNotificationRequest) error {
	pipeline := h.Notifications.Pipeline()
	audience := push.Audience{NotificationAudience: req.Audience}

	recipients, err := pipeline.Resolve(c.Context(), audience, req.Service)
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to resolve notification audience")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to resolve audience"})
	}
	withoutDevice := recipients.WithoutTarget
	if withoutDevice == nil {
		withoutDevice = []string{}
	}
	result := fiber.Map{
		"users":         recipients.Users,
		"devices":       len(recipients.Targets),
		"withoutDevice": withoutDevice,
		"filtered":      recipients.Filtered,
		"sent":          false,
	}
	if req.DryRun {
		return c.JSON(result)
	}
	if len(recipients.Targets) == 0 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": "None of the recipients has an eligible device"})
	}

	email, _ := c.Locals("email").(string)
	queued, err := pipeline.Enqueue(c.Context(), push.Message{
		Kind:       "admin",
		Category:   req.Service,
		Audience:   audience,
		Content:    push.Content{push.DefaultLanguage: {Title: req.Title, Body: req.Message}},
		Navigation: req.Navigation,
		ImageURL:   req.ImageURL,
		CreatedBy:  email,
	})
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to queue admin notification")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to queue the notification"})
	}

	utils.LogMessage(utils.LevelInfo, "Admin notification queued")
	utils.LogLineKeyValue(utils.LevelInfo, "Admin", email)
	utils.LogLineKeyValue(utils.LevelInfo, "Message", queued.ID)
	result["queued"] = true
	result["id"] = queued.ID
	return c.Status(fiber.StatusAccepted).JSON(result)
}

// GetPushMessages handles GET /admin/notifications/messages: the diagnostics of the latest queued messages.
func (h *AdminHandler) GetPushMessages(c *fiber.Ctx) error {
	pipeline := h.Notifications.Pipeline()
	if pipeline == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "The delivery pipeline is not enabled"})
	}
	messages, err := pipeline.Recent(c.Context(), c.QueryInt("limit", 20))
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to load push messages")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to load push messages"})
	}
	return c.JSON(fiber.Map{"messages": messages})
}

// GetPushMessage handles GET /admin/notifications/messages/:id.
func (h *AdminHandler) GetPushMessage(c *fiber.Ctx) error {
	pipeline := h.Notifications.Pipeline()
	if pipeline == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "The delivery pipeline is not enabled"})
	}
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil || id < 1 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid message id"})
	}
	summary, err := pipeline.Summary(c.Context(), id)
	if errors.Is(err, push.ErrMessageNotFound) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Message not found"})
	}
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to load push message")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to load the message"})
	}
	return c.JSON(summary)
}
