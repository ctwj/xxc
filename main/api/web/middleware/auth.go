package middleware

import (
	"errors"
	"moss/api/web/dto"
	"moss/domain/core/entity"

	"github.com/gofiber/fiber/v2"
)

func Auth(attrName string, predicate func(token string) (roleName string, ok bool)) func(*fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {

		if attrName == "" {
			return errors.New("attrName undefined")
		}

		token := ctx.Get(attrName) // header

		if token == "" {
			token = ctx.Get("Sec-WebSocket-Protocol") // 兼容 websocket
		}

		if token == "" {
			token = ctx.Query(attrName)
		}

		if token == "" {
			return ctx.Status(401).JSON(&dto.MessageResult{Message: "authorization failed"})
		}

		roleName, ok := predicate(token)
		if !ok {
			return ctx.Status(401).JSON(&dto.MessageResult{Message: "authorization failed"})
		}

		ctx.Locals("roleName", roleName)

		return ctx.Next()
	}
}

// UserAuth 用户认证中间件（用于前台用户API）
func UserAuth(ctx *fiber.Ctx) error {
	token := ctx.Get("token")
	if token == "" {
		token = ctx.Query("token")
	}

	if token == "" {
		return ctx.Status(401).JSON(&dto.MessageResult{Message: "token is required"})
	}

	userID, username, ok := entity.VerifyJwtToken(token)
	if !ok {
		return ctx.Status(401).JSON(&dto.MessageResult{Message: "invalid token"})
	}

	ctx.Locals("userID", userID)
	ctx.Locals("username", username)
	return ctx.Next()
}
