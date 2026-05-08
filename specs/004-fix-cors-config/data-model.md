# Data Model: Fix CORS Configuration

**Date**: 2026-05-08
**Feature**: 004-fix-cors-config

## Entity: Router Config

The Router config entity already exists. This document describes the CORS-related fields.

### Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `cors_origins` | string | Comma-separated list of allowed CORS origins | Optional, empty = localhost only |

### JSON Structure (Database)

```json
{
  "admin_path": "/admin",
  "sitemap_path": "/sitemap",
  "article_rule": "/article/{slug}",
  "category_rule": "/category/{slug}",
  "category_page_rule": "/category/{slug}/{page}",
  "tag_rule": "/tag/{slug}",
  "tag_page_rule": "/tag/{slug}/{page}",
  "cors_origins": "https://www.l9.lc,https://app.l9.lc",
  "compress_level": 0,
  "minify_code": true,
  "etag": true,
  "proxy_header": []
}
```

### State Transitions

```
[Empty cors_origins]
       ↓
[User enters origins in admin UI]
       ↓
[Save to database]
       ↓
[Reload config in memory]
       ↓
[CORS middleware uses new origins]
```

## Entity: CORS Middleware State

The CORS middleware is stateless. It reads the config on each request.

### Decision Flow

```
Request arrives with Origin header
       ↓
Check if Origin matches any in cors_origins
       ↓
    [Yes] → Return CORS headers
    [No]  → Check if localhost
                ↓
            [Yes] → Return CORS headers (dev mode)
            [No]  → No CORS headers (blocked)
```

## No New Tables Required

This feature uses the existing `config` table with the `router` entry.
