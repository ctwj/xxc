package service

import (
	"errors"
	"time"

	"moss/domain/core/entity"
	"moss/domain/core/repository"
	"moss/domain/core/repository/context"
)

var UserFavorite = new(UserFavoriteService)

type UserFavoriteService struct{}

// Add 添加收藏
func (s *UserFavoriteService) Add(userID, articleID uint) (*entity.UserFavorite, error) {
	if userID == 0 {
		return nil, errors.New("user id is required")
	}
	if articleID == 0 {
		return nil, errors.New("article id is required")
	}

	// 检查是否已收藏
	exists, err := repository.UserFavorite.ExistsByUserIDAndArticleID(userID, articleID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("already favorited")
	}

	favorite := &entity.UserFavorite{
		UserID:     userID,
		ArticleID:  articleID,
		CreateTime: time.Now().Unix(),
	}

	if err := repository.UserFavorite.Create(favorite); err != nil {
		return nil, err
	}

	return favorite, nil
}

// Remove 取消收藏
func (s *UserFavoriteService) Remove(userID, articleID uint) error {
	return repository.UserFavorite.DeleteByUserIDAndArticleID(userID, articleID)
}

// RemoveByID 根据收藏ID取消收藏
func (s *UserFavoriteService) RemoveByID(id uint, userID uint) error {
	favorite, err := repository.UserFavorite.Get(id)
	if err != nil {
		return err
	}
	if favorite.ID == 0 {
		return errors.New("favorite not found")
	}
	// 验证是否是当前用户的收藏
	if favorite.UserID != userID {
		return errors.New("not authorized")
	}
	return repository.UserFavorite.Delete(id)
}

// Toggle 切换收藏状态（已收藏则取消，未收藏则添加）
func (s *UserFavoriteService) Toggle(userID, articleID uint) (bool, error) {
	exists, err := repository.UserFavorite.ExistsByUserIDAndArticleID(userID, articleID)
	if err != nil {
		return false, err
	}

	if exists {
		if err := s.Remove(userID, articleID); err != nil {
			return false, err
		}
		return false, nil // 返回 false 表示已取消收藏
	}

	if _, err := s.Add(userID, articleID); err != nil {
		return false, err
	}
	return true, nil // 返回 true 表示已添加收藏
}

// IsFavorited 检查是否已收藏
func (s *UserFavoriteService) IsFavorited(userID, articleID uint) (bool, error) {
	return repository.UserFavorite.ExistsByUserIDAndArticleID(userID, articleID)
}

// ListByUserID 获取用户收藏列表
func (s *UserFavoriteService) ListByUserID(ctx *context.Context, userID uint) ([]entity.UserFavorite, error) {
	return repository.UserFavorite.ListByUserID(ctx, userID)
}

// CountByUserID 统计用户收藏数量
func (s *UserFavoriteService) CountByUserID(userID uint) (int64, error) {
	return repository.UserFavorite.CountByUserID(userID)
}

// CountByArticleID 统计文章收藏数量
func (s *UserFavoriteService) CountByArticleID(articleID uint) (int64, error) {
	return repository.UserFavorite.CountByArticleID(articleID)
}

// GetArticleIDsByUserID 获取用户收藏的文章ID列表
func (s *UserFavoriteService) GetArticleIDsByUserID(userID uint) ([]uint, error) {
	return repository.UserFavorite.ListArticleIDsByUserID(userID)
}
