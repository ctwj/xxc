# API Contract Changes: Fix Content Display Issues

**Branch**: `010-fix-content-display` | **Date**: 2026-05-11

## Existing Endpoints Modified

### GET /api/articles

**Change**: `description` field in response will no longer contain HTML tags.

**Before** (buggy):
```json
{
  "id": 1,
  "title": "Article Title",
  "description": "<p>Some article content that gets truncated...</p>",
  "thumbnail": "/api/telegram/media/123"
}
```

**After** (fixed):
```json
{
  "id": 1,
  "title": "Article Title",
  "description": "Some article content that gets truncated...",
  "thumbnail": "/api/telegram/media/123"
}
```

### GET /api/articles/:slug

**Change**: `description` field in response will no longer contain HTML tags. `content` field is unchanged.

**Before** (buggy):
```json
{
  "id": 1,
  "title": "Article Title",
  "description": "<p>Some article content that gets truncated...</p>",
  "content": "<p>Full article HTML content</p><p>Multiple paragraphs</p>"
}
```

**After** (fixed):
```json
{
  "id": 1,
  "title": "Article Title",
  "description": "Some article content that gets truncated...",
  "content": "<p>Full article HTML content</p><p>Multiple paragraphs</p>"
}
```

## New Endpoint

### GET /api/telegram/media/:mediaId

**Purpose**: Serve Telegram media files (images, videos, documents) to public API consumers.

**Note**: This route already exists at `/admin/api/telegram/media/:mediaId` (admin router). The new route is a duplicate registration in the public API router to make media accessible without the `/admin` path prefix.

**Authentication**: None required.

**Query Parameters**:
| Parameter | Type | Description |
|-----------|------|-------------|
| `thumb` | int | Optional. Set to `1` to get thumbnail instead of full image |

**Response**:
- **Success (200)**: Binary file content with appropriate `Content-Type` header (image/jpeg, video/mp4, etc.)
  - Headers: `Cache-Control: public, max-age=31536000`
- **Not Found (404)**: Media record not found in database
- **Error (500)**: Failed to download from Telegram servers

**CORS**: Handled by public API router middleware (allows configured origins)

**Caching**: Files cached locally at `./upload/telegram_cache/{mediaId}.{ext}` for 1 year
