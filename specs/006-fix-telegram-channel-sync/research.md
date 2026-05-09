# Research: 修复 TelegramSync 频道消息同步

**Feature**: 006-fix-telegram-channel-sync
**Date**: 2026-05-09

## 根本原因分析

### 问题现象

TelegramSync 插件在群组（supergroup/megagroup）中正常工作，但在广播频道（broadcast channel）中不产生文章。

### 根因：`updates.Manager` 缺少频道 Access Hash

**核心代码路径**（生产模式 `NewClientWithStorage`）：

```
Telegram Server → gotd/td client → updates.Manager → UpdateDispatcher → OnNewChannelMessage/OnNewMessage
```

**问题出在 `updates.Manager` 对 `UpdateNewChannelMessage` 的处理**：

1. **插件初始化代码**（`client.go:106-108`）：
   ```go
   c.updatesMgr = updates.New(updates.Config{
       Handler: c.dispatcher,
   })
   ```
   只设置了 `Handler`，没有设置 `AccessHasher` 和 `Storage`，使用默认的内存存储（空）。

2. **Manager.Run() 启动时**（gotd/td `manager.go`）：从 `Storage` 加载已知频道，从 `AccessHasher` 获取 access hash。由于都是空内存存储，`channels` map 为空。

3. **收到 `UpdateNewChannelMessage` 时**（gotd/td `state_apply.go:applyCombined`）：
   - 检测到是 `IsChannelPtsUpdate` → 调用 `handleChannel()`
   - `handleChannel()` 在 `channels` map 中找不到该频道
   - 调用 `AccessHasher.GetChannelAccessHash()` → 未找到（内存为空）
   - 尝试 `restoreAccessHash()` 通过 `getDifference` 恢复 → 可能失败
   - **最终结果：更新被静默丢弃**，日志级别为 Debug（"Failed to recover missing access hash, update ignored"）

4. **为什么群组正常**：群组消息以 `UpdateNewMessage` 到达，走 `handlePts → applyPts → handler.Handle()` 路径。此路径**不需要 access hash 查找**，直接传递给 `UpdateDispatcher.OnNewMessage` 回调。

### 次要 Bug：频道配置被意外覆盖

**代码位置**：`TelegramChannelSync.go:402`

```go
p.channels = enabledChannels  // 每次处理消息时，用启用频道覆盖完整列表
```

每次处理消息时，`p.channels` 被替换为仅包含 `Status == 1` 的频道列表。禁用的频道从内存中永久丢失，直到下次从 JSON 重新解析。

## 决策：修复方案

### 方案 A：自定义 `ChannelAccessHasher`（推荐）

**决策**：实现一个自定义的 `ChannelAccessHasher`，从频道配置中提供 access hash。

**理由**：
- 保持 `updates.Manager` 的 gap recovery（断线恢复）能力
- 最小化代码变更
- 与 gotd/td 库的设计意图一致
- 满足 FR-008（断线自动恢复）

**实现要点**：
1. 创建 `channelAccessHasher` 结构体，实现 `ChannelAccessHasher` 接口
2. 从频道配置（`GetUserChannels` API 返回的 `access_hash`）中读取并缓存 access hash
3. 在 `NewClientWithStorage` 中配置到 `updates.Config.AccessHasher`
4. 当频道配置变更时，更新缓存

### 方案 B：绕过 `updates.Manager`（备选）

**被否决原因**：
- 丢失 gap recovery 能力
- 需要自行实现消息去重和状态管理
- 增加维护复杂度

### 方案 C：使用 gotd/td `peers.Manager`（备选）

**被否决原因**：
- `peers.Manager` 是完整的 peer 管理系统，对当前需求过重
- 需要额外的存储后端
- 引入不必要的复杂性

## 修复范围

1. **主要修复**：实现自定义 `ChannelAccessHasher`，让 `updates.Manager` 能找到频道的 access hash
2. **附加修复**：修复 `p.channels = enabledChannels` 覆盖 bug
3. **不需要修改**：`handleNewChannelMessage` 和 `handleNewMessage` 回调逻辑本身是正确的

## 涉及文件

| 文件 | 变更类型 | 说明 |
|------|---------|------|
| `main/plugins/telegram_sync/client.go` | 修改 | 添加 `channelAccessHasher`，配置到 `updates.Config` |
| `main/plugins/TelegramChannelSync.go` | 修改 | 修复 `p.channels` 覆盖 bug，更新 access hash 缓存 |
