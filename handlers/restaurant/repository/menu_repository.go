package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/plugimt/transat-backend/handlers/restaurant/internal"
	"github.com/plugimt/transat-backend/models"
	"github.com/plugimt/transat-backend/services"
	"github.com/plugimt/transat-backend/utils"
)

type MenuRepository struct {
	DB           *sql.DB
	NotifService *services.NotificationService
}

func NewMenuRepository(db *sql.DB, notifService *services.NotificationService) *MenuRepository {
	return &MenuRepository{
		DB:           db,
		NotifService: notifService,
	}
}

// GetDishDetails retrieves detailed information about a specific dish
func (r *MenuRepository) GetDishDetails(dishID int) (interface{}, error) {
	query := `
		SELECT 
			ra.id_restaurant_articles,
			ra.name,
			ra.first_time_served,
			ra.last_time_served,
			COALESCE(AVG(ran.note), 0) as average_rating,
			COUNT(ran.note) as total_ratings,
			COUNT(DISTINCT rm.date_served) as times_served
		FROM restaurant_articles ra
		LEFT JOIN restaurant_articles_notes ran ON ra.id_restaurant_articles = ran.id_restaurant_articles
		LEFT JOIN restaurant_meals rm ON ra.id_restaurant_articles = rm.id_restaurant_articles
		WHERE ra.id_restaurant_articles = $1
		GROUP BY ra.id_restaurant_articles, ra.name, ra.first_time_served, ra.last_time_served
	`

	var dishDetails struct {
		ID              int        `json:"id" db:"id_restaurant_articles"`
		Name            string     `json:"name" db:"name"`
		FirstTimeServed time.Time  `json:"first_time_served" db:"first_time_served"`
		LastTimeServed  *time.Time `json:"last_time_served,omitempty" db:"last_time_served"`
		AverageRating   float64    `json:"average_rating" db:"average_rating"`
		TotalRatings    int        `json:"total_ratings" db:"total_ratings"`
		TimesServed     int        `json:"times_served" db:"times_served"`
	}

	err := r.DB.QueryRow(query, dishID).Scan(
		&dishDetails.ID,
		&dishDetails.Name,
		&dishDetails.FirstTimeServed,
		&dishDetails.LastTimeServed,
		&dishDetails.AverageRating,
		&dishDetails.TotalRatings,
		&dishDetails.TimesServed,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	dishDetails.Name = internal.CapitalizeItemName(dishDetails.Name)

	// Get all reviews ordered by date (most recent first)
	reviewsQuery := `
		SELECT n.first_name, n.last_name, n.profile_picture, ran.note, ran.comment, ran.date
		FROM restaurant_articles_notes ran
		JOIN newf n ON ran.email = n.email
		WHERE ran.id_restaurant_articles = $1
		ORDER BY ran.date DESC
	`

	reviewRows, err := r.DB.Query(reviewsQuery, dishID)
	if err != nil {
		utils.LogMessage(utils.LevelWarn, fmt.Sprintf("Failed to fetch reviews: %v", err))
		return dishDetails, nil
	}
	defer func(reviewRows *sql.Rows) {
		err := reviewRows.Close()
		if err != nil {
			utils.LogMessage(utils.LevelWarn, fmt.Sprintf("Failed to close review rows: %v", err))
		}
	}(reviewRows)

	var reviews []models.ReviewResponse

	for reviewRows.Next() {
		var review models.ReviewResponse

		err := reviewRows.Scan(&review.FirstName, &review.LastName, &review.ProfilePicture, &review.Rating, &review.Comment, &review.Date)
		if err != nil {
			utils.LogMessage(utils.LevelWarn, fmt.Sprintf("Failed to scan review: %v", err))
			continue
		}
		reviews = append(reviews, review)
	}

	// Create response with reviews
	response := struct {
		models.RestaurantArticle
		AverageRating float64                 `json:"average_rating"`
		TotalRatings  int                     `json:"total_ratings"`
		TimesServed   int                     `json:"times_served"`
		Reviews       []models.ReviewResponse `json:"recent_reviews"`
	}{
		RestaurantArticle: models.RestaurantArticle{
			ID:              dishDetails.ID,
			Name:            dishDetails.Name,
			FirstTimeServed: dishDetails.FirstTimeServed,
			LastTimeServed:  dishDetails.LastTimeServed,
		},
		AverageRating: dishDetails.AverageRating,
		TotalRatings:  dishDetails.TotalRatings,
		TimesServed:   dishDetails.TimesServed,
		Reviews:       reviews,
	}

	return response, nil
}

// shouldSendMenuNotification checks if a notification has already been sent today
func (r *MenuRepository) shouldSendMenuNotification(today string) (bool, error) {
	var notificationSent bool
	query := `SELECT notification_sent FROM restaurant_notifications WHERE DATE(date) = $1 ORDER BY date DESC LIMIT 1`

	err := r.DB.QueryRow(query, today).Scan(&notificationSent)
	if err != nil {
		if err == sql.ErrNoRows {
			return true, nil
		}
		return false, err
	}

	return !notificationSent, nil
}

// sendMenuUpdateNotification sends a notification for menu updates and records it
func (r *MenuRepository) sendMenuUpdateNotification(today string) error {
	// Get subscribers to RESTAURANT service
	subscribers, err := r.NotifService.GetSubscribedUsers("RESTAURANT")
	if err != nil {
		return fmt.Errorf("failed to get restaurant subscribers: %w", err)
	}

	if len(subscribers) == 0 {
		utils.LogMessage(utils.LevelInfo, "No users subscribed to RESTAURANT notifications")
		return nil
	}

	var tokens []string
	for _, sub := range subscribers {
		if sub.NotificationToken != "" {
			tokens = append(tokens, sub.NotificationToken)
		}
	}

	if len(tokens) == 0 {
		utils.LogMessage(utils.LevelInfo, "No valid notification tokens found")
		return nil
	}

	payload := models.NotificationPayload{
		NotificationTokens: tokens,
		Title:              "NOUVEAU MENU V2",
		Message:            "🤪 On peut laisser des review bientôt. Cette notif est un test merci de l'ignorer",
		Sound:              "default",
		ChannelID:          "default",
		Data: map[string]interface{}{
			"screen": "Restaurant",
		},
	}

	err = r.NotifService.SendPushNotification(payload)
	if err != nil {
		return fmt.Errorf("failed to send push notification: %w", err)
	}

	_, err = r.DB.Exec(`
		INSERT INTO restaurant_notifications (date, notification_sent) 
		VALUES (CURRENT_TIMESTAMP, true)
	`)
	if err != nil {
		return fmt.Errorf("failed to record notification: %w", err)
	}

	utils.LogMessage(utils.LevelInfo, fmt.Sprintf("Sent menu notification to %d users", len(tokens)))
	return nil
}

// SaveDishReview saves a dish review and returns the updated rating info
func (r *MenuRepository) SaveDishReview(dishID int, userEmail string, rating int, comment string) (*models.ReviewResult, error) {
	var dishName string
	err := r.DB.QueryRow("SELECT name FROM restaurant_articles WHERE id_restaurant_articles = $1", dishID).Scan(&dishName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("dish not found")
		}
		return nil, fmt.Errorf("database error checking dish: %w", err)
	}

	// Insert or update the review
	_, err = r.DB.Exec(`
		INSERT INTO restaurant_articles_notes (email, id_restaurant_articles, note, comment)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (email, id_restaurant_articles) 
		DO UPDATE SET 
			note = EXCLUDED.note,
			comment = EXCLUDED.comment
	`, userEmail, dishID, rating, comment)

	if err != nil {
		return nil, fmt.Errorf("failed to save review: %w", err)
	}

	var newAverageRating float64
	var totalRatings int
	err = r.DB.QueryRow(`
		SELECT COALESCE(AVG(note), 0), COUNT(note)
		FROM restaurant_articles_notes 
		WHERE id_restaurant_articles = $1
	`, dishID).Scan(&newAverageRating, &totalRatings)

	if err != nil {
		utils.LogMessage(utils.LevelWarn, fmt.Sprintf("Failed to calculate new average: %v", err))
	}

	return &models.ReviewResult{
		Message:       "Review saved successfully",
		DishName:      internal.CapitalizeItemName(dishName),
		AverageRating: newAverageRating,
		TotalRatings:  totalRatings,
		YourRating:    rating,
		YourComment:   comment,
	}, nil
}

// SyncTodaysMenu synchronizes today's menu with the database
func (r *MenuRepository) SyncTodaysMenu(fetchedItems []models.FetchedItem) error {
	utils.LogHeader("🔄 Syncing Today's Menu")
	today := time.Now().Format("2006-01-02")
	utils.LogMessage(utils.LevelInfo, "Checking for stale menu...")

	var fetchedNames []string
	for _, item := range fetchedItems {
		normalized := internal.NormalizeItemName(item.Name)
		if normalized != "" {
			fetchedNames = append(fetchedNames, normalized)
		}
	}

	var latestMenuNames []string
	var latestMenuDate sql.NullString

	err := r.DB.QueryRow(`
		SELECT MAX(date_served) 
		FROM restaurant_meals
	`).Scan(&latestMenuDate)

	if err != nil {
		utils.LogMessage(utils.LevelError, fmt.Sprintf("Failed to query latest menu date: %v", err))
		return err
	}

	if !latestMenuDate.Valid {
		utils.LogMessage(utils.LevelInfo, "No existing menu found, proceeding with sync")
	} else {
		query := `
			SELECT DISTINCT ra.name 
			FROM restaurant_meals rm 
			JOIN restaurant_articles ra ON rm.id_restaurant_articles = ra.id_restaurant_articles 
			WHERE rm.date_served = $1
		`

		rows, err := r.DB.Query(query, latestMenuDate.String)
		if err != nil {
			utils.LogMessage(utils.LevelError, fmt.Sprintf("Failed to query latest menu: %v", err))
			return err
		}
		defer func(rows *sql.Rows) {
			err := rows.Close()
			if err != nil {
				utils.LogMessage(utils.LevelWarn, fmt.Sprintf("Failed to close rows: %v", err))
			}
		}(rows)

		for rows.Next() {
			var name string
			if err := rows.Scan(&name); err != nil {
				utils.LogMessage(utils.LevelError, fmt.Sprintf("Failed to scan menu item: %v", err))
				continue
			}
			latestMenuNames = append(latestMenuNames, name)
		}
	}

	similarity := internal.CalculateSimilarity(fetchedNames, latestMenuNames)
	latestDateStr := "none"
	if latestMenuDate.Valid {
		latestDateStr = latestMenuDate.String
	}
	utils.LogMessage(utils.LevelInfo, fmt.Sprintf("Menu similarity with latest (%s): %.2f%% (fetched: %d items, latest: %d items)",
		latestDateStr, similarity*100, len(fetchedNames), len(latestMenuNames)))

	if len(fetchedNames) > 0 && len(latestMenuNames) > 0 {
		utils.LogMessage(utils.LevelInfo, fmt.Sprintf("Fetched items: %v", fetchedNames))
		utils.LogMessage(utils.LevelInfo, fmt.Sprintf("Latest items: %v", latestMenuNames))
	}

	if similarity == 1.0 {
		utils.LogMessage(utils.LevelWarn, "Similar menu detected (=100% similarity), skipping update")
		return nil
	}

	utils.LogMessage(utils.LevelInfo, "Processing fetched menu items...")

	var processedArticleIDs []int

	for _, item := range fetchedItems {
		normalizedName := internal.NormalizeItemName(item.Name)
		if normalizedName == "" {
			continue
		}

		var articleID int
		err := r.DB.QueryRow(`
			INSERT INTO restaurant_articles (name, last_time_served)
			VALUES ($1, CURRENT_DATE)
			ON CONFLICT (name) DO UPDATE SET 
				last_time_served = CURRENT_DATE
			RETURNING id_restaurant_articles
		`, normalizedName).Scan(&articleID)

		if err != nil {
			utils.LogMessage(utils.LevelError, fmt.Sprintf("Failed to upsert article '%s': %v", normalizedName, err))
			continue
		}

		processedArticleIDs = append(processedArticleIDs, articleID)

		_, err = r.DB.Exec(`
			INSERT INTO restaurant_meals (id_restaurant, id_restaurant_articles, date_served)
			VALUES ($1, $2, $3)
			ON CONFLICT (id_restaurant, id_restaurant_articles, date_served) DO NOTHING
		`, item.MenuTypeID, articleID, today)

		if err != nil {
			utils.LogMessage(utils.LevelError, fmt.Sprintf("Failed to insert meal for article ID %d: %v", articleID, err))
		}
	}

	utils.LogMessage(utils.LevelInfo, "Removing obsolete menu items for today...")

	if len(processedArticleIDs) > 0 {
		placeholders := make([]string, len(processedArticleIDs))
		args := make([]interface{}, len(processedArticleIDs)+1)
		args[0] = today

		for i, id := range processedArticleIDs {
			placeholders[i] = fmt.Sprintf("$%d", i+2)
			args[i+1] = id
		}

		deleteQuery := fmt.Sprintf(`
			DELETE FROM restaurant_meals 
			WHERE date_served = $1 
			AND id_restaurant_articles NOT IN (%s)
		`, strings.Join(placeholders, ","))

		result, err := r.DB.Exec(deleteQuery, args...)
		if err != nil {
			utils.LogMessage(utils.LevelError, fmt.Sprintf("Failed to remove obsolete items: %v", err))
		} else {
			rowsAffected, _ := result.RowsAffected()
			utils.LogMessage(utils.LevelInfo, fmt.Sprintf("Removed %d obsolete menu items", rowsAffected))
		}
	} else {
		result, err := r.DB.Exec(`DELETE FROM restaurant_meals WHERE date_served = $1`, today)
		if err != nil {
			utils.LogMessage(utils.LevelError, fmt.Sprintf("Failed to clear today's menu: %v", err))
		} else {
			rowsAffected, _ := result.RowsAffected()
			utils.LogMessage(utils.LevelInfo, fmt.Sprintf("Cleared all %d items from today's menu", rowsAffected))
		}
	}

	utils.LogMessage(utils.LevelInfo, fmt.Sprintf("Successfully synchronized menu with %d items", len(processedArticleIDs)))

	// Check if we should send a notification (similarity < 80% and no notification sent today)
	if similarity < 0.8 && len(processedArticleIDs) > 0 {
		shouldSendNotification, err := r.shouldSendMenuNotification(today)
		if err != nil {
			utils.LogMessage(utils.LevelError, fmt.Sprintf("Failed to check notification status: %v", err))
		} else if shouldSendNotification {
			utils.LogMessage(utils.LevelInfo, fmt.Sprintf("Menu similarity %.2f%% is below 80%%, sending notification", similarity*100))
			err := r.sendMenuUpdateNotification(today)
			if err != nil {
				utils.LogMessage(utils.LevelError, fmt.Sprintf("Failed to send menu notification: %v", err))
			} else {
				utils.LogMessage(utils.LevelInfo, "Menu update notification sent successfully")
			}
		}
	}

	utils.LogFooter()
	return nil
}

// GetTodaysMenuWithRatings retrieves today's complete menu with average ratings
func (r *MenuRepository) GetTodaysMenuWithRatings() (*models.MenuResponse, error) {
	today := time.Now().Format("2006-01-02")

	query := `
		SELECT 
			rm.id_restaurant_articles,
			ra.name,
			rm.id_restaurant,
			COALESCE(AVG(ran.note), 0) as average_rating
		FROM restaurant_meals rm
		JOIN restaurant_articles ra ON rm.id_restaurant_articles = ra.id_restaurant_articles
		LEFT JOIN restaurant_articles_notes ran ON ra.id_restaurant_articles = ran.id_restaurant_articles
		WHERE rm.date_served = $1
		GROUP BY rm.id_restaurant_articles, ra.name, rm.id_restaurant
		ORDER BY rm.id_restaurant, ra.name
	`

	rows, err := r.DB.Query(query, today)
	if err != nil {
		return nil, fmt.Errorf("failed to query today's menu: %w", err)
	}
	defer rows.Close()

	var menuEntries []models.MenuEntry
	for rows.Next() {
		var entry models.MenuEntry
		err := rows.Scan(&entry.ArticleID, &entry.Name, &entry.MenuTypeID, &entry.AverageRating)
		if err != nil {
			return nil, fmt.Errorf("failed to scan menu entry: %w", err)
		}
		entry.Name = internal.CapitalizeItemName(entry.Name)
		menuEntries = append(menuEntries, entry)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating menu rows: %w", err)
	}

	return &models.MenuResponse{
		Date:  today,
		Items: menuEntries,
	}, nil
}

// GetTodaysMenuCategorized retrieves today's menu categorized by menu type
func (r *MenuRepository) GetTodaysMenuCategorized() (*models.CategorizedMenuResponse, error) {
	today := time.Now().Format("2006-01-02")

	var menuCount int
	countQuery := `SELECT COUNT(*) FROM restaurant_meals WHERE date_served = $1`
	err := r.DB.QueryRow(countQuery, today).Scan(&menuCount)
	if err != nil {
		return nil, fmt.Errorf("failed to check today's menu count: %w", err)
	}

	response := &models.CategorizedMenuResponse{
		GrilladesMidi: []models.MenuItemWithRating{},
		Migrateurs:    []models.MenuItemWithRating{},
		Cibo:          []models.MenuItemWithRating{},
		AccompMidi:    []models.MenuItemWithRating{},
		GrilladesSoir: []models.MenuItemWithRating{},
		AccompSoir:    []models.MenuItemWithRating{},
		UpdatedDate:   time.Now().Format(time.RFC3339),
	}

	if menuCount == 0 {
		return response, nil
	}

	query := `
		SELECT 
			rm.id_restaurant_articles,
			ra.name,
			rm.id_restaurant,
			COALESCE(AVG(ran.note), 0) as average_rating
		FROM restaurant_meals rm
		JOIN restaurant_articles ra ON rm.id_restaurant_articles = ra.id_restaurant_articles
		LEFT JOIN restaurant_articles_notes ran ON ra.id_restaurant_articles = ran.id_restaurant_articles
		WHERE rm.date_served = $1
		GROUP BY rm.id_restaurant_articles, ra.name, rm.id_restaurant
		ORDER BY rm.id_restaurant, ra.name
	`

	rows, err := r.DB.Query(query, today)
	if err != nil {
		return nil, fmt.Errorf("failed to query today's menu: %w", err)
	}
	defer rows.Close()

	menuTypeMap := map[int]string{
		1: "grilladesMidi",
		2: "migrateurs",
		3: "cibo",
		4: "accompMidi",
		5: "grilladesSoir",
		6: "accompSoir",
	}

	for rows.Next() {
		var entry models.MenuEntry
		err := rows.Scan(&entry.ArticleID, &entry.Name, &entry.MenuTypeID, &entry.AverageRating)
		if err != nil {
			return nil, fmt.Errorf("failed to scan menu entry: %w", err)
		}

		menuItem := models.MenuItemWithRating{
			ID:            entry.ArticleID,
			Name:          internal.CapitalizeItemName(entry.Name),
			AverageRating: entry.AverageRating,
		}

		switch menuTypeMap[entry.MenuTypeID] {
		case "grilladesMidi":
			response.GrilladesMidi = append(response.GrilladesMidi, menuItem)
		case "migrateurs":
			response.Migrateurs = append(response.Migrateurs, menuItem)
		case "cibo":
			response.Cibo = append(response.Cibo, menuItem)
		case "accompMidi":
			response.AccompMidi = append(response.AccompMidi, menuItem)
		case "grilladesSoir":
			response.GrilladesSoir = append(response.GrilladesSoir, menuItem)
		case "accompSoir":
			response.AccompSoir = append(response.AccompSoir, menuItem)
		}
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating menu rows: %w", err)
	}

	return response, nil
}
