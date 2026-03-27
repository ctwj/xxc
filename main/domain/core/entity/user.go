package entity

import (
	"errors"
	"time"

	"github.com/brianvoe/sjwt"
	"github.com/duke-git/lancet/v2/random"
	"golang.org/x/crypto/bcrypt"
	"moss/infrastructure/utils/timex"
)

// User 用户实体
type User struct {
	ID            uint   `gorm:"type:int;size:32;primaryKey;autoIncrement" json:"id"`
	Username      string `gorm:"type:varchar(50);uniqueIndex;not null" json:"username"`
	Email         string `gorm:"type:varchar(100);uniqueIndex" json:"email"`
	Password      string `gorm:"type:varchar(255);not null" json:"-"`
	Nickname      string `gorm:"type:varchar(50);default:''" json:"nickname"`
	Avatar        string `gorm:"type:varchar(255);default:''" json:"avatar"`
	Bio           string `gorm:"type:varchar(500);default:''" json:"bio"`
	Status        int8   `gorm:"type:smallint;default:1;index" json:"status"` // 0:禁用 1:正常
	EmailVerified bool   `gorm:"type:boolean;default:false;index" json:"email_verified"` // 邮箱是否已验证
	InviteCode    string `gorm:"type:varchar(20);default:''" json:"invite_code"`
	InvitedBy     uint   `gorm:"type:int;size:32;default:0" json:"invited_by"`
	LastLogin     int64  `gorm:"type:bigint;default:0" json:"last_login"`
	CreateTime    int64  `gorm:"type:bigint;index" json:"create_time"`
	UpdateTime    int64  `gorm:"type:bigint;default:0" json:"update_time"`
}

func (User) TableName() string {
	return "users"
}

// UserStatus 用户状态常量
const (
	UserStatusDisabled = 0
	UserStatusActive   = 1
)

// JWT 配置（用户独立 JWT Key）
var userJwtKey = random.RandString(20)

// UserLoginExpire 登录过期时间配置
var UserLoginExpire = timex.Duration{Number: 7, Unit: "day"}

// EncryptPassword 加密密码
func EncryptPassword(password string) (string, error) {
	if password == "" {
		return "", errors.New("password is required")
	}
	b, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// VerifyPassword 验证密码
func (u *User) VerifyPassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)) == nil
}

// SetPassword 设置密码（加密存储）
func (u *User) SetPassword(password string) error {
	p, err := EncryptPassword(password)
	if err != nil {
		return err
	}
	u.Password = p
	return nil
}

// GenerateJwtToken 生成 JWT token
func (u *User) GenerateJwtToken() string {
	claims := sjwt.New()
	claims.Set("user_id", u.ID)
	claims.Set("username", u.Username)
	claims.SetIssuedAt(time.Now())
	if d := UserLoginExpire.Duration(); d > 0 {
		claims.SetExpiresAt(time.Now().Add(d))
	}
	return claims.Generate([]byte(userJwtKey))
}

// VerifyJwtToken 验证 JWT token
func VerifyJwtToken(token string) (userID uint, username string, ok bool) {
	if !sjwt.Verify(token, []byte(userJwtKey)) {
		return 0, "", false
	}
	claims, err := sjwt.Parse(token)
	if err != nil {
		return 0, "", false
	}
	if err := claims.Validate(); err != nil {
		return 0, "", false
	}
	userIDFloat, _ := claims.Get("user_id")
	usernameVal, _ := claims.Get("username")
	if userIDFloat == nil || usernameVal == nil {
		return 0, "", false
	}
	return uint(userIDFloat.(float64)), usernameVal.(string), true
}

// IsActive 检查用户是否激活
func (u *User) IsActive() bool {
	return u.Status == UserStatusActive
}

// LastLoginFormat 格式化最后登录时间
func (u *User) LastLoginFormat(layouts ...string) string {
	if u.LastLogin == 0 {
		return ""
	}
	layout := "2006-01-02 15:04:05"
	if len(layouts) > 0 && len(layouts[0]) > 0 {
		layout = layouts[0]
	}
	return time.Unix(u.LastLogin, 0).Format(layout)
}

// CreateTimeFormat 格式化创建时间
func (u *User) CreateTimeFormat(layouts ...string) string {
	if u.CreateTime == 0 {
		return ""
	}
	layout := "2006-01-02 15:04:05"
	if len(layouts) > 0 && len(layouts[0]) > 0 {
		layout = layouts[0]
	}
	return time.Unix(u.CreateTime, 0).Format(layout)
}
