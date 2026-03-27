package router

import (
	"moss/api/web/controller"
	"moss/api/web/middleware"

	"github.com/gofiber/fiber/v2"
)

// RegisterUserRoutes 注册前台用户路由
// 在 router.go 的 newFiber 函数中调用
func (r *Router) RegisterUserRoutes(route fiber.Router) {
	// 公开路由（无需认证）
	route.Post("/user/register", controller.UserRegister)
	route.Post("/user/login", controller.UserLogin)
	route.Get("/user/invite-code/validate", controller.InviteCodeValidate)
	route.Get("/user/verify-email", controller.UserVerifyEmail) // 邮箱验证

	// 需要认证的路由
	userGroup := route.Group("/user", middleware.UserAuth)
	userGroup.Get("/profile", controller.UserProfile)
	userGroup.Post("/profile", controller.UserUpdateProfile)
	userGroup.Post("/password", controller.UserUpdatePassword)
	userGroup.Post("/resend-verification", controller.UserResendVerificationEmail) // 重发验证邮件

	// 用户收藏路由
	userGroup.Post("/favorites", controller.UserFavoriteList)
	userGroup.Get("/favorites/count", controller.UserFavoriteCount)
	userGroup.Post("/favorite/add", controller.UserFavoriteAdd)
	userGroup.Post("/favorite/toggle", controller.UserFavoriteToggle)
	userGroup.Post("/favorite/remove/:article_id", controller.UserFavoriteRemove)
	userGroup.Get("/favorite/check/:article_id", controller.UserFavoriteIsFavorited)
}