package model

type ProductResponse struct {
	ID               uint64         `json:"id"`
	SKU              string         `json:"sku"`
	Name             string         `json:"name"`
	ShortDescription string         `json:"short_description"`
	Description      string         `json:"description"`
	Category         string         `json:"category"`
	PriceCents       int64          `json:"price_cents"`
	Currency         string         `json:"currency"`
	InventoryCount   int64          `json:"inventory_count"`
	Status           string         `json:"status"`
	Photos           []ProductPhoto `json:"photos"`
	Labels           []string       `json:"labels"`
	CreatedUser      uint64         `json:"created_user"`
	CreatedAt        string         `json:"created_at"`
	UpdatedUser      uint64         `json:"updated_user"`
	UpdatedAt        string         `json:"updated_at"`
}

type ProductPhoto struct {
	URL          string `json:"url"`
	ThumbnailURL string `json:"thumbnail_url"`
	AltText      string `json:"alt_text"`
	IsPrimary    bool   `json:"is_primary"`
}
