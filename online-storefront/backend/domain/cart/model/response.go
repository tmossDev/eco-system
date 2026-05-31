package model

type CartResponse struct {
	ID            uint64     `json:"id"`
	UserID        uint64     `json:"user_id"`
	Items         []CartItem `json:"items"`
	ItemCount     int64      `json:"item_count"`
	SubtotalCents int64      `json:"subtotal_cents"`
	Currency      string     `json:"currency"`
	CreatedAt     string     `json:"created_at"`
	UpdatedAt     string     `json:"updated_at"`
}

type CartItem struct {
	ProductID    uint64 `json:"product_id"`
	SKU          string `json:"sku"`
	Name         string `json:"name"`
	Quantity     int64  `json:"quantity"`
	PriceCents   int64  `json:"price_cents"`
	Currency     string `json:"currency"`
	LineTotal    int64  `json:"line_total_cents"`
	ThumbnailURL string `json:"thumbnail_url,omitempty"`
}
