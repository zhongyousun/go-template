package model

import (
	"go-template/internal/order/model"
	"time"
)

type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"createdAt"`
	Role      string    `json:"role"`
}

type UserWithOrders struct {
	User   *User          `json:"user"`
	Orders []*model.Order `json:"orders"`
}
