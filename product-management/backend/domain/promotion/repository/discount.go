package repository

import "tmossDev.github.com/eco-system/product-management/backend/domain/promotion/model"

type DiscountRepository interface {
	ListDiscounts() ([]model.Discount, error)
	GetDiscountByID(discountID uint64) (*model.Discount, error)
	CreateDiscount(discount *model.Discount) error
	UpdateDiscount(discount model.Discount) error
	DeleteDiscount(discountID uint64, deletingUserID uint64) error
	Shutdown()
}
