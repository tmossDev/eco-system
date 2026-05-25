package service

import (
	"tmossDev.github.com/eco-system/product-management/backend/domain/promotion/model"
	"tmossDev.github.com/eco-system/product-management/backend/domain/promotion/repository"
	"tmossDev.github.com/eco-system/shared-components/backend/package/types"
	"tmossDev.github.com/eco-system/shared-components/backend/package/validator"
)

type PromotionService interface {
	ListDiscounts() ([]model.Discount, error)
	GetDiscount(discountID uint64) (*model.Discount, error)
	CreateDiscount(body string, creatingUserID uint64) (*model.Discount, error)
	UpdateDiscount(discountID uint64, body string, updatingUserID uint64) (*model.Discount, error)
	DeleteDiscount(discountID uint64, deletingUserID uint64) error
	GetPromotionSettings() (*model.PromotionSettings, error)
	UpdatePromotionSettings(body string, updatingUserID uint64) (*model.PromotionSettings, error)
	Shutdown()
}

type PromotionServiceImpl struct {
	validator     validator.Validator
	discountRepo  repository.DiscountRepository
	promotionRepo repository.PromotionRepository
}

func NewPromotionService(validator validator.Validator, discountRepo repository.DiscountRepository, promotionRepo repository.PromotionRepository) PromotionService {
	return &PromotionServiceImpl{
		validator:     validator,
		discountRepo:  discountRepo,
		promotionRepo: promotionRepo,
	}
}

func (service *PromotionServiceImpl) ListDiscounts() ([]model.Discount, error) {
	return service.discountRepo.ListDiscounts()
}

func (service *PromotionServiceImpl) GetDiscount(discountID uint64) (*model.Discount, error) {
	return service.discountRepo.GetDiscountByID(discountID)
}

func (service *PromotionServiceImpl) CreateDiscount(body string, creatingUserID uint64) (*model.Discount, error) {
	var request model.DiscountRequest
	if err := service.validator.MarshalAndValidateREQ(body, &request); err != nil {
		return nil, err
	}

	if err := validateDiscountRequest(request); err != nil {
		return nil, err
	}

	discount := discountRequestToModel(request)
	discount.CreatedUser = creatingUserID
	if err := service.discountRepo.CreateDiscount(&discount); err != nil {
		return nil, err
	}

	return &discount, nil
}

func (service *PromotionServiceImpl) UpdateDiscount(discountID uint64, body string, updatingUserID uint64) (*model.Discount, error) {
	var request model.DiscountUpdateRequest
	if err := service.validator.MarshalAndValidateREQ(body, &request); err != nil {
		return nil, err
	}

	if err := validateDiscountRequest(request.DiscountRequest); err != nil {
		return nil, err
	}

	discount := discountRequestToModel(request.DiscountRequest)
	discount.ID = discountID
	discount.UpdatedUser = updatingUserID
	if err := service.discountRepo.UpdateDiscount(discount); err != nil {
		return nil, err
	}

	return &discount, nil
}

func (service *PromotionServiceImpl) DeleteDiscount(discountID uint64, deletingUserID uint64) error {
	return service.discountRepo.DeleteDiscount(discountID, deletingUserID)
}

func (service *PromotionServiceImpl) GetPromotionSettings() (*model.PromotionSettings, error) {
	return service.promotionRepo.GetPromotionSettings()
}

func (service *PromotionServiceImpl) UpdatePromotionSettings(body string, updatingUserID uint64) (*model.PromotionSettings, error) {
	var request model.PromotionSettingsRequest
	if err := service.validator.MarshalAndValidateREQ(body, &request); err != nil {
		return nil, err
	}

	settings := model.PromotionSettings{
		PromotionsEnabled: request.PromotionsEnabled,
		UpdatedUser:       updatingUserID,
	}
	if err := service.promotionRepo.UpdatePromotionSettings(settings); err != nil {
		return nil, err
	}

	return &settings, nil
}

func (service *PromotionServiceImpl) Shutdown() {
	service.discountRepo.Shutdown()
	service.promotionRepo.Shutdown()
}

func discountRequestToModel(discount model.DiscountRequest) model.Discount {
	return model.Discount{
		Name:                  discount.Name,
		Description:           discount.Description,
		DiscountType:          discount.DiscountType,
		Scope:                 discount.Scope,
		PercentageBasisPoints: discount.PercentageBasisPoints,
		AmountCents:           discount.AmountCents,
		Currency:              discount.Currency,
		BuyQuantity:           discount.BuyQuantity,
		FreeQuantity:          discount.FreeQuantity,
		MinProductCount:       discount.MinProductCount,
		StartsAt:              discount.StartsAt,
		EndsAt:                discount.EndsAt,
		Status:                discount.Status,
		ProductIDs:            discount.ProductIDs,
	}
}

func validateDiscountRequest(request model.DiscountRequest) error {
	if request.DiscountType == "Percentage" {
		if request.PercentageBasisPoints == nil || request.AmountCents != nil || request.Currency != "" {
			return validatorError()
		}
	}

	if request.DiscountType == "Amount" {
		if request.AmountCents == nil || request.PercentageBasisPoints != nil || request.Currency == "" {
			return validatorError()
		}
	}

	if request.DiscountType == "QuantityBonus" {
		if request.BuyQuantity < 1 || request.FreeQuantity < 1 || request.PercentageBasisPoints != nil || request.AmountCents != nil || request.Currency != "" {
			return validatorError()
		}

		if request.MinProductCount < request.BuyQuantity+request.FreeQuantity {
			return validatorError()
		}
	}

	if request.DiscountType != "QuantityBonus" && (request.BuyQuantity != 0 || request.FreeQuantity != 0) {
		return validatorError()
	}

	if request.Scope == "ProductSet" && len(request.ProductIDs) == 0 {
		return validatorError()
	}

	if request.Scope == "Global" && len(request.ProductIDs) > 0 {
		return validatorError()
	}

	return nil
}

func validatorError() error {
	return types.NewInvalidInputError()
}
