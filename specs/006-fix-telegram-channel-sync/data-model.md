# Data Model: 修复 TelegramSync 频道消息同步

**Feature**: 006-fix-telegram-channel-sync
**Date**: 2026-05-09

## 现有实体（无变更）

### TelegramChannel

| 字段 | 类型 | 说明 |
|------|------|------|
| ID | int | 主键，自增 |
| ChannelID | int64 | Telegram 频道 ID（唯一索引） |
| ChannelName | string | 频道名称 |
| ChannelLink | string | 频道链接 |
| Status | int | 1=启用, 0=禁用 |
| CategoryID | int | 目标分类 ID |
| ArticleStatus | bool | 同步文章发布状态 |
| ArticleAuthor | int | 文章作者 ID |
| FilterKeywords | string | 关键词过滤规则 JSON |
| FilterMessageTypes | string | 消息类型过滤 |
| FilterMinLength | int | 最小消息长度 |
| FilterMaxLength | int | 最大消息长度 |
| LastSyncTime | int64 | 最后同步时间戳 |
| LastMessageID | int64 | 最后处理的消息 ID |
| TotalSyncCount | int | 总同步数量 |
| ErrorCount | int | 错误计数 |
| CreateTime | int64 | 创建时间 |
| UpdateTime | int64 | 更新时间 |
| Remark | string | 备注 |

### ChannelConfig（运行时配置）

| 字段 | 类型 | 说明 |
|------|------|------|
| ChannelID | int64 | 频道 ID |
| ChannelName | string | 频道名称 |
| ChannelLink | string | 频道链接 |
| Status | int | 启用状态 |
| CategoryID | int | 目标分类 |
| ArticleStatus | bool | 发布状态 |
| FilterKeywords | string | 关键词过滤 |
| FilterMessageTypes | string | 消息类型过滤 |
| FilterMinLength | int | 最小长度 |
| FilterMaxLength | int | 最大长度 |

## 新增运行时结构（非持久化）

### channelAccessHasher

| 字段 | 类型 | 说明 |
|------|------|------|
| hashes | map[int64]int64 | 频道 ID → access hash 映射 |
| mu | sync.RWMutex | 读写锁 |

**接口方法**：
- `GetChannelAccessHash(ctx, userID, channelID) → (hash, found, error)`
- `SetChannelAccessHash(ctx, userID, channelID, hash) → error`

**说明**：这是纯内存结构，实现 gotd/td 的 `ChannelAccessHasher` 接口。不涉及数据库变更。

## 实体关系

```
TelegramChannel (DB) ──→ ChannelConfig (JSON in-memory)
                         ↓
                    channelAccessHasher (in-memory)
                         ↓
                    updates.Manager (gotd/td)
                         ↓
                    UpdateDispatcher (gotd/td)
                         ↓
                    handleNewChannelMessage / handleNewMessage
```
