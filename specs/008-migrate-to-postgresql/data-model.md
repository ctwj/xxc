# Data Model: 修复 MySQL 数据库编码问题

**Feature**: 008-migrate-to-postgresql
**Date**: 2026-05-10

## 概述

数据模型保持不变，本文档记录字符集配置的层级关系和转换策略。

## MySQL 字符集层级

```
┌─────────────────────────────────────┐
│         服务器级别 (my.cnf)          │  ← 默认设置
│   character-set-server = utf8mb4     │
│   collation-server = utf8mb4_0900_ai_ci │
└─────────────────────────────────────┘
              ↓ 继承
┌─────────────────────────────────────┐
│         数据库级别                   │  ← CREATE DATABASE
│   DEFAULT CHARACTER SET utf8mb4      │
│   DEFAULT COLLATE utf8mb4_0900_ai_ci │
└─────────────────────────────────────┘
              ↓ 继承
┌─────────────────────────────────────┐
│           表级别                     │  ← CREATE TABLE
│   DEFAULT CHARACTER SET utf8mb4      │
│   DEFAULT COLLATE utf8mb4_0900_ai_ci │
└─────────────────────────────────────┘
              ↓ 继承
┌─────────────────────────────────────┐
│          字段级别                    │  ← 列定义
│   CHARACTER SET utf8mb4             │
│   COLLATE utf8mb4_0900_ai_ci         │
└─────────────────────────────────────┘
```

## 连接字符集

应用程序连接时指定的字符集优先级最高：

```
DSN: user:pass@tcp(host:port)/dbname?charset=utf8mb4&collation=utf8mb4_0900_ai_ci
```

## 转换策略

### 1. 服务器级别配置

```ini
# /etc/mysql/mysql.conf.d/mysqld.cnf
[mysqld]
character-set-server = utf8mb4
collation-server = utf8mb4_0900_ai_ci

[client]
default-character-set = utf8mb4

[mysql]
default-character-set = utf8mb4
```

### 2. 数据库级别转换

```sql
ALTER DATABASE xxc CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
```

### 3. 表级别转换

```sql
-- 单个表
ALTER TABLE table_name CONVERT TO CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci;

-- 批量转换（生成语句）
SELECT CONCAT('ALTER TABLE ', TABLE_NAME, ' CONVERT TO CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci;')
FROM information_schema.TABLES
WHERE TABLE_SCHEMA = 'xxc' AND TABLE_TYPE = 'BASE TABLE';
```

## GORM 实体修改

### 修改前

```go
MediaUrls string `gorm:"type:text;default:''" json:"media_urls"`
```

### 修改后

```go
MediaUrls string `gorm:"type:text" json:"media_urls"`
```

## 验证查询

```sql
-- 查看数据库字符集
SHOW CREATE DATABASE xxc;

-- 查看表字符集
SELECT TABLE_NAME, TABLE_COLLATION 
FROM information_schema.TABLES 
WHERE TABLE_SCHEMA = 'xxc';

-- 查看连接字符集
SHOW VARIABLES LIKE 'character%';
SHOW VARIABLES LIKE 'collation%';
```
