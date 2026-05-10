# Data Model: 前端双平台部署

**Feature**: 009-dual-deploy-vercel-cloudflare
**Date**: 2026-05-10

## 概述

本功能不涉及数据模型，仅涉及部署配置文件。

## 配置文件结构

### Vercel 配置 (vercel.json)

```json
{
  "buildCommand": "pnpm build",
  "outputDirectory": "dist",
  "framework": "vite",
  "rewrites": [
    { "source": "/(.*)", "destination": "/index.html" }
  ]
}
```

### Cloudflare Pages 配置 (wrangler.toml)

```toml
name = "moss-frontend"
compatibility_date = "2024-01-01"

[build]
command = "pnpm build"

[build.environment]
NODE_VERSION = "18"
```

## 环境变量

| 变量名 | 说明 | 平台 |
|--------|------|------|
| VITE_API_URL | 后端 API 地址 | 两个平台 |
| VITE_SITE_URL | 前端站点地址 | 两个平台 |

## 域名配置

| 平台 | 域名 | DNS 提供商 |
|------|------|-----------|
| Vercel | admin.example.com | Cloudflare DNS |
| Cloudflare | cf.example.com | Cloudflare DNS |
