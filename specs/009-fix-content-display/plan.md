# Implementation Plan: Fix Content Display Issues

**Branch**: `010-fix-content-display` | **Date**: 2026-05-11 | **Spec**: [spec.md](spec.md)
**Input**: Feature specification from `/specs/009-fix-content-display/spec.md`

## Summary

Fix two critical display bugs on the production website (https://www.l9.lc): (1) HTML tags (`<p>`, `<br>`, etc.) appearing as visible raw text in article card descriptions on the homepage, and (2) all images and videos returning 404 because media route paths don't match the URLs stored in the database. The fix requires three targeted backend changes — no frontend changes or database migrations needed.

## Technical Context

**Language/Version**: Go 1.21+ (backend), TypeScript/React/Next.js (frontend)
**Primary Dependencies**: Fiber (gofiber/fiber/v2), React 18, Next.js 14
**Storage**: SQLite (default, configurable MySQL/PostgreSQL)
**Testing**: `go test ./...` (backend)
**Target Platform**: Linux server (production), Windows (development)
**Project Type**: Web application (CMS with Go backend + Next.js frontend)
**Performance Goals**: N/A (bug fix, no new features)
**Constraints**: No database migration; existing stored URLs must continue to work
**Scale/Scope**: 23 articles, ~15 with images, 1 with video

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

Constitution is a blank template — no specific principles defined. Gate: **PASS**.

## Project Structure

### Documentation (this feature)

```text
specs/009-fix-content-display/
├── plan.md              # This file
├── spec.md              # Feature specification
├── research.md          # Phase 0 output — root cause analysis and solution decisions
├── data-model.md        # Phase 1 output — data flow and processing changes
├── quickstart.md        # Phase 1 output — change summary and verification steps
├── contracts/
│   └── api-changes.md   # Phase 1 output — API contract modifications
└── tasks.md             # Phase 2 output (/speckit-tasks)
```

### Source Code (files to modify)

```text
main/
├── plugins/
│   └── TelegramChannelSync.go    # Fix truncateDescription() to strip HTML
├── api/web/
│   ├── controller/
│   │   └── api.go                 # Strip HTML from description in API responses
│   └── router/
│       └── api.go                 # Add public telegram media route
```

**Structure Decision**: Modify existing files only. No new files or directories needed.

## Implementation Approach

### Change 1: Strip HTML from Description at Storage Time
**File**: `main/plugins/TelegramChannelSync.go`
**Function**: `truncateDescription()` (line ~976)

Modify `truncateDescription` to strip HTML tags before truncating. This ensures new articles store clean plain-text descriptions.

### Change 2: Strip HTML from Description in API Responses
**File**: `main/api/web/controller/api.go`
**Functions**: `APIArticleList()` (line ~13), `APIArticleDetail()` (line ~126)

Add HTML tag stripping to the `description` field in both API endpoints. This fixes the display for all existing articles already in the database.

### Change 3: Add Public Telegram Media Route
**File**: `main/api/web/router/api.go`
**Function**: `RegisterAPIRoutes()` (line ~11)

Add `GET /telegram/media/:mediaId` route with the `TelegramGetMedia` controller handler. This makes media files accessible at the URL path that matches what's stored in the database (`/api/telegram/media/{id}`), with CORS support from the public API middleware.

## Complexity Tracking

No constitution violations — no entries needed.
