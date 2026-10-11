package models

// NotificationCategory is a kind of push notification a user can switch on or off.
// To add one: declare it here, list it in NotificationCategories, seed it in a migration
// (the `services` table) and add it to the app's notification settings.
type NotificationCategory string

const (
	NotificationRestaurant     NotificationCategory = "RESTAURANT"
	NotificationTraq           NotificationCategory = "TRAQ"
	NotificationEvents         NotificationCategory = "EVENTS"
	NotificationEventReminders NotificationCategory = "EVENT_REMINDERS"
	NotificationReservations   NotificationCategory = "RESERVATIONS"
	NotificationPhotos         NotificationCategory = "PHOTOS"
)

// NotificationCategories lists every category, in display order.
var NotificationCategories = []NotificationCategory{
	NotificationRestaurant,
	NotificationTraq,
	NotificationEvents,
	NotificationEventReminders,
	NotificationReservations,
	NotificationPhotos,
}

// OptInNotificationCategories are off until the user switches them on. They are never granted by a migration
// or by registration, which subscribes every other category by default.
var OptInNotificationCategories = []NotificationCategory{NotificationPhotos}

func (c NotificationCategory) IsOptIn() bool {
	for _, optIn := range OptInNotificationCategories {
		if c == optIn {
			return true
		}
	}
	return false
}

func (c NotificationCategory) IsValid() bool {
	for _, known := range NotificationCategories {
		if c == known {
			return true
		}
	}
	return false
}

// NotificationPreference is one category and whether the user receives it.
type NotificationPreference struct {
	Service NotificationCategory `json:"service"`
	Enabled bool                 `json:"enabled"`
}
