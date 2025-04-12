package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"

	"Transat_2.0_Backend/models"
	"github.com/gofiber/fiber/v2"
)

var (
	lastFetchTime time.Time
	cachedMenus   map[int]*models.MenuData // Map of language ID to menu data
	cacheMutex    sync.Mutex
)

func init() {
	cachedMenus = make(map[int]*models.MenuData)
}

func fetchMenuFromAPI() (*models.MenuData, error) {
	const targetURL = "https://toast-js.ew.r.appspot.com/coteresto?key=1ohdRUdCYo6e71aLuBh7ZfF2lc_uZqp9D78icU4DPufA"
	regex := regexp.MustCompile(`var loadingData = (\[\[.*?\]\])`)

	resp, err := http.Get(targetURL)
	if err != nil {
		return nil, fmt.Errorf("unable to fetch: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("failed to close response body: %v\n", err)
			return
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read response body: %w", err)
	}

	matches := regex.FindSubmatch(body)
	if len(matches) < 2 {
		return nil, fmt.Errorf("no match found in the response")
	}

	loadingData := string(matches[1])

	var nestedItems [][]models.MenuItem
	if err := json.Unmarshal([]byte(loadingData), &nestedItems); err != nil {
		return nil, fmt.Errorf("unable to parse JSON: %w", err)
	}

	if len(nestedItems) == 0 || len(nestedItems[0]) == 0 {
		return nil, fmt.Errorf("no menu items found in the response")
	}

	menuItems := nestedItems[0]

	menuData := models.MenuData{
		GrilladesMidi: []string{},
		Migrateurs:    []string{},
		Cibo:          []string{},
		AccompMidi:    []string{},
		GrilladesSoir: []string{},
		AccompSoir:    []string{},
	}

	for _, item := range menuItems {
		category := getMenuCategory(item.Pole, item.Accompagnement, item.Periode)
		if category != "" {
			menuItem := strings.TrimSpace(item.Nom + " " + item.Info1 + item.Info2)

			switch category {
			case "grilladesMidi":
				if !contains(menuData.GrilladesMidi, menuItem) {
					menuData.GrilladesMidi = append(menuData.GrilladesMidi, menuItem)
				}
			case "migrateurs":
				if !contains(menuData.Migrateurs, menuItem) {
					menuData.Migrateurs = append(menuData.Migrateurs, menuItem)
				}
			case "cibo":
				if !contains(menuData.Cibo, menuItem) {
					menuData.Cibo = append(menuData.Cibo, menuItem)
				}
			case "accompMidi":
				if !contains(menuData.AccompMidi, menuItem) {
					menuData.AccompMidi = append(menuData.AccompMidi, menuItem)
				}
			case "accompSoir":
				if !contains(menuData.AccompSoir, menuItem) {
					menuData.AccompSoir = append(menuData.AccompSoir, menuItem)
				}
			case "grilladesSoir":
				if !contains(menuData.GrilladesSoir, menuItem) {
					menuData.GrilladesSoir = append(menuData.GrilladesSoir, menuItem)
				}
			}
		}
	}

	return &menuData, nil
}

func getMenuCategory(pole string, accompagnement string, periode string) string {
	if accompagnement == "TRUE" {
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

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func updateRestaurantMenu(db *sql.DB, menuData *models.MenuData) (*models.Restaurant, error) {
	menuJSON, err := json.Marshal(menuData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal menu data: %w", err)
	}

	if len(string(menuJSON)) > 5000 {
		return nil, fmt.Errorf("menu data too large for database column (max 5000 characters)")
	}

	query := `INSERT INTO restaurant (articles) VALUES ($1) RETURNING id_restaurant, articles, updated_date`

	var restaurant models.Restaurant
	err = db.QueryRow(query, string(menuJSON)).Scan(&restaurant.ID, &restaurant.Articles, &restaurant.UpdatedDate)
	if err != nil {
		return nil, fmt.Errorf("failed to update restaurant menu: %w", err)
	}

	return &restaurant, nil
}

func getRestaurant(c *fiber.Ctx) error {
	log.Println("╔======== 🍽️ Get Restaurant 🍽️ ========╗")

	// Parse request body for language
	var req models.RestaurantRequest
	if err := c.BodyParser(&req); err != nil {
		// If no body provided, default to French
		req.Language = "fr"
	}

	// Convert language code to language ID
	langID := getLanguageID(req.Language)
	if langID == 0 {
		log.Println("║ 💥 Invalid language code: ", req.Language)
		log.Println("╚=========================================╝")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid language code",
		})
	}

	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	now := time.Now()
	isSameDay := now.Year() == lastFetchTime.Year() &&
		now.Month() == lastFetchTime.Month() &&
		now.Day() == lastFetchTime.Day()

	// Check if we have a cached menu for the requested language
	if isSameDay && cachedMenus[langID] != nil {
		log.Println("║ ✅ Returning cached menu data for language", req.Language)
		log.Println("╚=========================================╝")
		var menuToReturn models.FullMenuData
		menuToReturn.MenuData = *cachedMenus[langID]
		menuToReturn.UpdatedDate = lastFetchTime.Format("2006-01-02 15:04:05")
		return c.JSON(menuToReturn)
	}

	shouldFetchFromAPI := !isSameDay

	request := `
        SELECT restaurant.id_restaurant, restaurant.articles, restaurant.updated_date, restaurant.language 
        from restaurant
        WHERE restaurant.language = $1
        ORDER BY updated_date DESC
        LIMIT 1;
    `

	stmt, err := db.Prepare(request)
	if err != nil {
		log.Println("║ 💥 Failed to prepare statement: ", err)
		log.Println("╚=========================================╝")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Something went wrong",
		})
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			log.Println("║ 💥 Failed to close statement: ", err)
		}
	}(stmt)

	var restaurant models.Restaurant
	err = stmt.QueryRow(langID).Scan(&restaurant.ID, &restaurant.Articles, &restaurant.UpdatedDate, &restaurant.Language)
	if err != nil {
		log.Println("║ 💥 Failed to get restaurant: ", err)
		if err != sql.ErrNoRows {
			log.Println("╚=========================================╝")
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Something went wrong",
			})
		}

		shouldFetchFromAPI = true
	}

	if !shouldFetchFromAPI && restaurant.ID != 0 {
		var dbMenuData models.MenuData
		if err := json.Unmarshal([]byte(restaurant.Articles), &dbMenuData); err != nil {
			log.Println("║ 💥 Failed to parse database menu: ", err)
			// Need to fetch from API since DB data is invalid
			shouldFetchFromAPI = true
		} else {
			cachedMenus[langID] = &dbMenuData
			log.Println("║ ✅ Returning menu data from database for language", langID)
			log.Println("╚=========================================╝")
			return c.JSON(&dbMenuData)
		}
	}

	log.Println("║ 🔄 Fetching menu from API")
	menuData, err := fetchMenuFromAPI()
	if err != nil {
		log.Println("║ 💥 Failed to fetch menu from API: ", err)
		if restaurant.ID != 0 {
			var dbMenuData models.MenuData
			if err := json.Unmarshal([]byte(restaurant.Articles), &dbMenuData); err != nil {
				log.Println("║ 💥 Failed to parse database menu: ", err)
				log.Println("╚=========================================╝")
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Failed to fetch menu and existing data is invalid",
				})
			}
			log.Println("║ ⚠️ Returning existing menu data from database")
			log.Println("╚=========================================╝")
			return c.JSON(&dbMenuData)
		}
		log.Println("╚=========================================╝")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch menu",
		})
	}

	lastFetchTime = now
	cachedMenus[1] = menuData // Cache French version by default

	// If requested language is not French, translate the menu
	if langID != 1 {
		translationService, err := NewTranslationService()
		if err != nil {
			log.Println("║ 💥 Failed to create translation service: ", err)
			log.Println("╚=========================================╝")
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to initialize translation service",
			})
		}

		translatedMenu, err := translationService.TranslateMenu(menuData, req.Language)
		if err != nil {
			log.Println("║ 💥 Failed to translate menu: ", err)
			log.Println("╚=========================================╝")
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to translate menu",
			})
		}

		cachedMenus[langID] = translatedMenu
		menuData = translatedMenu
	}

	// Update database with translated menu
	menuJSON, err := json.Marshal(menuData)
	if err != nil {
		log.Println("║ 💥 Failed to marshal menu data: ", err)
		log.Println("╚=========================================╝")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to prepare menu data for database",
		})
	}

	updateQuery := `
		WITH check_duplicate AS (
  			SELECT id_restaurant
  			FROM restaurant
  			WHERE articles = $1
    		AND language = $2
    		AND DATE(updated_date) = DATE(CURRENT_TIMESTAMP)
		)
		INSERT INTO restaurant (articles, language)
		SELECT $1, $2
		WHERE NOT EXISTS (SELECT 1 FROM check_duplicate)
		RETURNING id_restaurant, articles, updated_date, language;
	`

	updateStmt, err := db.Prepare(updateQuery)
	if err != nil {
		log.Println("║ 💥 Failed to prepare update statement: ", err)
		log.Println("╚=========================================╝")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to prepare database update",
		})
	}
	defer updateStmt.Close()

	err = updateStmt.QueryRow(string(menuJSON), langID).Scan(&restaurant.ID, &restaurant.Articles, &restaurant.UpdatedDate, &restaurant.Language)
	if err != nil {
		log.Println("║ 💥 Failed to update restaurant menu in database: ", err)
		log.Println("║ ⚠️ Returning fetched menu data despite DB update failure")
		log.Println("╚=========================================╝")
		var menuToReturn models.FullMenuData
		menuToReturn.MenuData = *menuData
		menuToReturn.UpdatedDate = lastFetchTime.Format("2006-01-02 15:04:05")
		return c.JSON(menuToReturn)
	}

	log.Println("║ ✅ Restaurant menu updated successfully")

	if req.Language == "fr" {
		log.Println("║ 🔔 Sending notification to subscribers")
		notificationService := NewNotificationService(db)
		err = notificationService.SendDailyMenuNotification()
		if err != nil {
			log.Println("║ 💥 Failed to send notification: ", err)
			log.Println("╚=========================================╝")
		}
	}

	// Set flag to indicate menu has been updated for today
	menuCheckMutex.Lock()
	if now.Hour() > 12 {
		menuCheckedToday = true
	}
	menuCheckMutex.Unlock()

	log.Println("╚=========================================╝")

	var menuToReturn models.FullMenuData
	menuToReturn.MenuData = *menuData
	menuToReturn.UpdatedDate = lastFetchTime.Format("2006-01-02 15:04:05")
	return c.JSON(menuToReturn)
}

// Helper function to get language ID from language code
func getLanguageID(langCode string) int {
	switch langCode {
	case "fr":
		return 1 // French
	case "en":
		return 2 // English
	case "es":
		return 3 // Spanish
	case "de":
		return 4 // German
	case "it":
		return 5 // Italian
	case "pt":
		return 6 // Portuguese
	case "ru":
		return 7 // Russian
	case "zh":
		return 8 // Chinese
	case "ko":
		return 9 // Korean
	default:
		return 0 // Invalid
	}
}

func checkAndUpdateMenu(notificationService *NotificationService) (bool, error) {
	log.Println("Fetching latest menu from API")

	menuData, err := fetchMenuFromAPI()
	if err != nil {
		return false, fmt.Errorf("failed to fetch menu: %w", err)
	}

	// Get the latest menu from the database for comparison
	request := `
        SELECT restaurant.id_restaurant, restaurant.articles, restaurant.updated_date from restaurant
        ORDER BY updated_date DESC
        LIMIT 1;
    `

	stmt, err := db.Prepare(request)
	if err != nil {
		return false, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			return
		}
	}(stmt)

	var restaurant models.Restaurant
	err = stmt.QueryRow().Scan(&restaurant.ID, &restaurant.Articles, &restaurant.UpdatedDate)
	if err != nil && err != sql.ErrNoRows {
		return false, fmt.Errorf("failed to get restaurant: %w", err)
	}

	// Check if we need to update the menu
	needsUpdate := true
	if err != sql.ErrNoRows && restaurant.ID != 0 {
		var dbMenuData models.MenuData
		if err := json.Unmarshal([]byte(restaurant.Articles), &dbMenuData); err == nil {
			if compareMenus(&dbMenuData, menuData) {
				log.Println("Menu hasn't changed, no update needed")
				return false, nil
			}
		}
	}

	if needsUpdate {
		log.Println("Menu has changed, updating database")

		_, err = updateRestaurantMenu(db, menuData)
		if err != nil {
			return false, fmt.Errorf("failed to update menu in database: %w", err)
		}

		log.Println("Sending notifications about updated menu")
		if len(cachedMenus) == 0 {
			err = notificationService.SendDailyMenuNotification()
			if err != nil {
				return true, fmt.Errorf("menu updated but failed to send notification: %w", err)
			}
		}

		return true, nil
	}

	return false, nil
}

func compareMenus(menu1, menu2 *models.MenuData) bool {
	if !compareStringSlices(menu1.GrilladesMidi, menu2.GrilladesMidi) ||
		!compareStringSlices(menu1.Migrateurs, menu2.Migrateurs) ||
		!compareStringSlices(menu1.Cibo, menu2.Cibo) ||
		!compareStringSlices(menu1.AccompMidi, menu2.AccompMidi) ||
		!compareStringSlices(menu1.GrilladesSoir, menu2.GrilladesSoir) ||
		!compareStringSlices(menu1.AccompSoir, menu2.AccompSoir) {
		return false
	}
	return true
}

func compareStringSlices(slice1, slice2 []string) bool {
	if len(slice1) != len(slice2) {
		return false
	}

	map1 := make(map[string]int)
	for _, s := range slice1 {
		map1[s]++
	}

	map2 := make(map[string]int)
	for _, s := range slice2 {
		map2[s]++
	}

	for k, v := range map1 {
		if map2[k] != v {
			return false
		}
	}

	return true
}
