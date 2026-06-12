package models

// Association represents an association entity
type Association struct {
	ID          int    `json:"id_associations" db:"id_associations"`
	Name        string `json:"name" db:"name"`
	Picture     string `json:"picture" db:"picture"`
	Description string `json:"description" db:"description"`
	Location    string `json:"location" db:"location"`
	Link        string `json:"link" db:"link"`
}

// AssociationMember represents a association member relationship
type AssociationMember struct {
	Email         string `json:"email" db:"email"`
	AssociationID int    `json:"id_associations" db:"id_associations"`
	IsRespo       bool   `json:"is_respo,omitempty"` // Derived field, not stored directly
}

// AssociationWithMembers represents a association with its members count and member info
type AssociationWithMembers struct {
	Association
	MemberCount  int  `json:"member_count"`
	IsUserMember bool `json:"is_user_member,omitempty"`
	IsUserRespo  bool `json:"is_user_respo,omitempty"`
}

// CreateAssociationRequest represents the request body for creating a association
type CreateAssociationRequest struct {
	Name        string `json:"name" validate:"required,max=50"`
	Picture     string `json:"picture" validate:"required,max=500"`
	Description string `json:"description" validate:"max=500"`
	Location    string `json:"location" validate:"max=100"`
	Link        string `json:"link" validate:"max=500"`
}

// UpdateAssociationRequest represents the request body for updating a association
type UpdateAssociationRequest struct {
	Name        string `json:"name,omitempty" validate:"max=50"`
	Picture     string `json:"picture,omitempty" validate:"max=500"`
	Description string `json:"description,omitempty" validate:"max=500"`
	Location    string `json:"location,omitempty" validate:"max=100"`
	Link        string `json:"link,omitempty" validate:"max=500"`
}

// AddRespoAssociationRequest represents the request body for adding a association responsible
type AddRespoAssociationRequest struct {
	Email string `json:"email" validate:"required,email"`
}
