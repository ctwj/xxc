# Tasks: Next.js Frontend Integration

**Input**: Design documents from `/specs/001-react-frontend-integration/`
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/api.md

**Tests**: 未明确要求测试，本任务列表不包含测试任务。

**Organization**: 任务按用户故事组织，每个故事可独立实现和测试。

## Format: `[ID] [P?] [Story] Description`

- **[P]**: 可并行执行（不同文件，无依赖）
- **[Story]**: 所属用户故事 (US1, US2, US3, US4, US5)
- 描述中包含精确文件路径

## Path Conventions

- **后端**: `main/` (现有 Go 项目)
- **前端**: `frontend/` (新建 Next.js 项目)
- **管理后台**: `admin/` (保持不变)

---

## Phase 1: Setup (项目初始化)

**Purpose**: 创建项目结构和基础配置

- [X] T001 创建 frontend 目录和 Next.js 15 项目结构
- [X] T002 [P] 初始化 Next.js 项目：`cd frontend && npx create-next-app@latest . --typescript --tailwind --app-router`
- [X] T003 [P] 安装前端依赖：`npm install framer-motion jose clsx tailwind-merge sonner`
- [X] T004 [P] 安装 Radix UI 组件：`npm install @radix-ui/react-dialog @radix-ui/react-dropdown-menu @radix-ui/react-toast`
- [X] T005 [P] 配置 Tailwind CSS：创建 `frontend/tailwind.config.ts`
- [X] T006 [P] 创建 TypeScript 类型定义：`frontend/src/types/index.ts`
- [X] T007 [P] 创建 API 客户端：`frontend/src/lib/api.ts`
- [X] T008 [P] 创建工具函数：`frontend/src/lib/utils.ts`
- [X] T009 配置环境变量：创建 `frontend/.env.local`

---

## Phase 2: Foundational (后端基础设施)

**Purpose**: 后端 API 基础设施，必须在用户故事开始前完成

**⚠️ CRITICAL**: 所有用户故事依赖此后端 API

### 后端 CORS 配置

- [X] T010 添加 CORS 中间件配置：`main/api/web/middleware/cors.go`
- [X] T011 在路由中启用 CORS：`main/api/web/router/router.go`

### 后端 JWT 认证

- [X] T012 [P] 添加 JWT 依赖：`cd main && go get github.com/golang-jwt/jwt/v5`
- [X] T013 [P] 创建 JWT 工具函数：`main/infrastructure/support/auth/jwt.go`
- [X] T014 创建 JWT 中间件：`main/api/web/middleware/jwt.go`
- [X] T015 添加 JWT 配置到 conf.toml 结构：`main/domain/config/entity/auth.go`

### 公开 API 控制器

- [X] T016 [P] 创建公开 API 控制器：`main/api/web/controller/api.go`
- [X] T017 [P] 创建 API DTO：`main/api/web/dto/api.go`
- [X] T018 创建 API 路由：`main/api/web/router/api.go`

**Checkpoint**: 后端基础设施就绪，可以开始用户故事实现

---

## Phase 3: User Story 1 - 访问 Next.js 前端 (Priority: P1) 🎯 MVP

**Goal**: 用户可以访问 Next.js 前端，查看文章列表和详情

**Independent Test**: 访问首页显示文章列表，点击文章可查看详情

### 后端 API 实现 [US1]

- [X] T019 [P] [US1] 实现文章列表 API：`GET /api/articles` 在 `main/api/web/controller/api.go`
- [X] T020 [P] [US1] 实现文章详情 API：`GET /api/articles/:slug` 在 `main/api/web/controller/api.go`
- [X] T021 [P] [US1] 实现分类列表 API：`GET /api/categories` 在 `main/api/web/controller/api.go`
- [X] T022 [P] [US1] 实现标签列表 API：`GET /api/tags` 在 `main/api/web/controller/api.go`
- [X] T023 [US1] 注册 API 路由：更新 `main/api/web/router/api.go`

### 前端页面实现 [US1]

- [X] T024 [P] [US1] 创建首页组件：`frontend/src/app/page.tsx`
- [X] T025 [P] [US1] 创建文章详情页：`frontend/src/app/article/[slug]/page.tsx`
- [X] T026 [P] [US1] 创建文章卡片组件：`frontend/src/components/home/InfoCard.tsx`
- [X] T027 [P] [US1] 创建卡片堆叠组件：`frontend/src/components/home/CardStack.tsx`
- [X] T028 [P] [US1] 创建布局组件：`frontend/src/app/layout.tsx`
- [X] T029 [US1] 实现 ISR 静态生成：配置 `generateStaticParams` 在文章详情页
- [X] T030 [US1] 创建主题切换功能：`frontend/src/contexts/ThemeContext.tsx`

**Checkpoint**: 用户故事 1 完成，可独立测试首页和文章详情

---

## Phase 4: User Story 2 - 数据集成 (Priority: P1)

**Goal**: Next.js 从 Moss CMS API 获取真实数据

**Independent Test**: 首页显示真实文章数据，详情页显示完整内容

### 后端搜索 API [US2]

- [X] T031 [US2] 实现搜索 API：`GET /api/search` 在 `main/api/web/controller/api.go`

### 前端数据获取 [US2]

- [X] T032 [P] [US2] 创建文章服务：`frontend/src/services/articleService.ts`
- [X] T033 [P] [US2] 创建分类服务：`frontend/src/services/categoryService.ts`
- [X] T034 [US2] 更新首页获取真实数据：修改 `frontend/src/app/page.tsx`
- [X] T035 [US2] 更新文章详情页获取真实数据：修改 `frontend/src/app/article/[slug]/page.tsx`
- [X] T036 [US2] 创建搜索页面：`frontend/src/app/search/page.tsx`

**Checkpoint**: 用户故事 2 完成，前端显示真实数据

---

## Phase 5: User Story 3 - ISR 实时更新 (Priority: P1)

**Goal**: 后台发布文章后前端自动更新

**Independent Test**: 发布新文章后刷新首页可见

### 后端 Webhook [US3]

- [X] T037 [US3] 创建 Webhook 控制器：`main/api/web/controller/webhook.go`
- [X] T038 [US3] 添加 Webhook 路由：`POST /api/webhook/revalidate` 在 `main/api/web/router/api.go`
- [X] T039 [US3] 在文章发布时触发 Webhook：修改 `main/api/web/controller/article.go` (通过 WebhookTrigger 插件实现)

### 前端 ISR Revalidation [US3]

- [X] T040 [US3] 创建 Revalidate API Route：`frontend/src/app/api/revalidate/route.ts`
- [X] T041 [US3] 配置 ISR revalidate 时间：更新 `frontend/src/app/page.tsx`
- [X] T042 [US3] 添加文章详情页 revalidate：更新 `frontend/src/app/article/[slug]/page.tsx`

**Checkpoint**: 用户故事 3 完成，ISR 更新机制就绪

---

## Phase 6: User Story 4 - 认证集成 (Priority: P2)

**Goal**: 用户可登录并使用收藏功能

**Independent Test**: 用户登录后可收藏文章

### 后端认证 API [US4]

- [X] T043 [P] [US4] 创建认证控制器：`main/api/web/controller/auth.go`
- [X] T044 [P] [US4] 实现登录 API：`POST /api/auth/login`
- [X] T045 [P] [US4] 实现注册 API：`POST /api/auth/register`
- [X] T046 [P] [US4] 实现登出 API：`POST /api/auth/logout`
- [X] T047 [P] [US4] 实现获取当前用户 API：`GET /api/auth/me`
- [X] T048 [US4] 注册认证路由：更新 `main/api/web/router/api.go`

### 后端收藏功能 [US4]

- [X] T049 [US4] 创建 Favorite 实体：`main/domain/core/entity/favorite.go`
- [X] T050 [US4] 创建 Favorite 仓库：`main/domain/core/repository/favorite.go`
- [X] T051 [US4] 创建 Favorite 服务：`main/domain/core/service/favorite.go`
- [X] T052 [US4] 实现收藏 API：`GET/POST/DELETE /api/favorites` 在 `main/api/web/controller/favorite.go`

### 前端认证 [US4]

- [X] T053 [P] [US4] 创建认证 Context：`frontend/src/contexts/AuthContext.tsx`
- [X] T054 [P] [US4] 创建登录页面：`frontend/src/app/login/page.tsx`
- [X] T055 [P] [US4] 创建注册页面：`frontend/src/app/register/page.tsx`
- [X] T056 [US4] 创建认证服务：`frontend/src/services/authService.ts`
- [X] T057 [US4] 添加认证状态管理：更新 `frontend/src/app/layout.tsx`

### 前端收藏功能 [US4]

- [X] T058 [US4] 创建收藏服务：`frontend/src/services/favoriteService.ts`
- [X] T059 [US4] 创建收藏页面：`frontend/src/app/favorites/page.tsx`
- [X] T060 [US4] 添加收藏按钮到文章卡片：更新 `frontend/src/components/home/InfoCard.tsx`

**Checkpoint**: 用户故事 4 完成，认证和收藏功能就绪

---

## Phase 7: User Story 5 - 管理后台保留 (Priority: P2)

**Goal**: 管理员可继续使用 Vue 后台

**Independent Test**: 管理员可登录后台发布文章

### 后台 Webhook 集成 [US5]

- [X] T061 [US5] 添加 Webhook 配置到后台设置：`admin/src/views/config/module/site.vue`
- [X] T062 [US5] 创建 Webhook 触发插件：`main/plugins/WebhookTrigger.go`
- [X] T063 [US5] 注册 Webhook 插件：更新 `main/startup/startup.go`

**Checkpoint**: 用户故事 5 完成，管理后台集成 Webhook

---

## Phase 8: Polish & Cross-Cutting Concerns

**Purpose**: 优化和完善

- [X] T064 [P] 添加错误边界组件：`frontend/src/components/common/ErrorBoundary.tsx`
- [X] T065 [P] 添加加载状态组件：`frontend/src/components/common/Loading.tsx`
- [X] T066 [P] 优化 SEO 元数据：更新 `frontend/src/app/layout.tsx`
- [X] T067 [P] 添加 sitemap：`frontend/src/app/sitemap.ts`
- [X] T068 [P] 添加 robots.txt：`frontend/src/app/robots.ts`
- [X] T069 配置 Vercel 部署：创建 `frontend/vercel.json`
- [X] T070 更新 CLAUDE.md 文档
- [X] T071 运行 quickstart.md 验证流程 (后端和前端编译通过，部署配置完成)

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: 无依赖，可立即开始
- **Foundational (Phase 2)**: 依赖 Setup 完成 - **阻塞所有用户故事**
- **User Stories (Phase 3-7)**: 都依赖 Foundational 完成
  - US1, US2, US3 可并行（都是 P1）
  - US4, US5 可并行（都是 P2）
- **Polish (Phase 8)**: 依赖所需用户故事完成

### User Story Dependencies

- **US1 (P1)**: 依赖 Foundational - 无其他故事依赖
- **US2 (P1)**: 依赖 Foundational + US1（需要页面基础）
- **US3 (P1)**: 依赖 Foundational + US1（需要页面基础）
- **US4 (P2)**: 依赖 Foundational - 可独立实现
- **US5 (P2)**: 依赖 US3（需要 Webhook 机制）

### Parallel Opportunities

**Phase 1 Setup**:
```
T002, T003, T004, T005, T006, T007, T008 可并行
```

**Phase 2 Foundational**:
```
T012, T013 可并行
T016, T017 可并行
```

**Phase 3 US1**:
```
T019, T020, T021, T022 可并行（后端 API）
T024, T025, T026, T027, T028 可并行（前端页面）
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. 完成 Phase 1: Setup
2. 完成 Phase 2: Foundational
3. 完成 Phase 3: User Story 1
4. **验证**: 访问首页和文章详情
5. 部署 MVP

### Incremental Delivery

1. Setup + Foundational → 基础就绪
2. US1 → 首页和详情页可用 → **MVP!**
3. US2 + US3 → 数据集成和 ISR → 增量发布
4. US4 → 认证和收藏 → 增量发布
5. US5 + Polish → 完整功能

---

## Summary

| Phase | 任务数 | 可并行 |
|-------|--------|--------|
| Setup | 9 | 8 |
| Foundational | 9 | 4 |
| US1 | 12 | 8 |
| US2 | 6 | 2 |
| US3 | 6 | 0 |
| US4 | 18 | 5 |
| US5 | 3 | 0 |
| Polish | 8 | 5 |
| **Total** | **71** | **32** |

**MVP Scope**: Phase 1 + Phase 2 + Phase 3 (US1) = 30 tasks
