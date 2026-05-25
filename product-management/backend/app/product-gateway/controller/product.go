package controller

import (
	"errors"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/kataras/iris/v12"
	userClient "tmossDev.github.com/eco-system/product-management/backend/app/product-gateway/client"
	"tmossDev.github.com/eco-system/product-management/backend/domain/product/service"
	promotionService "tmossDev.github.com/eco-system/product-management/backend/domain/promotion/service"
	sharedConstants "tmossDev.github.com/eco-system/shared-components/backend/package/constants"
	"tmossDev.github.com/eco-system/shared-components/backend/package/types"
)

const systemUserID uint64 = 1

type GatewayController interface {
	Login() iris.Handler
	Logout() iris.Handler
	DashboardSummary() iris.Handler
	ListProducts() iris.Handler
	ProductDetails() iris.Handler
	CreateProduct() iris.Handler
	UpdateProduct() iris.Handler
	DeleteProduct() iris.Handler
	ListDiscounts() iris.Handler
	DiscountDetails() iris.Handler
	CreateDiscount() iris.Handler
	UpdateDiscount() iris.Handler
	DeleteDiscount() iris.Handler
	GetPromotionSettings() iris.Handler
	UpdatePromotionSettings() iris.Handler
	UploadProductPhoto() iris.Handler
	GetProductMedia() iris.Handler
	GetSettings() iris.Handler
	UpdateSettings() iris.Handler
}

type GatewayControllerImp struct {
	productService   service.ProductService
	promotionService promotionService.PromotionService
	authClient       userClient.AuthClient
}

type dashboardStatResponse struct {
	Label   string `json:"label"`
	Value   string `json:"value"`
	Caption string `json:"caption"`
}

type dashboardSummaryResponse struct {
	Stats          []dashboardStatResponse `json:"stats"`
	RecentActivity []string                `json:"recentActivity"`
}

type applicationSettingsResponse struct {
	ApplicationName     string `json:"applicationName"`
	DefaultCategory     string `json:"defaultCategory"`
	InventoryAlerts     bool   `json:"inventoryAlerts"`
	RequireReview       bool   `json:"requireReview"`
	DefaultCurrencyCode string `json:"defaultCurrencyCode"`
}

func NewGatewayControllerImp(productService service.ProductService, promotionService promotionService.PromotionService, authClient userClient.AuthClient) *GatewayControllerImp {
	return &GatewayControllerImp{
		productService:   productService,
		promotionService: promotionService,
		authClient:       authClient,
	}
}

func (controller *GatewayControllerImp) marshalErrorResponse(ctx iris.Context, err error) {
	var typedErr *types.SocketError
	ok := errors.As(err, &typedErr)
	if !ok {
		ctx.StatusCode(iris.StatusInternalServerError)
		_ = ctx.JSON(iris.Map{
			"error":   "Internal Server Error",
			"message": err.Error(),
		})
		return
	}

	ctx.StatusCode(typedErr.StatusCode())
	_ = ctx.JSON(iris.Map{
		"message": typedErr.Error(),
	})
}

func (controller *GatewayControllerImp) getActorID(ctx iris.Context) uint64 {
	return systemUserID
}

func (controller *GatewayControllerImp) getJwtTokenFromSession(ctx iris.Context) (string, error) {
	authHeader := ctx.GetHeader("Authorization")
	if authHeader != "" {
		parts := strings.Split(authHeader, " ")
		if len(parts) == 2 && parts[0] == "Bearer" {
			return parts[1], nil
		}

		return "", errors.New("authorization header format must be 'Bearer {token}'")
	}

	cookie := ctx.GetCookie("jwt")
	if cookie != "" {
		return cookie, nil
	}

	return "", errors.New("jwt token not found in headers or cookies")
}

func (controller *GatewayControllerImp) Login() iris.Handler {
	return func(ctx iris.Context) {
		body, err := ctx.GetBody()
		if err != nil {
			controller.marshalErrorResponse(ctx, types.NewInternalServerError())
			return
		}

		requestID := ctx.Values().GetString(sharedConstants.CTXRequestIdKey)
		loginResponse, err := controller.authClient.Login(requestID, body)
		if err != nil {
			controller.marshalErrorResponse(ctx, err)
			return
		}

		ctx.SetCookieKV(
			"jwt",
			loginResponse.Jwt,
			iris.CookieExpires(24*time.Hour),
			iris.CookieHTTPOnly(true),
		)

		_ = ctx.JSON(loginResponse)
	}
}

func (controller *GatewayControllerImp) Logout() iris.Handler {
	return func(ctx iris.Context) {
		jwt, err := controller.getJwtTokenFromSession(ctx)
		if err != nil {
			controller.marshalErrorResponse(ctx, types.NewUnauthorizedError())
			return
		}

		requestID := ctx.Values().GetString(sharedConstants.CTXRequestIdKey)
		err = controller.authClient.Logout(requestID, jwt)
		if err != nil {
			controller.marshalErrorResponse(ctx, err)
			return
		}

		ctx.RemoveCookie("jwt")
		ctx.StatusCode(iris.StatusOK)
		_ = ctx.JSON(iris.Map{
			"message": "logged out",
		})
	}
}

func (controller *GatewayControllerImp) DashboardSummary() iris.Handler {
	return func(ctx iris.Context) {
		products, err := controller.productService.ListProducts()
		if err != nil {
			controller.marshalErrorResponse(ctx, err)
			return
		}

		activeCount := 0
		inventoryCount := int64(0)
		for _, product := range products {
			if product.Status == "Active" {
				activeCount++
			}
			inventoryCount += product.InventoryCount
		}

		_ = ctx.JSON(dashboardSummaryResponse{
			Stats: []dashboardStatResponse{
				{Label: "Total products", Value: strconv.Itoa(len(products)), Caption: "Catalog records"},
				{Label: "Active products", Value: strconv.Itoa(activeCount), Caption: "Visible for sale"},
				{Label: "Inventory units", Value: strconv.FormatInt(inventoryCount, 10), Caption: "Available stock"},
				{Label: "Draft products", Value: strconv.Itoa(len(products) - activeCount), Caption: "Need review"},
			},
			RecentActivity: []string{
				"Product catalog is ready for administration",
				"Inventory tracking is available for all products",
				"Draft, active, and archived catalog states are supported",
			},
		})
	}
}

func (controller *GatewayControllerImp) ListProducts() iris.Handler {
	return func(ctx iris.Context) {
		products, err := controller.productService.ListProducts()
		if err != nil {
			controller.marshalErrorResponse(ctx, err)
			return
		}

		_ = ctx.JSON(products)
	}
}

func (controller *GatewayControllerImp) ProductDetails() iris.Handler {
	return func(ctx iris.Context) {
		productID, err := ctx.Params().GetUint64("id")
		if err != nil {
			controller.marshalErrorResponse(ctx, types.NewInvalidInputError())
			return
		}

		product, err := controller.productService.GetProduct(productID)
		if err != nil {
			controller.marshalErrorResponse(ctx, err)
			return
		}

		_ = ctx.JSON(product)
	}
}

func (controller *GatewayControllerImp) CreateProduct() iris.Handler {
	return func(ctx iris.Context) {
		body, err := ctx.GetBody()
		if err != nil {
			controller.marshalErrorResponse(ctx, types.NewInternalServerError())
			return
		}

		product, err := controller.productService.CreateProduct(string(body), controller.getActorID(ctx))
		if err != nil {
			controller.marshalErrorResponse(ctx, err)
			return
		}

		ctx.StatusCode(iris.StatusCreated)
		_ = ctx.JSON(product)
	}
}

func (controller *GatewayControllerImp) UpdateProduct() iris.Handler {
	return func(ctx iris.Context) {
		productID, err := ctx.Params().GetUint64("id")
		if err != nil {
			controller.marshalErrorResponse(ctx, types.NewInvalidInputError())
			return
		}

		body, err := ctx.GetBody()
		if err != nil {
			controller.marshalErrorResponse(ctx, types.NewInternalServerError())
			return
		}

		product, err := controller.productService.UpdateProduct(productID, string(body), controller.getActorID(ctx))
		if err != nil {
			controller.marshalErrorResponse(ctx, err)
			return
		}

		_ = ctx.JSON(product)
	}
}

func (controller *GatewayControllerImp) DeleteProduct() iris.Handler {
	return func(ctx iris.Context) {
		productID, err := ctx.Params().GetUint64("id")
		if err != nil {
			controller.marshalErrorResponse(ctx, types.NewInvalidInputError())
			return
		}

		if err := controller.productService.DeleteProduct(productID, controller.getActorID(ctx)); err != nil {
			controller.marshalErrorResponse(ctx, err)
			return
		}

		ctx.StatusCode(iris.StatusNoContent)
	}
}

func (controller *GatewayControllerImp) ListDiscounts() iris.Handler {
	return func(ctx iris.Context) {
		discounts, err := controller.promotionService.ListDiscounts()
		if err != nil {
			controller.marshalErrorResponse(ctx, err)
			return
		}

		_ = ctx.JSON(discounts)
	}
}

func (controller *GatewayControllerImp) DiscountDetails() iris.Handler {
	return func(ctx iris.Context) {
		discountID, err := ctx.Params().GetUint64("id")
		if err != nil {
			controller.marshalErrorResponse(ctx, types.NewInvalidInputError())
			return
		}

		discount, err := controller.promotionService.GetDiscount(discountID)
		if err != nil {
			controller.marshalErrorResponse(ctx, err)
			return
		}

		_ = ctx.JSON(discount)
	}
}

func (controller *GatewayControllerImp) CreateDiscount() iris.Handler {
	return func(ctx iris.Context) {
		body, err := ctx.GetBody()
		if err != nil {
			controller.marshalErrorResponse(ctx, types.NewInternalServerError())
			return
		}

		discount, err := controller.promotionService.CreateDiscount(string(body), controller.getActorID(ctx))
		if err != nil {
			controller.marshalErrorResponse(ctx, err)
			return
		}

		ctx.StatusCode(iris.StatusCreated)
		_ = ctx.JSON(discount)
	}
}

func (controller *GatewayControllerImp) UpdateDiscount() iris.Handler {
	return func(ctx iris.Context) {
		discountID, err := ctx.Params().GetUint64("id")
		if err != nil {
			controller.marshalErrorResponse(ctx, types.NewInvalidInputError())
			return
		}

		body, err := ctx.GetBody()
		if err != nil {
			controller.marshalErrorResponse(ctx, types.NewInternalServerError())
			return
		}

		discount, err := controller.promotionService.UpdateDiscount(discountID, string(body), controller.getActorID(ctx))
		if err != nil {
			controller.marshalErrorResponse(ctx, err)
			return
		}

		_ = ctx.JSON(discount)
	}
}

func (controller *GatewayControllerImp) DeleteDiscount() iris.Handler {
	return func(ctx iris.Context) {
		discountID, err := ctx.Params().GetUint64("id")
		if err != nil {
			controller.marshalErrorResponse(ctx, types.NewInvalidInputError())
			return
		}

		if err := controller.promotionService.DeleteDiscount(discountID, controller.getActorID(ctx)); err != nil {
			controller.marshalErrorResponse(ctx, err)
			return
		}

		ctx.StatusCode(iris.StatusNoContent)
	}
}

func (controller *GatewayControllerImp) GetPromotionSettings() iris.Handler {
	return func(ctx iris.Context) {
		settings, err := controller.promotionService.GetPromotionSettings()
		if err != nil {
			controller.marshalErrorResponse(ctx, err)
			return
		}

		_ = ctx.JSON(settings)
	}
}

func (controller *GatewayControllerImp) UpdatePromotionSettings() iris.Handler {
	return func(ctx iris.Context) {
		body, err := ctx.GetBody()
		if err != nil {
			controller.marshalErrorResponse(ctx, types.NewInternalServerError())
			return
		}

		settings, err := controller.promotionService.UpdatePromotionSettings(string(body), controller.getActorID(ctx))
		if err != nil {
			controller.marshalErrorResponse(ctx, err)
			return
		}

		_ = ctx.JSON(settings)
	}
}

func (controller *GatewayControllerImp) UploadProductPhoto() iris.Handler {
	return func(ctx iris.Context) {
		productID, err := ctx.Params().GetUint64("id")
		if err != nil {
			controller.marshalErrorResponse(ctx, types.NewInvalidInputError())
			return
		}

		file, fileHeader, err := ctx.FormFile("file")
		if err != nil {
			controller.marshalErrorResponse(ctx, types.NewInvalidInputError())
			return
		}
		defer file.Close()

		product, err := controller.productService.UploadProductPhoto(
			productID,
			fileHeader.Filename,
			file,
			controller.getActorID(ctx),
		)
		if err != nil {
			controller.marshalErrorResponse(ctx, err)
			return
		}

		ctx.StatusCode(iris.StatusCreated)
		_ = ctx.JSON(product)
	}
}

func (controller *GatewayControllerImp) GetProductMedia() iris.Handler {
	return func(ctx iris.Context) {
		objectKey := ctx.Params().Get("objectKey")
		mediaObject, err := controller.productService.GetProductMedia(objectKey)
		if err != nil {
			controller.marshalErrorResponse(ctx, err)
			return
		}
		defer mediaObject.Body.Close()

		ctx.Header("Content-Type", mediaObject.ContentType)
		ctx.Header("Cache-Control", "public, max-age=31536000, immutable")
		if mediaObject.ContentLength > 0 {
			ctx.Header("Content-Length", strconv.FormatInt(mediaObject.ContentLength, 10))
		}
		if _, err := io.Copy(ctx.ResponseWriter(), mediaObject.Body); err != nil {
			controller.marshalErrorResponse(ctx, types.NewInternalServerError())
			return
		}
	}
}

func (controller *GatewayControllerImp) GetSettings() iris.Handler {
	return func(ctx iris.Context) {
		_ = ctx.JSON(applicationSettingsResponse{
			ApplicationName:     "Product Admin Web App",
			DefaultCategory:     "General",
			InventoryAlerts:     true,
			RequireReview:       false,
			DefaultCurrencyCode: "USD",
		})
	}
}

func (controller *GatewayControllerImp) UpdateSettings() iris.Handler {
	return func(ctx iris.Context) {
		var settings applicationSettingsResponse
		if err := ctx.ReadJSON(&settings); err != nil {
			controller.marshalErrorResponse(ctx, types.NewInvalidInputError())
			return
		}

		_ = ctx.JSON(settings)
	}
}
