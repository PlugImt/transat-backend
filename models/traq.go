package models

type TraqArticle struct {
	ID           int     `json:"id_traq"`
	Name         string  `json:"name"`
	Disabled     bool    `json:"disabled"`
	Limited      bool    `json:"limited"`
	Alcohol      float32 `json:"alcohol"`
	OutOfStock   bool    `json:"out_of_stock"`
	CreationDate string  `json:"creation_date"`
	Picture      string  `json:"picture"`
	Description  string  `json:"description"`
	Price        float32 `json:"price"`
	PriceHalf    float32 `json:"price_half"`
	TraqType     string  `json:"traq_type"`
}

type TraqType struct {
	IDType int    `json:"id_traq_types"`
	Name   string `json:"name"`
}
