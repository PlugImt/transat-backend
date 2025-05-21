package routes

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"github.com/plugimt/transat-backend/middlewares"
	"github.com/plugimt/transat-backend/models"
	"github.com/plugimt/transat-backend/services"
)

// SetupNotificationRoutes configures routes for notifications
func SetupNotificationRoutes(router fiber.Router, db *sql.DB, notificationService *services.NotificationService) {
	// Create a notifications group
	notificationGroup := router.Group("/notifications")

	// Apply authentication middleware
	notificationGroup.Use(middlewares.JWTMiddleware)

	// route 1: Send notification to a specific user by email
	notificationGroup.Post("/send-to-user", func(c *fiber.Ctx) error {
		type SendToUserRequest struct {
			Email   string                 `json:"email"`
			Title   string                 `json:"title"`
			Message string                 `json:"message"`
			Data    map[string]interface{} `json:"data,omitempty"`
		}

		var req SendToUserRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request format",
			})
		}

		if req.Email == "" || req.Title == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Email and title are required",
			})
		}

		// Create notification payload
		payload := models.NotificationPayload{
			UserEmails: []string{req.Email},
			Title:      req.Title,
			Message:    req.Message,
			Data:       req.Data,
		}

		// Send notification using existing service method
		err := notificationService.SendPushNotification(payload)
		if err != nil {
			return c.Status(fiber.StatusPartialContent).JSON(fiber.Map{
				"success": false,
				"message": err.Error(),
			})
		}

		return c.JSON(fiber.Map{
			"success": true,
			"message": "Notification sent successfully",
		})
	})

	// route 2: Send notification to a group of users by service name
	notificationGroup.Post("/send-to-group", func(c *fiber.Ctx) error {
		type SendToGroupRequest struct {
			ServiceName string                 `json:"serviceName"`
			Title       string                 `json:"title"`
			Message     string                 `json:"message"`
			Data        map[string]interface{} `json:"data,omitempty"`
		}

		var req SendToGroupRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request format",
			})
		}

		if req.ServiceName == "" || req.Title == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Service name and title are required",
			})
		}

		// Get subscribers for the service
		subscribers, err := notificationService.GetSubscribedUsers(req.ServiceName)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   "Failed to get subscribers",
				"message": err.Error(),
			})
		}

		if len(subscribers) == 0 {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "No subscribers found for this service",
			})
		}

		// Extract notification tokens
		var tokens []string
		for _, sub := range subscribers {
			if sub.NotificationToken != "" {
				tokens = append(tokens, sub.NotificationToken)
			}
		}

		if len(tokens) == 0 {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "No valid notification tokens found",
			})
		}

		// Create and send notification
		payload := models.NotificationPayload{
			NotificationTokens: tokens,
			Title:              req.Title,
			Message:            req.Message,
			Data:               req.Data,
		}

		err = notificationService.SendPushNotification(payload)
		if err != nil {
			return c.Status(fiber.StatusPartialContent).JSON(fiber.Map{
				"success": false,
				"message": err.Error(),
				"count":   len(tokens),
			})
		}

		return c.JSON(fiber.Map{
			"success": true,
			"message": "Notification sent to group successfully",
			"count":   len(tokens),
		})
	})
}
