package models

import "encoding/json"

// InteSchedule stores the raw ITE calendar ICS and its parsed JSON representation.
type InteSchedule struct {
	RawICS  string          `json:"raw_ics,omitempty" db:"raw_ics"`
	IcsJSON json.RawMessage `json:"ics_json,omitempty" db:"ics_json"`
}
