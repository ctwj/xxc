package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"moss/infrastructure/support/auth"
	"moss/infrastructure/support/log"
)

// JWTMiddleware validates JWT token from cookie or Authorization header
func JWTMiddleware(c *fiber.Ctx) error {
	// Try to get token from cookie first
	tokenString := c.Cookies("token")

	// If no cookie, try Authorization header
	if tokenString == "" {
		authHeader := c.Get("Authorization")
		if authHeader != "" {
			// Remove "Bearer " prefix if present
			tokenString = strings.TrimPrefix(authHeader, "Bearer ")
		}
	}

	if tokenString == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
			"message": "No authentication token provided",
		})
	}

	// Validate token
	claims, err := auth.ValidateToken(tokenString)
	if err != nil {
		log.Warn("JWT validation failed", log.Err(err))
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
			"message": "Invalid or expired token",
		})
	}

	// Store user info in context locals
	c.Locals("userID", claims.UserID)
	c.Locals("username", claims.Username)
	c.Locals("role", claims.Role)

	return c.Next()
}

// OptionalJWTMiddleware validates JWT token if present but doesn't require it
func OptionalJWTMiddleware(c *fiber.Ctx) error {
	// Try to get token from cookie first
	tokenString := c.Cookies("token")

	// If no cookie, try Authorization header
	if tokenString == "" {
		authHeader := c.Get("Authorization")
		if authHeader != "" {
			tokenString = strings.TrimPrefix(authHeader, "Bearer ")
		}
	}

	// No token, continue without user info
	if tokenString == "" {
		return c.Next()
	}

	// Validate token
	claims, err := auth.ValidateToken(tokenString)
	if err != nil {
		// Invalid token, continue without user info
		return c.Next()
	}

	// Store user info in context locals
	c.Locals("userID", claims.UserID)
	c.Locals("username", claims.Username)
	c.Locals("role", claims.Role)

	return c.Next()
}

// GetUserID retrieves user ID from context
func GetUserID(c *fiber.Ctx) uint {
	userID := c.Locals("userID")
	if userID == nil {
		return 0
	}
	return userID.(uint)
}

// GetUsername retrieves username from context
func GetUsername(c *fiber.Ctx) string {
	username := c.Locals("username")
	if username == nil {
		return ""
	}
	return username.(string)
}

// GetRole retrieves user role from context
func GetRole(c *fiber.Ctx) string {
	role := c.Locals("role")
	if role == nil {
		return ""
	}
	return role.(string)
}