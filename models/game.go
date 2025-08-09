package models

// BassineScoreRequest represents the payload to set/update a user's bassine score
type BassineScoreRequest struct {
    Score int `json:"score"`
}

// BassineScoreResponse represents the response for a user's bassine score
type BassineScoreResponse struct {
    Email string `json:"email"`
    Score int    `json:"score"`
    TopLeaderboard []string `json:"top_leaderboard,omitempty"`
}


// BassineUserScore represents one user's score for the bassine game
type BassineUserScore struct {
    User  ReservationUser `json:"user"`
    Score int             `json:"score"`
}


