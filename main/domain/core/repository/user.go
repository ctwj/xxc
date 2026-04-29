package repository

import (
	"moss/domain/core/entity"
	"moss/domain/core/repository/context"
	"moss/infrastructure/persistent/db"
)

type UserRepository struct{}

var User = &UserRepository{}

// Create creates a new user
func (r *UserRepository) Create(user *entity.User) error {
	return db.DB.Create(user).Error
}

// Update updates a user
func (r *UserRepository) Update(user *entity.User) error {
	return db.DB.Save(user).Error
}

// Delete deletes a user by ID
func (r *UserRepository) Delete(id uint) error {
	return db.DB.Delete(&entity.User{}, id).Error
}

// GetByID gets a user by ID
func (r *UserRepository) GetByID(id uint) (*entity.User, error) {
	var user entity.User
	err := db.DB.First(&user, id).Error
	return &user, err
}

// GetByUsername gets a user by username
func (r *UserRepository) GetByUsername(username string) (*entity.User, error) {
	var user entity.User
	err := db.DB.Where("username = ?", username).First(&user).Error
	return &user, err
}

// GetByEmail gets a user by email
func (r *UserRepository) GetByEmail(email string) (*entity.User, error) {
	var user entity.User
	err := db.DB.Where("email = ?", email).First(&user).Error
	return &user, err
}

// List lists all users
func (r *UserRepository) List(ctx *context.Context) ([]entity.User, error) {
	var users []entity.User
	query := db.DB
	if ctx != nil {
		if ctx.Limit > 0 {
			query = query.Limit(ctx.Limit)
			if ctx.Page > 0 {
				query = query.Offset((ctx.Page - 1) * ctx.Limit)
			}
		}
		if ctx.Order != "" {
			query = query.Order(ctx.Order)
		}
	}
	err := query.Find(&users).Error
	return users, err
}

// Count counts all users
func (r *UserRepository) Count() (int64, error) {
	var count int64
	err := db.DB.Model(&entity.User{}).Count(&count).Error
	return count, err
}

// ExistsByUsername checks if a user exists by username
func (r *UserRepository) ExistsByUsername(username string) (bool, error) {
	var count int64
	err := db.DB.Model(&entity.User{}).Where("username = ?", username).Count(&count).Error
	return count > 0, err
}

// ExistsByEmail checks if a user exists by email
func (r *UserRepository) ExistsByEmail(email string) (bool, error) {
	var count int64
	err := db.DB.Model(&entity.User{}).Where("email = ?", email).Count(&count).Error
	return count > 0, err
}
