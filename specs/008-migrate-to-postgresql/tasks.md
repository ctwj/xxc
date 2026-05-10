# Tasks: 修复 MySQL 数据库编码问题

**Input**: Design documents from `/specs/008-migrate-to-postgresql/`
**Prerequisites**: plan.md (required), spec.md (required for user stories), research.md, data-model.md, contracts/

**Tests**: 本任务列表不包含测试任务，因为这是运维配置任务。

**Organization**: 任务按用户故事组织，每个故事可独立实施和验证。

## Format: `[ID] [P?] [Story] Description`

- **[P]**: 可并行执行（不同文件，无依赖）
- **[Story]**: 任务所属的用户故事（US1, US2, US3, US4）
- 描述中包含具体文件路径

---

## Phase 1: Setup (准备工作)

**Purpose**: 准备修复所需的脚本和备份

- [x] T001 创建数据库备份脚本 `scripts/backup-db.sh`
- [x] T002 [P] 创建批量转换脚本 `scripts/fix-mysql-encoding.sql`
- [x] T003 [P] 创建验证脚本 `scripts/verify-encoding.sh`

---

## Phase 2: Foundational (MySQL 服务器配置)

**Purpose**: 配置 MySQL 服务器默认字符集（预防机制）

**⚠️ CRITICAL**: 这是预防问题的核心，必须首先完成

- [ ] T004 修改 MySQL 配置文件 `/etc/mysql/mysql.conf.d/mysqld.cnf` 添加字符集设置
- [ ] T005 重启 MySQL 服务并验证默认字符集
- [ ] T006 验证 MySQL 服务器字符集配置正确

**Checkpoint**: MySQL 服务器配置完成，新数据库将自动使用正确编码

---

## Phase 3: User Story 1 - 修复现有数据库编码 (Priority: P1) 🎯 MVP

**Goal**: 统一现有数据库中所有表的字符集

**Independent Test**: 执行批量转换脚本后，检查所有表的字符集是否为 `utf8mb4_0900_ai_ci`

### Implementation for User Story 1

- [ ] T007 [US1] 执行数据库备份脚本 `scripts/backup-db.sh`
- [ ] T008 [US1] 修改数据库默认字符集 `ALTER DATABASE xxc CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci`
- [ ] T009 [US1] 执行批量转换脚本转换所有表
- [ ] T010 [US1] 验证所有表字符集已统一

**Checkpoint**: 现有数据库编码已修复，应用程序应能正常启动

---

## Phase 4: User Story 2 - 配置 MySQL 服务器默认编码 (Priority: P1)

**Goal**: 配置 MySQL 服务器默认使用 UTF-8 编码

**Independent Test**: 创建新数据库，检查其默认编码

### Implementation for User Story 2

- [ ] T011 [US2] 验证 Phase 2 的 MySQL 配置已生效
- [ ] T012 [US2] 创建测试数据库验证默认编码
- [ ] T013 [US2] 删除测试数据库

**Checkpoint**: MySQL 服务器配置验证完成，新数据库将自动使用正确编码

---

## Phase 5: User Story 3 - 配置应用程序连接编码 (Priority: P1)

**Goal**: 在应用程序 DSN 中明确指定字符集

**Independent Test**: 修改配置后启动应用程序，验证连接使用的字符集

### Implementation for User Story 3

- [ ] T014 [US3] 修改应用程序 DSN 配置 `/opt/moss/conf.toml`
- [ ] T015 [US3] 重启应用程序服务
- [ ] T016 [US3] 验证应用程序连接字符集正确

**Checkpoint**: 应用程序连接编码配置完成

---

## Phase 6: User Story 4 - 修复 GORM 实体定义 (Priority: P2)

**Goal**: 移除 TEXT 类型字段的默认值定义

**Independent Test**: 修改实体定义后重启服务，验证表迁移成功

### Implementation for User Story 4

- [x] T017 [P] [US4] 修复 ArticleDetail 实体定义 `main/domain/core/entity/article.go`
- [ ] T018 [US4] 重新编译应用程序
- [ ] T019 [US4] 重启应用程序验证迁移成功

**Checkpoint**: GORM 实体定义修复完成

---

## Phase 7: Polish & Verification

**Purpose**: 最终验证和清理

- [ ] T020 执行验证脚本 `scripts/verify-encoding.sh`
- [ ] T021 测试中文/emoji 数据读写
- [ ] T022 测试创建新表验证预防机制
- [ ] T023 清理临时脚本和备份文件

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: 无依赖，可立即开始
- **Foundational (Phase 2)**: 无依赖，可与 Phase 1 并行
- **User Story 1 (Phase 3)**: 依赖 Phase 2 完成
- **User Story 2 (Phase 4)**: 依赖 Phase 2 完成
- **User Story 3 (Phase 5)**: 依赖 Phase 3 完成
- **User Story 4 (Phase 6)**: 依赖 Phase 3 完成，可与 Phase 5 并行
- **Polish (Phase 7)**: 依赖所有用户故事完成

### User Story Dependencies

- **User Story 1 (P1)**: 依赖 MySQL 配置完成
- **User Story 2 (P1)**: 依赖 MySQL 配置完成，用于验证
- **User Story 3 (P1)**: 依赖数据库编码修复完成
- **User Story 4 (P2)**: 依赖数据库编码修复完成

### Parallel Opportunities

- T002, T003 可并行执行
- Phase 2 和 Phase 1 可并行执行
- User Story 4 和 User Story 3 可并行执行

---

## Parallel Example: Setup Phase

```bash
# 并行执行所有准备任务:
Task: "创建批量转换脚本 scripts/fix-mysql-encoding.sql"
Task: "创建验证脚本 scripts/verify-encoding.sh"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. 完成 Phase 1: Setup
2. 完成 Phase 2: Foundational
3. 完成 Phase 3: User Story 1
4. **STOP and VALIDATE**: 验证数据库编码已修复
5. 应用程序应能正常启动

### Complete Implementation

1. Setup + Foundational → MySQL 配置完成
2. User Story 1 → 数据库编码修复
3. User Story 2 → 验证预防机制
4. User Story 3 → 应用程序配置
5. User Story 4 → GORM 实体修复
6. Polish → 最终验证

---

## Notes

- [P] 任务 = 不同文件，无依赖
- [Story] 标签将任务映射到特定用户故事
- 每个用户故事应可独立完成和验证
- 在每个检查点验证故事独立性
- 避免模糊任务、同文件冲突、跨故事依赖
