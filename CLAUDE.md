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
## Active Feature: Telegram Channel Sync

**Branch**: `002-telegram-channel-sync`
**Plan**: [specs/002-telegram-channel-sync/plan.md](specs/002-telegram-channel-sync/plan.md)
**Status**: Planning Complete

### Key Documents
- [Specification](specs/002-telegram-channel-sync/spec.md)
- [Research](specs/002-telegram-channel-sync/research.md)
- [Data Model](specs/002-telegram-channel-sync/data-model.md)
- [API Contracts](specs/002-telegram-channel-sync/contracts/api.md)
- [Quickstart](specs/002-telegram-channel-sync/quickstart.md)

### Architecture Decision
Moss 插件模式：使用 gotd/td 库监听 Telegram 频道，自动同步消息为 CMS 文章

### Implementation Summary

#### Plugin Structure (简化架构)
- **Main Plugin**: `main/plugins/TelegramChannelSync.go` - 包含实体定义和核心逻辑
- **Sub Package**: `main/plugins/telegram_sync/` - 可选，用于拆分复杂逻辑
  - `client.go` - Telegram 客户端管理
  - `handler.go` - 消息处理与更新分发
  - `filter.go` - 消息过滤规则引擎
  - `media.go` - 媒体下载与处理
  - `session.go` - 会话持久化（加密存储）

#### Frontend Admin
- **Config Component**: `admin/src/views/plugin/options/TelegramChannelSync.vue`
- 集成到 Moss 现有的插件配置界面

**Note**: 所有后端代码集中在 `main/plugins/` 目录，复用现有 Article/Category 实体和服务，不创建独立的 entity/repository/service/controller 文件。前端只需添加一个配置表单组件。

### Key Features
1. Telegram 客户端连接与会话持久化
2. 频道配置管理（增删改查）
3. 消息过滤规则（关键词白名单/黑名单、类型、长度）
4. 图片媒体下载与上传
5. 同步状态监控与日志

### Dependencies
- gotd/td - Telegram MTProto client library

### Development Commands

```bash
# Add dependency
cd main && go get github.com/gotd/td@latest

# Run with plugin
task dev
```
<!-- SPECKIT END -->

