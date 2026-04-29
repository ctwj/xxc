package entity

import "time"

// ViewHistory represents a user's view history for an article
type ViewHistory struct {
	ID        uint      `gorm:"type:int;size:32;primaryKey;autoIncrement" json:"id"`
	UserID    uint      `gorm:"type:int;size:32;not null;index" json:"user_id"`
	ArticleID int       `gorm:"type:int;size:32;not null;index" json:"article_id"`
	ViewedAt  time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP;index" json:"viewed_at"`
}

func (ViewHistory) TableName() string {
	return "view_history"
}
