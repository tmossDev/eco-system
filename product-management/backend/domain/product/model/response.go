package model

type ProductResponse struct {
	ID             uint64 `json:"id"`
	SKU            string `json:"sku"`
	Name           string `json:"name"`
	Description    string `json:"description"`
	Category       string `json:"category"`
	PriceCents     int64  `json:"price_cents"`
	Currency       string `json:"currency"`
	InventoryCount int64  `json:"inventory_count"`
	Status         string `json:"status"`
	CreatedUser    uint64 `json:"created_user"`
	CreatedAt      string `json:"created_at"`
	UpdatedUser    uint64 `json:"updated_user"`
	UpdatedAt      string `json:"updated_at"`
}
