package telegram_sync

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gotd/td/telegram/downloader"
	"github.com/gotd/td/tg"
	"go.uber.org/zap"
	"moss/infrastructure/persistent/storage"
	"moss/infrastructure/support/upload"
)

// MediaHandler 媒体处理器
type MediaHandler struct {
	log          *zap.Logger
	maxImageSize int
	api          *tg.Client
}

// NewMediaHandler 创建媒体处理器
func NewMediaHandler(log *zap.Logger) *MediaHandler {
	return &MediaHandler{
		log:          log,
		maxImageSize: 10485760, // 10MB
	}
}

// SetMaxImageSize 设置最大图片大小
func (m *MediaHandler) SetMaxImageSize(size int) {
	m.maxImageSize = size
}

// SetAPI 设置 API 客户端
func (m *MediaHandler) SetAPI(api *tg.Client) {
	m.api = api
}

// ProcessMedia 处理媒体内容，返回媒体信息列表（只提取元数据，不下载）
func (h *MediaHandler) ProcessMedia(ctx context.Context, msg *tg.Message, channelID, messageID int64) ([]MessageMediaInfo, error) {
	if msg.Media == nil {
		return nil, nil
	}

	var mediaInfos []MessageMediaInfo

	switch media := msg.Media.(type) {
	case *tg.MessageMediaPhoto:
		info, err := h.extractPhotoMetadata(media, channelID, messageID)
		if err != nil {
			h.log.Error("提取图片元数据失败", zap.Error(err))
			return nil, err
		}
		if info != nil {
			mediaInfos = append(mediaInfos, *info)
		}

	case *tg.MessageMediaDocument:
		infos, err := h.extractDocumentMetadata(media, channelID, messageID)
		if err != nil {
			h.log.Error("提取文档元数据失败", zap.Error(err))
			return nil, err
		}
		mediaInfos = append(mediaInfos, infos...)

	default:
		h.log.Debug("不支持的媒体类型", zap.String("type", fmt.Sprintf("%T", media)))
	}

	return mediaInfos, nil
}

// extractPhotoMetadata 提取图片元数据（不下载）
func (h *MediaHandler) extractPhotoMetadata(photo *tg.MessageMediaPhoto, channelID, messageID int64) (*MessageMediaInfo, error) {
	p, ok := photo.Photo.(*tg.Photo)
	if !ok {
		return nil, fmt.Errorf("invalid photo type")
	}

	// 找到最大的图片尺寸
	var largestWidth, largestHeight int
	for _, size := range p.Sizes {
		switch s := size.(type) {
		case *tg.PhotoSize:
			if largestWidth*largestHeight < s.W*s.H {
				largestWidth = s.W
				largestHeight = s.H
			}
		case *tg.PhotoSizeProgressive:
			if largestWidth*largestHeight < s.W*s.H {
				largestWidth = s.W
				largestHeight = s.H
			}
		}
	}

	// 返回媒体 API URL 而不是直接下载
	mediaURL := fmt.Sprintf("/api/telegram/media/%d", p.ID)

	h.log.Info("提取图片元数据成功",
		zap.Int64("photo_id", p.ID),
		zap.Int("width", largestWidth),
		zap.Int("height", largestHeight))

	return &MessageMediaInfo{
		MediaID:      p.ID,
		MediaType:    "photo",
		URL:          mediaURL,
		AccessHash:   p.AccessHash,
		FileReference: p.FileReference,
		Filename:     fmt.Sprintf("telegram_photo_%d_%d.jpg", channelID, p.ID),
		Width:        largestWidth,
		Height:       largestHeight,
	}, nil
}

// extractDocumentMetadata 提取文档元数据（不下载）
func (h *MediaHandler) extractDocumentMetadata(doc *tg.MessageMediaDocument, channelID, messageID int64) ([]MessageMediaInfo, error) {
	d, ok := doc.Document.(*tg.Document)
	if !ok {
		return nil, fmt.Errorf("invalid document type")
	}

	// 获取文件属性
	var filename string
	var width, height int
	var isImage, isVideo bool
	var thumbURL string
	var thumbID int64
	var thumbAccessHash int64
	var thumbFileRef []byte

	for _, attr := range d.Attributes {
		switch a := attr.(type) {
		case *tg.DocumentAttributeFilename:
			filename = a.FileName
		case *tg.DocumentAttributeImageSize:
			isImage = true
			width = a.W
			height = a.H
		case *tg.DocumentAttributeVideo:
			isVideo = true
			width = a.W
			height = a.H
		}
	}

	// 只处理图片和视频
	if !isImage && !isVideo {
		h.log.Debug("文档不是图片或视频，跳过", zap.Int64("doc_id", d.ID))
		return nil, nil
	}

	// 提取视频封面（从 Thumbs）
	if isVideo && len(d.Thumbs) > 0 {
		// 找最大的缩略图
		var largestThumb tg.PhotoSizeClass
		var largestW, largestH int
		for _, thumb := range d.Thumbs {
			switch t := thumb.(type) {
			case *tg.PhotoSize:
				if largestW*largestH < t.W*t.H {
					largestW = t.W
					largestH = t.H
					largestThumb = t
				}
			case *tg.PhotoSizeProgressive:
				if largestW*largestH < t.W*t.H {
					largestW = t.W
					largestH = t.H
					largestThumb = t
				}
			}
		}

		// 如果找到了缩略图，记录其信息
		if largestThumb != nil {
			// 缩略图使用与文档相同的 ID 和 access_hash
			thumbID = d.ID
			thumbAccessHash = d.AccessHash
			thumbFileRef = d.FileReference
			thumbURL = fmt.Sprintf("/api/telegram/media/%d?thumb=1", d.ID)
			h.log.Info("找到视频封面缩略图",
				zap.Int64("doc_id", d.ID),
				zap.Int("width", largestW),
				zap.Int("height", largestH))
		}
	}

	// 返回媒体 API URL
	mediaURL := fmt.Sprintf("/api/telegram/media/%d", d.ID)

	h.log.Info("提取文档元数据成功",
		zap.Int64("doc_id", d.ID),
		zap.String("mime_type", d.MimeType),
		zap.Int64("size", d.Size),
		zap.Bool("has_thumb", thumbURL != ""))

	mediaType := "photo"
	if isVideo {
		mediaType = "video"
	}

	if filename == "" {
		ext := ".jpg"
		if isVideo {
			ext = ".mp4"
		}
		filename = fmt.Sprintf("telegram_media_%d_%d%s", channelID, d.ID, ext)
	}

	return []MessageMediaInfo{
		{
			MediaID:        d.ID,
			MediaType:      mediaType,
			URL:            mediaURL,
			AccessHash:     d.AccessHash,
			FileReference:  d.FileReference,
			Filename:       filename,
			Width:          width,
			Height:         height,
			ThumbURL:       thumbURL,
			ThumbID:        thumbID,
			ThumbAccessHash: thumbAccessHash,
			ThumbFileRef:   thumbFileRef,
		},
	}, nil
}

// processPhoto 处理图片
func (h *MediaHandler) processPhoto(ctx context.Context, photo *tg.MessageMediaPhoto, channelID, messageID int64) (*MessageMediaInfo, error) {
	p, ok := photo.Photo.(*tg.Photo)
	if !ok {
		return nil, fmt.Errorf("invalid photo type")
	}

	// 找到最大的图片尺寸
	var largestSize tg.PhotoSizeClass
	var largestWidth, largestHeight int

	for _, size := range p.Sizes {
		switch s := size.(type) {
		case *tg.PhotoSize:
			if largestWidth*largestHeight < s.W*s.H {
				largestWidth = s.W
				largestHeight = s.H
				largestSize = s
			}
		case *tg.PhotoSizeProgressive:
			if largestWidth*largestHeight < s.W*s.H {
				largestWidth = s.W
				largestHeight = s.H
				largestSize = s
			}
		}
	}

	if largestSize == nil {
		return nil, fmt.Errorf("no photo size found")
	}

	// 检查文件大小
	var fileSize int64
	switch s := largestSize.(type) {
	case *tg.PhotoSize:
		fileSize = int64(s.Size)
	case *tg.PhotoSizeProgressive:
		// Progressive size 没有直接的 Size 字段，使用最大尺寸估算
		fileSize = int64(largestWidth * largestHeight * 3 / 10) // 估算压缩后大小
	}

	if fileSize > int64(h.maxImageSize) {
		h.log.Warn("图片大小超过限制",
			zap.Int64("size", fileSize),
			zap.Int("max", h.maxImageSize))
		return nil, nil
	}

	// 下载图片
	d := downloader.NewDownloader()
	location := &tg.InputPhotoFileLocation{
		ID:            p.ID,
		AccessHash:    p.AccessHash,
		FileReference: p.FileReference,
		ThumbSize:     getPhotoSizeType(largestSize),
	}

	var buf bytes.Buffer
	_, err := d.Download(h.api, location).Stream(ctx, &buf)
	if err != nil {
		h.log.Error("下载图片失败", zap.Error(err), zap.Int64("photo_id", p.ID))
		return nil, fmt.Errorf("download failed: %w", err)
	}

	h.log.Info("图片下载成功",
		zap.Int64("photo_id", p.ID),
		zap.Int("width", largestWidth),
		zap.Int("height", largestHeight),
		zap.Int("size", buf.Len()))

	// 上传到存储
	filename := fmt.Sprintf("telegram_photo_%d_%d.jpg", channelID, p.ID)
	ext := ".jpg"

	result, err := upload.Upload(filename, ext, storage.NewSetValueBytes(buf.Bytes()))
	if err != nil {
		h.log.Error("上传图片失败", zap.Error(err))
		return nil, fmt.Errorf("upload failed: %w", err)
	}

	h.log.Info("图片上传成功",
		zap.String("url", result.URL),
		zap.String("path", result.FullPath))

	return &MessageMediaInfo{
		MediaID:   p.ID,
		MediaType: "photo",
		URL:       result.URL,
		Filename:  filename,
		Width:     largestWidth,
		Height:    largestHeight,
	}, nil
}

// processDocument 处理文档（包括图片文档）
func (h *MediaHandler) processDocument(ctx context.Context, doc *tg.MessageMediaDocument, channelID, messageID int64) ([]MessageMediaInfo, error) {
	d, ok := doc.Document.(*tg.Document)
	if !ok {
		return nil, fmt.Errorf("invalid document type")
	}

	// 检查文件大小
	if d.Size > int64(h.maxImageSize) {
		h.log.Warn("文档大小超过限制",
			zap.Int64("size", d.Size),
			zap.Int("max", h.maxImageSize))
		return nil, nil
	}

	// 获取文件属性
	var filename string
	var width, height int
	var isImage, isVideo bool

	for _, attr := range d.Attributes {
		switch a := attr.(type) {
		case *tg.DocumentAttributeFilename:
			filename = a.FileName
		case *tg.DocumentAttributeImageSize:
			isImage = true
			width = a.W
			height = a.H
		case *tg.DocumentAttributeVideo:
			isVideo = true
			width = a.W
			height = a.H
		}
	}

	// 只处理图片和视频
	if !isImage && !isVideo {
		h.log.Debug("文档不是图片或视频，跳过", zap.Int64("doc_id", d.ID))
		return nil, nil
	}

	// 下载文档
	dl := downloader.NewDownloader()
	location := d.AsInputDocumentFileLocation()

	var buf bytes.Buffer
	_, err := dl.Download(h.api, location).Stream(ctx, &buf)
	if err != nil {
		h.log.Error("下载文档失败", zap.Error(err), zap.Int64("doc_id", d.ID))
		return nil, fmt.Errorf("download failed: %w", err)
	}

	h.log.Info("文档下载成功",
		zap.Int64("doc_id", d.ID),
		zap.String("mime_type", d.MimeType),
		zap.Int("size", buf.Len()))

	// 确定文件扩展名
	ext := filepath.Ext(filename)
	if ext == "" {
		ext = filepath.Ext(d.MimeType)
		if ext == "" {
			if isImage {
				ext = ".jpg"
			} else if isVideo {
				ext = ".mp4"
			}
		}
	}

	if filename == "" {
		filename = fmt.Sprintf("telegram_media_%d_%d%s", channelID, d.ID, ext)
	}

	// 上传到存储
	result, err := upload.Upload(filename, ext, storage.NewSetValueBytes(buf.Bytes()))
	if err != nil {
		h.log.Error("上传文档失败", zap.Error(err))
		return nil, fmt.Errorf("upload failed: %w", err)
	}

	h.log.Info("文档上传成功",
		zap.String("url", result.URL),
		zap.String("path", result.FullPath))

	mediaType := "photo"
	if isVideo {
		mediaType = "video"
	}

	return []MessageMediaInfo{
		{
			MediaID:   d.ID,
			MediaType: mediaType,
			URL:       result.URL,
			Filename:  filename,
			Width:     width,
			Height:    height,
		},
	}, nil
}

// getPhotoSizeType 获取图片尺寸类型
func getPhotoSizeType(size tg.PhotoSizeClass) string {
	switch s := size.(type) {
	case *tg.PhotoSize:
		return s.Type
	case *tg.PhotoSizeProgressive:
		return s.Type
	default:
		return "w" // 默认使用最大尺寸
	}
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

// DownloadMediaByID 通过媒体 ID 下载媒体文件（用于 API）
func (h *MediaHandler) DownloadMediaByID(ctx context.Context, mediaID int64, accessHash int64, fileReference []byte, mediaType string) ([]byte, string, error) {
	if h.api == nil {
		return nil, "", fmt.Errorf("API client not initialized")
	}

	// 检查本地缓存
	cacheDir := "./upload/telegram_cache"
	cacheFile := fmt.Sprintf("%s/%d.%s", cacheDir, mediaID, getExtByType(mediaType))

	// 尝试从缓存读取
	if data, err := os.ReadFile(cacheFile); err == nil {
		mimeType := getMimeTypeByType(mediaType)
		h.log.Debug("从缓存读取媒体", zap.Int64("media_id", mediaID), zap.Int("size", len(data)))
		return data, mimeType, nil
	}

	// 使用并行下载优化速度
	d := downloader.NewDownloader().WithPartSize(512 * 1024) // 512KB 分片

	var location tg.InputFileLocationClass

	if mediaType == "photo" {
		location = &tg.InputPhotoFileLocation{
			ID:            mediaID,
			AccessHash:    accessHash,
			FileReference: fileReference,
			ThumbSize:     "w", // 最大尺寸
		}
	} else {
		location = &tg.InputDocumentFileLocation{
			ID:            mediaID,
			AccessHash:    accessHash,
			FileReference: fileReference,
		}
	}

	var buf bytes.Buffer
	// 使用 4 个并行线程下载
	_, err := d.Download(h.api, location).WithThreads(4).Stream(ctx, &buf)
	if err != nil {
		return nil, "", fmt.Errorf("download failed: %w", err)
	}

	data := buf.Bytes()
	mimeType := getMimeTypeByType(mediaType)

	// 保存到缓存
	if err := os.MkdirAll(cacheDir, 0755); err == nil {
		if err := os.WriteFile(cacheFile, data, 0644); err != nil {
			h.log.Warn("保存缓存失败", zap.Error(err), zap.String("file", cacheFile))
		} else {
			h.log.Debug("媒体已缓存", zap.Int64("media_id", mediaID), zap.String("file", cacheFile))
		}
	}

	return data, mimeType, nil
}

// getExtByType 根据媒体类型获取文件扩展名
func getExtByType(mediaType string) string {
	switch mediaType {
	case "photo":
		return "jpg"
	case "video":
		return "mp4"
	default:
		return "bin"
	}
}

// getMimeTypeByType 根据媒体类型获取 MIME 类型
func getMimeTypeByType(mediaType string) string {
	switch mediaType {
	case "photo":
		return "image/jpeg"
	case "video":
		return "video/mp4"
	default:
		return "application/octet-stream"
	}
}

// DownloadFromURL 从 URL 下载文件
func (h *MediaHandler) DownloadFromURL(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
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