package restaurant

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode"

	"github.com/plugimt/transat-backend/models"
	"github.com/plugimt/transat-backend/services"
	"github.com/plugimt/transat-backend/utils" // logger

	"github.com/gofiber/fiber/v2"
)

type RestaurantHandler struct {
	DB           *sql.DB
	TransService *services.TranslationService
	NotifService *services.NotificationService

	cacheMutex    sync.RWMutex
	lastFetchTime time.Time
	menuSourceURL string
	apiRegex      *regexp.Regexp

	lastNotificationDate    string
	menuSimilarityThreshold float64
}

func NewRestaurantHandler(db *sql.DB, transService *services.TranslationService, notifService *services.NotificationService) *RestaurantHandler {
	// Initialize the regex to match the loadingData array in the API response
	regex := regexp.MustCompile(`var loadingData\s*=\s*(\[\[.*?\]\]);?`)
	sourceURL := "https://toast-js.ew.r.appspot.com/coteresto?key=1ohdRUdCYo6e71aLuBh7ZfF2lc_uZqp9D78icU4DPufA"

	return &RestaurantHandler{
		DB:                      db,
		TransService:            transService,
		NotifService:            notifService,
		menuSourceURL:           sourceURL,
		apiRegex:                regex,
		menuSimilarityThreshold: 0.7,
	}
}

func (h *RestaurantHandler) GetRestaurantMenu(c *fiber.Ctx) error {
	utils.LogHeader("🍽️ Get Restaurant Menu")

	language := c.Query("language")

	if language == "" {
		utils.LogMessage(utils.LevelWarn, "No language specified, defaulting to French")
		language = "fr"
	}

	// Get today's menu categorized by type
	menuResponse, err := h.GetTodaysMenuCategorized()
	if err != nil {
		utils.LogMessage(utils.LevelError, fmt.Sprintf("Failed to get today's menu: %v", err))
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve menu",
		})
	}

	totalItems := len(menuResponse.GrilladesMidi) + len(menuResponse.Migrateurs) + len(menuResponse.Cibo) +
		len(menuResponse.AccompMidi) + len(menuResponse.GrilladesSoir) + len(menuResponse.AccompSoir)

	// If no menu items exist for today, return 204 No Content
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

	// Get dish details with ratings and statistics
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

	err = h.DB.QueryRow(query, dishID).Scan(
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
			utils.LogMessage(utils.LevelError, fmt.Sprintf("Dish not found: ID %d", dishID))
			utils.LogFooter()
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Dish not found",
			})
		}
		utils.LogMessage(utils.LevelError, fmt.Sprintf("Database error: %v", err))
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve dish details",
		})
	}

	// Capitalize the dish name
	dishDetails.Name = h.capitalizeItemName(dishDetails.Name)

	// Get all reviews ordered by date (most recent first)
	reviewsQuery := `
		SELECT n.first_name, n.last_name, n.profile_picture, ran.note, ran.comment, ran.date
		FROM restaurant_articles_notes ran
		JOIN newf n ON ran.email = n.email
		WHERE ran.id_restaurant_articles = $1
		ORDER BY ran.date DESC
	`

	reviewRows, err := h.DB.Query(reviewsQuery, dishID)
	if err != nil {
		utils.LogMessage(utils.LevelWarn, fmt.Sprintf("Failed to fetch reviews: %v", err))
	} else {
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

		utils.LogMessage(utils.LevelInfo, fmt.Sprintf("Retrieved details for dish: %s", dishDetails.Name))
		utils.LogFooter()
		return c.JSON(response)
	}

	utils.LogMessage(utils.LevelInfo, fmt.Sprintf("Retrieved details for dish: %s", dishDetails.Name))
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

	// Check if dish exists
	var dishName string
	err = h.DB.QueryRow("SELECT name FROM restaurant_articles WHERE id_restaurant_articles = $1", dishID).Scan(&dishName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.LogMessage(utils.LevelError, fmt.Sprintf("Dish not found: ID %d", dishID))
			utils.LogFooter()
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Dish not found",
			})
		}
		utils.LogMessage(utils.LevelError, fmt.Sprintf("Database error checking dish: %v", err))
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to verify dish",
		})
	}

	// Insert or update the review
	_, err = h.DB.Exec(`
		INSERT INTO restaurant_articles_notes (email, id_restaurant_articles, note, comment)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (email, id_restaurant_articles) 
		DO UPDATE SET 
			note = EXCLUDED.note,
			comment = EXCLUDED.comment
	`, userEmail, dishID, reviewRequest.Rating, reviewRequest.Comment)

	if err != nil {
		utils.LogMessage(utils.LevelError, fmt.Sprintf("Failed to save review: %v", err))
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save review",
		})
	}

	// Get updated average rating
	var newAverageRating float64
	var totalRatings int
	err = h.DB.QueryRow(`
		SELECT COALESCE(AVG(note), 0), COUNT(note)
		FROM restaurant_articles_notes 
		WHERE id_restaurant_articles = $1
	`, dishID).Scan(&newAverageRating, &totalRatings)

	if err != nil {
		utils.LogMessage(utils.LevelWarn, fmt.Sprintf("Failed to calculate new average: %v", err))
	}

	utils.LogMessage(utils.LevelInfo, fmt.Sprintf("Review saved for dish %s by %s: %d stars", dishName, userEmail, reviewRequest.Rating))
	utils.LogFooter()

	return c.JSON(fiber.Map{
		"message":        "Review saved successfully",
		"dish_name":      h.capitalizeItemName(dishName),
		"average_rating": newAverageRating,
		"total_ratings":  totalRatings,
		"your_rating":    reviewRequest.Rating,
		"your_comment":   reviewRequest.Comment,
	})
}

func (h *RestaurantHandler) fetchMenuFromAPI() (*models.MenuData, error) {
	resp, err := http.Get(h.menuSourceURL)
	if err != nil {
		return nil, fmt.Errorf("unable to fetch URL '%s': %w", h.menuSourceURL, err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("Error closing response body: %v", err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code %d from '%s'", resp.StatusCode, h.menuSourceURL)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read response body: %w", err)
	}

	matches := h.apiRegex.FindSubmatch(body)
	if len(matches) < 2 {
		log.Printf("Body sample for regex failure: %s", string(body[:200])) // Log beginning of body
		return nil, fmt.Errorf("regex did not find 'loadingData' array in response")
	}
	loadingDataJSON := string(matches[1])
	return h.parseLoadingData(loadingDataJSON)
}

func (h *RestaurantHandler) parseLoadingData(jsonData string) (*models.MenuData, error) {
	var nestedItems [][]models.MenuItemAPI
	if err := json.Unmarshal([]byte(jsonData), &nestedItems); err != nil {
		log.Printf("Invalid JSON structure from API: %v", err)
		return nil, fmt.Errorf("unable to parse menu JSON from API: %w", err)
	}

	if len(nestedItems) == 0 || len(nestedItems[0]) == 0 {
		return nil, fmt.Errorf("no menu items found in the parsed API data")
	}

	return h.processRawMenuItems(nestedItems[0]), nil
}

func (h *RestaurantHandler) processRawMenuItems(items []models.MenuItemAPI) *models.MenuData {
	menuData := models.MenuData{
		GrilladesMidi: []string{},
		Migrateurs:    []string{},
		Cibo:          []string{},
		AccompMidi:    []string{},
		GrilladesSoir: []string{},
		AccompSoir:    []string{},
	}
	itemMap := make(map[string]map[string]bool)

	for _, item := range items {
		category := getMenuCategory(item.Pole, item.Accompagnement, item.Periode)
		if category == "" {
			continue // Skip items that don't map to a known category
		}

		menuItemText := strings.TrimSpace(fmt.Sprintf("%s %s %s", item.Nom, item.Info1, item.Info2))
		menuItemText = strings.Join(strings.Fields(menuItemText), " ")
		if menuItemText == "" {
			continue // Skip empty items
		}

		if itemMap[category] == nil {
			itemMap[category] = make(map[string]bool)
		}

		if !itemMap[category][menuItemText] {
			itemMap[category][menuItemText] = true
			switch category {
			case "grilladesMidi":
				menuData.GrilladesMidi = append(menuData.GrilladesMidi, menuItemText)
			case "migrateurs":
				menuData.Migrateurs = append(menuData.Migrateurs, menuItemText)
			case "cibo":
				menuData.Cibo = append(menuData.Cibo, menuItemText)
			case "accompMidi":
				menuData.AccompMidi = append(menuData.AccompMidi, menuItemText)
			case "accompSoir":
				menuData.AccompSoir = append(menuData.AccompSoir, menuItemText)
			case "grilladesSoir":
				menuData.GrilladesSoir = append(menuData.GrilladesSoir, menuItemText)
			}
		}
	}

	return &menuData
}

func (h *RestaurantHandler) CheckAndUpdateMenuCron() (bool, error) {
	utils.LogHeader("🤖 Cron: Check & Update Restaurant Menu")

	// 1. Fetch latest base menu from API
	baseMenuData, err := h.fetchMenuFromAPI()
	if err != nil {
		utils.LogMessage(utils.LevelError, "Cron: Failed to fetch base menu from API")
		utils.LogMessage(utils.LevelError, fmt.Sprintf("Error: %v", err))
		utils.LogFooter()
		return false, err // Don't proceed if API fetch fails
	}

	// 2. Convert to FetchedItems format
	fetchedItems := h.convertMenuDataToFetchedItems(baseMenuData)
	utils.LogMessage(utils.LevelInfo, fmt.Sprintf("Converted %d menu items for synchronization", len(fetchedItems)))

	// 3. Synchronize with database
	err = h.SyncTodaysMenu(fetchedItems)
	if err != nil {
		utils.LogMessage(utils.LevelError, fmt.Sprintf("Failed to sync menu: %v", err))
		utils.LogFooter()
		return false, err
	}

	// 4. Check if notifications should be triggered
	if h.isNotificationTimeAllowed(time.Now()) && len(fetchedItems) > 0 {
		utils.LogMessage(utils.LevelInfo, "Menu updated successfully - notifications could be triggered here")
		// TODO: Implement notification logic if needed
	}

	utils.LogMessage(utils.LevelInfo, "Menu synchronization completed successfully")
	utils.LogFooter()
	return true, nil
}

func (h *RestaurantHandler) isNotificationTimeAllowed(t time.Time) bool {
	weekday := t.Weekday()
	if weekday == time.Saturday || weekday == time.Sunday {
		return false
	}

	hour := t.Hour()
	return hour >= 7 && hour < 16
}

func getMenuCategory(pole string, accompagnement string, periode string) string {
	isAccomp := strings.EqualFold(accompagnement, "TRUE")

	if isAccomp {
		if periode == "midi" {
			return "accompMidi"
		}
		return "accompSoir"
	}

	switch pole {
	case "Grillades / Plats traditions":
		if periode == "midi" {
			return "grilladesMidi"
		}
		return "grilladesSoir"
	case "Les Cuistots migrateurs":
		return "migrateurs"
	case "Le Végétarien":
		return "cibo"
	default:
		return ""
	}
}

// convertMenuDataToFetchedItems converts the old MenuData format to FetchedItem slice
func (h *RestaurantHandler) convertMenuDataToFetchedItems(menuData *models.MenuData) []models.FetchedItem {
	var fetchedItems []models.FetchedItem

	// Map category names to restaurant IDs based on the database inserts
	categoryToID := map[string]int{
		"grilladesMidi": 1,
		"migrateurs":    2,
		"cibo":          3,
		"accompMidi":    4,
		"grilladesSoir": 5,
		"accompSoir":    6,
	}

	// Process each category
	for category, items := range map[string][]string{
		"grilladesMidi": menuData.GrilladesMidi,
		"migrateurs":    menuData.Migrateurs,
		"cibo":          menuData.Cibo,
		"accompMidi":    menuData.AccompMidi,
		"grilladesSoir": menuData.GrilladesSoir,
		"accompSoir":    menuData.AccompSoir,
	} {
		menuTypeID := categoryToID[category]
		for _, item := range items {
			if item != "" {
				fetchedItems = append(fetchedItems, models.FetchedItem{
					Name:       item,
					Category:   category,
					MenuTypeID: menuTypeID,
				})
			}
		}
	}

	return fetchedItems
}

// normalizeItemName normalizes menu item names for consistent comparison
func (h *RestaurantHandler) normalizeItemName(name string) string {
	// Trim whitespace
	normalized := strings.TrimSpace(name)

	// Convert to lowercase
	normalized = strings.ToLower(normalized)

	// Remove extra whitespace between words
	normalized = strings.Join(strings.Fields(normalized), " ")

	// Remove common punctuation but keep essential characters
	result := strings.Builder{}
	for _, r := range normalized {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || unicode.IsSpace(r) {
			result.WriteRune(r)
		}
	}

	return strings.TrimSpace(result.String())
}

// capitalizeItemName capitalizes only the first character of the item name
func (h *RestaurantHandler) capitalizeItemName(name string) string {
	if name == "" {
		return name
	}

	// Convert to runes to handle Unicode characters properly
	runes := []rune(name)
	if len(runes) > 0 {
		runes[0] = unicode.ToUpper(runes[0])
	}

	return string(runes)
}

// calculateSimilarity calculates Jaccard similarity between two sets of strings
func (h *RestaurantHandler) calculateSimilarity(set1, set2 []string) float64 {
	if len(set1) == 0 && len(set2) == 0 {
		return 1.0 // Both empty sets are identical
	}

	// Convert slices to maps for set operations
	map1 := make(map[string]bool)
	map2 := make(map[string]bool)

	for _, item := range set1 {
		map1[item] = true
	}

	for _, item := range set2 {
		map2[item] = true
	}

	// Calculate intersection
	intersection := 0
	for item := range map1 {
		if map2[item] {
			intersection++
		}
	}

	// Calculate union
	union := len(map1) + len(map2) - intersection

	if union == 0 {
		return 1.0
	}

	return float64(intersection) / float64(union)
}

// SyncTodaysMenu synchronizes today's menu with the database
func (h *RestaurantHandler) SyncTodaysMenu(fetchedItems []models.FetchedItem) error {
	utils.LogHeader("🔄 Syncing Today's Menu")

	today := time.Now().Format("2006-01-02")

	// Step 1: Stale Menu Guard Clause
	utils.LogMessage(utils.LevelInfo, "Checking for stale menu...")

	// Get normalized names from fetched items
	var fetchedNames []string
	for _, item := range fetchedItems {
		normalized := h.normalizeItemName(item.Name)
		if normalized != "" {
			fetchedNames = append(fetchedNames, normalized)
		}
	}

	var latestMenuNames []string
	var latestMenuDate sql.NullString

	err := h.DB.QueryRow(`
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

		rows, err := h.DB.Query(query, latestMenuDate.String)
		if err != nil {
			utils.LogMessage(utils.LevelError, fmt.Sprintf("Failed to query latest menu: %v", err))
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var name string
			if err := rows.Scan(&name); err != nil {
				utils.LogMessage(utils.LevelError, fmt.Sprintf("Failed to scan menu item: %v", err))
				continue
			}
			latestMenuNames = append(latestMenuNames, name)
		}
	}

	similarity := h.calculateSimilarity(fetchedNames, latestMenuNames)
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

	if similarity > 0.8 {
		utils.LogMessage(utils.LevelWarn, "Stale menu detected (>80% similarity), skipping update")
		return nil
	}

	// Step 2: Process Fetched Menu Items
	utils.LogMessage(utils.LevelInfo, "Processing fetched menu items...")

	var processedArticleIDs []int

	for _, item := range fetchedItems {
		normalizedName := h.normalizeItemName(item.Name)
		if normalizedName == "" {
			continue
		}

		// Upsert into restaurant_articles and update last_time_served
		var articleID int
		err := h.DB.QueryRow(`
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

		// Insert/Update restaurant_meals for today
		_, err = h.DB.Exec(`
			INSERT INTO restaurant_meals (id_restaurant, id_restaurant_articles, date_served)
			VALUES ($1, $2, $3)
			ON CONFLICT (id_restaurant, id_restaurant_articles, date_served) DO NOTHING
		`, item.MenuTypeID, articleID, today)

		if err != nil {
			utils.LogMessage(utils.LevelError, fmt.Sprintf("Failed to insert meal for article ID %d: %v", articleID, err))
		}
	}

	// Step 3: Remove Obsolete Items
	utils.LogMessage(utils.LevelInfo, "Removing obsolete menu items for today...")

	if len(processedArticleIDs) > 0 {
		// Build the NOT IN clause
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

		result, err := h.DB.Exec(deleteQuery, args...)
		if err != nil {
			utils.LogMessage(utils.LevelError, fmt.Sprintf("Failed to remove obsolete items: %v", err))
		} else {
			rowsAffected, _ := result.RowsAffected()
			utils.LogMessage(utils.LevelInfo, fmt.Sprintf("Removed %d obsolete menu items", rowsAffected))
		}
	} else {
		// If no items were processed, remove all items for today
		result, err := h.DB.Exec(`DELETE FROM restaurant_meals WHERE date_served = $1`, today)
		if err != nil {
			utils.LogMessage(utils.LevelError, fmt.Sprintf("Failed to clear today's menu: %v", err))
		} else {
			rowsAffected, _ := result.RowsAffected()
			utils.LogMessage(utils.LevelInfo, fmt.Sprintf("Cleared all %d items from today's menu", rowsAffected))
		}
	}

	utils.LogMessage(utils.LevelInfo, fmt.Sprintf("Successfully synchronized menu with %d items", len(processedArticleIDs)))
	utils.LogFooter()

	return nil
}

// GetTodaysMenuWithRatings retrieves today's complete menu with average ratings
func (h *RestaurantHandler) GetTodaysMenuWithRatings() (*models.MenuResponse, error) {
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

	rows, err := h.DB.Query(query, today)
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
		// Capitalize the name for display
		entry.Name = h.capitalizeItemName(entry.Name)
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
func (h *RestaurantHandler) GetTodaysMenuCategorized() (*models.CategorizedMenuResponse, error) {
	today := time.Now().Format("2006-01-02")

	// First check if there's any menu for today
	var menuCount int
	countQuery := `SELECT COUNT(*) FROM restaurant_meals WHERE date_served = $1`
	err := h.DB.QueryRow(countQuery, today).Scan(&menuCount)
	if err != nil {
		return nil, fmt.Errorf("failed to check today's menu count: %w", err)
	}

	// If no menu exists for today, return empty response
	if menuCount == 0 {
		return &models.CategorizedMenuResponse{
			GrilladesMidi: []models.MenuItemWithRating{},
			Migrateurs:    []models.MenuItemWithRating{},
			Cibo:          []models.MenuItemWithRating{},
			AccompMidi:    []models.MenuItemWithRating{},
			GrilladesSoir: []models.MenuItemWithRating{},
			AccompSoir:    []models.MenuItemWithRating{},
			UpdatedDate:   time.Now().Format(time.RFC3339),
		}, nil
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

	rows, err := h.DB.Query(query, today)
	if err != nil {
		return nil, fmt.Errorf("failed to query today's menu: %w", err)
	}
	defer rows.Close()

	// Initialize the categorized response
	response := &models.CategorizedMenuResponse{
		GrilladesMidi: []models.MenuItemWithRating{},
		Migrateurs:    []models.MenuItemWithRating{},
		Cibo:          []models.MenuItemWithRating{},
		AccompMidi:    []models.MenuItemWithRating{},
		GrilladesSoir: []models.MenuItemWithRating{},
		AccompSoir:    []models.MenuItemWithRating{},
		UpdatedDate:   time.Now().Format(time.RFC3339),
	}

	// Map menu type IDs to category names
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
			Name:          h.capitalizeItemName(entry.Name),
			AverageRating: entry.AverageRating,
		}

		// Add to appropriate category based on menu type ID
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
