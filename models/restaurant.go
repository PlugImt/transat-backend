package models

import "time"

// Restaurant represents a record in the 'restaurant' table
type Restaurant struct {
	ID   int    `json:"id_restaurant" db:"id_restaurant"`
	Type string `json:"type" db:"type"`
}

// RestaurantArticle represents a record in the 'restaurant_articles' table
type RestaurantArticle struct {
	ID              int        `json:"id_restaurant_articles" db:"id_restaurant_articles"`
	FirstTimeServed time.Time  `json:"first_time_served" db:"first_time_served"`
	LastTimeServed  *time.Time `json:"last_time_served,omitempty" db:"last_time_served"`
	Name            string     `json:"name" db:"name"`
}

// RestaurantArticleNote represents a record in the 'restaurant_articles_notes' table
type RestaurantArticleNote struct {
	Email     string `json:"email" db:"email"`
	ArticleID int    `json:"id_restaurant_articles" db:"id_restaurant_articles"`
	Note      int    `json:"note" db:"note"`
	Comment   string `json:"comment,omitempty" db:"comment"`
}

// RestaurantMeal represents a record in the 'restaurant_meals' table
type RestaurantMeal struct {
	RestaurantID int       `json:"id_restaurant" db:"id_restaurant"`
	ArticleID    int       `json:"id_restaurant_articles" db:"id_restaurant_articles"`
	DateServed   time.Time `json:"date_served" db:"date_served"`
}

// MenuItemAPI represents the structure of an item returned by the V1 source API.
// Renamed from MenuItem for clarity.
type MenuItemAPI struct {
	Pole           string `json:"pole"`
	Accompagnement string `json:"accompagnement"` // Expects "TRUE" or "FALSE"
	Periode        string `json:"periode"`        // Expects "midi" or "soir"
	Nom            string `json:"nom"`
	Info1          string `json:"info1"`
	Info2          string `json:"info2"`
}

// FetchedItem represents a processed menu item ready for database synchronization
type FetchedItem struct {
	Name       string `json:"name"`
	Category   string `json:"category"`
	MenuTypeID int    `json:"menu_type_id"`
}

// MenuData holds the categorized menu items.
type MenuData struct {
	GrilladesMidi []string `json:"grilladesMidi"`
	Migrateurs    []string `json:"migrateurs"`
	Cibo          []string `json:"cibo"`
	AccompMidi    []string `json:"accompMidi"`
	GrilladesSoir []string `json:"grilladesSoir"`
	AccompSoir    []string `json:"accompSoir"`
}

// FullMenuData combines MenuData with the update timestamp.
type FullMenuData struct {
	MenuData
	UpdatedDate string `json:"updatedDate"`
}

// MenuEntry represents a menu item with its average rating for API responses
type MenuEntry struct {
	ArticleID     int     `json:"article_id" db:"id_restaurant_articles"`
	Name          string  `json:"name" db:"name"`
	MenuTypeID    int     `json:"menu_type_id" db:"id_restaurant"`
	AverageRating float64 `json:"average_rating" db:"average_rating"`
}

// MenuItemWithRating represents a menu item with its rating
type MenuItemWithRating struct {
	ID            int     `json:"id"`
	Name          string  `json:"name"`
	AverageRating float64 `json:"average_rating"`
}

// MenuResponse represents the complete response for today's menu
type MenuResponse struct {
	Date  string      `json:"date"`
	Items []MenuEntry `json:"items"`
}

// CategorizedMenuResponse represents the menu categorized by type
type CategorizedMenuResponse struct {
	GrilladesMidi []MenuItemWithRating `json:"grilladesMidi"`
	Migrateurs    []MenuItemWithRating `json:"migrateurs"`
	Cibo          []MenuItemWithRating `json:"cibo"`
	AccompMidi    []MenuItemWithRating `json:"accompMidi"`
	GrilladesSoir []MenuItemWithRating `json:"grilladesSoir"`
	AccompSoir    []MenuItemWithRating `json:"accompSoir"`
	UpdatedDate   string               `json:"updatedDate"`
}

// ReviewResponse represents a review with user details
type ReviewResponse struct {
	FirstName      string    `json:"first_name"`
	LastName       string    `json:"last_name"`
	ProfilePicture string    `json:"profile_picture"`
	Rating         int       `json:"rating"`
	Comment        string    `json:"comment,omitempty"`
	Date           time.Time `json:"date"`
}

// ReviewResult represents the result of saving a review
type ReviewResult struct {
	Message       string  `json:"message"`
	DishName      string  `json:"dish_name"`
	AverageRating float64 `json:"average_rating"`
	TotalRatings  int     `json:"total_ratings"`
	YourRating    int     `json:"your_rating"`
	YourComment   string  `json:"your_comment"`
}
