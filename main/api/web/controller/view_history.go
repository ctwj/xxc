package controller

import (
	"github.com/gofiber/fiber/v2"
	"moss/domain/core/service"
)

// ViewHistoryRecord records a view for an article
func ViewHistoryRecord(c *fiber.Ctx) error {
	userID := GetUserID(c)
	if userID == 0 {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}

	type RecordRequest struct {
		ArticleID int `json:"articleId"`
	}

	var req RecordRequest
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

	err = service.ViewHistory.Record(userID, req.ArticleID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to record view"})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "view recorded",
	})
}

// ViewHistoryList lists the user's view history
func ViewHistoryList(c *fiber.Ctx) error {
	userID := GetUserID(c)
	if userID == 0 {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}

	page := c.QueryInt("page", 1)
	pageSize := c.QueryInt("pageSize", 20)

	histories, err := service.ViewHistory.ListByUserID(nil, userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to get view history"})
	}

	// Build response with article details
	type HistoryWithArticle struct {
		ID        uint                   `json:"id"`
		ArticleID int                    `json:"articleId"`
		ViewedAt  string                 `json:"viewedAt"`
		Article   map[string]interface{} `json:"article,omitempty"`
	}

	result := make([]HistoryWithArticle, 0, len(histories))
	for _, h := range histories {
		item := HistoryWithArticle{
			ID:        h.ID,
			ArticleID: h.ArticleID,
			ViewedAt:  h.ViewedAt.Format("2006-01-02 15:04:05"),
		}

		// Get article details
		article, err := service.Article.Get(h.ArticleID)
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

// ViewHistoryClear clears all view history for the user
func ViewHistoryClear(c *fiber.Ctx) error {
	userID := GetUserID(c)
	if userID == 0 {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}

	err := service.ViewHistory.DeleteByUserID(userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to clear view history"})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "view history cleared",
	})
}