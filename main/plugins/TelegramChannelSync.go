package plugins

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/gotd/td/tg"
	"go.uber.org/zap"

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
	ctx           *pluginEntity.Plugin
	client        *telegram_sync.Client
	channels      []telegram_sync.ChannelConfig
	syncLogs      []telegram_sync.SyncLogItem
	filter        *telegram_sync.FilterEngine
	mediaHandler  *telegram_sync.MediaHandler
	mu            sync.RWMutex
	cancel        context.CancelFunc
	pluginCtx     context.Context
	wg            sync.WaitGroup
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

	// 调试输出
	fmt.Println("=== TelegramChannelSync 插件开始加载 ===")
	fmt.Printf("AppID: %d, AppHash: %s, SessionKey: %s\n", p.AppID, p.AppHash, p.SessionKey)

	// 创建插件级别的 context
	p.pluginCtx, p.cancel = context.WithCancel(context.Background())

	// 初始化过滤引擎
	p.filter = telegram_sync.NewFilterEngine(ctx.Log)
	p.mediaHandler = telegram_sync.NewMediaHandler(ctx.Log)
	p.mediaHandler.SetMaxImageSize(p.MaxImageSize)

	ctx.Log.Info("Telegram 频道同步插件开始加载",
		zap.Int("app_id", p.AppID),
		zap.Bool("has_app_hash", p.AppHash != ""),
		zap.Bool("has_session_key", p.SessionKey != ""))

	// 解析频道配置
	if err := p.parseChannels(); err != nil {
		ctx.Log.Warn("解析频道配置失败", zap.Error(err))
	}
	ctx.Log.Info("频道配置解析完成", zap.Int("channels", len(p.channels)))

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

	// 启动日志清理定时任务
	p.startLogCleanupTask()

	ctx.Log.Info("Telegram 频道同步插件加载完成",
		zap.Int("channels", len(p.channels)),
		zap.Bool("auto_reconnect", p.AutoReconnect),
		zap.Bool("connected", p.Connected),
		zap.Bool("authenticated", p.Authenticated))

	return nil
}

// startLogCleanupTask 启动日志清理定时任务
func (p *TelegramChannelSync) startLogCleanupTask() {
	p.wg.Add(1)
	go func() {
		defer p.wg.Done()

		// 每天凌晨 3 点清理日志
		ticker := time.NewTicker(24 * time.Hour)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				p.ctx.Log.Info("开始清理过期同步日志")
				if err := p.CleanupOldLogs(); err != nil {
					p.ctx.Log.Error("清理日志失败", zap.Error(err))
				} else {
					p.ctx.Log.Info("过期日志清理完成")
				}
			case <-p.pluginCtx.Done():
				p.ctx.Log.Info("日志清理任务已停止")
				return
			}
		}
	}()
}

// Run 手动执行（用于测试或触发同步）
func (p *TelegramChannelSync) Run(ctx *pluginEntity.Plugin) error {
	ctx.Log.Info("手动触发 Telegram 同步检查")

	// 重新解析频道配置（配置可能已更新）
	if err := p.parseChannels(); err != nil {
		ctx.Log.Warn("解析频道配置失败", zap.Error(err))
	}

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
		&telegram_sync.TelegramMedia{},
	)
}

// initClient 初始化 Telegram 客户端
func (p *TelegramChannelSync) initClient() error {
	fmt.Println("=== initClient 开始 ===")
	p.ctx.Log.Info("开始初始化 Telegram 客户端",
		zap.Int("app_id", p.AppID),
		zap.Bool("has_session_key", p.SessionKey != ""))

	config := &telegram_sync.ClientConfig{
		AppID:          p.AppID,
		AppHash:        p.AppHash,
		PhoneNumber:    p.PhoneNumber,
		SessionKey:     p.SessionKey,
		AutoReconnect:  p.AutoReconnect,
		ReconnectDelay: p.ReconnectDelay,
	}

	// 创建 session storage（如果配置了密钥）
	var storage *telegram_sync.DBStorage
	if p.SessionKey != "" && db.DB != nil {
		storage = telegram_sync.NewDBStorage(db.DB, p.SessionKey)
		p.ctx.Log.Info("已创建 session storage，将自动恢复会话")
	}

	// 创建客户端（带 session storage）
	var client *telegram_sync.Client
	var err error
	if storage != nil {
		client, err = telegram_sync.NewClientWithStorage(config, storage, p.ctx.Log)
	} else {
		client, err = telegram_sync.NewClient(config, p.ctx.Log)
	}
	if err != nil {
		return err
	}

	p.client = client

	// 设置消息处理回调
	p.ctx.Log.Info("设置消息处理回调")
	p.client.SetMessageHandler(p.handleChannelMessage)

	// 如果有 session storage，启动客户端时会自动恢复会话
	if storage != nil {
		p.ctx.Log.Info("启动 Telegram 客户端...")
		if err := p.client.Start(context.Background()); err != nil {
			p.ctx.Log.Warn("客户端启动失败", zap.Error(err))
		} else {
			// 检查是否已认证（会话恢复成功）
			if p.client.IsAuthenticated() {
				p.Authenticated = true
				p.Connected = true
				p.ctx.Log.Info("Telegram 会话自动恢复成功")
			} else {
				p.ctx.Log.Info("Telegram 客户端已连接，等待认证")
			}
		}
	} else {
		p.ctx.Log.Info("没有 session storage，客户端将在认证时启动")
	}

	return nil
}

// handleChannelMessage 处理频道消息
func (p *TelegramChannelSync) handleChannelMessage(ctx context.Context, channelID int64, msg *tg.Message) {
	fmt.Printf("=== handleChannelMessage 被调用: channelID=%d, messageID=%d ===\n", channelID, msg.ID)
	p.ctx.Log.Info("handleChannelMessage 被调用",
		zap.Int64("channel_id", channelID),
		zap.Int("message_id", msg.ID),
		zap.String("text", truncateMsgText(msg.Message, 50)))

	// 重新解析频道配置（确保最新）
	p.mu.Lock()
	if err := p.parseChannels(); err != nil {
		p.ctx.Log.Warn("解析频道配置失败", zap.Error(err))
		fmt.Printf("=== 解析频道配置失败: %v ===\n", err)
	}
	// 直接在持有锁的情况下获取启用的频道（避免死锁）
	var enabledChannels []telegram_sync.ChannelConfig
	for _, ch := range p.channels {
		if ch.Status == 1 {
			enabledChannels = append(enabledChannels, ch)
		}
	}
	p.channels = enabledChannels
	p.ctx.Log.Info("当前配置的频道列表", zap.Int("count", len(p.channels)))
	fmt.Printf("=== 配置的频道数量: %d ===\n", len(p.channels))
	for _, ch := range p.channels {
		p.ctx.Log.Info("频道配置", zap.Int64("channel_id", ch.ChannelID), zap.String("name", ch.ChannelName), zap.Int("status", ch.Status))
		fmt.Printf("=== 频道配置: channelID=%d, name=%s, status=%d ===\n", ch.ChannelID, ch.ChannelName, ch.Status)
	}
	p.mu.Unlock()

	// 1. 获取频道配置
	channelConfig := p.GetChannelByID(channelID)
	if channelConfig == nil {
		p.ctx.Log.Warn("频道未配置监听，跳过处理", zap.Int64("channel_id", channelID))
		fmt.Printf("=== 频道 %d 未配置监听，跳过 ===\n", channelID)
		return
	}

	p.ctx.Log.Info("找到频道配置", zap.Int64("channel_id", channelID), zap.String("name", channelConfig.ChannelName), zap.Int("status", channelConfig.Status))
	fmt.Printf("=== 找到频道配置: channelID=%d, name=%s, status=%d ===\n", channelID, channelConfig.ChannelName, channelConfig.Status)

	// 2. 检查频道是否启用
	if channelConfig.Status != 1 {
		p.ctx.Log.Warn("频道已禁用，跳过处理", zap.Int64("channel_id", channelID))
		return
	}

	// 3. 检查消息是否已处理（去重）
	if p.CheckMessageDuplicate(channelID, int64(msg.ID)) {
		p.ctx.Log.Warn("消息已处理，跳过", zap.Int("message_id", msg.ID))
		return
	}

	// 4. 应用过滤规则
	if p.filter != nil {
		if !p.filter.Apply(msg, channelConfig) {
			p.ctx.Log.Info("消息被过滤规则跳过",
				zap.Int64("channel_id", channelID),
				zap.Int("message_id", msg.ID),
				zap.String("filter_keywords", channelConfig.FilterKeywords),
				zap.String("filter_types", channelConfig.FilterMessageTypes))
			// 记录跳过日志（status=2 表示跳过）
			p.RecordSyncLog(channelID, int64(msg.ID), 0, 2, "filtered", truncateMsgText(msg.Message, 100))
			return
		}
	}

	p.ctx.Log.Info("开始创建文章", zap.Int64("channel_id", channelID), zap.Int("message_id", msg.ID))

	// 4. 处理媒体
	var mediaInfos []telegram_sync.MessageMediaInfo
	var thumbnailURL string
	if p.mediaHandler != nil && msg.Media != nil && p.DownloadMedia {
		p.mediaHandler.SetAPI(p.client.GetAPI())
		var err error
		mediaInfos, err = p.mediaHandler.ProcessMedia(ctx, msg, channelID, int64(msg.ID))
		if err != nil {
			p.ctx.Log.Warn("处理媒体失败", zap.Error(err))
		}
		// 使用第一个图片作为缩略图
		if len(mediaInfos) > 0 {
			thumbnailURL = mediaInfos[0].URL
			p.ctx.Log.Info("媒体处理成功", zap.Int("count", len(mediaInfos)), zap.String("thumbnail", thumbnailURL))
		}
	}

	// 5. 提取标题和内容
	title := extractTitleFromMessage(msg.Message)
	content := msg.Message

	// 6. 创建文章（包含媒体信息）
	article, err := p.CreateArticle(&telegram_sync.TelegramChannel{
		ChannelID:     channelID,
		CategoryID:    channelConfig.CategoryID,
		ArticleStatus: true, // 发布状态
	}, title, content, thumbnailURL)

	if err != nil {
		p.ctx.Log.Error("创建文章失败", zap.Error(err))
		p.RecordSyncLog(channelID, int64(msg.ID), 0, 0, err.Error(), title)
		return
	}

	// 7. 记录成功日志（包含媒体信息）
	p.RecordSyncLog(channelID, int64(msg.ID), article.ID, 1, "", title)

	p.ctx.Log.Info("文章创建成功",
		zap.Int64("channel_id", channelID),
		zap.Int("message_id", msg.ID),
		zap.Int("article_id", article.ID))
}

// extractTitleFromMessage 从消息文本提取标题
func extractTitleFromMessage(text string) string {
	if len(text) == 0 {
		return "无标题消息"
	}

	// 使用第一行作为标题
	lines := strings.Split(text, "\n")
	title := strings.TrimSpace(lines[0])

	// 限制标题长度
	if len(title) > 100 {
		title = title[:100] + "..."
	}

	return title
}

// truncateMsgText 截断消息文本
func truncateMsgText(text string, maxLen int) string {
	if len(text) <= maxLen {
		return text
	}
	return text[:maxLen] + "..."
}

// GetStatus 获取插件状态
func (p *TelegramChannelSync) GetStatus() map[string]interface{} {
	p.mu.RLock()
	defer p.mu.RUnlock()

	status := map[string]interface{}{
		"connected":           p.Connected,
		"authenticated":       p.Authenticated,
		"monitored_channels":  len(p.channels),
		"active_channels":     p.countActiveChannels(),
		"auto_reconnect":      p.AutoReconnect,
		"download_media":      p.DownloadMedia,
		"has_client":          p.client != nil,
		"has_filter":          p.filter != nil,
		"channels_json_len":   len(p.ChannelsJSON),
	}

	// 添加重连状态和消息处理状态
	if p.client != nil {
		reconnectStatus := p.client.GetReconnectStatus()
		status["reconnect_attempts"] = reconnectStatus["reconnect_attempts"]
		status["max_reconnect"] = reconnectStatus["max_reconnect"]
		status["reconnect_delay"] = reconnectStatus["reconnect_delay"]
		status["has_message_handler"] = p.client.HasMessageHandler()
		status["client_connected"] = p.client.IsConnected()
		status["client_authenticated"] = p.client.IsAuthenticated()
	}

	return status
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
func (p *TelegramChannelSync) GetRecentLogs(limit int) (interface{}, error) {
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

// SendAuthCode 发送 Telegram 验证码
func (p *TelegramChannelSync) SendAuthCode(phoneNumber string) (string, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	// 检查配置
	if p.AppID <= 0 || p.AppHash == "" {
		return "", fmt.Errorf("Telegram API 配置不完整，请先配置 App ID 和 App Hash")
	}

	// 更新手机号
	p.PhoneNumber = phoneNumber

	// 如果客户端未初始化，先初始化
	if p.client == nil {
		if err := p.initClient(); err != nil {
			return "", fmt.Errorf("初始化客户端失败: %w", err)
		}
	}

	// 启动客户端（如果未连接）
	if !p.client.IsConnected() {
		p.ctx.Log.Info("启动 Telegram 客户端...")
		if err := p.client.Start(context.Background()); err != nil {
			return "", fmt.Errorf("启动客户端失败: %w", err)
		}
		// 等待连接建立
		time.Sleep(3 * time.Second)
	}

	// 检查是否已连接
	if !p.client.IsConnected() {
		return "", fmt.Errorf("客户端连接超时，请重试")
	}

	// 发送验证码
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	phoneCodeHash, err := p.client.SendAuthCode(ctx, phoneNumber)
	if err != nil {
		p.ctx.Log.Error("发送验证码失败", zap.Error(err))
		return "", fmt.Errorf("发送验证码失败: %w", err)
	}

	p.ctx.Log.Info("验证码已发送", zap.String("phone", phoneNumber))
	return phoneCodeHash, nil
}

// VerifyAuthCode 验证 Telegram 验证码
func (p *TelegramChannelSync) VerifyAuthCode(code string) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.client == nil {
		return fmt.Errorf("客户端未初始化")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err := p.client.VerifyAuthCode(ctx, code)
	if err != nil {
		p.ctx.Log.Error("验证码验证失败", zap.Error(err))
		return fmt.Errorf("验证失败: %w", err)
	}

	// 更新认证状态
	p.Authenticated = true
	p.Connected = true

	// gotd/td 会自动保存会话到 SessionStorage
	p.ctx.Log.Info("Telegram 认证成功，会话已自动保存")

	p.ctx.Log.Info("Telegram 认证成功")
	return nil
}

// GetAuthStatus 获取认证状态
func (p *TelegramChannelSync) GetAuthStatus() map[string]interface{} {
	p.mu.RLock()
	defer p.mu.RUnlock()

	status := map[string]interface{}{
		"connected":           p.Connected,
		"authenticated":       p.Authenticated,
		"phone_number":        p.PhoneNumber,
		"monitored_channels":  len(p.channels),
		"active_channels":     p.countActiveChannels(),
		"has_client":          p.client != nil,
	}

	if p.client != nil {
		status["client_connected"] = p.client.IsConnected()
		status["client_authenticated"] = p.client.IsAuthenticated()
	}

	// 检查会话是否存在
	if db.DB != nil {
		var sessionCount int64
		db.DB.Model(&telegram_sync.TelegramSession{}).Count(&sessionCount)
		status["has_session_data"] = sessionCount > 0
	}

	return status
}

// CheckSessionStatus 检查会话状态（用于前端显示详细错误）
func (p *TelegramChannelSync) CheckSessionStatus() map[string]interface{} {
	p.mu.RLock()
	defer p.mu.RUnlock()

	result := map[string]interface{}{
		"status": "unknown",
		"message": "",
		"need_reauth": false,
	}

	// 检查是否有会话数据
	if db.DB != nil {
		var sessionCount int64
		db.DB.Model(&telegram_sync.TelegramSession{}).Count(&sessionCount)
		if sessionCount == 0 {
			result["status"] = "no_session"
			result["message"] = "没有保存的会话数据，请先完成认证"
			result["need_reauth"] = true
			return result
		}
	}

	// 检查客户端状态
	if p.client == nil {
		result["status"] = "client_not_initialized"
		result["message"] = "客户端未初始化"
		result["need_reauth"] = true
		return result
	}

	if !p.client.IsConnected() {
		result["status"] = "not_connected"
		result["message"] = "客户端未连接"
		result["need_reauth"] = false
		return result
	}

	// 尝试调用 API 检查认证状态
	if p.client.GetAPI() != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		_, err := p.client.GetAPI().UsersGetFullUser(ctx, &tg.InputUserSelf{})
		if err != nil {
			result["status"] = "session_expired"
			// 检查是否是 AUTH_KEY_UNREGISTERED 错误
			if strings.Contains(err.Error(), "AUTH_KEY_UNREGISTERED") ||
			   strings.Contains(err.Error(), "SESSION_REVOKED") ||
			   strings.Contains(err.Error(), "AUTH_KEY_INVALID") {
				result["message"] = "会话已过期或被撤销，请重新认证"
				result["need_reauth"] = true
			} else {
				result["message"] = fmt.Sprintf("认证检查失败: %s", err.Error())
				result["need_reauth"] = true
			}
			return result
		}
	}

	result["status"] = "ok"
	result["message"] = "会话正常"
	result["need_reauth"] = false
	return result
}

// GetUserChannels 获取用户的 Telegram 频道列表
func (p *TelegramChannelSync) GetUserChannels() ([]map[string]interface{}, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	// 检查客户端是否已连接
	if p.client == nil || !p.client.IsConnected() {
		return nil, fmt.Errorf("客户端未连接，请先完成认证")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	channels, err := p.client.GetUserChannels(ctx)
	if err != nil {
		p.ctx.Log.Error("获取频道列表失败", zap.Error(err))
		// 如果是认证错误，清除会话状态
		if err.Error() != "" && (strings.Contains(err.Error(), "AUTH_KEY_UNREGISTERED") ||
			strings.Contains(err.Error(), "SESSION_REVOKED") ||
			strings.Contains(err.Error(), "AUTH_KEY_INVALID")) {
			p.ctx.Log.Warn("检测到会话无效，清除认证状态")
			p.Authenticated = false
			p.Connected = false
			p.clearSession()
		}
		return nil, err
	}

	return channels, nil
}

// clearSession 清除数据库中的会话数据
func (p *TelegramChannelSync) clearSession() {
	if db.DB != nil {
		db.DB.Where("1 = 1").Delete(&telegram_sync.TelegramSession{})
		p.ctx.Log.Info("已清除数据库中的会话数据")
	}
}

// ClearAuth 清除认证状态（供 API 调用）
func (p *TelegramChannelSync) ClearAuth() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	// 停止客户端
	if p.client != nil {
		p.client.Stop()
	}

	// 清除会话数据
	p.clearSession()

	// 重置状态
	p.Authenticated = false
	p.Connected = false
	p.client = nil

	p.ctx.Log.Info("Telegram 认证已清除")
	return nil
}

// GetMediaURL 获取媒体文件 URL
func (p *TelegramChannelSync) GetMediaURL(mediaId int64) (string, error) {
	// 从数据库查询媒体信息
	var media telegram_sync.TelegramMedia
	if err := db.DB.Where("media_id = ?", mediaId).First(&media).Error; err != nil {
		return "", fmt.Errorf("media not found: %w", err)
	}
	
	if media.StorageURL == "" {
		return "", fmt.Errorf("media URL not available")
	}
	
	return media.StorageURL, nil
}
