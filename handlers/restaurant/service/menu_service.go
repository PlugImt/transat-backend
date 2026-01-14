package service

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/plugimt/transat-backend/models"
)

type MenuService struct {
	menuSourceURL string
	apiRegex      *regexp.Regexp
}

func NewMenuService() *MenuService {
	// Initialize the regex to match the loadingData array in the API response
	regex := regexp.MustCompile(`var loadingData\s*=\s*(\[\[.*?\]\]);?`)
	sourceURL := "https://toast-js.ew.r.appspot.com/coteresto?key=1ohdRUdCYo6e71aLuBh7ZfF2lc_uZqp9D78icU4DPufA"

	return &MenuService{
		menuSourceURL: sourceURL,
		apiRegex:      regex,
	}
}

// FetchRawMenuItems fetches and parses the raw MenuItemAPI array from the remote source.
func (s *MenuService) FetchRawMenuItems() ([]models.MenuItemAPI, error) {
	resp, err := http.Get(s.menuSourceURL)
	if err != nil {
		return nil, fmt.Errorf("unable to fetch URL '%s': %w", s.menuSourceURL, err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("Error closing response body: %v", err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code %d from '%s'", resp.StatusCode, s.menuSourceURL)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read response body: %w", err)
	}

	matches := s.apiRegex.FindSubmatch(body)
	if len(matches) < 2 {
		log.Printf("Body sample for regex failure: %s", string(body[:200]))
		return nil, fmt.Errorf("regex did not find 'loadingData' array in response")
	}
	loadingDataJSON := string(matches[1])

	var nestedItems [][]models.MenuItemAPI
	if err := json.Unmarshal([]byte(loadingDataJSON), &nestedItems); err != nil {
		log.Printf("Invalid JSON structure from API: %v", err)
		return nil, fmt.Errorf("unable to parse menu JSON from API: %w", err)
	}

	if len(nestedItems) == 0 || len(nestedItems[0]) == 0 {
		return nil, fmt.Errorf("no menu items found in the parsed API data")
	}

	return nestedItems[0], nil
}

// FetchMenuFromAPI returns the legacy MenuData structure for existing callers.
func (s *MenuService) FetchMenuFromAPI() (*models.MenuData, error) {
	items, err := s.FetchRawMenuItems()
	if err != nil {
		return nil, err
	}
	return s.processRawMenuItems(items), nil
}

func (s *MenuService) processRawMenuItems(items []models.MenuItemAPI) *models.MenuData {
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
		category := GetMenuCategory(item.Pole, item.Accompagnement, item.Periode)
		if category == "" {
			continue
		}

		menuItemText := strings.TrimSpace(fmt.Sprintf("%s %s %s", item.Nom, item.Info1, item.Info2))
		menuItemText = strings.Join(strings.Fields(menuItemText), " ")
		if menuItemText == "" {
			continue
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

// BuildFetchedItems converts raw MenuItemAPI entries to FetchedItem slice,
// preserving allergen codes so they can be stored with restaurant_articles.
func (s *MenuService) BuildFetchedItems(items []models.MenuItemAPI) []models.FetchedItem {
	var fetchedItems []models.FetchedItem

	categoryToID := map[string]int{
		"grilladesMidi": 1,
		"migrateurs":    2,
		"cibo":          3,
		"accompMidi":    4,
		"grilladesSoir": 5,
		"accompSoir":    6,
	}

	for _, item := range items {
		category := GetMenuCategory(item.Pole, item.Accompagnement, item.Periode)
		if category == "" {
			continue
		}

		menuItemText := strings.TrimSpace(fmt.Sprintf("%s %s %s", item.Nom, item.Info1, item.Info2))
		menuItemText = strings.Join(strings.Fields(menuItemText), " ")
		if menuItemText == "" {
			continue
		}

		menuTypeID, ok := categoryToID[category]
		if !ok {
			continue
		}

		// Collect allergen codes, skipping empty / placeholder values
		rawCodes := []string{
			item.Allergene1, item.Allergene2, item.Allergene3, item.Allergene4, item.Allergene5,
			item.Allergene6, item.Allergene7, item.Allergene8, item.Allergene9, item.Allergene10,
		}
		var allergenCodes []string
		seen := make(map[string]bool)
		for _, code := range rawCodes {
			code = strings.TrimSpace(code)
			if code == "" {
				continue
			}
			lower := strings.ToLower(code)
			if lower == "vide" || lower == "undefined" {
				continue
			}
			if !seen[code] {
				seen[code] = true
				allergenCodes = append(allergenCodes, code)
			}
		}

		// Collect marker codes (boolean fields that are "TRUE" or "VRAI")
		var markerCodes []string
		markerMap := map[string]string{
			"ardoise":    item.Ardoise,
			"formule":    item.Formule,
			"vitalite":   item.Vitalite,
			"vegetarien": item.Vegetarien,
			"bio":        item.Bio,
			"local":      item.Local,
			"saison":     item.Saison,
			"equitable":  item.Equitable,
			"weightWatcher": item.WW, // API uses "ww", DB uses "weightWatcher"
			"peche":      item.Peche,
			"france":     item.France,
		}

		for markerName, markerValue := range markerMap {
			if strings.EqualFold(markerValue, "TRUE") || strings.EqualFold(markerValue, "VRAI") {
				// Special case: API uses "ww" but DB uses "weightWatcher"
				if markerName == "weightWatcher" {
					markerCodes = append(markerCodes, "weightWatcher")
				} else {
					markerCodes = append(markerCodes, markerName)
				}
			}
		}

		fetchedItems = append(fetchedItems, models.FetchedItem{
			Name:          menuItemText,
			Category:      category,
			MenuTypeID:    menuTypeID,
			AllergenCodes: allergenCodes,
			MarkerCodes:   markerCodes,
		})
	}

	return fetchedItems
}

// GetMenuCategory determines the menu category based on pole, accompagnement, and periode
func GetMenuCategory(pole string, accompagnement string, periode string) string {
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

// ConvertMenuDataToFetchedItems converts the old MenuData format to FetchedItem slice
func (s *MenuService) ConvertMenuDataToFetchedItems(menuData *models.MenuData) []models.FetchedItem {
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
