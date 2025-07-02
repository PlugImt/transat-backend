package models

import (
	"time"
)

// Notification represents a record in the 'notifications' table
type Notification struct {
	Email       string `json:"email" db:"email"`
	IDServices  int    `json:"id_services" db:"id_services"`
	ServiceName string `json:"service_name,omitempty" db:"service_name"`
}

// NotificationTarget represents a user's email and notification token
type NotificationTarget struct {
	Email             string `json:"email"`
	NotificationToken string `json:"notification_token"`
}

// NotificationTargetWithLanguage represents a user's email, notification token, and language preference
type NotificationTargetWithLanguage struct {
	Email             string `json:"email"`
	NotificationToken string `json:"notification_token"`
	LanguageCode      string `json:"language_code"`
}

// NotificationPayload represents the data structure for sending push notifications
type NotificationPayload struct {
	NotificationTokens []string               `json:"notification_tokens,omitempty"`
	UserEmails         []string               `json:"user_emails,omitempty"`
	Groups             []string               `json:"groups,omitempty"`
	Title              string                 `json:"title"`
	Message            string                 `json:"message,omitempty"`
	Sound              string                 `json:"sound,omitempty"`
	ChannelID          string                 `json:"channel_id,omitempty"`
	Badge              int                    `json:"badge,omitempty"`
	Data               map[string]interface{} `json:"data,omitempty"`
	Subtitle           string                 `json:"subtitle,omitempty"`
	TTL                int                    `json:"ttl,omitempty"`
	ImageURL           string                 `json:"imageUrl,omitempty"`
}

}
