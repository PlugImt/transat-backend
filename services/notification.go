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

	"github.com/plugimt/transat-backend/models" // Assuming models are correctly placed
	"github.com/plugimt/transat-backend/utils"

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

// GetSubscribedUsersWithLanguage retrieves users subscribed to a specific service with their language preferences.
func (ns *NotificationService) GetSubscribedUsersWithLanguage(serviceName string) ([]models.NotificationTargetWithLanguage, error) {
	query := `
		SELECT n.email, COALESCE(newf.notification_token, '') as notification_token, l.code as language_code
		FROM notifications n
		JOIN newf ON n.email = newf.email
		JOIN services s ON n.id_services = s.id_services
		JOIN languages l ON newf.language = l.id_languages
		WHERE s.name = $1 AND newf.notification_token IS NOT NULL AND newf.notification_token != '';
	`
	rows, err := ns.db.Query(query, serviceName)
	if err != nil {
		log.Printf("Error querying subscribed users with language for %s: %v", serviceName, err)
		return nil, fmt.Errorf("failed to query subscribed users with language: %w", err)
	}
	defer rows.Close()

	var targets []models.NotificationTargetWithLanguage
	for rows.Next() {
		var target models.NotificationTargetWithLanguage
		if err := rows.Scan(&target.Email, &target.NotificationToken, &target.LanguageCode); err != nil {
			log.Printf("Error scanning notification target with language: %v", err)
			continue // Skip this user on error
		}
		targets = append(targets, target)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error iterating over subscribed users with language rows: %v", err)
		return nil, fmt.Errorf("failed to process subscribed users with language list: %w", err)
	}

	return targets, nil
}

// SendPushNotification sends a push notification via Expo.
func (ns *NotificationService) SendPushNotification(payload models.NotificationPayload) error {
	var tokens []string

	// If tokens are directly provided, use them
	if len(payload.NotificationTokens) > 0 {
		tokens = payload.NotificationTokens
	} else if len(payload.UserEmails) > 0 {
		// Fallback to resolving tokens from emails if needed
		request := `SELECT notification_token FROM newf WHERE email = ANY($1)`
		rows, err := ns.db.Query(request, pq.Array(payload.UserEmails))
		if err != nil {
			log.Printf("Error querying tokens from db: %v", err)
			return fmt.Errorf("failed to query tokens from db: %w", err)
		}
		defer rows.Close()

		for rows.Next() {
			var token string
			if err := rows.Scan(&token); err != nil {
				log.Printf("Error scanning token: %v", err)
				continue
			}
			tokens = append(tokens, token)
		}
	}

	if len(tokens) == 0 {
		return fmt.Errorf("no valid notification tokens found")
	}

	// Track errors to return a summary at the end
	var failedTokens []string
	var lastError error

	// Send to each token individually
	for _, token := range tokens {
		expoPayload := map[string]interface{}{
			"to":    token, // Send to single token instead of array
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
		if payload.Subtitle != "" {
			expoPayload["subtitle"] = payload.Subtitle
		}
		if payload.TTL != 0 {
			expoPayload["ttl"] = payload.TTL
		}

		payloadBytes, err := json.Marshal(expoPayload)
		if err != nil {
			log.Printf("Error marshalling Expo push notification payload for token %s: %v", token, err)
			failedTokens = append(failedTokens, token)
			lastError = err
			continue
		}

		req, err := http.NewRequest("POST", ns.expoPushURL, bytes.NewBuffer(payloadBytes))
		if err != nil {
			log.Printf("Error creating Expo push request for token %s: %v", token, err)
			failedTokens = append(failedTokens, token)
			lastError = err
			continue
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")
		req.Header.Set("Accept-Encoding", "gzip, deflate")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Printf("Error sending push notification to Expo for token %s: %v", token, err)
			failedTokens = append(failedTokens, token)
			lastError = err
			continue
		}

		// Read response body
		bodyBytes, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			log.Printf("Error response from Expo (%d) for token %s: %s", resp.StatusCode, token, string(bodyBytes))
			failedTokens = append(failedTokens, token)
			lastError = fmt.Errorf("expo push notification failed with status code %d", resp.StatusCode)
			continue
		}

		// Check for errors in the response
		var responseData map[string]interface{}
		if err := json.Unmarshal(bodyBytes, &responseData); err == nil {
			// Check if the response contains data about errors
			if data, ok := responseData["data"].([]interface{}); ok && len(data) > 0 {
				if dataItem, ok := data[0].(map[string]interface{}); ok {
					if status, ok := dataItem["status"].(string); ok && status != "ok" {
						errMsg := "unknown error"
						if msg, ok := dataItem["message"].(string); ok {
							errMsg = msg
						}
						log.Printf("Expo returned error for token %s: %s", token, errMsg)
						failedTokens = append(failedTokens, token)
						lastError = fmt.Errorf("expo returned error: %s", errMsg)
						continue
					}
				}
			}
		}

		log.Printf("Successfully sent push notification to token: %s", token)
	}

	// Report results
	totalTokens := len(tokens)
	failedCount := len(failedTokens)

	if failedCount > 0 {
		log.Printf("Completed sending notifications with %d/%d failures", failedCount, totalTokens)
		return fmt.Errorf("failed to send %d/%d notifications: %w", failedCount, totalTokens, lastError)
	}

	log.Printf("Successfully sent all %d push notifications", totalTokens)
	return nil
}

// SendNotification handles sending notifications based on payload details.
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

	// Check if either tokens, emails, or groups are provided
	if len(payload.NotificationTokens) == 0 && len(payload.UserEmails) == 0 && len(payload.Groups) == 0 {
		log.Println("║ 💥 Notification tokens, user emails, or groups must be specified")
		log.Println("╚=========================================╝")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Notification tokens, user emails, or groups must be specified",
		})
	}

	// If we already have tokens, use them directly
	if len(payload.NotificationTokens) > 0 {
		err := ns.SendPushNotification(payload)
		totalTokens := len(payload.NotificationTokens)

		if err != nil {
			log.Println("║ ⚠️ Some push notifications failed to send: ", err)
			log.Println("╚=========================================╝")
			return c.Status(fiber.StatusPartialContent).JSON(fiber.Map{
				"success": false,
				"message": err.Error(),
				"count":   totalTokens,
			})
		}

		log.Println("║ ✅ All notifications sent successfully")
		log.Println("║ 📱 Total Tokens: ", totalTokens)
		log.Println("╚=========================================╝")
		return c.JSON(fiber.Map{
			"success": true,
			"count":   totalTokens,
		})
	}

	// Otherwise, get notification targets from emails or groups
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

	// Extract tokens from targets
	var tokens []string
	for _, target := range targets {
		if target.NotificationToken != "" {
			tokens = append(tokens, target.NotificationToken)
		}
	}

	if len(tokens) == 0 {
		log.Println("║ 💥 No valid notification tokens found")
		log.Println("╚=========================================╝")
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "No valid notification tokens found",
		})
	}

	// Update payload with resolved tokens
	payload.NotificationTokens = tokens

	err = ns.SendPushNotification(payload)
	totalTokens := len(tokens)

	if err != nil {
		log.Println("║ ⚠️ Some push notifications failed to send: ", err)
		log.Println("╚=========================================╝")
		return c.Status(fiber.StatusPartialContent).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
			"count":   totalTokens,
		})
	}

	log.Println("║ ✅ All notifications sent successfully")
	log.Println("║ 📱 Total Tokens: ", totalTokens)
	log.Println("╚=========================================╝")
	return c.JSON(fiber.Map{
		"success": true,
		"count":   totalTokens,
	})
}

// SendDailyMenuNotification sends the daily menu notification to subscribed users.
// It reads messages from notifications.json, picks one randomly, and sends it.
func (ns *NotificationService) SendDailyMenuNotification() error {
	log.Println("╔======== 🍽️ Send Daily Menu Notification 🍽️ ========╗")

	subscribers, err := ns.GetSubscribedUsers("RESTAURANT") // Keep using existing function
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

	// Extract notification tokens from subscribers
	var tokens []string
	for _, sub := range subscribers {
		if sub.NotificationToken != "" {
			tokens = append(tokens, sub.NotificationToken)
		}
	}

	if len(tokens) == 0 {
		log.Println("║ ℹ️ No valid notification tokens found")
		log.Println("╚=========================================╝")
		return nil
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
	r := rand.New(rand.NewSource(utils.UnixNanoParis(utils.Now())))
	randomMessage := messages[r.Intn(len(messages))]

	// Prepare base payload
	payload := models.NotificationPayload{
		NotificationTokens: tokens,
		Title:              "Menu du jour disponible",
		Message:            randomMessage,
		Sound:              "default",
		ChannelID:          "default", // Ensure this channel exists on the client app
		Data: map[string]interface{}{
			"screen": "Restaurant",
		},
	}

	utils.LogLineKeyValue(utils.LevelInfo, "Sending daily menu notification to", len(tokens))

	err = ns.SendPushNotification(payload)
	if err != nil {
		// Log error but don't consider it a complete failure
		// as some notifications might have been sent successfully
		utils.LogLineKeyValue(utils.LevelWarn, "Some menu notifications failed", err.Error())
		utils.LogFooter()
		return err
	}

	utils.LogLineKeyValue(utils.LevelInfo, "All daily menu notifications processed successfully", len(tokens))
	utils.LogFooter()

	return nil
}
