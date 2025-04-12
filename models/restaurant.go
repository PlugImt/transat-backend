package models

type Restaurant struct {
	ID          int    `json:"id_restaurant"`
	Articles    string `json:"articles"`
	UpdatedDate string `json:"updated_date"`
	Language    int    `json:"language"`
}

type MenuItem struct {
	Pole           string `json:"pole"`
	Accompagnement string `json:"accompagnement"`
	Periode        string `json:"periode"`
	Nom            string `json:"nom"`
	Info1          string `json:"info1"`
	Info2          string `json:"info2"`
}

type MenuData struct {
	GrilladesMidi []string `json:"grilladesMidi"`
	Migrateurs    []string `json:"migrateurs"`
	Cibo          []string `json:"cibo"`
	AccompMidi    []string `json:"accompMidi"`
	GrilladesSoir []string `json:"grilladesSoir"`
	AccompSoir    []string `json:"accompSoir"`
}

type FullMenuData struct {
	MenuData
	UpdatedDate string `json:"updatedDate"`
}

type RestaurantRequest struct {
	Language string `json:"language"`
}
