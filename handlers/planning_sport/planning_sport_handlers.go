package planning_sport

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"github.com/plugimt/transat-backend/models"
	"github.com/plugimt/transat-backend/utils"
)

type PlanningSportHandler struct {
	db *sql.DB
}

func NewPlanningSportHandler(db *sql.DB) *PlanningSportHandler {
	return &PlanningSportHandler{db: db}
}

func (h *PlanningSportHandler) GetPlanning(c *fiber.Ctx) error {
	query := `
		SELECT id, day_of_week, activity, place, start_time, end_time
		FROM planning_sport
		ORDER BY id
	`

	rows, err := h.db.Query(query)
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to fetch planning")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch planning",
		})
	}
	defer rows.Close()

	var activities []models.PlanningSport
	for rows.Next() {
		var a models.PlanningSport
		var place sql.NullString

		err := rows.Scan(
			&a.ID,
			&a.Day,
			&a.Activity,
			&place,
			&a.StartTime,
			&a.EndTime,
		)
		if err != nil {
			utils.LogMessage(utils.LevelError, "Failed to scan planning row")
			continue
		}

		if place.Valid {
			a.Place = &place.String
		}

		activities = append(activities, a)
	}

	utils.LogMessage(utils.LevelInfo, "Successfully fetched planning")
	utils.LogLineKeyValue(utils.LevelInfo, "Count", len(activities))
	utils.LogFooter()

	return c.JSON(activities)
}