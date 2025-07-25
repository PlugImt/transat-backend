package repository

import (
	"database/sql"
	"fmt"
	"github.com/pkg/errors"
	"github.com/plugimt/transat-backend/models"
	"github.com/plugimt/transat-backend/utils"
)

type ReservationRepository struct {
	DB *sql.DB
}

func NewReservationRepository(db *sql.DB) *ReservationRepository {
	return &ReservationRepository{
		DB: db,
	}
}

func (r *ReservationRepository) CheckClubExists(IDClub int) (bool, error) {
	query := "SELECT COUNT(*) FROM clubs WHERE id_clubs = $1"
	row := r.DB.QueryRow(query, IDClub)
	var count int
	if err := row.Scan(&count); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.LogMessage(utils.LevelWarn, fmt.Sprintf("Club with ID %d does not exist", IDClub))
			return false, nil
		}
		utils.LogMessage(utils.LevelError, fmt.Sprintf("Failed to check club existence: %v", err))
		return false, err
	}
	if count > 0 {
		utils.LogMessage(utils.LevelInfo, fmt.Sprintf("Club with ID %d exists", IDClub))
		return true, nil
	} else {
		utils.LogMessage(utils.LevelWarn, fmt.Sprintf("Club with ID %d does not exist", IDClub))
	}
	// If no rows found, the club does not exist
	utils.LogMessage(utils.LevelInfo, fmt.Sprintf("Club with ID %d does not exist", IDClub))

	return false, nil
}

func (r *ReservationRepository) CheckCategoryExists(IDCategory int) (bool, error) {
	query := "SELECT COUNT(*) FROM reservation_category WHERE id_reservation_category = $1"
	row := r.DB.QueryRow(query, IDCategory)
	var count int
	if err := row.Scan(&count); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.LogMessage(utils.LevelWarn, fmt.Sprintf("Category with ID %d does not exist", IDCategory))
			return false, nil
		}
		utils.LogMessage(utils.LevelError, fmt.Sprintf("Failed to check category existence: %v", err))
		return false, err
	}
	if count > 0 {
		utils.LogMessage(utils.LevelInfo, fmt.Sprintf("Category with ID %d exists", IDCategory))
		return true, nil
	} else {
		utils.LogMessage(utils.LevelWarn, fmt.Sprintf("Category with ID %d does not exist", IDCategory))
	}
	// If no rows found, the category does not exist
	utils.LogMessage(utils.LevelInfo, fmt.Sprintf("Category with ID %d does not exist", IDCategory))

	return false, nil
}

func (r *ReservationRepository) CreateCategory(category models.ReservationCreateCategoryRequest) (bool, error) {
	if category.Name == "" {
		utils.LogMessage(utils.LevelError, "Category name is required")
		return false, fmt.Errorf("category name is required")
	}

	query := "INSERT INTO reservation_category (name, id_clubs, id_parent_category) VALUES ($1, $2, $3)"
	result, err := r.DB.Exec(query, category.Name, category.IDClubParent, category.Category)
	if err != nil {
		utils.LogMessage(utils.LevelError, fmt.Sprintf("Failed to create category: %v", err))
		return false, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		utils.LogMessage(utils.LevelError, fmt.Sprintf("Failed to get rows affected: %v", err))
		return false, err
	}

	if rowsAffected > 0 {
		utils.LogMessage(utils.LevelInfo, "Category created successfully")
		return true, nil
	}

	utils.LogMessage(utils.LevelWarn, "No rows affected while creating category")
	return false, nil
}

func (r *ReservationRepository) CreateItem(item models.ReservationCreateItemRequest) (bool, error) {
	if item.Name == "" {
		utils.LogMessage(utils.LevelError, "Item name is required")
		return false, fmt.Errorf("item name is required")
	}

	query := "INSERT INTO reservation_element (name, slot, description, location, id_clubs, id_reservation_category) VALUES ($1, $2, $3, $4, $5, $6)"
	result, err := r.DB.Exec(query, item.Name, item.Slot, item.Description, item.Location, item.IDClub, item.CategoryParent)
	if err != nil {
		utils.LogMessage(utils.LevelError, fmt.Sprintf("Failed to create item: %v", err))
		return false, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		utils.LogMessage(utils.LevelError, fmt.Sprintf("Failed to get rows affected: %v", err))
		return false, err
	}

	if rowsAffected > 0 {
		utils.LogMessage(utils.LevelInfo, "Item created successfully")
		return true, nil
	}

	utils.LogMessage(utils.LevelWarn, "No rows affected while creating item")
	return false, nil
}
