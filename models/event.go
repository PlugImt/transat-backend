package models

import "time"

// Event represents an event entity
type Event struct {
	ID           int        `json:"id_events" db:"id_events"`
	Name         string     `json:"name" db:"name"`
	Description  string     `json:"description" db:"description"`
	Link         string     `json:"link" db:"link"`
	StartDate    time.Time  `json:"start_date" db:"start_date"`
	EndDate      *time.Time `json:"end_date,omitempty" db:"end_date"`
	Location     string     `json:"location" db:"location"`
	CreationDate time.Time  `json:"creation_date" db:"creation_date"`
	Picture      string     `json:"picture" db:"picture"`
	Creator      string     `json:"creator" db:"creator"`
	ClubID       int        `json:"id_club" db:"id_club"`
}

// EventAttendee represents an event attendee relationship
type EventAttendee struct {
	Email   string `json:"email" db:"email"`
	EventID int    `json:"id_events" db:"id_events"`
}

// EventWithDetails represents an event with additional details
type EventWithDetails struct {
	Event
	AttendeeCount   int  `json:"attendee_count"`
	IsUserAttending bool `json:"is_user_attending,omitempty"`
}

// CreateEventRequest represents the request body for creating an event
type CreateEventRequest struct {
	Name        string     `json:"name" validate:"required,max=100"`
	Description string     `json:"description" validate:"max=200"`
	Link        string     `json:"link" validate:"max=100"`
	StartDate   time.Time  `json:"start_date" validate:"required"`
	EndDate     *time.Time `json:"end_date,omitempty"`
	Location    string     `json:"location" validate:"required,max=100"`
	Picture     string     `json:"picture" validate:"max=500"`
	ClubID      int        `json:"id_club" validate:"required"`
}

// UpdateEventRequest represents the request body for updating an event
type UpdateEventRequest struct {
	Name        string     `json:"name,omitempty" validate:"max=100"`
	Description string     `json:"description,omitempty" validate:"max=200"`
	Link        string     `json:"link,omitempty" validate:"max=100"`
	StartDate   *time.Time `json:"start_date,omitempty"`
	EndDate     *time.Time `json:"end_date,omitempty"`
	Location    string     `json:"location,omitempty" validate:"max=100"`
	Picture     string     `json:"picture,omitempty" validate:"max=500"`
}
