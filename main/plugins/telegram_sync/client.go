package telegram_sync

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/updates"
	"github.com/gotd/td/tg"
	"go.uber.org/zap"
)

// Client Telegram 客户端封装
type Client struct {
	config       *ClientConfig
	client       *telegram.Client
	api          *tg.Client
	sessionStore *DBStorage
	log          *zap.Logger

	// 更新管理器
	dispatcher   tg.UpdateDispatcher
	updatesMgr   *updates.Manager

	// 状态
	connected    bool
	authenticated bool

	// 认证流程
	phoneCodeHash string
	authPhone     string
	authPending   bool

	// 消息处理回调
	messageHandler func(ctx context.Context, channelID int64, msg *tg.Message)

	// 重连控制
	reconnectAttempts int
	maxReconnectAttempts int
	reconnectChan chan struct{}

	// 控制
	mu     sync.RWMutex
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
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

// NewClient 创建 Telegram 客户端（不包含 session storage）
func NewClient(config *ClientConfig, log *zap.Logger) (*Client, error) {
	if config.AppID <= 0 || config.AppHash == "" {
		return nil, errors.New("invalid Telegram API config: AppID and AppHash are required")
	}

	c := &Client{
		config:               config,
		log:                  log,
		maxReconnectAttempts: 5,
		reconnectChan:        make(chan struct{}, 1),
	}

	// 使用自定义 UpdateHandler 函数直接处理所有更新
	// 这比 UpdateDispatcher 更可靠，可以避免值接收者的问题
	opts := telegram.Options{
		UpdateHandler: telegram.UpdateHandlerFunc(c.handleUpdate),
	}

	// 创建客户端
	c.client = telegram.NewClient(config.AppID, config.AppHash, opts)
	log.Info("已设置自定义 UpdateHandler")

	return c, nil
}

// NewClientWithStorage 创建带 session storage 的客户端
func NewClientWithStorage(config *ClientConfig, storage *DBStorage, log *zap.Logger) (*Client, error) {
	if config.AppID <= 0 || config.AppHash == "" {
		return nil, errors.New("invalid Telegram API config: AppID and AppHash are required")
	}

	c := &Client{
		config:               config,
		sessionStore:         storage,
		log:                  log,
		maxReconnectAttempts: 5,
		reconnectChan:        make(chan struct{}, 1),
	}

	// 创建更新分发器
	c.dispatcher = tg.NewUpdateDispatcher()

	// 创建更新管理器（用于处理更新间隙和恢复）
	c.updatesMgr = updates.New(updates.Config{
		Handler: c.dispatcher,
	})

	// 设置更新处理回调
	c.dispatcher.OnNewChannelMessage(c.handleNewChannelMessage)
	c.dispatcher.OnNewMessage(c.handleNewMessage)

	// 创建客户端选项
	opts := telegram.Options{
		UpdateHandler:  c.updatesMgr, // 使用 updates.Manager 作为 UpdateHandler
		SessionStorage: storage,
	}

	// 创建客户端
	c.client = telegram.NewClient(config.AppID, config.AppHash, opts)
	log.Info("已设置 UpdateDispatcher 和 updates.Manager")

	return c, nil
}

// handleUpdate 处理所有 Telegram 更新（自定义 UpdateHandler）
func (c *Client) handleUpdate(ctx context.Context, u tg.UpdatesClass) error {
	fmt.Printf("=== handleUpdate 被调用: type=%T ===\n", u)
	c.log.Info("=== handleUpdate 被调用 ===", zap.String("type", fmt.Sprintf("%T", u)))

	switch update := u.(type) {
	case *tg.UpdatesCombined:
		c.log.Info("收到 UpdatesCombined", zap.Int("count", len(update.Updates)))
		for _, upd := range update.Updates {
			if err := c.handleSingleUpdate(ctx, upd); err != nil {
				c.log.Error("处理 UpdatesCombined 中的更新失败", zap.Error(err))
			}
		}
	case *tg.Updates:
		c.log.Info("收到 Updates", zap.Int("count", len(update.Updates)))
		for _, upd := range update.Updates {
			if err := c.handleSingleUpdate(ctx, upd); err != nil {
				c.log.Error("处理 Updates 中的更新失败", zap.Error(err))
			}
		}
	case *tg.UpdateShort:
		c.log.Info("收到 UpdateShort")
		return c.handleSingleUpdate(ctx, update.Update)
	default:
		c.log.Debug("忽略的更新类型", zap.String("type", fmt.Sprintf("%T", u)))
	}
	return nil
}

// handleSingleUpdate 处理单个更新
func (c *Client) handleSingleUpdate(ctx context.Context, u tg.UpdateClass) error {
	c.log.Info("=== handleSingleUpdate 被调用 ===", zap.String("type", fmt.Sprintf("%T", u)))

	switch update := u.(type) {
	case *tg.UpdateNewChannelMessage:
		// 使用空的 Entities，因为我们直接处理更新
		return c.handleNewChannelMessage(ctx, tg.Entities{}, update)
	case *tg.UpdateNewMessage:
		// 使用空的 Entities，因为我们直接处理更新
		return c.handleNewMessage(ctx, tg.Entities{}, update)
	default:
		c.log.Debug("忽略的单个更新类型", zap.String("type", fmt.Sprintf("%T", u)))
	}
	return nil
}

// handleNewChannelMessage 处理新频道消息（符合 UpdateDispatcher 回调格式）
func (c *Client) handleNewChannelMessage(ctx context.Context, e tg.Entities, u *tg.UpdateNewChannelMessage) error {
	fmt.Printf("=== handleNewChannelMessage 被调用: type=%T ===\n", u.Message)
	c.log.Info("=== handleNewChannelMessage 被调用 ===", zap.String("type", fmt.Sprintf("%T", u.Message)))

	// 检查是否为消息
	msg, ok := u.Message.(*tg.Message)
	if !ok {
		c.log.Debug("消息类型不是 *tg.Message", zap.String("type", fmt.Sprintf("%T", u.Message)))
		return nil
	}

	// 获取频道 ID
	peer := msg.GetPeerID()
	var channelID int64
	switch p := peer.(type) {
	case *tg.PeerChannel:
		channelID = p.GetChannelID()
	case *tg.PeerChat:
		// 群组消息，使用 Chat ID
		channelID = p.GetChatID()
	default:
		c.log.Warn("消息来源不是频道或群组", zap.String("type", fmt.Sprintf("%T", peer)))
		return nil
	}

	c.log.Info("收到频道消息详情",
		zap.Bool("out", msg.Out),
		zap.Int64("channel_id", channelID),
		zap.Int("message_id", msg.ID),
		zap.String("text", truncateMsgText(msg.Message, 50)))

	fmt.Printf("=== 收到频道消息: channelID=%d, messageID=%d, out=%v, text=%s ===\n",
		channelID, msg.ID, msg.Out, truncateMsgText(msg.Message, 50))

	// 忽略发出的消息
	if msg.Out {
		c.log.Debug("忽略发出的消息")
		return nil
	}

	// 调用消息处理回调
	c.mu.RLock()
	handler := c.messageHandler
	c.mu.RUnlock()

	if handler != nil {
		c.log.Info("调用消息处理回调", zap.Int64("channel_id", channelID))
		go handler(ctx, channelID, msg)
	} else {
		c.log.Warn("消息处理回调为空")
	}

	return nil
}

// handleNewMessage 处理新消息（群组/私聊）（符合 UpdateDispatcher 回调格式）
func (c *Client) handleNewMessage(ctx context.Context, e tg.Entities, u *tg.UpdateNewMessage) error {
	fmt.Printf("=== handleNewMessage 被调用: type=%T ===\n", u.Message)
	c.log.Info("=== handleNewMessage 被调用 ===", zap.String("type", fmt.Sprintf("%T", u.Message)))

	// 检查是否为消息
	msg, ok := u.Message.(*tg.Message)
	if !ok {
		c.log.Debug("消息类型不是 *tg.Message", zap.String("type", fmt.Sprintf("%T", u.Message)))
		return nil
	}

	c.log.Info("收到消息详情",
		zap.Bool("out", msg.Out),
		zap.String("peer_type", fmt.Sprintf("%T", msg.GetPeerID())),
		zap.Int("message_id", msg.ID),
		zap.String("text", truncateMsgText(msg.Message, 50)))

	// 忽略发出的消息
	if msg.Out {
		c.log.Debug("忽略发出的消息")
		return nil
	}

	// 获取消息来源
	peer := msg.GetPeerID()
	var channelID int64
	var isChannel bool

	switch p := peer.(type) {
	case *tg.PeerChannel:
		// 频道消息
		channelID = p.GetChannelID()
		isChannel = true
		c.log.Info("收到频道消息", zap.Int64("channel_id", channelID))

	case *tg.PeerChat:
		// 群组消息 - 群组 ID 需要转换
		// Telegram 群组使用 PeerChat，ID 是正数
		// 但在配置中可能使用不同的 ID 格式
		channelID = p.GetChatID()
		isChannel = false
		c.log.Info("收到群组消息", zap.Int64("chat_id", channelID))

	case *tg.PeerUser:
		// 私聊消息，忽略
		c.log.Info("收到私聊消息，忽略", zap.Int64("user_id", p.GetUserID()))
		return nil

	default:
		c.log.Warn("未知的消息来源类型", zap.String("type", fmt.Sprintf("%T", peer)))
		return nil
	}

	// 调用消息处理回调
	c.mu.RLock()
	handler := c.messageHandler
	c.mu.RUnlock()

	if handler != nil {
		c.log.Info("调用消息处理回调", zap.Int64("channel_id", channelID), zap.Bool("is_channel", isChannel))
		go handler(ctx, channelID, msg)
	} else {
		c.log.Warn("消息处理回调为空")
	}

	return nil
}

// Start 启动客户端（后台运行）
func (c *Client) Start(parentCtx context.Context) error {
	fmt.Println("=== Client.Start 开始 ===")
	c.mu.Lock()
	c.ctx, c.cancel = context.WithCancel(parentCtx)
	c.mu.Unlock()

	// 使用 channel 等待连接建立和认证检查完成
	connectedChan := make(chan struct{}, 1)
	authCheckChan := make(chan bool, 1)

	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		fmt.Println("=== 客户端 goroutine 开始运行 ===")
		err := c.client.Run(c.ctx, func(ctx context.Context) error {
			fmt.Println("=== Telegram 客户端 Run 回调开始 ===")
			c.api = tg.NewClient(c.client)
			c.mu.Lock()
			c.connected = true
			c.mu.Unlock()
			c.log.Info("Telegram 客户端已连接")
			fmt.Println("=== Telegram 客户端已连接 ===")

			// 通知连接已建立
			select {
			case connectedChan <- struct{}{}:
			default:
			}

			// 检查认证状态并获取用户 ID
			fmt.Println("=== 开始获取用户信息 ===")
			users, err := c.api.UsersGetUsers(ctx, []tg.InputUserClass{&tg.InputUserSelf{}})
			if err != nil {
				c.log.Error("获取用户信息失败", zap.Error(err))
				fmt.Printf("=== 获取用户信息失败: %v ===\n", err)
				// 不返回错误，继续运行
			} else {
				fmt.Printf("=== 获取到 %d 个用户 ===\n", len(users))
				var userID int64
				if len(users) > 0 {
					fmt.Printf("=== 用户类型: %T ===\n", users[0])
					if user, ok := users[0].(*tg.User); ok {
						userID = user.ID
						c.log.Info("获取用户ID成功", zap.Int64("user_id", userID))
						fmt.Printf("=== 用户ID: %d ===\n", userID)
					}
				}

				// 如果有 updates.Manager，启动它
				if c.updatesMgr != nil && userID > 0 {
					c.log.Info("启动 updates.Manager 进行更新处理")
					fmt.Println("=== 启动 updates.Manager ===")
					go func() {
						if err := c.updatesMgr.Run(ctx, c.api, userID, updates.AuthOptions{
							OnStart: func(ctx context.Context) {
								c.log.Info("updates.Manager 已启动")
								fmt.Println("=== updates.Manager 已启动 ===")
							},
						}); err != nil {
							c.log.Error("updates.Manager 运行错误", zap.Error(err))
							fmt.Printf("=== updates.Manager 运行错误: %v ===\n", err)
						}
					}()
				} else {
					fmt.Printf("=== updatesMgr=%v, userID=%d ===\n", c.updatesMgr != nil, userID)
				}
			}

			// 通知认证成功
			select {
			case authCheckChan <- true:
			default:
			}

			// 保持运行直到 context 取消
			<-c.ctx.Done()
			c.log.Info("Telegram 客户端正在关闭")
			return nil
		})
		if err != nil && err != context.Canceled {
			c.log.Error("客户端运行错误", zap.Error(err))
		}
		c.mu.Lock()
		c.connected = false
		c.mu.Unlock()
	}()

	// 等待连接建立或超时
	select {
	case <-connectedChan:
		c.log.Info("客户端连接成功")
	case <-time.After(10 * time.Second):
		return errors.New("连接超时")
	}

	// 等待认证检查完成
	select {
	case authenticated := <-authCheckChan:
		if authenticated {
			c.mu.Lock()
			c.authenticated = true
			c.mu.Unlock()
			c.log.Info("会话认证成功")
		} else {
			c.log.Info("需要重新认证")
		}
	case <-time.After(5 * time.Second):
		c.log.Warn("认证检查超时")
	}

	// 启动自动重连循环
	c.startReconnectLoop(c.ctx)

	return nil
}

// checkAuthStatus 检查认证状态
func (c *Client) checkAuthStatus(ctx context.Context, resultChan chan bool) {
	c.log.Info("开始检查认证状态")
	if c.api == nil {
		c.log.Warn("API 客户端为空，认证检查失败")
		resultChan <- false
		return
	}

	// 尝试获取当前用户信息来验证认证状态
	user, err := c.api.UsersGetFullUser(ctx, &tg.InputUserSelf{})

	if err == nil {
		c.log.Info("认证检查成功，用户已认证", zap.Any("user", user))
		resultChan <- true
	} else {
		c.log.Info("认证检查失败，需要重新认证", zap.Error(err))
		resultChan <- false
	}
}

// Stop 停止客户端
func (c *Client) Stop() {
	c.mu.Lock()
	if c.cancel != nil {
		c.cancel()
	}
	c.connected = false
	c.mu.Unlock()

	c.wg.Wait()
	c.log.Info("Telegram 客户端已停止")
}

// SetMessageHandler 设置消息处理回调
func (c *Client) SetMessageHandler(handler func(ctx context.Context, channelID int64, msg *tg.Message)) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.messageHandler = handler
}

// HasMessageHandler 检查是否设置了消息处理回调
func (c *Client) HasMessageHandler() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.messageHandler != nil
}

// startReconnectLoop 启动自动重连循环
func (c *Client) startReconnectLoop(parentCtx context.Context) {
	if !c.config.AutoReconnect {
		return
	}

	c.wg.Add(1)
	go func() {
		defer c.wg.Done()

		reconnectTimer := time.NewTicker(time.Duration(c.config.ReconnectDelay) * time.Second)
		defer reconnectTimer.Stop()

		for {
			select {
			case <-parentCtx.Done():
				c.log.Info("重连循环已停止")
				return
			case <-reconnectTimer.C:
				c.mu.RLock()
				connected := c.connected
				authenticated := c.authenticated
				c.mu.RUnlock()

				if !connected || !authenticated {
					c.log.Info("检测到连接断开，尝试重连...",
						zap.Bool("connected", connected),
						zap.Bool("authenticated", authenticated),
						zap.Int("attempts", c.reconnectAttempts))

					if c.reconnectAttempts >= c.maxReconnectAttempts {
						c.log.Error("重连尝试次数已达上限", zap.Int("max", c.maxReconnectAttempts))
						continue
					}

					if err := c.doReconnect(parentCtx); err != nil {
						c.log.Warn("重连失败", zap.Error(err))
						c.mu.Lock()
						c.reconnectAttempts++
						c.mu.Unlock()
					} else {
						c.mu.Lock()
						c.reconnectAttempts = 0
						c.mu.Unlock()
						c.log.Info("重连成功")
					}
				} else {
					// 连接正常，重置重连计数
					c.mu.Lock()
					c.reconnectAttempts = 0
					c.mu.Unlock()
				}
			}
		}
	}()
}

// doReconnect 执行重连
func (c *Client) doReconnect(ctx context.Context) error {
	c.log.Info("开始执行重连")

	// 停止当前客户端
	c.mu.Lock()
	c.connected = false
	c.authenticated = false
	c.mu.Unlock()

	// 重新启动客户端
	if err := c.Start(ctx); err != nil {
		return fmt.Errorf("重连启动失败: %w", err)
	}

	// 检查认证状态
	c.mu.RLock()
	authenticated := c.authenticated
	c.mu.RUnlock()

	if !authenticated {
		return errors.New("重连后认证状态无效")
	}

	return nil
}

// GetReconnectStatus 获取重连状态
func (c *Client) GetReconnectStatus() map[string]interface{} {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return map[string]interface{}{
		"auto_reconnect":       c.config.AutoReconnect,
		"reconnect_attempts":   c.reconnectAttempts,
		"max_reconnect":        c.maxReconnectAttempts,
		"reconnect_delay":      c.config.ReconnectDelay,
	}
}

// truncateMsgText 截断文本
func truncateMsgText(text string, maxLen int) string {
	if len(text) <= maxLen {
		return text
	}
	return text[:maxLen] + "..."
}

// SendAuthCode 发送验证码（需要客户端已启动）
func (c *Client) SendAuthCode(ctx context.Context, phone string) (string, error) {
	if c.api == nil {
		return "", errors.New("客户端未连接，请先启动客户端")
	}

	c.mu.Lock()
	c.authPhone = phone
	c.authPending = true
	c.mu.Unlock()

	// 发送验证码
	sentCode, err := c.api.AuthSendCode(ctx, &tg.AuthSendCodeRequest{
		PhoneNumber:   phone,
		APIID:         c.config.AppID,
		APIHash:       c.config.AppHash,
		Settings:      tg.CodeSettings{},
	})
	if err != nil {
		c.log.Error("发送验证码失败", zap.Error(err))
		return "", fmt.Errorf("send code failed: %w", err)
	}

	// 保存 phone_code_hash
	code, ok := sentCode.(*tg.AuthSentCode)
	if !ok {
		return "", errors.New("unexpected sent code type")
	}

	c.mu.Lock()
	c.phoneCodeHash = code.PhoneCodeHash
	c.mu.Unlock()

	c.log.Info("验证码已发送", zap.String("phone", phone))
	return c.phoneCodeHash, nil
}

// VerifyAuthCode 验证登录码（需要客户端已启动）
func (c *Client) VerifyAuthCode(ctx context.Context, code string) error {
	if c.api == nil {
		return errors.New("客户端未连接，请先启动客户端")
	}

	c.mu.RLock()
	phoneCodeHash := c.phoneCodeHash
	authPhone := c.authPhone
	c.mu.RUnlock()

	if phoneCodeHash == "" {
		return errors.New("no pending auth request")
	}

	// 验证码登录
	result, err := c.api.AuthSignIn(ctx, &tg.AuthSignInRequest{
		PhoneNumber:   authPhone,
		PhoneCodeHash: phoneCodeHash,
		PhoneCode:     code,
	})
	if err != nil {
		c.log.Error("验证码验证失败", zap.Error(err))
		return fmt.Errorf("sign in failed: %w", err)
	}

	// 检查是否需要 2FA
	switch authResult := result.(type) {
	case *tg.AuthAuthorization:
		c.mu.Lock()
		c.authenticated = true
		c.authPending = false
		c.mu.Unlock()
		c.log.Info("Telegram 认证成功")
		return nil
	case *tg.AuthAuthorizationSignUpRequired:
		return errors.New("sign up required")
	default:
		return fmt.Errorf("unexpected auth result type: %T", authResult)
	}
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

// Reconnect 重连
func (c *Client) Reconnect() error {
	c.mu.Lock()
	c.connected = false
	c.mu.Unlock()

	if c.config.ReconnectDelay > 0 {
		time.Sleep(time.Duration(c.config.ReconnectDelay) * time.Second)
	}

	c.log.Info("尝试重连 Telegram")
	return nil
}

// GetChannelInfo 获取频道信息
func (c *Client) GetChannelInfo(ctx context.Context, channelID int64) (*tg.Channel, error) {
	if c.api == nil {
		return nil, errors.New("client not initialized")
	}

	channels, err := c.api.ChannelsGetChannels(ctx, []tg.InputChannelClass{
		&tg.InputChannel{
			ChannelID:  channelID,
			AccessHash: 0,
		},
	})
	if err != nil {
		return nil, err
	}

	switch result := channels.(type) {
	case *tg.MessagesChats:
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

// GetUserChannels 获取用户加入的所有频道
func (c *Client) GetUserChannels(ctx context.Context) ([]map[string]interface{}, error) {
	if c.api == nil {
		return nil, errors.New("客户端未连接")
	}

	// 使用 MessagesGetDialogs 获取所有对话（包括频道）
	result, err := c.api.MessagesGetDialogs(ctx, &tg.MessagesGetDialogsRequest{
		Limit:      100,
		OffsetPeer: &tg.InputPeerEmpty{},
	})
	if err != nil {
		return nil, fmt.Errorf("获取对话列表失败: %w", err)
	}

	var channels []map[string]interface{}

	// 解析结果
	switch r := result.(type) {
	case *tg.MessagesDialogs:
		for _, chat := range r.Chats {
			if channel, ok := chat.(*tg.Channel); ok {
				channels = append(channels, map[string]interface{}{
					"id":          channel.ID,
					"access_hash": channel.AccessHash,
					"title":       channel.Title,
					"username":    channel.Username,
					"is_broadcast": channel.Broadcast,
					"is_megagroup": channel.Megagroup,
				})
			}
		}

	case *tg.MessagesDialogsSlice:
		for _, chat := range r.Chats {
			if channel, ok := chat.(*tg.Channel); ok {
				channels = append(channels, map[string]interface{}{
					"id":          channel.ID,
					"access_hash": channel.AccessHash,
					"title":       channel.Title,
					"username":    channel.Username,
					"is_broadcast": channel.Broadcast,
					"is_megagroup": channel.Megagroup,
				})
			}
		}
	}

	c.log.Info("获取频道列表成功", zap.Int("count", len(channels)))
	return channels, nil
}