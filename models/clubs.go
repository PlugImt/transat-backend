package models

// Club represents a record in the 'clubs' table
type Club struct {
	ID          int    `json:"id_clubs" db:"id_clubs"`
	Name        string `json:"name" db:"name"`
	Picture     string `json:"picture" db:"picture"`
	Description string `json:"description,omitempty" db:"description"`
	Location    string `json:"location,omitempty" db:"location"`
	Link        string `json:"link,omitempty" db:"link"`
}

// ClubMember represents a record in the 'clubs_members' table
type ClubMember struct {
	Email    string `json:"email" db:"email"`
	ClubID   int    `json:"id_clubs" db:"id_clubs"`
	ClubName string `json:"club_name,omitempty" db:"club_name"`
}

// ClubRespo represents a club responsable (managed through user_roles table)
type ClubRespo struct {
	Email    string `json:"email" db:"email"`
	ClubID   int    `json:"id_clubs" db:"id_clubs"`
	ClubName string `json:"club_name,omitempty" db:"club_name"`
}

// ClubMemberWithDetails represents a club member with user details
type ClubMemberWithDetails struct {
	Email          string `json:"email"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	ProfilePicture string `json:"profile_picture"`
	GraduationYear int    `json:"graduation_year"`
	Campus         string `json:"campus"`
}

// ClubWithStats represents a club with additional statistics
type ClubWithStats struct {
	Club
	MemberCount int `json:"member_count"`
	RespoCount  int `json:"respo_count"`
}

// --- Request Structs ---

// CreateClubRequest represents the request to create a new club
type CreateClubRequest struct {
	Name        string `json:"name" validate:"required,min=1,max=50"`
	Picture     string `json:"picture" validate:"required,max=500"`
	Description string `json:"description,omitempty" validate:"max=500"`
	Location    string `json:"location,omitempty" validate:"max=100"`
	Link        string `json:"link,omitempty" validate:"max=500"`
}

// UpdateClubRequest represents the request to update club details
type UpdateClubRequest struct {
	Name        string `json:"name,omitempty" validate:"omitempty,min=1,max=50"`
	Picture     string `json:"picture,omitempty" validate:"omitempty,max=500"`
	Description string `json:"description,omitempty" validate:"omitempty,max=500"`
	Location    string `json:"location,omitempty" validate:"omitempty,max=100"`
	Link        string `json:"link,omitempty" validate:"omitempty,max=500"`
}

// AddRespoRequest represents the request to add a respo to a club
type AddRespoRequest struct {
	Email string `json:"email" validate:"required,email"`
}

// RemoveRespoRequest represents the request to remove a respo from a club
type RemoveRespoRequest struct {
	Email string `json:"email" validate:"required,email"`
}

// --- Response Structs ---

// ClubListResponse represents the response for listing clubs
type ClubListResponse struct {
	Clubs []ClubWithStats `json:"clubs"`
	Total int             `json:"total"`
}

// ClubDetailsResponse represents detailed information about a club
type ClubDetailsResponse struct {
	Club
	MemberCount int                     `json:"member_count"`
	RespoCount  int                     `json:"respo_count"`
	Members     []ClubMemberWithDetails `json:"members"`
	Respos      []ClubMemberWithDetails `json:"respos"`
	IsMember    bool                    `json:"is_member"`
	IsRespo     bool                    `json:"is_respo"`
}

// UserClubsResponse represents the clubs a user belongs to
type UserClubsResponse struct {
	MemberOf []Club `json:"member_of"`
	RespoOf  []Club `json:"respo_of"`
}

// ClubMembersResponse represents the members of a club
type ClubMembersResponse struct {
	ClubID  int                     `json:"club_id"`
	Members []ClubMemberWithDetails `json:"members"`
	Total   int                     `json:"total"`
}

// ClubResposResponse represents the respos of all clubs
type ClubResposResponse struct {
	Respos []struct {
		ClubMemberWithDetails
		ClubName string `json:"club_name"`
		ClubID   int    `json:"club_id"`
	} `json:"respos"`
	Total int `json:"total"`
}

// ClubOperationResponse represents the result of club operations
type ClubOperationResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	ClubID  int    `json:"club_id,omitempty"`
}
