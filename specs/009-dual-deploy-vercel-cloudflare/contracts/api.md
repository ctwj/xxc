# Deployment Contract: 双平台部署配置

**Feature**: 009-dual-deploy-vercel-cloudflare
**Date**: 2026-05-10

## Cloudflare Pages 配置

### 项目设置

| 设置项 | 值 |
|--------|-----|
| 项目名称 | moss-frontend |
| 生产分支 | main |
| 构建命令 | pnpm build |
| 输出目录 | dist |
| Node.js 版本 | 18 |

### GitHub 集成

1. 在 Cloudflare Dashboard → Pages → Create a project
2. 选择 "Connect to Git"
3. 选择 GitHub 仓库
4. 配置构建设置

### 环境变量配置

在 Cloudflare Dashboard → Pages → [Project] → Settings → Environment variables:

```
VITE_API_URL = https://api.example.com
VITE_SITE_URL = https://cf.example.com
```

## Vercel 配置

### 项目设置（已存在）

| 设置项 | 值 |
|--------|-----|
| 项目名称 | moss-frontend |
| 生产分支 | main |
| 构建命令 | pnpm build |
| 输出目录 | dist |

### 环境变量配置

在 Vercel Dashboard → [Project] → Settings → Environment Variables:

```
VITE_API_URL = https://api.example.com
VITE_SITE_URL = https://admin.example.com
```

## 域名配置

### Cloudflare Pages 域名

1. 在 Cloudflare Dashboard → Pages → [Project] → Custom domains
2. 添加自定义域名: `cf.example.com`
3. DNS 自动配置（CNAME 记录）

### Vercel 域名

1. 在 Vercel Dashboard → [Project] → Settings → Domains
2. 添加自定义域名: `admin.example.com`
3. 在 Cloudflare DNS 添加 CNAME 记录指向 Vercel