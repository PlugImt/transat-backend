package middlewares

import (
	"github.com/gofiber/fiber/v2"
)

// SecurityHeadersMiddleware adds important security headers to all responses
func SecurityHeadersMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Content-Security-Policy - controls which resources can be loaded
		c.Set("Content-Security-Policy", "default-src 'self'; img-src 'self' data:; script-src 'self'; connect-src 'self'; frame-ancestors 'none'")
		
		// X-Content-Type-Options - prevents MIME type sniffing
		c.Set("X-Content-Type-Options", "nosniff")
		
		// X-Frame-Options - prevents clickjacking
		c.Set("X-Frame-Options", "DENY")
		
		// X-XSS-Protection - enables XSS filtering in browsers
		c.Set("X-XSS-Protection", "1; mode=block")
		
		// Strict-Transport-Security - enforces HTTPS
		c.Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
		
		// Referrer-Policy - controls referrer information
		c.Set("Referrer-Policy", "no-referrer-when-downgrade")
		
		// Feature-Policy - restricts browser features
		c.Set("Permissions-Policy", "geolocation=(), microphone=(), camera=()")
		
		// Cache-Control - prevents caching of API responses
		c.Set("Cache-Control", "no-store, no-cache, must-revalidate, proxy-revalidate")
		c.Set("Pragma", "no-cache")
		c.Set("Expires", "0")
		
		return c.Next()
	}
} 