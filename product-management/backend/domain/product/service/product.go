package service

import (
	"tmossDev.github.com/eco-system/product-management/backend/domain/product/model"
	"tmossDev.github.com/eco-system/product-management/backend/domain/product/repository"
	"tmossDev.github.com/eco-system/shared-components/backend/package/validator"
)

type ProductService interface {
	ListProducts() ([]model.ProductResponse, error)
	GetProduct(productID uint64) (*model.ProductResponse, error)
	CreateProduct(body string, creatingUserID uint64) (*model.ProductResponse, error)
	UpdateProduct(productID uint64, body string, updatingUserID uint64) (*model.ProductResponse, error)
	DeleteProduct(productID uint64, deletingUserID uint64) error
	Shutdown()
}

type ProductServiceImpl struct {
	validator   validator.Validator
	productRepo repository.ProductRepository
}

func NewProductService(validator validator.Validator, productRepo repository.ProductRepository) ProductService {
	return &ProductServiceImpl{
		validator:   validator,
		productRepo: productRepo,
	}
}

func (service *ProductServiceImpl) ListProducts() ([]model.ProductResponse, error) {
	return service.productRepo.List()
}

func (service *ProductServiceImpl) GetProduct(productID uint64) (*model.ProductResponse, error) {
	return service.productRepo.GetByID(productID)
}

func (service *ProductServiceImpl) CreateProduct(body string, creatingUserID uint64) (*model.ProductResponse, error) {
	var request model.ProductRequest
	if err := service.validator.MarshalAndValidateREQ(body, &request); err != nil {
		return nil, err
	}

	return service.productRepo.Create(request, creatingUserID)
}

func (service *ProductServiceImpl) UpdateProduct(productID uint64, body string, updatingUserID uint64) (*model.ProductResponse, error) {
	var request model.ProductUpdateRequest
	if err := service.validator.MarshalAndValidateREQ(body, &request); err != nil {
		return nil, err
	}

	return service.productRepo.Update(productID, request, updatingUserID)
}

func (service *ProductServiceImpl) DeleteProduct(productID uint64, deletingUserID uint64) error {
	return service.productRepo.Delete(productID, deletingUserID)
}

func (service *ProductServiceImpl) Shutdown() {
	service.productRepo.Shutdown()
}
