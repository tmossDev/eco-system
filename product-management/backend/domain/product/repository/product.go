package repository

import "tmossDev.github.com/eco-system/product-management/backend/domain/product/model"

type ProductRepository interface {
	List() ([]model.ProductResponse, error)
	GetByID(productID uint64) (*model.ProductResponse, error)
	Create(product *model.ProductResponse) error
	Update(product model.ProductResponse) error
	Delete(productID uint64, deletingUserID uint64) error
	Shutdown()
}
