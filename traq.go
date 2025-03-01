package main

import (
	"database/sql"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"log"
)

func createTraqTypes(c *fiber.Ctx) error {
	log.Println("╔======== 🍺 Create Traq Type 🍺 ========╗")
	var traqType TraqType
	if err := c.BodyParser(&traqType); err != nil {
		log.Println("║ 💥 Failed to parse request body: ", err)
		log.Println("╚=========================================╝")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if traqType.Name == "" {
		log.Println("║ 💥 Some fields are missing")
		log.Println("╚=========================================╝")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Fill in all fields",
		})
	}

	request := `
		INSERT INTO traq_types (name)
		VALUES ($1);
	`

	stmt, err := db.Prepare(request)
	if err != nil {
		log.Println("║ 💥 Failed to prepare statement: ", err)
		log.Println("╚=========================================╝")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Something went wrong",
		})
	}
	defer stmt.Close()

	_, err = stmt.Exec(traqType.Name)
	if err != nil {
		log.Println("║ 💥 Failed to create traq type: ", err)
		log.Println("║ 🍺 Name: ", traqType.Name)
		log.Println("╚=========================================╝")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Something went wrong",
		})
	}

	log.Println("║ ✅ Traq type created successfully")
	log.Println("║ 🍺 Name: ", traqType.Name)
	log.Println("╚=========================================╝")
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Traq type created successfully",
	})

}

func getAllTraqTypes(c *fiber.Ctx) error {
	log.Println("╔======== 🍺 Get All Traq Types 🍺 ========╗")

	request := `
		SELECT name
		FROM traq_types;
	`

	rows, err := db.Query(request)
	if err != nil {
		log.Println("║ 💥 Failed to query traq types: ", err)
		log.Println("╚=========================================╝")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Something went wrong",
		})
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Println("║ 💥 Failed to close rows: ", err)
			log.Println("╚=========================================╝")
		}
	}(rows)

	var traqTypes []TraqType
	for rows.Next() {
		var traqType TraqType
		if err := rows.Scan(&traqType.Name); err != nil {
			log.Println("║ 💥 Failed to scan row: ", err)
			log.Println("╚=========================================╝")
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Something went wrong",
			})
		}
		traqTypes = append(traqTypes, traqType)
	}

	log.Println("║ ✅ Traq types retrieved successfully")
	log.Println("╚=========================================╝")

	return c.JSON(traqTypes)
}

func createTraqArticle(c *fiber.Ctx) error {
	log.Println("╔======== 🍺 Create Article 🍺 ========╗")
	var article TraqArticle
	if err := c.BodyParser(&article); err != nil {
		log.Println("║ 💥 Failed to parse request body: ", err)
		log.Println("╚=========================================╝")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	fmt.Println(article)

	if article.Name == "" || article.Description == "" || article.Picture == "" || article.TraqType == "" {
		log.Println("║ 💥 Some fields are missing")
		log.Println("╚=========================================╝")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Fill in all fields",
		})
	}

	request := `
		INSERT INTO traq (name, description, picture, id_traq_types, limited, price, disabled, alcohol, out_of_stock, price_half) 
		VALUES ($1, $2, $3, (SELECT id_traq_types FROM traq_types WHERE name = $4), $5, $6, $7, $8, $9, $10);
	`

	stmt, err := db.Prepare(request)
	if err != nil {
		log.Println("║ 💥 Failed to prepare statement: ", err)
		log.Println("╚=========================================╝")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Something went wrong",
		})
	}

	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			log.Println("║ 💥 Failed to close statement: ", err)
			log.Println("╚=========================================╝")
		}
	}(stmt)

	_, err = stmt.Exec(article.Name, article.Description, article.Picture, article.TraqType, article.Limited, article.Price, article.Disabled, article.Alcohol, article.OutOfStock, article.PriceHalf)
	if err != nil {
		log.Println("║ 💥 Failed to create article: ", err)
		log.Println("║ 🍺 Name: ", article.Name)
		log.Println("╚=========================================╝")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Something went wrong",
		})
	}

	log.Println("║ ✅ Article created successfully")
	log.Println("║ 🍺 Name: ", article.Name)
	log.Println("╚=========================================╝")

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Article created successfully",
	})
}

func getAllTraqArticles(c *fiber.Ctx) error {
	log.Println("╔======== 🍺 Get All Articles 🍺 ========╗")

	request := `
		SELECT id_traq, name, description, picture, price, price_half, alcohol, creation_date, limited, out_of_stock, disabled, (SELECT name FROM traq_types WHERE traq.id_traq_types = traq_types.id_traq_types) as traq_type
		FROM traq;
	`

	rows, err := db.Query(request)
	if err != nil {
		log.Println("║ 💥 Failed to query articles: ", err)
		log.Println("╚=========================================╝")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Something went wrong",
		})
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Println("║ 💥 Failed to close rows: ", err)
			log.Println("╚=========================================╝")
		}
	}(rows)

	var articles []TraqArticle
	for rows.Next() {
		var article TraqArticle
		if err := rows.Scan(&article.ID, &article.Name, &article.Description, &article.Picture, &article.Price, &article.PriceHalf, &article.Alcohol, &article.CreationDate, &article.Limited, &article.OutOfStock, &article.Disabled, &article.TraqType); err != nil {
			log.Println("║ 💥 Failed to scan row: ", err)
			log.Println("╚=========================================╝")
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Something went wrong",
			})
		}
		articles = append(articles, article)
	}

	log.Println("║ ✅ Articles retrieved successfully")
	log.Println("╚=========================================╝")

	return c.JSON(articles)
}
