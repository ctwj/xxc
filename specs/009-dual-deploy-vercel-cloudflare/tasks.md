# Tasks: 前端双平台部署 (Vercel + Cloudflare)

**Input**: Design documents from `/specs/009-dual-deploy-vercel-cloudflare/`
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/api.md, quickstart.md

**Tests**: 无测试要求（部署配置任务）

**Organization**: 任务按用户故事分组，支持独立实现和验证。

## Format: `[ID] [P?] [Story] Description`

- **[P]**: 可并行执行（不同文件，无依赖）
- **[Story]**: 所属用户故事（US1, US2, US3, US4）
- 包含精确文件路径

---

## Phase 1: Setup (项目初始化)

**Purpose**: 确认现有配置和依赖

- [x] T001 确认 Vercel 项目配置正常，检查 frontend/vercel.json 存在
- [x] T002 确认前端构建配置正常，检查 frontend/next.config.ts
- [x] T003 [P] 确认 package.json 构建脚本配置正确

---

## Phase 2: Foundational (基础配置)

**Purpose**: 创建 Cloudflare Pages 配置文件

**⚠️ CRITICAL**: 必须完成此阶段才能进行用户故事实现

- [x] T004 创建 GitHub Actions 部署脚本 .github/workflows/deploy-frontend-cloudflare.yml
- [x] T005 更新 .gitignore 添加 Cloudflare 相关忽略项

**Checkpoint**: Cloudflare 配置文件就绪，可以开始 Cloudflare Pages 项目创建

---

## Phase 3: User Story 1 - 配置 Cloudflare Pages 部署 (Priority: P1) 🎯 MVP

**Goal**: 配置 Cloudflare Pages 部署，使前端可以部署到 Cloudflare

**Independent Test**: 推送代码后，Cloudflare Pages 自动构建并部署成功

### Implementation for User Story 1

- [ ] T006 [US1] 在 Cloudflare Dashboard 创建 Pages 项目（手动操作 - Git 集成）
- [ ] T007 [US1] 连接 GitHub 仓库 ctwj/xxc 到 Cloudflare Pages（手动操作）
- [ ] T008 [US1] 配置构建设置：Build command = `cd frontend && pnpm build`，Output directory = `frontend/.next`
- [ ] T009 [US1] 配置 Node.js 版本为 18
- [ ] T010 [US1] 触发首次部署并验证成功

**手动操作指南**: 参考 [quickstart.md](./quickstart.md) Step 1-2

**Checkpoint**: Cloudflare Pages 项目创建完成，可以自动部署

---

## Phase 4: User Story 2 - 配置双平台 CI/CD (Priority: P1)

**Goal**: 配置双平台的 CI/CD 流水线，代码推送后两个平台都能自动部署

**Independent Test**: 推送代码后，Vercel 和 Cloudflare 都成功部署

### Implementation for User Story 2

- [ ] T011 [US2] 验证 Vercel 自动部署配置正常
- [ ] T012 [US2] 验证 Cloudflare Pages 自动部署配置正常
- [ ] T013 [US2] 推送测试代码到 main 分支触发双平台构建
- [ ] T014 [US2] 验证两个平台部署成功且页面内容一致

**Checkpoint**: 双平台 CI/CD 配置完成，推送代码自动触发双平台部署

---

## Phase 5: User Story 3 - 配置自定义域名 (Priority: P2)

**Goal**: 为两个平台配置不同的自定义域名

**Independent Test**: 两个域名都能正常访问前端应用

### Implementation for User Story 3

- [ ] T015 [US3] 在 Cloudflare Pages 配置自定义域名（如 cf.l9.lc）
- [ ] T016 [US3] 验证 Cloudflare 域名 DNS 配置正确
- [ ] T017 [US3] 验证 Vercel 域名配置正常（如 admin.l9.lc）
- [ ] T018 [US3] 验证两个域名 SSL 证书有效

**Checkpoint**: 双平台域名配置完成，可以通过不同域名访问

---

## Phase 6: User Story 4 - 环境变量同步 (Priority: P2)

**Goal**: 确保两个平台的环境变量配置一致

**Independent Test**: 两个平台的应用行为一致

### Implementation for User Story 4

- [ ] T019 [US4] 在 Cloudflare Dashboard 配置环境变量 VITE_API_URL
- [ ] T020 [US4] 验证 Vercel 环境变量配置与 Cloudflare 一致
- [ ] T021 [US4] 触发重新部署验证环境变量生效
- [ ] T022 [US4] 验证两个平台应用行为一致

**Checkpoint**: 环境变量同步完成，双平台应用行为一致

---

## Phase 7: Polish & Cross-Cutting Concerns

**Purpose**: 文档更新和最终验证

- [ ] T023 [P] 更新 CLAUDE.md 添加双平台部署说明
- [ ] T024 [P] 更新 quickstart.md 添加实际域名和配置信息
- [ ] T025 运行 quickstart.md 验证流程完整性
- [ ] T026 提交所有配置更改到 Git

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: 无依赖，可立即开始
- **Foundational (Phase 2)**: 依赖 Setup 完成
- **User Story 1 (Phase 3)**: 依赖 Foundational 完成 - Cloudflare 配置文件就绪
- **User Story 2 (Phase 4)**: 依赖 User Story 1 完成 - 需要 Cloudflare 项目存在
- **User Story 3 (Phase 5)**: 依赖 User Story 2 完成 - 需要双平台部署正常
- **User Story 4 (Phase 6)**: 依赖 User Story 2 完成 - 需要双平台部署正常
- **Polish (Phase 7)**: 依赖所有用户故事完成

### User Story Dependencies

- **User Story 1 (P1)**: 可在 Foundational 完成后开始 - 无其他故事依赖
- **User Story 2 (P1)**: 依赖 US1 完成 - 需要 Cloudflare 项目存在
- **User Story 3 (P2)**: 依赖 US2 完成 - 需要双平台部署正常
- **User Story 4 (P2)**: 依赖 US2 完成 - 可与 US3 并行

### Parallel Opportunities

- Phase 1 所有任务可并行执行
- Phase 5 和 Phase 6 可并行执行（US3 和 US4）
- Phase 7 中 T023 和 T024 可并行执行

---

## Parallel Example: Phase 5 & 6

```bash
# 可以同时进行:
Task: "在 Cloudflare Pages 配置自定义域名" (US3)
Task: "在 Cloudflare Dashboard 配置环境变量" (US4)
```

---

## Implementation Strategy

### MVP First (User Story 1 & 2)

1. 完成 Phase 1: Setup
2. 完成 Phase 2: Foundational
3. 完成 Phase 3: User Story 1 (Cloudflare Pages 项目创建)
4. 完成 Phase 4: User Story 2 (双平台 CI/CD)
5. **STOP and VALIDATE**: 推送代码验证双平台自动部署

### Incremental Delivery

1. Setup + Foundational → 配置文件就绪
2. User Story 1 → Cloudflare Pages 项目创建 → 验证部署
3. User Story 2 → 双平台 CI/CD → 验证自动部署
4. User Story 3 → 自定义域名 → 验证域名访问
5. User Story 4 → 环境变量同步 → 验证行为一致

---

## Notes

- **手动操作**: T006-T009 需要在 Cloudflare Dashboard 手动操作
- **验证优先**: 每个阶段完成后立即验证
- **域名配置**: 需要域名 DNS 管理权限
- **环境变量**: 敏感信息通过 Dashboard 配置，不提交到代码仓库
