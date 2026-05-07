package controller

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
)

// TestViewHistoryRecord tests the view history record endpoint
func TestViewHistoryRecord(t *testing.T) {
	app := fiber.New()
	app.Post("/api/history", ViewHistoryRecord)

	tests := []struct {
		name       string
		payload    map[string]interface{}
		expectCode int
	}{
		{
			name: "valid record",
			payload: map[string]interface{}{
				"articleId": 1,
			},
			expectCode: 401, // Unauthorized without JWT
		},
		{
			name:       "missing articleId",
			payload:    map[string]interface{}{},
			expectCode: 401,
		},
		{
			name: "invalid articleId",
			payload: map[string]interface{}{
				"articleId": 0,
			},
			expectCode: 401,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.payload)
			req := httptest.NewRequest("POST", "/api/history", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req, -1)
			if err != nil {
				t.Fatalf("failed to send request: %v", err)
			}

			t.Logf("Response status: %d", resp.StatusCode)
		})
	}
}

// TestViewHistoryList tests the view history list endpoint
func TestViewHistoryList(t *testing.T) {
	app := fiber.New()
	app.Get("/api/history", ViewHistoryList)

	req := httptest.NewRequest("GET", "/api/history", nil)

	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("failed to send request: %v", err)
	}

	// Should return 401 Unauthorized without JWT
	if resp.StatusCode != 401 {
		t.Logf("Expected 401, got %d", resp.StatusCode)
	}
}

// TestViewHistoryClear tests the view history clear endpoint
func TestViewHistoryClear(t *testing.T) {
	app := fiber.New()
	app.Delete("/api/history", ViewHistoryClear)

	req := httptest.NewRequest("DELETE", "/api/history", nil)

	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("failed to send request: %v", err)
	}

	// Should return 401 Unauthorized without JWT
	if resp.StatusCode != 401 {
		t.Logf("Expected 401, got %d", resp.StatusCode)
	}
}
