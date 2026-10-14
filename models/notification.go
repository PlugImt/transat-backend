package models

// NotificationTargetWithLanguage is a device token with its owner and preferred language.
type NotificationTargetWithLanguage struct {
	Email             string
	NotificationToken string
	LanguageCode      string
}

// NotificationPayload is a push notification ready to be sent to a list of device tokens.
type NotificationPayload struct {
	NotificationTokens []string
	Title              string
	Message            string
	// Kind says what produced the notification; it is sent to the app as data.type.
	Kind string
	// Navigation is the optional screen opened when the notification is tapped.
	Navigation *NavigationTarget
}
