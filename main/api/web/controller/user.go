package controller

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"moss/api/web/dto"
	"moss/api/web/mapper"
	"moss/domain/core/entity"
	"moss/domain/core/repository/context"
	"moss/domain/core/service"
)

// ============ 前台用户 API ============

type UserRegisterRequest struct {
	Username   string `json:"username"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	Nickname   string `json:"nickname"`
	InviteCode string `json:"invite_code"`
}

type UserLoginRequest struct {
	Account  string `json:"account"`
	Password string `json:"password"`
}

type UserUpdateProfileRequest struct {
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	Bio      string `json:"bio"`
}

type UserUpdatePasswordRequest struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

type UserLoginResponse struct {
	Token string      `json:"token"`
	User  *entity.User `json:"user"`
}

// UserRegister 用户注册
func UserRegister(ctx *fiber.Ctx) error {
	var req UserRegisterRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.JSON(mapper.MessageFail("invalid request body"))
	}

	user, err := service.User.Register(req.Username, req.Email, req.Password, req.Nickname, req.InviteCode)
	if err != nil {
		return ctx.JSON(mapper.MessageResult(err))
	}

	return ctx.JSON(mapper.MessageResultData(user, nil))
}

// UserLogin 用户登录
func UserLogin(ctx *fiber.Ctx) error {
	time.Sleep(800 * time.Millisecond) // 防止暴力破解

	var req UserLoginRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.JSON(mapper.MessageFail("invalid request body"))
	}

	user, token, err := service.User.Login(req.Account, req.Password)
	if err != nil {
		return ctx.JSON(mapper.MessageResult(err))
	}

	return ctx.JSON(mapper.MessageResultData(&UserLoginResponse{
		Token: token,
		User:  user,
	}, nil))
}

// UserProfile 获取当前用户信息
func UserProfile(ctx *fiber.Ctx) error {
	userID := ctx.Locals("userID").(uint)
	user, err := service.User.GetByID(userID)
	if err != nil {
		return ctx.JSON(mapper.MessageResult(err))
	}
	return ctx.JSON(mapper.MessageResultData(user, nil))
}

// UserUpdateProfile 更新个人资料
func UserUpdateProfile(ctx *fiber.Ctx) error {
	userID := ctx.Locals("userID").(uint)

	var req UserUpdateProfileRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.JSON(mapper.MessageFail("invalid request body"))
	}

	err := service.User.UpdateProfile(userID, req.Nickname, req.Avatar, req.Bio)
	return ctx.JSON(mapper.MessageResult(err))
}

// UserUpdatePassword 修改密码
func UserUpdatePassword(ctx *fiber.Ctx) error {
	userID := ctx.Locals("userID").(uint)

	var req UserUpdatePasswordRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.JSON(mapper.MessageFail("invalid request body"))
	}

	err := service.User.UpdatePassword(userID, req.OldPassword, req.NewPassword)
	return ctx.JSON(mapper.MessageResult(err))
}

// ============ 后台管理 API ============

// UserList 获取用户列表
func UserList(ctx *fiber.Ctx) error {
	var req struct {
		context.Context
		Keyword string `json:"keyword"`
	}
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.JSON(mapper.MessageFail("invalid request body"))
	}

	if req.Limit == 0 {
		req.Limit = 20
	}

	users, err := service.User.ListByKeyword(&req.Context, req.Keyword)
	if err != nil {
		return ctx.JSON(mapper.MessageResult(err))
	}

	return ctx.JSON(mapper.MessageResultData(users, nil))
}

// UserCount 统计用户数量
func UserCount(ctx *fiber.Ctx) error {
	var req struct {
		Keyword string `json:"keyword"`
	}
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.JSON(mapper.MessageFail("invalid request body"))
	}

	count, err := service.User.CountByKeyword(req.Keyword)
	if err != nil {
		return ctx.JSON(mapper.MessageResult(err))
	}

	return ctx.JSON(mapper.MessageResultData(count, nil))
}

// UserDetail 获取用户详情
func UserDetail(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		return ctx.JSON(mapper.MessageFail("invalid user id"))
	}

	user, err := service.User.GetByID(uint(id))
	if err != nil {
		return ctx.JSON(mapper.MessageResult(err))
	}

	return ctx.JSON(mapper.MessageResultData(user, nil))
}

// UserEnable 启用用户
func UserEnable(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		return ctx.JSON(mapper.MessageFail("invalid user id"))
	}

	err = service.User.Enable(uint(id))
	return ctx.JSON(mapper.MessageResult(err))
}

// UserDisable 禁用用户
func UserDisable(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		return ctx.JSON(mapper.MessageFail("invalid user id"))
	}

	err = service.User.Disable(uint(id))
	return ctx.JSON(mapper.MessageResult(err))
}

// UserDelete 删除用户
func UserDelete(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		return ctx.JSON(mapper.MessageFail("invalid user id"))
	}

	err = service.User.Delete(uint(id))
	return ctx.JSON(mapper.MessageResult(err))
}

// UserStats 用户统计信息
func UserStats(ctx *fiber.Ctx) error {
	total, err := service.User.Count()
	if err != nil {
		return ctx.JSON(mapper.MessageResult(err))
	}

	today, err := service.User.CountToday()
	if err != nil {
		return ctx.JSON(mapper.MessageResult(err))
	}

	return ctx.JSON(mapper.MessageResultData(&dto.UserStats{
		Total: total,
		Today: today,
	}, nil))
}

// ============ 邀请码管理 API ============

type InviteCodeCreateRequest struct {
	MaxUses    int `json:"max_uses"`    // 最大使用次数，0表示无限制
	ExpireDays int `json:"expire_days"` // 过期天数，0表示永不过期
}

// InviteCodeList 获取邀请码列表
func InviteCodeList(ctx *fiber.Ctx) error {
	var req context.Context
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.JSON(mapper.MessageFail("invalid request body"))
	}

	if req.Limit == 0 {
		req.Limit = 20
	}

	codes, err := service.InviteCode.List(&req)
	if err != nil {
		return ctx.JSON(mapper.MessageResult(err))
	}

	return ctx.JSON(mapper.MessageResultData(codes, nil))
}

// InviteCodeCount 统计邀请码数量
func InviteCodeCount(ctx *fiber.Ctx) error {
	count, err := service.InviteCode.Count()
	if err != nil {
		return ctx.JSON(mapper.MessageResult(err))
	}
	return ctx.JSON(mapper.MessageResultData(count, nil))
}

// InviteCodeCreate 创建邀请码
func InviteCodeCreate(ctx *fiber.Ctx) error {
	var req InviteCodeCreateRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.JSON(mapper.MessageFail("invalid request body"))
	}

	// createdBy 暂时使用 0，后续可以从 ctx.Locals 获取管理员 ID
	code, err := service.InviteCode.Create(req.MaxUses, req.ExpireDays, 0)
	if err != nil {
		return ctx.JSON(mapper.MessageResult(err))
	}

	return ctx.JSON(mapper.MessageResultData(code, nil))
}

// InviteCodeDelete 删除邀请码
func InviteCodeDelete(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		return ctx.JSON(mapper.MessageFail("invalid invite code id"))
	}

	err = service.InviteCode.Delete(uint(id))
	return ctx.JSON(mapper.MessageResult(err))
}

// InviteCodeValidate 验证邀请码（公开接口，用于注册时校验）
func InviteCodeValidate(ctx *fiber.Ctx) error {
	code := ctx.Query("code")
	if code == "" {
		return ctx.JSON(mapper.MessageFail("invite code is required"))
	}

	inviteCode, err := service.InviteCode.Validate(code)
	if err != nil {
		return ctx.JSON(mapper.MessageResult(err))
	}

	return ctx.JSON(mapper.MessageResultData(map[string]interface{}{
		"valid":      true,
		"code":       inviteCode.Code,
		"max_uses":   inviteCode.MaxUses,
		"used_count": inviteCode.UsedCount,
		"expire_at":  inviteCode.ExpireAt,
	}, nil))
}

// ============ 用户收藏 API ============

// UserFavoriteList 获取收藏列表
func UserFavoriteList(ctx *fiber.Ctx) error {
	userID := ctx.Locals("userID").(uint)

	var req context.Context
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.JSON(mapper.MessageFail("invalid request body"))
	}

	if req.Limit == 0 {
		req.Limit = 20
	}

	favorites, err := service.UserFavorite.ListByUserID(&req, userID)
	if err != nil {
		return ctx.JSON(mapper.MessageResult(err))
	}

	return ctx.JSON(mapper.MessageResultData(favorites, nil))
}

// UserFavoriteCount 统计收藏数量
func UserFavoriteCount(ctx *fiber.Ctx) error {
	userID := ctx.Locals("userID").(uint)

	count, err := service.UserFavorite.CountByUserID(userID)
	if err != nil {
		return ctx.JSON(mapper.MessageResult(err))
	}

	return ctx.JSON(mapper.MessageResultData(count, nil))
}

// UserFavoriteAdd 添加收藏
func UserFavoriteAdd(ctx *fiber.Ctx) error {
	userID := ctx.Locals("userID").(uint)

	var req struct {
		ArticleID uint `json:"article_id"`
	}
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.JSON(mapper.MessageFail("invalid request body"))
	}

	favorite, err := service.UserFavorite.Add(userID, req.ArticleID)
	if err != nil {
		return ctx.JSON(mapper.MessageResult(err))
	}

	return ctx.JSON(mapper.MessageResultData(favorite, nil))
}

// UserFavoriteRemove 取消收藏
func UserFavoriteRemove(ctx *fiber.Ctx) error {
	userID := ctx.Locals("userID").(uint)

	articleID, err := ctx.ParamsInt("article_id")
	if err != nil {
		return ctx.JSON(mapper.MessageFail("invalid article id"))
	}

	err = service.UserFavorite.Remove(userID, uint(articleID))
	return ctx.JSON(mapper.MessageResult(err))
}

// UserFavoriteToggle 切换收藏状态
func UserFavoriteToggle(ctx *fiber.Ctx) error {
	userID := ctx.Locals("userID").(uint)

	var req struct {
		ArticleID uint `json:"article_id"`
	}
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.JSON(mapper.MessageFail("invalid request body"))
	}

	isFavorited, err := service.UserFavorite.Toggle(userID, req.ArticleID)
	if err != nil {
		return ctx.JSON(mapper.MessageResult(err))
	}

	return ctx.JSON(mapper.MessageResultData(map[string]bool{
		"is_favorited": isFavorited,
	}, nil))
}

// UserFavoriteIsFavorited 检查是否已收藏
func UserFavoriteIsFavorited(ctx *fiber.Ctx) error {
	userID := ctx.Locals("userID").(uint)

	articleID, err := ctx.ParamsInt("article_id")
	if err != nil {
		return ctx.JSON(mapper.MessageFail("invalid article id"))
	}

	isFavorited, err := service.UserFavorite.IsFavorited(userID, uint(articleID))
	if err != nil {
		return ctx.JSON(mapper.MessageResult(err))
	}

	return ctx.JSON(mapper.MessageResultData(map[string]bool{
		"is_favorited": isFavorited,
	}, nil))
}

// ============ 邮箱验证 API ============

// UserVerifyEmail 验证邮箱（公开接口）
func UserVerifyEmail(ctx *fiber.Ctx) error {
	token := ctx.Query("token")
	if token == "" {
		return ctx.JSON(mapper.MessageFail("token is required"))
	}

	err := service.User.VerifyEmail(token)
	if err != nil {
		return ctx.JSON(mapper.MessageResult(err))
	}

	return ctx.JSON(mapper.MessageSuccess("email verified successfully"))
}

// UserResendVerificationEmail 重新发送验证邮件（需要登录）
func UserResendVerificationEmail(ctx *fiber.Ctx) error {
	userID := ctx.Locals("userID").(uint)

	user, err := service.User.GetByID(userID)
	if err != nil {
		return ctx.JSON(mapper.MessageResult(err))
	}

	if user.EmailVerified {
		return ctx.JSON(mapper.MessageFail("email already verified"))
	}

	if user.Email == "" {
		return ctx.JSON(mapper.MessageFail("no email address"))
	}

	err = service.User.ResendVerificationEmail(user.Email)
	if err != nil {
		return ctx.JSON(mapper.MessageResult(err))
	}

	return ctx.JSON(mapper.MessageSuccess("verification email sent"))
}
