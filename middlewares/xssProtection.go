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

		// Function to check if a string is a URL
		isURL := func(input string) bool {
			// Simple URL pattern check
			urlPattern := regexp.MustCompile(`^(https?:\/\/)?(www\.)?[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b([-a-zA-Z0-9()@:%_\+.~#?&//=]*)$`)
			return urlPattern.MatchString(input) || strings.HasPrefix(input, "/api/data/")
		}

		// Function to sanitize input from XSS
		sanitizeXSS := func(input string, isURLField bool) string {
			// Skip sanitizing URLs
			if isURLField {
				return input
			}

			// Replace potentially dangerous characters
			input = strings.ReplaceAll(input, "<", "&lt;")
			input = strings.ReplaceAll(input, ">", "&gt;")
			input = strings.ReplaceAll(input, "\"", "&quot;")
			input = strings.ReplaceAll(input, "'", "&#x27;")
			input = strings.ReplaceAll(input, "(", "&#40;")
			input = strings.ReplaceAll(input, ")", "&#41;")
			input = strings.ReplaceAll(input, "/", "&#x2F;")
			input = regexp.MustCompile(`(?i)javascript:`).ReplaceAllString(input, "")
			input = regexp.MustCompile(`(?i)data:`).ReplaceAllString(input, "")

			return input
		}

		// Check and sanitize all string values for XSS
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

			// Check if this field contains a URL
			isURLField := isURL(strValue) ||
				strings.Contains(strings.ToLower(field), "url") ||
				strings.Contains(strings.ToLower(field), "link") ||
				strings.Contains(strings.ToLower(field), "image") ||
				strings.Contains(strings.ToLower(field), "picture") ||
				strings.Contains(strings.ToLower(field), "avatar") ||
				strings.Contains(strings.ToLower(field), "photo")

			// Sanitize the input, but don't escape URLs
			requestBody[field] = sanitizeXSS(strValue, isURLField)
		}

		// Replace the request body with the sanitized one
		sanitizedBody, _ := json.Marshal(requestBody)
		c.Request().SetBody(sanitizedBody)

		return c.Next()
	}
}
