# Quickstart: 前端双平台部署

**Feature**: 009-dual-deploy-vercel-cloudflare
**Date**: 2026-05-10

## 快速开始

### Step 1: 创建 Cloudflare Pages 项目

1. 登录 [Cloudflare Dashboard](https://dash.cloudflare.com/)
2. 进入 **Pages** → **Create a project**
3. 选择 **Connect to Git**
4. 选择 GitHub 仓库 `ctwj/xxc`
5. 配置构建设置：
   - **Project name**: `moss-frontend`
   - **Production branch**: `main`
   - **Build command**: `cd frontend && pnpm build`
   - **Build output directory**: `frontend/.next`
   - **Node.js version**: `18`

### Step 2: 配置环境变量

在 Cloudflare Dashboard → Pages → moss-frontend → Settings → Environment variables:

```
VITE_API_URL = https://api.l9.lc
```

### Step 3: 配置自定义域名

1. 在 Pages → moss-frontend → Custom domains
2. 添加域名: `cf.l9.lc`（或你选择的域名）
3. DNS 自动配置

### Step 4: 验证部署

推送代码到 main 分支后：

```bash
# 检查 Cloudflare Pages 构建状态
# Dashboard → Pages → moss-admin → Deployments

# 访问 Cloudflare Pages URL
# https://moss-frontend.pages.dev 或 https://cf.l9.lc
```

## Vercel 配置（已存在）

确认 Vercel 项目配置正确：

1. 登录 [Vercel Dashboard](https://vercel.com/)
2. 检查项目构建设置
3. 确认环境变量配置

## 双平台验证

```bash
# 推送代码触发双平台构建
git push origin main

# 检查两个平台部署状态
# Vercel: https://vercel.com/dashboard
# Cloudflare: https://dash.cloudflare.com/pages

# 验证页面一致性
curl -s https://admin.l9.lc | head -20
curl -s https://cf.l9.lc | head -20
```

## 常见问题

### Q: Cloudflare 构建失败 "pnpm not found"

在 Cloudflare Pages 设置中添加 `packageManager: pnpm@8.x`

### Q: 两个平台页面不一致

检查环境变量配置是否相同，特别是 `VITE_API_URL`

### Q: 域名 DNS 解析失败

确认 DNS CNAME 记录配置正确：
- Cloudflare Pages: CNAME → `[project].pages.dev`
- Vercel: CNAME → `cname.vercel-dns.com`