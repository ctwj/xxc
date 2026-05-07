package telegram_sync

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"sync"
	"time"

	gotdsession "github.com/gotd/td/session"
	"gorm.io/gorm"
)

// DBStorage 实现 session.Storage 接口，使用数据库存储会话
type DBStorage struct {
	db  *gorm.DB
	key []byte
	mu  sync.Mutex
}

// NewDBStorage 创建数据库会话存储
func NewDBStorage(db *gorm.DB, encryptKey string) *DBStorage {
	// 确保密钥长度为 32 字节 (AES-256)
	keyBytes := []byte(encryptKey)
	if len(keyBytes) < 32 {
		paddedKey := make([]byte, 32)
		copy(paddedKey, keyBytes)
		keyBytes = paddedKey
	} else if len(keyBytes) > 32 {
		keyBytes = keyBytes[:32]
	}

	return &DBStorage{
		db:  db,
		key: keyBytes,
	}
}

// LoadSession 实现 session.Storage 接口
func (s *DBStorage) LoadSession(ctx context.Context) ([]byte, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var sess TelegramSession
	err := s.db.Where("status = ?", 1).First(&sess).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gotdsession.ErrNotFound // 返回 ErrNotFound 让 gotd/td 知道需要新会话
		}
		return nil, err
	}

	// 解密会话数据
	return s.decrypt(sess.SessionData)
}

// StoreSession 实现 session.Storage 接口
func (s *DBStorage) StoreSession(ctx context.Context, data []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 加密会话数据
	encrypted, err := s.encrypt(data)
	if err != nil {
		return err
	}

	session := &TelegramSession{
		SessionData: encrypted,
		SessionHash: s.generateHash(data),
		Status:      1,
		CreateTime:  time.Now().Unix(),
		UpdateTime:  time.Now().Unix(),
	}

	// 删除旧会话
	s.db.Where("1 = 1").Delete(&TelegramSession{})

	return s.db.Create(session).Error
}

// encrypt 加密数据
func (s *DBStorage) encrypt(plaintext []byte) ([]byte, error) {
	if len(plaintext) == 0 {
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

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return ciphertext, nil
}

// decrypt 解密数据
func (s *DBStorage) decrypt(ciphertext []byte) ([]byte, error) {
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

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

// generateHash 生成会话哈希
func (s *DBStorage) generateHash(data []byte) string {
	encoded := base64.StdEncoding.EncodeToString(data)
	if len(encoded) > 64 {
		return encoded[:64]
	}
	return encoded
}
