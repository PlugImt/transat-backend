package models

import "time"

// SupportRequest represents a user support request
type SupportRequest struct {
	ID        int       `json:"id_support"`
	Email     string    `json:"email"`
	Subject   string    `json:"subject"`
	Message   string    `json:"message"`
	Status    string    `json:"status"`
	ImageURL  string    `json:"image_url,omitempty"`
	Response  string    `json:"response,omitempty"`
	CreatedAt time.Time `json:"creation_date"`
}

// SupportMedia represents attached media to a support request
type SupportMedia struct {
	ID       int `json:"id_support_media"`
	Support  int `json:"id_support"`
	FileID   int `json:"id_files"`
}

// SupportRequestInput is used for creating new support requests
type SupportRequestInput struct {
	Subject  string `json:"subject" validate:"required,min=2,max=200"`
	Message  string `json:"message" validate:"required,min=5,max=5000"`
	FileID   int    `json:"file_id,omitempty"`
}

// SupportRequestResponse is the structure returned when listing support requests
type SupportRequestResponse struct {
	Requests []SupportRequest `json:"requests"`
}

// SupportRequestCreateResponse is the structure returned when creating a support request
type SupportRequestCreateResponse struct {
	Request SupportRequest `json:"request"`
} 