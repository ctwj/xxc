# Quickstart: 修复 MySQL 数据库编码问题

**Feature**: 008-migrate-to-postgresql
**Date**: 2026-05-10

## 快速开始

### Step 1: 配置 MySQL 服务器

编辑 `/etc/mysql/mysql.conf.d/mysqld.cnf`：

```bash
sudo vim /etc/mysql/mysql.conf.d/mysqld.cnf
```

添加以下内容：

```ini
[mysqld]
character-set-server = utf8mb4
collation-server = utf8mb4_0900_ai_ci

[client]
default-character-set = utf8mb4

[mysql]
default-character-set = utf8mb4
```

重启 MySQL：

```bash
sudo systemctl restart mysql
```

### Step 2: 配置应用程序 DSN

编辑 `/opt/moss/conf.toml`：

```bash
sudo vim /opt/moss/conf.toml
```

修改 DSN：

```toml
dsn = 'xxc:your_password@tcp(127.0.0.1:3306)/xxc?charset=utf8mb4&collation=utf8mb4_0900_ai_ci&parseTime=True'
```

### Step 3: 批量转换现有表

执行以下命令：

```bash
# 1. 修改数据库默认字符集
mysql -u xxc -p -e "ALTER DATABASE xxc CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci;"

# 2. 生成并执行表转换语句
mysql -u xxc -p -N -e "
SELECT CONCAT('ALTER TABLE ', TABLE_NAME, ' CONVERT TO CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci;')
FROM information_schema.TABLES
WHERE TABLE_SCHEMA = 'xxc' AND TABLE_TYPE = 'BASE TABLE';
" | mysql -u xxc -p xxc
```

### Step 4: 重启应用程序

```bash
sudo systemctl restart moss
```

### Step 5: 验证

```bash
# 检查服务状态
sudo systemctl status moss

# 检查字符集
mysql -u xxc -p -e "SHOW VARIABLES LIKE 'character%'; SHOW VARIABLES LIKE 'collation%';"

# 检查表字符集
mysql -u xxc -p xxc -e "SELECT TABLE_NAME, TABLE_COLLATION FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'xxc';"
```

## 一键脚本

创建脚本 `fix-encoding.sh`：

```bash
#!/bin/bash
DB_USER="xxc"
DB_NAME="xxc"

echo "=== Step 1: 修改数据库默认字符集 ==="
mysql -u $DB_USER -p -e "ALTER DATABASE $DB_NAME CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci;"

echo "=== Step 2: 批量转换表 ==="
mysql -u $DB_USER -p -N -e "
SELECT CONCAT('ALTER TABLE ', TABLE_NAME, ' CONVERT TO CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci;')
FROM information_schema.TABLES
WHERE TABLE_SCHEMA = '$DB_NAME' AND TABLE_TYPE = 'BASE TABLE';
" | mysql -u $DB_USER -p $DB_NAME

echo "=== Step 3: 验证结果 ==="
mysql -u $DB_USER -p -e "SHOW VARIABLES LIKE 'character%';"
mysql -u $DB_USER -p $DB_NAME -e "SELECT TABLE_NAME, TABLE_COLLATION FROM information_schema.TABLES WHERE TABLE_SCHEMA = '$DB_NAME';"

echo "=== 完成 ==="
```

## 常见问题

### Q: 转换时报错 "Conversion impossible"

检查是否有数据包含无效编码：

```sql
-- 查找无效数据
SELECT * FROM table_name WHERE column_name != CONVERT(column_name USING utf8mb4);
```

### Q: 转换后中文乱码

可能是原数据使用了错误的编码，需要手动修复：

```sql
-- 尝试修复
UPDATE table_name SET column_name = CONVERT(CAST(CONVERT(column_name USING latin1) AS BINARY) USING utf8mb4);
```

### Q: MySQL 版本低于 8.0

使用 `utf8mb4_unicode_ci` 替代 `utf8mb4_0900_ai_ci`：

```ini
[mysqld]
character-set-server = utf8mb4
collation-server = utf8mb4_unicode_ci
```
