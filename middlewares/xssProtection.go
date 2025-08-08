package middlewares

import (
	"encoding/json"
	"regexp"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// XSSProtectionMiddleware sanitizes request inputs to prevent XSS attacks
func XSSProtectionMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Skip for GET and DELETE requests
		if c.Method() == "GET" || c.Method() == "DELETE" {
			return c.Next()
		}

		// Skip for multipart/form-data requests (file uploads)
		contentType := c.Get("Content-Type")
		if strings.HasPrefix(contentType, "multipart/form-data") {
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

		// XSS attack patterns to check for
		xssPatterns := []*regexp.Regexp{
			regexp.MustCompile(`(?i)<script.*?>.*?</script.*?>`),
			regexp.MustCompile(`(?i)javascript:`),
			regexp.MustCompile(`(?i)on\w+\s*=`),
			regexp.MustCompile(`(?i)<iframe.*?>.*?</iframe.*?>`),
			regexp.MustCompile(`(?i)<img.*?onerror.*?=`),
			regexp.MustCompile(`(?i)data:.*?base64`),
			regexp.MustCompile(`(?i)document\.cookie`),
			regexp.MustCompile(`(?i)document\.location`),
			regexp.MustCompile(`(?i)document\.write`),
			regexp.MustCompile(`(?i)document\.evaluate`),
			regexp.MustCompile(`(?i)window\.`),
			regexp.MustCompile(`(?i)eval\(`),
			regexp.MustCompile(`(?i)setTimeout\(`),
			regexp.MustCompile(`(?i)setInterval\(`),
			regexp.MustCompile(`(?i)new\s+Function\(`),
		}

		// Parse the request body into a map
		var requestBody map[string]interface{}
		if err := json.Unmarshal(c.Body(), &requestBody); err != nil {
			// If body isn't valid JSON, just pass it along
			return c.Next()
		}

		// Function to check for XSS attacks
		containsXSS := func(input string) bool {
			for _, pattern := range xssPatterns {
				if pattern.MatchString(input) {
					return true
				}
			}
			return false
		}

		for field, value := range requestBody {
			strValue, ok := value.(string)
			if !ok {
				continue
			}

			// Check for XSS in the string value
			if containsXSS(strValue) {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": "Invalid content detected",
				})
			}
			_ = field
		}

		return c.Next()
	}
}
