package service

import (
	"io"

	"tmossDev.github.com/eco-system/product-management/backend/domain/product/model"
	"tmossDev.github.com/eco-system/product-management/backend/domain/product/repository"
	"tmossDev.github.com/eco-system/shared-components/backend/package/validator"
)

type ProductService interface {
	ListProducts() ([]model.ProductResponse, error)
	GetProduct(productID uint64) (*model.ProductResponse, error)
	CreateProduct(body string, creatingUserID uint64) (*model.ProductResponse, error)
	UpdateProduct(productID uint64, body string, updatingUserID uint64) (*model.ProductResponse, error)
	DeleteProduct(productID uint64, deletingUserID uint64) error
	UploadProductPhoto(productID uint64, fileName string, body io.Reader, updatingUserID uint64) (*model.ProductResponse, error)
	GetProductMedia(objectKey string) (*repository.ProductMediaObject, error)
	Shutdown()
}

type ProductServiceImpl struct {
	validator   validator.Validator
	productRepo repository.ProductRepository
	mediaStore  repository.ProductMediaStore
}

func NewProductService(validator validator.Validator, productRepo repository.ProductRepository) ProductService {
	return &ProductServiceImpl{
		validator:   validator,
		productRepo: productRepo,
		mediaStore:  repository.NewS3ProductMediaStoreFromEnv(),
	}
}

func (service *ProductServiceImpl) ListProducts() ([]model.ProductResponse, error) {
	return service.productRepo.List()
}

func (service *ProductServiceImpl) GetProduct(productID uint64) (*model.ProductResponse, error) {
	return service.productRepo.GetByID(productID)
}

func (service *ProductServiceImpl) CreateProduct(body string, creatingUserID uint64) (*model.ProductResponse, error) {
	var request model.ProductRequest
	if err := service.validator.MarshalAndValidateREQ(body, &request); err != nil {
		return nil, err
	}

	product := productRequestToResponse(request)
	product.CreatedUser = creatingUserID
	if err := service.productRepo.Create(&product); err != nil {
		return nil, err
	}

	return &product, nil
}

func (service *ProductServiceImpl) UpdateProduct(productID uint64, body string, updatingUserID uint64) (*model.ProductResponse, error) {
	var request model.ProductUpdateRequest
	if err := service.validator.MarshalAndValidateREQ(body, &request); err != nil {
		return nil, err
	}

	product := productRequestToResponse(request.ProductRequest)
	product.ID = productID
	product.UpdatedUser = updatingUserID
	if err := service.productRepo.Update(product); err != nil {
		return nil, err
	}

	return &product, nil
}

func (service *ProductServiceImpl) DeleteProduct(productID uint64, deletingUserID uint64) error {
	return service.productRepo.Delete(productID, deletingUserID)
}

func (service *ProductServiceImpl) UploadProductPhoto(productID uint64, fileName string, body io.Reader, updatingUserID uint64) (*model.ProductResponse, error) {
	product, err := service.productRepo.GetByID(productID)
	if err != nil {
		return nil, err
	}

	thumbnailURL, detailURL, err := service.mediaStore.SaveProductImage(productID, fileName, body)
	if err != nil {
		return nil, err
	}

	photos := product.Photos
	photos = append(photos, model.ProductPhoto{
		URL:          detailURL,
		ThumbnailURL: thumbnailURL,
		AltText:      product.Name,
		IsPrimary:    len(photos) == 0,
	})

	product.Photos = photos

	product.UpdatedUser = updatingUserID
	if err := service.productRepo.Update(*product); err != nil {
		return nil, err
	}

	return product, nil
}

func (service *ProductServiceImpl) GetProductMedia(objectKey string) (*repository.ProductMediaObject, error) {
	return service.mediaStore.GetObject(objectKey)
}

func (service *ProductServiceImpl) Shutdown() {
	service.productRepo.Shutdown()
}

func productRequestToResponse(product model.ProductRequest) model.ProductResponse {
	return model.ProductResponse{
		SKU:              product.SKU,
		Name:             product.Name,
		ShortDescription: product.ShortDescription,
		Description:      product.Description,
		Category:         product.Category,
		PriceCents:       product.PriceCents,
		Currency:         product.Currency,
		InventoryCount:   product.InventoryCount,
		Status:           product.Status,
		Photos:           photosFromRequest(product.Photos),
		Labels:           product.Labels,
	}
}

func photosFromRequest(photos []model.ProductPhotoRequest) []model.ProductPhoto {
	responses := make([]model.ProductPhoto, 0, len(photos))
	for _, photo := range photos {
		responses = append(responses, model.ProductPhoto{
			URL:          photo.URL,
			ThumbnailURL: photo.ThumbnailURL,
			AltText:      photo.AltText,
			IsPrimary:    photo.IsPrimary,
		})
	}

	return responses
}
