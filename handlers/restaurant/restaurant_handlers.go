package restaurant

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/plugimt/transat-backend/handlers/restaurant/repository"
	"github.com/plugimt/transat-backend/handlers/restaurant/service"
	"github.com/plugimt/transat-backend/services"
	"github.com/plugimt/transat-backend/utils"

	"github.com/gofiber/fiber/v2"
)

type RestaurantHandler struct {
	MenuRepository *repository.MenuRepository
	MenuService    *service.MenuService
	Scheduler      *Scheduler
}

func NewRestaurantHandler(db *sql.DB, transService *services.TranslationService, notifService *services.NotificationService) *RestaurantHandler {
	menuRepo := repository.NewMenuRepository(db)
	menuService := service.NewMenuService()
	scheduler := NewScheduler(menuRepo, menuService, notifService)

	return &RestaurantHandler{
		MenuRepository: menuRepo,
		MenuService:    menuService,
		Scheduler:      scheduler,
	}
}

func (h *RestaurantHandler) GetRestaurantMenu(c *fiber.Ctx) error {
	utils.LogHeader("🍽️ Get Restaurant Menu")

	language := c.Query("language")

	if language == "" {
		utils.LogMessage(utils.LevelWarn, "No language specified, defaulting to French")
		language = "fr"
	}

	menuResponse, err := h.MenuRepository.GetTodaysMenuCategorized()
	if err != nil {
		utils.LogMessage(utils.LevelError, fmt.Sprintf("Failed to get today's menu: %v", err))
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve menu",
		})
	}

	totalItems := len(menuResponse.GrilladesMidi) + len(menuResponse.Migrateurs) + len(menuResponse.Cibo) +
		len(menuResponse.AccompMidi) + len(menuResponse.GrilladesSoir) + len(menuResponse.AccompSoir)

	if totalItems == 0 {
		utils.LogMessage(utils.LevelInfo, "No menu found for today")
		utils.LogFooter()
		return c.Status(fiber.StatusNoContent).JSON(menuResponse)
	}

	utils.LogMessage(utils.LevelInfo, fmt.Sprintf("Retrieved categorized menu with %d total items", totalItems))
	utils.LogFooter()

	return c.JSON(menuResponse)
}

// GetDishDetails retrieves detailed information about a specific dish
func (h *RestaurantHandler) GetDishDetails(c *fiber.Ctx) error {
	utils.LogHeader("🍽️ Get Dish Details")

	dishIDParam := c.Params("id")
	if dishIDParam == "" {
		utils.LogMessage(utils.LevelError, "Missing dish ID parameter")
		utils.LogFooter()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Dish ID is required",
		})
	}

	dishID, err := strconv.Atoi(dishIDParam)
	if err != nil {
		utils.LogMessage(utils.LevelError, fmt.Sprintf("Invalid dish ID: %s", dishIDParam))
		utils.LogFooter()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid dish ID format",
		})
	}

	dishDetails, err := h.MenuRepository.GetDishDetails(dishID)
	if err != nil {
		utils.LogMessage(utils.LevelError, fmt.Sprintf("Failed to get dish details: %v", err))
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve dish details",
		})
	}

	if dishDetails == nil {
		utils.LogMessage(utils.LevelError, fmt.Sprintf("Dish not found: ID %d", dishID))
		utils.LogFooter()
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Dish not found",
		})
	}

	utils.LogMessage(utils.LevelInfo, fmt.Sprintf("Retrieved details for dish"))
	utils.LogFooter()
	return c.JSON(dishDetails)
}

// PostDishReview allows users to post a review/rating for a specific dish
func (h *RestaurantHandler) PostDishReview(c *fiber.Ctx) error {
	utils.LogHeader("🌟 Post Dish Review")

	dishIDParam := c.Params("id")
	if dishIDParam == "" {
		utils.LogMessage(utils.LevelError, "Missing dish ID parameter")
		utils.LogFooter()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Dish ID is required",
		})
	}

	dishID, err := strconv.Atoi(dishIDParam)
	if err != nil {
		utils.LogMessage(utils.LevelError, fmt.Sprintf("Invalid dish ID: %s", dishIDParam))
		utils.LogFooter()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid dish ID format",
		})
	}

	userEmail, ok := c.Locals("email").(string)
	if !ok || userEmail == "" {
		utils.LogMessage(utils.LevelError, "No email found in JWT token")
		utils.LogFooter()
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authentication required",
		})
	}

	var reviewRequest struct {
		Rating  int    `json:"rating" validate:"required,min=1,max=5"`
		Comment string `json:"comment,omitempty"`
	}

	if err := c.BodyParser(&reviewRequest); err != nil {
		utils.LogMessage(utils.LevelError, fmt.Sprintf("Failed to parse request body: %v", err))
		utils.LogFooter()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request format",
		})
	}

	if reviewRequest.Rating < 1 || reviewRequest.Rating > 5 {
		utils.LogMessage(utils.LevelError, fmt.Sprintf("Invalid rating: %d", reviewRequest.Rating))
		utils.LogFooter()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Rating must be between 1 and 5",
		})
	}

	result, err := h.MenuRepository.SaveDishReview(dishID, userEmail, reviewRequest.Rating, reviewRequest.Comment)
	if err != nil {
		utils.LogMessage(utils.LevelError, fmt.Sprintf("Failed to save review: %v", err))
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save review",
		})
	}

	utils.LogMessage(utils.LevelInfo, fmt.Sprintf("Review saved for dish %s by %s: %d stars", result.DishName, userEmail, reviewRequest.Rating))
	utils.LogFooter()

	return c.JSON(result)
}

// CheckAndUpdateMenuCron is the cron endpoint
func (h *RestaurantHandler) CheckAndUpdateMenuCron() (bool, error) {
	return h.Scheduler.CheckAndUpdateMenuCron()
}
