package entity

import "time"

// LikeType represents the type of like (like or dislike)
type LikeType int

const (
	LikeTypeNone    LikeType = 0
	LikeTypeLike    LikeType = 1
	LikeTypeDislike LikeType = 2
)

// Like represents a user's like/dislike for an article
type Like struct {
	ID        uint      `gorm:"type:int;size:32;primaryKey;autoIncrement" json:"id"`
	UserID    uint      `gorm:"type:int;size:32;not null;uniqueIndex:idx_user_article" json:"user_id"`
	ArticleID int       `gorm:"type:int;size:32;not null;uniqueIndex:idx_user_article" json:"article_id"`
	Type      LikeType  `gorm:"type:tinyint;default:0" json:"type"` // 0=none, 1=like, 2=dislike
	CreatedAt time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (Like) TableName() string {
	return "like"
}
