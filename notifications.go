package main

import (
	"github.com/go-resty/resty/v2"
	"github.com/gofiber/fiber/v2"
	"log"
)

var tokenGroups = make(map[string][]string)

func registerToken(c *fiber.Ctx) error {
	var body struct {
		Token string `json:"token"`
		Group string `json:"group"`
	}
	if err := c.BodyParser(&body); err != nil {
		log.Println("Error parsing body:", err)
		return c.Status(400).SendString("Invalid request")
	}

	tokenGroups[body.Group] = append(tokenGroups[body.Group], body.Token)
	log.Println("Registered token:", body.Token, "for group:", body.Group)
	return c.SendString("Token registered")
}

func sendNotification(c *fiber.Ctx) error {
	var body struct {
		Group   string `json:"group"`
		Title   string `json:"title"`
		Message string `json:"message"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).SendString("Invalid request")
	}

	tokens, exists := tokenGroups[body.Group]
	if !exists || len(tokens) == 0 {
		return c.Status(400).SendString("No users in this group")
	}

	client := resty.New()
	for _, token := range tokens {
		_, err := client.R().
			SetHeader("Content-Type", "application/json").
			SetBody(map[string]interface{}{
				"to":    token,
				"title": body.Title,
				"body":  body.Message,
			}).
			Post("https://exp.host/--/api/v2/push/send")

		if err != nil {
			log.Println("Error sending notification:", err)
		}
	}

	return c.SendString("Notification sent")
}
