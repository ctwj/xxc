# Research: Fix Content Display Issues

**Branch**: `010-fix-content-display` | **Date**: 2026-05-11

## Research Task 1: HTML Tags Appearing in Article Description

### Decision
Strip HTML tags from the `description` field at two points: (1) during article creation/storage via `truncateDescription`, and (2) in the API response for backward compatibility with existing data.

### Rationale
- The website (https://www.l9.lc) confirms all 23 articles show `<p>`, `</p>` tags as visible text in the description area
- The `description` field in `ArticleBase` is varchar(250), intended as a plain-text summary
- The card component (`InfoCard.tsx`) renders description as plain text via React JSX — HTML auto-escapes
- Fixing at storage prevents the issue for new articles; fixing at API response handles existing data
- No database migration needed — `description` is display-only metadata

### Root Cause
`truncateDescription()` in `TelegramChannelSync.go:976` truncates the content which may contain HTML, producing partial tags like `<p>Some text that got trunca...`. The `convertToHTML()` function wraps content in `<p>` tags, and the description is set from this HTML content.

### Alternatives Considered
| Alternative | Rejected Because |
|------------|-----------------|
| Use `dangerouslySetInnerHTML` in InfoCard | Spec requires plain text on cards; XSS risk with user content |
| Database migration to clean existing data | Overkill for display-only field; can be handled at API layer |
| Strip HTML only in frontend | Backend should return clean data; multiple frontend clients would all need the fix |

---

## Research Task 2: Media URLs Returning 404

### Decision
Add the Telegram media route to the public API router (`api.go`) so it is accessible at `/api/telegram/media/:mediaId`, matching the URLs already stored in the database.

### Rationale
- Media URLs generated in `media.go` (lines 104, 188, 197) use `/api/telegram/media/{id}`
- The route is only registered in `admin.go:32` under the admin path prefix, resolving to `/admin/api/telegram/media/:mediaId`
- No authentication is required (registered before `auth()` middleware at `admin.go:34`)
- Adding to public API router provides CORS support (critical for cross-origin frontend requests)
- Existing stored URLs in the database will work without data migration
- The admin route continues to work (no breaking change)

### Root Cause
Route registration path mismatch:
- Route registered: `/admin/api/telegram/media/:mediaId` (via `admin.go`)
- URL generated: `/api/telegram/media/{id}` (via `media.go`)
- Public API router (`api.go`) has no telegram media route → 404

### Alternatives Considered
| Alternative | Rejected Because |
|------------|-----------------|
| Change URL generation to include `/admin` prefix | Would break existing stored URLs; requires data migration |
| Move route from admin to public | Could break admin panel media features; dual registration is safer |
| Add URL rewrite/proxy rule | Adds complexity; route duplication is simpler and more maintainable |

---

## Research Task 3: CORS Configuration for Media

### Decision
By adding the media route to the public API router, CORS is automatically handled by the existing `CORSConfig()` middleware applied to `/api/*` routes.

### Rationale
- Public API routes in `api.go` use `app.Group("/api", middleware.CORSConfig())`
- Admin routes do NOT have CORS middleware
- Frontend (www.l9.lc) and API (api.l9.lc) are on different origins → CORS is required
- No additional CORS configuration needed

---

## Research Task 4: Article Detail Page Content Rendering

### Decision
No change needed for article detail page. The `dangerouslySetInnerHTML` in `article/[slug]/page.tsx:108` correctly renders HTML content.

### Rationale
- Article detail page uses `dangerouslySetInnerHTML={{ __html: article.content || "" }}` which is correct
- The `Content` field (full HTML body) renders properly as rich text
- Only the `description` field (used on cards) has the HTML tag visibility issue
