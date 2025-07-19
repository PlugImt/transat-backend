package restaurant

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/plugimt/transat-backend/models"

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
	menuRepo := repository.NewMenuRepository(db, notifService)
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
		return c.Status(fiber.StatusOK).JSON(menuResponse)
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

func (h *RestaurantHandler) GetRestaurantTestMenu(c *fiber.Ctx) error {
	utils.LogHeader("🍽️ Get Restaurant Test Menu")

	utils.LogMessage(utils.LevelInfo, "Returning test menu data")
	utils.LogFooter()
	
	return c.JSON(&models.CategorizedMenuResponse{
		GrilladesMidi: []models.MenuItemWithRating{
			{ID: 20, Name: "poulet rôti", AverageRating: 4.5},
			{ID: 84, Name: "spaghetti à la bolognaise", AverageRating: 4.7},
			{ID: 68, Name: "vol au vent de colin", AverageRating: 4.2},
		},
		Migrateurs: []models.MenuItemWithRating{
			{ID: 165, Name: "estouffade de calamars à la nazaré", AverageRating: 4.8},
			{ID: 164, Name: "boulette de boeuf à la crème de champignons", AverageRating: 4.3},
		},
		Cibo: []models.MenuItemWithRating{
			{ID: 52, Name: "risotto aux champignons et courgettes", AverageRating: 4.6},
			{ID: 67, Name: "pizza végétarienne", AverageRating: 4.4},
		},
		AccompMidi: []models.MenuItemWithRating{
			{ID: 12, Name: "gâteau de choux fleurs noix et mimolette", AverageRating: 4.1},
			{ID: 15, Name: "fenouil rôti", AverageRating: 4.0},
		},
		GrilladesSoir: []models.MenuItemWithRating{
			{ID: 66, Name: "filet deiglegin croûte dolives", AverageRating: 4.5},
			{ID: 59, Name: "filet de poulet tandoori", AverageRating: 4.7},
			{ID: 86, Name: "pizza au chèvre", AverageRating: 4.2},
		},
		AccompSoir: []models.MenuItemWithRating{
			{ID: 103, Name: "haricots beurre persillées", AverageRating: 4.1},
			{ID: 50, Name: "lentilles vertes bio", AverageRating: 4.0},
		},
		UpdatedDate: utils.FormatParis(h.MenuRepository.GetLastMenuUpdateTime(), time.RFC3339),
	})
}

// CheckAndUpdateMenuCron is the cron endpoint
func (h *RestaurantHandler) CheckAndUpdateMenuCron() (bool, error) {
	return h.Scheduler.CheckAndUpdateMenuCron()
}
