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

	// Like APIs - authentication required
	likes := api.Group("/likes", middleware.JWTMiddleware)
	likes.Post("/", controller.LikeSet)
	likes.Get("/", controller.LikeList)
	likes.Get("/status", controller.LikeGet)

	// View History APIs - authentication required
	viewHistory := api.Group("/history", middleware.JWTMiddleware)
	viewHistory.Post("/", controller.ViewHistoryRecord)
	viewHistory.Get("/", controller.ViewHistoryList)
	viewHistory.Delete("/", controller.ViewHistoryClear)

	// Webhook API - for ISR revalidation
	api.Post("/webhook/revalidate", controller.WebhookRevalidate)
}