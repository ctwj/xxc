package telegram_sync

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gotd/td/tg"
	"go.uber.org/zap"
)

// MediaHandler 媒体处理器
type MediaHandler struct {
	log         *zap.Logger
	maxImageSize int
	httpClient  *http.Client
}

// NewMediaHandler 创建媒体处理器
func NewMediaHandler(log *zap.Logger) *MediaHandler {
	return &MediaHandler{
		log:          log,
		maxImageSize: 10485760, // 10MB
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// SetMaxImageSize 设置最大图片大小
func (m *MediaHandler) SetMaxImageSize(size int) {
	m.maxImageSize = size
}

// ProcessMedia 处理媒体内容
func (h *MediaHandler) ProcessMedia(ctx context.Context, media tg.MessageMediaClass) (string, error) {
	if media == nil {
		return "", nil
	}

	switch m := media.(type) {
	case *tg.MessageMediaPhoto:
		return h.processPhoto(ctx, m)
	case *tg.MessageMediaDocument:
		return h.processDocument(ctx, m)
	default:
		h.log.Debug("不支持的媒体类型", zap.String("type", fmt.Sprintf("%T", media)))
		return "", nil
	}
}

// processPhoto 处理图片
func (h *MediaHandler) processPhoto(ctx context.Context, photo *tg.MessageMediaPhoto) (string, error) {
	p, ok := photo.Photo.(*tg.Photo)
	if !ok {
		return "", fmt.Errorf("invalid photo type")
	}

	// 找到最大的图片尺寸
	var largest *tg.PhotoSize
	for _, size := range p.Sizes {
		switch s := size.(type) {
		case *tg.PhotoSize:
			if largest == nil || s.W*s.H > largest.W*largest.H {
				largest = s
			}
		case *tg.PhotoSizeProgressive:
			if largest == nil || s.W*s.H > largest.W*largest.H {
				largest = &tg.PhotoSize{
					Type: s.Type,
					W:    s.W,
					H:    s.H,
				}
			}
		}
	}

	if largest == nil {
		return "", fmt.Errorf("no photo size found")
	}

	// TODO: 使用 Telegram API 下载图片
	// 需要通过 gotd/td 的 download API 实现

	h.log.Info("处理图片",
		zap.Int64("photo_id", p.ID),
		zap.Int("width", largest.W),
		zap.Int("height", largest.H))

	return "", nil // 返回下载后的 URL
}

// processDocument 处理文档
func (h *MediaHandler) processDocument(ctx context.Context, doc *tg.MessageMediaDocument) (string, error) {
	d, ok := doc.Document.(*tg.Document)
	if !ok {
		return "", fmt.Errorf("invalid document type")
	}

	// 检查是否为图片
	isImage := false
	for _, attr := range d.Attributes {
		if _, ok := attr.(*tg.DocumentAttributeImageSize); ok {
			isImage = true
			break
		}
	}

	if !isImage {
		h.log.Debug("文档不是图片，跳过", zap.Int64("doc_id", d.ID))
		return "", nil
	}

	// 检查文件大小
	if d.Size > int64(h.maxImageSize) {
		h.log.Warn("图片大小超过限制",
			zap.Int64("size", d.Size),
			zap.Int("max", h.maxImageSize))
		return "", nil
	}

	// TODO: 使用 Telegram API 下载文档
	h.log.Info("处理图片文档",
		zap.Int64("doc_id", d.ID),
		zap.Int64("size", d.Size))

	return "", nil // 返回下载后的 URL
}

// DownloadFromURL 从 URL 下载文件
func (h *MediaHandler) DownloadFromURL(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := h.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("download failed: status %d", resp.StatusCode)
	}

	// 检查大小
	if resp.ContentLength > int64(h.maxImageSize) {
		return nil, fmt.Errorf("file too large: %d bytes (max: %d)", resp.ContentLength, h.maxImageSize)
	}

	// 读取内容
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// UploadToStorage 上传到存储
func (h *MediaHandler) UploadToStorage(ctx context.Context, data []byte, filename string) (string, error) {
	// TODO: 使用 Moss 的上传基础设施
	// 需要集成 upload.Upload 函数

	h.log.Info("上传文件", zap.String("filename", filename), zap.Int("size", len(data)))

	return "", nil // 返回上传后的 URL
}

// DownloadAndUpload 下载并上传
func (h *MediaHandler) DownloadAndUpload(ctx context.Context, url string, filename string) (string, error) {
	// 下载
	data, err := h.DownloadFromURL(ctx, url)
	if err != nil {
		return "", err
	}

	// 上传
	return h.UploadToStorage(ctx, data, filename)
}

// GetMediaType 获取媒体类型
func (h *MediaHandler) GetMediaType(media tg.MessageMediaClass) string {
	if media == nil {
		return "none"
	}

	switch media.(type) {
	case *tg.MessageMediaPhoto:
		return "photo"
	case *tg.MessageMediaDocument:
		return "document"
	default:
		return "unknown"
	}
}