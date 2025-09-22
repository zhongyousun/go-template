package repository

import (
	"go-template/internal/order/model"

	"gorm.io/gorm"
)

// OrderRepository defines the contract for order data access using GORM
type OrderRepository struct {
	DB *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{DB: db}
}

func (r *OrderRepository) GetOrderByID(id int64) (*model.Order, error) {
	var order model.Order
	result := r.DB.First(&order, id)
	if result.Error == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return &order, nil
}

func (r *OrderRepository) GetOrdersByUserID(userID int64) ([]*model.Order, error) {
	var orders []*model.Order
	result := r.DB.Where("user_id = ?", userID).Find(&orders)
	if result.Error != nil {
		return nil, result.Error
	}
	return orders, nil
}
