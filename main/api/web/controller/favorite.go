package controller

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"moss/domain/core/entity"
	coreCtx "moss/domain/core/repository/context"
	"moss/domain/core/service"
)

// FavoriteList returns the user's favorite articles
func FavoriteList(c *fiber.Ctx) error {
	userID := GetUserID(c)
	if userID == 0 {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}

	page := c.QueryInt("page", 1)
	pageSize := c.QueryInt("pageSize", 20)

	ctx := &coreCtx.Context{
		Limit: pageSize,
		Order: "id desc",
		Page:  page,
	}

	favorites, err := service.Favorite.ListByUserID(ctx, userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to get favorites"})
	}

	// Get article details for each favorite
	type FavoriteWithArticle struct {
		ID        uint                   `json:"id"`
		ArticleID int                    `json:"articleId"`
		CreatedAt string                 `json:"createdAt"`
		Article   map[string]interface{} `json:"article,omitempty"`
	}

	result := make([]FavoriteWithArticle, 0, len(favorites))
	for _, fav := range favorites {
		item := FavoriteWithArticle{
			ID:        fav.ID,
			ArticleID: fav.ArticleID,
			CreatedAt: fav.CreatedAt.Format("2006-01-02 15:04:05"),
		}

		// Get article details
		article, err := service.Article.Get(fav.ArticleID)
		if err == nil {
			item.Article = map[string]interface{}{
				"id":          article.ID,
				"slug":        article.Slug,
				"title":       article.Title,
				"thumbnail":   article.Thumbnail,
				"description": article.Description,
				"views":       article.Views,
				"createTime":  article.CreateTime,
			}
		}

		result = append(result, item)
	}

	total, _ := service.Favorite.CountByUserID(userID)

	return c.JSON(fiber.Map{
		"data":     result,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

// FavoriteAdd adds an article to favorites
func FavoriteAdd(c *fiber.Ctx) error {
	userID := GetUserID(c)
	if userID == 0 {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}

	type AddRequest struct {
		ArticleID int `json:"articleId"`
	}

	var req AddRequest
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

	// Check if already favorited
	exists, _ := service.Favorite.ExistsByUserIDAndArticleID(userID, req.ArticleID)
	if exists {
		return c.Status(400).JSON(fiber.Map{"error": "already favorited"})
	}

	// Create favorite
	fav := &entity.Favorite{
		UserID:    userID,
		ArticleID: req.ArticleID,
	}

	err = service.Favorite.Create(fav)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to add favorite"})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "added to favorites",
	})
}

// FavoriteRemove removes an article from favorites
func FavoriteRemove(c *fiber.Ctx) error {
	userID := GetUserID(c)
	if userID == 0 {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}

	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid id"})
	}

	// Verify the favorite belongs to the user
	favorite, err := service.Favorite.GetByID(uint(id))
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "favorite not found"})
	}

	if favorite.UserID != userID {
		return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
	}

	err = service.Favorite.Delete(uint(id))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to remove favorite"})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "removed from favorites",
	})
}