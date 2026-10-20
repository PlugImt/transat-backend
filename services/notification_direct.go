package services

import (
	"errors"
	"fmt"
	"time"

	appI18n "github.com/plugimt/transat-backend/i18n"
	"github.com/plugimt/transat-backend/models"
)

// announceDirect is Announce without the delivery pipeline: it resolves the devices from the legacy token table
// and sends one request per language to Expo. It reports an error only when nothing could be delivered, so that the
// caller (a reminder, the waiting list) knows to try again or to pass the turn.
func (ns *NotificationService) announceDirect(a Announcement) error {
	if !a.ExpiresAt.IsZero() && time.Now().After(a.ExpiresAt) {
		return nil
	}
	targets, err := ns.resolveTargets(a.Audience.NotificationAudience, a.Category)
	if err != nil {
		return err
	}
	tokensByLang := map[string][]string{}
	for _, t := range targets {
		if t.NotificationToken == "" {
			continue
		}
		lang := t.LanguageCode
		if lang == "" {
			lang = defaultLang
		}
		tokensByLang[lang] = append(tokensByLang[lang], t.NotificationToken)
	}

	var errs []error
	accepted := 0
	for lang, tokens := range tokensByLang {
		title, body := a.Compose(lang, appI18n.GetLocalizer(lang))
		result, err := ns.SendPush(models.NotificationPayload{
			NotificationTokens: tokens,
			Title:              title,
			Message:            body,
			Kind:               a.Kind,
			Navigation:         a.Navigation,
			ActionCategory:     a.ActionCategory,
		})
		accepted += result.Accepted
		if err != nil {
			errs = append(errs, fmt.Errorf("%s: %w", lang, err))
		}
	}
	if accepted == 0 && len(errs) > 0 {
		return errors.Join(errs...)
	}
	return nil
}
