package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	appI18n "github.com/plugimt/transat-backend/i18n"
	"github.com/plugimt/transat-backend/models"
	"github.com/plugimt/transat-backend/services/push"
)

// SetPipeline switches the service from sending straight to Expo to queueing in the delivery pipeline, which
// delivers from its own worker. It is called once at startup.
func (ns *NotificationService) SetPipeline(p *push.Pipeline) { ns.pipeline = p }

// Pipeline returns the delivery pipeline, nil when the service still sends directly.
func (ns *NotificationService) Pipeline() *push.Pipeline { return ns.pipeline }

// OutboxEnabled reports whether notifications are queued in the delivery pipeline.
func (ns *NotificationService) OutboxEnabled() bool { return ns.pipeline != nil }

// Announcement is a notification written once and localized for every language.
type Announcement struct {
	// Kind labels the notification in the delivery log.
	Kind string
	// Category is the preference that applies: users who switched it off are left out.
	Category models.NotificationCategory
	// DedupeKey makes queueing idempotent; the same key is never queued twice.
	DedupeKey string
	// Audience defaults to every user of the category.
	Audience push.Audience
	// Compose writes the notification in one language; it is called once per language at queueing time.
	Compose    func(lang string, localizer *i18n.Localizer) (title, message string)
	Navigation *models.NavigationTarget
	// ImageURL is an optional https image shown with the notification; empty for text only.
	ImageURL string
	// ActionCategory is the interactive category (action buttons) the app shows with the notification; empty for none.
	ActionCategory string
	// ExpiresAt is the deadline after which the notification is no longer worth sending; zero uses the default.
	ExpiresAt time.Time
}

// Announce queues an announcement and returns at once. The recipients are not chosen here: they are worked out
// by the pipeline's worker, which also checks every device again right before sending.
func (ns *NotificationService) Announce(ctx context.Context, a Announcement) (push.EnqueueResult, error) {
	if ns.pipeline == nil {
		return push.EnqueueResult{}, errors.New("the delivery pipeline is not enabled")
	}
	languages, err := ns.languageCodes(ctx)
	if err != nil {
		return push.EnqueueResult{}, err
	}
	content := push.Content{}
	for _, lang := range languages {
		title, body := a.Compose(lang, appI18n.GetLocalizer(lang))
		content[lang] = push.Text{Title: title, Body: body}
	}
	if _, ok := content[push.DefaultLanguage]; !ok {
		title, body := a.Compose(push.DefaultLanguage, appI18n.GetLocalizer(push.DefaultLanguage))
		content[push.DefaultLanguage] = push.Text{Title: title, Body: body}
	}
	audience := a.Audience
	if audience.Type == "" {
		audience.Type = models.AudienceAll
	}
	return ns.pipeline.Enqueue(ctx, push.Message{
		Kind: a.Kind, Category: a.Category, DedupeKey: a.DedupeKey, Audience: audience,
		Content: content, Navigation: a.Navigation, ImageURL: a.ImageURL, ActionCategory: a.ActionCategory, ExpiresAt: a.ExpiresAt,
	})
}

// CanReach reports whether a notification of the category sent to this user would go to at least one device now:
// the user did not switch the category off and has a device with a valid session, a token and a permission.
func (ns *NotificationService) CanReach(ctx context.Context, email string, category models.NotificationCategory) (bool, error) {
	if ns.pipeline == nil {
		return false, nil
	}
	res, err := ns.pipeline.Resolve(ctx, push.Audience{
		NotificationAudience: models.NotificationAudience{Type: models.AudienceUsers, Emails: []string{email}},
	}, category)
	if err != nil {
		return false, err
	}
	return len(res.Targets) > 0, nil
}

func (ns *NotificationService) languageCodes(ctx context.Context) ([]string, error) {
	rows, err := ns.db.QueryContext(ctx, `SELECT code FROM languages ORDER BY code`)
	if err != nil {
		return nil, fmt.Errorf("failed to list languages: %w", err)
	}
	defer rows.Close()
	var codes []string
	for rows.Next() {
		var code string
		if err := rows.Scan(&code); err != nil {
			return nil, err
		}
		codes = append(codes, code)
	}
	return codes, rows.Err()
}

// queuePhotoAlbum is NotifyNewPhotoAlbum for the delivery pipeline. Queueing is the durable step: once it
// succeeds the announcement is final and the ledger of the photo service can close it. The counts describe the
// devices that qualify now; the pipeline checks them again when it sends, and Accepted and Failed stay zero
// because nothing has been sent yet.
func (ns *NotificationService) queuePhotoAlbum(albumID int64, title string, eventDate *time.Time) (PhotoAlbumDelivery, error) {
	ctx := context.Background()
	if _, err := ns.Announce(ctx, Announcement{
		Kind:      "photo_album",
		Category:  models.NotificationPhotos,
		DedupeKey: fmt.Sprintf("photo_album:%d", albumID),
		Compose: func(lang string, l *i18n.Localizer) (string, string) {
			return composePhotoAlbum(l, lang, title, eventDate)
		},
		Navigation: models.PhotoAlbumNavigation(albumID, truncateRunes(title, 120)),
	}); err != nil {
		return PhotoAlbumDelivery{}, err
	}
	res, err := ns.pipeline.Resolve(ctx, push.Audience{NotificationAudience: models.NotificationAudience{Type: models.AudienceAll}}, models.NotificationPhotos)
	if err != nil {
		// The announcement is queued; only its counts are missing.
		return PhotoAlbumDelivery{}, nil
	}
	return PhotoAlbumDelivery{Users: res.Users, Devices: len(res.Targets)}, nil
}
