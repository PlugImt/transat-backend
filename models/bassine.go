package models

import "time"

type BassineUser struct {
	*ReservationUser
	Rank         int `json:"rank" bson:"rank"`
	BassineCount int `json:"bassine_count" bson:"bassine_count"`
}

type BassineCount struct {
	Email        string      `json:"email" bson:"email"`
	Rank         int         `json:"rank" bson:"rank"`
	BassineCount int         `json:"bassine_count" bson:"bassine_count"`
	UserAbove    BassineUser `json:"user_above,omitempty" bson:"user_above,omitempty"`
	UserBelow    BassineUser `json:"user_below,omitempty" bson:"user_below,omitempty"`
}

type Leaderboard struct {
	Users []BassineUser `json:"users" bson:"users"`
}

type BassineHistoryItem struct {
	*ReservationUser
	Dates []time.Time `json:"dates" bson:"dates"`
}

// Admin-specific types
type AdminBassineScore struct {
	ID                int         `json:"id"`
	UserEmail         string      `json:"user_email"`
	UserFirstName     string      `json:"user_first_name"`
	UserLastName      string      `json:"user_last_name"`
	CurrentScore      int         `json:"current_score"`
	TotalGamesPlayed  int         `json:"total_games_played"`
	CreationDate      time.Time   `json:"creation_date"`
	LastUpdated       time.Time   `json:"last_updated"`
}

type AdminBassineHistory struct {
	ID          int       `json:"id"`
	UserEmail   string    `json:"user_email"`
	ScoreChange int       `json:"score_change"`
	NewTotal    int       `json:"new_total"`
	GameDate    time.Time `json:"game_date"`
	Notes       *string   `json:"notes,omitempty"`
	AdminEmail  *string   `json:"admin_email,omitempty"`
}

type UpdateScoreRequest struct {
	UserEmail   string  `json:"userEmail" validate:"required,email"`
	ScoreChange int     `json:"scoreChange" validate:"required"`
	Notes       *string `json:"notes,omitempty"`
}
