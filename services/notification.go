package services

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"Transat_2.0_Backend/models" // Assuming models are correctly placed
	"Transat_2.0_Backend/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/lib/pq"
)

// NotificationService handles sending notifications.
type NotificationService struct {
	db          *sql.DB
	expoPushURL string
}

// NewNotificationService creates a new NotificationService.
func NewNotificationService(db *sql.DB) *NotificationService {
	return &NotificationService{
		db:          db,
		expoPushURL: "https://api.expo.dev/v2/push/send",
	}
}

func (ns *NotificationService) GetNotificationTargets(userEmails []string, groups []string) ([]models.NotificationTarget, error) {
	log.Println("╔======== 📧 Get Notification Targets 📧 ========╗")
	var targets []models.NotificationTarget

	if len(userEmails) > 0 {
		query := `
            SELECT email, notification_token 
            FROM newf 
            WHERE email = ANY($1) 
            AND notification_token IS NOT NULL
        `
		rows, err := ns.db.Query(query, pq.Array(userEmails))
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

		stmt, err := ns.db.Prepare(query)
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

// GetSubscribedUsers retrieves users subscribed to a specific service.
func (ns *NotificationService) GetSubscribedUsers(serviceName string) ([]models.NotificationTarget, error) {
	query := `
		SELECT n.email, COALESCE(newf.notification_token, '') as notification_token
		FROM notifications n
		JOIN newf ON n.email = newf.email
		JOIN services s ON n.id_services = s.id_services
		WHERE s.name = $1 AND newf.notification_token IS NOT NULL AND newf.notification_token != '';
	`
	rows, err := ns.db.Query(query, serviceName)
	if err != nil {
		log.Printf("Error querying subscribed users for %s: %v", serviceName, err)
		return nil, fmt.Errorf("failed to query subscribed users: %w", err)
	}
	defer rows.Close()

	var targets []models.NotificationTarget
	for rows.Next() {
		var target models.NotificationTarget
		if err := rows.Scan(&target.Email, &target.NotificationToken); err != nil {
			log.Printf("Error scanning notification target: %v", err)
			continue // Skip this user on error
		}
		targets = append(targets, target)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error iterating over subscribed users rows: %v", err)
		return nil, fmt.Errorf("failed to process subscribed users list: %w", err)
	}

	return targets, nil
}

// SendPushNotification sends a push notification via Expo.
func (ns *NotificationService) SendPushNotification(payload models.NotificationPayload) error {
	// resolve the tokens from the db from the userEmails
	request := `SELECT notification_token FROM newf WHERE email = ANY($1)`
	rows, err := ns.db.Query(request, pq.Array(payload.UserEmails))
	if err != nil {
		log.Printf("Error querying tokens from db: %v", err)
		return fmt.Errorf("failed to query tokens from db: %w", err)
	}

	var tokens []string
	for rows.Next() {
		var token string
		if err := rows.Scan(&token); err != nil {
			log.Printf("Error scanning token: %v", err)
			continue
		}
		tokens = append(tokens, token)
	}

	expoPayload := map[string]interface{}{
		"to":    tokens,
		"title": payload.Title,
	}

	// Only include optional fields if they have values
	if payload.Message != "" {
		expoPayload["body"] = payload.Message
	}
	if payload.Sound != "" {
		expoPayload["sound"] = payload.Sound
	}
	if payload.ChannelID != "" {
		expoPayload["channelId"] = payload.ChannelID
	}
	if payload.Badge != 0 {
		expoPayload["badge"] = payload.Badge
	}
	if payload.Data != nil {
		expoPayload["data"] = payload.Data
	}

	payloadBytes, err := json.Marshal(expoPayload)
	if err != nil {
		log.Printf("Error marshalling Expo push notification payload: %v", err)
		return fmt.Errorf("failed to create notification request: %w", err)
	}

	req, err := http.NewRequest("POST", ns.expoPushURL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		log.Printf("Error creating Expo push request: %v", err)
		return fmt.Errorf("failed to create push request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Accept-Encoding", "gzip, deflate")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error sending push notification to Expo: %v", err)
		return fmt.Errorf("failed to send push notification: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// Log response body for debugging
		bodyBytes, _ := io.ReadAll(resp.Body)
		log.Printf("Error response from Expo (%d): %s", resp.StatusCode, string(bodyBytes))
		return fmt.Errorf("expo push notification failed with status code %d", resp.StatusCode)
	}

	log.Printf("Successfully sent push notification(s)")
	// Optionally parse response body for details on individual message status
	// var expoResponse map[string]interface{}
	// if err := json.NewDecoder(resp.Body).Decode(&expoResponse); err == nil {
	//  log.Printf("Expo Response: %+v", expoResponse)
	// }

	return nil
}

// SendNotification handles sending notifications based on payload details.
// This function remains a placeholder as its original implementation was complex
// and likely tied to specific handler logic (permissions, data fetching etc.)
func (ns *NotificationService) SendNotification(c *fiber.Ctx) error {
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

	targets, err := ns.GetNotificationTargets(payload.UserEmails, payload.Groups)
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
		if err := ns.SendPushNotification(payload); err != nil {
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

// SendDailyMenuNotification sends the daily menu notification to subscribed users.
// It reads messages from notifications.json, picks one randomly, and sends it.
func (ns *NotificationService) SendDailyMenuNotification() error {
	log.Println("╔======== 🍽️ Send Daily Menu Notification 🍽️ ========╗")

	subscribers, err := ns.GetSubscribedUsers("Restaurant") // Keep using existing function
	if err != nil {
		log.Println("║ 💥 Failed to get notification targets: ", err)
		log.Println("╚=========================================╝")
		return err
	}

	if len(subscribers) == 0 {
		log.Println("║ ℹ️ No users subscribed to RESTAURANT notifications")
		log.Println("╚=========================================╝")
		return nil
	}

	log.Printf("║ ℹ️ Found %d users subscribed to RESTAURANT notifications", len(subscribers))

	// Read messages from external JSON file
	file, err := os.ReadFile("notifications.json")
	if err != nil {
		log.Println("║ 💥 Failed to read notifications.json: ", err)
		log.Println("╚=========================================╝")
		// Consider default message or different error handling?
		return fmt.Errorf("failed to read notifications.json: %w", err)
	}

	var emails []string
	for _, sub := range subscribers {
		emails = append(emails, sub.Email)
	}

	// Parse messages
	var messages []string
	if err := json.Unmarshal(file, &messages); err != nil {
		log.Println("║ 💥 Failed to parse menu messages from notifications.json: ", err)
		log.Println("╚=========================================╝")
		return fmt.Errorf("failed to parse notifications.json: %w", err)
	}

	if len(messages) == 0 {
		log.Println("║ ⚠️ No messages found in notifications.json")
		log.Println("╚=========================================╝")
		return fmt.Errorf("no messages available in notifications.json")
	}

	// Randomize selection
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	randomMessage := messages[r.Intn(len(messages))]

	// Prepare base payload
	payload := models.NotificationPayload{
		UserEmails: emails,
		Title:      "Menu du jour disponible",
		Message:    randomMessage,
		Sound:      "default",
		ChannelID:  "default", // Ensure this channel exists on the client app
		Data: map[string]interface{}{
			"screen": "Restaurant",
		},
	}

	utils.LogLineKeyValue(utils.LevelInfo, "Sending daily menu notification to", len(subscribers))

	err = ns.SendPushNotification(payload)
	if err != nil {
		// Log individual failures but continue sending to others
		utils.LogLineKeyValue(utils.LevelError, "Failed to send menu notification", err)
		utils.LogFooter()
		return err
	}

	utils.LogLineKeyValue(utils.LevelInfo, "Daily menu notifications processed", len(subscribers))
	utils.LogFooter()

	return nil
}
