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

	query := `
		INSERT INTO reservation_category 
			(name, id_clubs%s)
		VALUES 
			($1, $2%s)
		RETURNING id_reservation_category
	`

	args := []interface{}{category.Name, category.IDClubParent}
	extraCols, extraVals := "", ""

	if category.IDCategoryParent != nil {
		extraCols = ", id_parent_category"
		extraVals = ", $3"
		args = append(args, *category.IDCategoryParent)
	}

	finalQuery := fmt.Sprintf(query, extraCols, extraVals)
	row := r.DB.QueryRow(finalQuery, args...)

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
			(name, slot, id_clubs%s) 
		VALUES 
			($1, $2, $3%s)
		RETURNING id_reservation_element
	`

	args := []interface{}{item.Name, item.Slot, *item.IDClubParent}
	extraCols, extraVals := "", ""

	if item.IDCategoryParent != nil {
		extraCols = ", id_reservation_category"
		extraVals = ", $4"
		args = append(args, *item.IDCategoryParent)
	}
	if item.Description != nil {
		extraCols += ", description"
		extraVals += ", $5"
		args = append(args, *item.Description)
	}
	if item.Location != nil {
		extraCols += ", location"
		extraVals += ", $6"
		args = append(args, *item.Location)
	}

	finalQuery := fmt.Sprintf(fmt.Sprintf(query, extraCols, extraVals) + ";")

	fmt.Printf("Executing query: %s with args: %v\n", finalQuery, args)

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

func (r *ReservationRepository) GetCategoryList(IDCategoryParent *int) ([]models.ReservationCategory, error) {
	var categories []models.ReservationCategory

	query := "SELECT id_reservation_category, name FROM reservation_category"
	args := []interface{}{}

	if IDCategoryParent != nil {
		query += " WHERE id_parent_category = $1;"
		args = append(args, *IDCategoryParent)
	} else {
		query += " WHERE id_parent_category IS NULL;"
	}

	rows, err := r.DB.Query(query, args...)
	if err != nil {
		utils.LogMessage(utils.LevelError, fmt.Sprintf("Failed to get categories: %v", err))
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			utils.LogMessage(utils.LevelError, fmt.Sprintf("Failed to close rows: %v", err))
		} else {
			utils.LogMessage(utils.LevelInfo, "Rows closed successfully")
		}
	}(rows)

	for rows.Next() {
		var category models.ReservationCategory
		if err := rows.Scan(&category.ID, &category.Name); err != nil {
			utils.LogMessage(utils.LevelError, fmt.Sprintf("Failed to scan category: %v", err))
			return nil, err
		}
		categories = append(categories, category)
	}

	if err := rows.Err(); err != nil {
		utils.LogMessage(utils.LevelError, fmt.Sprintf("Row error: %v", err))
		return nil, err
	}

	return categories, nil
}

func (r *ReservationRepository) GetItemList(IDCategoryParent *int) ([]models.ReservationItem, error) {
	var items []models.ReservationItem

	query := "SELECT id_reservation_element, name, slot FROM reservation_element"
	args := []interface{}{}

	if IDCategoryParent != nil {
		query += " WHERE id_reservation_category = $1;"
		args = append(args, *IDCategoryParent)
	} else {
		query += " WHERE id_reservation_category IS NULL;"
	}

	rows, err := r.DB.Query(query, args...)
	if err != nil {
		utils.LogMessage(utils.LevelError, fmt.Sprintf("Failed to get items: %v", err))
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			utils.LogMessage(utils.LevelError, fmt.Sprintf("Failed to close rows: %v", err))
		} else {
			utils.LogMessage(utils.LevelInfo, "Rows closed successfully")
		}
	}(rows)

	for rows.Next() {
		var item models.ReservationItem
		if err := rows.Scan(&item.ID, &item.Name, &item.Slot); err != nil {
			utils.LogMessage(utils.LevelError, fmt.Sprintf("Failed to scan item: %v", err))
			return nil, err
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		utils.LogMessage(utils.LevelError, fmt.Sprintf("Row error: %v", err))
		return nil, err
	}

	// TODO; add user details if item is reserved

	return items, nil
}
