package service

import (
	"testing"
)

func TestUserService_Register(t *testing.T) {
	// Note: This test requires a database connection
	// In a real test environment, you would use a test database or mock

	t.Run("empty username", func(t *testing.T) {
		_, err := User.Register("", "test@example.com", "password123")
		if err == nil {
			t.Error("expected error for empty username")
		}
	})

	t.Run("empty email", func(t *testing.T) {
		_, err := User.Register("testuser", "", "password123")
		if err == nil {
			t.Error("expected error for empty email")
		}
	})

	t.Run("empty password", func(t *testing.T) {
		_, err := User.Register("testuser", "test@example.com", "")
		if err == nil {
			t.Error("expected error for empty password")
		}
	})
}

func TestUserService_Login(t *testing.T) {
	t.Run("empty credentials", func(t *testing.T) {
		_, err := User.Login("", "password")
		if err == nil {
			t.Error("expected error for empty username")
		}
	})
}
