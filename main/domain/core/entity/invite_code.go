package entity

import (
	"time"

	"github.com/duke-git/lancet/v2/random"
)

// InviteCode 邀请码实体
type InviteCode struct {
	ID         uint   `gorm:"type:int;size:32;primaryKey;autoIncrement" json:"id"`
	Code       string `gorm:"type:varchar(20);uniqueIndex;not null" json:"code"`
	MaxUses    int    `gorm:"type:int;default:0" json:"max_uses"`     // 最大使用次数，0表示无限制
	UsedCount  int    `gorm:"type:int;default:0" json:"used_count"`   // 已使用次数
	ExpireAt   int64  `gorm:"type:int;size:32;default:0" json:"expire_at"` // 过期时间，0表示永不过期
	CreatedBy  uint   `gorm:"type:int;size:32;default:0" json:"created_by"` // 创建人（管理员ID）
	CreateTime int64  `gorm:"type:int;size:32" json:"create_time"`
}

func (InviteCode) TableName() string {
	return "invite_code"
}

// GenerateCode 生成随机邀请码
func GenerateInviteCode() string {
	return random.RandString(8)
}

// IsExpired 检查邀请码是否过期
func (i *InviteCode) IsExpired() bool {
	if i.ExpireAt == 0 {
		return false
	}
	return time.Now().Unix() > i.ExpireAt
}

// IsUsedUp 检查邀请码是否已用完
func (i *InviteCode) IsUsedUp() bool {
	if i.MaxUses == 0 {
		return false
	}
	return i.UsedCount >= i.MaxUses
}

// IsValid 检查邀请码是否有效
func (i *InviteCode) IsValid() bool {
	return !i.IsExpired() && !i.IsUsedUp()
}

// IncrementUsed 增加使用次数
func (i *InviteCode) IncrementUsed() {
	i.UsedCount++
}

// ExpireAtFormat 格式化过期时间
func (i *InviteCode) ExpireAtFormat(layouts ...string) string {
	if i.ExpireAt == 0 {
		return "永不过期"
	}
	layout := "2006-01-02 15:04:05"
	if len(layouts) > 0 && len(layouts[0]) > 0 {
		layout = layouts[0]
	}
	return time.Unix(i.ExpireAt, 0).Format(layout)
}

// CreateTimeFormat 格式化创建时间
func (i *InviteCode) CreateTimeFormat(layouts ...string) string {
	if i.CreateTime == 0 {
		return ""
	}
	layout := "2006-01-02 15:04:05"
	if len(layouts) > 0 && len(layouts[0]) > 0 {
		layout = layouts[0]
	}
	return time.Unix(i.CreateTime, 0).Format(layout)
}
