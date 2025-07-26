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

func (r *ReservationRepository) GetIDClubParent(IDCategory int) (int, error) {
	query := "SELECT id_clubs FROM reservation_category WHERE id_reservation_category = $1"
	row := r.DB.QueryRow(query, IDCategory)
	var idClub int
	if err := row.Scan(&idClub); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.LogMessage(utils.LevelWarn, fmt.Sprintf("No club found for category ID %d", IDCategory))
			return 0, nil
		}
		utils.LogMessage(utils.LevelError, fmt.Sprintf("Failed to get club parent: %v", err))
		return 0, err
	}
	utils.LogMessage(utils.LevelInfo, fmt.Sprintf("Found club ID %d for category ID %d", idClub, IDCategory))
	return idClub, nil
}

func (r *ReservationRepository) CreateCategory(category models.ReservationCreateCategoryRequest) (models.ReservationCategoryComplete, error) {
	var res models.ReservationCategoryComplete
	if category.Name == "" {
		utils.LogMessage(utils.LevelError, "Category name is required")
		return res, fmt.Errorf("category name is required")
	}

	query := "INSERT INTO reservation_category (name, id_clubs, id_parent_category) VALUES ($1, $2, $3) RETURNING id_reservation_category"
	row := r.DB.QueryRow(query, category.Name, category.IDClubParent, category.IDCategoryParent)

	if err := row.Scan(&res.ID); err != nil {
		utils.LogMessage(utils.LevelError, fmt.Sprintf("Failed to create category: %v", err))
		return res, err
	}

	utils.LogMessage(utils.LevelInfo, "Category created successfully")
	res.Name = category.Name
	res.IDClubParent = category.IDClubParent
	res.IDCategoryParent = category.IDCategoryParent
	return res, nil
}

func (r *ReservationRepository) CreateItem(item models.ReservationCreateItemRequest) (models.ReservationItemComplete, error) {
	var res models.ReservationItemComplete

	if item.Name == "" {
		utils.LogMessage(utils.LevelError, "Item name is required")
		return res, fmt.Errorf("item name is required")
	}

	query := `
	INSERT INTO reservation_element 
		(name, slot, description, location, id_clubs%s) 
	VALUES 
		($1, $2, $3, $4, $5%s)
	RETURNING id_reservation_element
`

	args := []interface{}{item.Name, item.Slot, item.Description, item.Location, *item.IDClubParent}
	extraCols, extraVals := "", ""

	if item.IDCategoryParent != nil {
		extraCols = ", id_reservation_category"
		extraVals = ", $6"
		args = append(args, *item.IDCategoryParent)
	}

	finalQuery := fmt.Sprintf(query, extraCols, extraVals)
	row := r.DB.QueryRow(finalQuery, args...)

	if err := row.Scan(&res.ID); err != nil {
		utils.LogMessage(utils.LevelError, fmt.Sprintf("Failed to create item: %v", err))
		return res, err
	}

	utils.LogMessage(utils.LevelInfo, "Item created successfully")
	res.Name = item.Name
	res.Slot = item.Slot
	res.Description = item.Description
	res.Location = item.Location
	res.IDClubParent = item.IDClubParent
	res.IDCategoryParent = item.IDCategoryParent
	return res, nil
}
