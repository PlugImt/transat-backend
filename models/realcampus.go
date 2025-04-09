package models

import (
	"database/sql"
	"time"
)

// RealCampusPost represents a post in the RealCampus system
type RealCampusPost struct {
	ID           int       `json:"id"`
	FileID1      int       `json:"file_id_1"`
	FileID2      int       `json:"file_id_2"`
	AuthorEmail  string    `json:"author_email"`
	Location     string    `json:"location"`
	Privacy      string    `json:"privacy"`
	CreationDate time.Time `json:"creation_date"`
}

// RealCampusReaction represents a reaction to a post
type RealCampusReaction struct {
	ID           int       `json:"id"`
	PostID       int       `json:"post_id"`
	FileID       int       `json:"file_id"`
	AuthorEmail  string    `json:"author_email"`
	CreationDate time.Time `json:"creation_date"`
}

// RealCampusFriendship represents a friendship connection
type RealCampusFriendship struct {
	ID             int          `json:"id"`
	UserID         string       `json:"user_id"`
	FriendID       string       `json:"friend_id"`
	Status         string       `json:"status"`
	RequestDate    time.Time    `json:"request_date"`
	AcceptanceDate sql.NullTime `json:"acceptance_date"`
}

// RealCampusFile extends the base File model for RealCampus functionality
type RealCampusFile struct {
	ID             int       `json:"id"`
	Name           string    `json:"name"`
	Path           string    `json:"path"`
	Email          string    `json:"email"`
	CreationDate   time.Time `json:"creation_date"`
	LastAccessDate time.Time `json:"last_access_date"`
}
