# Data Model: Fix Content Display Issues

**Branch**: `010-fix-content-display` | **Date**: 2026-05-11

## Entity Changes

**No schema changes required.** This is a bug fix that modifies data processing logic, not the data model.

## Entities Involved

### Article (existing, no changes)

| Field | Type | Description |
|-------|------|-------------|
| `ArticleBase.Description` | varchar(250) | Plain-text summary of the article. **Bug**: currently stores truncated HTML |
| `ArticleDetail.Content` | text | Full HTML body of the article. **Working correctly** |
| `ArticleDetail.ContentType` | varchar(20) | Content format ("html" or "markdown") |
| `ArticleBase.Thumbnail` | varchar(250) | Thumbnail image URL (relative path) |
| `ArticleDetail.MediaUrls` | text | JSON array of media URLs |
| `ArticleDetail.VideoUrl` | varchar(250) | Video URL (relative path) |

### TelegramMedia (existing, no changes)

| Field | Type | Description |
|-------|------|-------------|
| `ID` | int64 | Telegram media file ID |
| `MediaId` | int64 | Unique identifier used in URL paths |
| `MediaType` | string | Type of media (photo, video, document) |
| `AccessHash` | int64 | Telegram access hash for downloading |
| `FileReference` | []byte | Telegram file reference bytes |
| `StorageURL` | string | Cached/processed URL if available |

## Data Flow

```
Telegram Message
    │
    ▼
processSingleMessage() ─── convertToHTML() ──→ Content (HTML, correct)
    │
    ▼
truncateDescription() ──→ Description (BUG: truncated HTML)
    │                         Should be: plain text
    ▼
API Response
    ├── article list → description field (shows HTML tags on frontend)
    ├── article detail → content field (renders HTML correctly)
    └── media URLs → /api/telegram/media/{id} (404 - route not found)
```

## Processing Changes

### 1. Description Sanitization

**Where**: `TelegramChannelSync.go` - `truncateDescription()` function

**Before**: Truncates raw content (may contain HTML) to 200 runes
**After**: Strip HTML tags first, then truncate to 200 runes

### 2. API Response Sanitization

**Where**: `api.go` - `APIArticleList()` and `APIArticleDetail()`

**Before**: Returns `description` field as-is from database
**After**: Strip HTML tags from `description` before returning in API response

### 3. Media Route Registration

**Where**: `api.go` - `RegisterAPIRoutes()`

**Before**: No telegram media route in public API
**After**: Add `GET /api/telegram/media/:mediaId` route with CORS support
