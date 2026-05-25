package model

type PromotionSettings struct {
	PromotionsEnabled bool   `json:"promotions_enabled"`
	UpdatedUser       uint64 `json:"updated_user"`
	UpdatedAt         string `json:"updated_at"`
}

type PromotionSettingsRequest struct {
	PromotionsEnabled bool `json:"promotions_enabled"`
}
