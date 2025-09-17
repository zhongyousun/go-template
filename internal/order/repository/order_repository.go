package repository

import (
	"database/sql"
	"go-template/internal/order/model"
)

// OrderRepository defines the contract for order data access
type OrderRepository interface {
	GetOrderByID(id int64) (*model.Order, error)
}

type orderRepository struct {
	DB *sql.DB
}

func NewOrderRepository(db *sql.DB) OrderRepository {
	return &orderRepository{DB: db}
}

func (r *orderRepository) GetOrderByID(id int64) (*model.Order, error) {
	var order model.Order
	err := r.DB.QueryRow(
		`SELECT "id", "product", "price", "userId", "createdAt" FROM public."order" 
		WHERE "id" = $1`,
		id,
	).Scan(&order.ID, &order.Product, &order.Price, &order.UserID, &order.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return &order, nil
}
