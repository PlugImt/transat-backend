package main

import (
	"github.com/gofiber/fiber/v2"
)

func login(c *fiber.Ctx) error {
	return c.SendString("User Logged In")
}
