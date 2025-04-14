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
	"strings"
	"sync"
	"time"

	"Transat_2.0_Backend/models"
	"Transat_2.0_Backend/services"
	"Transat_2.0_Backend/utils" // logger

	"github.com/gofiber/fiber/v2"
)

// RestaurantHandler handles menu fetching, caching, and potentially updates.
type RestaurantHandler struct {
	DB           *sql.DB
	TransService *services.TranslationService  // Service for translating menus
	NotifService *services.NotificationService // Service for sending notifications

	// Cache - Consider a more robust cache implementation (e.g., Redis, dedicated cache library)
	cacheMutex    sync.RWMutex             // Mutex for thread-safe cache access
	lastFetchTime time.Time                // Timestamp of the last successful API fetch
	cachedMenus   map[int]*models.MenuData // Map[languageID] -> MenuData
	menuSourceURL string                   // URL to fetch the menu from
	apiRegex      *regexp.Regexp           // Compiled regex for parsing API response
}

// NewRestaurantHandler creates a new RestaurantHandler.
func NewRestaurantHandler(db *sql.DB, transService *services.TranslationService, notifService *services.NotificationService) *RestaurantHandler {
	// Compile regex once
	// regex := regexp.MustCompile(`var loadingData = (\s*\{\s*.*?\s*\}\s*);?`) // Updated Regex to capture object
	// Corrected regex for the [[...]] format found in data.html
	regex := regexp.MustCompile(`var loadingData\s*=\s*(\[\[.*?\]\]);?`)
	// Make source URL configurable
	sourceURL := "https://toast-js.ew.r.appspot.com/coteresto?key=1ohdRUdCYo6e71aLuBh7ZfF2lc_uZqp9D78icU4DPufA"

	return &RestaurantHandler{
		DB:            db,
		TransService:  transService,
		NotifService:  notifService,
		cachedMenus:   make(map[int]*models.MenuData),
		menuSourceURL: sourceURL,
		apiRegex:      regex,
	}
}

// GetRestaurantMenu handles requests for the restaurant menu.
// It uses caching and fetches/translates the menu if needed.
func (h *RestaurantHandler) GetRestaurantMenu(c *fiber.Ctx) error {
	utils.LogHeader("🍽️ Get Restaurant Menu")

	// decode RestaurantRequest from query params
	language := c.Query("language")

	if language == "" {
		utils.LogMessage(utils.LevelWarn, "No language specified, defaulting to French")
		language = "fr"
	}

	utils.LogLineKeyValue(utils.LevelInfo, "Requested Language", language)

	// Convert language code to language ID (using helper)
	langID := getLanguageID(language)
	if langID == 0 {
		utils.LogMessage(utils.LevelWarn, "Invalid language code requested")
		utils.LogFooter()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fmt.Sprintf("Invalid language code: %s", language),
		})
	}
	utils.LogLineKeyValue(utils.LevelInfo, "Language ID", langID)

	// --- Cache Check ---
	h.cacheMutex.RLock() // Read lock for checking cache
	cachedMenu := h.cachedMenus[langID]
	lastFetch := h.lastFetchTime
	h.cacheMutex.RUnlock()

	now := time.Now()
	// Check if cache is valid for today
	// Allow fetching new data if cache is old, even if it's the same day (e.g., force refresh)
	// Let's refresh if cache is older than, say, 1 hour? Or just check date?
	// For simplicity, let's check if it's from today first.
	isCacheFromToday := !lastFetch.IsZero() &&
		now.Year() == lastFetch.Year() &&
		now.Month() == lastFetch.Month() &&
		now.Day() == lastFetch.Day()

	utils.LogLineKeyValue(utils.LevelDebug, "isCacheFromToday", isCacheFromToday)
	utils.LogLineKeyValue(utils.LevelDebug, "lastFetch", lastFetch)

	if cachedMenu != nil && isCacheFromToday {
		utils.LogMessage(utils.LevelInfo, "Returning cached menu data")
		utils.LogFooter()
		return c.JSON(models.FullMenuData{
			MenuData:    *cachedMenu,
			UpdatedDate: lastFetch.Format(time.RFC3339), // Use standard format
		})
	}

	// --- Cache Miss or Stale: Fetch/Update ---
	utils.LogMessage(utils.LevelInfo, "Cache miss or stale, proceeding to fetch/update")

	// Acquire write lock for potential cache update
	h.cacheMutex.Lock()
	defer h.cacheMutex.Unlock() // Ensure unlock happens

	// Double-check cache after acquiring write lock (another request might have updated it)
	cachedMenu = h.cachedMenus[langID]
	lastFetch = h.lastFetchTime
	isCacheFromToday = !lastFetch.IsZero() &&
		now.Year() == lastFetch.Year() &&
		now.Month() == lastFetch.Month() &&
		now.Day() == lastFetch.Day()

	if cachedMenu != nil && isCacheFromToday {
		utils.LogMessage(utils.LevelInfo, "Returning cached menu data (populated by another request)")
		utils.LogFooter()
		return c.JSON(models.FullMenuData{
			MenuData:    *cachedMenu,
			UpdatedDate: lastFetch.Format(time.RFC3339),
		})
	}

	// --- Fetch from Source API ---
	utils.LogMessage(utils.LevelInfo, "Fetching fresh menu data from source API")
	baseMenuData, err := h.fetchMenuFromAPI() // Fetch the base (French) menu
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to fetch menu from source API")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		// Fallback: Try to return the latest data from the database if fetch fails
		latestDbMenu, dbErr := h.getLatestMenuFromDB(langID)
		if dbErr == nil && latestDbMenu != nil {
			utils.LogMessage(utils.LevelWarn, "API fetch failed, returning latest menu from database")
			utils.LogFooter()
			// Update cache with DB data? Or leave cache stale? Let's update.
			h.cachedMenus[langID] = &latestDbMenu.MenuData
			// Don't update h.lastFetchTime from DB data to encourage refetching API later
			return c.JSON(latestDbMenu)
		}
		utils.LogFooter()
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"error": "Failed to fetch latest menu data",
		})
	}
	// Successfully fetched base menu
	fetchedTime := time.Now()

	// --- Translation (if needed) ---
	finalMenuData := baseMenuData
	if langID != 1 { // Language ID 1 is assumed to be French (source)
		utils.LogMessage(utils.LevelInfo, "Translating base menu")
		utils.LogLineKeyValue(utils.LevelInfo, "Target Language", language)
		translatedMenu, transErr := h.TransService.TranslateMenu(baseMenuData, language)
		if transErr != nil {
			utils.LogMessage(utils.LevelError, "Failed to translate menu")
			utils.LogLineKeyValue(utils.LevelError, "Error", transErr)
			// Return the base (French) menu with an error? Or just the base menu?
			// Let's return the base menu and log the error.
			finalMenuData = baseMenuData // Use base menu as fallback
		} else {
			finalMenuData = translatedMenu // Use translated menu
			utils.LogMessage(utils.LevelInfo, "Menu translated successfully")
		}
	}

	// --- Update Cache ---
	utils.LogMessage(utils.LevelInfo, "Updating menu cache")
	h.cachedMenus[langID] = finalMenuData
	// Only update lastFetchTime if the base fetch was successful
	h.lastFetchTime = fetchedTime

	// --- Update Database (Asynchronously?) ---
	// We update the DB with the potentially translated menu for the specific language ID
	// Run this in a goroutine so it doesn't block the response.
	go func(menuToSave *models.MenuData, langToSave int, fetchTimestamp time.Time) {
		saveErr := h.saveMenuToDB(menuToSave, langToSave, fetchTimestamp)
		if saveErr != nil {
			log.Printf("Error saving menu (lang %d) to DB asynchronously: %v", langToSave, saveErr)
		} else {
			log.Printf("Successfully saved menu (lang %d) to DB asynchronously", langToSave)
			// --- Send Notifications (only if French menu was updated) ---
			// This logic should ideally live within the save/update function or be triggered by it.
			// Only notify if the *base* French menu changed and was successfully saved.
			if langToSave == 1 { // If we just saved the French menu
				// Check if menu actually changed compared to previous DB version? (Requires extra logic)
				// For now, notify if base menu was fetched & saved.
				if h.NotifService != nil {
					utils.LogMessage(utils.LevelInfo, "Triggering daily menu notification send")
					notifErr := h.NotifService.SendDailyMenuNotification() // Send the base menu
					if notifErr != nil {
						log.Printf("Error sending daily menu notification asynchronously: %v", notifErr)
					}
				} else {
					log.Println("Warning: NotificationService not available for sending menu update.")
				}
			}
		}

	}(finalMenuData, langID, fetchedTime)

	utils.LogMessage(utils.LevelInfo, "Returning fetched/translated menu data")
	utils.LogFooter()

	// Return the fresh data
	return c.JSON(models.FullMenuData{
		MenuData:    *finalMenuData,
		UpdatedDate: fetchedTime.Format(time.RFC3339),
	})
}

// --- Internal Methods ---

// fetchMenuFromAPI fetches and parses the base menu (French) from the source URL.
func (h *RestaurantHandler) fetchMenuFromAPI() (*models.MenuData, error) {
	resp, err := http.Get(h.menuSourceURL)
	if err != nil {
		return nil, fmt.Errorf("unable to fetch URL '%s': %w", h.menuSourceURL, err)
	}
	defer resp.Body.Close()

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

// parse `var loadingData = [[...]]`
func (h *RestaurantHandler) parseLoadingData(jsonData string) (*models.MenuData, error) {
	var nestedItems [][]models.MenuItemAPI // Use specific struct for API parsing
	if err := json.Unmarshal([]byte(jsonData), &nestedItems); err != nil {
		log.Printf("Invalid JSON structure from API: %v", err)
		return nil, fmt.Errorf("unable to parse menu JSON from API: %w", err)
	}

	if len(nestedItems) == 0 || len(nestedItems[0]) == 0 {
		return nil, fmt.Errorf("no menu items found in the parsed API data")
	}

	// Process raw items into structured MenuData
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
	itemMap := make(map[string]map[string]bool) // category -> item -> exists

	for _, item := range items {
		category := getMenuCategory(item.Pole, item.Accompagnement, item.Periode)
		if category == "" {
			continue // Skip items that don't map to a known category
		}

		// Combine name and info fields, trimming space
		menuItemText := strings.TrimSpace(fmt.Sprintf("%s %s %s", item.Nom, item.Info1, item.Info2))
		// remove all double or more spaces
		menuItemText = strings.Join(strings.Fields(menuItemText), " ")
		if menuItemText == "" {
			continue // Skip empty items
		}

		// Ensure category map exists
		if itemMap[category] == nil {
			itemMap[category] = make(map[string]bool)
		}

		// Add item if not already present in the map for this category
		if !itemMap[category][menuItemText] {
			itemMap[category][menuItemText] = true // Mark as added
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

// saveMenuToDB saves the provided menu data for the given language ID.
// It checks if an identical menu for the same language exists for today before inserting.
func (h *RestaurantHandler) saveMenuToDB(menuData *models.MenuData, langID int, fetchTimestamp time.Time) error {
	menuJSON, err := json.Marshal(menuData)
	if err != nil {
		return fmt.Errorf("failed to marshal menu data for DB: %w", err)
	}

	// Check if menu data exceeds reasonable size (e.g., TEXT column limit in PG is large, but good to have a sanity check)
	// const maxMenuSize = 10000 // Example limit
	// if len(menuJSON) > maxMenuSize {
	// 	return fmt.Errorf("menu data JSON size (%d) exceeds limit (%d)", len(menuJSON), maxMenuSize)
	// }

	// Use transaction to check existence and insert atomically
	tx, err := h.DB.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin DB transaction: %w", err)
	}
	defer tx.Rollback() // Rollback if commit doesn't happen

	// Check if this exact menu for this language already exists today
	var exists bool
	checkQuery := `
		SELECT EXISTS (
			SELECT 1 FROM restaurant
			WHERE language = $1
			  AND articles::jsonb = $2::jsonb -- Compare JSON content effectively
			  AND DATE(updated_date) = DATE($3) -- Check if entry exists for today
		);`
	// Use fetchTimestamp for the date check
	err = tx.QueryRow(checkQuery, langID, string(menuJSON), fetchTimestamp).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check for existing menu in DB: %w", err)
	}

	if exists {
		log.Printf("Identical menu for language %d already exists in DB for today. Skipping insert.", langID)
		return nil // Not an error, just no need to insert
	}

	// Insert the new menu record
	insertQuery := `
		INSERT INTO restaurant (articles, language, updated_date)
		VALUES ($1, $2, $3);
	`
	_, err = tx.Exec(insertQuery, string(menuJSON), langID, fetchTimestamp)
	if err != nil {
		return fmt.Errorf("failed to insert new menu into DB: %w", err)
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit menu insert transaction: %w", err)
	}

	return nil
}

// getLatestMenuFromDB retrieves the most recent menu entry for a given language from the database.
func (h *RestaurantHandler) getLatestMenuFromDB(langID int) (*models.FullMenuData, error) {
	query := `
		SELECT id_restaurant, articles, updated_date
		FROM restaurant
		WHERE language = $1
		ORDER BY updated_date DESC
		LIMIT 1;
	`
	var idRestaurant int
	var articlesJSON string
	var updatedDate time.Time

	err := h.DB.QueryRow(query, langID).Scan(&idRestaurant, &articlesJSON, &updatedDate)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // No menu found in DB for this language
		}
		return nil, fmt.Errorf("failed to query latest DB menu (lang %d): %w", langID, err)
	}

	var menuData models.MenuData
	if err := json.Unmarshal([]byte(articlesJSON), &menuData); err != nil {
		log.Printf("Failed to parse menu JSON from database (lang %d, date %s): %v", langID, updatedDate, err)
		// Data in DB is corrupt, return nil or error?
		return nil, fmt.Errorf("failed to parse menu data from DB: %w", err)
	}

	utils.LogLineKeyValue(utils.LevelDebug, "Latest DB Menu ID", idRestaurant)

	return &models.FullMenuData{
		MenuData:    menuData,
		UpdatedDate: updatedDate.Format(time.RFC3339),
	}, nil
}

// CheckAndUpdateMenuCron is intended to be called by the cron job.
// It fetches the base menu, compares it with the last known base menu in DB,
// saves if different, and triggers notifications.
func (h *RestaurantHandler) CheckAndUpdateMenuCron() (bool, error) {
	utils.LogHeader("🤖 Cron: Check & Update Restaurant Menu")

	// 1. Fetch latest base menu (French, langID=1) from API
	baseMenuData, err := h.fetchMenuFromAPI()
	if err != nil {
		utils.LogMessage(utils.LevelError, "Cron: Failed to fetch base menu from API")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		utils.LogFooter()
		return false, err // Don't proceed if API fetch fails
	}
	fetchedTime := time.Now()

	// 2. Get latest base menu from DB for comparison
	latestDbMenu, err := h.getLatestMenuFromDB(1) // Get French menu (langID 1)
	utils.LogLineKeyValue(utils.LevelDebug, "Latest DB Menu", latestDbMenu)
	needsUpdate := true // Assume update needed unless proven otherwise
	if err != nil {
		utils.LogMessage(utils.LevelWarn, "Cron: Failed to get latest base menu from DB for comparison")
		utils.LogLineKeyValue(utils.LevelWarn, "Error", err)
		// Proceed with saving the fetched menu as we can't compare
	} else if latestDbMenu != nil {
		// Compare fetched menu with DB menu
		if compareMenus(&latestDbMenu.MenuData, baseMenuData) {
			utils.LogMessage(utils.LevelInfo, "Cron: Fetched menu is identical to latest in DB. No update needed.")
			needsUpdate = false
		} else {
			utils.LogMessage(utils.LevelInfo, "Cron: Menu change detected.")
		}
	} else {
		utils.LogMessage(utils.LevelInfo, "Cron: No existing base menu found in DB. Saving fetched menu.")
	}

	// 3. Save to DB if changed (or no previous version)
	updated := false
	if needsUpdate {
		utils.LogMessage(utils.LevelInfo, "Cron: Saving updated base menu to DB")
		saveErr := h.saveMenuToDB(baseMenuData, 1, fetchedTime) // Save French menu
		if saveErr != nil {
			utils.LogMessage(utils.LevelError, "Cron: Failed to save updated base menu to DB")
			utils.LogLineKeyValue(utils.LevelError, "Error", saveErr)
			utils.LogFooter()
			// Return false as update failed, but log error
			return false, saveErr
		}
		updated = true // Mark as updated in this cycle

		// 4. Update cache with the new base menu
		h.cacheMutex.Lock()
		h.cachedMenus[1] = baseMenuData // Update French cache
		h.lastFetchTime = fetchedTime   // Update fetch time since we got new data
		// Invalidate caches for other languages? Or let them fetch/translate on demand?
		// Let's invalidate other languages for simplicity.
		for langID := range h.cachedMenus {
			if langID != 1 {
				delete(h.cachedMenus, langID)
			}
		}
		h.cacheMutex.Unlock()
		utils.LogMessage(utils.LevelInfo, "Cron: Updated menu cache")

		// 5. Trigger notifications (only if base menu was updated)
		if h.NotifService != nil {
			utils.LogMessage(utils.LevelInfo, "Cron: Triggering daily menu notification send")
			notifErr := h.NotifService.SendDailyMenuNotification()
			if notifErr != nil {
				utils.LogMessage(utils.LevelError, "Cron: Failed to send daily menu notification")
				utils.LogLineKeyValue(utils.LevelError, "Error", notifErr)
				// Log error, but don't necessarily fail the whole check-update cycle because of notification failure
			}
		} else {
			utils.LogMessage(utils.LevelWarn, "Cron: NotificationService not available.")
		}
	}

	utils.LogFooter()
	return updated, nil // Return true if an update was saved
}

// --- Helper Functions ---

// getMenuCategory maps API fields to internal category names.
func getMenuCategory(pole string, accompagnement string, periode string) string {
	// Handle boolean string conversion carefully
	isAccomp := strings.EqualFold(accompagnement, "TRUE")

	if isAccomp {
		if periode == "midi" {
			return "accompMidi"
		}
		// Assume anything not "midi" is "soir" for accompaniments
		return "accompSoir"
	}

	// Map main dishes based on pole and period
	switch pole {
	case "Grillades / Plats traditions":
		if periode == "midi" {
			return "grilladesMidi"
		}
		return "grilladesSoir" // Assume soir if not midi
	case "Les Cuistots migrateurs":
		return "migrateurs" // Assume migrateurs are same midi/soir? API data might clarify.
	case "Le Végétarien":
		return "cibo" // Assume cibo/vegetarian is same midi/soir?
	default:
		// Log unknown poles?
		// log.Printf("Warning: Unknown menu pole encountered: '%s'", pole)
		return "" // Ignore unknown poles
	}
}

// getLanguageID converts language code string to internal ID.
// TODO: Fetch this mapping from the database 'languages' table instead of hardcoding.
func getLanguageID(langCode string) int {
	switch strings.ToLower(langCode) {
	case "fr":
		return 1
	case "en":
		return 2
	case "es":
		return 3
	case "de":
		return 4
	case "it":
		return 5
	case "pt":
		return 6
	case "ru":
		return 7
	case "zh":
		return 8
	case "ko":
		return 9
	default:
		return 0 // Invalid/Unknown
	}
}

// compareMenus compares two MenuData structs for equality.
func compareMenus(menu1, menu2 *models.MenuData) bool {
	// Use helper to compare slices ignoring order
	return compareStringSlicesIgnoreOrder(menu1.GrilladesMidi, menu2.GrilladesMidi) &&
		compareStringSlicesIgnoreOrder(menu1.Migrateurs, menu2.Migrateurs) &&
		compareStringSlicesIgnoreOrder(menu1.Cibo, menu2.Cibo) &&
		compareStringSlicesIgnoreOrder(menu1.AccompMidi, menu2.AccompMidi) &&
		compareStringSlicesIgnoreOrder(menu1.GrilladesSoir, menu2.GrilladesSoir) &&
		compareStringSlicesIgnoreOrder(menu1.AccompSoir, menu2.AccompSoir)
}

// compareStringSlicesIgnoreOrder checks if two string slices contain the same elements, regardless of order.
func compareStringSlicesIgnoreOrder(slice1, slice2 []string) bool {
	if len(slice1) != len(slice2) {
		return false
	}
	if len(slice1) == 0 {
		return true // Both empty
	}

	map1 := make(map[string]int)
	for _, s := range slice1 {
		map1[strings.TrimSpace(s)]++ // Trim space for comparison
	}

	map2 := make(map[string]int)
	for _, s := range slice2 {
		map2[strings.TrimSpace(s)]++ // Trim space
	}

	// Compare maps
	if len(map1) != len(map2) {
		return false // Different number of unique items
	}
	for k, v := range map1 {
		if map2[k] != v {
			return false // Different count for an item
		}
	}

	return true
}
