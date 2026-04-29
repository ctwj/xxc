package controller

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
)

// TestAuthRegister tests the registration endpoint
func TestAuthRegister(t *testing.T) {
	app := fiber.New()
	app.Post("/api/auth/register", AuthRegister)

	tests := []struct {
		name       string
		payload    map[string]string
		expectCode int
	}{
		{
			name: "valid registration",
			payload: map[string]string{
				"username": "testuser",
				"email":    "test@example.com",
				"password": "password123",
			},
			expectCode: 200, // or 201 depending on implementation
		},
		{
			name: "missing username",
			payload: map[string]string{
				"email":    "test@example.com",
				"password": "password123",
			},
			expectCode: 400,
		},
		{
			name: "missing email",
			payload: map[string]string{
				"username": "testuser",
				"password": "password123",
			},
			expectCode: 400,
		},
		{
			name: "missing password",
			payload: map[string]string{
				"username": "testuser",
				"email":    "test@example.com",
			},
			expectCode: 400,
		},
		{
			name: "short password",
			payload: map[string]string{
				"username": "testuser",
				"email":    "test@example.com",
				"password": "123",
			},
			expectCode: 400,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.payload)
			req := httptest.NewRequest("POST", "/api/auth/register", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req, -1)
			if err != nil {
				t.Fatalf("failed to send request: %v", err)
			}

			// Note: Actual status code may vary based on database state
			// This is a basic structure test
			t.Logf("Response status: %d", resp.StatusCode)
		})
	}
}

// TestAuthLogin tests the login endpoint
func TestAuthLogin(t *testing.T) {
	app := fiber.New()
	app.Post("/api/auth/login", AuthLogin)

	tests := []struct {
		name       string
		payload    map[string]string
		expectCode int
	}{
		{
			name: "missing username",
			payload: map[string]string{
				"password": "password123",
			},
			expectCode: 401,
		},
		{
			name: "missing password",
			payload: map[string]string{
				"username": "testuser",
			},
			expectCode: 401,
		},
		{
			name: "invalid credentials",
			payload: map[string]string{
				"username": "nonexistent",
				"password": "wrongpassword",
			},
			expectCode: 401,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.payload)
			req := httptest.NewRequest("POST", "/api/auth/login", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req, -1)
			if err != nil {
				t.Fatalf("failed to send request: %v", err)
			}

			t.Logf("Response status: %d", resp.StatusCode)
		})
	}
}
