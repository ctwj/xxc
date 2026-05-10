# Implementation Plan: 修复 MySQL 数据库编码问题

**Branch**: `008-migrate-to-postgresql` | **Date**: 2026-05-10 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/008-migrate-to-postgresql/spec.md`

## Summary

修复 MySQL 数据库字符集编码不一致问题，并建立三层防护机制防止问题再次发生：
1. MySQL 服务器级别：配置默认字符集
2. 应用程序级别：DSN 指定连接字符集
3. 数据库级别：批量转换现有表

## Technical Context

**Language/Version**: Go 1.21+
**Primary Dependencies**: GORM (gorm.io/gorm), MySQL Driver (gorm.io/driver/mysql)
**Storage**: MySQL 8.0.x
**Testing**: Go testing framework (go test)
**Target Platform**: Linux server
**Project Type**: Web service (CMS)
**Performance Goals**: 批量转换脚本执行时间 < 5 分钟
**Constraints**: UTF-8 编码，支持中文和 emoji
**Scale/Scope**: 中小型 CMS 系统

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

Constitution 文件为模板状态，无具体约束。修复工作符合项目现有架构。

## Project Structure

### Documentation (this feature)

```text
specs/008-migrate-to-postgresql/
├── plan.md              # This file
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 output
├── quickstart.md        # Phase 1 output
├── contracts/           # Phase 1 output
└── tasks.md             # Phase 2 output (NOT created yet)
```

### Source Code (repository root)

```text
main/
├── domain/core/entity/article.go  # 需要修复的 GORM 实体
├── conf.toml                       # 应用程序配置文件
└── ...

scripts/
└── fix-mysql-encoding.sql          # 批量转换脚本
```

**Structure Decision**: 使用现有项目结构，主要修改配置文件和 GORM 实体定义。

## Complexity Tracking

无需额外复杂性，修复工作符合项目现有架构。

## Implementation Phases

### Phase 1: MySQL 服务器配置

1. 修改 `/etc/mysql/mysql.conf.d/mysqld.cnf`
2. 重启 MySQL 服务
3. 验证默认字符集

### Phase 2: 应用程序配置

1. 修改 `conf.toml` 中的 DSN
2. 修复 GORM 实体定义（移除 TEXT 默认值）
3. 重启应用程序

### Phase 3: 批量转换现有表

1. 创建转换脚本
2. 执行转换
3. 验证结果

### Phase 4: 验证测试

1. 应用程序启动测试
2. 数据读写测试
3. 中文/emoji 测试
