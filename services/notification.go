package services

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/lib/pq"
	appI18n "github.com/plugimt/transat-backend/i18n"
	"github.com/plugimt/transat-backend/models"
	"github.com/plugimt/transat-backend/services/push"

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
	// A request that fails on the network, with 429 or with 5xx is retried this many times in all.
	maxAttempts int
	retryDelay  time.Duration
	// pipeline, when set, replaces direct sending by queueing; see SetPipeline.
	pipeline *push.Pipeline
}

// NewNotificationService creates a new NotificationService.
func NewNotificationService(db *sql.DB) *NotificationService {
	return &NotificationService{
		db:          db,
		expoPushURL: expoPushURL,
		client:      &http.Client{Timeout: 15 * time.Second},
		maxAttempts: 3,
		retryDelay:  2 * time.Second,
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
	Details struct {
		Error string `json:"error"`
	} `json:"details"`
}

// PushResult counts what Expo accepted. InvalidTokens are devices Expo says no longer exist; they are removed.
type PushResult struct {
	Accepted      int
	Failed        int
	InvalidTokens []string
}

// SendPushNotification sends a push notification to the payload's tokens through Expo, in batches.
func (ns *NotificationService) SendPushNotification(payload models.NotificationPayload) error {
	tokens := dedupe(payload.NotificationTokens)
	if len(tokens) == 0 {
		return fmt.Errorf("no valid notification tokens found")
	}
	result, err := ns.SendPush(payload)
	if result.Failed > 0 {
		return fmt.Errorf("failed to send %d/%d notifications: %w", result.Failed, len(tokens), err)
	}
	return err
}

// SendPush is SendPushNotification with the outcome per device. A device that fails does not stop the others;
// tokens Expo reports as unregistered are deleted so they are not tried again.
func (ns *NotificationService) SendPush(payload models.NotificationPayload) (PushResult, error) {
	tokens := dedupe(payload.NotificationTokens)
	var result PushResult
	var lastErr error
	data := payload.ExpoData()

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
			if payload.ActionCategory != "" {
				message["categoryId"] = payload.ActionCategory
			}
			if data != nil {
				message["data"] = data
			}
			messages[i] = message
		}

		accepted, failed, invalid, err := ns.sendBatch(messages, batch)
		result.Accepted += accepted
		result.Failed += failed
		result.InvalidTokens = append(result.InvalidTokens, invalid...)
		if err != nil {
			lastErr = err
		}
	}

	ns.pruneTokens(result.InvalidTokens)
	return result, lastErr
}

func (ns *NotificationService) pruneTokens(tokens []string) {
	if len(tokens) == 0 {
		return
	}
	res, err := ns.db.Exec(`DELETE FROM user_notification_tokens WHERE token = ANY($1)`, pq.Array(tokens))
	if err != nil {
		log.Printf("Failed to remove %d unregistered push tokens: %v", len(tokens), err)
		return
	}
	removed, _ := res.RowsAffected()
	log.Printf("Removed %d unregistered push token rows", removed)
}

// sendBatch posts one Expo request, retrying transient failures, and returns how many of its messages were
// accepted or failed, and which tokens are no longer valid.
func (ns *NotificationService) sendBatch(messages []map[string]interface{}, tokens []string) (accepted, failed int, invalid []string, err error) {
	body, err := json.Marshal(messages)
	if err != nil {
		return 0, len(messages), nil, fmt.Errorf("failed to encode expo messages: %w", err)
	}

	var raw []byte
	for attempt := 1; ; attempt++ {
		var retry bool
		raw, retry, err = ns.post(body)
		if err == nil {
			break
		}
		if !retry || attempt >= ns.maxAttempts {
			return 0, len(messages), nil, err
		}
		time.Sleep(time.Duration(attempt) * ns.retryDelay)
	}

	var parsed struct {
		Data []expoTicket `json:"data"`
	}
	if err := json.Unmarshal(raw, &parsed); err != nil || len(parsed.Data) != len(messages) {
		return len(messages), 0, nil, nil // Accepted by Expo; the receipt format is not worth failing the send for.
	}

	var lastErr error
	for i, ticket := range parsed.Data {
		if ticket.Status == "ok" {
			accepted++
			continue
		}
		failed++
		lastErr = fmt.Errorf("expo returned error: %s", redactTokens(ticket.Message))
		if ticket.Details.Error == "DeviceNotRegistered" {
			invalid = append(invalid, tokens[i])
		}
	}
	return accepted, failed, invalid, lastErr
}

// post sends one request. retry says whether the failure is worth another attempt.
func (ns *NotificationService) post(body []byte) (raw []byte, retry bool, err error) {
	req, err := http.NewRequest(http.MethodPost, ns.expoPushURL, bytes.NewReader(body))
	if err != nil {
		return nil, false, fmt.Errorf("failed to build expo request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := ns.client.Do(req)
	if err != nil {
		return nil, true, fmt.Errorf("expo request failed: %w", err)
	}
	defer resp.Body.Close()
	raw, _ = io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		transient := resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode >= 500
		return nil, transient, fmt.Errorf("expo push failed with status %d: %s", resp.StatusCode, redactTokens(string(raw)))
	}
	return raw, false, nil
}

var pushTokenPattern = regexp.MustCompile(`Expo(?:nent)?PushToken\[[^\]]*\]`)

// redactTokens keeps device tokens out of logs and error messages.
func redactTokens(text string) string {
	return pushTokenPattern.ReplaceAllString(text, "[push token]")
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
