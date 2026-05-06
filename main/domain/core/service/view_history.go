package service

import (
	"time"

	"moss/domain/core/entity"
	"moss/domain/core/repository"
	"moss/domain/core/repository/context"
)

type ViewHistoryService struct{}

var ViewHistory = &ViewHistoryService{}

// Record records a view for an article
func (s *ViewHistoryService) Record(userID uint, articleID int) error {
	// Check if already viewed
	existing, err := repository.ViewHistory.GetByUserIDAndArticleID(userID, articleID)
	if err != nil {
		// Not found, create new
		history := &entity.ViewHistory{
			UserID:    userID,
			ArticleID: articleID,
			ViewedAt:  time.Now(),
		}
		return repository.ViewHistory.Create(history)
	}

	// Update viewed_at timestamp
	return repository.ViewHistory.UpdateViewedAt(existing.ID, time.Now())
}

// Delete deletes a view history by ID
func (s *ViewHistoryService) Delete(id uint) error {
	return repository.ViewHistory.Delete(id)
}

// DeleteByUserID deletes all view history for a user
func (s *ViewHistoryService) DeleteByUserID(userID uint) error {
	return repository.ViewHistory.DeleteByUserID(userID)
}

// GetByUserIDAndArticleID gets a view history by user ID and article ID
func (s *ViewHistoryService) GetByUserIDAndArticleID(userID uint, articleID int) (*entity.ViewHistory, error) {
	return repository.ViewHistory.GetByUserIDAndArticleID(userID, articleID)
}

// ListByUserID lists all view history for a user
func (s *ViewHistoryService) ListByUserID(ctx *context.Context, userID uint) ([]entity.ViewHistory, error) {
	return repository.ViewHistory.ListByUserID(ctx, userID)
}

// CountByUserID counts view history for a user
func (s *ViewHistoryService) CountByUserID(userID uint) (int64, error) {
	return repository.ViewHistory.CountByUserID(userID)
}