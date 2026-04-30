package plugins

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"moss/domain/core/entity"
	"moss/domain/core/repository"
	"moss/domain/core/service"
	pluginEntity "moss/domain/support/entity"
	"moss/infrastructure/persistent/db"
	"moss/plugins/telegram_sync"
)

// TelegramChannelSync Telegram 频道同步插件
type TelegramChannelSync struct {
	// Telegram API 配置
	AppID       int    `json:"app_id"`
	AppHash     string `json:"app_hash"`
	PhoneNumber string `json:"phone_number"`

	// 会话配置
	SessionKey string `json:"session_key"`

	// 同步配置
	AutoReconnect  bool `json:"auto_reconnect"`
	ReconnectDelay int  `json:"reconnect_delay"`
	SyncDelay      int  `json:"sync_delay"`
	MaxRetries     int  `json:"max_retries"`

	// 媒体配置
	DownloadMedia bool `json:"download_media"`
	MaxImageSize  int  `json:"max_image_size"`

	// 日志配置
	LogLevel    string `json:"log_level"`
	KeepLogDays int    `json:"keep_log_days"`

	// 频道配置（JSON 存储）
	ChannelsJSON string `json:"channels_json"`

	// 运行时状态
	Connected    bool `json:"connected"`
	Authenticated bool `json:"authenticated"`

	// 内部字段
	ctx          *pluginEntity.Plugin
	client       *telegram_sync.Client
	channels     []telegram_sync.ChannelConfig
	syncLogs     []telegram_sync.SyncLogItem
	mu           sync.RWMutex
	cancel       context.CancelFunc
	wg           sync.WaitGroup
}

// NewTelegramChannelSync 创建插件实例
func NewTelegramChannelSync() *TelegramChannelSync {
	return &TelegramChannelSync{
		AutoReconnect:  true,
		ReconnectDelay: 5,
		SyncDelay:      1,
		MaxRetries:     3,
		DownloadMedia:  true,
		MaxImageSize:   10485760, // 10MB
		LogLevel:       "info",
		KeepLogDays:    30,
		ChannelsJSON:   "[]",
		Connected:      false,
		Authenticated:  false,
	}
}

// Info 返回插件信息
func (p *TelegramChannelSync) Info() *pluginEntity.PluginInfo {
	return &pluginEntity.PluginInfo{
		ID:         "TelegramChannelSync",
		About:      "Telegram 频道消息同步插件 - 自动监听频道消息并发布为 CMS 文章",
		RunEnable:  true,
		CronEnable: false, // 持续运行，不需要定时任务
	}
}

// Load 插件加载
func (p *TelegramChannelSync) Load(ctx *pluginEntity.Plugin) error {
	p.ctx = ctx

	// 解析频道配置
	if err := p.parseChannels(); err != nil {
		ctx.Log.Warn("解析频道配置失败", zap.Error(err))
	}

	// 自动迁移数据库表
	if err := p.autoMigrate(); err != nil {
		return fmt.Errorf("数据库迁移失败: %w", err)
	}

	// 初始化 Telegram 客户端
	if p.AppID > 0 && p.AppHash != "" {
		if err := p.initClient(); err != nil {
			ctx.Log.Warn("初始化 Telegram 客户端失败", zap.Error(err))
		}
	} else {
		ctx.Log.Info("Telegram API 配置不完整，请先配置 App ID 和 App Hash")
	}

	ctx.Log.Info("Telegram 频道同步插件加载完成",
		zap.Int("channels", len(p.channels)),
		zap.Bool("auto_reconnect", p.AutoReconnect))

	return nil
}

// Run 手动执行（用于测试或触发同步）
func (p *TelegramChannelSync) Run(ctx *pluginEntity.Plugin) error {
	ctx.Log.Info("手动触发 Telegram 同步检查")

	// 检查客户端状态
	if p.client == nil {
		return fmt.Errorf("Telegram 客户端未初始化")
	}

	// 检查认证状态
	if !p.Authenticated {
		return fmt.Errorf("Telegram 未认证，请先完成认证")
	}

	// 检查频道配置
	if len(p.channels) == 0 {
		return fmt.Errorf("未配置监听频道")
	}

	ctx.Log.Info("同步检查完成", zap.Int("monitored_channels", len(p.channels)))
	return nil
}

// Unload 插件卸载
func (p *TelegramChannelSync) Unload() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	// 停止客户端
	if p.cancel != nil {
		p.cancel()
	}

	// 等待所有 goroutine 完成
	p.wg.Wait()

	// 关闭客户端
	if p.client != nil {
		p.client.Stop()
	}

	if p.ctx != nil {
		p.ctx.Log.Info("Telegram 频道同步插件已卸载")
	}

	return nil
}

// parseChannels 解析频道配置
func (p *TelegramChannelSync) parseChannels() error {
	if p.ChannelsJSON == "" || p.ChannelsJSON == "[]" {
		p.channels = []telegram_sync.ChannelConfig{}
		return nil
	}

	var channels []telegram_sync.ChannelConfig
	if err := json.Unmarshal([]byte(p.ChannelsJSON), &channels); err != nil {
		return err
	}

	p.channels = channels
	return nil
}

// SaveChannels 保存频道配置
func (p *TelegramChannelSync) SaveChannels(channels []telegram_sync.ChannelConfig) error {
	data, err := json.Marshal(channels)
	if err != nil {
		return err
	}

	p.mu.Lock()
	p.ChannelsJSON = string(data)
	p.channels = channels
	p.mu.Unlock()

	return nil
}

// autoMigrate 自动迁移数据库表
func (p *TelegramChannelSync) autoMigrate() error {
	db := db.DB
	if db == nil {
		return fmt.Errorf("数据库连接不可用")
	}

	return db.AutoMigrate(
		&telegram_sync.TelegramChannel{},
		&telegram_sync.TelegramSyncLog{},
		&telegram_sync.TelegramSession{},
	)
}

// initClient 初始化 Telegram 客户端
func (p *TelegramChannelSync) initClient() error {
	config := &telegram_sync.ClientConfig{
		AppID:          p.AppID,
		AppHash:        p.AppHash,
		PhoneNumber:    p.PhoneNumber,
		SessionKey:     p.SessionKey,
		AutoReconnect:  p.AutoReconnect,
		ReconnectDelay: p.ReconnectDelay,
	}

	client, err := telegram_sync.NewClient(config, p.ctx.Log)
	if err != nil {
		return err
	}

	p.client = client

	// 尝试恢复会话
	if p.SessionKey != "" {
		session, err := p.loadSession()
		if err == nil && session != nil {
			p.ctx.Log.Info("发现已保存的会话，尝试恢复")
			if err := p.client.RestoreSession(session); err != nil {
				p.ctx.Log.Warn("会话恢复失败，需要重新认证", zap.Error(err))
			} else {
				p.Authenticated = true
				p.ctx.Log.Info("会话恢复成功")
			}
		}
	}

	return nil
}

// loadSession 从数据库加载会话
func (p *TelegramChannelSync) loadSession() (*telegram_sync.TelegramSession, error) {
	db := db.DB
	if db == nil {
		return nil, fmt.Errorf("数据库连接不可用")
	}

	var session telegram_sync.TelegramSession
	if err := db.Where("status = ?", 1).First(&session).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &session, nil
}

// saveSession 保存会话到数据库
func (p *TelegramChannelSync) saveSession(sessionData []byte) error {
	db := db.DB
	if db == nil {
		return fmt.Errorf("数据库连接不可用")
	}

	storage := telegram_sync.NewSessionStorage(p.SessionKey, "default")
	encrypted, err := storage.Encrypt(sessionData)
	if err != nil {
		return err
	}

	session := &telegram_sync.TelegramSession{
		SessionData: encrypted,
		SessionHash: storage.GenerateSessionHash(sessionData),
		Status:      1,
		CreateTime:  time.Now().Unix(),
		UpdateTime:  time.Now().Unix(),
	}

	// 删除旧会话
	db.Where("1 = 1").Delete(&telegram_sync.TelegramSession{})

	return db.Create(session).Error
}

// GetStatus 获取插件状态
func (p *TelegramChannelSync) GetStatus() map[string]interface{} {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return map[string]interface{}{
		"connected":           p.Connected,
		"authenticated":       p.Authenticated,
		"monitored_channels":  len(p.channels),
		"active_channels":     p.countActiveChannels(),
		"auto_reconnect":      p.AutoReconnect,
		"download_media":      p.DownloadMedia,
	}
}

// countActiveChannels 统计启用的频道数量
func (p *TelegramChannelSync) countActiveChannels() int {
	count := 0
	for _, ch := range p.channels {
		if ch.Status == 1 {
			count++
		}
	}
	return count
}

// CreateArticle 创建文章
func (p *TelegramChannelSync) CreateArticle(channel *telegram_sync.TelegramChannel, title, content string, thumbnail string) (*entity.Article, error) {
	article := &entity.Article{
		ArticleBase: entity.ArticleBase{
			Slug:        fmt.Sprintf("tg-%d-%d", channel.ChannelID, time.Now().Unix()),
			Title:       title,
			CreateTime:  time.Now().Unix(),
			CategoryID:  channel.CategoryID,
			Status:      channel.ArticleStatus,
			Thumbnail:   thumbnail,
			Description: truncateDescription(content, 200),
		},
		ArticleDetail: entity.ArticleDetail{
			Content: content,
		},
	}

	if err := service.Article.Create(article); err != nil {
		return nil, err
	}

	return article, nil
}

// RecordSyncLog 记录同步日志
func (p *TelegramChannelSync) RecordSyncLog(channelID int64, messageID int64, articleID int, status int, errMsg string, title string) error {
	db := db.DB
	if db == nil {
		return fmt.Errorf("数据库连接不可用")
	}

	log := &telegram_sync.TelegramSyncLog{
		ChannelID:     channelID,
		MessageID:     messageID,
		ArticleID:     articleID,
		Status:        status,
		ErrorMessage:  errMsg,
		MessageTitle:  title,
		CreateTime:    time.Now().Unix(),
	}

	return db.Create(log).Error
}

// CleanupOldLogs 清理过期日志
func (p *TelegramChannelSync) CleanupOldLogs() error {
	if p.KeepLogDays <= 0 {
		return nil
	}

	db := db.DB
	if db == nil {
		return fmt.Errorf("数据库连接不可用")
	}

	cutoff := time.Now().AddDate(0, 0, -p.KeepLogDays).Unix()
	return db.Where("create_time < ?", cutoff).Delete(&telegram_sync.TelegramSyncLog{}).Error
}

// GetRecentLogs 获取最近的同步日志
func (p *TelegramChannelSync) GetRecentLogs(limit int) ([]telegram_sync.TelegramSyncLog, error) {
	db := db.DB
	if db == nil {
		return nil, fmt.Errorf("数据库连接不可用")
	}

	var logs []telegram_sync.TelegramSyncLog
	if err := db.Order("create_time DESC").Limit(limit).Find(&logs).Error; err != nil {
		return nil, err
	}

	return logs, nil
}

// truncateDescription 截断描述
func truncateDescription(content string, maxLen int) string {
	if len(content) <= maxLen {
		return content
	}
	return content[:maxLen] + "..."
}

// GetChannelByID 根据频道 ID 获取频道配置
func (p *TelegramChannelSync) GetChannelByID(channelID int64) *telegram_sync.ChannelConfig {
	p.mu.RLock()
	defer p.mu.RUnlock()

	for i := range p.channels {
		if p.channels[i].ChannelID == channelID {
			return &p.channels[i]
		}
	}
	return nil
}

// GetEnabledChannels 获取所有启用的频道
func (p *TelegramChannelSync) GetEnabledChannels() []telegram_sync.ChannelConfig {
	p.mu.RLock()
	defer p.mu.RUnlock()

	var result []telegram_sync.ChannelConfig
	for _, ch := range p.channels {
		if ch.Status == 1 {
			result = append(result, ch)
		}
	}
	return result
}

// UpdateChannelSyncInfo 更新频道同步信息
func (p *TelegramChannelSync) UpdateChannelSyncInfo(channelID int64, lastMessageID int64) {
	p.mu.Lock()
	defer p.mu.Unlock()

	for i := range p.channels {
		if p.channels[i].ChannelID == channelID {
			// 更新数据库中的频道记录
			db := db.DB
			if db != nil {
				db.Model(&telegram_sync.TelegramChannel{}).
					Where("channel_id = ?", channelID).
					Updates(map[string]interface{}{
						"last_message_id": lastMessageID,
						"last_sync_time":  time.Now().Unix(),
					})
			}
			return
		}
	}
}

// CheckMessageDuplicate 检查消息是否已同步
func (p *TelegramChannelSync) CheckMessageDuplicate(channelID int64, messageID int64) bool {
	db := db.DB
	if db == nil {
		return false
	}

	var count int64
	db.Model(&telegram_sync.TelegramSyncLog{}).
		Where("channel_id = ? AND message_id = ?", channelID, messageID).
		Count(&count)

	return count > 0
}

// GetCategoryName 获取分类名称
func (p *TelegramChannelSync) GetCategoryName(categoryID int) string {
	if categoryID <= 0 {
		return ""
	}

	category, err := repository.Category.Get(categoryID)
	if err != nil {
		return ""
	}
	return category.Title
}
