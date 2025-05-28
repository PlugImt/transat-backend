package services

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"time"

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
func NewNotificationService(db *sql.DB, expoPushURL string) *NotificationService {
	return &NotificationService{
		db:          db,
		expoPushURL: expoPushURL,
	}
}

func (ns *NotificationService) GetNotificationTargets(userEmails []string, groups []string) ([]models.NotificationTarget, error) {
	utils.LogHeader("📧 Get Notification Targets 📧")
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
			utils.LogMessage(utils.LevelError, fmt.Sprintf("Failed to query emails: %v", err))
			utils.LogFooter()
			return nil, fmt.Errorf("failed to query notification targets by email: %w", err)
		}
		defer func(rows *sql.Rows) {
			err := rows.Close()
			if err != nil {
				utils.LogMessage(utils.LevelError, fmt.Sprintf("Failed to close rows for email query: %v", err))
				// utils.LogFooter() // Avoid double footer if error occurs in main flow too
			}
		}(rows)

		for rows.Next() {
			var target models.NotificationTarget
			if err := rows.Scan(&target.Email, &target.NotificationToken); err != nil {
				utils.LogMessage(utils.LevelError, fmt.Sprintf("Failed to scan row for email target: %v", err))
				utils.LogFooter()
				return nil, fmt.Errorf("failed to scan notification target: %w", err)
			}
			targets = append(targets, target)
		}
		if err = rows.Err(); err != nil { // Check for errors during rows.Next() iteration
			utils.LogMessage(utils.LevelError, fmt.Sprintf("Error iterating email target rows: %v", err))
			utils.LogFooter()
			return nil, fmt.Errorf("error iterating notification target rows: %w", err)
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
			utils.LogMessage(utils.LevelError, fmt.Sprintf("Failed to prepare statement for group query: %v", err))
			utils.LogFooter()
			return nil, fmt.Errorf("failed to prepare statement for group targets: %w", err)
		}
		defer stmt.Close() // Ensure statement is closed

		rows, err := stmt.Query(pq.Array(groups))
		if err != nil {
			utils.LogMessage(utils.LevelError, fmt.Sprintf("Failed to query groups: %v", err))
			utils.LogFooter()
			return nil, fmt.Errorf("failed to query notification targets by group: %w", err)
		}

		defer func(rows *sql.Rows) {
			err := rows.Close()
			if err != nil {
				utils.LogMessage(utils.LevelError, fmt.Sprintf("Failed to close rows for group query: %v", err))
			}
		}(rows)

		for rows.Next() {
			var target models.NotificationTarget
			if err := rows.Scan(&target.Email, &target.NotificationToken); err != nil {
				utils.LogMessage(utils.LevelError, fmt.Sprintf("Failed to scan row for group target: %v", err))
				utils.LogFooter()
				return nil, fmt.Errorf("failed to scan notification target from group: %w", err)
			}
			targets = append(targets, target)
		}

		if err := rows.Err(); err != nil {
			utils.LogMessage(utils.LevelError, fmt.Sprintf("Error iterating group target rows: %v", err))
			utils.LogFooter()
			return nil, fmt.Errorf("error iterating group target rows: %w", err)
		}
	}

	utils.LogMessage(utils.LevelInfo, "Notification targets retrieved successfully")
	utils.LogFooter()
	return targets, nil
}

// GetSubscribedUsers retrieves users subscribed to a specific service.
func (ns *NotificationService) GetSubscribedUsers(serviceName string) ([]models.NotificationTarget, error) {
	utils.LogHeader("📧 Get Subscribed Users 📧")
	query := `
		SELECT n.email, COALESCE(newf.notification_token, '') as notification_token
		FROM notifications n
		JOIN newf ON n.email = newf.email
		JOIN services s ON n.id_services = s.id_services
		WHERE s.name = $1 AND newf.notification_token IS NOT NULL AND newf.notification_token != '';
	`
	rows, err := ns.db.Query(query, serviceName)
	if err != nil {
		utils.LogMessage(utils.LevelError, fmt.Sprintf("Error querying subscribed users for %s: %v", serviceName, err))
		utils.LogFooter()
		return nil, fmt.Errorf("failed to query subscribed users: %w", err)
	}
	defer rows.Close()

	var targets []models.NotificationTarget
	for rows.Next() {
		var target models.NotificationTarget
		if err := rows.Scan(&target.Email, &target.NotificationToken); err != nil {
			utils.LogMessage(utils.LevelWarn, fmt.Sprintf("Error scanning notification target for user %s, skipping: %v", target.Email, err))
			continue // Skip this user on error
		}
		targets = append(targets, target)
	}

	if err = rows.Err(); err != nil {
		utils.LogMessage(utils.LevelError, fmt.Sprintf("Error iterating over subscribed users rows: %v", err))
		utils.LogFooter()
		return nil, fmt.Errorf("failed to process subscribed users list: %w", err)
	}

	utils.LogMessage(utils.LevelInfo, fmt.Sprintf("Successfully retrieved %d subscribed users for service %s", len(targets), serviceName))
	utils.LogFooter()
	return targets, nil
}

// SendPushNotification sends a push notification via Expo.
func (ns *NotificationService) SendPushNotification(payload models.NotificationPayload) error {
	utils.LogHeader("🚀 Sending Push Notification 🚀")
	var tokens []string

	if len(payload.NotificationTokens) > 0 {
		tokens = payload.NotificationTokens
	} else if len(payload.UserEmails) > 0 {
		utils.LogMessage(utils.LevelInfo, fmt.Sprintf("No direct tokens provided, resolving from %d user emails", len(payload.UserEmails)))
		request := `SELECT notification_token FROM newf WHERE email = ANY($1) AND notification_token IS NOT NULL AND notification_token != ''`
		rows, err := ns.db.Query(request, pq.Array(payload.UserEmails))
		if err != nil {
			utils.LogMessage(utils.LevelError, fmt.Sprintf("Error querying tokens from db: %v", err))
			utils.LogFooter()
			return fmt.Errorf("failed to query tokens from db: %w", err)
		}
		defer rows.Close()

		for rows.Next() {
			var token string
			if err := rows.Scan(&token); err != nil {
				utils.LogMessage(utils.LevelWarn, fmt.Sprintf("Error scanning token, skipping: %v", err))
				continue
			}
			tokens = append(tokens, token)
		}
		if err = rows.Err(); err != nil { // Check for errors during rows.Next() iteration
			utils.LogMessage(utils.LevelError, fmt.Sprintf("Error iterating token rows: %v", err))
			utils.LogFooter()
			return fmt.Errorf("error iterating token rows: %w", err)
		}
		utils.LogMessage(utils.LevelInfo, fmt.Sprintf("Resolved %d tokens from emails", len(tokens)))
	}

	if len(tokens) == 0 {
		utils.LogMessage(utils.LevelWarn, "No valid notification tokens found after resolving.")
		utils.LogFooter()
		return fmt.Errorf("no valid notification tokens found") // No need to wrap, it's a specific condition
	}

	var failedTokens []string
	var firstError error // Changed from lastError to firstError for more conventional error reporting
	var successfulSends int

	utils.LogMessage(utils.LevelInfo, fmt.Sprintf("Attempting to send notifications to %d tokens.", len(tokens)))

	for _, token := range tokens {
		expoPayload := map[string]interface{}{
			"to":    token,
			"title": payload.Title,
		}

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
			utils.LogMessage(utils.LevelError, fmt.Sprintf("Error marshalling Expo payload for token %s: %v", token, err))
			failedTokens = append(failedTokens, token)
			if firstError == nil {
				firstError = fmt.Errorf("failed to marshal payload for token %s: %w", token, err)
			}
			continue
		}

		req, err := http.NewRequest("POST", ns.expoPushURL, bytes.NewBuffer(payloadBytes))
		if err != nil {
			utils.LogMessage(utils.LevelError, fmt.Sprintf("Error creating Expo push request for token %s: %v", token, err))
			failedTokens = append(failedTokens, token)
			if firstError == nil {
				firstError = fmt.Errorf("failed to create request for token %s: %w", token, err)
			}
			continue
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")
		req.Header.Set("Accept-Encoding", "gzip, deflate")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			utils.LogMessage(utils.LevelError, fmt.Sprintf("Error sending push notification to Expo for token %s: %v", token, err))
			failedTokens = append(failedTokens, token)
			if firstError == nil {
				firstError = fmt.Errorf("failed to send push notification to Expo for token %s: %w", token, err)
			}
			continue
		}

		// Read response body
		bodyBytes, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			utils.LogMessage(utils.LevelError, fmt.Sprintf("Error response from Expo (%d) for token %s: %s", resp.StatusCode, token, string(bodyBytes)))
			failedTokens = append(failedTokens, token)
			if firstError == nil {
				firstError = fmt.Errorf("expo push notification failed with status code %d", resp.StatusCode)
			}
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
						utils.LogMessage(utils.LevelError, fmt.Sprintf("Expo returned error for token %s: %s", token, errMsg))
						failedTokens = append(failedTokens, token)
						if firstError == nil {
							firstError = fmt.Errorf("expo returned error: %s", errMsg)
						}
						continue
					}
				}
			}
		}

		utils.LogMessage(utils.LevelInfo, fmt.Sprintf("Successfully sent push notification to token: %s", token))
		successfulSends++
	}

	// Report results
	totalTokens := len(tokens)
	failedCount := len(failedTokens)

	if failedCount > 0 {
		utils.LogMessage(utils.LevelWarn, fmt.Sprintf("Completed sending notifications with %d/%d failures", failedCount, totalTokens))
		utils.LogFooter()
		return fmt.Errorf("failed to send %d/%d notifications: %w", failedCount, totalTokens, firstError)
	}

	utils.LogMessage(utils.LevelInfo, fmt.Sprintf("Successfully sent all %d push notifications", totalTokens))
	utils.LogFooter()
	return nil
}

// SendNotification handles sending notifications based on payload details.
func (ns *NotificationService) SendNotification(c *fiber.Ctx) error {
	utils.LogHeader("📤 Send Notification 📤")
	var payload models.NotificationPayload
	if err := c.BodyParser(&payload); err != nil {
		utils.LogMessage(utils.LevelError, fmt.Sprintf("Failed to parse request body: %v", err))
		utils.LogFooter()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if payload.Title == "" {
		utils.LogMessage(utils.LevelWarn, "Title is required for notification")
		utils.LogFooter()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Title is required",
		})
	}

	if len(payload.NotificationTokens) == 0 && len(payload.UserEmails) == 0 && len(payload.Groups) == 0 {
		utils.LogMessage(utils.LevelWarn, "Notification tokens, user emails, or groups must be specified")
		utils.LogFooter()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Notification tokens, user emails, or groups must be specified",
		})
	}

	if len(payload.NotificationTokens) > 0 {
		err := ns.SendPushNotification(payload)
		totalTokens := len(payload.NotificationTokens)

		if err != nil {
			utils.LogMessage(utils.LevelWarn, fmt.Sprintf("Some push notifications failed to send: %v", err))
			utils.LogFooter()
			return c.Status(fiber.StatusPartialContent).JSON(fiber.Map{
				"success": false,
				"message": err.Error(),
				"count":   totalTokens,
			})
		}

		utils.LogMessage(utils.LevelInfo, fmt.Sprintf("All %d notifications sent successfully via direct tokens", totalTokens))
		utils.LogFooter()
		return c.JSON(fiber.Map{
			"success": true,
			"count":   totalTokens,
		})
	}

	targets, err := ns.GetNotificationTargets(payload.UserEmails, payload.Groups)
	if err != nil {
		utils.LogMessage(utils.LevelError, fmt.Sprintf("Failed to get notification targets: %v", err))
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get notification targets",
		})
	}

	if len(targets) == 0 {
		utils.LogMessage(utils.LevelWarn, "No valid notification targets found from emails/groups")
		utils.LogFooter()
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "No valid notification targets found",
		})
	}

	var tokens []string
	for _, target := range targets {
		if target.NotificationToken != "" {
			tokens = append(tokens, target.NotificationToken)
		}
	}

	if len(tokens) == 0 {
		utils.LogMessage(utils.LevelWarn, "No valid notification tokens extracted from targets")
		utils.LogFooter()
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "No valid notification tokens found",
		})
	}

	payload.NotificationTokens = tokens

	err = ns.SendPushNotification(payload)
	totalTokensAfterResolve := len(tokens)

	if err != nil {
		utils.LogMessage(utils.LevelWarn, fmt.Sprintf("Some push notifications failed to send after resolving targets: %v", err))
		utils.LogFooter()
		return c.Status(fiber.StatusPartialContent).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
			"count":   totalTokensAfterResolve,
		})
	}

	utils.LogMessage(utils.LevelInfo, fmt.Sprintf("All %d notifications sent successfully after resolving targets", totalTokensAfterResolve))
	utils.LogFooter()
	return c.JSON(fiber.Map{
		"success": true,
		"count":   totalTokensAfterResolve,
	})
}

// SendDailyMenuNotification sends the daily menu notification to subscribed users.
func (ns *NotificationService) SendDailyMenuNotification() error {
	utils.LogHeader("🍽️ Send Daily Menu Notification 🍽️")

	subscribers, err := ns.GetSubscribedUsers("RESTAURANT")
	if err != nil {
		utils.LogMessage(utils.LevelError, fmt.Sprintf("Failed to get RESTAURANT notification subscribers: %v", err))
		utils.LogFooter()
		return fmt.Errorf("failed to get RESTAURANT subscribers: %w", err)
	}

	if len(subscribers) == 0 {
		utils.LogMessage(utils.LevelInfo, "No users subscribed to RESTAURANT notifications")
		utils.LogFooter()
		return nil
	}

	utils.LogMessage(utils.LevelInfo, fmt.Sprintf("Found %d users subscribed to RESTAURANT notifications", len(subscribers)))

	file, err := os.ReadFile("notifications.json")
	if err != nil {
		utils.LogMessage(utils.LevelError, fmt.Sprintf("Failed to read notifications.json: %v", err))
		utils.LogFooter()
		return fmt.Errorf("failed to read notifications.json: %w", err)
	}

	var tokens []string
	for _, sub := range subscribers {
		if sub.NotificationToken != "" {
			tokens = append(tokens, sub.NotificationToken)
		}
	}

	if len(tokens) == 0 {
		utils.LogMessage(utils.LevelInfo, "No valid notification tokens found for RESTAURANT subscribers")
		utils.LogFooter()
		return nil
	}

	var messages []string
	if err := json.Unmarshal(file, &messages); err != nil {
		utils.LogMessage(utils.LevelError, fmt.Sprintf("Failed to parse menu messages from notifications.json: %v", err))
		utils.LogFooter()
		return fmt.Errorf("failed to parse notifications.json: %w", err)
	}

	if len(messages) == 0 {
		utils.LogMessage(utils.LevelWarn, "No messages found in notifications.json")
		utils.LogFooter()
		return fmt.Errorf("no messages available in notifications.json")
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	randomMessage := messages[r.Intn(len(messages))]

	payload := models.NotificationPayload{
		NotificationTokens: tokens,
		Title:              "Menu du jour disponible",
		Message:            randomMessage,
		Sound:              "default",
		ChannelID:          "default",
		Data: map[string]interface{}{
			"screen": "Restaurant",
		},
	}

	utils.LogMessage(utils.LevelInfo, fmt.Sprintf("Sending daily menu notification to %d tokens", len(tokens)))

	err = ns.SendPushNotification(payload)
	if err != nil {
		utils.LogMessage(utils.LevelWarn, fmt.Sprintf("Some daily menu notifications may have failed: %v", err.Error()))
		utils.LogFooter()
		return err // Return the error from SendPushNotification
	}

	utils.LogMessage(utils.LevelInfo, fmt.Sprintf("Daily menu notifications processed for %d tokens.", len(tokens)))
	utils.LogFooter()

	return nil
}
