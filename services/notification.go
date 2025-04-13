package services

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"Transat_2.0_Backend/models" // Assuming models are correctly placed
	"Transat_2.0_Backend/utils"

	"github.com/gofiber/fiber/v2"
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
	// Expo expects the 'to' field to contain the list of push tokens
	// We map UserEmails (which contains tokens in our current structure) to the 'to' field.
	expoPayload := map[string]interface{}{
		"to":        payload.UserEmails, // Use the tokens here
		"title":     payload.Title,
		"body":      payload.Message,
		"sound":     payload.Sound,
		"channelId": payload.ChannelID,
		"badge":     payload.Badge,
		"data":      payload.Data,
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
	// Add Authorization if needed (Expo might require tokens for certain operations)
	// req.Header.Set("Authorization", "Bearer YOUR_EXPO_ACCESS_TOKEN")

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
	log.Println("╔======== 🚀 Send Notification Handler (Placeholder) 🚀 ========╗")
	// TODO: Implement the logic from the original `sendNotification` handler.
	// This involves:
	// 1. Parsing the request body (title, message, targets like emails/groups).
	// 2. Checking permissions (e.g., is the user allowed to send notifications?).
	// 3. Resolving targets (emails/groups) to actual push tokens.
	// 4. Constructing the `models.NotificationPayload`.
	// 5. Calling `ns.SendPushNotification(payload)`.

	var reqPayload models.NotificationPayload // Example structure
	if err := c.BodyParser(&reqPayload); err != nil {
		utils.LogMessage(utils.LevelError, "Failed to parse notification request")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
	}

	// --- Placeholder for target resolution and permission checks ---
	// Example: Get tokens for specified emails or groups
	targetTokens := []string{} // Replace with actual token fetching logic
	if len(reqPayload.UserEmails) > 0 {
		// Fetch tokens for these emails...
	}
	if len(reqPayload.Groups) > 0 {
		// Fetch tokens for users in these groups...
	}
	// --- End Placeholder ---

	if len(targetTokens) == 0 {
		utils.LogMessage(utils.LevelWarn, "No target tokens found for notification")
		utils.LogFooter()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "No recipients found for the notification"})
	}

	pushPayload := models.NotificationPayload{
		UserEmails: targetTokens, // Use the fetched tokens here
		Title:      reqPayload.Title,
		Message:    reqPayload.Message,
		Sound:      "default", // Or from request
		ChannelID:  "default", // Or from request
		Badge:      reqPayload.Badge,
		Data:       reqPayload.Data,
	}

	err := ns.SendPushNotification(pushPayload)
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to send push notification")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to send notification"})
	}

	utils.LogMessage(utils.LevelInfo, "Notification sent successfully")
	utils.LogFooter()
	return c.SendStatus(fiber.StatusOK)
}

// SendDailyMenuNotification sends the daily menu notification to subscribed users.
func (ns *NotificationService) SendDailyMenuNotification(menuData *models.MenuData) error {
	utils.LogHeader("🍽️ Sending Daily Menu Notification")

	subscribers, err := ns.GetSubscribedUsers("Restaurant")
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to get restaurant subscribers")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return err
	}

	if len(subscribers) == 0 {
		utils.LogMessage(utils.LevelInfo, "No users subscribed to restaurant notifications")
		utils.LogFooter()
		return nil
	}

	utils.LogLineKeyValue(utils.LevelInfo, "Subscribers found", len(subscribers))

	// Construct the notification message from menuData
	var menuSummary strings.Builder
	menuSummary.WriteString("Au menu aujourd'hui:\n")
	if len(menuData.GrilladesMidi) > 0 {
		menuSummary.WriteString(fmt.Sprintf("- Grill/Trad: %s\n", strings.Join(menuData.GrilladesMidi, ", ")))
	}
	if len(menuData.Migrateurs) > 0 {
		menuSummary.WriteString(fmt.Sprintf("- Migrateurs: %s\n", strings.Join(menuData.Migrateurs, ", ")))
	}
	if len(menuData.Cibo) > 0 {
		menuSummary.WriteString(fmt.Sprintf("- Végé: %s\n", strings.Join(menuData.Cibo, ", ")))
	}
	// Add more sections as needed

	// Collect tokens
	tokens := []string{}
	for _, sub := range subscribers {
		if sub.NotificationToken != "" {
			tokens = append(tokens, sub.NotificationToken)
		}
	}

	if len(tokens) == 0 {
		utils.LogMessage(utils.LevelInfo, "No valid push tokens found for subscribers")
		utils.LogFooter()
		return nil
	}

	// Prepare Expo payload
	payload := models.NotificationPayload{
		UserEmails: tokens, // Expo uses 'to' field for tokens
		Title:      "Menu du Restaurant",
		Message:    menuSummary.String(),
		Sound:      "default",
		ChannelID:  "default", // Ensure this channel exists on the client app
		// Add other fields like Badge, Data if needed
	}

	// Send via Expo
	err = ns.SendPushNotification(payload)
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to send push notifications")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return err
	}

	utils.LogMessage(utils.LevelInfo, "Daily menu notifications sent successfully")
	utils.LogFooter()
	return nil
}
