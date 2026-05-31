package service

import (
	"tmossDev.github.com/eco-system/online-storefront/backend/domain/cart/model"
	"tmossDev.github.com/eco-system/online-storefront/backend/domain/cart/repository"
	"tmossDev.github.com/eco-system/shared-components/backend/package/validator"
)

type CartService interface {
	GetCurrent(userID uint64) (*model.CartResponse, error)
	AddItem(userID uint64, body string) (*model.CartResponse, error)
	UpdateItem(userID uint64, productID uint64, body string) (*model.CartResponse, error)
	RemoveItem(userID uint64, productID uint64) (*model.CartResponse, error)
	Clear(userID uint64) (*model.CartResponse, error)
}

type CartServiceImpl struct {
	validator validator.Validator
	cartRepo  repository.CartRepository
}

func NewCartService(validator validator.Validator, cartRepo repository.CartRepository) CartService {
	return &CartServiceImpl{validator: validator, cartRepo: cartRepo}
}

func (service *CartServiceImpl) GetCurrent(userID uint64) (*model.CartResponse, error) {
	return service.cartRepo.GetCurrent(userID)
}

func (service *CartServiceImpl) AddItem(userID uint64, body string) (*model.CartResponse, error) {
	var request model.AddItemRequest
	if err := service.validator.MarshalAndValidateREQ(body, &request); err != nil {
		return nil, err
	}

	return service.cartRepo.AddItem(userID, request.ProductID, request.Quantity)
}

func (service *CartServiceImpl) UpdateItem(userID uint64, productID uint64, body string) (*model.CartResponse, error) {
	var request model.UpdateItemRequest
	if err := service.validator.MarshalAndValidateREQ(body, &request); err != nil {
		return nil, err
	}

	return service.cartRepo.UpdateItem(userID, productID, request.Quantity)
}

func (service *CartServiceImpl) RemoveItem(userID uint64, productID uint64) (*model.CartResponse, error) {
	return service.cartRepo.RemoveItem(userID, productID)
}

func (service *CartServiceImpl) Clear(userID uint64) (*model.CartResponse, error) {
	return service.cartRepo.Clear(userID)
}
