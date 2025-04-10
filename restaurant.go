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

	"github.com/gofiber/fiber/v2"
)

var (
	lastFetchTime time.Time
	cachedMenus   map[int]*MenuData // Map of language ID to menu data
	cacheMutex    sync.Mutex
)

func init() {
	cachedMenus = make(map[int]*MenuData)
}

func fetchMenuFromAPI() (*MenuData, error) {
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

	var nestedItems [][]MenuItem
	if err := json.Unmarshal([]byte(loadingData), &nestedItems); err != nil {
		return nil, fmt.Errorf("unable to parse JSON: %w", err)
	}

	if len(nestedItems) == 0 || len(nestedItems[0]) == 0 {
		return nil, fmt.Errorf("no menu items found in the response")
	}

	menuItems := nestedItems[0]

	menuData := MenuData{
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

func updateRestaurantMenu(db *sql.DB, menuData *MenuData) (*Restaurant, error) {
	menuJSON, err := json.Marshal(menuData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal menu data: %w", err)
	}

	if len(string(menuJSON)) > 5000 {
		return nil, fmt.Errorf("menu data too large for database column (max 5000 characters)")
	}

	query := `INSERT INTO restaurant (articles) VALUES ($1) RETURNING id_restaurant, articles, updated_date`

	var restaurant Restaurant
	err = db.QueryRow(query, string(menuJSON)).Scan(&restaurant.ID, &restaurant.Articles, &restaurant.UpdatedDate)
	if err != nil {
		return nil, fmt.Errorf("failed to update restaurant menu: %w", err)
	}

	return &restaurant, nil
}

// isMenuFromToday checks if a menu's updated date is from today
func isMenuFromToday(updatedDate string) bool {
	menuTime, err := time.Parse("2006-01-02 15:04:05", updatedDate)
	if err != nil {
		return false
	}

	now := time.Now()
	return menuTime.Year() == now.Year() &&
		menuTime.Month() == now.Month() &&
		menuTime.Day() == now.Day()
}

func getRestaurant(c *fiber.Ctx) error {
	log.Println("╔======== 🍽️ Get Restaurant 🍽️ ========╗")

	// Parse request body for language
	var req RestaurantRequest
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

	// Check if we have a cached menu for today in the requested language
	if cachedMenus[langID] != nil {
		log.Println("║ ✅ Returning cached menu data for language", req.Language)
		log.Println("╚=========================================╝")
		var menuToReturn FullMenuData
		menuToReturn.MenuData = *cachedMenus[langID]
		menuToReturn.UpdatedDate = time.Now().Format("2006-01-02 15:04:05")
		return c.JSON(menuToReturn)
	}

	// If no cached menu, fetch from database
	request := `
		SELECT restaurant.id_restaurant, restaurant.articles, restaurant.updated_date, restaurant.language 
		FROM restaurant 
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
	defer stmt.Close()

	var restaurant Restaurant
	err = stmt.QueryRow(langID).Scan(&restaurant.ID, &restaurant.Articles, &restaurant.UpdatedDate, &restaurant.Language)
	if err == nil {
		var dbMenuData MenuData
		if err := json.Unmarshal([]byte(restaurant.Articles), &dbMenuData); err == nil {
			cachedMenus[langID] = &dbMenuData
			log.Println("║ ✅ Returning menu data from database for language", req.Language)
			log.Println("╚=========================================╝")
			var menuToReturn FullMenuData
			menuToReturn.MenuData = dbMenuData
			menuToReturn.UpdatedDate = restaurant.UpdatedDate
			return c.JSON(menuToReturn)
		}
	}

	// If no menu in database, fetch from API
	log.Println("║ 🔄 Fetching menu from API")
	menuData, err := fetchMenuFromAPI()
	if err != nil {
		log.Println("║ 💥 Failed to fetch menu from API: ", err)
		log.Println("╚=========================================╝")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch menu",
		})
	}

	// Cache French menu
	cachedMenus[1] = menuData

	// Check if French menu exists in database
	var frenchRestaurant Restaurant
	err = stmt.QueryRow(1).Scan(&frenchRestaurant.ID, &frenchRestaurant.Articles, &frenchRestaurant.UpdatedDate, &frenchRestaurant.Language)
	if err != nil {
		// If no French menu in database, insert it and send notification
		menuJSON, err := json.Marshal(menuData)
		if err != nil {
			log.Println("║ 💥 Failed to marshal menu data: ", err)
			log.Println("╚=========================================╝")
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to prepare menu data for database",
			})
		}

		updateQuery := `
			INSERT INTO restaurant (articles, language) 
			VALUES ($1, $2)
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

		_, err = updateStmt.Exec(string(menuJSON), 1)
		if err != nil {
			log.Println("║ 💥 Failed to update restaurant menu in database: ", err)
		} else {
			log.Println("║ ✅ Restaurant menu updated successfully")
			log.Println("║ 🔔 Sending notification to subscribers")
			notificationService := NewNotificationService(db)
			err = notificationService.SendDailyMenuNotification()
			if err != nil {
				log.Println("║ ⚠️ Failed to send notification: ", err)
			}
		}
	}

	// If requested language is not French, translate the menu
	if langID != 1 {
		// Check if translated menu already exists in database
		var translatedRestaurant Restaurant
		err = stmt.QueryRow(langID).Scan(&translatedRestaurant.ID, &translatedRestaurant.Articles, &translatedRestaurant.UpdatedDate, &translatedRestaurant.Language)
		if err != nil {
			// If no translated menu in database, create and insert it
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

			// Update database with translated menu
			translatedMenuJSON, err := json.Marshal(translatedMenu)
			if err != nil {
				log.Println("║ 💥 Failed to marshal translated menu data: ", err)
				log.Println("╚=========================================╝")
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Failed to prepare translated menu data for database",
				})
			}

			updateQuery := `
			INSERT INTO restaurant (articles, language)
				VALUES ($1, $2)
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

			_, err = updateStmt.Exec(string(translatedMenuJSON), langID)
			if err != nil {
				log.Println("║ 💥 Failed to update translated menu in database: ", err)
			}

			cachedMenus[langID] = translatedMenu
			menuData = translatedMenu
		} else {
			// If translated menu exists, use it
			var dbMenuData MenuData
			if err := json.Unmarshal([]byte(translatedRestaurant.Articles), &dbMenuData); err == nil {
				cachedMenus[langID] = &dbMenuData
				menuData = &dbMenuData
			}
		}
	}

	log.Println("╚=========================================╝")

	var menuToReturn FullMenuData
	menuToReturn.MenuData = *menuData
	menuToReturn.UpdatedDate = time.Now().Format("2006-01-02 15:04:05")
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

	var restaurant Restaurant
	err = stmt.QueryRow().Scan(&restaurant.ID, &restaurant.Articles, &restaurant.UpdatedDate)
	if err != nil && err != sql.ErrNoRows {
		return false, fmt.Errorf("failed to get restaurant: %w", err)
	}

	// Check if we need to update the menu
	needsUpdate := true
	if err != sql.ErrNoRows && restaurant.ID != 0 {
		var dbMenuData MenuData
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
		err = notificationService.SendDailyMenuNotification()
		if err != nil {
			return true, fmt.Errorf("menu updated but failed to send notification: %w", err)
		}

		return true, nil
	}

	return false, nil
}

func compareMenus(menu1, menu2 *MenuData) bool {
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

// Reset cache at midnight
func resetCache() {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()
	cachedMenus = make(map[int]*MenuData)
	lastFetchTime = time.Time{}
}
