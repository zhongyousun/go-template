package service

import (
	"go-template/internal/order/model"
	"go-template/internal/order/repository"
)

type OrderService struct {
	Repo repository.OrderRepository
}

func NewOrderService(repo repository.OrderRepository) *OrderService {
	return &OrderService{Repo: repo}
}

func (s *OrderService) GetOrderByID(id int64) (*model.Order, error) {
	return s.Repo.GetOrderByID(id)
}

func (s *OrderService) CreateOrder(order *model.Order) error {
	return s.Repo.CreateOrder(order)
}
