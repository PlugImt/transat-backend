package reservation

import (
	"database/sql"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/plugimt/transat-backend/handlers/reservation/repository"
	"github.com/plugimt/transat-backend/models"
	"github.com/plugimt/transat-backend/utils"
)

type ReservationHandler struct {
	ReservationRepository *repository.ReservationRepository
}

func NewReservationHandler(db *sql.DB) *ReservationHandler {
	return &ReservationHandler{
		ReservationRepository: repository.NewReservationRepository(db),
	}
}

// GetReservationItems handles GET /reservation - returns root categories and items
func (h *ReservationHandler) GetReservationItems(c *fiber.Ctx) error {
	utils.LogHeader("📅 Get Root Reservation Items")
	return c.JSON(fiber.Map{})
}

// GetReservationCategoryItemsByID handles GET /reservation/category/{id}
func (h *ReservationHandler) GetReservationCategoryItemsByID(c *fiber.Ctx) error {
	utils.LogHeader("📅 Get Category Reservation Items")
	return c.JSON(fiber.Map{})
}

// CreateReservationCategory handles POST /reservation/category
func (h *ReservationHandler) CreateReservationCategory(c *fiber.Ctx) error {
	utils.LogHeader("📅 Create Reservation Category")

	var categoryRequest models.ReservationCreateCategoryRequest

	if err := c.BodyParser(&categoryRequest); err != nil {
		utils.LogMessage(utils.LevelError, fmt.Sprintf("Failed to parse request body: %v", err))
		utils.LogFooter()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request format",
		})
	}

	if categoryRequest.Name == "" {
		utils.LogMessage(utils.LevelError, "Category name is required")
		utils.LogFooter()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Category name is required",
		})
	}
	if categoryRequest.IDClubParent == nil && categoryRequest.IDCategoryParent == nil {
		utils.LogMessage(utils.LevelError, "Club ID or Category ID is required")
		utils.LogFooter()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Club ID or Category ID is required",
		})
	}

	if categoryRequest.IDClubParent != nil {
		clubExists, err := h.ReservationRepository.CheckClubExists(*categoryRequest.IDClubParent)
		if err != nil {
			utils.LogMessage(utils.LevelError, fmt.Sprintf("Failed to check club existence: %v", err))
			utils.LogFooter()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to check club existence",
			})
		}

		if !clubExists {
			utils.LogMessage(utils.LevelError, fmt.Sprintf("Club with ID %d does not exist", categoryRequest.IDClubParent))
			utils.LogFooter()
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": fmt.Sprintf("Club with ID %d does not exist", categoryRequest.IDClubParent),
			})
		}
	}

	if categoryRequest.IDCategoryParent != nil {
		categoryExists, err := h.ReservationRepository.CheckCategoryExists(*categoryRequest.IDCategoryParent)
		if err != nil {
			utils.LogMessage(utils.LevelError, fmt.Sprintf("Failed to check category existence: %v", err))
			utils.LogFooter()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to check category existence",
			})
		}

		if !categoryExists {
			utils.LogMessage(utils.LevelError, fmt.Sprintf("Category with ID %d does not exist", *categoryRequest.IDCategoryParent))
			utils.LogFooter()
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": fmt.Sprintf("Category with ID %d does not exist", *categoryRequest.IDCategoryParent),
			})
		}
	}

	if categoryRequest.IDClubParent == nil {
		IDClubParent, err := h.ReservationRepository.GetIDClubParent(*categoryRequest.IDCategoryParent)
		if err != nil {
			utils.LogMessage(utils.LevelError, fmt.Sprintf("Failed to get parent club ID: %v", err))
			utils.LogFooter()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to get parent club ID",
			})
		}

		categoryRequest.IDClubParent = &IDClubParent
	}

	created, err := h.ReservationRepository.CreateCategory(categoryRequest)
	if err != nil {
		utils.LogMessage(utils.LevelError, fmt.Sprintf("Failed to create category: %v", err))
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create category",
		})
	}

	utils.LogMessage(utils.LevelInfo, "Category created successfully")
	utils.LogFooter()

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":  "Category created successfully",
		"category": created,
	})
}

// GetItemDetails handles GET /reservation/items/{id}
func (h *ReservationHandler) GetItemDetails(c *fiber.Ctx) error {
	utils.LogHeader("📅 Get Item Details")
	return c.JSON(fiber.Map{})
}

// UpdateReservationItem handles PATCH /reservation/{id} - make or end reservations
func (h *ReservationHandler) UpdateReservationItem(c *fiber.Ctx) error {
	utils.LogHeader("📅 Update Reservation Item")
	return c.JSON(fiber.Map{})
}

// CreateReservationItem handles POST /reservation/item/
func (h *ReservationHandler) CreateReservationItem(c *fiber.Ctx) error {
	utils.LogHeader("📅 Create Reservation Item")

	var itemRequest models.ReservationCreateItemRequest

	if err := c.BodyParser(&itemRequest); err != nil {
		utils.LogMessage(utils.LevelError, fmt.Sprintf("Failed to parse request body: %v", err))
		utils.LogFooter()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request format",
		})
	}

	if itemRequest.Name == "" {
		utils.LogMessage(utils.LevelError, "Item name is required")
		utils.LogFooter()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Item name is required",
		})
	}
	if itemRequest.IDClubParent == nil && itemRequest.IDCategoryParent == nil {
		utils.LogMessage(utils.LevelError, "Either Club ID or Category Parent is required")
		utils.LogFooter()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Either Club ID or Category Parent is required",
		})
	}
	if itemRequest.IDClubParent != nil {
		clubExists, err := h.ReservationRepository.CheckClubExists(*itemRequest.IDClubParent)
		if err != nil {
			utils.LogMessage(utils.LevelError, fmt.Sprintf("Failed to check club existence: %v", err))
			utils.LogFooter()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to check club existence",
			})
		}

		if !clubExists {
			utils.LogMessage(utils.LevelError, fmt.Sprintf("Club with ID %d does not exist", *itemRequest.IDClubParent))
			utils.LogFooter()
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": fmt.Sprintf("Club with ID %d does not exist", *itemRequest.IDClubParent),
			})
		}
	}
	if itemRequest.IDCategoryParent != nil {
		categoryExists, err := h.ReservationRepository.CheckCategoryExists(*itemRequest.IDCategoryParent)
		if err != nil {
			utils.LogMessage(utils.LevelError, fmt.Sprintf("Failed to check category existence: %v", err))
			utils.LogFooter()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to check category existence",
			})
		}

		if !categoryExists {
			utils.LogMessage(utils.LevelError, fmt.Sprintf("Category with ID %d does not exist", *itemRequest.IDCategoryParent))
			utils.LogFooter()
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": fmt.Sprintf("Category with ID %d does not exist", *itemRequest.IDCategoryParent),
			})
		}
	}

	if itemRequest.IDClubParent == nil {
		IDClubParent, err := h.ReservationRepository.GetIDClubParent(*itemRequest.IDCategoryParent)
		if err != nil {
			utils.LogMessage(utils.LevelError, fmt.Sprintf("Failed to get parent club ID: %v", err))
			utils.LogFooter()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to get parent club ID",
			})
		}
		itemRequest.IDClubParent = &IDClubParent
	}

	createdItem, err := h.ReservationRepository.CreateItem(itemRequest)
	if err != nil {
		utils.LogMessage(utils.LevelError, fmt.Sprintf("Failed to create item: %v", err))
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create item",
		})
	}

	utils.LogMessage(utils.LevelInfo, "Item created successfully")
	utils.LogFooter()
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Item created successfully",
		"item":    createdItem,
	})
}

// GetReservationCategories handles GET /reservation/category
func (h *ReservationHandler) GetReservationCategories(c *fiber.Ctx) error {
	utils.LogHeader("📅 Get All Reservation Categories")
	return c.JSON(fiber.Map{})
}

// DeleteReservationItem handles DELETE /reservation/{id}
func (h *ReservationHandler) DeleteReservationItem(c *fiber.Ctx) error {
	utils.LogHeader("📅 Delete Reservation Item")
	return c.JSON(fiber.Map{})
}
