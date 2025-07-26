package models

// ReservationUser is used in multiple places for user info
type ReservationUser struct {
	Email          string `json:"email"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	ProfilePicture string `json:"profile_picture,omitempty"`
}

// ReservationCategory represents a reservation category
type ReservationCategory struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// ReservationItem represents an item that can be reserved
type ReservationItem struct {
	ID   int              `json:"id"`
	Name string           `json:"name"`
	Slot bool             `json:"slot"`
	User *ReservationUser `json:"user,omitempty"` // Optional: only shown when item is reserved
}

// ReservationOverviewResponse represents the response for /reservation
type ReservationOverviewResponse struct {
	Categories []ReservationCategory `json:"categories"`
	Items      []ReservationItem     `json:"items"`
}

// ReservationCreateCategoryRequest is the body for POST /reservation/category/
type ReservationCreateCategoryRequest struct {
	Name             string `json:"name"`
	IDClubParent     *int   `json:"id_club_parent,omitempty"`     // Optional if category is defined
	IDCategoryParent *int   `json:"id_category_parent,omitempty"` // Optional parent category
}

// ReservationCategoryComplete represents a complete reservation category with items
type ReservationCategoryComplete struct {
	ID int `json:"id"`
	ReservationCreateCategoryRequest
}

// ReservationSlotDetail represents one reservation slot for an item
type ReservationSlotDetail struct {
	ID        int             `json:"id"`
	User      ReservationUser `json:"user"`
	StartDate string          `json:"start_date"`
	EndDate   string          `json:"end_date"`
}

// ReservationItemDetailResponse represents the response for /reservation/items/{id}
type ReservationItemDetailResponse struct {
	ID                int                     `json:"id"`
	Name              string                  `json:"name"`
	Slot              bool                    `json:"slot"`
	Reservation       []ReservationSlotDetail `json:"reservation"`
	ReservationBefore []ReservationSlotDetail `json:"reservation_before"`
	ReservationAfter  []ReservationSlotDetail `json:"reservation_after"`
}

// ReservationStartRequest is used to start a reservation
type ReservationStartRequest struct {
	StartDate string `json:"start_date"`
}

// ReservationEndRequest is used to end a reservation (non-slot-based)
type ReservationEndRequest struct {
	EndDate string `json:"end_date"`
}

// ReservationCreateItemRequest is the body for POST /reservation/item/
type ReservationCreateItemRequest struct {
	Name           string `json:"name"`
	Slot           bool   `json:"slot"`
	Description    string `json:"description,omitempty"`
	Location       string `json:"location,omitempty"`
	IDClub         *int   `json:"id_club,omitempty"`
	CategoryParent *int   `json:"category_parent,omitempty"`
}
