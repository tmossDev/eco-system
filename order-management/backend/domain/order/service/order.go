package service

import (
	"tmossDev.github.com/eco-system/order-management/backend/domain/order/model"
	"tmossDev.github.com/eco-system/order-management/backend/domain/order/repository"
	"tmossDev.github.com/eco-system/shared-components/backend/package/validator"
)

type OrderService interface {
	ListOrders() ([]model.OrderResponse, error)
	GetOrder(orderID uint64) (*model.OrderResponse, error)
	UpdateStatus(orderID uint64, body string) (*model.OrderResponse, error)
}

type OrderServiceImpl struct {
	validator validator.Validator
	orderRepo repository.OrderRepository
}

func NewOrderService(validator validator.Validator, orderRepo repository.OrderRepository) OrderService {
	return &OrderServiceImpl{validator: validator, orderRepo: orderRepo}
}

func (service *OrderServiceImpl) ListOrders() ([]model.OrderResponse, error) {
	return service.orderRepo.ListOrders()
}

func (service *OrderServiceImpl) GetOrder(orderID uint64) (*model.OrderResponse, error) {
	return service.orderRepo.GetOrder(orderID)
}

func (service *OrderServiceImpl) UpdateStatus(orderID uint64, body string) (*model.OrderResponse, error) {
	var request model.UpdateOrderStatusRequest
	if err := service.validator.MarshalAndValidateREQ(body, &request); err != nil {
		return nil, err
	}

	return service.orderRepo.UpdateStatus(orderID, request.Status)
}
