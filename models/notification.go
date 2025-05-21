package models

type NotificationPayload struct {
	UserEmails         []string               `json:"userEmails,omitempty"`
	NotificationTokens []string               `json:"notificationTokens,omitempty"`
	Groups             []string               `json:"groups,omitempty"`
	Title              string                 `json:"title"`
	Message            string                 `json:"body,omitempty"`
	ImageURL           string                 `json:"imageUrl,omitempty"`
	TTL                int                    `json:"ttl,omitempty"`
	Subtitle           string                 `json:"subtitle,omitempty"`
	Sound              string                 `json:"sound,omitempty"`
	ChannelID          string                 `json:"channelId,omitempty"`
	Badge              int                    `json:"badge,omitempty"`
	Data               map[string]interface{} `json:"data,omitempty"`
}

type NotificationTarget struct {
	Email             string
	NotificationToken string
}
