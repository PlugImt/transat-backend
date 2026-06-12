package models

import (
	"encoding/json"
	"time"
)

// UserSchedule stores a Newf's ICS calendar subscription and synced data.
type UserSchedule struct {
	UserID       int             `json:"user_id" db:"user_id"`
	IcsURL       string          `json:"ics_url,omitempty" db:"ics_url"`
	LastSyncAt   *time.Time      `json:"last_sync_at,omitempty" db:"last_sync_at"`
	CalendarData json.RawMessage `json:"calendar_data,omitempty" db:"calendar_data"`
}

// CalendarEvent is a single parsed ICS event entry.
type CalendarEvent struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
	Location  string `json:"location"`
}

// CalendarData is date-keyed events: {"2025-09-03": [...]}.
type CalendarData map[string][]CalendarEvent

// UpdateUserScheduleRequest is the body for PATCH /schedule/me.
type UpdateUserScheduleRequest struct {
	IcsURL string `json:"ics_url"`
}
