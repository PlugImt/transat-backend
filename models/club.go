package models

// Club represents a club entity
type Club struct {
	ID          int    `json:"id_clubs" db:"id_clubs"`
	Name        string `json:"name" db:"name"`
	Picture     string `json:"picture" db:"picture"`
	Description string `json:"description" db:"description"`
	Location    string `json:"location" db:"location"`
	Link        string `json:"link" db:"link"`
}

// ClubMember represents a club member relationship
type ClubMember struct {
	Email   string `json:"email" db:"email"`
	ClubID  int    `json:"id_clubs" db:"id_clubs"`
	IsRespo bool   `json:"is_respo,omitempty"` // Derived field, not stored directly
}

// ClubWithMembers represents a club with its members count and member info
type ClubWithMembers struct {
	Club
	MemberCount  int  `json:"member_count"`
	IsUserMember bool `json:"is_user_member,omitempty"`
	IsUserRespo  bool `json:"is_user_respo,omitempty"`
}

// CreateClubRequest represents the request body for creating a club
type CreateClubRequest struct {
	Name        string `json:"name" validate:"required,max=50"`
	Picture     string `json:"picture" validate:"required,max=500"`
	Description string `json:"description" validate:"max=500"`
	Location    string `json:"location" validate:"max=100"`
	Link        string `json:"link" validate:"max=500"`
}

// UpdateClubRequest represents the request body for updating a club
type UpdateClubRequest struct {
	Name        string `json:"name,omitempty" validate:"max=50"`
	Picture     string `json:"picture,omitempty" validate:"max=500"`
	Description string `json:"description,omitempty" validate:"max=500"`
	Location    string `json:"location,omitempty" validate:"max=100"`
	Link        string `json:"link,omitempty" validate:"max=500"`
}

// AddRespoRequest represents the request body for adding a club responsible
type AddRespoRequest struct {
	Email string `json:"email" validate:"required,email"`
}
