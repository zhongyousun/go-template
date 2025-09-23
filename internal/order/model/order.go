package model

import "time"

// TableName sets the table name for GORM to 'order' (not the default 'orders')
func (Order) TableName() string {
	return "order"
}

type Order struct {
	ID        int64     `json:"id"`
	Product   string    `json:"product"`
	Price     float64   `json:"price"`
	UserID    int64     `json:"user_id" gorm:"column:userId"`
	CreatedAt time.Time `json:"created_at" gorm:"column:createdAt"`
	// Items is a slice of order items, ignored by GORM (not stored in DB)
	Items []OrderItem `json:"items" gorm:"-"`
}

type OrderItem struct {
	ProductID int64   `json:"product_id"`
	Name      string  `json:"name"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}
