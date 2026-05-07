package repository

import (
	"time"

	"moss/domain/core/entity"
	"moss/domain/core/repository/context"
	"moss/infrastructure/persistent/db"
)

type ViewHistoryRepository struct{}

var ViewHistory = &ViewHistoryRepository{}

// Create creates a new view history
func (r *ViewHistoryRepository) Create(history *entity.ViewHistory) error {
	return db.DB.Create(history).Error
}

// Delete deletes a view history by ID
func (r *ViewHistoryRepository) Delete(id uint) error {
	return db.DB.Delete(&entity.ViewHistory{}, id).Error
}

// DeleteByUserID deletes all view history for a user
func (r *ViewHistoryRepository) DeleteByUserID(userID uint) error {
	return db.DB.Where("user_id = ?", userID).Delete(&entity.ViewHistory{}).Error
}

// GetByID gets a view history by ID
func (r *ViewHistoryRepository) GetByID(id uint) (*entity.ViewHistory, error) {
	var history entity.ViewHistory
	err := db.DB.First(&history, id).Error
	return &history, err
}

// GetByUserIDAndArticleID gets a view history by user ID and article ID
func (r *ViewHistoryRepository) GetByUserIDAndArticleID(userID uint, articleID int) (*entity.ViewHistory, error) {
	var history entity.ViewHistory
	err := db.DB.Where("user_id = ? AND article_id = ?", userID, articleID).First(&history).Error
	return &history, err
}

// ListByUserID lists all view history for a user
func (r *ViewHistoryRepository) ListByUserID(ctx *context.Context, userID uint) ([]entity.ViewHistory, error) {
	var histories []entity.ViewHistory
	query := db.DB.Where("user_id = ?", userID)
	if ctx != nil {
		if ctx.Limit > 0 {
			query = query.Limit(ctx.Limit)
			if ctx.Page > 0 {
				query = query.Offset((ctx.Page - 1) * ctx.Limit)
			}
		}
		if ctx.Order != "" {
			query = query.Order(ctx.Order)
		}
	} else {
		query = query.Order("viewed_at desc")
	}
	err := query.Find(&histories).Error
	return histories, err
}

// CountByUserID counts view history for a user
func (r *ViewHistoryRepository) CountByUserID(userID uint) (int64, error) {
	var count int64
	err := db.DB.Model(&entity.ViewHistory{}).Where("user_id = ?", userID).Count(&count).Error
	return count, err
}

// UpdateViewedAt updates the viewed_at timestamp for a history entry
func (r *ViewHistoryRepository) UpdateViewedAt(id uint, viewedAt time.Time) error {
	return db.DB.Model(&entity.ViewHistory{}).Where("id = ?", id).Update("viewed_at", viewedAt).Error
}