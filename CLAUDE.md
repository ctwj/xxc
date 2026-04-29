# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Moss is a lightweight content management system (CMS) built with Go backend and Vue 3 frontend. It features a plugin architecture, multi-database support, and internationalization.

## Architecture

### Backend Structure (Go)
- **Entry Point**: `main/cmd/web/main.go`
- **Web Framework**: Fiber (gofiber/fiber/v2)
- **Layered Architecture**:
  - `api/web/` - Web API layer (controllers, DTOs, mappers, middleware, routers)
  - `application/` - Application services (orchestrates domain logic)
  - `domain/` - Domain models and business logic
    - `domain/core/` - Core entities (Article, Category, Tag, Link, Store)
    - `domain/config/` - Configuration entities
  - `infrastructure/` - Infrastructure layer
    - `persistent/` - Database, storage drivers (local, S3, OSS, COS, FTP, B2)
    - `support/` - Cache (Badger, Redis, Memcached), logging, templates, upload
    - `utils/` - Utility functions
  - `plugins/` - Plugin system for extensible functionality
  - `resources/` - Embedded static resources (admin, themes)
  - `startup/` - Application initialization and plugin registration

### Frontend Structure (Vue 3)
- **Admin Panel**: `admin/` - Vue 3 + Vite + Tailwind CSS + Arco Design
- **Build Output**: `main/resources/admin/` (embedded in binary)
- **API Proxy**: Dev server proxies `/admin/api/*` to backend at `http://127.0.0.1:9008`

## Development Commands

```bash
# Install dependencies
task init-admin          # Frontend dependencies
cd main && go mod tidy   # Backend dependencies

# Development
task dev                 # Start full development environment (both frontend and backend with hot reload)
task run                 # Start backend only (no hot reload)
task admin               # Start frontend only

# Testing
cd main && go test ./...                     # Run all backend tests
cd main && go test -run TestFunctionName ./path/to/package  # Run specific test

# Build
task build               # Build both frontend and backend for production
task build-admin         # Build frontend only
task build-main          # Cross-compile backend for multiple platforms

# Utilities
task status              # Check development environment status
task reset-admin         # Reset admin credentials (admin/admin123)
```

### Development Environment
- **Backend Hot Reload**: Uses Air tool (config: `main/.air.toml`), monitors `.go`, `.tpl`, `.tmpl`, `.html`, `.toml` files
- **Frontend Hot Reload**: Vite dev server with HMR
- **Default Ports**:
  - Backend: 9008
  - Frontend: 3000

## Plugin System

Plugins are located in `main/plugins/` and implement specific interfaces. Key plugin types:
- **Content Processing**: ArticleSanitizer, GenerateSlug, GenerateDescription
- **Media Processing**: SaveArticleImages, MakeCarousel
- **SEO**: PushToBaidu, PushToBing, PushToSearchEngine
- **Automation**: GnDownSpider, AISeoPlugin
- **Cloud Transfer**: BaiduCloudTransfer, QuarkCloudTransfer

Plugins are registered in `main/startup/startup.go`.

## Database Support

Supports SQLite (default), MySQL, and PostgreSQL. Configuration via `main/conf.toml`:
- SQLite: `./moss.db?_pragma=journal_mode(WAL)`
- MySQL: `user:password@tcp(host:port)/dbname?charset=utf8mb4&parseTime=True`
- PostgreSQL: `host=127.0.0.1 port=5432 user=postgres password=123456 dbname=moss sslmode=disable`

GORM handles migrations automatically.

## Key Development Notes

1. **Frontend Hot Reload**: Changes in `admin/src/` are automatically applied - no restart needed
2. **Plugin Development**: Create new plugins in `main/plugins/` and register in `startup.go`
3. **Configuration**: Use `main/conf.toml` for runtime configuration (created on first run)
4. **Multi-language**: Admin panel supports 12 languages (see `admin/src/locale/lang/`)
5. **Template Engine**: Uses Jet template engine (`infrastructure/support/template/engine/jet.go`)
6. **Storage Drivers**: Supports local, S3, OSS, COS, FTP, B2 - configured via admin panel

<!-- SPECKIT START -->
## Active Feature: Next.js Frontend Integration

**Branch**: `001-react-frontend-integration`
**Plan**: [specs/001-react-frontend-integration/plan.md](specs/001-react-frontend-integration/plan.md)
**Status**: Implementation Complete

### Key Documents
- [Specification](specs/001-react-frontend-integration/spec.md)
- [Research](specs/001-react-frontend-integration/research.md)
- [Data Model](specs/001-react-frontend-integration/data-model.md)
- [API Contracts](specs/001-react-frontend-integration/contracts/api.md)
- [Quickstart](specs/001-react-frontend-integration/quickstart.md)

### Architecture Decision
Headless CMS 模式：Next.js (ISR) 前端 + Go REST API 后端

### Implementation Summary

#### Backend API (Go/Fiber)
- **CORS**: `main/api/web/middleware/cors.go` - Enables cross-origin requests from frontend
- **JWT Auth**: `main/infrastructure/support/auth/jwt.go` - Token generation and validation
- **Public API**: `main/api/web/controller/api.go` - Article, category, tag, search endpoints
- **Auth API**: `main/api/web/controller/auth.go` - Login, logout, current user
- **Favorites**: `main/domain/core/entity/favorite.go`, `repository/favorite.go`, `service/favorite.go`
- **Webhook**: `main/api/web/controller/webhook.go` - ISR revalidation trigger
- **Routes**: `main/api/web/router/api.go` - All `/api/*` endpoints

#### Frontend (Next.js 15)
- **Location**: `frontend/`
- **Framework**: Next.js 15 with App Router, TypeScript, Tailwind CSS
- **ISR**: 60-second revalidation for static pages
- **Pages**: Home, Article detail, Categories, Tags, Search, Login, Register, Favorites
- **API Client**: `frontend/src/lib/api.ts`
- **Auth Context**: `frontend/src/contexts/AuthContext.tsx`

#### Deployment
- **Frontend**: Vercel (configure `NEXT_PUBLIC_API_URL` and `REVALIDATE_SECRET`)
- **Backend**: Any Go hosting (set `CORSOrigins` in config)

### Development Commands

```bash
# Start backend
cd main && go run cmd/web/main.go

# Start frontend dev server
cd frontend && npm run dev

# Build frontend for production
cd frontend && npm run build
```

### Environment Variables

**Frontend (.env.local)**:
- `NEXT_PUBLIC_API_URL` - Backend API URL (e.g., http://localhost:9008)
- `REVALIDATE_SECRET` - Secret for webhook authentication

**Backend (conf.toml)**:
- `CORSOrigins` - Allowed origins for CORS (e.g., http://localhost:3000,https://your-domain.com)
<!-- SPECKIT END -->

