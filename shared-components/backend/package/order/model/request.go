package model

type UpdateOrderStatusRequest struct {
	Status string `json:"status" validate:"required"`
}
