package repository

import (
	"database/sql"
	"encoding/json"

	"github.com/plugimt/transat-backend/models"
)

type InteScheduleRepository struct {
	DB *sql.DB
}

func NewInteScheduleRepository(db *sql.DB) *InteScheduleRepository {
	return &InteScheduleRepository{DB: db}
}

func (r *InteScheduleRepository) Get() (*models.InteSchedule, error) {
	query := `
		SELECT raw_ics, ics_json
		FROM inte_schedule
		LIMIT 1
	`

	var schedule models.InteSchedule
	var rawICS sql.NullString
	var icsJSON []byte

	err := r.DB.QueryRow(query).Scan(&rawICS, &icsJSON)
	if err != nil {
		return nil, err
	}

	if rawICS.Valid {
		schedule.RawICS = rawICS.String
	}
	if len(icsJSON) > 0 {
		schedule.IcsJSON = json.RawMessage(icsJSON)
	}

	return &schedule, nil
}

func (r *InteScheduleRepository) UpdateIcsJSON(icsJSON json.RawMessage) error {
	query := `
		UPDATE inte_schedule
		SET ics_json = $1
		WHERE raw_ics IS NOT NULL
	`
	res, err := r.DB.Exec(query, icsJSON)
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
