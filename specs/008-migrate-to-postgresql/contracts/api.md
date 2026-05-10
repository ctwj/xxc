# API Contract: MySQL 编码配置

**Feature**: 008-migrate-to-postgresql
**Date**: 2026-05-10

## 配置规范

### MySQL 服务器配置

**文件位置**: `/etc/mysql/mysql.conf.d/mysqld.cnf`

```ini
[mysqld]
character-set-server = utf8mb4
collation-server = utf8mb4_0900_ai_ci

[client]
default-character-set = utf8mb4

[mysql]
default-character-set = utf8mb4
```

### 应用程序 DSN 配置

**文件位置**: `/opt/moss/conf.toml`

```toml
db = 'mysql'
dsn = 'user:password@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&collation=utf8mb4_0900_ai_ci&parseTime=True'
```

### DSN 参数说明

| 参数 | 值 | 说明 |
|------|-----|------|
| charset | utf8mb4 | 连接字符集 |
| collation | utf8mb4_0900_ai_ci | 排序规则 |
| parseTime | True | 解析时间类型 |

## 批量转换脚本

**文件位置**: `scripts/fix-mysql-encoding.sql`

```sql
-- 1. 修改数据库默认字符集
ALTER DATABASE xxc CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci;

-- 2. 批量转换所有表
-- 以下语句需要动态生成并执行
```

### 生成转换语句

```sql
SELECT CONCAT('ALTER TABLE ', TABLE_NAME, ' CONVERT TO CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci;') AS stmt
FROM information_schema.TABLES
WHERE TABLE_SCHEMA = 'xxc' AND TABLE_TYPE = 'BASE TABLE';
```

## 验证命令

### 检查字符集配置

```bash
mysql -u root -p -e "SHOW VARIABLES LIKE 'character%'; SHOW VARIABLES LIKE 'collation%';"
```

### 检查数据库字符集

```bash
mysql -u root -p -e "SHOW CREATE DATABASE xxc;"
```

### 检查表字符集

```bash
mysql -u root -p xxc -e "SELECT TABLE_NAME, TABLE_COLLATION FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'xxc';"
```
