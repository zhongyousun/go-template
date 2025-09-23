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

// UserRepository defines the contract for user data access (interface)
// This allows you to abstract the data layer and easily switch implementations (e.g., GORM, SQL, mock).
type UserRepository interface {
	CreateUser(user *model.User) error
	GetUserByID(id int64) (*model.User, error)
	GetUserByIDWithCache(id int64) (*model.User, error)
	GetUserByEmail(email string) (*model.User, error)
	UpdateUser(user *model.User) error
	DeleteUser(id int64) error
}

// GormUserRepository implements UserRepository using GORM
type GormUserRepository struct {
	DB *gorm.DB
}

// NewUserRepository returns a UserRepository implemented with GORM.
// Pass a *gorm.DB instance to use as the database connection.
func NewUserRepository(db *gorm.DB) UserRepository {
	return &GormUserRepository{DB: db}
}

func (r *GormUserRepository) CreateUser(user *model.User) error {
	result := r.DB.Create(user)
	return result.Error
}

func (r *GormUserRepository) GetUserByID(id int64) (*model.User, error) {
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

func (r *GormUserRepository) GetUserByIDWithCache(id int64) (*model.User, error) {
	cacheKey := fmt.Sprintf("user:%d", id)
	if redisclient.Rdb != nil {
		if val, err := redisclient.Rdb.Get(context.Background(), cacheKey).Result(); err == nil {
			var u model.User
			if err := json.Unmarshal([]byte(val), &u); err == nil {
				fmt.Println("[GetUserByIDWithCache] source: cache")
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
	fmt.Println("[GetUserByIDWithCache] source: db")
	return &u, nil
}

func (r *GormUserRepository) GetUserByEmail(email string) (*model.User, error) {
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

func (r *GormUserRepository) UpdateUser(user *model.User) error {
	result := r.DB.Save(user)
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return result.Error
}

func (r *GormUserRepository) DeleteUser(id int64) error {
	result := r.DB.Delete(&model.User{}, id)
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return result.Error
}
