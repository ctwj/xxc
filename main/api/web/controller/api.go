package controller

import (
	"github.com/gofiber/fiber/v2"
	"moss/domain/config"
	coreCtx "moss/domain/core/repository/context"
	"moss/domain/core/entity"
	"moss/domain/core/service"
	"moss/infrastructure/support/log"
)

// APIArticleList returns a list of articles for the frontend API
func APIArticleList(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	pageSize := c.QueryInt("pageSize", 20)
	if pageSize > 100 {
		pageSize = 100
	}
	categorySlug := c.Query("category")
	tagSlug := c.Query("tag")

	ctx := &coreCtx.Context{
		Limit:   pageSize,
		Order:   "id desc",
		Page:    page,
		Comment: "API.ArticleList",
		Where:   &coreCtx.Where{Field: "status", Operator: coreCtx.WhereOperatorEqualTrue},
	}

	var list []entity.ArticleBase
	var total int64
	var err error

	if categorySlug != "" {
		// Filter by category
		category, catErr := service.Category.GetBySlug(categorySlug)
		if catErr != nil {
			return c.JSON(fiber.Map{
				"data":     []entity.ArticleBase{},
				"total":    0,
				"page":     page,
				"pageSize": pageSize,
				"hasMore":  false,
			})
		}
		list, err = service.Article.ListByCategoryID(ctx, category.ID)
		total, _ = service.Article.CountByCategoryID(category.ID)
	} else if tagSlug != "" {
		// Filter by tag
		tag, tagErr := service.Tag.GetBySlug(tagSlug)
		if tagErr != nil {
			return c.JSON(fiber.Map{
				"data":     []entity.ArticleBase{},
				"total":    0,
				"page":     page,
				"pageSize": pageSize,
				"hasMore":  false,
			})
		}
		list, err = service.Article.ListByTagID(ctx, tag.ID)
		// Count articles by tag - use mapping service
		total, _ = service.Mapping.CountByTagID(tag.ID)
	} else {
		list, err = service.Article.List(ctx)
		total, _ = service.Article.CountByWhere(&coreCtx.Where{Field: "status", Operator: coreCtx.WhereOperatorEqualTrue})
	}

	if err != nil {
		log.Error("API article list failed", log.Err(err))
		return c.Status(500).JSON(fiber.Map{"error": "failed to get articles"})
	}

	hasMore := page*pageSize < int(total)

	return c.JSON(fiber.Map{
		"data":     list,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
		"hasMore":  hasMore,
	})
}

// APIArticleDetail returns a single article by slug
func APIArticleDetail(c *fiber.Ctx) error {
	slug := c.Params("slug")

	article, err := service.Article.GetBySlug(slug)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "article not found"})
	}

	// Check if article is published
	if !article.Status {
		return c.Status(404).JSON(fiber.Map{"error": "article not found"})
	}

	// Get category
	var category interface{}
	if article.CategoryID > 0 {
		cat, _ := service.Category.Get(article.CategoryID)
		category = cat
	}

	// Get tags
	tags, _ := service.Tag.ListByArticleID(&coreCtx.Context{}, article.ID)

	return c.JSON(fiber.Map{
		"id":          article.ID,
		"slug":        article.Slug,
		"title":       article.Title,
		"content":     article.Content,
		"description": article.Description,
		"thumbnail":   article.Thumbnail,
		"keywords":    article.Keywords,
		"views":       article.Views,
		"createTime":  article.CreateTime,
		"category":    category,
		"tags":        tags,
		"extends":     article.Extends,
		"res":         article.Res,
	})
}

// APICategoryList returns all categories
func APICategoryList(c *fiber.Ctx) error {
	categories, err := service.Category.List(&coreCtx.Context{
		Order: "id asc",
	})
	if err != nil {
		log.Error("API category list failed", log.Err(err))
		return c.Status(500).JSON(fiber.Map{"error": "failed to get categories"})
	}

	// Add article count to each category
	result := make([]map[string]interface{}, len(categories))
	for i, cat := range categories {
		count, _ := service.Article.CountByCategoryID(cat.ID)
		result[i] = map[string]interface{}{
			"id":           cat.ID,
			"slug":         cat.Slug,
			"name":         cat.Name,
			"title":        cat.Title,
			"description":  cat.Description,
			"articleCount": count,
		}
	}

	return c.JSON(result)
}

// APITagList returns all tags
func APITagList(c *fiber.Ctx) error {
	tags, err := service.Tag.List(&coreCtx.Context{
		Order: "id asc",
	})
	if err != nil {
		log.Error("API tag list failed", log.Err(err))
		return c.Status(500).JSON(fiber.Map{"error": "failed to get tags"})
	}

	// Add article count to each tag
	result := make([]map[string]interface{}, len(tags))
	for i, tag := range tags {
		count, _ := service.Mapping.CountByTagID(tag.ID)
		result[i] = map[string]interface{}{
			"id":           tag.ID,
			"slug":         tag.Slug,
			"name":         tag.Name,
			"title":        tag.Title,
			"description":  tag.Description,
			"articleCount": count,
		}
	}

	return c.JSON(result)
}

// APISearch searches articles by keyword
func APISearch(c *fiber.Ctx) error {
	keyword := c.Query("keyword")
	if keyword == "" {
		return c.JSON(fiber.Map{
			"data":    []entity.ArticleBase{},
			"keyword": "",
			"total":   0,
		})
	}

	page := c.QueryInt("page", 1)
	limit := config.Config.Template.IndexList.Limit
	if limit <= 0 {
		limit = 20
	}

	ctx := &coreCtx.Context{
		Limit:   limit,
		Order:   "id desc",
		Page:    page,
		Comment: "API.Search",
		Where:   &coreCtx.Where{Field: "status", Operator: coreCtx.WhereOperatorEqualTrue},
	}

	list, err := service.Article.ListByKeyword(ctx, keyword)
	if err != nil {
		log.Error("API search failed", log.Err(err))
		return c.Status(500).JSON(fiber.Map{"error": "search failed"})
	}

	// Get total count
	total, _ := service.Article.CountByKeyword(keyword)

	return c.JSON(fiber.Map{
		"data":    list,
		"keyword": keyword,
		"total":   total,
	})
}