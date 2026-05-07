# Feature Specification: React Frontend Integration Analysis

**Feature Branch**: `feat/react-frontend-integration`
**Created**: 2026-04-29
**Status**: Draft
**Input**: User description: "当前系统使用 模版的方式 显示的前端，但是UI效果不太好，我现在获得了一个 react 的前端，效果更好，请帮我分析，是否可以使用当前提供的 代码作为新的前端，有哪些方案，推荐什么方案"

## Clarifications

### Session 2026-04-29

- Q: Next.js 渲染模式选择？ → A: ISR (Incremental Static Regeneration) - 静态页面 + 定时重新验证，适合"频繁发布新文章、内容稳定"的场景
- Q: Next.js 部署架构？ → A: Next.js 独立部署 (Vercel/Netlify) - 前后端完全分离，最佳性能
- Q: 认证机制集成？ → A: JWT Token 认证 - 无状态，适合前后端分离架构
- Q: ISR 重新验证触发方式？ → A: Webhook 按需触发 - 发布文章时主动通知 Next.js 重新生成，实时性最佳
- Q: API 设计风格？ → A: REST API - 复用现有 API 基础，对 ISR 缓存友好

## Executive Summary

本文档分析将提供的 React 前端（xxc.zip）集成到现有 Moss CMS 系统的可行性方案。经过详细分析，我们识别出三种主要集成方案。

**最终选择**: 方案三 - Headless CMS 模式，使用 Next.js 作为前端框架，采用 ISR 渲染模式。

## Current System Analysis

### Moss CMS 现有前端架构

**技术栈**:
- 模板引擎: Jet Template Engine (`main/infrastructure/support/template/engine/jet.go`)
- 模板位置: `main/resources/themes/{theme_name}/template/`
- 渲染方式: 服务端渲染 (SSR)，通过 `RenderService` 渲染 HTML 模板
- 路由: Fiber 框架处理请求，返回渲染后的 HTML

**现有功能**:
- 首页展示 (`template/index.html`)
- 文章详情页 (`template/article.html`)
- 分类页面 (`template/category.html`)
- 标签页面 (`template/tag.html`)
- 搜索功能 (`template/search.html`)
- 主题系统支持多主题切换

### React 前端分析 (xxc.zip)

**技术栈**:
- 框架: React 18 + TypeScript
- 构建工具: Vite (rolldown-vite)
- UI 组件: Radix UI + Tailwind CSS + shadcn/ui
- 动画: Framer Motion
- 路由: React Router 7
- 认证: miaoda-auth-react + Supabase
- 存储: Supabase

**页面结构**:
- 首页 (`HomePage.tsx`): 卡片堆叠式信息流展示
- 登录页 (`LoginPage.tsx`): 用户登录
- 注册页 (`RegisterPage.tsx`): 用户注册
- 收藏页 (`FavoritePage.tsx`): 用户收藏内容

**数据类型** (`types/index.ts`):
- `InfoCardData`: 内容卡片数据结构
- 支持: text, image, video, image-text, video-text, images, images-text, images-video-text, long-text

**特点**:
- 现代化 UI 设计，卡片堆叠交互
- 暗色/亮色主题切换
- 使用 Mock 数据 (`mockData.ts`) 和 localStorage 认证
- 依赖 Supabase 作为后端服务

## Integration Options Analysis

### 方案一: 完全替换 (Full Replacement)

**描述**: 使用 React 前端完全替换现有的 Jet 模板系统

**实现方式**:
1. React 应用独立部署或嵌入 Go 二进制
2. 后端提供 REST API 供 React 调用
3. 移除 Jet 模板渲染逻辑

**优点**:
- 现代化用户体验，交互流畅
- 组件化开发，易于维护
- 前后端分离，独立部署

**缺点**:
- 需要重构后端 API（当前部分 API 已存在）
- 失去 SEO 优势（SSR → CSR）
- 增加部署复杂度
- 需要处理认证系统集成

**工作量**: 高 (约 4-6 周)

### 方案二: 混合模式 (Hybrid Mode) - 推荐

**描述**: React 前端作为新主题/模块，与现有模板系统共存

**实现方式**:
1. 将 React 构建产物嵌入 `main/resources/` 目录
2. 新增路由 `/app/*` 指向 React 应用
3. 保留现有模板系统用于 SEO 关键页面
4. 共享后端 API 和认证系统

**架构图**:
```
用户请求
    ├── / (首页) → Jet 模板渲染 (SEO)
    ├── /article/* → Jet 模板渲染 (SEO)
    ├── /app/* → React SPA (交互体验)
    └── /admin/* → Vue Admin (管理后台)
```

**优点**:
- 保留 SEO 优势
- 渐进式迁移，风险可控
- 用户可选择传统页面或现代化界面
- 复用现有 API 和认证

**缺点**:
- 需要维护两套前端
- 路由和状态管理复杂度增加

**工作量**: 中 (约 2-3 周)

### 方案三: API 模式 (Headless CMS)

**描述**: Moss CMS 作为纯 API 后端，React 作为独立前端应用

**实现方式**:
1. 完善 REST API 覆盖所有数据操作
2. React 应用独立部署（如 Vercel、Netlify）
3. 后端仅提供 API 和管理后台

**优点**:
- 完全解耦，灵活部署
- 可支持多前端（Web、Mobile、小程序）
- 符合现代 Headless CMS 架构

**缺点**:
- SEO 需要额外处理（预渲染/SSR）
- 部署架构变更
- 需要处理跨域、认证等问题

**工作量**: 高 (约 4-5 周)

## Recommended Solution

**选择方案**: 方案三 - Headless CMS 模式

**技术选型**:
- 前端框架: Next.js (替代原 React + Vite)
- 渲染模式: ISR (Incremental Static Regeneration)
- 部署方式: Next.js 独立部署于 Vercel/Netlify
- 认证机制: JWT Token
- API 风格: REST API
- ISR 触发: Webhook 按需触发

**架构图**:
```
┌─────────────────────────────────────────────────────────────┐
│                        用户请求                              │
└─────────────────────────────────────────────────────────────┘
                              │
              ┌───────────────┼───────────────┐
              ▼               ▼               ▼
     ┌─────────────┐  ┌─────────────┐  ┌─────────────┐
     │  Next.js    │  │  Go API     │  │  Vue Admin  │
     │  (Vercel)   │  │  (后端)     │  │  (管理后台) │
     │  ISR 页面   │  │  REST API   │  │             │
     └─────────────┘  └─────────────┘  └─────────────┘
              │               │               │
              └───────────────┴───────────────┘
                              │
                    ┌─────────▼─────────┐
                    │   Moss CMS 数据库  │
                    │   (SQLite/MySQL)  │
                    └───────────────────┘

发布流程:
后台发布文章 → Webhook → Vercel ISR 重新生成
```

**理由**:
1. **SEO 完美支持**: ISR 在构建时生成完整 HTML，搜索引擎直接抓取
2. **性能最优**: 静态页面 + CDN 分发，首屏加载极快
3. **实时更新**: Webhook 触发 ISR，新文章立即生效
4. **架构解耦**: 前后端独立部署，灵活扩展
5. **成本可控**: Vercel 免费额度足够小型站点使用

## User Scenarios & Testing *(mandatory)*

### User Story 1 - 访问 Next.js 前端 (Priority: P1)

用户可以通过域名访问 Next.js 前端界面，获得现代化的交互体验和快速的页面加载。

**Why this priority**: 这是核心功能入口，决定了用户是否能使用新界面。

**Independent Test**: 访问站点首页，验证 Next.js ISR 页面正常加载，显示文章卡片列表。

**Acceptance Scenarios**:

1. **Given** 用户访问站点首页, **When** 页面加载完成, **Then** Next.js 预渲染的静态页面正常显示，包含文章列表
2. **Given** 用户在首页, **When** 点击文章卡片, **Then** 跳转到文章详情页，内容完整显示
3. **Given** 用户在任意页面, **When** 切换主题, **Then** 暗色/亮色主题正确切换

---

### User Story 2 - 数据集成 (Priority: P1)

Next.js 前端能够从 Moss CMS 后端 API 获取真实数据。

**Why this priority**: 没有真实数据，前端只是展示壳，无实际价值。

**Independent Test**: Next.js 首页显示来自 Moss CMS 的真实文章数据。

**Acceptance Scenarios**:

1. **Given** 后端有文章数据, **When** Next.js 在构建时调用 `/api/articles`, **Then** 返回正确的文章列表并生成静态页面
2. **Given** 用户查看文章详情, **When** 访问 `/article/{slug}`, **Then** 显示正确的文章内容
3. **Given** 用户搜索内容, **When** 提交搜索关键词, **Then** 返回匹配的搜索结果

---

### User Story 3 - ISR 实时更新 (Priority: P1)

管理员在后台发布新文章后，Next.js 前端能够自动更新。

**Why this priority**: ISR 的核心价值在于按需更新，这是 Headless CMS 模式的关键能力。

**Independent Test**: 在后台发布新文章后，前端首页能够显示新文章。

**Acceptance Scenarios**:

1. **Given** 管理员在后台发布新文章, **When** 点击发布, **Then** Webhook 触发 Vercel ISR 重新生成
2. **Given** ISR 重新生成完成, **When** 用户刷新首页, **Then** 新文章出现在列表中
3. **Given** 文章被更新, **When** Webhook 触发, **Then** 文章详情页内容同步更新

---

### User Story 4 - 认证集成 (Priority: P2)

用户可以使用 Moss CMS 现有账户登录 Next.js 前端。

**Why this priority**: 认证是收藏、个性化等功能的基础，但可以延后实现。

**Independent Test**: 用户能够使用现有账户登录，并访问收藏功能。

**Acceptance Scenarios**:

1. **Given** 用户有 Moss CMS 账户, **When** 在 Next.js 应用中登录, **Then** JWT 认证成功，显示用户信息
2. **Given** 已登录用户, **When** 收藏文章, **Then** 收藏状态正确保存到后端
3. **Given** 已登录用户, **When** 访问收藏页, **Then** 显示已收藏的文章列表

---

### User Story 5 - 管理后台保留 (Priority: P2)

管理员能够继续使用现有的 Vue 管理后台进行内容管理。

**Why this priority**: 管理后台是内容运营的基础，必须保持可用。

**Independent Test**: 管理员能够登录后台，发布、编辑、删除文章。

**Acceptance Scenarios**:

1. **Given** 管理员访问 `/admin/`, **When** 登录成功, **Then** 进入管理后台
2. **Given** 管理员在后台发布文章, **When** 点击发布, **Then** 文章保存成功并触发 Webhook

### Edge Cases

- 当 Vercel 部署失败时，是否有回退机制？
- 当 API 请求超时时，Next.js 如何优雅降级（显示缓存内容或错误页面）？
- Webhook 触发失败时，如何确保 ISR 最终一致性？
- JWT Token 过期时，如何无感刷新？
- 跨域请求如何处理（Next.js API Route 代理）？

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: Next.js 前端必须部署在 Vercel/Netlify，使用 ISR 渲染模式
- **FR-002**: Next.js 必须在构建时从 Moss CMS API 获取文章列表数据
- **FR-003**: Next.js 必须支持文章详情页的静态生成和按需更新
- **FR-004**: Next.js 必须支持搜索功能，调用后端搜索 API
- **FR-005**: Moss CMS 后端必须提供完整的 REST API 供 Next.js 调用
- **FR-006**: Next.js 前端必须支持暗色/亮色主题切换
- **FR-007**: 用户必须能够使用 JWT Token 认证登录 Next.js 前端
- **FR-008**: 已登录用户必须能够收藏文章
- **FR-009**: 后台发布文章时必须通过 Webhook 触发 Next.js ISR 更新
- **FR-010**: Vue 管理后台必须保持可用，用于内容管理
- **FR-011**: Go 后端必须支持 CORS 配置，允许 Next.js 跨域请求

### Key Entities

- **Article**: 文章实体，包含标题、内容、摘要、缩略图、分类、标签、浏览量等
- **Category**: 分类实体，文章的归属分类
- **Tag**: 标签实体，文章的关联标签
- **User**: 用户实体，用于认证和收藏功能
- **Favorite**: 收藏关系，用户与文章的收藏关联

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Next.js 页面首次加载时间小于 1.5 秒（ISR 静态页面 + CDN）
- **SC-002**: Lighthouse 性能评分达到 90+ 分
- **SC-003**: SEO 完全支持，所有页面可被搜索引擎正确索引
- **SC-004**: ISR 更新延迟小于 30 秒（从发布到前端可见）
- **SC-005**: 用户能够在 3 分钟内完成从访问到收藏文章的完整流程
- **SC-006**: 管理后台发布文章后，Webhook 触发成功率 99%+

## Assumptions

- 用户浏览器支持现代 JavaScript 特性（ES2020+）
- Moss CMS 后端 API 可以扩展支持完整的 REST API
- Vercel 免费额度足够站点使用（或愿意付费）
- 现有认证系统可以扩展支持 JWT Token
- 管理员愿意使用 Webhook 机制触发 ISR 更新
- Next.js 前端将完全替代现有的 Jet 模板前端（传统模板将被移除）
