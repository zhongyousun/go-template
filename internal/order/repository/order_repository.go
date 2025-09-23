package repository

import (
	"go-template/internal/order/model"

	"gorm.io/gorm"
)

// OrderRepository defines the contract for order data access.
// This interface allows you to abstract the data layer and easily switch implementations (e.g., GORM, SQL, mock).
type OrderRepository interface {
	// GetOrderByID returns an order by its ID. Returns nil if not found.
	GetOrderByID(id int64) (*model.Order, error)
	// GetOrdersByUserID returns all orders for a given user ID.
	GetOrdersByUserID(userID int64) ([]*model.Order, error)
	// CreateOrder creates a new order in the database.
	CreateOrder(order *model.Order) error
}

// GormOrderRepository is a GORM-based implementation of the OrderRepository interface.
// It holds a *gorm.DB instance and provides methods to access order data using GORM ORM.
type GormOrderRepository struct {
	DB *gorm.DB
}

// NewOrderRepository returns an OrderRepository implemented with GORM.
// Pass a *gorm.DB instance to use as the database connection.
// This allows you to depend on the interface rather than a concrete implementation.
func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &GormOrderRepository{DB: db}
}

func (r *GormOrderRepository) GetOrderByID(id int64) (*model.Order, error) {
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

func (r *GormOrderRepository) GetOrdersByUserID(userID int64) ([]*model.Order, error) {
	var orders []*model.Order
	result := r.DB.Where("user_id = ?", userID).Find(&orders)
	if result.Error != nil {
		return nil, result.Error
	}
	return orders, nil
}

// CreateOrder inserts a new order into the database using GORM
func (r *GormOrderRepository) CreateOrder(order *model.Order) error {
	result := r.DB.Create(order)
	return result.Error
}
