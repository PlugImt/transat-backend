package repository

import (
	"database/sql"
	"time"

	"github.com/plugimt/transat-backend/models"
)

type BassineRepository struct {
	DB *sql.DB
}

func NewBassineRepository(db *sql.DB) *BassineRepository {
	return &BassineRepository{
		DB: db,
	}
}

func (r *BassineRepository) AddBassine(email string) error {
	query := `INSERT INTO bassine_history (email, date) VALUES ($1, $2)`
	_, err := r.DB.Exec(query, email, time.Now())
	return err
}

func (r *BassineRepository) RemoveBassine(email string) (bool, error) {
	// Delete latest entry for this user
	query := `
        DELETE FROM bassine_history
        WHERE id_bassine_history = (
            SELECT id_bassine_history
            FROM bassine_history
            WHERE email = $1
            ORDER BY date DESC
            LIMIT 1
        )
    `
	res, err := r.DB.Exec(query, email)
	if err != nil {
		return false, err
	}
	affected, _ := res.RowsAffected()
	return affected > 0, nil
}

func (r *BassineRepository) GetLeaderboard() ([]models.BassineUser, error) {
	query := `
        SELECT
            n.email,
            n.first_name,
            n.last_name,
            COALESCE(n.profile_picture, '') AS profile_picture,
            s.score,
            RANK() OVER (ORDER BY s.score DESC, n.email ASC) AS rank
        FROM bassine_scores s
        JOIN newf n ON n.email = s.email
        ORDER BY s.score DESC, n.email ASC;
    `

	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]models.BassineUser, 0)
	for rows.Next() {
		var user models.BassineUser
		var email, firstName, lastName, picture string
		var score, rank int
		if err := rows.Scan(&email, &firstName, &lastName, &picture, &score, &rank); err != nil {
			return nil, err
		}
		user.ReservationUser = &models.ReservationUser{
			Email:          email,
			FirstName:      firstName,
			LastName:       lastName,
			ProfilePicture: picture,
		}
		user.BassineCount = score
		user.Rank = rank
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}

func (r *BassineRepository) GetUserBassineOverview(email string) (models.BassineCount, error) {
	// Fetch the full leaderboard in one query and then compute neighbors in memory for simplicity
	users, err := r.GetLeaderboard()
	if err != nil {
		return models.BassineCount{}, err
	}
	index := -1
	for i, u := range users {
		if u.ReservationUser != nil && u.ReservationUser.Email == email {
			index = i
			break
		}
	}
	if index == -1 {
		var firstName, lastName, picture string
		err := r.DB.QueryRow(`SELECT first_name, last_name, COALESCE(profile_picture,'') FROM newf WHERE email = $1`, email).Scan(&firstName, &lastName, &picture)
		if err != nil {
			return models.BassineCount{}, err
		}
		rank := len(users) + 1
		result := models.BassineCount{
			Email:        email,
			Rank:         rank,
			BassineCount: 0,
		}
		if len(users) > 0 {
			result.UserAbove = users[len(users)-1]
		}
		_ = firstName
		_ = lastName
		_ = picture
		return result, nil
	}

	result := models.BassineCount{
		Email:        email,
		Rank:         users[index].Rank,
		BassineCount: users[index].BassineCount,
	}
	if index-1 >= 0 {
		result.UserAbove = users[index-1]
	}
	if index+1 < len(users) {
		result.UserBelow = users[index+1]
	}
	return result, nil
}

func (r *BassineRepository) GetUserHistory(email string) (models.BassineHistoryItem, error) {
	var firstName, lastName, picture string
	err := r.DB.QueryRow(`SELECT first_name, last_name, COALESCE(profile_picture,'') FROM newf WHERE email = $1`, email).Scan(&firstName, &lastName, &picture)
	if err != nil {
		return models.BassineHistoryItem{}, err
	}

	rows, err := r.DB.Query(`SELECT date FROM bassine_history WHERE email = $1 ORDER BY date ASC`, email)
	if err != nil {
		return models.BassineHistoryItem{}, err
	}
	defer rows.Close()

	dates := make([]time.Time, 0)
	for rows.Next() {
		var d time.Time
		if err := rows.Scan(&d); err != nil {
			return models.BassineHistoryItem{}, err
		}
		dates = append(dates, d)
	}
	if err := rows.Err(); err != nil {
		return models.BassineHistoryItem{}, err
	}

	history := models.BassineHistoryItem{
		ReservationUser: &models.ReservationUser{
			Email:          email,
			FirstName:      firstName,
			LastName:       lastName,
			ProfilePicture: picture,
		},
		Dates: dates,
	}
	return history, nil
}
