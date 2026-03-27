package repository

import (
	"time"

	"moss/domain/core/entity"
	"moss/domain/core/repository/context"
	"moss/domain/core/repository/gormx"
	"moss/infrastructure/persistent/db"
)

var InviteCode = new(InviteCodeRepo)

type InviteCodeRepo struct{}

// MigrateTable 迁移表结构
func (r *InviteCodeRepo) MigrateTable() error {
	return db.DB.AutoMigrate(&entity.InviteCode{})
}

// Create 创建邀请码
func (r *InviteCodeRepo) Create(code *entity.InviteCode) error {
	return db.DB.Create(code).Error
}

// Update 更新邀请码
func (r *InviteCodeRepo) Update(code *entity.InviteCode) error {
	return db.DB.Save(code).Error
}

// Delete 删除邀请码
func (r *InviteCodeRepo) Delete(id uint) error {
	return db.DB.Delete(&entity.InviteCode{}, id).Error
}

// Get 根据ID获取邀请码
func (r *InviteCodeRepo) Get(id uint) (*entity.InviteCode, error) {
	var code entity.InviteCode
	err := db.DB.First(&code, id).Error
	return &code, err
}

// GetByCode 根据邀请码字符串获取
func (r *InviteCodeRepo) GetByCode(code string) (*entity.InviteCode, error) {
	var inviteCode entity.InviteCode
	err := db.DB.Where("code = ?", code).First(&inviteCode).Error
	return &inviteCode, err
}

// ExistsByCode 检查邀请码是否存在
func (r *InviteCodeRepo) ExistsByCode(code string) (bool, error) {
	var count int64
	err := db.DB.Model(&entity.InviteCode{}).Where("code = ?", code).Count(&count).Error
	return count > 0, err
}

// Count 统计邀请码总数
func (r *InviteCodeRepo) Count() (int64, error) {
	var count int64
	err := db.DB.Model(&entity.InviteCode{}).Count(&count).Error
	return count, err
}

// List 获取邀请码列表
func (r *InviteCodeRepo) List(ctx *context.Context) ([]entity.InviteCode, error) {
	var codes []entity.InviteCode
	err := db.DB.Model(&entity.InviteCode{}).Scopes(gormx.Context(ctx)).Find(&codes).Error
	return codes, err
}

// IncrementUsed 增加使用次数
func (r *InviteCodeRepo) IncrementUsed(id uint) error {
	return db.DB.Model(&entity.InviteCode{}).Where("id = ?", id).
		UpdateColumn("used_count", db.DB.Raw("used_count + 1")).Error
}

// GetValidCodes 获取有效的邀请码列表
func (r *InviteCodeRepo) GetValidCodes(ctx *context.Context) ([]entity.InviteCode, error) {
	var codes []entity.InviteCode
	now := time.Now().Unix()
	err := db.DB.Model(&entity.InviteCode{}).
		Where("expire_at = 0 OR expire_at > ?", now).
		Where("max_uses = 0 OR used_count < max_uses").
		Scopes(gormx.Context(ctx)).
		Find(&codes).Error
	return codes, err
}
