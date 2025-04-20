package models

import "time"

// RequestStatistic represents a single request to the API
type RequestStatistic struct {
	ID              int       `json:"id"`
	UserEmail       string    `json:"user_email"`
	Endpoint        string    `json:"endpoint"`
	Method          string    `json:"method"`
	RequestReceived time.Time `json:"request_received"`
	ResponseSent    time.Time `json:"response_sent"`
	StatusCode      int       `json:"status_code"`
	DurationMs      int       `json:"duration_ms"`
} 