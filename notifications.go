package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"Transat_2.0_Backend/models"
	"github.com/gofiber/fiber/v2"
	"github.com/lib/pq"
)

type NotificationService struct {
	db *sql.DB
}

func NewNotificationService(db *sql.DB) *NotificationService {
	return &NotificationService{
		db: db,
	}
}

func (s *NotificationService) SendNotification(target models.NotificationTarget, payload models.NotificationPayload) error {
	if payload.Data == nil {
		payload.Data = map[string]interface{}{
			"sentTime": time.Now().Format(time.RFC3339),
		}
	}

	log.Println("╔======== 📤 Send Notification 📤 ========╗")
	notificationBody := map[string]interface{}{
		"to":        target.NotificationToken,
		"title":     payload.Title,
		"ttl":       payload.TTL,
		"subtitle":  payload.Subtitle,
		"sound":     payload.Sound,
		"body":      payload.Message,
		"channelId": payload.ChannelID,
		"badge":     payload.Badge,
		"data":      payload.Data,
	}

	notificationBodyBytes, err := json.Marshal(notificationBody)
	if err != nil {
		return err
	}

	response, err := http.Post("https://api.expo.dev/v2/push/send", "application/json", bytes.NewBuffer(notificationBodyBytes))

	if response != nil {
		defer func(response *http.Response) {
			err := response.Body.Close()
			if err != nil {
				log.Printf("║ 💥 Error closing response body: %v\n", err)
			}
		}(response)

		if response.StatusCode != http.StatusOK {
			log.Printf("║ 💥 Failed to send notification to %s | %s | %s\n", target.Email, notificationBody, response.Status)
			log.Println("╚=========================================╝")
			return err
		}

	} else {
		if err != nil {
			log.Printf("║ 💥 Error sending notification to %s: %v\n", target.Email, err)
			log.Println("╚=========================================╝")
			return err

		}
	}

	log.Printf("║ ✅ Notification sent successfully to %s | %s\n", target.Email, notificationBody)
	log.Println("║ 📧 Email: ", target.Email)
	log.Println("╚=========================================╝")
	return nil
}

func (s *NotificationService) GetNotificationTargets(emails []string, groups []string) ([]models.NotificationTarget, error) {
	log.Println("╔======== 📧 Get Notification Targets 📧 ========╗")
	var targets []models.NotificationTarget

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
			var target models.NotificationTarget
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
			var target models.NotificationTarget
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
	var payload models.NotificationPayload
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

	// Read messages from external JSON file
	file, err := os.ReadFile("notifications.json")
	if err != nil {
		log.Println("║ 💥 Failed to read notification.json: ", err)
		log.Println("╚=========================================╝")
		return err
	}

	// Parse messages
	var messages []string
	if err := json.Unmarshal(file, &messages); err != nil {
		log.Println("║ 💥 Failed to parse menu messages: ", err)
		log.Println("╚=========================================╝")
		return err
	}

	// Randomize selection
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	randomMessage := messages[r.Intn(len(messages))]

	payload := models.NotificationPayload{
		Title:   "Menu du jour disponible",
		Message: randomMessage,
		Data: map[string]interface{}{
			"screen": "Restaurant",
		},
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
