package repository

import "tmossDev.github.com/eco-system/order-management/backend/domain/order/model"

type OrderRepository interface {
	ListOrders() ([]model.OrderResponse, error)
	GetOrder(orderID uint64) (*model.OrderResponse, error)
	UpdateStatus(orderID uint64, status string) (*model.OrderResponse, error)
}
