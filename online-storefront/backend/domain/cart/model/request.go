package model

type AddItemRequest struct {
	ProductID uint64 `json:"product_id" validate:"required"`
	Quantity  int64  `json:"quantity" validate:"required,gt=0"`
}

type UpdateItemRequest struct {
	Quantity int64 `json:"quantity" validate:"required,gt=0"`
}
