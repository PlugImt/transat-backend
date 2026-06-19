package schedule

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"github.com/plugimt/transat-backend/handlers/schedule/repository"
	"github.com/plugimt/transat-backend/handlers/schedule/service"
	"github.com/plugimt/transat-backend/models"
	"github.com/plugimt/transat-backend/utils"
)

type Handler struct {
	userRepo    *repository.UserScheduleRepository
	icsService  *service.IcsService
	inteService *service.InteScheduleService
}

func NewHandler(db *sql.DB, icsService *service.IcsService) *Handler {
	return &Handler{
		userRepo:    repository.NewUserScheduleRepository(db),
		icsService:  icsService,
		inteService: service.NewInteScheduleService(db),
	}
}

// GetMySchedule handles GET /schedule/me
func (h *Handler) GetMySchedule(c *fiber.Ctx) error {
	utils.LogHeader("📅 Get User Schedule")

	email, ok := c.Locals("email").(string)
	if !ok || email == "" {
		utils.LogMessage(utils.LevelWarn, "User email not found in token")
		utils.LogFooter()
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	schedule, err := h.userRepo.GetByEmail(email)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.LogMessage(utils.LevelWarn, "No schedule found for user")
			utils.LogFooter()
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Schedule not found"})
		}
		utils.LogMessage(utils.LevelError, "Failed to get user schedule")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve schedule"})
	}

	utils.LogMessage(utils.LevelInfo, "Successfully retrieved user schedule")
	utils.LogFooter()
	return c.JSON(schedule)
}

// UpdateMySchedule handles PATCH /schedule/me
func (h *Handler) UpdateMySchedule(c *fiber.Ctx) error {
	utils.LogHeader("📅 Update User Schedule")

	email, ok := c.Locals("email").(string)
	if !ok || email == "" {
		utils.LogMessage(utils.LevelWarn, "User email not found in token")
		utils.LogFooter()
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	var req models.UpdateUserScheduleRequest
	if err := c.BodyParser(&req); err != nil {
		utils.LogMessage(utils.LevelWarn, "Invalid request body")
		utils.LogLineKeyValue(utils.LevelWarn, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}
	if req.IcsURL == "" {
		utils.LogMessage(utils.LevelWarn, "ics_url is required")
		utils.LogFooter()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ics_url is required"})
	}

	if err := h.userRepo.UpsertIcsURL(email, req.IcsURL); err != nil {
		utils.LogMessage(utils.LevelError, "Failed to update user schedule")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update schedule"})
	}

	schedule, err := h.userRepo.GetByEmail(email)
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to fetch updated user schedule")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve updated schedule"})
	}

	go func(userID int, icsURL string) {
		if err := h.icsService.SyncUserSchedule(userID, icsURL); err != nil {
			utils.LogMessage(utils.LevelError, "Background ICS sync failed")
			utils.LogLineKeyValue(utils.LevelError, "UserID", userID)
			utils.LogLineKeyValue(utils.LevelError, "Error", err)
		}
	}(schedule.UserID, req.IcsURL)

	utils.LogMessage(utils.LevelInfo, "Successfully updated user schedule")
	utils.LogFooter()
	return c.JSON(schedule)
}

// DeleteMySchedule handles DELETE /schedule/me
func (h *Handler) DeleteMySchedule(c *fiber.Ctx) error {
	utils.LogHeader("📅 Delete User Schedule")

	email, ok := c.Locals("email").(string)
	if !ok || email == "" {
		utils.LogMessage(utils.LevelWarn, "User email not found in token")
		utils.LogFooter()
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	deleted, err := h.userRepo.Delete(email)
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to delete user schedule")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete schedule"})
	}
	if !deleted {
		utils.LogMessage(utils.LevelWarn, "No schedule found to delete")
		utils.LogFooter()
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Schedule not found"})
	}

	utils.LogMessage(utils.LevelInfo, "Successfully deleted user schedule")
	utils.LogFooter()
	return c.SendStatus(fiber.StatusNoContent)
}

// GetInteSchedule handles GET /schedule/inte
func (h *Handler) GetInteSchedule(c *fiber.Ctx) error {
	utils.LogHeader("📅 Get ITE Schedule")

	calendarData, err := h.inteService.GetCalendarData()
	if err != nil {
		if err == sql.ErrNoRows {
			utils.LogMessage(utils.LevelWarn, "No ITE schedule found")
			utils.LogFooter()
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "INTE schedule not found"})
		}
		utils.LogMessage(utils.LevelError, "Failed to get ITE schedule")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve ITE schedule"})
	}

	utils.LogMessage(utils.LevelInfo, "Successfully retrieved ITE schedule")
	utils.LogFooter()
	return c.JSON(calendarData)
}
