# Research: Telegram Channel Sync

**Date**: 2026-04-30
**Feature**: 002-telegram-channel-sync

## 技术选型研究

### 1. Telegram 客户端库选择

**Decision**: 使用 gotd/td 库

**Rationale**:
- 纯 Go 实现，与 Moss 后端技术栈一致
- 提供 MTProto API 完整支持
- 内置会话持久化机制
- 自动重连和数据中心迁移
- 活跃维护，文档完善

**Alternatives Considered**:
- **telegram-bot-api**: 仅支持 Bot API，无法监听频道更新（Bot 无法主动获取频道消息）
- **telethon/python-telegram-bot**: Python 实现，与项目技术栈不匹配
- **tdlib**: C++ 库，需要 CGO 绑定，部署复杂

### 2. 会话持久化方案

**Decision**: 使用 gotd/td 内置的 `session.FileStorage` + AES 加密

**Rationale**:
- gotd/td 提供开箱即用的会话存储接口
- `FileStorage` 支持文件持久化
- 可通过自定义加密层保护敏感数据
- 支持从 Telethon/TDesktop 导入会话

**Implementation Approach**:
```go
// 自定义加密存储
type EncryptedSessionStorage struct {
    storage *session.FileStorage
    key     []byte // 从配置中获取的加密密钥
}

func (s *EncryptedSessionStorage) LoadSession(ctx context.Context) ([]byte, error) {
    data, err := s.storage.LoadSession(ctx)
    if err != nil {
        return nil, err
    }
    return decrypt(data, s.key)
}

func (s *EncryptedSessionStorage) StoreSession(ctx context.Context, data []byte) error {
    encrypted, err := encrypt(data, s.key)
    if err != nil {
        return err
    }
    return s.storage.StoreSession(ctx, encrypted)
}
```

### 3. 频道消息监听方案

**Decision**: 使用 `tg.NewUpdateDispatcher()` + `OnNewMessage` 处理器

**Rationale**:
- gotd/td 提供统一的更新分发机制
- 支持过滤特定类型的更新
- 可区分频道消息、群组消息、私聊消息

**Implementation Approach**:
```go
dispatcher := tg.NewUpdateDispatcher()

dispatcher.OnNewMessage(func(ctx context.Context, entities tg.Entities, u *tg.UpdateNewMessage) error {
    // 检查是否为频道消息
    msg, ok := u.Message.(*tg.Message)
    if !ok {
        return nil
    }
    
    // 获取消息来源 Peer
    peer := msg.GetPeerID()
    channelPeer, ok := peer.(*tg.PeerChannel)
    if !ok {
        return nil // 非频道消息
    }
    
    // 检查频道是否在监听列表中
    channelID := channelPeer.GetChannelID()
    if !isChannelMonitored(channelID) {
        return nil
    }
    
    // 处理消息...
    return processMessage(ctx, channelID, msg)
})
```

### 4. 消息过滤规则引擎

**Decision**: 实现可配置的规则引擎，支持关键词白名单/黑名单、消息类型、长度限制

**Rationale**:
- 需要灵活的过滤规则配置
- 支持多种过滤条件组合
- 规则存储在数据库中，可通过管理界面配置

**Implementation Approach**:
```go
type FilterRule struct {
    Type         string   // whitelist, blacklist
    Keywords     []string // 关键词列表
    MessageTypes []string // text, photo, video, document
    MinLength    int      // 最小长度
    MaxLength    int      // 最大长度
}

func (r *FilterRule) Match(msg *tg.Message) bool {
    // 长度检查
    if r.MinLength > 0 && len(msg.Message) < r.MinLength {
        return false
    }
    if r.MaxLength > 0 && len(msg.Message) > r.MaxLength {
        return false
    }
    
    // 消息类型检查
    if len(r.MessageTypes) > 0 {
        if !containsMessageType(msg, r.MessageTypes) {
            return false
        }
    }
    
    // 关键词检查
    if r.Type == "whitelist" && len(r.Keywords) > 0 {
        return containsAnyKeyword(msg.Message, r.Keywords)
    }
    if r.Type == "blacklist" && len(r.Keywords) > 0 {
        return !containsAnyKeyword(msg.Message, r.Keywords)
    }
    
    return true
}
```

### 5. 媒体处理方案

**Decision**: 复用 Moss 现有的上传基础设施 (`infrastructure/support/upload`)

**Rationale**:
- Moss 已有完善的文件上传和存储机制
- 支持本地存储和云存储（S3、OSS、COS 等）
- 可复用 `SaveArticleImages` 插件的图片处理逻辑

**Implementation Approach**:
```go
func downloadAndUploadMedia(ctx context.Context, media *tg.MessageMedia) (string, error) {
    // 1. 从 Telegram 下载媒体文件
    // 2. 使用 Moss 上传基础设施上传
    // 3. 返回上传后的 URL
}
```

### 6. 插件生命周期管理

**Decision**: 遵循 Moss 插件接口 (`Load`, `Run`, `Unload`)

**Rationale**:
- 与现有插件架构一致
- 支持优雅启动和关闭
- 可通过管理界面控制插件状态

**Implementation Approach**:
```go
type TelegramChannelSync struct {
    client   *telegram.Client
    ctx      *pluginEntity.Plugin
    cancel   context.CancelFunc
    wg       sync.WaitGroup
    // 配置字段...
}

func (p *TelegramChannelSync) Info() *pluginEntity.PluginInfo {
    return &pluginEntity.PluginInfo{
        ID:         "TelegramChannelSync",
        About:      "Telegram 频道消息同步插件",
        RunEnable:  true,
        CronEnable: false, // 持续运行，不需要定时任务
    }
}

func (p *TelegramChannelSync) Load(ctx *pluginEntity.Plugin) error {
    p.ctx = ctx
    // 初始化 Telegram 客户端
    // 加载会话
    // 启动监听
    return nil
}

func (p *TelegramChannelSync) Unload() error {
    // 停止监听
    // 保存会话
    // 清理资源
    return nil
}
```

## 风险与缓解措施

### 1. Telegram API 限制

**风险**: Telegram 对 API 调用有频率限制，频繁调用可能导致账号被限制。

**缓解措施**:
- 使用官方推荐的调用间隔
- 实现请求队列和限流
- 监控 API 响应状态，遇到限制时自动降速

### 2. 会话安全

**风险**: 会话文件泄露可能导致账号被盗用。

**缓解措施**:
- 使用 AES-256 加密存储会话
- 加密密钥从配置中获取，不硬编码
- 会话文件权限设置为仅当前用户可读

### 3. 网络稳定性

**风险**: 网络不稳定可能导致连接断开，错过消息。

**缓解措施**:
- 利用 gotd/td 内置的自动重连机制
- 使用 `updates` 包恢复错过的更新
- 记录最后处理的消息 ID，重连后从断点继续

### 4. 消息去重

**风险**: 网络抖动或重连可能导致重复接收消息。

**缓解措施**:
- 使用 Telegram 消息 ID 作为唯一标识
- 在数据库中记录已处理的消息 ID
- 实现基于 Redis 或内存的去重缓存

## 依赖版本

| 依赖 | 版本 | 用途 |
|------|------|------|
| gotd/td | latest | Telegram MTProto 客户端 |
| go.uber.org/zap | v1.27.0 | 日志（已集成） |
| gorm.io/gorm | v1.25.12 | ORM（已集成） |
| gofiber/fiber/v2 | v2.52.5 | Web 框架（已集成） |

## 参考资料

- [gotd/td GitHub](https://github.com/gotd/td)
- [gotd/td 文档](https://pkg.go.dev/github.com/gotd/td)
- [Telegram MTProto API](https://core.telegram.org/mtproto)
- [Telegram API 文档](https://core.telegram.org/api)
- [How To Not Get Banned Guide](https://core.telegram.org/api/obtaining_api_id)
