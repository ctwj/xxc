package telegram_sync

import (
	"strings"

	"github.com/gotd/td/tg"
	"go.uber.org/zap"
)

// FilterEngine 过滤规则引擎
type FilterEngine struct {
	log *zap.Logger
}

// NewFilterEngine 创建过滤引擎
func NewFilterEngine(log *zap.Logger) *FilterEngine {
	return &FilterEngine{log: log}
}

// Apply 应用过滤规则
func (f *FilterEngine) Apply(msg *tg.Message, config *ChannelConfig) bool {
	// 1. 消息类型过滤
	if !f.checkMessageType(msg, config.FilterMessageTypes) {
		f.log.Debug("消息类型不符合", zap.String("types", config.FilterMessageTypes))
		return false
	}

	// 2. 长度过滤
	if !f.checkLength(msg, config.FilterMinLength, config.FilterMaxLength) {
		f.log.Debug("消息长度不符合",
			zap.Int("min", config.FilterMinLength),
			zap.Int("max", config.FilterMaxLength),
			zap.Int("actual", len(msg.Message)))
		return false
	}

	// 3. 关键词过滤
	if config.FilterKeywords != nil {
		if !f.checkKeywords(msg, config.FilterKeywords) {
			f.log.Debug("关键词过滤不符合")
			return false
		}
	}

	return true
}

// checkMessageType 检查消息类型
func (f *FilterEngine) checkMessageType(msg *tg.Message, allowedTypes string) bool {
	if allowedTypes == "" || allowedTypes == "text,photo,video" {
		return true // 默认允许所有类型
	}

	types := strings.Split(allowedTypes, ",")

	// 检查消息是否包含指定类型的内容
	hasText := len(msg.Message) > 0
	hasPhoto := f.hasPhotoMedia(msg)
	hasVideo := f.hasVideoMedia(msg)

	for _, t := range types {
		t = strings.TrimSpace(t)
		switch t {
		case "text":
			if hasText {
				return true
			}
		case "photo":
			if hasPhoto {
				return true
			}
		case "video":
			if hasVideo {
				return true
			}
		}
	}

	return false
}

// hasPhotoMedia 检查是否包含图片
func (f *FilterEngine) hasPhotoMedia(msg *tg.Message) bool {
	if msg.Media == nil {
		return false
	}

	switch media := msg.Media.(type) {
	case *tg.MessageMediaPhoto:
		return true
	case *tg.MessageMediaDocument:
		// 检查文档是否为图片
		if doc, ok := media.Document.(*tg.Document); ok {
			for _, attr := range doc.Attributes {
				if _, ok := attr.(*tg.DocumentAttributeImageSize); ok {
					return true
				}
			}
		}
	}

	return false
}

// hasVideoMedia 检查是否包含视频
func (f *FilterEngine) hasVideoMedia(msg *tg.Message) bool {
	if msg.Media == nil {
		return false
	}

	switch media := msg.Media.(type) {
	case *tg.MessageMediaDocument:
		if doc, ok := media.Document.(*tg.Document); ok {
			for _, attr := range doc.Attributes {
				if _, ok := attr.(*tg.DocumentAttributeVideo); ok {
					return true
				}
			}
		}
	}

	return false
}

// checkLength 检查消息长度
func (f *FilterEngine) checkLength(msg *tg.Message, minLen, maxLen int) bool {
	textLen := len(msg.Message)

	if minLen > 0 && textLen < minLen {
		return false
	}

	if maxLen > 0 && textLen > maxLen {
		// 超长消息可以截断，不算过滤失败
		return true
	}

	return true
}

// checkKeywords 检查关键词
func (f *FilterEngine) checkKeywords(msg *tg.Message, rule *FilterRule) bool {
	if len(rule.Keywords) == 0 {
		return true
	}

	text := msg.Message
	if !rule.CaseSensitive {
		text = strings.ToLower(text)
	}

	switch rule.Type {
	case "whitelist":
		// 白名单：必须包含至少一个关键词（或全部，取决于 MatchAll）
		return f.matchWhitelist(text, rule.Keywords, rule.CaseSensitive, rule.MatchAll)

	case "blacklist":
		// 黑名单：不能包含任何关键词
		return f.matchBlacklist(text, rule.Keywords, rule.CaseSensitive)

	default:
		return true
	}
}

// matchWhitelist 白名单匹配
func (f *FilterEngine) matchWhitelist(text string, keywords []string, caseSensitive bool, matchAll bool) bool {
	if !caseSensitive {
		text = strings.ToLower(text)
	}

	matchedCount := 0
	for _, kw := range keywords {
		if !caseSensitive {
			kw = strings.ToLower(kw)
		}
		if strings.Contains(text, kw) {
			matchedCount++
		}
	}

	if matchAll {
		return matchedCount == len(keywords)
	}
	return matchedCount > 0
}

// matchBlacklist 黑名单匹配
func (f *FilterEngine) matchBlacklist(text string, keywords []string, caseSensitive bool) bool {
	if !caseSensitive {
		text = strings.ToLower(text)
	}

	for _, kw := range keywords {
		if !caseSensitive {
			kw = strings.ToLower(kw)
		}
		if strings.Contains(text, kw) {
			return false // 包含黑名单关键词，过滤掉
		}
	}

	return true // 不包含任何黑名单关键词，通过
}

// KeywordWhitelistFilter 关键词白名单过滤器
func (f *FilterEngine) KeywordWhitelistFilter(text string, keywords []string, caseSensitive bool) bool {
	return f.matchWhitelist(text, keywords, caseSensitive, false)
}

// KeywordBlacklistFilter 关键词黑名单过滤器
func (f *FilterEngine) KeywordBlacklistFilter(text string, keywords []string, caseSensitive bool) bool {
	return f.matchBlacklist(text, keywords, caseSensitive)
}

// MessageTypeFilter 消息类型过滤器
func (f *FilterEngine) MessageTypeFilter(msg *tg.Message, allowedTypes []string) bool {
	for _, t := range allowedTypes {
		switch t {
		case "text":
			if len(msg.Message) > 0 {
				return true
			}
		case "photo":
			if f.hasPhotoMedia(msg) {
				return true
			}
		case "video":
			if f.hasVideoMedia(msg) {
				return true
			}
		}
	}
	return false
}

// LengthFilter 长度过滤器
func (f *FilterEngine) LengthFilter(text string, minLen, maxLen int) bool {
	textLen := len(text)
	if minLen > 0 && textLen < minLen {
		return false
	}
	if maxLen > 0 && textLen > maxLen {
		return false
	}
	return true
}

// TruncateContent 截断内容
func (f *FilterEngine) TruncateContent(text string, maxLen int) string {
	if len(text) <= maxLen {
		return text
	}
	return text[:maxLen] + "..."
}