package repository

import (
	"time"

	"moss/domain/core/entity"
	"moss/domain/core/repository/context"
	"moss/domain/core/repository/gormx"
	"moss/infrastructure/persistent/db"
)

var User = new(UserRepo)

type UserRepo struct{}

// MigrateTable 迁移表结构
func (r *UserRepo) MigrateTable() error {
	return db.DB.AutoMigrate(&entity.User{})
}

// Create 创建用户
func (r *UserRepo) Create(user *entity.User) error {
	return db.DB.Create(user).Error
}

// Update 更新用户
func (r *UserRepo) Update(user *entity.User) error {
	return db.DB.Save(user).Error
}

// UpdateFields 更新指定字段
func (r *UserRepo) UpdateFields(id uint, fields map[string]interface{}) error {
	return db.DB.Model(&entity.User{}).Where("id = ?", id).Updates(fields).Error
}

// Delete 删除用户
func (r *UserRepo) Delete(id uint) error {
	return db.DB.Delete(&entity.User{}, id).Error
}

// Get 根据ID获取用户
func (r *UserRepo) Get(id uint) (*entity.User, error) {
	var user entity.User
	err := db.DB.First(&user, id).Error
	return &user, err
}

// GetByUsername 根据用户名获取用户
func (r *UserRepo) GetByUsername(username string) (*entity.User, error) {
	var user entity.User
	err := db.DB.Where("username = ?", username).First(&user).Error
	return &user, err
}

// GetByEmail 根据邮箱获取用户
func (r *UserRepo) GetByEmail(email string) (*entity.User, error) {
	var user entity.User
	err := db.DB.Where("email = ?", email).First(&user).Error
	return &user, err
}

// GetByUsernameOrEmail 根据用户名或邮箱获取用户
func (r *UserRepo) GetByUsernameOrEmail(account string) (*entity.User, error) {
	var user entity.User
	err := db.DB.Where("username = ? OR email = ?", account, account).First(&user).Error
	return &user, err
}

// ExistsByUsername 检查用户名是否存在
func (r *UserRepo) ExistsByUsername(username string) (bool, error) {
	var count int64
	err := db.DB.Model(&entity.User{}).Where("username = ?", username).Count(&count).Error
	return count > 0, err
}

// ExistsByEmail 检查邮箱是否存在
func (r *UserRepo) ExistsByEmail(email string) (bool, error) {
	var count int64
	err := db.DB.Model(&entity.User{}).Where("email = ?", email).Count(&count).Error
	return count > 0, err
}

// Count 统计用户总数
func (r *UserRepo) Count() (int64, error) {
	var count int64
	err := db.DB.Model(&entity.User{}).Count(&count).Error
	return count, err
}

// CountByStatus 根据状态统计用户数量
func (r *UserRepo) CountByStatus(status int8) (int64, error) {
	var count int64
	err := db.DB.Model(&entity.User{}).Where("status = ?", status).Count(&count).Error
	return count, err
}

// CountToday 统计今日新增用户数
func (r *UserRepo) CountToday() (int64, error) {
	var count int64
	today := time.Now().Truncate(24 * time.Hour)
	err := db.DB.Model(&entity.User{}).Where("create_time >= ?", today.Unix()).Count(&count).Error
	return count, err
}

// List 获取用户列表
func (r *UserRepo) List(ctx *context.Context) ([]entity.User, error) {
	var users []entity.User
	err := db.DB.Model(&entity.User{}).Scopes(gormx.Context(ctx)).Find(&users).Error
	return users, err
}

// ListByStatus 根据状态获取用户列表
func (r *UserRepo) ListByStatus(ctx *context.Context, status int8) ([]entity.User, error) {
	var users []entity.User
	err := db.DB.Model(&entity.User{}).Where("status = ?", status).Scopes(gormx.Context(ctx)).Find(&users).Error
	return users, err
}

// ListByKeyword 根据关键词搜索用户
func (r *UserRepo) ListByKeyword(ctx *context.Context, keyword string) ([]entity.User, error) {
	var users []entity.User
	like := "%" + keyword + "%"
	err := db.DB.Model(&entity.User{}).
		Where("username LIKE ? OR email LIKE ? OR nickname LIKE ?", like, like, like).
		Scopes(gormx.Context(ctx)).
		Find(&users).Error
	return users, err
}

// CountByKeyword 根据关键词统计用户数
func (r *UserRepo) CountByKeyword(keyword string) (int64, error) {
	var count int64
	like := "%" + keyword + "%"
	err := db.DB.Model(&entity.User{}).
		Where("username LIKE ? OR email LIKE ? OR nickname LIKE ?", like, like, like).
		Count(&count).Error
	return count, err
}

// UpdateLastLogin 更新最后登录时间
func (r *UserRepo) UpdateLastLogin(id uint) error {
	return db.DB.Model(&entity.User{}).Where("id = ?", id).UpdateColumn("last_login", time.Now().Unix()).Error
}

// UpdateStatus 更新用户状态
func (r *UserRepo) UpdateStatus(id uint, status int8) error {
	return db.DB.Model(&entity.User{}).Where("id = ?", id).UpdateColumn("status", status).Error
}

// Enable 启用用户
func (r *UserRepo) Enable(id uint) error {
	return r.UpdateStatus(id, entity.UserStatusActive)
}

// Disable 禁用用户
func (r *UserRepo) Disable(id uint) error {
	return r.UpdateStatus(id, entity.UserStatusDisabled)
}

// GetByInviteCode 根据邀请码获取使用该邀请码的用户列表
func (r *UserRepo) GetByInviteCode(ctx *context.Context, inviteCode string) ([]entity.User, error) {
	var users []entity.User
	err := db.DB.Model(&entity.User{}).
		Where("invite_code = ?", inviteCode).
		Scopes(gormx.Context(ctx)).
		Find(&users).Error
	return users, err
}