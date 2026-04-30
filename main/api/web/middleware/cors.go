package middleware

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"moss/domain/config"
)

// CORSConfig returns CORS middleware handler
func CORSConfig() fiber.Handler {
	return func(c *fiber.Ctx) error {
		origin := c.Get("Origin")

		// If no origin, skip CORS
		if origin == "" {
			return c.Next()
		}

		// Check if origin is allowed
		allowed := isOriginAllowed(origin)
		fmt.Printf("[CORS] origin=%s, allowed=%v\n", origin, allowed)

		// Handle preflight request
		if c.Method() == fiber.MethodOptions {
			if allowed {
				c.Set("Access-Control-Allow-Origin", origin)
				c.Set("Access-Control-Allow-Credentials", "true")
				c.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
				c.Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, X-Requested-With")
				c.Set("Access-Control-Max-Age", "86400")
				c.Set("Access-Control-Expose-Headers", "Set-Cookie")
				fmt.Printf("[CORS] Preflight headers set for origin=%s\n", origin)
			}
			c.Vary("Origin")
			c.Vary("Access-Control-Request-Method")
			c.Vary("Access-Control-Request-Headers")
			return c.SendStatus(fiber.StatusNoContent)
		}

		// Handle actual request
		if allowed {
			c.Set("Access-Control-Allow-Origin", origin)
			c.Set("Access-Control-Allow-Credentials", "true")
			c.Set("Access-Control-Expose-Headers", "Set-Cookie")
		}
		c.Vary("Origin")

		return c.Next()
	}
}

// isOriginAllowed checks if the origin is allowed
func isOriginAllowed(origin string) bool {
	// Get configured origins
	origins := config.Config.Router.CORSOrigins

	// If no origins configured, allow localhost for development
	if origins == "" {
		// Allow any localhost origin for development
		if strings.HasPrefix(origin, "http://localhost:") ||
			strings.HasPrefix(origin, "http://127.0.0.1:") {
			return true
		}
		return false
	}

	// Check against configured origins
	for _, o := range strings.Split(origins, ",") {
		if strings.TrimSpace(o) == origin {
			return true
		}
	}

	// Also allow localhost for development even if other origins are configured
	if strings.HasPrefix(origin, "http://localhost:") ||
		strings.HasPrefix(origin, "http://127.0.0.1:") {
		return true
	}

	return false
}
