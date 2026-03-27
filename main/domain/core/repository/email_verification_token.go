package repository

import (
	"time"

	"moss/domain/core/entity"
	"moss/infrastructure/persistent/db"
)

var EmailVerificationToken = new(EmailVerificationTokenRepo)

type EmailVerificationTokenRepo struct{}

// MigrateTable 迁移表结构
func (r *EmailVerificationTokenRepo) MigrateTable() error {
	return db.DB.AutoMigrate(&entity.EmailVerificationToken{})
}

// Create 创建令牌
func (r *EmailVerificationTokenRepo) Create(token *entity.EmailVerificationToken) error {
	return db.DB.Create(token).Error
}

// Delete 删除令牌
func (r *EmailVerificationTokenRepo) Delete(id uint) error {
	return db.DB.Delete(&entity.EmailVerificationToken{}, id).Error
}

// DeleteByUserID 根据用户ID删除所有令牌
func (r *EmailVerificationTokenRepo) DeleteByUserID(userID uint) error {
	return db.DB.Where("user_id = ?", userID).Delete(&entity.EmailVerificationToken{}).Error
}

// GetByToken 根据令牌获取
func (r *EmailVerificationTokenRepo) GetByToken(token string) (*entity.EmailVerificationToken, error) {
	var t entity.EmailVerificationToken
	err := db.DB.Where("token = ?", token).First(&t).Error
	return &t, err
}

// DeleteExpired 删除过期的令牌
func (r *EmailVerificationTokenRepo) DeleteExpired() error {
	return db.DB.Where("expire_at < ?", time.Now().Unix()).Delete(&entity.EmailVerificationToken{}).Error
}