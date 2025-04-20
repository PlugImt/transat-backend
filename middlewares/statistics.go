package middlewares

import (
	"time"

	"Transat_2.0_Backend/services"
	"github.com/gofiber/fiber/v2"
)

// StatisticsMiddleware captures request timing and logs statistics
func StatisticsMiddleware(statisticsService *services.StatisticsService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Record request time
		requestReceived := time.Now()
		
		// Store original URL path
		endpoint := c.Path()
		method := c.Method()
		
		// Process request
		err := c.Next()
		
		// Record response time
		responseSent := time.Now()
		
		// Extract user email if available
		userEmail := ""
		if email, ok := c.Locals("email").(string); ok && email != "" {
			userEmail = email
		}
		
		// Log the request to the statistics service
		go statisticsService.LogRequest(
			userEmail,
			endpoint,
			method,
			requestReceived,
			responseSent,
			c.Response().StatusCode(),
		)
		
		return err
	}
} 