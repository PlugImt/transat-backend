package repository

import (
	"database/sql"
	"fmt"
	"github.com/pkg/errors"
	"github.com/plugimt/transat-backend/models"
	"github.com/plugimt/transat-backend/utils"
	"time"
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

func (r *ReservationRepository) CheckItemExists(IDItem int) (bool, error) {
	query := "SELECT COUNT(*) FROM reservation_element WHERE id_reservation_element = $1"
	row := r.DB.QueryRow(query, IDItem)
	var count int
	if err := row.Scan(&count); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.LogMessage(utils.LevelWarn, fmt.Sprintf("Item with ID %d does not exist", IDItem))
			return false, nil
		}
		utils.LogMessage(utils.LevelError, fmt.Sprintf("Failed to check item existence: %v", err))
		return false, err
	}
	if count > 0 {
		utils.LogMessage(utils.LevelInfo, fmt.Sprintf("Item with ID %d exists", IDItem))
		return true, nil
	} else {
		utils.LogMessage(utils.LevelWarn, fmt.Sprintf("Item with ID %d does not exist", IDItem))
	}
	utils.LogMessage(utils.LevelInfo, fmt.Sprintf("Item with ID %d does not exist", IDItem))

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

func (r *ReservationRepository) GetCategoryList(IDCategoryParent *int, ClubID *int) ([]models.ReservationCategory, error) {
	var categories []models.ReservationCategory

	query := "SELECT id_reservation_category, name FROM reservation_category"
	args := []interface{}{}

	if IDCategoryParent != nil {
		query += " WHERE id_parent_category = $1;"
		args = append(args, *IDCategoryParent)
	} else if ClubID != nil {
		query += " WHERE id_clubs = $1 AND id_parent_category IS NULL;"
		args = append(args, *ClubID)
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

func (r *ReservationRepository) GetItemList(IDCategoryParent *int, ClubID *int) ([]models.ReservationItem, error) {
	var items []models.ReservationItem

	query := `
		SELECT
		    re.id_reservation_element,
		    re.name,
		    re.slot,
		    n.email,
		    n.first_name,
		    n.last_name,
		    n.profile_picture
		FROM reservation_element re
		         LEFT JOIN reservation r
		                   ON r.id_reservation_element = re.id_reservation_element
		                       AND re.slot = TRUE
		                       AND r.end_date IS NULL
		         LEFT JOIN newf n
		                   ON n.email = r.email
		
	`
	args := []interface{}{}

	if IDCategoryParent != nil {
		query += " WHERE id_reservation_category = $1;"
		args = append(args, *IDCategoryParent)
	} else if ClubID != nil {
		query += " WHERE re.id_clubs = $1 AND id_reservation_category IS NULL;"
		args = append(args, *ClubID)
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

		var email, firstName, lastName, profilePicture sql.NullString

		if err := rows.Scan(&item.ID, &item.Name, &item.Slot, &email, &firstName, &lastName, &profilePicture); err != nil {
			utils.LogMessage(utils.LevelError, fmt.Sprintf("Failed to scan item: %v", err))
			return nil, err
		}

		if email.Valid {
			item.User = &models.ReservationUser{
				Email:          email.String,
				FirstName:      firstName.String,
				LastName:       lastName.String,
				ProfilePicture: profilePicture.String,
			}
		} else {
			item.User = nil
		}

		items = append(items, item)
	}

	return items, nil
}

func (r *ReservationRepository) GetItemDetails(itemID int, date time.Time) (models.ReservationItemDetailResponse, error) {

	res := models.ReservationItemDetailResponse{
		ID:                itemID,
		Name:              "",    // WIll be filled later
		Slot:              false, // Will be filled later
		Reservation:       []models.ReservationSlotDetail{},
		ReservationBefore: []models.ReservationSlotDetail{},
		ReservationAfter:  []models.ReservationSlotDetail{},
	}

	dateBefore := date.AddDate(0, 0, -1)
	dateAfter := date.AddDate(0, 0, 1)
	dateBeforeStr := dateBefore.Format("2006-01-02")
	dateAfterStr := dateAfter.Format("2006-01-02")
	dateStr := date.Format("2006-01-02")

	query := `
		SELECT re.name,
		       re.slot,
		       n.email,
		       n.first_name,
		       n.last_name,
		       n.profile_picture,
		       r.id_reservation_element,
		       r.start_date,
		       r.end_date
		FROM reservation r
		         JOIN newf n ON r.email = n.email
		         JOIN reservation_element re on r.id_reservation_element = re.id_reservation_element
		WHERE r.id_reservation_element = $1
		  AND DATE(r.start_date) IN (CAST($2 AS DATE), CAST($3 AS DATE), CAST($4 AS DATE));
	`

	args := []interface{}{itemID, dateBeforeStr, dateStr, dateAfterStr}

	rows, err := r.DB.Query(query, args...)
	if err != nil {
		utils.LogMessage(utils.LevelError, fmt.Sprintf("Failed to get item details: %v", err))
		return res, err
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
		var item models.ReservationSlotDetail
		var user models.ReservationUser
		if err := rows.Scan(&res.Name, &res.Slot, &user.Email, &user.FirstName, &user.LastName, &user.ProfilePicture, &item.ID, &item.StartDate, &item.EndDate); err != nil {
			utils.LogMessage(utils.LevelError, fmt.Sprintf("Failed to scan item details: %v", err))
			return res, err
		}
		item.User = user

		parsedStartDate, err := time.Parse(time.RFC3339, item.StartDate)
		if err != nil {
			utils.LogMessage(utils.LevelError, fmt.Sprintf("Failed to parse start date: %v", err))
			continue // skip invalid row
		}

		itemDateStr := parsedStartDate.Format("2006-01-02")

		if itemDateStr == dateStr {
			res.Reservation = append(res.Reservation, item)
		} else if itemDateStr < dateStr {
			res.ReservationBefore = append(res.ReservationBefore, item)
		} else {
			res.ReservationAfter = append(res.ReservationAfter, item)
		}

	}

	if err := rows.Err(); err != nil {
		utils.LogMessage(utils.LevelError, fmt.Sprintf("Row error: %v", err))
		return res, err
	}

	utils.LogMessage(utils.LevelInfo, fmt.Sprintf("Item details retrieved for item ID %d on date %s", itemID, dateStr))
	return res, nil
}

func (r *ReservationRepository) IsItemAvailable(IDItem int, TimeStamp time.Time) (bool, error) {
	query := `
		WITH element_slot AS (
		    SELECT slot
		    FROM reservation_element
		    WHERE id_reservation_element = $1
		),
		     availability_check AS (
		         SELECT
		             CASE
		                 WHEN slot = false THEN (
		                     SELECT COUNT(*) = 0
		                     FROM reservation
		                     WHERE id_reservation_element = $1
		                       AND end_date IS NULL
		                 )
		                 WHEN slot = true THEN (
		                     SELECT COUNT(*) = 0
		                     FROM reservation
		                     WHERE id_reservation_element = $1
		                       AND start_date = $2
		                 )
		                 ELSE true
		                 END AS is_available
		         FROM element_slot
		     )
		SELECT is_available FROM availability_check;
	`
	row := r.DB.QueryRow(query, IDItem, TimeStamp)
	var isAvailable bool
	if err := row.Scan(&isAvailable); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.LogMessage(utils.LevelWarn, fmt.Sprintf("No availability check found for item ID %d at time %s", IDItem, TimeStamp))
			return false, nil
		}
		utils.LogMessage(utils.LevelError, fmt.Sprintf("Failed to check item availability: %v", err))
		return false, err
	}

	if isAvailable {
		utils.LogMessage(utils.LevelInfo, fmt.Sprintf("Item ID %d is available at time %s", IDItem, TimeStamp))
	} else {
		utils.LogMessage(utils.LevelWarn, fmt.Sprintf("Item ID %d is not available at time %s", IDItem, TimeStamp))
	}

	return isAvailable, nil
}

func (r *ReservationRepository) IsItemPerSlot(IDItem int) (bool, error) {
	query := "SELECT slot FROM reservation_element WHERE id_reservation_element = $1"
	row := r.DB.QueryRow(query, IDItem)
	var isSlot bool
	if err := row.Scan(&isSlot); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.LogMessage(utils.LevelWarn, fmt.Sprintf("No item found with ID %d", IDItem))
			return false, nil
		}
		utils.LogMessage(utils.LevelError, fmt.Sprintf("Failed to check if item is per slot: %v", err))
		return false, err
	}

	if isSlot {
		utils.LogMessage(utils.LevelInfo, fmt.Sprintf("Item ID %d is a slot-based item", IDItem))
	} else {
		utils.LogMessage(utils.LevelInfo, fmt.Sprintf("Item ID %d is not a slot-based item", IDItem))
	}

	return isSlot, nil
}

func (r *ReservationRepository) DoesReservationCanBeEnded(ItemID int, UserEmail string) (bool, error) {
	query := `
		SELECT COUNT(*) > 0
		FROM reservation
		WHERE id_reservation_element = $1 AND email = $2 AND end_date IS NULL;
	`
	row := r.DB.QueryRow(query, ItemID, UserEmail)
	var canBeEnded bool
	if err := row.Scan(&canBeEnded); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.LogMessage(utils.LevelWarn, fmt.Sprintf("No reservation found for item ID %d by user %s", ItemID, UserEmail))
			return false, nil
		}
		utils.LogMessage(utils.LevelError, fmt.Sprintf("Failed to check if reservation can be ended: %v", err))
		return false, err
	}

	if canBeEnded {
		utils.LogMessage(utils.LevelInfo, fmt.Sprintf("Reservation for item ID %d by user %s can be ended", ItemID, UserEmail))
	} else {
		utils.LogMessage(utils.LevelInfo, fmt.Sprintf("Reservation for item ID %d by user %s cannot be ended", ItemID, UserEmail))
	}

	return canBeEnded, nil
}

func (r *ReservationRepository) CreateReservation(item models.ReservationManagementRequestTime, IDItem int, ItemPerSlot bool, UserEmail string) (models.ReservationItem, error) {
	var res models.ReservationItem

	isAvailable, err := r.IsItemAvailable(IDItem, *item.StartDate)
	if err != nil {
		utils.LogMessage(utils.LevelError, fmt.Sprintf("Error checking item availability: %v", err))
		return res, err
	}
	if !isAvailable {
		utils.LogMessage(utils.LevelWarn, fmt.Sprintf("Item ID %d is not available at %s", IDItem, *item.StartDate))
		return res, fmt.Errorf("item is not available")
	}

	columns := "email, id_reservation_element, start_date"
	values := "$1, $2, $3"
	args := []interface{}{UserEmail, IDItem, *item.StartDate}

	if ItemPerSlot {
		endDate := item.StartDate.Add(time.Hour)
		item.EndDate = &endDate
		utils.LogMessage(utils.LevelInfo, fmt.Sprintf("End date set to %s for item ID %d", endDate, IDItem))

		columns += ", end_date"
		values += ", $4"
		args = append(args, *item.EndDate)
	}

	query := fmt.Sprintf(`
		WITH inserted AS (
			INSERT INTO reservation (%s)
			VALUES (%s)
			RETURNING email
		)
		SELECT i.email, n.first_name, n.last_name, n.profile_picture
		FROM inserted i
		JOIN newf n ON i.email = n.email;
	`, columns, values)

	res.User = &models.ReservationUser{}

	row := r.DB.QueryRow(query, args...)
	if err := row.Scan(&res.User.Email, &res.User.FirstName, &res.User.LastName, &res.User.ProfilePicture); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.LogMessage(utils.LevelWarn, fmt.Sprintf("No reservation created for item ID %d at %s", IDItem, *item.StartDate))
			return res, nil
		}
		utils.LogMessage(utils.LevelError, fmt.Sprintf("Failed to create reservation: %v", err))
		return res, err
	}

	utils.LogMessage(utils.LevelInfo, fmt.Sprintf("Reservation created for item ID %d at %s", IDItem, *item.StartDate))

	res.ID = IDItem
	res.Name = "" // Optional
	res.Slot = ItemPerSlot

	return res, nil
}

func (r *ReservationRepository) EndReservation(item models.ReservationManagementRequestTime, IDItem int, ItemPerSlot bool, UserEmail string) (models.ReservationItem, error) {
	var res models.ReservationItem

	canBeEnded, err := r.DoesReservationCanBeEnded(IDItem, UserEmail)
	if err != nil {
		utils.LogMessage(utils.LevelError, fmt.Sprintf("Error checking if reservation can be ended: %v", err))
		return res, err
	}
	if !canBeEnded {
		utils.LogMessage(utils.LevelWarn, fmt.Sprintf("Reservation for item ID %d by user %s cannot be ended", IDItem, UserEmail))
		return res, fmt.Errorf("reservation cannot be ended")
	}
	if item.EndDate == nil {
		utils.LogMessage(utils.LevelError, "End date is required to end a reservation")
		return res, fmt.Errorf("end date is required")
	}
	if ItemPerSlot {
		utils.LogMessage(utils.LevelError, "Item is per slot, cannot end reservation")
		return res, fmt.Errorf("item is per slot, cannot end reservation")
	}

	query := `
		WITH updated AS (
			UPDATE reservation
			SET end_date = $1
			WHERE id_reservation_element = $2 AND email = $3 AND end_date IS NULL
			RETURNING email
		)
		SELECT i.email, n.first_name, n.last_name, n.profile_picture
		FROM updated i
		JOIN newf n ON i.email = n.email;
	`

	row := r.DB.QueryRow(query, *item.EndDate, IDItem, UserEmail)
	res.User = &models.ReservationUser{}

	if err := row.Scan(&res.User.Email, &res.User.FirstName, &res.User.LastName, &res.User.ProfilePicture); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.LogMessage(utils.LevelWarn, fmt.Sprintf("No reservation ended for item ID %d at %s", IDItem, *item.EndDate))
			return res, nil
		}
		utils.LogMessage(utils.LevelError, fmt.Sprintf("Failed to end reservation: %v", err))
		return res, err
	}

	utils.LogMessage(utils.LevelInfo, fmt.Sprintf("Reservation ended for item ID %d at %s", IDItem, *item.EndDate))

	res.ID = IDItem
	res.Name = "" // Optional
	res.Slot = ItemPerSlot

	return res, nil
}

func (r *ReservationRepository) DeleteReservation(item models.ReservationManagementRequestTime, IDItem int, UserEmail string) (bool, error) {
	query := `
		DELETE FROM reservation
		WHERE id_reservation_element = $1 AND email = $2 AND start_date = $3
	`

	args := []interface{}{IDItem, UserEmail, *item.StartDate}
	result, err := r.DB.Exec(query, args...)
	if err != nil {
		utils.LogMessage(utils.LevelError, fmt.Sprintf("Failed to delete reservation: %v", err))
		return false, err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		utils.LogMessage(utils.LevelError, fmt.Sprintf("Failed to get rows affected: %v", err))
		return false, err
	}
	if rowsAffected == 0 {
		utils.LogMessage(utils.LevelWarn, fmt.Sprintf("No reservation found for item ID %d by user %s at %s", IDItem, UserEmail, *item.StartDate))
		return false, nil
	}
	utils.LogMessage(utils.LevelInfo, fmt.Sprintf("Reservation deleted for item ID %d by user %s at %s", IDItem, UserEmail, *item.StartDate))
	return true, nil
}
