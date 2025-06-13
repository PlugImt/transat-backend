package models

import "time"

// Course représente un cours dans le planning
type Course struct {
	ID        int       `json:"id,omitempty" db:"id" example:"1"`
	Date      string    `json:"date" db:"date" example:"2024-01-15" binding:"required"`
	Title     string    `json:"title" db:"title" example:"Mathématiques" binding:"required"`
	StartTime string    `json:"start_time" db:"start_time" example:"09:00" binding:"required"`
	EndTime   string    `json:"end_time" db:"end_time" example:"10:30" binding:"required"`
	Teacher   string    `json:"teacher" db:"teacher" example:"Prof. Dupont"`
	Room      string    `json:"room" db:"room" example:"A101"`
	Group     string    `json:"group" db:"group" example:"L3-INFO"`
	UserEmail string    `json:"user_email" db:"user_email" example:"student@example.com" binding:"required"`
	CreatedAt time.Time `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at,omitempty" db:"updated_at"`
}

// CreateCourseRequest représente la requête pour créer un nouveau cours
type CreateCourseRequest struct {
	Date      string `json:"date" example:"2024-01-15" binding:"required"`
	Title     string `json:"title" example:"Mathématiques" binding:"required"`
	StartTime string `json:"start_time" example:"09:00" binding:"required"`
	EndTime   string `json:"end_time" example:"10:30" binding:"required"`
	Teacher   string `json:"teacher" example:"Prof. Dupont"`
	Room      string `json:"room" example:"A101"`
	Group     string `json:"group" example:"L3-INFO"`
	UserEmail string `json:"user_email" example:"student@example.com" binding:"required"`
}

// UserWithPassID représente un utilisateur avec son Pass ID pour la planification
type UserWithPassID struct {
	ID        int    `json:"id" example:"1"`
	FirstName string `json:"first_name" example:"Jean"`
	LastName  string `json:"last_name" example:"Dupont"`
	PassID    *int64 `json:"pass_id" example:"123456789"`
}

// UpdatePassIDRequest représente la requête pour mettre à jour un Pass ID
type UpdatePassIDRequest struct {
	PassID int `json:"pass_id" example:"123456789" binding:"required"`
}
