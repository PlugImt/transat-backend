package models

import "time"

// Carpool represents a carpool offer entity
type Carpool struct {
	ID             int       `json:"id_carpools" db:"id_carpools"`
	CreatorEmail   string    `json:"creator_email" db:"creator_email"`
	TripType       string    `json:"trip_type" db:"trip_type"`
	DeparturePlace string    `json:"departure_place" db:"departure_place"`
	Destination    string    `json:"destination" db:"destination"`
	DepartureTime  time.Time `json:"departure_time" db:"departure_time"`
	ContactDetails string    `json:"contact_details" db:"contact_details"`
	Description    string    `json:"description" db:"description"`
	Status         string    `json:"status" db:"status"`
	CreationDate   time.Time `json:"creation_date" db:"creation_date"`
	UpdatedDate    time.Time `json:"updated_date" db:"updated_date"`
}

// CarpoolWithCreator represents a carpool with the creator's basic info (for GET requests)
type CarpoolWithCreator struct {
	Carpool
	CreatorFirstName  string `json:"creator_first_name"`
	CreatorLastName   string `json:"creator_last_name"`
	CreatorProfilePic string `json:"creator_profile_picture"`
	CreatorGradYear   int64  `json:"creator_graduation_year,omitempty"`
}

// CreateCarpoolRequest represents the request body for creating a carpool offer
type CreateCarpoolRequest struct {
	TripType       string    `json:"trip_type" validate:"required,oneof=SHOPPING WEEKEND OTHER"`
	DeparturePlace string    `json:"departure_place" validate:"required,max=100"`
	Destination    string    `json:"destination" validate:"required,max=100"`
	DepartureTime  time.Time `json:"departure_time" validate:"required"`
	ContactDetails string    `json:"contact_details" validate:"required,max=100"`
	Description    string    `json:"description,omitempty" validate:"max=500"`
}

// UpdateCarpoolStatusRequest represents the request body for changing the status
type UpdateCarpoolStatusRequest struct {
	Status string `json:"status" validate:"required,oneof=OPEN FULL ARCHIVED"`
}
