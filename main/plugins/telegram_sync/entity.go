package telegram_sync

import (
	"time"
)

// TelegramChannel 频道配置
type TelegramChannel struct {
	ID             int    `gorm:"type:int;size:32;primaryKey;autoIncrement" json:"id"`
	ChannelID      int64  `gorm:"type:bigint;uniqueIndex;not null" json:"channel_id"`       // Telegram 频道 ID
	ChannelName    string `gorm:"type:varchar(200);default:''" json:"channel_name"`         // 频道名称
	ChannelLink    string `gorm:"type:varchar(300);default:''" json:"channel_link"`         // 频道链接
	Status         int    `gorm:"type:tinyint;default:1;index" json:"status"`               // 1=启用, 0=禁用
	CategoryID     int    `gorm:"type:int;size:32;default:0;index" json:"category_id"`      // 目标分类 ID
	ArticleStatus  bool   `gorm:"type:boolean;default:true" json:"article_status"`          // 同步文章发布状态
	ArticleAuthor  int    `gorm:"type:int;size:32;default:0" json:"article_author"`         // 文章作者 ID

	// 过滤规则
	FilterKeywords     string `gorm:"type:text" json:"filter_keywords"`                         // 关键词过滤规则 JSON
	FilterMessageTypes string `gorm:"type:varchar(100);default:'text,photo'" json:"filter_message_types"` // 消息类型过滤
	FilterMinLength    int    `gorm:"type:int;default:0" json:"filter_min_length"`              // 最小消息长度
	FilterMaxLength    int    `gorm:"type:int;default:0" json:"filter_max_length"`              // 最大消息长度

	// 统计信息
	LastSyncTime   int64 `gorm:"type:bigint;default:0" json:"last_sync_time"`   // 最后同步时间戳
	LastMessageID  int64 `gorm:"type:bigint;default:0" json:"last_message_id"`  // 最后处理的消息 ID
	TotalSyncCount int   `gorm:"type:int;default:0" json:"total_sync_count"`    // 总同步数量
	ErrorCount     int   `gorm:"type:int;default:0" json:"error_count"`         // 错误计数

	// 元数据
	CreateTime int64  `gorm:"type:bigint;default:0" json:"create_time"`
	UpdateTime int64  `gorm:"type:bigint;default:0" json:"update_time"`
	Remark     string `gorm:"type:varchar(500);default:''" json:"remark"`
}

func (TelegramChannel) TableName() string {
	return "telegram_channel"
}

// TelegramSyncLog 同步日志
type TelegramSyncLog struct {
	ID             int    `gorm:"type:int;size:32;primaryKey;autoIncrement" json:"id"`
	ChannelID      int64  `gorm:"type:bigint;index;not null" json:"channel_id"`
	MessageID      int64  `gorm:"type:bigint;index;not null" json:"message_id"`
	ArticleID      int    `gorm:"type:int;size:32;default:0;index" json:"article_id"`
	Status         int    `gorm:"type:tinyint;default:0;index" json:"status"` // 0=失败, 1=成功, 2=跳过
	ErrorMessage   string `gorm:"type:varchar(500);default:''" json:"error_message"`
	MessageTitle   string `gorm:"type:varchar(250);default:''" json:"message_title"`
	MessageContent string `gorm:"type:text" json:"message_content"`
	CreateTime     int64  `gorm:"type:bigint;default:0;index" json:"create_time"`
}

func (TelegramSyncLog) TableName() string {
	return "telegram_sync_log"
}

// TelegramSession 会话信息
type TelegramSession struct {
	ID          int    `gorm:"type:int;size:32;primaryKey;autoIncrement" json:"id"`
	SessionData []byte `gorm:"type:blob" json:"-"`                              // 加密的会话数据
	SessionHash string `gorm:"type:varchar(64);uniqueIndex" json:"session_hash"` // 会话哈希
	Status      int    `gorm:"type:tinyint;default:1" json:"status"`             // 1=有效, 0=过期
	CreateTime  int64  `gorm:"type:bigint;default:0" json:"create_time"`
	UpdateTime  int64  `gorm:"type:bigint;default:0" json:"update_time"`
}

func (TelegramSession) TableName() string {
	return "telegram_session"
}

// TelegramMedia 媒体文件存储
type TelegramMedia struct {
	ID           int    `gorm:"type:int;size:32;primaryKey;autoIncrement" json:"id"`
	MediaID      int64  `gorm:"type:bigint;uniqueIndex;not null" json:"media_id"`       // Telegram 文件 ID (photo ID 或 document ID)
	MediaType    string `gorm:"type:varchar(20);not null" json:"media_type"`            // photo, document, video
	AccessHash   int64  `gorm:"type:bigint;not null" json:"access_hash"`                // 访问哈希（用于下载）
	FileReference []byte `gorm:"type:blob" json:"-"`                                     // 文件引用（用于下载）

	// 来源信息
	ChannelID    int64  `gorm:"type:bigint;index" json:"channel_id"`                    // 来源频道 ID
	MessageID    int64  `gorm:"type:bigint;index" json:"message_id"`                    // 来源消息 ID

	// 文件信息
	Filename     string `gorm:"type:varchar(250);default:''" json:"filename"`           // 原始文件名
	MimeType     string `gorm:"type:varchar(100);default:''" json:"mime_type"`          // MIME 类型
	Size         int64  `gorm:"type:bigint;default:0" json:"size"`                      // 文件大小
	Width        int    `gorm:"type:int;default:0" json:"width"`                        // 图片宽度
	Height       int    `gorm:"type:int;default:0" json:"height"`                       // 图片高度

	// 存储信息
	StoragePath  string `gorm:"type:varchar(500);default:''" json:"storage_path"`       // 存储路径
	StorageURL   string `gorm:"type:varchar(500);default:''" json:"storage_url"`        // 访问 URL

	// 元数据
	CreateTime   int64  `gorm:"type:bigint;default:0;index" json:"create_time"`
}

func (TelegramMedia) TableName() string {
	return "telegram_media"
}

// FilterRule 过滤规则 JSON 结构
type FilterRule struct {
	Type          string   `json:"type"`            // whitelist, blacklist
	Keywords      []string `json:"keywords"`        // 关键词列表
	CaseSensitive bool     `json:"case_sensitive"`  // 是否区分大小写
	MatchAll      bool     `json:"match_all"`       // 白名单模式是否需要匹配所有
}

// ChannelConfig 频道配置（用于 JSON 存储）
type ChannelConfig struct {
	ChannelID          int64  `json:"channel_id"`
	ChannelName        string `json:"channel_name"`
	ChannelLink        string `json:"channel_link"`
	Status             int    `json:"status"`
	CategoryID         int    `json:"category_id"`
	ArticleStatus      bool   `json:"article_status"`
	FilterKeywords     string `json:"filter_keywords"`      // 改为字符串，前端存储原始值
	FilterMessageTypes string `json:"filter_message_types"`
	FilterMinLength    int    `json:"filter_min_length"`
	FilterMaxLength    int    `json:"filter_max_length"`
	AccessHash         int64  `json:"access_hash"` // 频道 access hash，用于 updates.Manager
}

// SyncLogItem 同步日志项（用于前端显示）
type SyncLogItem struct {
	ChannelName   string `json:"channel_name"`
	MessageID     int64  `json:"message_id"`
	ArticleID     int    `json:"article_id"`
	Status        int    `json:"status"`
	MessageTitle  string `json:"message_title"`
	CreateTime    int64  `json:"create_time"`
}

// MessageMediaInfo 消息媒体信息（用于文章创建）
type MessageMediaInfo struct {
	MediaID        int64  `json:"media_id"`
	MediaType      string `json:"media_type"`
	URL            string `json:"url"`
	AccessHash     int64  `json:"access_hash"`
	FileReference  []byte `json:"file_reference"`
	Filename       string `json:"filename"`
	Width          int    `json:"width"`
	Height         int    `json:"height"`
	// 视频封面信息
	ThumbURL        string `json:"thumb_url"`
	ThumbID         int64  `json:"thumb_id"`
	ThumbAccessHash int64  `json:"thumb_access_hash"`
	ThumbFileRef    []byte `json:"thumb_file_ref"`
}

// GetCreateTimeFormat 格式化创建时间
func (c *TelegramChannel) GetCreateTimeFormat(layouts ...string) string {
	if c.CreateTime == 0 {
		return ""
	}
	layout := "2006-01-02 15:04:05"
	if len(layouts) > 0 && layouts[0] != "" {
		layout = layouts[0]
	}
	return time.Unix(c.CreateTime, 0).Format(layout)
}

// GetLastSyncTimeFormat 格式化最后同步时间
func (c *TelegramChannel) GetLastSyncTimeFormat(layouts ...string) string {
	if c.LastSyncTime == 0 {
		return ""
	}
	layout := "2006-01-02 15:04:05"
	if len(layouts) > 0 && layouts[0] != "" {
		layout = layouts[0]
	}
	return time.Unix(c.LastSyncTime, 0).Format(layout)
}
