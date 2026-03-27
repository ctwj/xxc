package entity

import (
	"time"

	"github.com/duke-git/lancet/v2/random"
)

// EmailVerificationToken 邮箱验证令牌
type EmailVerificationToken struct {
	ID        uint   `gorm:"type:int;size:32;primaryKey;autoIncrement" json:"id"`
	UserID    uint   `gorm:"type:int;size:32;not null;index" json:"user_id"`
	Email     string `gorm:"type:varchar(100);not null;index" json:"email"`
	Token     string `gorm:"type:varchar(64);uniqueIndex;not null" json:"token"`
	ExpireAt  int64  `gorm:"type:bigint;not null" json:"expire_at"`
	CreatedAt int64  `gorm:"type:bigint;not null" json:"created_at"`
}

func (EmailVerificationToken) TableName() string {
	return "email_verification_tokens"
}

// GenerateToken 生成验证令牌
func GenerateVerificationToken(userID uint, email string) *EmailVerificationToken {
	return &EmailVerificationToken{
		UserID:    userID,
		Email:     email,
		Token:     random.RandString(32),
		ExpireAt:  time.Now().Add(24 * time.Hour).Unix(), // 24小时有效期
		CreatedAt: time.Now().Unix(),
	}
}

// IsExpired 检查令牌是否过期
func (t *EmailVerificationToken) IsExpired() bool {
	return time.Now().Unix() > t.ExpireAt
}

// IsValid 检查令牌是否有效
func (t *EmailVerificationToken) IsValid() bool {
	return !t.IsExpired()
}
