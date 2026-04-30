# Data Model: Telegram Channel Sync

**Date**: 2026-04-30
**Feature**: 002-telegram-channel-sync

## 架构说明

**重要**: 采用简化架构，所有实体定义放在 `main/plugins/TelegramChannelSync.go` 或 `main/plugins/telegram_sync/` 子包中，不创建独立的 `domain/core/entity/` 文件。这与 Moss 现有插件（如 `SaveArticleImages`、`GnDownSpider`）的模式一致。

## 实体定义

### 1. TelegramChannel（频道配置）

代表一个要监听的 Telegram 频道配置。

```go
type TelegramChannel struct {
    ID             int       `gorm:"type:int;size:32;primaryKey;autoIncrement" json:"id"`
    ChannelID      int64     `gorm:"type:bigint;uniqueIndex;not null" json:"channel_id"`      // Telegram 频道 ID
    ChannelName    string    `gorm:"type:varchar(200);default:''" json:"channel_name"`        // 频道名称（用于显示）
    ChannelLink    string    `gorm:"type:varchar(300);default:''" json:"channel_link"`        // 频道链接（可选，用于识别）
    Status         int       `gorm:"type:tinyint;default:1;index" json:"status"`              // 状态: 1=启用, 0=禁用
    CategoryID     int       `gorm:"type:int;size:32;default:0;index" json:"category_id"`     // 目标文章分类 ID
    ArticleStatus  bool      `gorm:"type:boolean;default:true" json:"article_status"`         // 同步文章的发布状态
    ArticleAuthor  int       `gorm:"type:int;size:32;default:0" json:"article_author"`        // 文章作者 ID
    
    // 过滤规则（JSON 存储）
    FilterKeywords     string `gorm:"type:text" json:"filter_keywords"`           // 关键词过滤规则 JSON
    FilterMessageTypes string `gorm:"type:varchar(100);default:'text,photo'" json:"filter_message_types"` // 消息类型过滤
    FilterMinLength    int    `gorm:"type:int;default:0" json:"filter_min_length"` // 最小消息长度
    FilterMaxLength    int    `gorm:"type:int;default:0" json:"filter_max_length"` // 最大消息长度
    
    // 统计信息
    LastSyncTime   int64     `gorm:"type:bigint;default:0" json:"last_sync_time"`    // 最后同步时间戳
    LastMessageID  int64     `gorm:"type:bigint;default:0" json:"last_message_id"`   // 最后处理的消息 ID
    TotalSyncCount int       `gorm:"type:int;default:0" json:"total_sync_count"`     // 总同步数量
    ErrorCount     int       `gorm:"type:int;default:0" json:"error_count"`          // 错误计数
    
    // 元数据
    CreateTime     int64     `gorm:"type:bigint;default:0" json:"create_time"`       // 创建时间
    UpdateTime     int64     `gorm:"type:bigint;default:0" json:"update_time"`       // 更新时间
    Remark         string    `gorm:"type:varchar(500);default:''" json:"remark"`     // 备注
}

func (TelegramChannel) TableName() string {
    return "telegram_channel"
}
```

**字段说明**:

| 字段 | 类型 | 说明 |
|------|------|------|
| ChannelID | int64 | Telegram 频道 ID（负数，如 -1001234567890） |
| ChannelName | string | 频道显示名称 |
| Status | int | 1=启用监听, 0=禁用 |
| CategoryID | int | 同步文章的目标分类 |
| FilterKeywords | string | JSON 格式的关键词过滤规则 |
| LastMessageID | int64 | 用于断点续传 |

### 2. TelegramSyncLog（同步日志）

记录每次同步操作的详细信息。

```go
type TelegramSyncLog struct {
    ID           int    `gorm:"type:int;size:32;primaryKey;autoIncrement" json:"id"`
    ChannelID    int64  `gorm:"type:bigint;index;not null" json:"channel_id"`      // 频道 ID
    MessageID    int64  `gorm:"type:bigint;index;not null" json:"message_id"`      // Telegram 消息 ID
    ArticleID    int    `gorm:"type:int;size:32;default:0;index" json:"article_id"` // 生成的文章 ID
    
    Status       int    `gorm:"type:tinyint;default:0;index" json:"status"`        // 0=失败, 1=成功, 2=过滤跳过
    ErrorMessage string `gorm:"type:varchar(500);default:''" json:"error_message"` // 错误信息
    
    // 消息快照（用于调试）
    MessageTitle   string `gorm:"type:varchar(250);default:''" json:"message_title"`   // 消息标题/摘要
    MessageContent string `gorm:"type:text" json:"message_content"`                    // 消息内容摘要
    
    CreateTime   int64  `gorm:"type:bigint;default:0;index" json:"create_time"`     // 同步时间
}

func (TelegramSyncLog) TableName() string {
    return "telegram_sync_log"
}
```

**状态说明**:

| Status | 说明 |
|--------|------|
| 0 | 同步失败 |
| 1 | 同步成功 |
| 2 | 被过滤规则跳过 |

### 3. TelegramSession（会话信息）

存储 Telegram 认证会话（加密）。

```go
type TelegramSession struct {
    ID           int    `gorm:"type:int;size:32;primaryKey;autoIncrement" json:"id"`
    SessionData  []byte `gorm:"type:blob" json:"-"`                    // 加密的会话数据
    SessionHash  string `gorm:"type:varchar(64);uniqueIndex" json:"session_hash"` // 会话哈希（用于验证）
    
    Status       int    `gorm:"type:tinyint;default:1" json:"status"`  // 1=有效, 0=过期
    CreateTime   int64  `gorm:"type:bigint;default:0" json:"create_time"`
    UpdateTime   int64  `gorm:"type:bigint;default:0" json:"update_time"`
}

func (TelegramSession) TableName() string {
    return "telegram_session"
}
```

**安全说明**:
- `SessionData` 存储加密后的会话数据
- 使用 AES-256-GCM 加密
- 加密密钥从配置文件中获取

### 4. TelegramPluginConfig（插件配置）

存储插件全局配置。

```go
type TelegramPluginConfig struct {
    ID   int `gorm:"type:int;size:32;primaryKey;autoIncrement" json:"id"`
    
    // Telegram API 配置
    AppID       int    `gorm:"type:int;size:32" json:"app_id"`           // Telegram App ID
    AppHash     string `gorm:"type:varchar(64)" json:"app_hash"`         // Telegram App Hash
    PhoneNumber string `gorm:"type:varchar(20)" json:"phone_number"`     // 手机号（用户认证）
    
    // 会话配置
    SessionKey string `gorm:"type:varchar(64)" json:"session_key"` // 会话加密密钥
    
    // 同步配置
    AutoReconnect   bool `gorm:"type:boolean;default:true" json:"auto_reconnect"`   // 自动重连
    ReconnectDelay  int  `gorm:"type:int;default:5" json:"reconnect_delay"`         // 重连延迟（秒）
    SyncDelay       int  `gorm:"type:int;default:1" json:"sync_delay"`              // 同步延迟（秒）
    MaxRetries      int  `gorm:"type:int;default:3" json:"max_retries"`             // 最大重试次数
    
    // 媒体配置
    DownloadMedia   bool `gorm:"type:boolean;default:true" json:"download_media"`   // 下载媒体
    MaxImageSize    int  `gorm:"type:int;default:10485760" json:"max_image_size"`   // 最大图片大小（字节）
    
    // 日志配置
    LogLevel       string `gorm:"type:varchar(10);default:'info'" json:"log_level"` // 日志级别
    KeepLogDays    int    `gorm:"type:int;default:30" json:"keep_log_days"`         // 日志保留天数
    
    CreateTime int64 `gorm:"type:bigint;default:0" json:"create_time"`
    UpdateTime int64 `gorm:"type:bigint;default:0" json:"update_time"`
}

func (TelegramPluginConfig) TableName() string {
    return "telegram_plugin_config"
}
```

## 实体关系

```
TelegramChannel (1) -----> (N) TelegramSyncLog
     |
     +---> Category (N:1)  [关联到现有分类表]
     |
     +---> Article (N:1)   [关联到现有文章表]

TelegramPluginConfig (1) -----> (1) TelegramSession
```

## 索引设计

### telegram_channel
- `idx_channel_id`: (channel_id) - 唯一索引，快速查找频道
- `idx_status`: (status) - 按状态过滤
- `idx_category`: (category_id) - 按分类查询

### telegram_sync_log
- `idx_channel_time`: (channel_id, create_time) - 按频道查询日志
- `idx_message`: (channel_id, message_id) - 消息去重检查
- `idx_status`: (status) - 按状态统计

### telegram_session
- `idx_hash`: (session_hash) - 唯一索引，会话验证

## 数据迁移

GORM 自动迁移将在插件初始化时执行：

```go
func autoMigrate(db *gorm.DB) error {
    return db.AutoMigrate(
        &TelegramChannel{},
        &TelegramSyncLog{},
        &TelegramSession{},
        &TelegramPluginConfig{},
    )
}
```

## 过滤规则 JSON 格式

`TelegramChannel.FilterKeywords` 字段存储的 JSON 格式：

```json
{
  "type": "whitelist",
  "keywords": ["技术", "教程", "编程"],
  "case_sensitive": false,
  "match_all": false
}
```

或黑名单模式：

```json
{
  "type": "blacklist",
  "keywords": ["广告", "推广", "赞助"],
  "case_sensitive": false
}
```

**字段说明**:
- `type`: "whitelist" 或 "blacklist"
- `keywords`: 关键词数组
- `case_sensitive`: 是否区分大小写
- `match_all`: 白名单模式下是否需要匹配所有关键词
