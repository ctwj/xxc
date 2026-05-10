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
## Active Feature: 修复 MySQL 数据库编码问题

**Branch**: `008-migrate-to-postgresql`
**Plan**: [specs/008-migrate-to-postgresql/plan.md](specs/008-migrate-to-postgresql/plan.md)
**Status**: Planning Complete

### Key Documents
- [Specification](specs/008-migrate-to-postgresql/spec.md)
- [Research](specs/008-migrate-to-postgresql/research.md)
- [Data Model](specs/008-migrate-to-postgresql/data-model.md)
- [API Contracts](specs/008-migrate-to-postgresql/contracts/api.md)
- [Quickstart](specs/008-migrate-to-postgresql/quickstart.md)

### Problem
创建数据库时没有指定编码，GORM migrate 后不同表使用了不同的字符集/排序规则，导致编码错误。

### Root Cause
1. 创建数据库时未指定默认字符集和排序规则
2. MySQL 服务器、数据库、表三层级字符集配置不一致
3. GORM 实体定义中 TEXT 类型字段设置了无效的默认值

### Solution
1. 配置 MySQL 服务器默认字符集（预防机制）
2. 应用程序 DSN 指定连接字符集
3. 批量转换现有数据库和表
4. 修复 GORM 实体定义
<!-- SPECKIT END -->

