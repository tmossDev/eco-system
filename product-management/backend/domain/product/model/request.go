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
