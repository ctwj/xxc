package repository

import (
	"moss/domain/core/entity"
	"moss/domain/core/repository/context"
	"moss/domain/core/repository/gormx"
	"moss/infrastructure/persistent/db"
)

var UserFavorite = new(UserFavoriteRepo)

type UserFavoriteRepo struct{}

// MigrateTable 迁移表结构
func (r *UserFavoriteRepo) MigrateTable() error {
	return db.DB.AutoMigrate(&entity.UserFavorite{})
}

// Create 创建收藏
func (r *UserFavoriteRepo) Create(favorite *entity.UserFavorite) error {
	return db.DB.Create(favorite).Error
}

// Delete 删除收藏
func (r *UserFavoriteRepo) Delete(id uint) error {
	return db.DB.Delete(&entity.UserFavorite{}, id).Error
}

// DeleteByUserIDAndArticleID 根据用户ID和文章ID删除收藏
func (r *UserFavoriteRepo) DeleteByUserIDAndArticleID(userID, articleID uint) error {
	return db.DB.Where("user_id = ? AND article_id = ?", userID, articleID).Delete(&entity.UserFavorite{}).Error
}

// Get 根据ID获取收藏
func (r *UserFavoriteRepo) Get(id uint) (*entity.UserFavorite, error) {
	var favorite entity.UserFavorite
	err := db.DB.First(&favorite, id).Error
	return &favorite, err
}

// GetByUserIDAndArticleID 根据用户ID和文章ID获取收藏
func (r *UserFavoriteRepo) GetByUserIDAndArticleID(userID, articleID uint) (*entity.UserFavorite, error) {
	var favorite entity.UserFavorite
	err := db.DB.Where("user_id = ? AND article_id = ?", userID, articleID).First(&favorite).Error
	return &favorite, err
}

// ExistsByUserIDAndArticleID 检查用户是否已收藏某文章
func (r *UserFavoriteRepo) ExistsByUserIDAndArticleID(userID, articleID uint) (bool, error) {
	var count int64
	err := db.DB.Model(&entity.UserFavorite{}).Where("user_id = ? AND article_id = ?", userID, articleID).Count(&count).Error
	return count > 0, err
}

// CountByUserID 统计用户收藏数量
func (r *UserFavoriteRepo) CountByUserID(userID uint) (int64, error) {
	var count int64
	err := db.DB.Model(&entity.UserFavorite{}).Where("user_id = ?", userID).Count(&count).Error
	return count, err
}

// CountByArticleID 统计文章收藏数量
func (r *UserFavoriteRepo) CountByArticleID(articleID uint) (int64, error) {
	var count int64
	err := db.DB.Model(&entity.UserFavorite{}).Where("article_id = ?", articleID).Count(&count).Error
	return count, err
}

// ListByUserID 获取用户收藏列表
func (r *UserFavoriteRepo) ListByUserID(ctx *context.Context, userID uint) ([]entity.UserFavorite, error) {
	var favorites []entity.UserFavorite
	err := db.DB.Model(&entity.UserFavorite{}).
		Where("user_id = ?", userID).
		Scopes(gormx.Context(ctx)).
		Find(&favorites).Error
	return favorites, err
}

// ListArticleIDsByUserID 获取用户收藏的文章ID列表
func (r *UserFavoriteRepo) ListArticleIDsByUserID(userID uint) ([]uint, error) {
	var articleIDs []uint
	err := db.DB.Model(&entity.UserFavorite{}).
		Where("user_id = ?", userID).
		Pluck("article_id", &articleIDs).Error
	return articleIDs, err
}
