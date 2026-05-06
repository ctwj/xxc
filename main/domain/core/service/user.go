package service

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
	"moss/domain/core/entity"
	"moss/domain/core/repository"
	"moss/domain/core/repository/context"
)

type UserService struct{}

var User = &UserService{}

// Register registers a new user
func (s *UserService) Register(username, email, password string) (*entity.User, error) {
	// Validate input
	if username == "" {
		return nil, errors.New("username is required")
	}
	if email == "" {
		return nil, errors.New("email is required")
	}
	if password == "" {
		return nil, errors.New("password is required")
	}

	// Check if username exists
	exists, _ := repository.User.ExistsByUsername(username)
	if exists {
		return nil, errors.New("username already exists")
	}

	// Check if email exists
	exists, _ = repository.User.ExistsByEmail(email)
	if exists {
		return nil, errors.New("email already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	user := &entity.User{
		Username: username,
		Email:    email,
		Password: string(hashedPassword),
		Role:     "user",
	}

	err = repository.User.Create(user)
	if err != nil {
		return nil, errors.New("failed to create user")
	}

	return user, nil
}

// Login authenticates a user
func (s *UserService) Login(username, password string) (*entity.User, error) {
	// Try to find user by username or email
	user, err := repository.User.GetByUsername(username)
	if err != nil {
		// Try email
		user, err = repository.User.GetByEmail(username)
		if err != nil {
			return nil, errors.New("invalid credentials")
		}
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}

// GetByID gets a user by ID
func (s *UserService) GetByID(id uint) (*entity.User, error) {
	return repository.User.GetByID(id)
}

// GetByUsername gets a user by username
func (s *UserService) GetByUsername(username string) (*entity.User, error) {
	return repository.User.GetByUsername(username)
}

// GetByEmail gets a user by email
func (s *UserService) GetByEmail(email string) (*entity.User, error) {
	return repository.User.GetByEmail(email)
}

// List lists all users
func (s *UserService) List(ctx *context.Context) ([]entity.User, error) {
	return repository.User.List(ctx)
}

// Count counts all users
func (s *UserService) Count() (int64, error) {
	return repository.User.Count()
}

// UpdatePassword updates a user's password
func (s *UserService) UpdatePassword(id uint, oldPassword, newPassword string) error {
	user, err := repository.User.GetByID(id)
	if err != nil {
		return errors.New("user not found")
	}

	// Verify old password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword))
	if err != nil {
		return errors.New("invalid old password")
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("failed to hash password")
	}

	user.Password = string(hashedPassword)
	return repository.User.Update(user)
}

// Delete deletes a user
func (s *UserService) Delete(id uint) error {
	return repository.User.Delete(id)
}
