package main

import (
	"database/sql"
	"github.com/gofiber/fiber/v2"
	"log"
)

func login(c *fiber.Ctx) error {
	return c.SendString("User Logged In")
}

func traq(c *fiber.Ctx) error {
	request := `SELECT 
    			COALESCE(id, 0),
    			COALESCE(name, ''),
    			COALESCE(description, ''),
    			COALESCE(limited, false),
    			COALESCE(limited_expiration, ''),
    			COALESCE(creation_date, ''),
    			COALESCE(disabled, false),
    			COALESCE(out_of_stock, false),
    			COALESCE(price_half, 0.0),
    			COALESCE(price_full, 0.0),
    			COALESCE(id_type, 0),
    			COALESCE(alcool, 0.0),
    			COALESCE(image, '')
    FROM transat.traq_articles`

	rows, err := db.Query(request)
	if err != nil {
		log.Fatalf("Error querying database: %v", err)
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Fatalf("Error closing rows: %v", err)
		}
	}(rows)

	var result []TraqArticle

	for rows.Next() {
		var id int
		var name string
		var description string
		var limited bool
		var limited_expiration string
		var creation_date string
		var disabled bool
		var out_of_stock bool
		var price_half float64
		var price_full float64
		var id_type int
		var alcool float32
		var image string

		err := rows.Scan(&id, &name, &description, &limited, &limited_expiration, &creation_date, &disabled, &out_of_stock, &price_half, &price_full, &id_type, &alcool, &image)
		if err != nil {
			log.Fatalf("Error scanning rows: %v", err)
		}

		result = append(result, TraqArticle{
			ID:                id,
			Name:              name,
			Description:       description,
			Limited:           limited,
			LimitedExpiration: limited_expiration,
			CreationDate:      creation_date,
			Disabled:          disabled,
			OutOfStock:        out_of_stock,
			PriceHalf:         price_half,
			PriceFull:         price_full,
			IDType:            id_type,
			Alcool:            alcool,
			Image:             image,
		})
	}

	return c.JSON(result)
}
