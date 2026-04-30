package telegram_sync

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/gotd/td/tg"
	"go.uber.org/zap"
)

// MessageHandler 消息处理器
type MessageHandler struct {
	client    *Client
	plugin    interface{} // 指向主插件的引用
	log       *zap.Logger
	filter    *FilterEngine
	media     *MediaHandler
}

// NewMessageHandler 创建消息处理器
func NewMessageHandler(client *Client, log *zap.Logger) *MessageHandler {
	return &MessageHandler{
		client: client,
		log:    log,
		filter: NewFilterEngine(log),
		media:  NewMediaHandler(log),
	}
}

// SetPlugin 设置插件引用
func (h *MessageHandler) SetPlugin(plugin interface{}) {
	h.plugin = plugin
}

// HandleMessage 处理频道消息
func (h *MessageHandler) HandleMessage(ctx context.Context, channelID int64, msg *tg.Message) error {
	h.log.Info("处理频道消息",
		zap.Int64("channel_id", channelID),
		zap.Int("message_id", msg.ID))

	// 1. 获取频道配置
	channelConfig := h.getChannelConfig(channelID)
	if channelConfig == nil {
		h.log.Debug("频道未配置监听", zap.Int64("channel_id", channelID))
		return nil
	}

	// 2. 检查频道是否启用
	if channelConfig.Status != 1 {
		h.log.Debug("频道已禁用", zap.Int64("channel_id", channelID))
		return nil
	}

	// 3. 检查消息是否已处理（去重）
	if h.isDuplicate(channelID, msg.ID) {
		h.log.Debug("消息已处理，跳过", zap.Int("message_id", msg.ID))
		return nil
	}

	// 4. 应用过滤规则
	if !h.filter.Apply(msg, channelConfig) {
		h.log.Debug("消息被过滤规则跳过",
			zap.Int("message_id", msg.ID),
			zap.String("text", truncateText(msg.Message, 50)))
		// 记录跳过日志
		h.recordLog(channelID, msg.ID, 0, 2, "", truncateText(msg.Message, 100))
		return nil
	}

	// 5. 处理媒体
	thumbnail := ""
	if h.media != nil && msg.Media != nil {
		mediaURL, err := h.media.ProcessMedia(ctx, msg.Media)
		if err != nil {
			h.log.Warn("处理媒体失败", zap.Error(err))
		} else {
			thumbnail = mediaURL
		}
	}

	// 6. 转换为文章
	title := h.extractTitle(msg)
	content := h.formatContent(msg)

	// 7. 创建文章
	articleID, err := h.createArticle(channelConfig, title, content, thumbnail)
	if err != nil {
		h.log.Error("创建文章失败", zap.Error(err))
		h.recordLog(channelID, msg.ID, 0, 0, err.Error(), title)
		return err
	}

	// 8. 记录成功日志
	h.recordLog(channelID, msg.ID, articleID, 1, "", title)

	h.log.Info("文章创建成功",
		zap.Int64("channel_id", channelID),
		zap.Int("message_id", msg.ID),
		zap.Int("article_id", articleID))

	return nil
}

// getChannelConfig 获取频道配置
func (h *MessageHandler) getChannelConfig(channelID int64) *ChannelConfig {
	// 通过插件接口获取
	// 实际实现需要主插件提供 GetChannelByID 方法
	return nil
}

// isDuplicate 检查消息是否重复
func (h *MessageHandler) isDuplicate(channelID int64, messageID int) bool {
	// 通过插件接口检查
	// 实际实现需要主插件提供 CheckMessageDuplicate 方法
	return false
}

// extractTitle 提取标题
func (h *MessageHandler) extractTitle(msg *tg.Message) string {
	// 尝试从消息中提取标题
	text := msg.Message
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

// formatContent 格式化内容
func (h *MessageHandler) formatContent(msg *tg.Message) string {
	// 将消息内容转换为 HTML 格式
	content := msg.Message

	// 处理实体（链接、粗体等）
	if len(msg.Entities) > 0 {
		content = h.formatEntities(content, msg.Entities)
	}

	// 添加媒体信息
	if msg.Media != nil {
		content += "\n\n<!-- 媒体内容 -->"
	}

	return content
}

// formatEntities 格式化消息实体
func (h *MessageHandler) formatEntities(text string, entities []tg.MessageEntityClass) string {
	// 简化实现：保留原始文本
	// 实际实现可以处理粗体、链接、斜体等
	return text
}

// truncateText 截断文本
func truncateText(text string, maxLen int) string {
	if len(text) <= maxLen {
		return text
	}
	return text[:maxLen] + "..."
}

// createArticle 创建文章
func (h *MessageHandler) createArticle(config *ChannelConfig, title, content, thumbnail string) (int, error) {
	// 通过插件接口创建文章
	// 实际实现需要主插件提供 CreateArticle 方法
	return 0, nil
}

// recordLog 记录同步日志
func (h *MessageHandler) recordLog(channelID int64, messageID int, articleID int, status int, errMsg string, title string) {
	// 通过插件接口记录日志
	// 实际实现需要主插件提供 RecordSyncLog 方法
}

// SetupDispatcher 设置更新分发器
func (h *MessageHandler) SetupDispatcher(dispatcher *tg.UpdateDispatcher) {
	dispatcher.OnNewMessage(func(ctx context.Context, entities tg.Entities, u *tg.UpdateNewMessage) error {
		msg, ok := u.Message.(*tg.Message)
		if !ok || msg.Out {
			return nil
		}

		// 检查是否为频道消息
		peer := msg.GetPeerID()
		channelPeer, ok := peer.(*tg.PeerChannel)
		if !ok {
			return nil
		}

		channelID := channelPeer.GetChannelID()
		return h.HandleMessage(ctx, channelID, msg)
	})
}

// FilterConfigJSON 解析过滤配置 JSON
func ParseFilterConfig(jsonStr string) (*FilterRule, error) {
	if jsonStr == "" {
		return nil, nil
	}

	var rule FilterRule
	if err := json.Unmarshal([]byte(jsonStr), &rule); err != nil {
		return nil, err
	}

	return &rule, nil
}