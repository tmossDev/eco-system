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
	GetProductMedia(objectKey string) (*ProductMediaObject, error)
	Shutdown()
}

type ProductServiceImpl struct {
	validator   validator.Validator
	productRepo repository.ProductRepository
	mediaStore  ProductMediaStore
}

func NewProductService(validator validator.Validator, productRepo repository.ProductRepository) ProductService {
	return &ProductServiceImpl{
		validator:   validator,
		productRepo: productRepo,
		mediaStore:  NewS3ProductMediaStoreFromEnv(),
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

	return service.productRepo.Create(request, creatingUserID)
}

func (service *ProductServiceImpl) UpdateProduct(productID uint64, body string, updatingUserID uint64) (*model.ProductResponse, error) {
	var request model.ProductUpdateRequest
	if err := service.validator.MarshalAndValidateREQ(body, &request); err != nil {
		return nil, err
	}

	return service.productRepo.Update(productID, request, updatingUserID)
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

	request := productResponseToUpdateRequest(*product)
	request.Photos = photosToRequest(photos)

	return service.productRepo.Update(productID, request, updatingUserID)
}

func (service *ProductServiceImpl) GetProductMedia(objectKey string) (*ProductMediaObject, error) {
	return service.mediaStore.GetObject(objectKey)
}

func (service *ProductServiceImpl) Shutdown() {
	service.productRepo.Shutdown()
}

func productResponseToUpdateRequest(product model.ProductResponse) model.ProductUpdateRequest {
	return model.ProductUpdateRequest{
		ProductRequest: model.ProductRequest{
			SKU:              product.SKU,
			Name:             product.Name,
			ShortDescription: product.ShortDescription,
			Description:      product.Description,
			Category:         product.Category,
			PriceCents:       product.PriceCents,
			Currency:         product.Currency,
			InventoryCount:   product.InventoryCount,
			Status:           product.Status,
			Photos:           photosToRequest(product.Photos),
		},
	}
}

func photosToRequest(photos []model.ProductPhoto) []model.ProductPhotoRequest {
	requests := make([]model.ProductPhotoRequest, 0, len(photos))
	for _, photo := range photos {
		requests = append(requests, model.ProductPhotoRequest{
			URL:          photo.URL,
			ThumbnailURL: photo.ThumbnailURL,
			AltText:      photo.AltText,
			IsPrimary:    photo.IsPrimary,
		})
	}

	return requests
}
