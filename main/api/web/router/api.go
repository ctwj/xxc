package router

import (
	"github.com/gofiber/fiber/v2"
	"moss/api/web/controller"
	"moss/api/web/middleware"
)

// RegisterAPIRoutes registers public API routes for the frontend
func RegisterAPIRoutes(app fiber.Router) {
	api := app.Group("/api")

	// Public APIs - no authentication required
	api.Get("/articles", controller.APIArticleList)
	api.Get("/articles/:slug", controller.APIArticleDetail)
	api.Get("/categories", controller.APICategoryList)
	api.Get("/tags", controller.APITagList)
	api.Get("/search", controller.APISearch)

	// Auth APIs
	auth := api.Group("/auth")
	auth.Post("/login", controller.AuthLogin)
	auth.Post("/register", controller.AuthRegister)
	auth.Post("/logout", controller.AuthLogout)
	auth.Get("/me", middleware.JWTMiddleware, controller.AuthMe)

	// Favorites APIs - authentication required
	favorites := api.Group("/favorites", middleware.JWTMiddleware)
	favorites.Get("/", controller.FavoriteList)
	favorites.Post("/", controller.FavoriteAdd)
	favorites.Delete("/:id", controller.FavoriteRemove)

	// Webhook API - for ISR revalidation
	api.Post("/webhook/revalidate", controller.WebhookRevalidate)
}