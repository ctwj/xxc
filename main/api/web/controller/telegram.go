package controller

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"moss/api/web/mapper"
	"moss/domain/support/service"
)

// TelegramSendCodeRequest 发送验证码请求
type TelegramSendCodeRequest struct {
	PhoneNumber string `json:"phone_number"`
}

// TelegramVerifyCodeRequest 验证码验证请求
type TelegramVerifyCodeRequest struct {
	Code string `json:"code"`
}

// TelegramAuthInterface 定义 Telegram 认证接口
type TelegramAuthInterface interface {
	SendAuthCode(phoneNumber string) (string, error)
	VerifyAuthCode(code string) error
	GetAuthStatus() map[string]interface{}
	GetUserChannels() ([]map[string]interface{}, error)
	ClearAuth() error
	GetStatus() map[string]interface{}
	CheckSessionStatus() map[string]interface{}
}

// TelegramLogInterface 定义获取日志接口
type TelegramLogInterface interface {
	GetRecentLogs(limit int) (interface{}, error)
}

// TelegramMediaInterface 定义媒体接口
type TelegramMediaInterface interface {
	GetMediaURL(mediaId int64) (string, error)
}

// TelegramMediaDownloadInterface 定义媒体下载接口
type TelegramMediaDownloadInterface interface {
	DownloadMediaFile(mediaId int64) ([]byte, string, error)
}

// TelegramSendCode 发送 Telegram 验证码
func TelegramSendCode(ctx *fiber.Ctx) error {
	var req TelegramSendCodeRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.JSON(mapper.MessageResult(err))
	}

	// 获取插件实例
	plugin, err := service.Plugin.Get("TelegramChannelSync")
	if err != nil {
		return ctx.JSON(mapper.MessageResultData(nil, err))
	}

	// 类型断言获取 Telegram 认证接口
	telegramPlugin, ok := plugin.Entry.(TelegramAuthInterface)
	if !ok {
		return ctx.JSON(mapper.MessageResultData(nil, errors.New("plugin does not support telegram auth")))
	}

	// 调用插件的发送验证码方法
	result, err := telegramPlugin.SendAuthCode(req.PhoneNumber)
	if err != nil {
		return ctx.JSON(mapper.MessageResultData(nil, err))
	}

	return ctx.JSON(mapper.MessageResultData(map[string]string{"phone_code_hash": result}, nil))
}

// TelegramVerifyCode 验证 Telegram 验证码
func TelegramVerifyCode(ctx *fiber.Ctx) error {
	var req TelegramVerifyCodeRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.JSON(mapper.MessageResult(err))
	}

	// 获取插件实例
	plugin, err := service.Plugin.Get("TelegramChannelSync")
	if err != nil {
		return ctx.JSON(mapper.MessageResultData(nil, err))
	}

	// 类型断言获取 Telegram 认证接口
	telegramPlugin, ok := plugin.Entry.(TelegramAuthInterface)
	if !ok {
		return ctx.JSON(mapper.MessageResultData(nil, errors.New("plugin does not support telegram auth")))
	}

	// 调用插件的验证方法
	err = telegramPlugin.VerifyAuthCode(req.Code)
	if err != nil {
		return ctx.JSON(mapper.MessageResultData(nil, err))
	}

	return ctx.JSON(mapper.MessageResultData(map[string]bool{"success": true}, nil))
}

// TelegramAuthStatus 获取 Telegram 认证状态
func TelegramAuthStatus(ctx *fiber.Ctx) error {
	// 获取插件实例
	plugin, err := service.Plugin.Get("TelegramChannelSync")
	if err != nil {
		return ctx.JSON(mapper.MessageResultData(nil, err))
	}

	// 类型断言获取 Telegram 认证接口
	telegramPlugin, ok := plugin.Entry.(TelegramAuthInterface)
	if !ok {
		return ctx.JSON(mapper.MessageResultData(nil, errors.New("plugin does not support telegram auth")))
	}

	// 获取认证状态
	status := telegramPlugin.GetAuthStatus()
	return ctx.JSON(mapper.MessageResultData(status, nil))
}

// TelegramGetChannels 获取用户的 Telegram 频道列表
func TelegramGetChannels(ctx *fiber.Ctx) error {
	// 获取插件实例
	plugin, err := service.Plugin.Get("TelegramChannelSync")
	if err != nil {
		return ctx.JSON(mapper.MessageResultData(nil, err))
	}

	// 类型断言获取 Telegram 认证接口
	telegramPlugin, ok := plugin.Entry.(TelegramAuthInterface)
	if !ok {
		return ctx.JSON(mapper.MessageResultData(nil, errors.New("plugin does not support telegram auth")))
	}

	// 获取频道列表
	channels, err := telegramPlugin.GetUserChannels()
	if err != nil {
		return ctx.JSON(mapper.MessageResultData(nil, err))
	}

	return ctx.JSON(mapper.MessageResultData(channels, nil))
}

// TelegramClearAuth 清除 Telegram 认证
func TelegramClearAuth(ctx *fiber.Ctx) error {
	// 获取插件实例
	plugin, err := service.Plugin.Get("TelegramChannelSync")
	if err != nil {
		return ctx.JSON(mapper.MessageResultData(nil, err))
	}

	// 类型断言获取 Telegram 认证接口
	telegramPlugin, ok := plugin.Entry.(TelegramAuthInterface)
	if !ok {
		return ctx.JSON(mapper.MessageResultData(nil, errors.New("plugin does not support telegram auth")))
	}

	// 清除认证
	err = telegramPlugin.ClearAuth()
	if err != nil {
		return ctx.JSON(mapper.MessageResultData(nil, err))
	}

	return ctx.JSON(mapper.MessageResultData(map[string]bool{"success": true}, nil))
}

// TelegramDebugStatus 获取调试状态信息
func TelegramDebugStatus(ctx *fiber.Ctx) error {
	// 获取插件实例
	plugin, err := service.Plugin.Get("TelegramChannelSync")
	if err != nil {
		return ctx.JSON(mapper.MessageResultData(nil, err))
	}

	// 类型断言获取 Telegram 认证接口
	telegramPlugin, ok := plugin.Entry.(TelegramAuthInterface)
	if !ok {
		return ctx.JSON(mapper.MessageResultData(nil, errors.New("plugin does not support telegram auth")))
	}

	// 获取状态
	status := telegramPlugin.GetStatus()
	return ctx.JSON(mapper.MessageResultData(status, nil))
}

// TelegramCheckSession 检查会话状态
func TelegramCheckSession(ctx *fiber.Ctx) error {
	// 获取插件实例
	plugin, err := service.Plugin.Get("TelegramChannelSync")
	if err != nil {
		return ctx.JSON(mapper.MessageResultData(nil, err))
	}

	// 类型断言获取 Telegram 认证接口
	telegramPlugin, ok := plugin.Entry.(TelegramAuthInterface)
	if !ok {
		return ctx.JSON(mapper.MessageResultData(nil, errors.New("plugin does not support telegram auth")))
	}

	// 检查会话状态
	status := telegramPlugin.CheckSessionStatus()
	return ctx.JSON(mapper.MessageResultData(status, nil))
}

// TelegramGetLogs 获取同步日志
func TelegramGetLogs(ctx *fiber.Ctx) error {
	// 获取插件实例
	plugin, err := service.Plugin.Get("TelegramChannelSync")
	if err != nil {
		return ctx.JSON(mapper.MessageResultData(nil, err))
	}

	// 类型断言获取日志接口
	logPlugin, ok := plugin.Entry.(TelegramLogInterface)
	if !ok {
		return ctx.JSON(mapper.MessageResultData(nil, errors.New("plugin does not support log interface")))
	}

	// 获取日志
	limit := 50
	logs, err := logPlugin.GetRecentLogs(limit)
	if err != nil {
		return ctx.JSON(mapper.MessageResultData(nil, err))
	}

	return ctx.JSON(mapper.MessageResultData(logs, nil))
}

// TelegramPublicDebug 公开的调试 API（仅用于开发调试）
func TelegramPublicDebug(ctx *fiber.Ctx) error {
	// 获取插件实例
	plugin, err := service.Plugin.Get("TelegramChannelSync")
	if err != nil {
		return ctx.JSON(mapper.MessageResultData(nil, err))
	}

	// 类型断言获取 Telegram 认证接口
	telegramPlugin, ok := plugin.Entry.(TelegramAuthInterface)
	if !ok {
		return ctx.JSON(mapper.MessageResultData(nil, errors.New("plugin does not support telegram auth")))
	}

	// 获取状态
	status := telegramPlugin.GetStatus()
	return ctx.JSON(mapper.MessageResultData(status, nil))
}

// TelegramPublicChannels 公开的获取频道列表 API（仅用于开发调试）
func TelegramPublicChannels(ctx *fiber.Ctx) error {
	// 获取插件实例
	plugin, err := service.Plugin.Get("TelegramChannelSync")
	if err != nil {
		return ctx.JSON(mapper.MessageResultData(nil, err))
	}

	// 类型断言获取 Telegram 认证接口
	telegramPlugin, ok := plugin.Entry.(TelegramAuthInterface)
	if !ok {
		return ctx.JSON(mapper.MessageResultData(nil, errors.New("plugin does not support telegram auth")))
	}

	// 获取频道列表
	channels, err := telegramPlugin.GetUserChannels()
	if err != nil {
		return ctx.JSON(mapper.MessageResultData(nil, err))
	}

	return ctx.JSON(mapper.MessageResultData(channels, nil))
}

// TelegramGetMedia 获取媒体文件（返回下载代理）
func TelegramGetMedia(ctx *fiber.Ctx) error {
	mediaIdStr := ctx.Params("mediaId")
	mediaId, err := strconv.ParseInt(mediaIdStr, 10, 64)
	if err != nil {
		return ctx.JSON(mapper.MessageResultData(nil, errors.New("invalid media id")))
	}

	// 获取插件实例
	plugin, err := service.Plugin.Get("TelegramChannelSync")
	if err != nil {
		return ctx.JSON(mapper.MessageResultData(nil, err))
	}

	// 尝试下载接口
	downloadPlugin, ok := plugin.Entry.(TelegramMediaDownloadInterface)
	if ok {
		data, mimeType, err := downloadPlugin.DownloadMediaFile(mediaId)
		if err != nil {
			return ctx.Status(404).JSON(mapper.MessageResultData(nil, err))
		}
		// 设置 Content-Type 并返回文件内容
		ctx.Set("Content-Type", mimeType)
		ctx.Set("Cache-Control", "public, max-age=31536000") // 缓存一年
		return ctx.Send(data)
	}

	// 降级到 URL 接口
	mediaPlugin, ok := plugin.Entry.(TelegramMediaInterface)
	if !ok {
		return ctx.JSON(mapper.MessageResultData(nil, errors.New("plugin does not support media interface")))
	}

	// 获取媒体 URL
	url, err := mediaPlugin.GetMediaURL(mediaId)
	if err != nil {
		return ctx.JSON(mapper.MessageResultData(nil, err))
	}

	return ctx.JSON(mapper.MessageResultData(map[string]string{"url": url}, nil))
}
