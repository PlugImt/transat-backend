package services

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	appI18n "github.com/plugimt/transat-backend/i18n"
	"github.com/plugimt/transat-backend/models"
	"github.com/plugimt/transat-backend/utils"
)

type EventNotificationType string

const (
	EventNotificationCreated    EventNotificationType = "CREATED"
	EventNotificationReminder1h EventNotificationType = "REMINDER_1H"
)

// EventNotificationService encapsulates business logic for event notifications
type EventNotificationService struct {
	db                  *sql.DB
	notificationService *NotificationService
}

func NewEventNotificationService(db *sql.DB, notificationService *NotificationService) *EventNotificationService {
	return &EventNotificationService{db: db, notificationService: notificationService}
}

// LogNotification creates a deduped log entry for a sent notification
func (s *EventNotificationService) LogNotification(eventID int, email string, notifType EventNotificationType) error {
	_, err := s.db.Exec(`
        INSERT INTO event_notifications_log (id_events, email, notification_type)
        VALUES ($1, $2, $3)
        ON CONFLICT (id_events, email, notification_type) DO NOTHING
    `, eventID, email, string(notifType))
	return err
}

// AlreadySent checks whether a notification of a given type was already sent to a user for an event
func (s *EventNotificationService) AlreadySent(eventID int, email string, notifType EventNotificationType) (bool, error) {
	var exists bool
	err := s.db.QueryRow(`
        SELECT EXISTS(
            SELECT 1 FROM event_notifications_log
            WHERE id_events = $1 AND email = $2 AND notification_type = $3
        )
    `, eventID, email, string(notifType)).Scan(&exists)
	return exists, err
}

// getEventCore fetches event core fields used for notifications
func (s *EventNotificationService) getEventCore(eventID int) (name string, clubID int, startDate time.Time, err error) {
	err = s.db.QueryRow(`
        SELECT name, id_club, start_date
        FROM events WHERE id_events = $1
    `, eventID).Scan(&name, &clubID, &startDate)
	return
}

// getRecipients finds user emails interested in the club or the event itself
func (s *EventNotificationService) getRecipients(eventID int, clubID int) ([]string, error) {
	// Users who are members of the club
	clubQuery := `SELECT email FROM clubs_members WHERE id_clubs = $1`
	rowsClub, err := s.db.Query(clubQuery, clubID)
	if err != nil {
		return nil, fmt.Errorf("failed to query club members: %w", err)
	}
	defer rowsClub.Close()

	emailSet := map[string]struct{}{}
	for rowsClub.Next() {
		var email string
		if err := rowsClub.Scan(&email); err == nil {
			emailSet[email] = struct{}{}
		}
	}
	if err := rowsClub.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate club member rows: %w", err)
	}

	// Users who marked interest in the event (attendees table)
	evtQuery := `SELECT email FROM events_attendents WHERE id_events = $1`
	rowsEvt, err := s.db.Query(evtQuery, eventID)
	if err != nil {
		return nil, fmt.Errorf("failed to query event interested users: %w", err)
	}
	defer rowsEvt.Close()
	for rowsEvt.Next() {
		var email string
		if err := rowsEvt.Scan(&email); err == nil {
			emailSet[email] = struct{}{}
		}
	}
	if err := rowsEvt.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate event interested rows: %w", err)
	}

	// materialize unique list
	recipients := make([]string, 0, len(emailSet))
	for e := range emailSet {
		recipients = append(recipients, e)
	}
	return recipients, nil
}

// SendEventCreated sends creation notifications to all relevant recipients, deduped and localized
func (s *EventNotificationService) SendEventCreated(eventID int) error {
	name, clubID, startDate, err := s.getEventCore(eventID)
	if err != nil {
		return err
	}
	recipients, err := s.getRecipients(eventID, clubID)
	if err != nil {
		return err
	}
	if len(recipients) == 0 {
		return nil
	}

	// fetch tokens+language
	targets, err := s.notificationService.GetUsersWithLanguageByEmails(recipients)
	if err != nil {
		return err
	}
	if len(targets) == 0 {
		return nil
	}

	// group by language and filter out already-sent
	languageToTokens := map[string][]string{}
	languageToEmails := map[string][]string{}
	for _, t := range targets {
		sent, err := s.AlreadySent(eventID, t.Email, EventNotificationCreated)
		if err != nil {
			continue
		}
		if sent || t.NotificationToken == "" {
			continue
		}
		languageToTokens[t.LanguageCode] = append(languageToTokens[t.LanguageCode], t.NotificationToken)
		languageToEmails[t.LanguageCode] = append(languageToEmails[t.LanguageCode], t.Email)
	}

	// send per language
	total := 0
	for lang, tokens := range languageToTokens {
		if len(tokens) == 0 {
			continue
		}
		localizer := appI18n.GetLocalizer(lang)

		title := localizer.MustLocalize(&i18n.LocalizeConfig{
			MessageID:      "event.created.title",
			DefaultMessage: &i18n.Message{ID: "event.created.title", Other: "New event created"},
		})

		// include event id and open screen
		message := localizer.MustLocalize(&i18n.LocalizeConfig{
			MessageID: "event.created.message",
			TemplateData: map[string]interface{}{
				"Name":      name,
				"StartTime": utils.FormatParis(startDate, "02/01 15:04"),
				"EventID":   eventID,
			},
			DefaultMessage: &i18n.Message{ID: "event.created.message", Other: "Event '{{.Name}}' created for {{.StartTime}}."},
		})

		payload := models.NotificationPayload{
			NotificationTokens: tokens,
			Title:              title,
			Message:            message,
			Sound:              "default",
			ChannelID:          "default",
			Data: map[string]interface{}{
				"screen":  "Event",
				"eventId": eventID,
			},
		}

		if err := s.notificationService.SendPushNotification(payload); err != nil {
			// continue to next language group
			continue
		}

		// log all emails for which we sent
		for _, email := range languageToEmails[lang] {
			_ = s.LogNotification(eventID, email, EventNotificationCreated)
		}
		total += len(tokens)
	}

	_ = total // kept for potential metrics/logging
	return nil
}

// SendDueOneHourReminders scans for events starting within the next hour and sends reminders if not already sent
func (s *EventNotificationService) SendDueOneHourReminders() error {
	now := utils.Now()
	lower := now.Add(1 * time.Hour)
	upper := lower.Add(5 * time.Minute)

	rows, err := s.db.Query(`
        SELECT id_events FROM events
        WHERE start_date >= $1 AND start_date < $2
    `, lower, upper)
	if err != nil {
		return fmt.Errorf("failed to query upcoming events: %w", err)
	}
	defer rows.Close()

	var eventIDs []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err == nil {
			eventIDs = append(eventIDs, id)
		}
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("failed to iterate events: %w", err)
	}

	for _, eventID := range eventIDs {
		if err := s.sendOneHourReminderForEvent(eventID); err != nil {
			// continue other events
			continue
		}
	}
	return nil
}

func (s *EventNotificationService) sendOneHourReminderForEvent(eventID int) error {
	name, clubID, startDate, err := s.getEventCore(eventID)
	if err != nil {
		return err
	}

	// If event was created less than 1 hour before start, skip reminder per requirements
	// Note: creation_date is stored on events; use it if needed. We'll guard at send time too.
	var creationDate time.Time
	_ = s.db.QueryRow(`SELECT creation_date FROM events WHERE id_events = $1`, eventID).Scan(&creationDate)
	if !creationDate.IsZero() && startDate.Sub(creationDate) < time.Hour {
		return nil
	}

	recipients, err := s.getRecipients(eventID, clubID)
	if err != nil {
		return err
	}
	if len(recipients) == 0 {
		return nil
	}

	targets, err := s.notificationService.GetUsersWithLanguageByEmails(recipients)
	if err != nil {
		return err
	}
	if len(targets) == 0 {
		return nil
	}

	languageToTokens := map[string][]string{}
	languageToEmails := map[string][]string{}
	for _, t := range targets {
		sent, err := s.AlreadySent(eventID, t.Email, EventNotificationReminder1h)
		if err != nil {
			continue
		}
		if sent || t.NotificationToken == "" {
			continue
		}
		languageToTokens[t.LanguageCode] = append(languageToTokens[t.LanguageCode], t.NotificationToken)
		languageToEmails[t.LanguageCode] = append(languageToEmails[t.LanguageCode], t.Email)
	}

	for lang, tokens := range languageToTokens {
		if len(tokens) == 0 {
			continue
		}
		localizer := appI18n.GetLocalizer(lang)
		title := localizer.MustLocalize(&i18n.LocalizeConfig{
			MessageID:      "event.reminder.title",
			DefaultMessage: &i18n.Message{ID: "event.reminder.title", Other: "Event starts in 1 hour"},
		})
		message := localizer.MustLocalize(&i18n.LocalizeConfig{
			MessageID: "event.reminder.message",
			TemplateData: map[string]interface{}{
				"Name":    name,
				"EventID": eventID,
			},
			DefaultMessage: &i18n.Message{ID: "event.reminder.message", Other: "'{{.Name}}' starts in 1 hour."},
		})

		payload := models.NotificationPayload{
			NotificationTokens: tokens,
			Title:              title,
			Message:            message,
			Sound:              "default",
			ChannelID:          "default",
			Data: map[string]interface{}{
				"screen":  "Event",
				"eventId": eventID,
			},
		}
		if err := s.notificationService.SendPushNotification(payload); err != nil {
			continue
		}
		for _, email := range languageToEmails[lang] {
			_ = s.LogNotification(eventID, email, EventNotificationReminder1h)
		}
	}
	return nil
}
