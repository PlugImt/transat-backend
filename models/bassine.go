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
