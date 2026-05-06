package telegram_sync

import (
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// LogManager 日志管理器
type LogManager struct {
	log         *zap.Logger
	db          *gorm.DB
	keepLogDays int
}

// NewLogManager 创建日志管理器
func NewLogManager(db *gorm.DB, log *zap.Logger) *LogManager {
	return &LogManager{
		db:  db,
		log: log,
	}
}

// SetKeepLogDays 设置日志保留天数
func (l *LogManager) SetKeepLogDays(days int) {
	l.keepLogDays = days
}

// Record 记录同步日志
func (l *LogManager) Record(channelID int64, messageID int64, articleID int, status int, errMsg string, title string, content string) error {
	logEntry := &TelegramSyncLog{
		ChannelID:      channelID,
		MessageID:      messageID,
		ArticleID:      articleID,
		Status:         status,
		ErrorMessage:   errMsg,
		MessageTitle:   title,
		MessageContent: content,
		CreateTime:     time.Now().Unix(),
	}

	if l.db == nil {
		l.log.Warn("数据库不可用，无法记录日志")
		return nil
	}

	return l.db.Create(logEntry).Error
}

// RecordSuccess 记录成功日志
func (l *LogManager) RecordSuccess(channelID int64, messageID int64, articleID int, title string) error {
	return l.Record(channelID, messageID, articleID, 1, "", title, "")
}

// RecordFailure 记录失败日志
func (l *LogManager) RecordFailure(channelID int64, messageID int64, errMsg string, title string) error {
	return l.Record(channelID, messageID, 0, 0, errMsg, title, "")
}

// RecordSkipped 记录跳过日志
func (l *LogManager) RecordSkipped(channelID int64, messageID int64, reason string, title string) error {
	return l.Record(channelID, messageID, 0, 2, reason, title, "")
}

// GetRecent 获取最近的日志
func (l *LogManager) GetRecent(limit int) ([]TelegramSyncLog, error) {
	if l.db == nil {
		return nil, nil
	}

	var logs []TelegramSyncLog
	if err := l.db.Order("create_time DESC").Limit(limit).Find(&logs).Error; err != nil {
		return nil, err
	}

	return logs, nil
}

// GetByChannel 获取指定频道的日志
func (l *LogManager) GetByChannel(channelID int64, limit int) ([]TelegramSyncLog, error) {
	if l.db == nil {
		return nil, nil
	}

	var logs []TelegramSyncLog
	if err := l.db.Where("channel_id = ?", channelID).
		Order("create_time DESC").
		Limit(limit).
		Find(&logs).Error; err != nil {
		return nil, err
	}

	return logs, nil
}

// GetByStatus 获取指定状态的日志
func (l *LogManager) GetByStatus(status int, limit int) ([]TelegramSyncLog, error) {
	if l.db == nil {
		return nil, nil
	}

	var logs []TelegramSyncLog
	if err := l.db.Where("status = ?", status).
		Order("create_time DESC").
		Limit(limit).
		Find(&logs).Error; err != nil {
		return nil, err
	}

	return logs, nil
}

// Cleanup 清理过期日志
func (l *LogManager) Cleanup() error {
	if l.db == nil {
		return nil
	}

	if l.keepLogDays <= 0 {
		return nil
	}

	cutoff := time.Now().AddDate(0, 0, -l.keepLogDays).Unix()

	result := l.db.Where("create_time < ?", cutoff).Delete(&TelegramSyncLog{})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected > 0 {
		l.log.Info("清理过期日志", zap.Int64("deleted", result.RowsAffected))
	}

	return nil
}

// GetStats 获取统计信息
func (l *LogManager) GetStats(days int) (*SyncStats, error) {
	if l.db == nil {
		return nil, nil
	}

	cutoff := time.Now().AddDate(0, 0, -days).Unix()

	stats := &SyncStats{}

	// 总数
	var total int64
	l.db.Model(&TelegramSyncLog{}).Where("create_time >= ?", cutoff).Count(&total)
	stats.TotalSynced = int(total)

	// 成功数
	var success int64
	l.db.Model(&TelegramSyncLog{}).Where("create_time >= ? AND status = ?", cutoff, 1).Count(&success)
	stats.TotalSuccess = int(success)

	// 失败数
	var failed int64
	l.db.Model(&TelegramSyncLog{}).Where("create_time >= ? AND status = ?", cutoff, 0).Count(&failed)
	stats.TotalFailed = int(failed)

	// 跳过数
	var skipped int64
	l.db.Model(&TelegramSyncLog{}).Where("create_time >= ? AND status = ?", cutoff, 2).Count(&skipped)
	stats.TotalSkipped = int(skipped)

	// 成功率
	if total > 0 {
		stats.SuccessRate = float64(success) / float64(total) * 100
	}

	return stats, nil
}

// SyncStats 同步统计
type SyncStats struct {
	TotalSynced  int     `json:"total_synced"`
	TotalSuccess int     `json:"total_success"`
	TotalFailed  int     `json:"total_failed"`
	TotalSkipped int     `json:"total_skipped"`
	SuccessRate  float64 `json:"success_rate"`
}

// CheckDuplicate 检查消息是否已同步
func (l *LogManager) CheckDuplicate(channelID int64, messageID int64) bool {
	if l.db == nil {
		return false
	}

	var count int64
	l.db.Model(&TelegramSyncLog{}).
		Where("channel_id = ? AND message_id = ?", channelID, messageID).
		Count(&count)

	return count > 0
}

// GetLastMessageID 获取频道最后处理的消息 ID
func (l *LogManager) GetLastMessageID(channelID int64) int64 {
	if l.db == nil {
		return 0
	}

	var log TelegramSyncLog
	if err := l.db.Where("channel_id = ?", channelID).
		Order("message_id DESC").
		First(&log).Error; err != nil {
		return 0
	}

	return log.MessageID
}
