package service

import (
	"moss/domain/core/entity"
	"moss/domain/core/repository"
	"moss/domain/core/repository/context"
)

type LikeService struct{}

var Like = &LikeService{}

// SetLike sets a like or dislike for an article
func (s *LikeService) SetLike(userID uint, articleID int, likeType entity.LikeType) error {
	// Check if already exists
	existing, err := repository.Like.GetByUserIDAndArticleID(userID, articleID)
	if err != nil {
		// Not found, create new
		like := &entity.Like{
			UserID:    userID,
			ArticleID: articleID,
			Type:      likeType,
		}
		return repository.Like.Create(like)
	}

	// Update existing
	existing.Type = likeType
	return repository.Like.Update(existing)
}

// RemoveLike removes a like/dislike for an article
func (s *LikeService) RemoveLike(userID uint, articleID int) error {
	return repository.Like.DeleteByUserIDAndArticleID(userID, articleID)
}

// GetByUserIDAndArticleID gets a like by user ID and article ID
func (s *LikeService) GetByUserIDAndArticleID(userID uint, articleID int) (*entity.Like, error) {
	return repository.Like.GetByUserIDAndArticleID(userID, articleID)
}

// GetUserLikeType gets the user's like type for an article (0 if not liked)
func (s *LikeService) GetUserLikeType(userID uint, articleID int) entity.LikeType {
	like, err := repository.Like.GetByUserIDAndArticleID(userID, articleID)
	if err != nil {
		return entity.LikeTypeNone
	}
	return like.Type
}

// ListByUserID lists all likes for a user
func (s *LikeService) ListByUserID(ctx *context.Context, userID uint) ([]entity.Like, error) {
	return repository.Like.ListByUserID(ctx, userID)
}

// CountByArticleID counts likes for an article
func (s *LikeService) CountLikes(articleID int) (int64, error) {
	return repository.Like.CountByArticleIDAndType(articleID, entity.LikeTypeLike)
}

// CountDislikes counts dislikes for an article
func (s *LikeService) CountDislikes(articleID int) (int64, error) {
	return repository.Like.CountByArticleIDAndType(articleID, entity.LikeTypeDislike)
}

// CountByUserID counts likes for a user
func (s *LikeService) CountByUserID(userID uint) (int64, error) {
	return repository.Like.CountByUserID(userID)
}

// ExistsByUserIDAndArticleID checks if a like exists
func (s *LikeService) ExistsByUserIDAndArticleID(userID uint, articleID int) (bool, error) {
	return repository.Like.ExistsByUserIDAndArticleID(userID, articleID)
}