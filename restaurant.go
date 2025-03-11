package main

import (
	"database/sql"
	"github.com/gofiber/fiber/v2"
	"log"
)

func getRestaurant(c *fiber.Ctx) error {
	log.Println("╔======== 🍽️ Get Restaurant 🍽️ ========╗")
	request := `
		SELECT restaurant.id_restaurant, restaurant.articles, restaurant.updated_date from restaurant
		ORDER BY updated_date DESC
		LIMIT 1;
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
			return
		}
	}(stmt)

	var restaurant Restaurant
	err = stmt.QueryRow().Scan(&restaurant.ID, &restaurant.Articles, &restaurant.UpdatedDate)
	if err != nil {
		log.Println("║ 💥 Failed to get restaurant: ", err)
		log.Println("╚=========================================╝")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Something went wrong",
		})
	}

	log.Println("║ ✅ Restaurant retrieved successfully")
	log.Println("╚=========================================╝")
	return c.JSON(restaurant)
}
