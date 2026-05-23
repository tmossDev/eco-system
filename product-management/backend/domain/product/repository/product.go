package repository

import "tmossDev.github.com/eco-system/product-management/backend/domain/product/model"

type ProductRepository interface {
	List() ([]model.ProductResponse, error)
	GetByID(productID uint64) (*model.ProductResponse, error)
	Create(product model.ProductRequest, creatingUserID uint64) (*model.ProductResponse, error)
	Update(productID uint64, product model.ProductUpdateRequest, updatingUserID uint64) (*model.ProductResponse, error)
	Delete(productID uint64, deletingUserID uint64) error
	ListDiscounts() ([]model.Discount, error)
	GetDiscountByID(discountID uint64) (*model.Discount, error)
	CreateDiscount(discount model.DiscountRequest, creatingUserID uint64) (*model.Discount, error)
	UpdateDiscount(discountID uint64, discount model.DiscountUpdateRequest, updatingUserID uint64) (*model.Discount, error)
	DeleteDiscount(discountID uint64, deletingUserID uint64) error
	Shutdown()
}
