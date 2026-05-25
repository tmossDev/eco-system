package model

type Discount struct {
	ID                    uint64   `json:"id"`
	Name                  string   `json:"name"`
	Description           string   `json:"description"`
	DiscountType          string   `json:"discount_type"`
	Scope                 string   `json:"scope"`
	PercentageBasisPoints *int64   `json:"percentage_basis_points"`
	AmountCents           *int64   `json:"amount_cents"`
	Currency              string   `json:"currency"`
	BuyQuantity           int64    `json:"buy_quantity"`
	FreeQuantity          int64    `json:"free_quantity"`
	MinProductCount       int64    `json:"min_product_count"`
	StartsAt              string   `json:"starts_at"`
	EndsAt                string   `json:"ends_at"`
	Status                string   `json:"status"`
	ProductIDs            []uint64 `json:"product_ids"`
	CreatedUser           uint64   `json:"created_user"`
	CreatedAt             string   `json:"created_at"`
	UpdatedUser           uint64   `json:"updated_user"`
	UpdatedAt             string   `json:"updated_at"`
}
