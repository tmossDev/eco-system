package model

type DiscountRequest struct {
	Name                  string   `json:"name" validate:"required,gt=0,lte=140"`
	Description           string   `json:"description" validate:"omitempty,lte=2000"`
	DiscountType          string   `json:"discount_type" validate:"required,oneof=Percentage Amount QuantityBonus"`
	Scope                 string   `json:"scope" validate:"required,oneof=Global ProductSet"`
	PercentageBasisPoints *int64   `json:"percentage_basis_points" validate:"omitempty,gte=1,lte=10000"`
	AmountCents           *int64   `json:"amount_cents" validate:"omitempty,gt=0"`
	Currency              string   `json:"currency" validate:"omitempty,len=3"`
	BuyQuantity           int64    `json:"buy_quantity" validate:"omitempty,gte=0"`
	FreeQuantity          int64    `json:"free_quantity" validate:"omitempty,gte=0"`
	MinProductCount       int64    `json:"min_product_count" validate:"gte=1"`
	StartsAt              string   `json:"starts_at" validate:"omitempty"`
	EndsAt                string   `json:"ends_at" validate:"omitempty"`
	Status                string   `json:"status" validate:"required,oneof=Draft Active Archived"`
	ProductIDs            []uint64 `json:"product_ids" validate:"omitempty,dive,gt=0"`
}

type DiscountUpdateRequest struct {
	DiscountRequest
}
