package model

type ProductRequest struct {
	SKU              string                `json:"sku" validate:"required,gt=0,lte=80"`
	Name             string                `json:"name" validate:"required,gt=0,lte=140"`
	ShortDescription string                `json:"short_description" validate:"omitempty,lte=280"`
	Description      string                `json:"description" validate:"omitempty,lte=2000"`
	Category         string                `json:"category" validate:"required,gt=0,lte=120"`
	PriceCents       int64                 `json:"price_cents" validate:"gte=0"`
	Currency         string                `json:"currency" validate:"required,len=3"`
	InventoryCount   int64                 `json:"inventory_count" validate:"gte=0"`
	Status           string                `json:"status" validate:"required,oneof=Draft Active Archived"`
	Photos           []ProductPhotoRequest `json:"photos" validate:"omitempty,dive"`
	Labels           []string              `json:"labels" validate:"omitempty,dive,gt=0,lte=40"`
}

type ProductPhotoRequest struct {
	URL          string `json:"url" validate:"required,gt=0,lte=2000"`
	ThumbnailURL string `json:"thumbnail_url" validate:"omitempty,lte=2000"`
	AltText      string `json:"alt_text" validate:"omitempty,lte=240"`
	IsPrimary    bool   `json:"is_primary"`
}

type ProductUpdateRequest struct {
	ProductRequest
}

type DiscountRequest struct {
	Name                  string   `json:"name" validate:"required,gt=0,lte=140"`
	Description           string   `json:"description" validate:"omitempty,lte=2000"`
	DiscountType          string   `json:"discount_type" validate:"required,oneof=Percentage Amount"`
	Scope                 string   `json:"scope" validate:"required,oneof=Global ProductSet"`
	PercentageBasisPoints *int64   `json:"percentage_basis_points" validate:"omitempty,gte=1,lte=10000"`
	AmountCents           *int64   `json:"amount_cents" validate:"omitempty,gt=0"`
	Currency              string   `json:"currency" validate:"omitempty,len=3"`
	MinProductCount       int64    `json:"min_product_count" validate:"gte=1"`
	StartsAt              string   `json:"starts_at" validate:"omitempty"`
	EndsAt                string   `json:"ends_at" validate:"omitempty"`
	Status                string   `json:"status" validate:"required,oneof=Draft Active Archived"`
	ProductIDs            []uint64 `json:"product_ids" validate:"omitempty,dive,gt=0"`
}

type DiscountUpdateRequest struct {
	DiscountRequest
}
