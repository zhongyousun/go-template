package model

import "time"

type Order struct {
	ID        int64       `json:"id"`
	Product   string      `json:"product"`
	Price     float64     `json:"price"`
	UserID    int64       `json:"user_id"`
	Amount    float64     `json:"amount"`
	Status    string      `json:"status"`
	CreatedAt time.Time   `json:"created_at"`
	Items     []OrderItem `json:"items"`
}

type OrderItem struct {
	ProductID int64   `json:"product_id"`
	Name      string  `json:"name"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}
