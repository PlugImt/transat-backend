package bassine

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"github.com/plugimt/transat-backend/handlers/bassine/repository"
	"github.com/plugimt/transat-backend/models"
	"github.com/plugimt/transat-backend/utils"
)

type BassineHandler struct {
	BassineRepository *repository.BassineRepository
}

func NewBassineHandler(db *sql.DB) *BassineHandler {
	return &BassineHandler{
		BassineRepository: repository.NewBassineRepository(db),
	}
}

// IncrementBassine handles PATCH /bassine?type=up|down
func (h *BassineHandler) IncrementBassine(c *fiber.Ctx) error {
	utils.LogHeader("🍺 Update Bassine Counter")

	email, ok := c.Locals("email").(string)
	if !ok || email == "" {
		utils.LogMessage(utils.LevelWarn, "User email not found in token")
		utils.LogFooter()
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	actionType := c.Query("type", "up")
	switch actionType {
	case "up":
		if err := h.BassineRepository.AddBassine(email); err != nil {
			utils.LogMessage(utils.LevelError, "Failed to increment bassine counter")
			utils.LogLineKeyValue(utils.LevelError, "Error", err)
			utils.LogFooter()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to increment"})
		}
	case "down":
		deleted, err := h.BassineRepository.RemoveBassine(email)
		if err != nil {
			utils.LogMessage(utils.LevelError, "Failed to decrement bassine counter")
			utils.LogLineKeyValue(utils.LevelError, "Error", err)
			utils.LogFooter()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to decrement"})
		}
		if !deleted {
			utils.LogMessage(utils.LevelWarn, "No bassine history to decrement")
			utils.LogFooter()
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "No bassine to decrement"})
		}
	default:
		utils.LogMessage(utils.LevelWarn, "Invalid type parameter for bassine update")
		utils.LogLineKeyValue(utils.LevelWarn, "type", actionType)
		utils.LogFooter()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "type must be 'up' or 'down'"})
	}

	overview, err := h.BassineRepository.GetUserBassineOverview(email)
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to fetch updated bassine overview")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch updated overview"})
	}

	utils.LogMessage(utils.LevelInfo, "Bassine counter updated successfully")
	utils.LogFooter()
	return c.JSON(overview)
}

// GetMyBassine handles GET /bassine/me
func (h *BassineHandler) GetMyBassine(c *fiber.Ctx) error {
	utils.LogHeader("🍺 Get My Bassine Overview")

	email, ok := c.Locals("email").(string)
	if !ok || email == "" {
		utils.LogMessage(utils.LevelWarn, "User email not found in token")
		utils.LogFooter()
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	overview, err := h.BassineRepository.GetUserBassineOverview(email)
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to get my bassine overview")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve bassine overview"})
	}

	utils.LogMessage(utils.LevelInfo, "Successfully retrieved bassine overview")
	utils.LogFooter()
	return c.JSON(overview)
}

// GetBassineLeaderboard handles GET /bassine/leaderboard
func (h *BassineHandler) GetBassineLeaderboard(c *fiber.Ctx) error {
	utils.LogHeader("🍺 Get Bassine Leaderboard")

	users, err := h.BassineRepository.GetLeaderboard()
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to get bassine leaderboard")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve leaderboard"})
	}

	utils.LogMessage(utils.LevelInfo, "Successfully retrieved bassine leaderboard")
	utils.LogFooter()
	return c.JSON(models.Leaderboard{Users: users})
}

// GetBassineHistory handles GET /bassine/history/:email
func (h *BassineHandler) GetBassineHistory(c *fiber.Ctx) error {
	utils.LogHeader("🍺 Get Bassine History")

	email := c.Params("email")
	if email == "" {
		utils.LogMessage(utils.LevelWarn, "Missing email parameter, using JWT email")
		email = c.Locals("email").(string)
		utils.LogFooter()
	}

	historyItem, err := h.BassineRepository.GetUserHistory(email)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.LogMessage(utils.LevelWarn, "User not found for history")
			utils.LogFooter()
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
		}
		utils.LogMessage(utils.LevelError, "Failed to get bassine history")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve history"})
	}

	utils.LogMessage(utils.LevelInfo, "Successfully retrieved bassine history")
	utils.LogFooter()
	return c.JSON(historyItem)
}
