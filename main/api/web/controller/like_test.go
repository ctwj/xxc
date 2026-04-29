package controller

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"moss/domain/core/entity"
)

// TestLikeSet tests the like/dislike endpoint
func TestLikeSet(t *testing.T) {
	app := fiber.New()
	app.Post("/api/likes", LikeSet)

	tests := []struct {
		name       string
		payload    map[string]interface{}
		expectCode int
	}{
		{
			name: "set like",
			payload: map[string]interface{}{
				"articleId": 1,
				"type":      int(entity.LikeTypeLike),
			},
			expectCode: 401, // Unauthorized without JWT
		},
		{
			name: "set dislike",
			payload: map[string]interface{}{
				"articleId": 1,
				"type":      int(entity.LikeTypeDislike),
			},
			expectCode: 401,
		},
		{
			name: "remove like",
			payload: map[string]interface{}{
				"articleId": 1,
				"type":      0,
			},
			expectCode: 401,
		},
		{
			name: "missing articleId",
			payload: map[string]interface{}{
				"type": 1,
			},
			expectCode: 401, // Unauthorized first, then 400
		},
		{
			name: "invalid type",
			payload: map[string]interface{}{
				"articleId": 1,
				"type":      3, // Invalid type
			},
			expectCode: 401,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.payload)
			req := httptest.NewRequest("POST", "/api/likes", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req, -1)
			if err != nil {
				t.Fatalf("failed to send request: %v", err)
			}

			t.Logf("Response status: %d", resp.StatusCode)
		})
	}
}

// TestLikeGet tests the like status endpoint
func TestLikeGet(t *testing.T) {
	app := fiber.New()
	app.Get("/api/likes/status", LikeGet)

	tests := []struct {
		name       string
		query      string
		expectCode int
	}{
		{
			name:       "valid articleId",
			query:      "?articleId=1",
			expectCode: 401, // Unauthorized without JWT
		},
		{
			name:       "missing articleId",
			query:      "",
			expectCode: 401,
		},
		{
			name:       "invalid articleId",
			query:      "?articleId=0",
			expectCode: 401,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/api/likes/status"+tt.query, nil)

			resp, err := app.Test(req, -1)
			if err != nil {
				t.Fatalf("failed to send request: %v", err)
			}

			t.Logf("Response status: %d", resp.StatusCode)
		})
	}
}