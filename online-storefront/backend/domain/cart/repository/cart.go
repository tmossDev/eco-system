package repository

import "tmossDev.github.com/eco-system/online-storefront/backend/domain/cart/model"

type CartRepository interface {
	GetCurrent(userID uint64) (*model.CartResponse, error)
	AddItem(userID uint64, productID uint64, quantity int64) (*model.CartResponse, error)
	UpdateItem(userID uint64, productID uint64, quantity int64) (*model.CartResponse, error)
	RemoveItem(userID uint64, productID uint64) (*model.CartResponse, error)
	Clear(userID uint64) (*model.CartResponse, error)
}
