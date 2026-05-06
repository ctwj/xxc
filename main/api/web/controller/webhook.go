package controller

import (
	"github.com/gofiber/fiber/v2"
	"moss/domain/core/service"
	"moss/infrastructure/support/log"
)

// WebhookRevalidate handles ISR revalidation webhook calls from the frontend
func WebhookRevalidate(c *fiber.Ctx) error {
	type RevalidateRequest struct {
		Secret string `json:"secret"`
		Slug   string `json:"slug"`
		Type   string `json:"type"` // article, category, tag, all
	}

	var req RevalidateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
	}

	// Validate secret (should match a configured secret)
	// TODO: Add proper secret validation from config
	if req.Secret == "" {
		return c.Status(401).JSON(fiber.Map{"error": "secret required"})
	}

	// Process based on type
	switch req.Type {
	case "article":
		if req.Slug == "" {
			return c.Status(400).JSON(fiber.Map{"error": "slug required for article revalidation"})
		}
		// Verify article exists
		_, err := service.Article.GetBySlug(req.Slug)
		if err != nil {
			return c.Status(404).JSON(fiber.Map{"error": "article not found"})
		}
		log.Info("Article revalidation requested", log.String("slug", req.Slug))

	case "category":
		if req.Slug == "" {
			return c.Status(400).JSON(fiber.Map{"error": "slug required for category revalidation"})
		}
		_, err := service.Category.GetBySlug(req.Slug)
		if err != nil {
			return c.Status(404).JSON(fiber.Map{"error": "category not found"})
		}
		log.Info("Category revalidation requested", log.String("slug", req.Slug))

	case "tag":
		if req.Slug == "" {
			return c.Status(400).JSON(fiber.Map{"error": "slug required for tag revalidation"})
		}
		_, err := service.Tag.GetBySlug(req.Slug)
		if err != nil {
			return c.Status(404).JSON(fiber.Map{"error": "tag not found"})
		}
		log.Info("Tag revalidation requested", log.String("slug", req.Slug))

	case "all":
		log.Info("Full site revalidation requested")

	default:
		return c.Status(400).JSON(fiber.Map{"error": "invalid type, must be article, category, tag, or all"})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "revalidation triggered",
		"type":    req.Type,
		"slug":    req.Slug,
	})
}