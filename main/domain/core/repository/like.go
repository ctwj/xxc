package repository

import (
	"moss/domain/core/entity"
	"moss/domain/core/repository/context"
	"moss/infrastructure/persistent/db"
)

type LikeRepository struct{}

var Like = &LikeRepository{}

// Create creates a new like
func (r *LikeRepository) Create(like *entity.Like) error {
	return db.DB.Create(like).Error
}

// Update updates a like
func (r *LikeRepository) Update(like *entity.Like) error {
	return db.DB.Save(like).Error
}

// Upsert creates or updates a like (based on unique index)
func (r *LikeRepository) Upsert(like *entity.Like) error {
	return db.DB.Save(like).Error
}

// Delete deletes a like by ID
func (r *LikeRepository) Delete(id uint) error {
	return db.DB.Delete(&entity.Like{}, id).Error
}

// DeleteByUserIDAndArticleID deletes a like by user ID and article ID
func (r *LikeRepository) DeleteByUserIDAndArticleID(userID uint, articleID int) error {
	return db.DB.Where("user_id = ? AND article_id = ?", userID, articleID).Delete(&entity.Like{}).Error
}

// GetByID gets a like by ID
func (r *LikeRepository) GetByID(id uint) (*entity.Like, error) {
	var like entity.Like
	err := db.DB.First(&like, id).Error
	return &like, err
}

// GetByUserIDAndArticleID gets a like by user ID and article ID
func (r *LikeRepository) GetByUserIDAndArticleID(userID uint, articleID int) (*entity.Like, error) {
	var like entity.Like
	err := db.DB.Where("user_id = ? AND article_id = ?", userID, articleID).First(&like).Error
	return &like, err
}

// ListByUserID lists all likes for a user
func (r *LikeRepository) ListByUserID(ctx *context.Context, userID uint) ([]entity.Like, error) {
	var likes []entity.Like
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
	}
	err := query.Find(&likes).Error
	return likes, err
}

// CountByArticleIDAndType counts likes/dislikes for an article
func (r *LikeRepository) CountByArticleIDAndType(articleID int, likeType entity.LikeType) (int64, error) {
	var count int64
	err := db.DB.Model(&entity.Like{}).Where("article_id = ? AND type = ?", articleID, likeType).Count(&count).Error
	return count, err
}

// CountByUserID counts likes for a user
func (r *LikeRepository) CountByUserID(userID uint) (int64, error) {
	var count int64
	err := db.DB.Model(&entity.Like{}).Where("user_id = ?", userID).Count(&count).Error
	return count, err
}

// ExistsByUserIDAndArticleID checks if a like exists
func (r *LikeRepository) ExistsByUserIDAndArticleID(userID uint, articleID int) (bool, error) {
	var count int64
	err := db.DB.Model(&entity.Like{}).Where("user_id = ? AND article_id = ?", userID, articleID).Count(&count).Error
	return count > 0, err
}
