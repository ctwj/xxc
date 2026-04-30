package telegram_sync

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/gotd/td/telegram"
	"github.com/gotd/td/tg"
	"go.uber.org/zap"
)

// Client Telegram 客户端封装
type Client struct {
	config       *ClientConfig
	client       *telegram.Client
	api          *tg.Client
	dispatcher   *tg.UpdateDispatcher
	log          *zap.Logger

	// 状态
	connected    bool
	authenticated bool

	// 认证流程
	phoneCodeHash string
	authPhone     string

	// 控制
	mu     sync.RWMutex
	ctx    context.Context
	cancel context.CancelFunc
}

// ClientConfig 客户端配置
type ClientConfig struct {
	AppID          int
	AppHash        string
	PhoneNumber    string
	SessionKey     string
	AutoReconnect  bool
	ReconnectDelay int
}

// NewClient 创建 Telegram 客户端
func NewClient(config *ClientConfig, log *zap.Logger) (*Client, error) {
	if config.AppID <= 0 || config.AppHash == "" {
		return nil, errors.New("invalid Telegram API config: AppID and AppHash are required")
	}

	c := &Client{
		config: config,
		log:    log,
	}

	// 创建更新分发器
	dispatcher := tg.NewUpdateDispatcher()
	c.dispatcher = &dispatcher

	// 创建客户端选项
	opts := telegram.Options{
		UpdateHandler: dispatcher,
	}

	// 创建客户端 (gotd/td v0.108.0 API)
	c.client = telegram.NewClient(config.AppID, config.AppHash, opts)

	return c, nil
}

// Start 启动客户端
func (c *Client) Start(ctx context.Context) error {
	c.mu.Lock()
	c.ctx, c.cancel = context.WithCancel(ctx)
	c.mu.Unlock()

	return c.client.Run(c.ctx, func(ctx context.Context) error {
		c.api = tg.NewClient(c.client)
		c.connected = true
		c.log.Info("Telegram 客户端已连接")

		// 设置更新处理器
		c.setupHandlers()

		// 保持运行
		<-c.ctx.Done()
		return nil
	})
}

// Stop 停止客户端
func (c *Client) Stop() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.cancel != nil {
		c.cancel()
	}
	c.connected = false
	c.log.Info("Telegram 客户端已停止")
}

// setupHandlers 设置更新处理器
func (c *Client) setupHandlers() {
	c.dispatcher.OnNewMessage(c.onNewMessage)
}

// onNewMessage 处理新消息
func (c *Client) onNewMessage(ctx context.Context, entities tg.Entities, u *tg.UpdateNewMessage) error {
	// 检查是否为频道消息
	msg, ok := u.Message.(*tg.Message)
	if !ok {
		return nil
	}

	// 忽略发出的消息
	if msg.Out {
		return nil
	}

	// 获取消息来源
	peer := msg.GetPeerID()
	channelPeer, ok := peer.(*tg.PeerChannel)
	if !ok {
		// 不是频道消息，忽略
		return nil
	}

	channelID := channelPeer.GetChannelID()
	c.log.Debug("收到频道消息",
		zap.Int64("channel_id", channelID),
		zap.Int("message_id", msg.ID),
		zap.String("text", msg.Message))

	// 触发消息处理回调（由插件设置）
	// 实际处理在 handler.go 中实现

	return nil
}

// SendAuthCode 发送验证码
func (c *Client) SendAuthCode(ctx context.Context, phone string) (string, error) {
	if c.api == nil {
		return "", errors.New("client not initialized")
	}

	c.authPhone = phone

	// 发送验证码
	sentCode, err := c.api.AuthSendCode(ctx, &tg.AuthSendCodeRequest{
		PhoneNumber:   phone,
		APIID:         c.config.AppID,
		APIHash:       c.config.AppHash,
		Settings:      tg.CodeSettings{},
	})
	if err != nil {
		return "", fmt.Errorf("send code failed: %w", err)
	}

	// 保存 phone_code_hash
	code, ok := sentCode.(*tg.AuthSentCode)
	if !ok {
		return "", errors.New("unexpected sent code type")
	}
	c.phoneCodeHash = code.PhoneCodeHash

	return c.phoneCodeHash, nil
}

// VerifyAuthCode 验证登录码
func (c *Client) VerifyAuthCode(ctx context.Context, code string) error {
	if c.api == nil {
		return errors.New("client not initialized")
	}
	if c.phoneCodeHash == "" {
		return errors.New("no pending auth request")
	}

	// 验证码登录
	result, err := c.api.AuthSignIn(ctx, &tg.AuthSignInRequest{
		PhoneNumber:   c.authPhone,
		PhoneCodeHash: c.phoneCodeHash,
		PhoneCode:     code,
	})
	if err != nil {
		return fmt.Errorf("sign in failed: %w", err)
	}

	// 检查是否需要 2FA
	switch result.(type) {
	case *tg.AuthAuthorization:
		c.authenticated = true
		c.log.Info("Telegram 认证成功")
		return nil
	case *tg.AuthAuthorizationSignUpRequired:
		return errors.New("sign up required")
	default:
		return fmt.Errorf("unexpected auth result type: %T", result)
	}
}

// RestoreSession 恢复会话
func (c *Client) RestoreSession(session *TelegramSession) error {
	if session == nil || len(session.SessionData) == 0 {
		return errors.New("invalid session data")
	}

	// 解密会话数据
	storage := NewSessionStorage(c.config.SessionKey, "default")
	data, err := storage.Decrypt(session.SessionData)
	if err != nil {
		return fmt.Errorf("decrypt session failed: %w", err)
	}

	// TODO: 使用 gotd/td 的 session.Loader 恢复会话
	// 这需要更详细的实现
	c.log.Info("会话数据已解密", zap.Int("size", len(data)))

	return nil
}

// SaveSession 保存当前会话
func (c *Client) SaveSession() ([]byte, error) {
	// TODO: 使用 gotd/td 的 session.Loader 导出会话
	// 这需要更详细的实现
	return nil, errors.New("not implemented")
}

// IsConnected 检查是否已连接
func (c *Client) IsConnected() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.connected
}

// IsAuthenticated 检查是否已认证
func (c *Client) IsAuthenticated() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.authenticated
}

// GetAPI 获取 API 客户端
func (c *Client) GetAPI() *tg.Client {
	return c.api
}

// GetDispatcher 获取更新分发器
func (c *Client) GetDispatcher() *tg.UpdateDispatcher {
	return c.dispatcher
}

// Reconnect 重连
func (c *Client) Reconnect() error {
	c.mu.Lock()
	c.connected = false
	c.mu.Unlock()

	// 等待重连延迟
	if c.config.ReconnectDelay > 0 {
		time.Sleep(time.Duration(c.config.ReconnectDelay) * time.Second)
	}

	// TODO: 实现重连逻辑
	c.log.Info("尝试重连 Telegram")

	return nil
}

// GetChannelInfo 获取频道信息
func (c *Client) GetChannelInfo(ctx context.Context, channelID int64) (*tg.Channel, error) {
	if c.api == nil {
		return nil, errors.New("client not initialized")
	}

	// 获取频道信息
	channels, err := c.api.ChannelsGetChannels(ctx, []tg.InputChannelClass{
		&tg.InputChannel{
			ChannelID:  channelID,
			AccessHash: 0, // 需要正确的 access hash
		},
	})
	if err != nil {
		return nil, err
	}

	// 提取频道信息
	switch result := channels.(type) {
	case *tg.MessagesChats:
		// 从 chats 中获取频道信息
		for _, chat := range result.Chats {
			if channel, ok := chat.(*tg.Channel); ok {
				if channel.ID == channelID {
					return channel, nil
				}
			}
		}
	}

	return nil, errors.New("channel not found")
}

// JoinChannel 加入频道
func (c *Client) JoinChannel(ctx context.Context, channelID int64, accessHash int64) error {
	if c.api == nil {
		return errors.New("client not initialized")
	}

	_, err := c.api.ChannelsJoinChannel(ctx, &tg.InputChannel{
		ChannelID:  channelID,
		AccessHash: accessHash,
	})

	return err
}
