package model

import (
	"go-template/internal/order/model"
	"time"
)

// TableName sets the table name for GORM to 'user' (not the default 'users')
func (User) TableName() string {
	return "user"
}

type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"createdAt" gorm:"column:createdAt"`
	Role      string    `json:"role"`
}

type UserWithOrders struct {
	User   *User          `json:"user"`
	Orders []*model.Order `json:"orders"`
}
