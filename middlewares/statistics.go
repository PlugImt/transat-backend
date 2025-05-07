package middlewares

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/plugimt/transat-backend/services"
	"github.com/plugimt/transat-backend/utils"
)

// StatisticsMiddleware captures request timing and logs statistics
func StatisticsMiddleware(statisticsService *services.StatisticsService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Record request time
		requestReceived := time.Now()

		// Store original URL path
		endpoint := c.Path()
		method := c.Method()

		// Process request and capture any error
		err := c.Next()

		// Record response time
		responseSent := time.Now()

		// Extract user email if available
		userEmail := ""
		if email, ok := c.Locals("email").(string); ok && email != "" {
			userEmail = email
		}

		// Get the actual status code from the response
		statusCode := c.Response().StatusCode()

		// Handle 404s and other errors correctly
		// If Fiber sets the status code properly, use it, otherwise derive from error
		if err != nil {
			errMessage := fmt.Sprintf("%v", err)
			log.Printf("Error during request processing: %s", errMessage)
			if statusCode == 200 {
				// If status code wasn't set but there was an error, use 500 as default
				statusCode = http.StatusInternalServerError
				utils.LogMessage(utils.LevelWarn, "Error occurred but status code was 200, adjusting to 500")
			}
		}

		// For 404 routes specifically
		if string(c.Response().Body()) == fmt.Sprintf("Cannot %s %s", method, endpoint) {
			if statusCode == 200 {
				statusCode = http.StatusNotFound
				utils.LogMessage(utils.LevelWarn, "404 detected from response body but status code was 200, adjusting to 404")
			}
		}

		// Debug logging for status code
		utils.LogHeader("📊 Statistics Middleware Debug")
		utils.LogLineKeyValue(utils.LevelDebug, "Path", endpoint)
		utils.LogLineKeyValue(utils.LevelDebug, "Method", method)
		utils.LogLineKeyValue(utils.LevelDebug, "Status Code", statusCode)
		utils.LogLineKeyValue(utils.LevelDebug, "Response Body", string(c.Response().Body()))
		utils.LogFooter()

		// Log the request to the statistics service
		statisticsService.LogRequest(
			userEmail,
			endpoint,
			method,
			requestReceived,
			responseSent,
			statusCode,
		)

		return err
	}
}
