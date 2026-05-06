package entity

import "time"

// Favorite represents a user's favorite article
type Favorite struct {
	ID        uint      `gorm:"type:int;size:32;primaryKey;autoIncrement" json:"id"`
	UserID    uint      `gorm:"type:int;size:32;not null;index"           json:"user_id"`
	ArticleID int       `gorm:"type:int;size:32;not null;index"           json:"article_id"`
	CreatedAt time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP"  json:"created_at"`
}

func (Favorite) TableName() string {
	return "favorite"
}