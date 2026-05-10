# Research: 修复 MySQL 数据库编码问题

**Feature**: 008-migrate-to-postgresql
**Date**: 2026-05-10

## 研究任务

### 1. MySQL 字符集层级结构

**Decision**: MySQL 字符集配置有三个层级，按优先级从高到低：
1. **连接级别**: DSN 参数或 SET NAMES 语句
2. **数据库级别**: CREATE DATABASE 时的默认设置
3. **服务器级别**: my.cnf 配置文件

**Rationale**: 理解层级关系是解决问题的关键。当上层未指定时，会继承下层设置。

**Alternatives considered**: 无。

### 2. MySQL 8.0 推荐字符集

**Decision**: 使用 `utf8mb4` 字符集 + `utf8mb4_0900_ai_ci` 排序规则。

**Rationale**:
- `utf8mb4` 是真正的 UTF-8，支持 4 字节字符（包括 emoji）
- `utf8mb4_0900_ai_ci` 是 MySQL 8.0 的默认排序规则，性能更好
- `ai` = accent insensitive（不区分重音）
- `ci` = case insensitive（不区分大小写）

**Alternatives considered**:
- `utf8mb4_unicode_ci`: MySQL 5.7 默认，兼容性更好但性能略差
- `utf8`: 不是真正的 UTF-8，只支持 3 字节字符

### 3. GORM TEXT 类型默认值问题

**Decision**: MySQL 不允许 TEXT/BLOB 类型有默认值，需要移除 GORM tag 中的 `default:''`。

**Rationale**: 这是 MySQL 的限制，不是 GORM 的问题。PostgreSQL 也不支持 TEXT 默认值。

**Alternatives considered**:
- 使用 VARCHAR 类型替代 TEXT（但 VARCHAR 有长度限制）
- 使用 NULL 作为默认值（但语义不同）

### 4. 批量转换脚本策略

**Decision**: 使用 SQL 脚本批量转换数据库和所有表。

**Rationale**:
- `ALTER DATABASE` 修改数据库默认字符集
- `ALTER TABLE ... CONVERT TO CHARACTER SET` 转换表和数据
- 可以在一条命令中生成所有转换语句

**Alternatives considered**:
- 使用 mysqldump 导出后重新导入（耗时太长）
- 使用第三方工具（不必要）

### 5. 预防机制设计

**Decision**: 三层防护机制：

| 层级 | 措施 | 配置位置 |
|------|------|----------|
| 服务器 | 默认字符集 | my.cnf |
| 连接 | DSN 参数 | conf.toml |
| 验证 | 启动检查 | 应用程序 |

**Rationale**: 多层防护确保即使某一层配置错误，其他层仍能保证正确性。

**Alternatives considered**:
- 仅依赖服务器配置（不够可靠）
- 仅依赖应用程序配置（新数据库可能出问题）

## 结论

解决方案清晰可行：
1. 配置 MySQL 服务器默认字符集
2. 应用程序 DSN 指定连接字符集
3. 批量转换现有数据库和表
4. 修复 GORM 实体定义
