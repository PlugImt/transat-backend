package main

type TurnstileResponse struct {
	Success bool     `json:"success"`
	Error   []string `json:"error-codes"`
}

type TraqArticle struct {
	ID                int     `json:"id"`
	Name              string  `json:"name"`
	Description       string  `json:"description"`
	Limited           bool    `json:"limited"`
	LimitedExpiration string  `json:"limited_expiration"`
	CreationDate      string  `json:"creation_date"`
	Disabled          bool    `json:"disabled"`
	OutOfStock        bool    `json:"out_of_stock"`
	PriceHalf         float64 `json:"price_half"`
	PriceFull         float64 `json:"price_full"`
	IDType            int     `json:"id_type"`
	Alcool            float32 `json:"alcool"`
	Image             string  `json:"image"`
}
