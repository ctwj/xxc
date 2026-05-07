package entity

import "time"

// User represents a registered user
type User struct {
	ID        uint      `gorm:"type:int;size:32;primaryKey;autoIncrement" json:"id"`
	Username  string    `gorm:"type:varchar(100);uniqueIndex;not null" json:"username"`
	Email     string    `gorm:"type:varchar(150);uniqueIndex;not null" json:"email"`
	Password  string    `gorm:"type:varchar(250);not null" json:"-"`
	Role      string    `gorm:"type:varchar(20);default:'user'" json:"role"` // user, admin
	CreatedAt time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (User) TableName() string {
	return "user"
}
