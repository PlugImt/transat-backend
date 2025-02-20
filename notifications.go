package main

import (
	"database/sql"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/lib/pq"
	"log"
)

func NewNotificationService(db *sql.DB) *NotificationService {
	return &NotificationService{
		db:     db,
		client: resty.New(),
	}
}

func (s *NotificationService) SendNotification(target NotificationTarget, payload NotificationPayload) error {
	notificationBody := map[string]interface{}{
		"to": target.NotificationToken,
		"notification": map[string]interface{}{
			"title": payload.Title,
			"body":  payload.Message,
			"android": map[string]interface{}{
				"channelId": "default",
				"imageUrl":  payload.ImageURL,
			},
			"data": map[string]interface{}{
				"screen": payload.Screen,
			},
		},
	}

	fmt.Println(notificationBody)

	_, err := s.client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(notificationBody).
		Post("https://exp.host/--/api/v2/push/send")

	if err != nil {
		log.Printf("Error sending notification to %s: %v\n", target.Email, err)
		return err
	}

	return nil
}

func (s *NotificationService) GetNotificationTargets(emails []string, groups []string) ([]NotificationTarget, error) {
	var targets []NotificationTarget

	if len(emails) > 0 {
		query := `
            SELECT email, notification_token 
            FROM newf 
            WHERE email = ANY($1) 
            AND notification_token IS NOT NULL
        `
		rows, err := s.db.Query(query, pq.Array(emails))
		if err != nil {
			return nil, err
		}
		defer func(rows *sql.Rows) {
			err := rows.Close()
			if err != nil {
				log.Println(err)
			}
		}(rows)

		for rows.Next() {
			var target NotificationTarget
			if err := rows.Scan(&target.Email, &target.NotificationToken); err != nil {
				return nil, err
			}
			targets = append(targets, target)
		}
	}

	if len(groups) > 0 {
		query := `
            SELECT DISTINCT n.email, nf.notification_token 
            FROM notifications n
            JOIN newf nf ON n.email = nf.email
            JOIN services s ON n.id_services = s.id_services
            WHERE s.name = ANY($1)
            AND nf.notification_token IS NOT NULL
        `
		rows, err := s.db.Query(query, pq.Array(groups))
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		for rows.Next() {
			var target NotificationTarget
			if err := rows.Scan(&target.Email, &target.NotificationToken); err != nil {
				return nil, err
			}
			targets = append(targets, target)
		}
	}

	return targets, nil
}

func sendNotification(c *fiber.Ctx) error {
	var payload NotificationPayload
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if payload.Title == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Title is required",
		})
	}

	if len(payload.UserEmails) == 0 && len(payload.Groups) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "At least one user email or group must be specified",
		})
	}

	notificationService := NewNotificationService(db)

	targets, err := notificationService.GetNotificationTargets(payload.UserEmails, payload.Groups)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get notification targets",
		})
	}

	if len(targets) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "No valid notification targets found",
		})
	}

	var failedTargets []string
	for _, target := range targets {
		if err := notificationService.SendNotification(target, payload); err != nil {
			failedTargets = append(failedTargets, target.Email)
		}
	}

	return c.JSON(fiber.Map{
		"success":       true,
		"totalTargets":  len(targets),
		"failedTargets": failedTargets,
	})
}
