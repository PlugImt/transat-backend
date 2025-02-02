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
	request := `
		SELECT id_traq,
			traq.name,
			COALESCE(disabled, false)     AS disabled,
			COALESCE(limited, false)      AS limited,
			COALESCE(alcool, 0)           AS alcool,
			COALESCE(out_of_stock, false) AS out_of_stock,
			creation_date,
			picture,
			COALESCE(description, '')     AS description,
			price,
			COALESCE(price_half, 0)       AS price_half,
			traq_types.name               AS type
		FROM traq
			join traq_types on traq.id_traq_types = traq_types.id_traq_types;
`

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
		var disabled bool
		var limited bool
		var alcool float64
		var out_of_stock bool
		var creation_date string
		var picture string
		var description string
		var price float64
		var price_half float64
		var tpe string

		err := rows.Scan(&id, &name, &disabled, &limited, &alcool, &out_of_stock, &creation_date, &picture, &description, &price, &price_half, &tpe)
		if err != nil {
			log.Fatalf("Error scanning rows: %v", err)
		}

		result = append(result, TraqArticle{
			ID:           id,
			Name:         name,
			Disabled:     disabled,
			Limited:      limited,
			Alcool:       alcool,
			OutOfStock:   out_of_stock,
			CreationDate: creation_date,
			Picture:      picture,
			Description:  description,
			Price:        price,
			PriceHalf:    price_half,
			Type:         tpe,
		})
	}

	return c.JSON(result)
}
