package controller

import (
	"time"

	"github.com/gofiber/fiber/v2"
	appService "moss/application/service"
	"moss/infrastructure/support/auth"
)

// AuthLogin handles user login
func AuthLogin(c *fiber.Ctx) error {
	type LoginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
	}

	// Login using application service (uses config-based admin)
	token, err := appService.AdminLogin(req.Username, req.Password, "", "")
	if err != nil {
		return c.Status(401).JSON(fiber.Map{
			"success": false,
			"error":   "invalid credentials",
		})
	}

	// Set cookie
	c.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    token,
		Expires:  time.Now().Add(auth.JWTExpireTime),
		HTTPOnly: true,
		Secure:   false, // Set to true in production with HTTPS
		SameSite: "Lax",
		Path:     "/",
	})

	return c.JSON(fiber.Map{
		"success": true,
		"user": fiber.Map{
			"username": req.Username,
			"role":     "admin",
		},
	})
}

// AuthRegister handles user registration (disabled for single-admin system)
func AuthRegister(c *fiber.Ctx) error {
	return c.Status(403).JSON(fiber.Map{
		"success": false,
		"error":   "registration is disabled. This system uses a single admin account.",
	})
}

// AuthLogout handles user logout
func AuthLogout(c *fiber.Ctx) error {
	// Clear the token cookie
	c.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
		Path:     "/",
	})

	return c.JSON(fiber.Map{"success": true})
}

// AuthMe returns the current authenticated user
func AuthMe(c *fiber.Ctx) error {
	userID := GetUserID(c)
	if userID == 0 {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}

	username := GetUsername(c)

	return c.JSON(fiber.Map{
		"id":       userID,
		"username": username,
		"role":     "admin",
	})
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
