# Quickstart: Next.js Frontend Integration

**Date**: 2026-04-29
**Feature**: 001-react-frontend-integration

## Prerequisites

- Node.js 20+
- Go 1.25+
- Vercel 账户
- Moss CMS 后端运行中

## Quick Setup (5 minutes)

### 1. 克隆项目

```bash
git clone <repository-url>
cd moss-cms
git checkout 001-react-frontend-integration
```

### 2. 启动后端

```bash
cd main
go mod tidy
go run cmd/web/main.go
```

后端运行在 `http://localhost:9008`

### 3. 启动前端开发服务器

```bash
cd frontend
npm install
npm run dev
```

前端运行在 `http://localhost:3000`

### 4. 访问应用

- 前端首页: http://localhost:3000
- 管理后台: http://localhost:9008/admin
- API 文档: 参见 [contracts/api.md](./contracts/api.md)

## Environment Variables

### 后端 (.env 或 conf.toml)

```toml
# JWT 配置
jwt_secret = "your-super-secret-key-change-in-production"
jwt_expire = "168h"  # 7天

# CORS 配置
cors_origins = "http://localhost:3000,https://your-domain.vercel.app"

# Webhook 配置
webhook_secret = "webhook-secret-key"
```

### 前端 (.env.local)

```env
# API 地址
NEXT_PUBLIC_API_URL=http://localhost:9008

# ISR Webhook Secret (用于 API Route)
REVALIDATE_SECRET=webhook-secret-key
```

### Vercel 环境变量

在 Vercel 项目设置中添加：

```
NEXT_PUBLIC_API_URL=https://api.your-domain.com
REVALIDATE_SECRET=webhook-secret-key
```

## Deployment

### 后端部署

1. 构建二进制文件：
```bash
cd main
go build -o moss cmd/web/main.go
```

2. 部署到服务器，确保 CORS 配置正确

### 前端部署 (Vercel)

1. 连接 GitHub 仓库到 Vercel

2. 配置项目：
   - Root Directory: `frontend`
   - Framework Preset: Next.js
   - Build Command: `npm run build`
   - Output Directory: `.next`

3. 添加环境变量

4. 部署

### 配置 Webhook

在 Moss CMS 管理后台配置 Webhook：

- URL: `https://your-domain.vercel.app/api/revalidate`
- Secret: 与 `REVALIDATE_SECRET` 一致
- 触发事件: 文章发布、更新、删除

## Testing

### 后端测试

```bash
cd main
go test ./...
```

### 前端测试

```bash
cd frontend
npm run test
npm run test:e2e
```

## Troubleshooting

### CORS 错误

确保后端 CORS 配置包含前端域名：

```go
// main.go
app.Use(cors.New(cors.Config{
    AllowOrigins: "http://localhost:3000,https://your-domain.vercel.app",
    AllowCredentials: true,
}))
```

### ISR 不更新

1. 检查 Webhook secret 是否一致
2. 检查 Vercel 日志是否有 revalidate 调用
3. 手动触发：`curl -X POST https://your-domain.vercel.app/api/revalidate -d '{"secret":"xxx"}'`

### JWT Token 无效

1. 检查后端和前端的 JWT secret 是否一致
2. 检查 Cookie 设置（Secure, SameSite）
3. 清除浏览器 Cookie 重新登录

## Next Steps

1. 完成功能开发后，运行 `/speckit-tasks` 生成任务列表
2. 按任务列表逐步实现功能
3. 每完成一个任务，运行测试验证
