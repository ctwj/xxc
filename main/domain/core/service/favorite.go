package service

import (
	"errors"
	"moss/domain/core/entity"
	"moss/domain/core/repository"
	"moss/domain/core/repository/context"
)

type FavoriteService struct{}

var Favorite = &FavoriteService{}

// Create creates a new favorite
func (s *FavoriteService) Create(favorite *entity.Favorite) error {
	// Check if already exists
	exists, err := repository.Favorite.ExistsByUserIDAndArticleID(favorite.UserID, favorite.ArticleID)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("favorite already exists")
	}
	return repository.Favorite.Create(favorite)
}

// Delete deletes a favorite by ID
func (s *FavoriteService) Delete(id uint) error {
	return repository.Favorite.Delete(id)
}

// DeleteByUserIDAndArticleID deletes a favorite by user ID and article ID
func (s *FavoriteService) DeleteByUserIDAndArticleID(userID uint, articleID int) error {
	return repository.Favorite.DeleteByUserIDAndArticleID(userID, articleID)
}

// GetByID gets a favorite by ID
func (s *FavoriteService) GetByID(id uint) (*entity.Favorite, error) {
	return repository.Favorite.GetByID(id)
}

// ListByUserID lists all favorites for a user
func (s *FavoriteService) ListByUserID(ctx *context.Context, userID uint) ([]entity.Favorite, error) {
	return repository.Favorite.ListByUserID(ctx, userID)
}

// CountByUserID counts favorites for a user
func (s *FavoriteService) CountByUserID(userID uint) (int64, error) {
	return repository.Favorite.CountByUserID(userID)
}

// ExistsByUserIDAndArticleID checks if a favorite exists
func (s *FavoriteService) ExistsByUserIDAndArticleID(userID uint, articleID int) (bool, error) {
	return repository.Favorite.ExistsByUserIDAndArticleID(userID, articleID)
}