package repository

import (
	"database/sql"
	"encoding/json"

	"github.com/plugimt/transat-backend/models"
)

type UserScheduleRepository struct {
	DB *sql.DB
}

func NewUserScheduleRepository(db *sql.DB) *UserScheduleRepository {
	return &UserScheduleRepository{DB: db}
}

func (r *UserScheduleRepository) GetByEmail(email string) (*models.UserSchedule, error) {
	query := `
		SELECT us.user_id, us.ics_url, us.last_sync_at, us.calendar_data
		FROM user_schedule us
		JOIN newf n ON n.id_newf = us.user_id
		WHERE n.email = $1
	`

	var schedule models.UserSchedule
	var icsURL sql.NullString
	var lastSyncAt sql.NullTime
	var calendarData []byte

	err := r.DB.QueryRow(query, email).Scan(
		&schedule.UserID,
		&icsURL,
		&lastSyncAt,
		&calendarData,
	)
	if err != nil {
		return nil, err
	}

	if icsURL.Valid {
		schedule.IcsURL = icsURL.String
	}
	if lastSyncAt.Valid {
		schedule.LastSyncAt = &lastSyncAt.Time
	}
	if len(calendarData) > 0 {
		schedule.CalendarData = json.RawMessage(calendarData)
	}

	return &schedule, nil
}

func (r *UserScheduleRepository) UpsertIcsURL(email string, icsURL string) error {
	query := `
		INSERT INTO user_schedule (user_id, ics_url)
		SELECT id_newf, $2 FROM newf WHERE email = $1
		ON CONFLICT (user_id) DO UPDATE SET ics_url = EXCLUDED.ics_url
	`
	res, err := r.DB.Exec(query, email, icsURL)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *UserScheduleRepository) UpdateCalendarData(userID int, data json.RawMessage) error {
	query := `
		UPDATE user_schedule
		SET calendar_data = $2, last_sync_at = NOW()
		WHERE user_id = $1
	`
	res, err := r.DB.Exec(query, userID, data)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *UserScheduleRepository) ListAllWithIcsURL() ([]models.UserSchedule, error) {
	query := `
		SELECT user_id, ics_url
		FROM user_schedule
		WHERE ics_url IS NOT NULL AND ics_url != ''
	`
	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var schedules []models.UserSchedule
	for rows.Next() {
		var schedule models.UserSchedule
		var icsURL sql.NullString
		if err := rows.Scan(&schedule.UserID, &icsURL); err != nil {
			return nil, err
		}
		if icsURL.Valid {
			schedule.IcsURL = icsURL.String
		}
		schedules = append(schedules, schedule)
	}
	return schedules, rows.Err()
}

func (r *UserScheduleRepository) Delete(email string) (bool, error) {
	query := `
		DELETE FROM user_schedule
		WHERE user_id = (SELECT id_newf FROM newf WHERE email = $1)
	`
	res, err := r.DB.Exec(query, email)
	if err != nil {
		return false, err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return false, err
	}
	return affected > 0, nil
}
