package repository

import "tmossDev.github.com/eco-system/product-management/backend/domain/promotion/model"

type PromotionRepository interface {
	GetPromotionSettings() (*model.PromotionSettings, error)
	UpdatePromotionSettings(settings model.PromotionSettings) error
	Shutdown()
}
