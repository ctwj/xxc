# Quickstart: Fix Content Display Issues

**Branch**: `010-fix-content-display` | **Date**: 2026-05-11

## Overview

Fix two critical display bugs on the public-facing website (https://www.l9.lc):
1. HTML tags appearing as raw text in article card descriptions
2. Images and videos returning 404 due to route path mismatch

## Changes Required

### Backend (Go)

**File 1: `main/plugins/TelegramChannelSync.go`**
- Modify `truncateDescription()` to strip HTML tags before truncating
- This fixes new articles created via Telegram sync

**File 2: `main/api/web/controller/api.go`**
- Add HTML tag stripping to `description` field in `APIArticleList()` and `APIArticleDetail()`
- This fixes the API response for existing articles already in the database

**File 3: `main/api/web/router/api.go`**
- Add `GET /telegram/media/:mediaId` route pointing to `controller.TelegramGetMedia`
- This makes media accessible at `/api/telegram/media/{id}` (matching stored URLs)

### Frontend (No Changes Required)
- Frontend rendering logic is correct
- Fixing the backend API responses resolves all display issues

## Verification Steps

1. Start dev environment: `task dev`
2. Visit http://localhost:3000 — check card descriptions have no HTML tags
3. Check article cards with images — verify thumbnails load
4. Open an article detail page — verify inline images and video load
5. Visit https://www.l9.lc after deployment to confirm production fix

## Deployment

Standard deployment — no database migration, no config changes needed.
