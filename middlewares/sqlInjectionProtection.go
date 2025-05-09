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

		// List of SQL injection patterns to check for
		sqlInjectionPatterns := []*regexp.Regexp{
			regexp.MustCompile(`(?i)(\s|^)(SELECT|INSERT|UPDATE|DELETE|DROP|ALTER|CREATE|TRUNCATE|UNION|JOIN|WHERE|AND|OR|LIKE|HAVING|GROUP\s+BY|ORDER\s+BY)(\s|$)`),
			regexp.MustCompile(`(?i)'(''|[^'])*'`),                              // SQL strings
			regexp.MustCompile(`(?i);`),                                         // SQL statement terminator
			regexp.MustCompile(`(?i)--`),                                        // SQL comment
			regexp.MustCompile(`(?i)\/\*.*?\*\/`),                               // SQL block comment
			regexp.MustCompile(`(?i)(\s|^)(EXEC|EXECUTE|EXEC\(|EXECUTE\()(\s|$)`), // SQL execution commands
			regexp.MustCompile(`(?i)xp_`),                                       // SQL Server extended stored procedures
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

		// Function to sanitize input
		sanitizeInput := func(input string) string {
			// Replace single quotes with two single quotes (SQL escape)
			input = strings.ReplaceAll(input, "'", "''")
			// Remove comments and multiple spaces
			input = regexp.MustCompile(`/\*.*?\*/`).ReplaceAllString(input, "")
			input = regexp.MustCompile(`--.*`).ReplaceAllString(input, "")
			input = regexp.MustCompile(`\s+`).ReplaceAllString(input, " ")
			return input
		}

		// Check all string values for SQL injection
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

			// Sanitize the input
			requestBody[field] = sanitizeInput(strValue)
		}

		// Replace the request body with the sanitized one
		sanitizedBody, _ := json.Marshal(requestBody)
		c.Request().SetBody(sanitizedBody)

		return c.Next()
	}
} 