package main

import (
	"database/sql"
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
	log.Println("╔======== 📤 Send Notification 📤 ========╗")
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
	
	_, err := s.client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(notificationBody).
		Post("https://exp.host/--/api/v2/push/send")

	if err != nil {
		log.Printf("║ 💥 Error sending notification to %s: %v\n", target.Email, err)
		log.Println("╚=========================================╝")
		return err
	}

	log.Println("║ ✅ Notification sent successfully")
	log.Println("║ 📧 Email: ", target.Email)
	log.Println("╚=========================================╝")
	return nil
}

func (s *NotificationService) GetNotificationTargets(emails []string, groups []string) ([]NotificationTarget, error) {
	log.Println("╔======== 📧 Get Notification Targets 📧 ========╗")
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
			log.Println("║ 💥 Failed to query emails: ", err)
			log.Println("╚=========================================╝")
			return nil, err
		}
		defer func(rows *sql.Rows) {
			err := rows.Close()
			if err != nil {
				log.Println("║ 💥 Failed to close rows: ", err)
				log.Println("╚=========================================╝")
			}
		}(rows)

		for rows.Next() {
			var target NotificationTarget
			if err := rows.Scan(&target.Email, &target.NotificationToken); err != nil {
				log.Println("║ 💥 Failed to scan row: ", err)
				log.Println("╚=========================================╝")
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

		stmt, err := s.db.Prepare(query)
		if err != nil {
			log.Println("║ 💥 Failed to prepare statement: ", err)
			log.Println("╚=========================================╝")
			return nil, err
		}

		rows, err := stmt.Query(pq.Array(groups))
		if err != nil {
			log.Println("║ 💥 Failed to query groups: ", err)
			log.Println("╚=========================================╝")
			return nil, err
		}

		defer func(rows *sql.Rows) {
			err := rows.Close()
			if err != nil {
				log.Println("║ 💥 Failed to close rows: ", err)
				log.Println("╚=========================================╝")
			}
		}(rows)

		for rows.Next() {
			var target NotificationTarget
			if err := rows.Scan(&target.Email, &target.NotificationToken); err != nil {
				log.Println("║ 💥 Failed to scan row: ", err)
				log.Println("╚=========================================╝")
				return nil, err
			}
			targets = append(targets, target)
		}

		if err := rows.Err(); err != nil {
			log.Println("║ 💥 Failed to iterate rows: ", err)
			log.Println("╚=========================================╝")
			return nil, err
		}

		if err := stmt.Close(); err != nil {
			log.Println("║ 💥 Failed to close statement: ", err)
			log.Println("╚=========================================╝")
			return nil, err
		}

	}

	log.Println("║ ✅ Notification targets retrieved successfully")
	log.Println("╚=========================================╝")
	return targets, nil
}

func sendNotification(c *fiber.Ctx) error {
	log.Println("╔======== 📤 Send Notification 📤 ========╗")
	var payload NotificationPayload
	if err := c.BodyParser(&payload); err != nil {
		log.Println("║ 💥 Failed to parse request body: ", err)
		log.Println("╚=========================================╝")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if payload.Title == "" {
		log.Println("║ 💥 Title is required")
		log.Println("╚=========================================╝")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Title is required",
		})
	}

	if len(payload.UserEmails) == 0 && len(payload.Groups) == 0 {
		log.Println("║ 💥 At least one user email or group must be specified")
		log.Println("╚=========================================╝")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "At least one user email or group must be specified",
		})
	}

	notificationService := NewNotificationService(db)

	targets, err := notificationService.GetNotificationTargets(payload.UserEmails, payload.Groups)
	if err != nil {
		log.Println("║ 💥 Failed to get notification targets: ", err)
		log.Println("╚=========================================╝")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get notification targets",
		})
	}

	if len(targets) == 0 {
		log.Println("║ 💥 No valid notification targets found")
		log.Println("╚=========================================╝")
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

	log.Println("║ ✅ Notification sent successfully")
	log.Println("║ 📧 Total Targets: ", len(targets))
	log.Println("║ 📧 Failed Targets: ", failedTargets)
	log.Println("╚=========================================╝")
	return c.JSON(fiber.Map{
		"success":       true,
		"totalTargets":  len(targets),
		"failedTargets": failedTargets,
	})
}

func (s *NotificationService) SendDailyMenuNotification() error {
	log.Println("╔======== 🍽️ Send Daily Menu Notification 🍽️ ========╗")

	targets, err := s.GetNotificationTargets([]string{}, []string{"RESTAURANT"})
	if err != nil {
		log.Println("║ 💥 Failed to get notification targets: ", err)
		log.Println("╚=========================================╝")
		return err
	}

	log.Printf("║ ℹ️ Found %d users subscribed to RESTAURANT notifications", len(targets))

	payload := NotificationPayload{
		Title:   "Menu du jour disponible !",
		Message: "Le menu du RU est maintenant disponible. Découvrez-le sur Transat !",
		Screen:  "Restaurant",
	}

	for _, target := range targets {
		err := s.SendNotification(target, payload)
		if err != nil {
			log.Printf("║ ⚠️ Failed to send menu notification to %s: %v", target.Email, err)
		}
	}

	log.Println("║ ✅ Daily menu notifications sent successfully")
	log.Println("╚=========================================╝")
	return nil
}
