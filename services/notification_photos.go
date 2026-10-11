package services

import (
	"fmt"
	"log"
	"regexp"
	"sort"
	"strings"
	"time"

	appI18n "github.com/plugimt/transat-backend/i18n"
	"github.com/plugimt/transat-backend/models"

	"github.com/nicksnyder/go-i18n/v2/i18n"
)

var expoTokenPattern = regexp.MustCompile(`^Expo(?:nent)?PushToken\[[^\]]+\]$`)

// PhotoAlbumDelivery says how an announcement went. It holds counts only, never emails or tokens.
type PhotoAlbumDelivery struct {
	Users    int
	Devices  int
	Accepted int
	Failed   int
	// Malformed device tokens that were skipped and removed.
	Skipped int
}

// photoAlbumTargets lists the devices that may receive a new-album notification, one row per token: users who
// switched Photos on, with a registered token. A token also registered to an account that did not opt in (a
// phone shared after a logout) is left out, so a device never gets a notification its current user declined.
func (ns *NotificationService) photoAlbumTargets() ([]models.NotificationTargetWithLanguage, error) {
	return ns.queryTargets(`
		SELECT unt.email, unt.token, COALESCE(l.code, 'fr')
		FROM user_notification_tokens unt
		JOIN notifications nt ON nt.email = unt.email
		JOIN services s ON s.id_services = nt.id_services AND s.name = $1
		LEFT JOIN newf n ON n.email = unt.email
		LEFT JOIN languages l ON l.id_languages = n.language
		WHERE unt.token <> ''
		  AND NOT EXISTS (
			SELECT 1 FROM user_notification_tokens other
			WHERE other.token = unt.token
			  AND NOT EXISTS (
				SELECT 1 FROM notifications nt2
				JOIN services s2 ON s2.id_services = nt2.id_services
				WHERE nt2.email = other.email AND s2.name = $1))
		ORDER BY unt.email, unt.token
	`, string(models.NotificationPhotos))
}

// NotifyNewPhotoAlbum announces an album to the users who opted in. It returns an error only when nothing could
// be delivered at all, which is when the caller should try again; a partly delivered announcement is final.
func (ns *NotificationService) NotifyNewPhotoAlbum(albumID int64, title string, eventDate *time.Time) (PhotoAlbumDelivery, error) {
	targets, err := ns.photoAlbumTargets()
	if err != nil {
		return PhotoAlbumDelivery{}, err
	}

	var delivery PhotoAlbumDelivery
	users := map[string]bool{}
	seen := map[string]bool{}
	byLang := map[string][]string{}
	var malformed []string
	for _, target := range targets {
		users[target.Email] = true
		// A token shared by several opted-in accounts is addressed once.
		if seen[target.NotificationToken] {
			continue
		}
		seen[target.NotificationToken] = true
		if !expoTokenPattern.MatchString(target.NotificationToken) {
			malformed = append(malformed, target.NotificationToken)
			continue
		}
		lang := target.LanguageCode
		if lang == "" {
			lang = defaultLang
		}
		byLang[lang] = append(byLang[lang], target.NotificationToken)
	}
	delivery.Users = len(users)
	delivery.Skipped = len(malformed)
	ns.pruneTokens(malformed)

	langs := make([]string, 0, len(byLang))
	for lang := range byLang {
		langs = append(langs, lang)
		delivery.Devices += len(byLang[lang])
	}
	sort.Strings(langs)

	var lastErr error
	for _, lang := range langs {
		heading, message := composePhotoAlbum(appI18n.GetLocalizer(lang), lang, title, eventDate)
		result, err := ns.SendPush(models.NotificationPayload{
			NotificationTokens: byLang[lang],
			Title:              heading,
			Message:            message,
			Navigation:         models.PhotoAlbumNavigation(albumID, truncateRunes(title, 120)),
		})
		delivery.Accepted += result.Accepted
		delivery.Failed += result.Failed
		if err != nil {
			lastErr = err
		}
	}

	if delivery.Accepted == 0 && delivery.Failed > 0 {
		return delivery, fmt.Errorf("no photo album notification was accepted: %w", lastErr)
	}
	if lastErr != nil {
		log.Printf("Photo album notification partly delivered: %d accepted, %d failed", delivery.Accepted, delivery.Failed)
	}
	return delivery, nil
}

func composePhotoAlbum(localizer *i18n.Localizer, lang, title string, eventDate *time.Time) (string, string) {
	name := truncateRunes(title, 100)
	heading := appI18n.Localize(localizer, "photo_notification.new_album_title", "New photo album", nil)
	if eventDate == nil {
		return heading, appI18n.Localize(localizer, "photo_notification.new_album_message", "{{.AlbumName}}",
			map[string]interface{}{"AlbumName": name})
	}
	return heading, appI18n.Localize(localizer, "photo_notification.new_album_message_dated", "{{.AlbumName}} · {{.Date}}",
		map[string]interface{}{"AlbumName": name, "Date": formatAlbumDate(lang, *eventDate)})
}

// formatAlbumDate writes an event date the way each supported language does; it is a calendar day, so no zone.
func formatAlbumDate(lang string, d time.Time) string {
	switch lang {
	case "en":
		return d.Format("Jan 2, 2006")
	case "de":
		return d.Format("02.01.2006")
	default:
		return d.Format("02/01/2006")
	}
}

func truncateRunes(s string, max int) string {
	s = strings.TrimSpace(s)
	runes := []rune(s)
	if len(runes) <= max {
		return s
	}
	return string(runes[:max-1]) + "…"
}
