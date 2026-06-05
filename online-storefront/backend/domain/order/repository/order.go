package repository

import "tmossDev.github.com/eco-system/online-storefront/backend/domain/order/model"

type OrderRepository interface {
	Checkout(userID uint64) (*model.OrderResponse, error)
	ListOrders(userID uint64) ([]model.OrderResponse, error)
}
