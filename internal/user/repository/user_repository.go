package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"go-template/internal/user/model"
	"go-template/pkg/redisclient"
	"time"

	"gorm.io/gorm"
)

// UserRepository defines the contract for user data access using GORM
type UserRepository struct {
	DB *gorm.DB
}

// NewUserRepository creates a new UserRepository instance
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) CreateUser(user *model.User) error {
	result := r.DB.Create(user)
	return result.Error
}

func (r *UserRepository) GetUserByID(id int64) (*model.User, error) {
	var user model.User
	result := r.DB.First(&user, id)
	if result.Error == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (r *UserRepository) GetUserByIDWithCache(id int64) (*model.User, error) {
	cacheKey := fmt.Sprintf("user:%d", id)
	if redisclient.Rdb != nil {
		if val, err := redisclient.Rdb.Get(context.Background(), cacheKey).Result(); err == nil {
			var u model.User
			if err := json.Unmarshal([]byte(val), &u); err == nil {
				return &u, nil // cache hit
			}
		}
	}
	var u model.User
	result := r.DB.First(&u, id)
	if result.Error == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}
	if redisclient.Rdb != nil {
		if bytes, err := json.Marshal(u); err == nil {
			_ = redisclient.Rdb.Set(context.Background(), cacheKey, bytes, 10*time.Minute).Err()
		}
	}
	return &u, nil
}

func (r *UserRepository) GetUserByEmail(email string) (*model.User, error) {
	var user model.User
	result := r.DB.Where("email = ?", email).First(&user)
	if result.Error == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (r *UserRepository) UpdateUser(user *model.User) error {
	result := r.DB.Save(user)
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return result.Error
}

func (r *UserRepository) DeleteUser(id int64) error {
	result := r.DB.Delete(&model.User{}, id)
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return result.Error
}
