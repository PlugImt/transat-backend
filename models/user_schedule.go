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

// UpdateUserScheduleRequest is the body for PATCH /schedule/me.
type UpdateUserScheduleRequest struct {
	IcsURL string `json:"ics_url"`
}
