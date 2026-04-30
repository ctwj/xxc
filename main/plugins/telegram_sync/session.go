package telegram_sync

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
)

// SessionStorage 会话存储管理器
type SessionStorage struct {
	key       []byte // AES 加密密钥
	sessionID string // 会话标识
}

// NewSessionStorage 创建会话存储管理器
func NewSessionStorage(key string, sessionID string) *SessionStorage {
	// 确保密钥长度为 32 字节 (AES-256)
	keyBytes := []byte(key)
	if len(keyBytes) < 32 {
		// 填充到 32 字节
		paddedKey := make([]byte, 32)
		copy(paddedKey, keyBytes)
		keyBytes = paddedKey
	} else if len(keyBytes) > 32 {
		keyBytes = keyBytes[:32]
	}

	return &SessionStorage{
		key:       keyBytes,
		sessionID: sessionID,
	}
}

// Encrypt 加密数据
func (s *SessionStorage) Encrypt(plaintext []byte) ([]byte, error) {
	if len(plaintext) == 0 {
		return nil, nil
	}

	block, err := aes.NewCipher(s.key)
	if err != nil {
		return nil, err
	}

	// 使用 GCM 模式
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// 生成随机 nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	// 加密并附加 nonce
	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return ciphertext, nil
}

// Decrypt 解密数据
func (s *SessionStorage) Decrypt(ciphertext []byte) ([]byte, error) {
	if len(ciphertext) == 0 {
		return nil, nil
	}

	block, err := aes.NewCipher(s.key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	// 提取 nonce 和密文
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

// EncryptToString 加密并转为 base64 字符串
func (s *SessionStorage) EncryptToString(plaintext []byte) (string, error) {
	ciphertext, err := s.Encrypt(plaintext)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// DecryptFromString 从 base64 字符串解密
func (s *SessionStorage) DecryptFromString(encoded string) ([]byte, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, err
	}
	return s.Decrypt(ciphertext)
}

// GenerateSessionHash 生成会话哈希（用于验证）
func (s *SessionStorage) GenerateSessionHash(sessionData []byte) string {
	// 简单实现：使用 base64 编码的前 64 字符
	encoded := base64.StdEncoding.EncodeToString(sessionData)
	if len(encoded) > 64 {
		return encoded[:64]
	}
	return encoded
}

// ValidateSessionKey 验证会话密钥是否有效
func ValidateSessionKey(key string) error {
	if len(key) < 16 {
		return errors.New("session key must be at least 16 characters")
	}
	return nil
}

// GenerateRandomKey 生成随机密钥
func GenerateRandomKey() (string, error) {
	key := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(key), nil
}

// Context key for storing session
type contextKey string

const (
	SessionContextKey contextKey = "telegram_session"
	AuthContextKey    contextKey = "telegram_auth"
)

// SetSessionInContext 在 context 中设置会话
func SetSessionInContext(ctx context.Context, session *TelegramSession) context.Context {
	return context.WithValue(ctx, SessionContextKey, session)
}

// GetSessionFromContext 从 context 获取会话
func GetSessionFromContext(ctx context.Context) *TelegramSession {
	if session, ok := ctx.Value(SessionContextKey).(*TelegramSession); ok {
		return session
	}
	return nil
}
