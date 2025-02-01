package main

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"io"
	"log"
	"net/http"
	"strings"
)

func TurnstileMiddleware() fiber.Handler {
	secretKey := "0x4AAAAAAA4L7ghKZum5gV9UyrY85Fg-C7g"

	return func(c *fiber.Ctx) error {
		token := c.Get("X-Turnstile-Token")
		if token == "" {
			log.Printf("Missing Turnstile token")
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Missing Turnstile verification",
			})
		}

		formData := fmt.Sprintf("secret=%s&response=%s", secretKey, token)

		resp, err := http.Post(
			"https://challenges.cloudflare.com/turnstile/v0/siteverify",
			"application/x-www-form-urlencoded",
			strings.NewReader(formData),
		)
		if err != nil {
			log.Printf("Error making request to Turnstile API: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to verify Turnstile token",
			})
		}
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				log.Printf("Error closing response body: %v", err)
			}
		}(resp.Body)

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Error reading Turnstile response: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to verify Turnstile token",
			})
		}

		var result TurnstileResponse
		if err := json.Unmarshal(body, &result); err != nil {
			log.Printf("Error parsing Turnstile response: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to verify Turnstile token",
			})
		}

		if !result.Success {
			log.Printf("Invalid Turnstile token. Errors: %v", result.Error)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid Turnstile verification",
			})
		}

		return c.Next()
	}
}

func authMiddleware(c *fiber.Ctx) error {
	token := c.Get("Authorization")
	if token == "" {
		return c.Status(401).JSON(fiber.Map{
			"error": "No token provided",
		})
	}

	// Verify JWT token (implementation omitted)
	if !verifyJWT(token) {
		return c.Status(401).JSON(fiber.Map{
			"error": "Invalid token",
		})
	}

	return c.Next()
}

func verifyJWT(token string) bool {
	// JWT verification logic would go here
	// This is just a placeholder that always returns true
	return true
}
