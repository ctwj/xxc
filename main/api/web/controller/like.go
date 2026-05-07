package controller

import (
	"github.com/gofiber/fiber/v2"
	"moss/domain/core/entity"
	"moss/domain/core/service"
)

// LikeSet sets a like or dislike for an article
func LikeSet(c *fiber.Ctx) error {
	userID := GetUserID(c)
	if userID == 0 {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}

	type SetRequest struct {
		ArticleID int `json:"articleId"`
		Type      int `json:"type"` // 1=like, 2=dislike, 0=remove
	}

	var req SetRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
	}

	if req.ArticleID <= 0 {
		return c.Status(400).JSON(fiber.Map{"error": "articleId is required"})
	}

	// Check if article exists
	_, err := service.Article.Get(req.ArticleID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "article not found"})
	}

	if req.Type == 0 {
		// Remove like
		err = service.Like.RemoveLike(userID, req.ArticleID)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "failed to remove like"})
		}
		return c.JSON(fiber.Map{
			"success": true,
			"message": "like removed",
		})
	}

	// Validate type
	if req.Type != int(entity.LikeTypeLike) && req.Type != int(entity.LikeTypeDislike) {
		return c.Status(400).JSON(fiber.Map{"error": "invalid type, must be 0, 1, or 2"})
	}

	err = service.Like.SetLike(userID, req.ArticleID, entity.LikeType(req.Type))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to set like"})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "like set",
	})
}

// LikeGet gets the user's like status for an article
func LikeGet(c *fiber.Ctx) error {
	userID := GetUserID(c)
	if userID == 0 {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}

	articleID := c.QueryInt("articleId", 0)
	if articleID <= 0 {
		return c.Status(400).JSON(fiber.Map{"error": "articleId is required"})
	}

	likeType := service.Like.GetUserLikeType(userID, articleID)

	// Get like and dislike counts
	likes, _ := service.Like.CountLikes(articleID)
	dislikes, _ := service.Like.CountDislikes(articleID)

	return c.JSON(fiber.Map{
		"type":     likeType,
		"likes":    likes,
		"dislikes": dislikes,
	})
}

// LikeList lists all articles the user has liked
func LikeList(c *fiber.Ctx) error {
	userID := GetUserID(c)
	if userID == 0 {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}

	page := c.QueryInt("page", 1)
	pageSize := c.QueryInt("pageSize", 20)
	likeType := c.QueryInt("type", 0) // 0=all, 1=likes, 2=dislikes

	likes, err := service.Like.ListByUserID(nil, userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to get likes"})
	}

	// Filter by type if specified
	var filtered []entity.Like
	for _, like := range likes {
		if likeType == 0 || int(like.Type) == likeType {
			filtered = append(filtered, like)
		}
	}

	// Build response with article details
	type LikeWithArticle struct {
		ID        uint                   `json:"id"`
		ArticleID int                    `json:"articleId"`
		Type      int                    `json:"type"`
		CreatedAt string                 `json:"createdAt"`
		Article   map[string]interface{} `json:"article,omitempty"`
	}

	result := make([]LikeWithArticle, 0, len(filtered))
	for _, like := range filtered {
		item := LikeWithArticle{
			ID:        like.ID,
			ArticleID: like.ArticleID,
			Type:      int(like.Type),
			CreatedAt: like.CreatedAt.Format("2006-01-02 15:04:05"),
		}

		// Get article details
		article, err := service.Article.Get(like.ArticleID)
		if err == nil {
			item.Article = map[string]interface{}{
				"id":          article.ID,
				"slug":        article.Slug,
				"title":       article.Title,
				"thumbnail":   article.Thumbnail,
				"description": article.Description,
			}
		}

		result = append(result, item)
	}

	// Pagination
	start := (page - 1) * pageSize
	end := start + pageSize
	if start > len(result) {
		start = len(result)
	}
	if end > len(result) {
		end = len(result)
	}

	return c.JSON(fiber.Map{
		"data":     result[start:end],
		"total":    len(result),
		"page":     page,
		"pageSize": pageSize,
	})
}