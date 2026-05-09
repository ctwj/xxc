# API Contracts: 修复 TelegramSync 频道消息同步

**Feature**: 006-fix-telegram-channel-sync
**Date**: 2026-05-09

## 概述

本次修复不涉及外部 API 变更，所有变更均为插件内部实现。现有 API 行为保持不变。

## 涉及的内部接口

### ChannelAccessHasher 接口（新增实现）

gotd/td 库定义的接口，本次新增实现：

```go
type ChannelAccessHasher interface {
    GetChannelAccessHash(ctx context.Context, userID int64, channelID int64) (hash int64, found bool, err error)
}
```

**实现**：`channelAccessHasher` 结构体，使用内存 map 存储频道 access hash。

**数据来源**：`GetUserChannels()` API 返回的 `access_hash` 字段。

### 现有 API（无变更）

| API | 方法 | 说明 |
|-----|------|------|
| 频道列表获取 | `GetUserChannels` | 返回频道信息（含 access_hash），行为不变 |
| 频道配置保存 | `SaveChannels` | 保存 JSON 配置，行为不变 |
| 消息同步 | `handleChannelMessage` | 处理消息回调，行为不变 |
