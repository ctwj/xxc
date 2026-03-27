package entity

import "time"

// UserFavorite 用户收藏实体
type UserFavorite struct {
	ID         uint `gorm:"type:int;size:32;primaryKey;autoIncrement" json:"id"`
	UserID     uint `gorm:"type:int;size:32;not null;index" json:"user_id"`
	ArticleID  uint `gorm:"type:int;size:32;not null;index" json:"article_id"`
	CreateTime int64 `gorm:"type:int;size:32" json:"create_time"`
}

func (UserFavorite) TableName() string {
	return "user_favorite"
}

// CreateTimeFormat 格式化收藏时间
func (f *UserFavorite) CreateTimeFormat(layouts ...string) string {
	if f.CreateTime == 0 {
		return ""
	}
	layout := "2006-01-02 15:04:05"
	if len(layouts) > 0 && len(layouts[0]) > 0 {
		layout = layouts[0]
	}
	return time.Unix(f.CreateTime, 0).Format(layout)
}

// UserFavoriteWithArticle 用户收藏带文章信息
type UserFavoriteWithArticle struct {
	UserFavorite
	Article *ArticleBase `json:"article" gorm:"-"`
}
