package models

type PlanningSport struct {
	ID        int    `json:"id"`
	Day string `json:"day_of_week"`
	Activity  string `json:"activity"`
	Place      *string `json:"place,omitempty"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}