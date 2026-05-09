package repository

import (
	"moss/domain/core/entity"
	"moss/domain/core/repository/context"
	"moss/infrastructure/persistent/db"
)

type FavoriteRepository struct{}

var Favorite = &FavoriteRepository{}

func (r *FavoriteRepository) MigrateTable() error {
	return db.DB.AutoMigrate(&entity.Favorite{})
}

// Create creates a new favorite
func (r *FavoriteRepository) Create(favorite *entity.Favorite) error {
	return db.DB.Create(favorite).Error
}

// Delete deletes a favorite by ID
func (r *FavoriteRepository) Delete(id uint) error {
	return db.DB.Delete(&entity.Favorite{}, id).Error
}

// DeleteByUserIDAndArticleID deletes a favorite by user ID and article ID
func (r *FavoriteRepository) DeleteByUserIDAndArticleID(userID uint, articleID int) error {
	return db.DB.Where("user_id = ? AND article_id = ?", userID, articleID).Delete(&entity.Favorite{}).Error
}

// GetByID gets a favorite by ID
func (r *FavoriteRepository) GetByID(id uint) (*entity.Favorite, error) {
	var favorite entity.Favorite
	err := db.DB.First(&favorite, id).Error
	return &favorite, err
}

// GetByUserIDAndArticleID gets a favorite by user ID and article ID
func (r *FavoriteRepository) GetByUserIDAndArticleID(userID uint, articleID int) (*entity.Favorite, error) {
	var favorite entity.Favorite
	err := db.DB.Where("user_id = ? AND article_id = ?", userID, articleID).First(&favorite).Error
	return &favorite, err
}

// ListByUserID lists all favorites for a user
func (r *FavoriteRepository) ListByUserID(ctx *context.Context, userID uint) ([]entity.Favorite, error) {
	var favorites []entity.Favorite
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
	err := query.Find(&favorites).Error
	return favorites, err
}

// CountByUserID counts favorites for a user
func (r *FavoriteRepository) CountByUserID(userID uint) (int64, error) {
	var count int64
	err := db.DB.Model(&entity.Favorite{}).Where("user_id = ?", userID).Count(&count).Error
	return count, err
}

// ExistsByUserIDAndArticleID checks if a favorite exists
func (r *FavoriteRepository) ExistsByUserIDAndArticleID(userID uint, articleID int) (bool, error) {
	var count int64
	err := db.DB.Model(&entity.Favorite{}).Where("user_id = ? AND article_id = ?", userID, articleID).Count(&count).Error
	return count > 0, err
}
