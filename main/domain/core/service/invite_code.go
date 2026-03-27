package service

import (
	"errors"
	"time"

	"moss/domain/core/entity"
	"moss/domain/core/repository"
	"moss/domain/core/repository/context"
)

var InviteCode = new(InviteCodeService)

type InviteCodeService struct{}

// Create 创建邀请码
func (s *InviteCodeService) Create(maxUses int, expireDays int, createdBy uint) (*entity.InviteCode, error) {
	code := &entity.InviteCode{
		Code:       entity.GenerateInviteCode(),
		MaxUses:    maxUses,
		UsedCount:  0,
		CreatedBy:  createdBy,
		CreateTime: time.Now().Unix(),
	}

	// 设置过期时间
	if expireDays > 0 {
		code.ExpireAt = time.Now().AddDate(0, 0, expireDays).Unix()
	}

	if err := repository.InviteCode.Create(code); err != nil {
		return nil, err
	}

	return code, nil
}

// Validate 验证邀请码
func (s *InviteCodeService) Validate(code string) (*entity.InviteCode, error) {
	if code == "" {
		return nil, errors.New("invite code is required")
	}

	inviteCode, err := repository.InviteCode.GetByCode(code)
	if err != nil {
		return nil, errors.New("invalid invite code")
	}

	if inviteCode.ID == 0 {
		return nil, errors.New("invalid invite code")
	}

	if inviteCode.IsExpired() {
		return nil, errors.New("invite code has expired")
	}

	if inviteCode.IsUsedUp() {
		return nil, errors.New("invite code has been used up")
	}

	return inviteCode, nil
}

// Use 使用邀请码（增加使用次数）
func (s *InviteCodeService) Use(code string) error {
	inviteCode, err := s.Validate(code)
	if err != nil {
		return err
	}

	return repository.InviteCode.IncrementUsed(inviteCode.ID)
}

// GetByID 根据ID获取邀请码
func (s *InviteCodeService) GetByID(id uint) (*entity.InviteCode, error) {
	return repository.InviteCode.Get(id)
}

// GetByCode 根据邀请码字符串获取
func (s *InviteCodeService) GetByCode(code string) (*entity.InviteCode, error) {
	return repository.InviteCode.GetByCode(code)
}

// List 获取邀请码列表
func (s *InviteCodeService) List(ctx *context.Context) ([]entity.InviteCode, error) {
	return repository.InviteCode.List(ctx)
}

// Count 统计邀请码数量
func (s *InviteCodeService) Count() (int64, error) {
	return repository.InviteCode.Count()
}

// Delete 删除邀请码
func (s *InviteCodeService) Delete(id uint) error {
	return repository.InviteCode.Delete(id)
}
