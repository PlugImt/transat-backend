package main

type TurnstileResponse struct {
	Success bool     `json:"success"`
	Error   []string `json:"error-codes"`
}

type TraqArticle struct {
	ID           int     `json:"id"`
	Name         string  `json:"name"`
	Disabled     bool    `json:"disabled"`
	Limited      bool    `json:"limited"`
	Alcool       float64 `json:"alcool"`
	OutOfStock   bool    `json:"out_of_stock"`
	CreationDate string  `json:"creation_date"`
	Picture      string  `json:"picture"`
	Description  string  `json:"description"`
	Price        float64 `json:"price"`
	PriceHalf    float64 `json:"price_half"`
	Type         string  `json:"type"`
}
