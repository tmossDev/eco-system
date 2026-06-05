package service

import (
	"tmossDev.github.com/eco-system/online-storefront/backend/domain/order/model"
	"tmossDev.github.com/eco-system/online-storefront/backend/domain/order/repository"
	"tmossDev.github.com/eco-system/shared-components/backend/package/constants"
	"tmossDev.github.com/eco-system/shared-components/backend/package/logger"
)

type OrderService interface {
	Checkout(userID uint64) (*model.OrderResponse, error)
	ListOrders(userID uint64) ([]model.OrderResponse, error)
}

type OrderServiceImpl struct {
	orderRepo repository.OrderRepository
}

func NewOrderService(orderRepo repository.OrderRepository) OrderService {
	return &OrderServiceImpl{orderRepo: orderRepo}
}

func (service *OrderServiceImpl) Checkout(userID uint64) (*model.OrderResponse, error) {
	order, err := service.orderRepo.Checkout(userID)
	if err != nil {
		return nil, err
	}
	logger.Infof(constants.DefaultRequestId, "Completed checkout for user %d with order %d", userID, order.ID)
	return order, nil
}

func (service *OrderServiceImpl) ListOrders(userID uint64) ([]model.OrderResponse, error) {
	return service.orderRepo.ListOrders(userID)
}
