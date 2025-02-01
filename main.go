package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/proxy"
	"log"
	"time"
)

func main() {

	app := fiber.New()

	// Update CORS configuration to allow WebSocket
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization, X-Turnstile-Token",
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
	}))

	loginRegisterLimiter := limiter.New(limiter.Config{
		Max:        3,
		Expiration: 300 * time.Second,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "Too many requests",
			})
		},
	})

	app.Get("/status", func(c *fiber.Ctx) error {
		return c.SendString("API is up and running")
	})

	// Route to user service
	app.All("/api/users/*", loginRegisterLimiter, func(c *fiber.Ctx) error {
		return proxy.Do(c, "http://localhost:3001"+c.Path())
	})

	// Route to food service
	app.All("/api/food/*", func(c *fiber.Ctx) error {
		return proxy.Do(c, "http://localhost:3002"+c.Path())
	})

	log.Fatal(app.Listen(":3000"))
}
