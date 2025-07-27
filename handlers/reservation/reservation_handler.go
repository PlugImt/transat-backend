package reservation

import (
	"database/sql"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/plugimt/transat-backend/handlers/reservation/repository"
	"github.com/plugimt/transat-backend/models"
	"github.com/plugimt/transat-backend/utils"
	"strconv"
	"time"
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

	var res models.ReservationOverviewResponse
	var categoryIDInt *int

	categoryID := c.Params("id")
	if categoryID != "" {
		parsedID, err := strconv.Atoi(categoryID)
		if err != nil {
			utils.LogMessage(utils.LevelError, fmt.Sprintf("Invalid category ID: %v", err))
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid category ID",
			})
		}
		categoryIDInt = &parsedID
	} else {
		categoryIDQuery := c.Query("category")
		if categoryIDQuery != "" {
			parsedID, err := strconv.Atoi(categoryIDQuery)
			if err != nil {
				utils.LogMessage(utils.LevelError, fmt.Sprintf("Invalid category ID: %v", err))
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": "Invalid category ID",
				})
			}
			categoryIDInt = &parsedID
		}
	}

	categoryList, err := h.ReservationRepository.GetCategoryList(categoryIDInt)
	if err != nil {
		utils.LogMessage(utils.LevelError, fmt.Sprintf("Failed to get category list: %v", err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve categories",
		})
	}
	utils.LogMessage(utils.LevelInfo, fmt.Sprintf("Retrieved %d categories", len(categoryList)))

	itemList, err := h.ReservationRepository.GetItemList(categoryIDInt)
	if err != nil {
		utils.LogMessage(utils.LevelError, fmt.Sprintf("Failed to get item list: %v", err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve items",
		})
	}

	utils.LogMessage(utils.LevelInfo, fmt.Sprintf("Retrieved %d items", len(itemList)))
	if len(categoryList) == 0 && len(itemList) == 0 {
		utils.LogMessage(utils.LevelInfo, "No categories or items found")
		utils.LogFooter()
		return c.JSON(res)
	}

	res.Categories = categoryList
	res.Items = itemList

	utils.LogMessage(utils.LevelInfo, "Successfully retrieved root reservation items")
	utils.LogFooter()
	return c.JSON(res)
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
	var res models.ReservationItem

	var StartDateDate time.Time
	var EndDateDate time.Time

	id := c.Params("id")
	if id == "" {
		utils.LogMessage(utils.LevelError, "Item ID is required")
		utils.LogFooter()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Item ID is required",
		})
	}
	itemID, err := strconv.Atoi(id)
	if err != nil {
		utils.LogMessage(utils.LevelError, fmt.Sprintf("Invalid item ID: %v", err))
		utils.LogFooter()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid item ID",
		})
	}

	var reservationRequest models.ReservationManagementRequest
	if err := c.BodyParser(&reservationRequest); err != nil {
		utils.LogMessage(utils.LevelError, fmt.Sprintf("Failed to parse request body: %v", err))
		utils.LogFooter()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request format",
		})
	}

	if reservationRequest.StartDate == nil && reservationRequest.EndDate == nil {
		utils.LogMessage(utils.LevelError, "At least one of StartDate or EndDate is required")
		utils.LogFooter()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "At least one of StartDate or EndDate is required",
		})
	}

	if reservationRequest.StartDate != nil && reservationRequest.EndDate != nil {
		utils.LogMessage(utils.LevelError, "Both StartDate and EndDate cannot be provided at the same time")
		utils.LogFooter()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Both StartDate and EndDate cannot be provided at the same time",
		})

	}

	if reservationRequest.StartDate != nil {
		StartDateDate, err = time.Parse("2006-01-02 15:04:05", *reservationRequest.StartDate)
		if err != nil {
			utils.LogMessage(utils.LevelError, fmt.Sprintf("Failed to parse StartDate: %v", err))
			utils.LogFooter()
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid StartDate format",
			})
		}
	}
	if reservationRequest.EndDate != nil {
		EndDateDate, err = time.Parse("2006-01-02 15:04:05", *reservationRequest.EndDate)
		if err != nil {
			utils.LogMessage(utils.LevelError, fmt.Sprintf("Failed to parse EndDate: %v", err))
			utils.LogFooter()
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid EndDate format",
			})
		}
	}

	ItemExists, err := h.ReservationRepository.CheckItemExists(itemID)
	if err != nil || !ItemExists {
		utils.LogMessage(utils.LevelError, fmt.Sprintf("Failed to check item existence: %v ", err))
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "The item does not exist",
		})
	}

	ItemPerSlot, err := h.ReservationRepository.IsItemPerSlot(itemID)
	if err != nil {
		utils.LogMessage(utils.LevelError, fmt.Sprintf("Failed to check if item is per slot: %v", err))
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to check if item is per slot",
		})
	}

	if ItemPerSlot && (StartDateDate.Minute() != 0 || StartDateDate.Second() != 0) {
		utils.LogMessage(utils.LevelError, "StartDate for per slot items must be on plain hours (e.g., 2023-10-01 14:00:00)")
		utils.LogFooter()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "StartDate for per slot items must be on plain hours (e.g., 2023-10-01 14:00:00)",
		})
	}
	if ItemPerSlot && reservationRequest.EndDate != nil {
		utils.LogMessage(utils.LevelError, "EndDate cannot be provided for per slot items")
		utils.LogFooter()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "EndDate cannot be provided for per slot items",
		})
	}

	reservationItemTime := models.ReservationManagementRequestTime{
		StartDate: &StartDateDate,
		EndDate:   &EndDateDate,
	}

	if reservationRequest.StartDate != nil {
		res, err = h.ReservationRepository.CreateReservation(reservationItemTime, itemID, ItemPerSlot, c.Locals("email").(string))
		if err != nil {
			utils.LogMessage(utils.LevelError, fmt.Sprintf("Failed to create reservation: %v", err))
			utils.LogFooter()

			if err.Error() == "item is not available" {
				return c.Status(fiber.StatusConflict).JSON(fiber.Map{
					"error": "Item is not available for the requested time",
				})
			}

			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to create reservation",
			})
		}
	} else if reservationRequest.EndDate != nil {
		res, err = h.ReservationRepository.EndReservation(reservationItemTime, itemID, ItemPerSlot, c.Locals("email").(string))
		if err != nil {
			utils.LogMessage(utils.LevelError, fmt.Sprintf("Failed to end reservation: %v", err))
			utils.LogFooter()

			if err.Error() == "item is not reserved" {
				return c.Status(fiber.StatusConflict).JSON(fiber.Map{
					"error": "Item is not reserved for the requested time",
				})
			}

			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to end reservation",
			})
		}
	}

	utils.LogMessage(utils.LevelInfo, "Reservation updated successfully")
	utils.LogFooter()

	return c.JSON(fiber.Map{
		"message":     "Reservation updated successfully",
		"reservation": res,
	})
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

	var StartDateDate time.Time

	id := c.Params("id")
	if id == "" {
		utils.LogMessage(utils.LevelError, "Item ID is required")
		utils.LogFooter()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Item ID is required",
		})
	}

	itemID, err := strconv.Atoi(id)
	if err != nil {
		utils.LogMessage(utils.LevelError, fmt.Sprintf("Invalid item ID: %v", err))
		utils.LogFooter()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid item ID",
		})
	}

	var reservationRequest models.ReservationManagementRequest
	if err := c.BodyParser(&reservationRequest); err != nil {
		utils.LogMessage(utils.LevelError, fmt.Sprintf("Failed to parse request body: %v", err))
		utils.LogFooter()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request format",
		})
	}

	itemExists, err := h.ReservationRepository.CheckItemExists(itemID)
	if err != nil {
		utils.LogMessage(utils.LevelError, fmt.Sprintf("Failed to check item existence: %v", err))
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to check item existence",
		})
	}

	if !itemExists {
		utils.LogMessage(utils.LevelError, fmt.Sprintf("Item with ID %d does not exist", itemID))
		utils.LogFooter()
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": fmt.Sprintf("Item with ID %d does not exist", itemID),
		})
	}

	if reservationRequest.StartDate != nil {
		StartDateDate, err = time.Parse("2006-01-02 15:04:05", *reservationRequest.StartDate)
		if err != nil {
			utils.LogMessage(utils.LevelError, fmt.Sprintf("Failed to parse StartDate: %v", err))
			utils.LogFooter()
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid StartDate format",
			})
		}
	}

	ItemPerSlot, err := h.ReservationRepository.IsItemPerSlot(itemID)
	if err != nil {
		utils.LogMessage(utils.LevelError, fmt.Sprintf("Failed to check if item is per slot: %v", err))
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to check if item is per slot",
		})
	}

	if !ItemPerSlot {
		utils.LogMessage(utils.LevelError, "Item is not a per slot item, cannot delete reservation")
		utils.LogFooter()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Item is not a per slot item, cannot delete reservation",
		})
	}

	reservationItemTime := models.ReservationManagementRequestTime{
		StartDate: &StartDateDate,
	}

	deleted, err := h.ReservationRepository.DeleteReservation(reservationItemTime, itemID, ItemPerSlot, c.Locals("email").(string))
	if err != nil {
		utils.LogMessage(utils.LevelError, fmt.Sprintf("Failed to delete reservation: %v", err))
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete reservation",
		})
	}

	utils.LogMessage(utils.LevelInfo, "Reservation deleted successfully")
	utils.LogFooter()
	return c.JSON(fiber.Map{
		"message":     "Reservation deleted successfully",
		"reservation": deleted,
	})
}
