package models

// Restaurant represents a record in the 'restaurant' table
type Restaurant struct {
	ID          int    `json:"id_restaurant"`
	Articles    string `json:"articles"` // JSON string of MenuData
	UpdatedDate string `json:"updated_date"`
	Language    int    `json:"language"`
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
