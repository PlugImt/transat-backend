package middlewares

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// Common field limitations with reasonable limits for different types of data
var fieldLimits = map[string]int{
	// Personal information
	"email":            254, // Max length according to RFC 5321
	"password":         72,  // BCrypt has a 72-byte limit
	"first_name":       50,
	"last_name":        50,
	"phone_number":     20,
	"campus":           100,
	"language":         10,
	
	// Content fields
	"title":            100,
	"content":          500,
	"description":      300,
	"message":          500,
	"comment":          250,
	
	// Tokens and codes
	"verification_code":  10,
	"notification_token": 255,
	
	// Default limit for any unspecified field
	"default": 500,
}

// ValidationMiddleware checks request body size and field lengths
func ValidationMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Skip validation for GET and DELETE requests which don't have a body
		if c.Method() == "GET" || c.Method() == "DELETE" {
			return c.Next()
		}
		
		// Parse the request body into a map
		var requestBody map[string]interface{}
		if err := json.Unmarshal(c.Body(), &requestBody); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid JSON format",
			})
		}
		
		// Validate field lengths
		for field, value := range requestBody {
			// Skip validation for non-string fields
			strValue, ok := value.(string)
			if !ok {
				continue
			}
			
			// Get the limit for this field or use default
			limit := fieldLimits["default"]
			// Check for field-specific limits using exact match or pattern matching
			if specificLimit, exists := fieldLimits[field]; exists {
				limit = specificLimit
			} else {
				// Try pattern matching for fields not explicitly defined
				for pattern, patternLimit := range fieldLimits {
					if strings.Contains(field, pattern) {
						limit = patternLimit
						break
					}
				}
			}
			
			// Check if the field exceeds the limit
			if len(strValue) > limit {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": fmt.Sprintf("Field '%s' exceeds maximum length of %d characters", field, limit),
				})
			}
		}
		
		return c.Next()
	}
} 