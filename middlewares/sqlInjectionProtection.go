package middlewares

import (
	"encoding/json"
	"regexp"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// SQLInjectionProtectionMiddleware sanitizes request inputs to prevent SQL injection
func SQLInjectionProtectionMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Skip for GET and DELETE requests
		if c.Method() == "GET" || c.Method() == "DELETE" {
			return c.Next()
		}

		// Skip for routes that don't require a body
		path := c.Path()
		if strings.HasSuffix(path, "/join") || strings.HasSuffix(path, "/leave") {
			return c.Next()
		}

		// If body is empty, skip validation
		if len(c.Body()) == 0 {
			return c.Next()
		}

		// List of SQL injection patterns to check for (reduced to avoid false positives)
		sqlInjectionPatterns := []*regexp.Regexp{
			regexp.MustCompile(`(?i)UNION\s+SELECT`), // classic union-based injection
			regexp.MustCompile(`--`),                 // inline comment
			regexp.MustCompile(`/\*.*?\*/`),          // block comment
			regexp.MustCompile(`;\s*`),               // statement terminator
		}

		// Parse the request body into a map
		var requestBody map[string]interface{}
		if err := json.Unmarshal(c.Body(), &requestBody); err != nil {
			// If body isn't valid JSON, just pass it along (other middleware will handle this)
			return c.Next()
		}

		// Function to check for SQL injection in a string
		containsSQLInjection := func(input string) bool {
			for _, pattern := range sqlInjectionPatterns {
				if pattern.MatchString(input) {
					return true
				}
			}
			return false
		}

			// Replace single quotes with two single quotes (SQL escape)
		// No sanitization: we rely on parameterized queries in handlers/repositories

		// Check all string values for SQL injection (do not mutate content)
		for field, value := range requestBody {
			strValue, ok := value.(string)
			if !ok {
				continue
			}

			// Check for SQL injection in the string value
			if containsSQLInjection(strValue) {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": "Invalid input detected",
				})
			}
			_ = field
		}
		// Do not modify the request body; pass through as-is
		return c.Next()
	}
}
