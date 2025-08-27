package middlewares

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

// GeneralRateLimiter is a middleware that limits repeated requests to all API endpoints
var GeneralRateLimiter = limiter.New(limiter.Config{
	Max:        1000,            // 100 requests
	Expiration: 1 * time.Minute, // per minute
	KeyGenerator: func(c *fiber.Ctx) string {
		// Use both IP and path to create more targeted rate limiting
		return c.IP() + "-" + c.Path()
	},
	LimitReached: func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
			"error": "Rate limit exceeded. Please try again later.",
		})
	},
})

// Create more strict rate limiters for sensitive endpoints

// AccountRateLimiter protects account-related endpoints
var AccountRateLimiter = limiter.New(limiter.Config{
	Max:        50,              // 5 requests
	Expiration: 1 * time.Minute, // per minute
	KeyGenerator: func(c *fiber.Ctx) string {
		return c.IP() + "-account"
	},
	LimitReached: func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
			"error": "Too many requests to account endpoints. Please try again later.",
		})
	},
})

// AccessRateLimiter for public access points (even more strict)
var AccessRateLimiter = limiter.New(limiter.Config{
	Max:        100,             // 10 requests
	Expiration: 5 * time.Minute, // per 5 minutes
	KeyGenerator: func(c *fiber.Ctx) string {
		return c.IP() + "-access"
	},
	LimitReached: func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
			"error": "Too many access attempts. Please try again later.",
		})
	},
	SkipFailedRequests:     false, // Count failed requests against the limit
	SkipSuccessfulRequests: false, // Count successful requests against the limit
})
