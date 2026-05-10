# Feature Specification: 修复 MySQL 数据库编码问题

**Feature Branch**: `008-migrate-to-postgresql`
**Created**: 2026-05-10
**Updated**: 2026-05-10
**Status**: Draft
**Input**: User description: "现在mysql数据庫出现了编码问题，不需要迁移过程，期望解决编码问题"

## Clarifications

### Session 2026-05-10

- Q: 是否需要迁移到 PostgreSQL？ → A: 不需要，保持 MySQL，只解决编码问题
- Q: 问题具体表现是什么？ → A: 创建数据库时没有指定编码，migrate 后不同表编码不同，需要统一解决方案并预防以后出现相同问题

## 问题背景

当前 MySQL 数据库存在字符集编码问题：
- **问题描述**: 创建数据库时没有指定编码，GORM migrate 后不同表使用了不同的字符集/排序规则
- **错误表现**: `Error 3988 (HY000): Conversion from collation utf8mb4_0900_ai_ci into utf8mb3_general_ci impossible for parameter`

**根本原因**：
1. 创建数据库时未指定默认字符集和排序规则
2. MySQL 服务器、数据库、表三层级字符集配置不一致
3. GORM migrate 时继承了不一致的字符集设置
4. GORM 实体定义中 TEXT 类型字段设置了无效的默认值

**目标**：
1. 修复现有数据库的编码问题
2. 建立预防机制，避免以后出现相同问题

## User Scenarios & Testing *(mandatory)*

### User Story 1 - 修复现有数据库编码 (Priority: P1)

作为系统管理员，我需要统一现有数据库中所有表的字符集，以便应用程序能够正常操作数据。

**Why this priority**: 这是当前最紧急的问题，必须立即修复。

**Independent Test**: 执行批量转换脚本后，检查所有表的字符集是否一致。

**Acceptance Scenarios**:

1. **Given** 数据库中存在编码不一致的表, **When** 执行批量转换脚本, **Then** 所有表转换为 `utf8mb4_0900_ai_ci`
2. **Given** 表已转换编码, **When** 查询历史数据, **Then** 数据正常显示无乱码
3. **Given** 编码已统一, **When** 应用程序启动, **Then** 无编码相关错误

---

### User Story 2 - 配置 MySQL 服务器默认编码 (Priority: P1)

作为系统管理员，我需要配置 MySQL 服务器默认使用 UTF-8 编码，以便新创建的数据库和表自动使用正确的编码。

**Why this priority**: 这是预防机制的核心，防止问题再次发生。

**Independent Test**: 修改配置后创建新数据库，检查其默认编码。

**Acceptance Scenarios**:

1. **Given** MySQL 配置文件已修改, **When** 重启 MySQL 服务, **Then** 服务正常运行
2. **Given** MySQL 已配置默认编码, **When** 创建新数据库, **Then** 数据库自动使用 `utf8mb4_0900_ai_ci`
3. **Given** MySQL 已配置默认编码, **When** 创建新表, **Then** 表自动继承正确的编码

---

### User Story 3 - 配置应用程序连接编码 (Priority: P1)

作为开发者，我需要在应用程序 DSN 中明确指定字符集，确保连接时使用正确的编码。

**Why this priority**: 即使服务器配置正确，应用程序连接参数错误也会导致问题。

**Independent Test**: 修改配置后启动应用程序，验证连接使用的字符集。

**Acceptance Scenarios**:

1. **Given** DSN 已添加字符集参数, **When** 应用程序连接数据库, **Then** 连接使用 `utf8mb4_0900_ai_ci`
2. **Given** 连接编码正确, **When** GORM 执行 migrate, **Then** 新表使用正确的编码

---

### User Story 4 - 修复 GORM 实体定义 (Priority: P2)

作为开发者，我需要移除 TEXT 类型字段的默认值定义，以便 GORM 迁移能够成功执行。

**Why this priority**: 这是导致迁移失败的次要原因，修复后可避免类似问题。

**Independent Test**: 修改实体定义后重启服务，验证表迁移成功。

**Acceptance Scenarios**:

1. **Given** `media_urls` 字段有 `default:''` 定义, **When** 移除默认值, **Then** GORM 生成的 DDL 不包含 DEFAULT 子句
2. **Given** 实体定义已修复, **When** 执行自动迁移, **Then** 迁移无错误

---

### Edge Cases

- 如果 MySQL 版本低于 8.0，`utf8mb4_0900_ai_ci` 可能不可用，需要使用 `utf8mb4_unicode_ci`
- 如果表中已有数据包含无效编码的字符，转换可能失败或产生乱码
- 如果数据量很大，批量转换可能需要较长时间
- 如果应用程序使用了 MySQL 特有的 SQL 语法，可能需要调整

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: MySQL 服务器配置必须设置默认字符集为 `utf8mb4` 和排序规则为 `utf8mb4_0900_ai_ci`
- **FR-002**: 应用程序 DSN 必须包含 `charset=utf8mb4&collation=utf8mb4_0900_ai_ci` 参数
- **FR-003**: 必须提供批量转换脚本，将现有数据库和所有表转换为统一编码
- **FR-004**: TEXT/BLOB 类型字段不得在 GORM tag 中定义默认值
- **FR-005**: 系统必须正确处理中文、emoji 等多字节字符
- **FR-006**: 新创建的数据库和表必须自动使用正确的编码（预防机制）

### Key Entities

- **MySQL 服务器配置**: `/etc/mysql/my.cnf` 或 `/etc/mysql/mysql.conf.d/mysqld.cnf`
- **数据库配置**: 数据库级别的字符集和排序规则
- **表配置**: 表级别的字符集和排序规则
- **GORM 实体定义**: Go 结构体中的字段 tag 定义
- **DSN 连接字符串**: 应用程序连接数据库时指定的字符集参数

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: 应用程序启动时数据库迁移无错误
- **SC-002**: 所有数据库表的字符集为 `utf8mb4`，排序规则为 `utf8mb4_0900_ai_ci`
- **SC-003**: 中文和 emoji 字符能够正确存储和读取
- **SC-004**: 历史数据查询无乱码
- **SC-005**: 新创建的数据库自动使用正确的编码（预防验证）
- **SC-006**: 新创建的表自动继承正确的编码（预防验证）

## Assumptions

- MySQL 版本为 8.0.x，支持 `utf8mb4_0900_ai_ci` 排序规则
- 应用程序使用 GORM 作为 ORM 框架
- 数据库连接使用 MySQL 驱动
- 现有数据不包含严重损坏的编码数据
- 有权限修改 MySQL 服务器配置文件