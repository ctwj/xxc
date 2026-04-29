package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"moss/domain/config"
)

// CORSConfig returns CORS middleware configuration
func CORSConfig() fiber.Handler {
	// Get allowed origins from config
	origins := config.Config.Router.CORSOrigins

	// If no origins configured, use permissive CORS without credentials
	// This is safe for development but should be configured for production
	if origins == "" {
		return cors.New(cors.Config{
			AllowOrigins: "*",
			AllowHeaders: "Origin, Content-Type, Accept, Authorization, X-Requested-With",
			AllowMethods: "GET, POST, PUT, DELETE, OPTIONS, PATCH",
			// AllowCredentials must be false when AllowOrigins is "*"
			AllowCredentials: false,
			ExposeHeaders:    "Set-Cookie",
			MaxAge:           86400, // 24 hours
		})
	}

	// When specific origins are configured, enable credentials
	return cors.New(cors.Config{
		AllowOrigins:     origins,
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization, X-Requested-With",
		AllowMethods:     "GET, POST, PUT, DELETE, OPTIONS, PATCH",
		AllowCredentials: true,
		ExposeHeaders:    "Set-Cookie",
		MaxAge:           86400, // 24 hours
	})
}
