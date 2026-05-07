# Tasks: Telegram Channel Sync

**Input**: Design documents from `/specs/002-telegram-channel-sync/`
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/api.md

**Tests**: 未明确要求测试，本任务列表不包含测试任务。

**Organization**: 任务按用户故事组织，每个故事可独立实现和测试。

## Format: `[ID] [P?] [Story] Description`

- **[P]**: 可并行执行（不同文件，无依赖）
- **[Story]**: 任务所属用户故事（US1, US2, US3, US4, US5）
- 描述中包含精确文件路径

## Path Conventions

- **Backend**: `main/plugins/` (插件主目录)
- **Frontend**: `admin/src/views/plugin/options/` (插件配置组件)

---

## Phase 1: Setup (项目初始化)

**Purpose**: 添加依赖和创建基础结构

- [x] T001 Add gotd/td dependency to `main/go.mod`
- [x] T002 [P] Create plugin sub-package directory `main/plugins/telegram_sync/`
- [x] T003 [P] Create frontend config component `admin/src/views/plugin/options/TelegramChannelSync.vue`

---

## Phase 2: Foundational (基础架构)

**Purpose**: 核心实体定义和基础设施，必须在用户故事之前完成

**⚠️ CRITICAL**: 此阶段完成前，用户故事无法开始

- [x] T004 Define TelegramChannel entity in `main/plugins/telegram_sync/entity.go`
- [x] T005 [P] Define TelegramSyncLog entity in `main/plugins/telegram_sync/entity.go`
- [x] T006 [P] Define TelegramSession entity in `main/plugins/telegram_sync/entity.go`
- [x] T007 [P] Define TelegramPluginConfig struct (plugin JSON fields) in `main/plugins/TelegramChannelSync.go`
- [x] T008 Implement database auto-migration in `main/plugins/TelegramChannelSync.go` Load method
- [x] T009 Implement session encryption/decryption utilities in `main/plugins/telegram_sync/session.go`

**Checkpoint**: 基础实体和数据库迁移就绪，用户故事实现可以开始

---

## Phase 3: User Story 5 - Telegram 认证与会话持久化 (Priority: P1) 🎯 MVP

**Goal**: 实现首次认证和会话持久化，确保插件重启后自动恢复认证状态

**Independent Test**: 完成首次认证后，重启服务，验证系统自动恢复连接

### Implementation for User Story 5

- [x] T010 [US5] Implement Telegram client initialization in `main/plugins/telegram_sync/client.go`
- [x] T011 [US5] Implement authentication flow (send code, verify code) in `main/plugins/telegram_sync/client.go`
- [x] T012 [US5] Implement session storage with encryption in `main/plugins/telegram_sync/session.go`
- [x] T013 [US5] Implement session restoration on plugin load in `main/plugins/telegram_sync/session.go`
- [x] T014 [US5] Add auth status check method in `main/plugins/telegram_sync/client.go`
- [x] T015 [US5] Integrate auth flow into main plugin `main/plugins/TelegramChannelSync.go`

**Checkpoint**: 认证功能完整，可独立测试认证流程和会话持久化

---

## Phase 4: User Story 3 - 频道配置管理 (Priority: P1)

**Goal**: 通过管理界面添加、编辑、删除监听频道配置

**Independent Test**: 通过管理界面完成频道的增删改查操作

### Implementation for User Story 3

- [x] T016 [US3] Implement channel CRUD operations in `main/plugins/telegram_sync/entity.go` (ChannelConfig struct)
- [x] T017 [US3] Add channel config fields to plugin struct in `main/plugins/TelegramChannelSync.go` (ChannelsJSON)
- [x] T018 [US3] Implement channel list storage (JSON in plugin config via ChannelsJSON field)
- [x] T019 [US3] Create frontend config form for channels in `admin/src/views/plugin/options/TelegramChannelSync.vue`
- [x] T020 [US3] Add channel enable/disable toggle in frontend component
- [x] T020b [US3] Implement runtime dynamic channel addition (parseChannels called in handleChannelMessage)

**Checkpoint**: 频道配置管理功能完整，支持运行时动态添加

---

## Phase 5: User Story 1 - 频道消息自动同步发布 (Priority: P1)

**Goal**: 监听频道消息并自动发布为 CMS 文章

**Independent Test**: 订阅测试频道，发送消息，验证文章自动创建

### Implementation for User Story 1

- [x] T021 [US1] Implement update dispatcher setup in `main/plugins/telegram_sync/client.go` (handleUpdate)
- [x] T022 [US1] Implement OnNewMessage handler for channel messages in `main/plugins/telegram_sync/client.go` (handleNewChannelMessage)
- [x] T023 [US1] Implement message-to-article conversion in `main/plugins/TelegramChannelSync.go` (handleChannelMessage)
- [x] T024 [US1] Implement message ID deduplication in `main/plugins/TelegramChannelSync.go` (CheckMessageDuplicate)
- [ ] T025 [US1] Implement media download and upload in `main/plugins/telegram_sync/media.go`
- [x] T026 [US1] Integrate with Moss Article service for article creation (CreateArticle method)
- [ ] T027 [US1] Implement auto-reconnect logic in `main/plugins/telegram_sync/client.go`
- [x] T028 [US1] Start Telegram client in plugin Load method (initClient with session storage)

**Checkpoint**: 核心同步功能基本完整，可测试消息同步（媒体处理待完善）

---

## Phase 6: User Story 2 - 消息过滤规则配置 (Priority: P2)

**Goal**: 支持关键词白名单/黑名单、消息类型、长度限制等过滤规则

**Independent Test**: 配置过滤规则后，验证只有符合规则的消息被发布

### Implementation for User Story 2

- [x] T029 [US2] Implement filter rule engine in `main/plugins/telegram_sync/filter.go`
- [x] T030 [US2] Implement keyword whitelist filter in `main/plugins/telegram_sync/filter.go`
- [x] T031 [US2] Implement keyword blacklist filter in `main/plugins/telegram_sync/filter.go`
- [x] T032 [US2] Implement message type filter (text, photo, video) in `main/plugins/telegram_sync/filter.go`
- [x] T033 [US2] Implement message length filter in `main/plugins/telegram_sync/filter.go`
- [x] T034 [US2] Integrate filter into message handler in `main/plugins/TelegramChannelSync.go` (handleChannelMessage)
- [x] T035 [US2] Add filter config UI in `admin/src/views/plugin/options/TelegramChannelSync.vue`

**Checkpoint**: 过滤功能完整，可独立测试各种过滤规则 ✅

---

## Phase 7: User Story 4 - 同步状态监控 (Priority: P3)

**Goal**: 查看同步状态、日志、连接状态

**Independent Test**: 触发同步操作后，检查监控页面显示正确状态

### Implementation for User Story 4

- [x] T036 [US4] Implement sync log recording in `main/plugins/TelegramChannelSync.go` (RecordSyncLog)
- [x] T037 [US4] Implement connection status tracking in `main/plugins/telegram_sync/client.go` (GetReconnectStatus)
- [x] T038 [US4] Add status getter methods to plugin struct in `main/plugins/TelegramChannelSync.go` (GetStatus)
- [x] T039 [US4] Add log cleanup (by keep_days config) in `main/plugins/TelegramChannelSync.go` (startLogCleanupTask)
- [x] T040 [US4] Add status display in frontend config component (refreshLogs, refreshStatus)

**Checkpoint**: 监控功能完整，可独立测试状态查看 ✅

---

## Phase 8: Polish & Integration

**Purpose**: 最终集成和优化

- [x] T041 Register plugin in `main/startup/startup.go`
- [x] T042 [P] Add plugin info (ID, About) in `main/plugins/TelegramChannelSync.go`
- [x] T043 [P] Implement plugin Unload method for graceful shutdown
- [ ] T044 Add error handling and logging throughout the plugin
- [ ] T045 Test end-to-end flow: auth → add channel → sync message → verify article

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: 无依赖，可立即开始
- **Foundational (Phase 2)**: 依赖 Setup 完成 - 阻塞所有用户故事
- **User Story 5 (Phase 3)**: 依赖 Foundational 完成 - 认证是其他功能的基础
- **User Story 3 (Phase 4)**: 依赖 US5 完成 - 需要认证才能测试频道监听
- **User Story 1 (Phase 5)**: 依赖 US3 和 US5 完成 - 核心同步功能
- **User Story 2 (Phase 6)**: 依赖 US1 完成 - 过滤功能依赖消息处理
- **User Story 4 (Phase 7)**: 可与 US1 并行或之后
- **Polish (Phase 8)**: 依赖所有用户故事完成

### User Story Dependencies

```
US5 (认证) ──► US3 (频道配置) ──► US1 (消息同步) ──► US2 (过滤规则)
                     │
                     └──────────────────────► US4 (监控)
```

### Parallel Opportunities

- T002, T003 可并行（不同目录）
- T004-T007 可并行（同一文件但不同实体定义）
- US4 可与 US2 并行开发

---

## Parallel Example: Foundational Phase

```bash
# 并行执行实体定义:
Task: "Define TelegramChannel entity in main/plugins/telegram_sync/entity.go"
Task: "Define TelegramSyncLog entity in main/plugins/telegram_sync/entity.go"
Task: "Define TelegramSession entity in main/plugins/telegram_sync/entity.go"
Task: "Define TelegramPluginConfig struct in main/plugins/TelegramChannelSync.go"
```

---

## Implementation Strategy

### MVP First (User Story 5 + US3 + US1)

1. Complete Phase 1: Setup
2. Complete Phase 2: Foundational
3. Complete Phase 3: User Story 5 (认证)
4. Complete Phase 4: User Story 3 (频道配置)
5. Complete Phase 5: User Story 1 (消息同步)
6. **STOP and VALIDATE**: 测试完整同步流程
7. 可部署演示

### Incremental Delivery

1. Setup + Foundational → 基础就绪
2. Add US5 → 认证可用
3. Add US3 → 频道配置可用
4. Add US1 → 核心同步可用 (MVP!)
5. Add US2 → 过滤功能增强
6. Add US4 → 监控功能完善

---

## Notes

- [P] 任务 = 不同文件，无依赖，可并行
- [Story] 标签映射任务到具体用户故事
- 每个用户故事应可独立完成和测试
- 每个任务或逻辑组完成后提交
- 在任何检查点停止以独立验证故事
- 避免：模糊任务、同文件冲突、破坏独立性的跨故事依赖
