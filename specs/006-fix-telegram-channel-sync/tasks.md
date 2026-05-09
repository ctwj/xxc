# Tasks: 修复 TelegramSync 频道消息同步

**Input**: Design documents from `/specs/006-fix-telegram-channel-sync/`
**Prerequisites**: plan.md, spec.md, research.md, data-model.md

**Tests**: 未明确要求，跳过测试任务。

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

---

## Phase 1: Setup

**Purpose**: 无需额外项目初始化，在现有代码基础上修改。

No setup tasks needed - this is a bug fix on existing codebase.

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: 实现 `channelAccessHasher` 基础设施，所有 User Story 依赖此组件。

**⚠️ CRITICAL**: US1 和 US3 都依赖此阶段完成

- [x] T001 在 `main/plugins/telegram_sync/client.go` 中添加 `channelAccessHasher` 结构体，实现 `ChannelAccessHasher` 接口的 `GetChannelAccessHash` 和 `SetChannelAccessHash` 方法
- [x] T002 在 `main/plugins/telegram_sync/client.go` 的 `Client` 结构体中添加 `accessHasher *channelAccessHasher` 字段
- [x] T003 在 `main/plugins/telegram_sync/client.go` 的 `NewClientWithStorage` 函数中，创建 `channelAccessHasher` 实例并配置到 `updates.Config.AccessHasher`
- [x] T004 在 `main/plugins/telegram_sync/client.go` 中添加 `LoadAccessHashes(channels []ChannelConfig)` 方法，从频道配置批量写入 access hash；同时在 `Client` 结构体或接收方上添加 `SetAccessHash(channelID, hash int64)` 方法供外部调用

**Checkpoint**: `channelAccessHasher` 基础设施就绪，`updates.Manager` 可通过它查找频道 access hash

---

## Phase 3: User Story 1 - 频道消息同步为文章 (Priority: P1) 🎯 MVP

**Goal**: 修复广播频道消息同步功能，使频道新消息能自动创建 CMS 文章

**Independent Test**: 在已绑定的 Telegram 广播频道发布一条消息，60 秒内检查 CMS 后台是否自动创建了对应文章

### Implementation for User Story 1

- [x] T005 [US1] 在 `main/plugins/TelegramChannelSync.go` 的 `initClient` 方法中，客户端创建后、启动前，调用 `LoadAccessHashes` 加载已配置频道的 access hash（需要从 `GetUserChannels` 或已保存的配置中获取 access hash）
- [x] T006 [US1] 在 `main/plugins/TelegramChannelSync.go` 的频道配置保存流程（`SaveChannels` 或相关 API handler）中，保存配置后调用 `LoadAccessHashes` 刷新 access hash 缓存
- [x] T007 [US1] 在 `main/plugins/telegram_sync/client.go` 的 `handleNewChannelMessage` 回调中，当收到频道消息但 access hash 未缓存时，尝试从消息实体（`tg.Entities`）中提取并缓存频道的 access hash，确保后续消息处理正常

**Checkpoint**: 广播频道消息能正确同步为 CMS 文章。此时应在测试频道中发送消息验证。

---

## Phase 4: User Story 2 - 群组消息继续正常工作 (Priority: P2)

**Goal**: 确保修复频道同步后，群组消息同步功能不受影响

**Independent Test**: 在已绑定的 Telegram 群组发送一条消息，确认仍能正常生成文章

### Implementation for User Story 2

- [x] T008 [US2] 验证 `main/plugins/telegram_sync/client.go` 中 `handleNewMessage` 回调逻辑未被修改，群组消息路径（`PeerChat` → `messageHandler`）完整保留
- [x] T009 [US2] 确认 `main/plugins/TelegramChannelSync.go` 的 `handleChannelMessage` 中群组消息的频道 ID 匹配逻辑不受 access hash 变更影响（群组消息不走 `updates.Manager` 的 channel state 路径）

**Checkpoint**: 群组消息同步功能保持 100% 不受影响

---

## Phase 5: User Story 3 - 频道配置不被意外清除 (Priority: P3)

**Goal**: 修复消息处理时频道配置被意外覆盖的 bug

**Independent Test**: 配置 3 个频道（2 启用、1 禁用），触发消息处理后检查配置是否完整保留

### Implementation for User Story 3

- [x] T010 [US3] 修复 `main/plugins/TelegramChannelSync.go` 第 402 行 `p.channels = enabledChannels` 覆盖 bug：删除该赋值语句，改为仅使用局部变量 `enabledChannels` 进行启用状态过滤，不修改 `p.channels`
- [x] T011 [US3] 在 `main/plugins/TelegramChannelSync.go` 的 `handleChannelMessage` 中，`GetChannelByID` 查找前不再过滤 enabled channels，而是直接在查找后检查 `channelConfig.Status == 1`（当前代码 422-424 行已有此检查）

**Checkpoint**: 消息处理后，所有频道配置（含禁用状态）保持不变

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: 代码清理和质量改进

- [x] T012 [P] 清理 `main/plugins/telegram_sync/client.go` 中所有 `fmt.Printf` 调试语句，保留 zap 结构化日志
- [x] T013 [P] 清理 `main/plugins/TelegramChannelSync.go` 中所有 `fmt.Printf` 调试语句，保留 zap 结构化日志
- [x] T014 运行 `quickstart.md` 中的验证步骤，确认所有验收场景通过

---

## Dependencies & Execution Order

### Phase Dependencies

- **Foundational (Phase 2)**: No dependencies - can start immediately
- **User Story 1 (Phase 3)**: Depends on Foundational (Phase 2) completion
- **User Story 2 (Phase 4)**: Depends on User Story 1 (Phase 3) - 需要先完成修改才能验证回归
- **User Story 3 (Phase 5)**: Depends on Foundational (Phase 2) completion - 可与 US1 并行
- **Polish (Phase 6)**: Depends on all user stories being complete

### User Story Dependencies

- **US1 (P1)**: Depends on Phase 2 - No dependencies on other stories
- **US2 (P2)**: Depends on US1 completion (回归验证)
- **US3 (P3)**: Depends on Phase 2 - Independent of US1/US2

### Task Dependencies

```
T001 → T002 → T003 → T004 (sequential, same file)
T004 → T005 (US1 needs LoadAccessHashes)
T004 → T006 (US1 needs save integration)
T005, T006 → T007 (can verify access hash flow)
T003 → T008, T009 (US2 regression check after core fix)
T001 → T010, T011 (US3 independent of access hash)
T010, T011 → T012, T013 (polish after implementation)
```

### Parallel Opportunities

- T010, T011 (US3) 可与 T005, T006, T007 (US1) 并行（不同逻辑，但同在 TelegramChannelSync.go，需注意合并冲突）
- T012, T013 可并行（不同文件）

---

## Parallel Example: Phase 2 + US3

```bash
# After T001-T004 complete, can proceed with:
# Track A: US1 implementation
Task T005: "Load access hashes after client init in TelegramChannelSync.go"
Task T006: "Refresh access hashes after channel config save"
Task T007: "Cache access hash from message entities"

# Track B: US3 fix (independent)
Task T010: "Fix p.channels overwrite bug in TelegramChannelSync.go:402"
Task T011: "Remove enabled channels filtering before GetChannelByID"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 2: Foundational (T001-T004)
2. Complete Phase 3: User Story 1 (T005-T007)
3. **STOP and VALIDATE**: 在广播频道发送测试消息，确认文章生成
4. 如验证通过，部署修复

### Incremental Delivery

1. Phase 2 → channelAccessHasher 基础就绪
2. + US1 → 频道消息同步恢复（MVP!）
3. + US2 → 确认群组不受影响
4. + US3 → 配置完整性修复
5. + Polish → 代码清理
