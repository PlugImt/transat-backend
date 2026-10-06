package services

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	appI18n "github.com/plugimt/transat-backend/i18n"
	"github.com/plugimt/transat-backend/models"

	"github.com/nicksnyder/go-i18n/v2/i18n"
)

const (
	expoPushURL   = "https://api.expo.dev/v2/push/send"
	expoBatchSize = 100 // Expo's per-request limit
	defaultLang   = "fr"
)

// NotificationService resolves who gets a push notification and sends it through Expo.
type NotificationService struct {
	db          *sql.DB
	expoPushURL string
	client      *http.Client
}

// NewNotificationService creates a new NotificationService.
func NewNotificationService(db *sql.DB) *NotificationService {
	return &NotificationService{
		db:          db,
		expoPushURL: expoPushURL,
		client:      &http.Client{Timeout: 15 * time.Second},
	}
}

// GetSubscribedUsersWithLanguage retrieves users subscribed to a category with their language.
func (ns *NotificationService) GetSubscribedUsersWithLanguage(category models.NotificationCategory) ([]models.NotificationTargetWithLanguage, error) {
	return ns.queryTargets(`
		SELECT n.email, unt.token, l.code
		FROM notifications n
		JOIN user_notification_tokens unt ON n.email = unt.email
		JOIN services s ON n.id_services = s.id_services
		JOIN newf ON n.email = newf.email
		JOIN languages l ON newf.language = l.id_languages
		WHERE s.name = $1
	`, string(category))
}

// GetUsersInterestedInEventOrClubWithLanguage retrieves users attending an event or in its club
// who are subscribed to the category, with their language.
func (ns *NotificationService) GetUsersInterestedInEventOrClubWithLanguage(eventID int, clubID int, category models.NotificationCategory) ([]models.NotificationTargetWithLanguage, error) {
	return ns.queryTargets(`
		SELECT DISTINCT unt.email, unt.token, COALESCE(l.code, 'fr')
		FROM user_notification_tokens unt
		LEFT JOIN newf n ON unt.email = n.email
		LEFT JOIN languages l ON n.language = l.id_languages
		WHERE unt.token IS NOT NULL AND unt.token != ''
		AND (
			unt.email IN (SELECT email FROM events_attendents WHERE id_events = $1)
			OR unt.email IN (SELECT email FROM clubs_members WHERE id_clubs = $2)
		)
		AND unt.email IN (
			SELECT n2.email FROM notifications n2
			JOIN services s ON s.id_services = n2.id_services
			WHERE s.name = $3
		)
	`, eventID, clubID, string(category))
}

func (ns *NotificationService) queryTargets(query string, args ...interface{}) ([]models.NotificationTargetWithLanguage, error) {
	rows, err := ns.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query notification targets: %w", err)
	}
	defer rows.Close()

	var targets []models.NotificationTargetWithLanguage
	for rows.Next() {
		var target models.NotificationTargetWithLanguage
		if err := rows.Scan(&target.Email, &target.NotificationToken, &target.LanguageCode); err != nil {
			return nil, fmt.Errorf("failed to scan notification target: %w", err)
		}
		targets = append(targets, target)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to read notification targets: %w", err)
	}
	return targets, nil
}

// SendLocalized sends one notification per language, composing title and message in each,
// and returns how many devices were addressed.
func (ns *NotificationService) SendLocalized(
	targets []models.NotificationTargetWithLanguage,
	compose func(localizer *i18n.Localizer) (title, message string),
	navigation *models.NavigationTarget,
) int {
	tokensByLang := make(map[string][]string)
	for _, target := range targets {
		if target.NotificationToken == "" {
			continue
		}
		lang := target.LanguageCode
		if lang == "" {
			lang = defaultLang
		}
		tokensByLang[lang] = append(tokensByLang[lang], target.NotificationToken)
	}

	sent := 0
	for lang, tokens := range tokensByLang {
		title, message := compose(appI18n.GetLocalizer(lang))
		err := ns.SendPushNotification(models.NotificationPayload{
			NotificationTokens: tokens,
			Title:              title,
			Message:            message,
			Navigation:         navigation,
		})
		if err != nil {
			log.Printf("Failed to send notification to %s users: %v", lang, err)
			continue
		}
		sent += len(tokens)
	}
	return sent
}

type expoTicket struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// SendPushNotification sends a push notification to the payload's tokens through Expo, in batches.
func (ns *NotificationService) SendPushNotification(payload models.NotificationPayload) error {
	tokens := dedupe(payload.NotificationTokens)
	if len(tokens) == 0 {
		return fmt.Errorf("no valid notification tokens found")
	}

	data := payload.ExpoData()
	failed := 0
	var lastErr error

	for start := 0; start < len(tokens); start += expoBatchSize {
		batch := tokens[start:min(start+expoBatchSize, len(tokens))]

		messages := make([]map[string]interface{}, len(batch))
		for i, token := range batch {
			message := map[string]interface{}{
				"to":        token,
				"title":     payload.Title,
				"sound":     "default",
				"channelId": "default",
			}
			if payload.Message != "" {
				message["body"] = payload.Message
			}
			if data != nil {
				message["data"] = data
			}
			messages[i] = message
		}

		batchFailed, err := ns.sendBatch(messages)
		if err != nil {
			lastErr = err
		}
		failed += batchFailed
	}

	if failed > 0 {
		return fmt.Errorf("failed to send %d/%d notifications: %w", failed, len(tokens), lastErr)
	}
	return nil
}

// sendBatch posts one Expo request and returns how many of its messages failed.
func (ns *NotificationService) sendBatch(messages []map[string]interface{}) (int, error) {
	body, err := json.Marshal(messages)
	if err != nil {
		return len(messages), fmt.Errorf("failed to encode expo messages: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, ns.expoPushURL, bytes.NewReader(body))
	if err != nil {
		return len(messages), fmt.Errorf("failed to build expo request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := ns.client.Do(req)
	if err != nil {
		return len(messages), fmt.Errorf("expo request failed: %w", err)
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return len(messages), fmt.Errorf("expo push failed with status %d: %s", resp.StatusCode, raw)
	}

	var parsed struct {
		Data []expoTicket `json:"data"`
	}
	if err := json.Unmarshal(raw, &parsed); err != nil {
		return 0, nil // Accepted by Expo; the receipt format is not worth failing the send for.
	}

	failed := 0
	var lastErr error
	for _, ticket := range parsed.Data {
		if ticket.Status != "ok" {
			failed++
			lastErr = fmt.Errorf("expo returned error: %s", ticket.Message)
		}
	}
	return failed, lastErr
}

func dedupe(values []string) []string {
	seen := make(map[string]struct{}, len(values))
	unique := make([]string, 0, len(values))
	for _, v := range values {
		if v == "" {
			continue
		}
		if _, dup := seen[v]; dup {
			continue
		}
		seen[v] = struct{}{}
		unique = append(unique, v)
	}
	return unique
}
