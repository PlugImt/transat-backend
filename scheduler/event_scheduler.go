package scheduler

import (
	"database/sql"
	"log"
	"time"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	appI18n "github.com/plugimt/transat-backend/i18n"
	"github.com/plugimt/transat-backend/models"
	"github.com/plugimt/transat-backend/services"
)

// EventScheduler handles scheduled tasks for events
type EventScheduler struct {
	db                 *sql.DB
	notificationService *services.NotificationService
	stopChan           chan bool
}

// NewEventScheduler creates a new event scheduler
func NewEventScheduler(db *sql.DB, notificationService *services.NotificationService) *EventScheduler {
	return &EventScheduler{
		db:                 db,
		notificationService: notificationService,
		stopChan:           make(chan bool),
	}
}

// Start starts the event scheduler
func (es *EventScheduler) Start() {
	go es.run()
	log.Println("✅ Event scheduler started")
}

// Stop stops the event scheduler
func (es *EventScheduler) Stop() {
	es.stopChan <- true
	log.Println("🛑 Event scheduler stopped")
}

// run executes the scheduler loop
func (es *EventScheduler) run() {
	ticker := time.NewTicker(1 * time.Minute) // Check every minute
	defer ticker.Stop()

	for {
		select {
		case <-es.stopChan:
			return
		case <-ticker.C:
			es.checkAndSendEventNotifications()
		}
	}
}

// checkAndSendEventNotifications checks for events starting in 1 hour and sends notifications
func (es *EventScheduler) checkAndSendEventNotifications() {
	// Get current time + 1 hour
	oneHourFromNow := time.Now().Add(1 * time.Hour)
	// Use a 2-minute window to account for timing variations
	windowStart := oneHourFromNow.Add(-1 * time.Minute)
	windowEnd := oneHourFromNow.Add(1 * time.Minute)

	query := `
		SELECT 
			e.id_events,
			e.name,
			e.start_date,
			e.location,
			e.id_club,
			COALESCE(c.name, '') as club_name
		FROM events e
		LEFT JOIN clubs c ON e.id_club = c.id_clubs
		WHERE e.start_date >= $1 
		AND e.start_date <= $2
		AND e.start_date > NOW()
		AND NOT EXISTS (
			SELECT 1 FROM event_notification_sent ens 
			WHERE ens.id_events = e.id_events 
			AND ens.notification_type = '1h_before'
		)
	`

	rows, err := es.db.Query(query, windowStart, windowEnd)
	if err != nil {
		log.Printf("Error querying events for 1h notification: %v", err)
		return
	}
	defer rows.Close()

	var eventsToNotify []struct {
		ID        int
		Name      string
		StartDate time.Time
		Location  string
		ClubID    int
		ClubName  string
	}

	for rows.Next() {
		var event struct {
			ID        int
			Name      string
			StartDate time.Time
			Location  string
			ClubID    int
			ClubName  string
		}
		if err := rows.Scan(&event.ID, &event.Name, &event.StartDate, &event.Location, &event.ClubID, &event.ClubName); err != nil {
			log.Printf("Error scanning event row: %v", err)
			continue
		}
		eventsToNotify = append(eventsToNotify, event)
	}

	// Send notifications for each event
	for _, event := range eventsToNotify {
		es.sendEventReminderNotification(event.ID, event.Name, event.Location, event.ClubID)
	}
}

// sendEventReminderNotification sends a reminder notification 1h before an event
func (es *EventScheduler) sendEventReminderNotification(eventID int, eventName string, location string, clubID int) {
	// Get users interested in the event or club with their language preferences
	users, err := es.notificationService.GetUsersInterestedInEventOrClubWithLanguage(eventID, clubID)
	if err != nil {
		log.Printf("Failed to get users interested in event %d: %v", eventID, err)
		return
	}

	if len(users) == 0 {
		log.Printf("No users to notify for event %d", eventID)
		// Still mark as sent to avoid retrying
		es.markNotificationSent(eventID, "1h_before")
		return
	}

	// Group users by language
	languageGroups := make(map[string][]models.NotificationTargetWithLanguage)
	for _, user := range users {
		if user.NotificationToken != "" {
			langCode := user.LanguageCode
			if langCode == "" {
				langCode = "fr" // Default to French
			}
			languageGroups[langCode] = append(languageGroups[langCode], user)
		}
	}

	totalSent := 0
	// Send notifications to each language group
	for langCode, langUsers := range languageGroups {
		localizer := appI18n.GetLocalizer(langCode)

		title := localizer.MustLocalize(&i18n.LocalizeConfig{
			MessageID: "event_notification.reminder_title",
			DefaultMessage: &i18n.Message{
				ID:    "event_notification.reminder_title",
				Other: "Event in 1h",
			},
		})

		message := localizer.MustLocalize(&i18n.LocalizeConfig{
			MessageID: "event_notification.reminder_message",
			TemplateData: map[string]interface{}{
				"EventName": eventName,
				"Location":  location,
			},
			DefaultMessage: &i18n.Message{
				ID:    "event_notification.reminder_message",
				Other: "{{.EventName}} - {{.Location}}",
			},
		})

		var tokens []string
		for _, user := range langUsers {
			tokens = append(tokens, user.NotificationToken)
		}

		payload := models.NotificationPayload{
			NotificationTokens: tokens,
			Title:              title,
			Message:            message,
			Sound:              "default",
			ChannelID:          "default",
			Data: map[string]interface{}{
				"screen":  "Events",
				"eventId": eventID,
			},
		}

		if err := es.notificationService.SendPushNotification(payload); err != nil {
			log.Printf("Failed to send 1h reminder notification to %s users for event %d: %v", langCode, eventID, err)
			continue
		}

		log.Printf("Successfully sent 1h reminder notification for event %d to %d users in %s", eventID, len(tokens), langCode)
		totalSent += len(tokens)
	}

	log.Printf("Successfully sent 1h reminder notification for event %d to %d users across %d languages", eventID, totalSent, len(languageGroups))

	// Mark notification as sent
	es.markNotificationSent(eventID, "1h_before")
}

// markNotificationSent marks that a notification has been sent for an event
func (es *EventScheduler) markNotificationSent(eventID int, notificationType string) {
	// Create table if it doesn't exist (should be done via migration, but this is a safety check)
	createTableQuery := `
		CREATE TABLE IF NOT EXISTS event_notification_sent (
			id SERIAL PRIMARY KEY,
			id_events INTEGER NOT NULL,
			notification_type VARCHAR(50) NOT NULL,
			sent_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(id_events, notification_type),
			FOREIGN KEY (id_events) REFERENCES events(id_events) ON DELETE CASCADE
		)
	`
	_, err := es.db.Exec(createTableQuery)
	if err != nil {
		log.Printf("Error creating event_notification_sent table: %v", err)
	}

	// Insert or update the notification record
	insertQuery := `
		INSERT INTO event_notification_sent (id_events, notification_type)
		VALUES ($1, $2)
		ON CONFLICT (id_events, notification_type) DO NOTHING
	`
	_, err = es.db.Exec(insertQuery, eventID, notificationType)
	if err != nil {
		log.Printf("Error marking notification as sent for event %d: %v", eventID, err)
	}
}
