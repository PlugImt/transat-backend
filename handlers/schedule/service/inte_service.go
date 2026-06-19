package service

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/plugimt/transat-backend/handlers/schedule/repository"
	"github.com/plugimt/transat-backend/models"
)

type InteScheduleService struct {
	repo *repository.InteScheduleRepository
}

func NewInteScheduleService(db *sql.DB) *InteScheduleService {
	return &InteScheduleService{
		repo: repository.NewInteScheduleRepository(db),
	}
}

func (s *InteScheduleService) GetCalendarData() (models.CalendarData, error) {
	schedule, err := s.repo.Get()
	if err != nil {
		return nil, err
	}

	if len(schedule.IcsJSON) > 0 {
		var calendarData models.CalendarData
		if err := json.Unmarshal(schedule.IcsJSON, &calendarData); err != nil {
			return nil, fmt.Errorf("unmarshal ics_json: %w", err)
		}
		return calendarData, nil
	}

	if schedule.RawICS == "" {
		return nil, sql.ErrNoRows
	}

	return s.processAndSave([]byte(schedule.RawICS))
}

func (s *InteScheduleService) processAndSave(rawICS []byte) (models.CalendarData, error) {
	calendarData, err := ParseICS(rawICS)
	if err != nil {
		return nil, err
	}

	raw, err := json.Marshal(calendarData)
	if err != nil {
		return nil, fmt.Errorf("marshal calendar data: %w", err)
	}

	if err := s.repo.UpdateIcsJSON(raw); err != nil {
		return nil, err
	}

	return calendarData, nil
}
