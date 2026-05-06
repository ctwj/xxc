# Implementation Plan: Next.js Frontend Integration (Headless CMS)

**Branch**: `001-react-frontend-integration` | **Date**: 2026-04-29 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/001-react-frontend-integration/spec.md`

## Summary

将 Moss CMS 从传统的 Jet 模板前端迁移到 Headless CMS 架构，使用 Next.js 作为前端框架，采用 ISR 渲染模式。前端独立部署于 Vercel，通过 REST API 与 Go 后端通信，使用 JWT Token 认证，Webhook 触发 ISR 按需更新。

## Technical Context

**Language/Version**:
- 后端: Go 1.25+
- 前端: TypeScript 5.9+, Node.js 20+

**Primary Dependencies**:
- 后端: Fiber (Go Web Framework), GORM, JWT (golang-jwt/jwt)
- 前端: Next.js 15+, React 18, Tailwind CSS, Radix UI, Framer Motion

**Storage**:
- 后端: SQLite (默认) / MySQL / PostgreSQL
- 前端: 无本地存储（JWT 存储在 HttpOnly Cookie）

**Testing**:
- 后端: Go standard testing (`go test`)
- 前端: Jest + React Testing Library, Playwright (E2E)

**Target Platform**:
- 后端: Linux/Windows/macOS 服务器
- 前端: Vercel (Edge Network)

**Project Type**: Web Application (前后端分离架构)

**Performance Goals**:
- 页面首次加载 < 1.5 秒
- Lighthouse 性能评分 90+
- ISR 更新延迟 < 30 秒

**Constraints**:
- 前端必须支持 SEO（ISR 静态生成）
- 跨域请求必须通过 CORS 配置
- JWT Token 必须存储在 HttpOnly Cookie

**Scale/Scope**:
- 小型内容站点
- 预计日访问量 < 10,000 PV

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

由于项目没有自定义 constitution.md（使用模板占位符），采用以下默认原则：

| 原则 | 状态 | 说明 |
|------|------|------|
| 简单优先 | ✅ Pass | 使用成熟的 Next.js + ISR 方案，避免过度设计 |
| 测试覆盖 | ⚠️ 注意 | 需要为 API 和前端组件编写测试 |
| 文档完善 | ✅ Pass | 本计划文档 + API 文档 |
| 安全性 | ✅ Pass | JWT + HttpOnly Cookie，CORS 配置 |

## Project Structure

### Documentation (this feature)

```text
specs/001-react-frontend-integration/
├── plan.md              # This file
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 output
├── quickstart.md        # Phase 1 output
├── contracts/           # Phase 1 output (API contracts)
│   └── api.md
└── tasks.md             # Phase 2 output (/speckit-tasks command)
```

### Source Code (repository root)

```text
# 现有后端结构保持不变
main/
├── api/web/
│   ├── controller/      # 现有控制器
│   ├── router/          # 路由配置
│   ├── middleware/      # 中间件 (新增 CORS, JWT)
│   └── dto/             # 数据传输对象
├── application/
│   └── service/         # 业务逻辑
├── domain/
│   └── core/            # 领域模型
├── infrastructure/
│   └── support/         # 基础设施
├── plugins/             # 插件系统
├── resources/           # 静态资源
└── startup/             # 启动配置

# 新增 Next.js 前端项目
frontend/
├── src/
│   ├── app/             # Next.js App Router
│   │   ├── page.tsx     # 首页
│   │   ├── article/
│   │   │   └── [slug]/page.tsx  # 文章详情
│   │   ├── search/page.tsx      # 搜索页
│   │   ├── favorites/page.tsx   # 收藏页
│   │   └── api/                 # API Routes
│   │       ├── auth/            # 认证相关
│   │       └── revalidate/      # ISR Webhook
│   ├── components/      # React 组件
│   │   ├── home/        # 首页组件 (来自 xxc.zip)
│   │   ├── article/     # 文章组件
│   │   └── common/      # 通用组件
│   ├── lib/             # 工具函数
│   │   ├── api.ts       # API 客户端
│   │   └── auth.ts      # 认证工具
│   ├── hooks/           # React Hooks
│   ├── contexts/        # React Context
│   └── types/           # TypeScript 类型
├── public/              # 静态资源
├── next.config.js       # Next.js 配置
├── tailwind.config.js   # Tailwind 配置
└── package.json

# 管理后台保持不变
admin/
└── [现有 Vue Admin 结构]
```

**Structure Decision**:
- 采用前后端分离架构
- 后端保持现有 Go 项目结构
- 新增 `frontend/` 目录存放 Next.js 项目
- 管理后台 `admin/` 保持不变

## Complexity Tracking

> 无重大架构违规需要记录

## Implementation Phases

### Phase 1: 后端 API 扩展 (Week 1)

**目标**: 完善 REST API，支持 JWT 认证，配置 CORS

**任务**:
1. 扩展现有 API，添加公开接口
   - `GET /api/articles` - 文章列表
   - `GET /api/articles/:slug` - 文章详情
   - `GET /api/categories` - 分类列表
   - `GET /api/tags` - 标签列表
   - `GET /api/search` - 搜索

2. 实现 JWT 认证
   - `POST /api/auth/login` - 登录
   - `POST /api/auth/register` - 注册
   - `GET /api/auth/me` - 当前用户
   - JWT 中间件

3. CORS 配置
   - 允许 Next.js 域名跨域访问

4. Webhook 端点
   - `POST /api/webhook/revalidate` - 触发 ISR 更新

### Phase 2: Next.js 项目初始化 (Week 1-2)

**目标**: 创建 Next.js 项目，迁移 xxc.zip 组件

**任务**:
1. 初始化 Next.js 15 项目
   - App Router
   - TypeScript
   - Tailwind CSS

2. 迁移 xxc.zip 组件
   - 适配 Next.js App Router
   - 移除 Supabase 依赖
   - 替换为 Moss CMS API 调用

3. 实现 ISR 页面
   - 首页静态生成
   - 文章详情页静态生成 + 按需更新

### Phase 3: 认证与收藏功能 (Week 2-3)

**目标**: 实现 JWT 认证和收藏功能

**任务**:
1. 认证流程
   - 登录/注册页面
   - JWT Token 管理
   - 受保护路由

2. 收藏功能
   - 收藏 API
   - 收藏页面
   - 收藏状态同步

### Phase 4: Webhook 与 ISR 集成 (Week 3)

**目标**: 实现发布文章后自动更新前端

**任务**:
1. 后台集成 Webhook
   - 发布文章时触发
   - 调用 Vercel API

2. Next.js API Route
   - 接收 Webhook
   - 调用 `revalidatePath` / `revalidateTag`

### Phase 5: 测试与部署 (Week 3-4)

**目标**: 完成测试，部署上线

**任务**:
1. 单元测试
   - API 测试
   - 组件测试

2. E2E 测试
   - 关键用户流程

3. 部署配置
   - Vercel 项目配置
   - 环境变量设置
   - 域名配置

## Risk Mitigation

| 风险 | 影响 | 缓解措施 |
|------|------|----------|
| ISR 更新延迟 | 中 | Webhook 失败时提供手动刷新按钮 |
| API 超时 | 中 | Next.js 实现缓存降级机制 |
| JWT 安全 | 高 | HttpOnly Cookie + CSRF 保护 |
| SEO 回归 | 高 | 部署前验证 Lighthouse 评分 |
