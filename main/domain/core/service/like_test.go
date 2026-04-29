package service

import (
	"testing"

	"moss/domain/core/entity"
)

func TestLikeService_SetLike(t *testing.T) {
	// Note: This test requires a database connection
	// These are unit tests for the service logic

	t.Run("like type validation", func(t *testing.T) {
		// Test that LikeType constants are correct
		if entity.LikeTypeNone != 0 {
			t.Errorf("expected LikeTypeNone to be 0, got %d", entity.LikeTypeNone)
		}
		if entity.LikeTypeLike != 1 {
			t.Errorf("expected LikeTypeLike to be 1, got %d", entity.LikeTypeLike)
		}
		if entity.LikeTypeDislike != 2 {
			t.Errorf("expected LikeTypeDislike to be 2, got %d", entity.LikeTypeDislike)
		}
	})
}

func TestLikeService_GetUserLikeType(t *testing.T) {
	// Test that GetUserLikeType returns LikeTypeNone for non-existent like
	// This test doesn't require database since it should return 0 for errors
	result := Like.GetUserLikeType(999999, 999999)
	if result != entity.LikeTypeNone {
		t.Errorf("expected LikeTypeNone for non-existent like, got %d", result)
	}
}
