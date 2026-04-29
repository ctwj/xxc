package controller

import (
	"time"

	"github.com/gofiber/fiber/v2"
	appService "moss/application/service"
	"moss/domain/core/service"
	"moss/infrastructure/support/auth"
)

// AuthLogin handles user login (supports both admin and regular users)
func AuthLogin(c *fiber.Ctx) error {
	type LoginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
	}

	// First try admin login (using application service)
	token, err := appService.AdminLogin(req.Username, req.Password, "", "")
	if err == nil {
		// Admin login successful
		c.Cookie(&fiber.Cookie{
			Name:     "token",
			Value:    token,
			Expires:  time.Now().Add(auth.JWTExpireTime),
			HTTPOnly: true,
			Secure:   false,
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

	// Try regular user login
	user, err := service.User.Login(req.Username, req.Password)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{
			"success": false,
			"error":   "invalid credentials",
		})
	}

	// Generate JWT token for user
	token, err = auth.GenerateJWTToken(user.ID, user.Username, user.Role)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to generate token"})
	}

	c.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    token,
		Expires:  time.Now().Add(auth.JWTExpireTime),
		HTTPOnly: true,
		Secure:   false,
		SameSite: "Lax",
		Path:     "/",
	})

	return c.JSON(fiber.Map{
		"success": true,
		"user": fiber.Map{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
			"role":     user.Role,
		},
	})
}

// AuthRegister handles user registration
func AuthRegister(c *fiber.Ctx) error {
	type RegisterRequest struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var req RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
	}

	// Validate input
	if req.Username == "" || req.Email == "" || req.Password == "" {
		return c.Status(400).JSON(fiber.Map{"error": "username, email, and password are required"})
	}

	if len(req.Password) < 6 {
		return c.Status(400).JSON(fiber.Map{"error": "password must be at least 6 characters"})
	}

	// Register user
	user, err := service.User.Register(req.Username, req.Email, req.Password)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	// Generate JWT token
	token, err := auth.GenerateJWTToken(user.ID, user.Username, user.Role)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to generate token"})
	}

	c.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    token,
		Expires:  time.Now().Add(auth.JWTExpireTime),
		HTTPOnly: true,
		Secure:   false,
		SameSite: "Lax",
		Path:     "/",
	})

	return c.JSON(fiber.Map{
		"success": true,
		"user": fiber.Map{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
			"role":     user.Role,
		},
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
	role := GetRole(c)

	return c.JSON(fiber.Map{
		"id":       userID,
		"username": username,
		"role":     role,
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

// GetRole retrieves role from context
func GetRole(c *fiber.Ctx) string {
	role := c.Locals("role")
	if role == nil {
		return ""
	}
	return role.(string)
}
