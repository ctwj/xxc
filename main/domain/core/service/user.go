package service

import (
	"errors"
	"fmt"
	"time"

	"moss/domain/config"
	"moss/domain/core/entity"
	"moss/domain/core/repository"
	"moss/domain/core/repository/context"
	supportService "moss/domain/support/service"
	"moss/infrastructure/general/message"
)

var User = new(UserService)

type UserService struct{}

// Register 用户注册
func (s *UserService) Register(username, email, password, nickname, inviteCode string) (*entity.User, error) {
	// 验证必填字段
	if username == "" {
		return nil, errors.New("username is required")
	}
	if password == "" {
		return nil, errors.New("password is required")
	}

	// 检查用户名是否已存在
	exists, err := repository.User.ExistsByUsername(username)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("username already exists")
	}

	// 检查邮箱是否已存在（如果提供了邮箱）
	if email != "" {
		exists, err = repository.User.ExistsByEmail(email)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, errors.New("email already exists")
		}
	}

	// 如果启用了邮箱验证但没有提供邮箱
	emailVerifyEnabled := config.Config.Email.Enable && config.Config.Email.VerifyEmail
	if emailVerifyEnabled && email == "" {
		return nil, errors.New("email is required for verification")
	}

	// 创建用户实体
	user := &entity.User{
		Username:      username,
		Email:         email,
		Nickname:      nickname,
		Status:        entity.UserStatusActive,
		EmailVerified: !emailVerifyEnabled, // 如果不需要验证，则直接设为已验证
		InviteCode:    inviteCode,
		CreateTime:    time.Now().Unix(),
	}

	// 设置密码
	if err := user.SetPassword(password); err != nil {
		return nil, err
	}

	// 保存用户
	if err := repository.User.Create(user); err != nil {
		return nil, err
	}

	// 如果需要邮箱验证，发送验证邮件
	if emailVerifyEnabled && email != "" {
		if err := s.SendVerificationEmail(user.ID, email); err != nil {
			// 发送失败不影响注册，但记录错误
			fmt.Printf("Failed to send verification email: %v\n", err)
		}
	}

	return user, nil
}

// SendVerificationEmail 发送验证邮件
func (s *UserService) SendVerificationEmail(userID uint, email string) error {
	// 删除该用户之前的验证令牌
	_ = repository.EmailVerificationToken.DeleteByUserID(userID)

	// 生成新的验证令牌
	token := entity.GenerateVerificationToken(userID, email)
	if err := repository.EmailVerificationToken.Create(token); err != nil {
		return err
	}

	// 构建验证链接
	verifyLink := fmt.Sprintf("%s/api/user/verify-email?token=%s", config.Config.Site.GetURL(), token.Token)

	// 发送邮件
	return supportService.SendVerificationEmail(email, verifyLink)
}

// VerifyEmail 验证邮箱
func (s *UserService) VerifyEmail(tokenStr string) error {
	if tokenStr == "" {
		return errors.New("token is required")
	}

	// 查找令牌
	token, err := repository.EmailVerificationToken.GetByToken(tokenStr)
	if err != nil {
		return errors.New("invalid or expired token")
	}

	if token.ID == 0 {
		return errors.New("invalid token")
	}

	if token.IsExpired() {
		return errors.New("token has expired")
	}

	// 获取用户
	user, err := repository.User.Get(token.UserID)
	if err != nil {
		return errors.New("user not found")
	}

	if user.ID == 0 {
		return errors.New("user not found")
	}

	// 更新用户邮箱验证状态
	if err := repository.User.UpdateFields(user.ID, map[string]interface{}{
		"email_verified": true,
		"update_time":    time.Now().Unix(),
	}); err != nil {
		return err
	}

	// 删除验证令牌
	_ = repository.EmailVerificationToken.Delete(token.ID)

	return nil
}

// ResendVerificationEmail 重新发送验证邮件
func (s *UserService) ResendVerificationEmail(email string) error {
	if email == "" {
		return errors.New("email is required")
	}

	user, err := repository.User.GetByEmail(email)
	if err != nil {
		return errors.New("user not found")
	}

	if user.ID == 0 {
		return errors.New("user not found")
	}

	if user.EmailVerified {
		return errors.New("email already verified")
	}

	return s.SendVerificationEmail(user.ID, email)
}

// Login 用户登录
func (s *UserService) Login(account, password string) (*entity.User, string, error) {
	if account == "" {
		return nil, "", errors.New("account is required")
	}
	if password == "" {
		return nil, "", errors.New("password is required")
	}

	// 查找用户
	user, err := repository.User.GetByUsernameOrEmail(account)
	if err != nil {
		return nil, "", errors.New("account or password error")
	}

	// 检查用户是否存在
	if user.ID == 0 {
		return nil, "", errors.New("account or password error")
	}

	// 检查用户状态
	if !user.IsActive() {
		return nil, "", errors.New("account is disabled")
	}

	// 检查邮箱验证
	emailVerifyEnabled := config.Config.Email.Enable && config.Config.Email.VerifyEmail
	if emailVerifyEnabled && !user.EmailVerified {
		return nil, "", errors.New("please verify your email first")
	}

	// 验证密码
	if !user.VerifyPassword(password) {
		return nil, "", errors.New("account or password error")
	}

	// 更新最后登录时间
	_ = repository.User.UpdateLastLogin(user.ID)

	// 生成 JWT token
	token := user.GenerateJwtToken()

	return user, token, nil
}

// GetByID 根据ID获取用户
func (s *UserService) GetByID(id uint) (*entity.User, error) {
	if id == 0 {
		return nil, message.ErrIdRequired
	}
	return repository.User.Get(id)
}

// GetByUsername 根据用户名获取用户
func (s *UserService) GetByUsername(username string) (*entity.User, error) {
	if username == "" {
		return nil, errors.New("username is required")
	}
	return repository.User.GetByUsername(username)
}

// UpdateProfile 更新用户资料
func (s *UserService) UpdateProfile(id uint, nickname, avatar, bio string) error {
	user, err := repository.User.Get(id)
	if err != nil {
		return err
	}
	if user.ID == 0 {
		return message.ErrRecordNotFound
	}

	fields := map[string]interface{}{
		"update_time": time.Now().Unix(),
	}

	if nickname != "" {
		fields["nickname"] = nickname
	}
	if avatar != "" {
		fields["avatar"] = avatar
	}
	if bio != "" {
		fields["bio"] = bio
	}

	return repository.User.UpdateFields(id, fields)
}

// UpdatePassword 修改密码
func (s *UserService) UpdatePassword(id uint, oldPassword, newPassword string) error {
	if oldPassword == "" {
		return errors.New("old password is required")
	}
	if newPassword == "" {
		return errors.New("new password is required")
	}

	user, err := repository.User.Get(id)
	if err != nil {
		return err
	}
	if user.ID == 0 {
		return message.ErrRecordNotFound
	}

	// 验证旧密码
	if !user.VerifyPassword(oldPassword) {
		return errors.New("old password is incorrect")
	}

	// 设置新密码
	if err := user.SetPassword(newPassword); err != nil {
		return err
	}

	return repository.User.UpdateFields(id, map[string]interface{}{
		"password":    user.Password,
		"update_time": time.Now().Unix(),
	})
}

// List 获取用户列表
func (s *UserService) List(ctx *context.Context) ([]entity.User, error) {
	return repository.User.List(ctx)
}

// ListByKeyword 根据关键词搜索用户
func (s *UserService) ListByKeyword(ctx *context.Context, keyword string) ([]entity.User, error) {
	if keyword == "" {
		return s.List(ctx)
	}
	return repository.User.ListByKeyword(ctx, keyword)
}

// Count 用户总数
func (s *UserService) Count() (int64, error) {
	return repository.User.Count()
}

// CountByKeyword 根据关键词统计用户数
func (s *UserService) CountByKeyword(keyword string) (int64, error) {
	if keyword == "" {
		return s.Count()
	}
	return repository.User.CountByKeyword(keyword)
}

// CountToday 今日新增用户数
func (s *UserService) CountToday() (int64, error) {
	return repository.User.CountToday()
}

// Enable 启用用户
func (s *UserService) Enable(id uint) error {
	return repository.User.Enable(id)
}

// Disable 禁用用户
func (s *UserService) Disable(id uint) error {
	return repository.User.Disable(id)
}

// Delete 删除用户
func (s *UserService) Delete(id uint) error {
	return repository.User.Delete(id)
}