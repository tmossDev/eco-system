package controller

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/kataras/iris/v12"
	"tmossDev.github.com/eco-system/product-management/backend/domain/product/service"
	"tmossDev.github.com/eco-system/shared-components/backend/package/types"
	"tmossDev.github.com/eco-system/shared-components/backend/package/utils"
	userConstants "tmossDev.github.com/eco-system/user-management/backend/domain/user/constants"
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
	GetSettings() iris.Handler
	UpdateSettings() iris.Handler
}

type GatewayControllerImp struct {
	productService service.ProductService
}

type loginRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginResponse struct {
	AccessToken string       `json:"accessToken"`
	Jwt         string       `json:"jwt"`
	ExpireAt    int64        `json:"expire_at"`
	User        authUserInfo `json:"user"`
}

type authUserInfo struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
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

func NewGatewayControllerImp(productService service.ProductService) *GatewayControllerImp {
	return &GatewayControllerImp{
		productService: productService,
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

func (controller *GatewayControllerImp) Login() iris.Handler {
	return func(ctx iris.Context) {
		var request loginRequest
		if err := ctx.ReadJSON(&request); err != nil {
			controller.marshalErrorResponse(ctx, types.NewInvalidInputError())
			return
		}

		email := request.Email
		if email == "" {
			email = request.Username
		}

		if strings.ToLower(strings.TrimSpace(email)) != "admin@example.com" || request.Password != "password" {
			controller.marshalErrorResponse(ctx, types.NewUnauthorizedError())
			return
		}

		token, expireAt, err := utils.GenerateJwt(strconv.FormatUint(systemUserID, 10), userConstants.PASSWORD_SECRET_HASHING_KEY)
		if err != nil {
			controller.marshalErrorResponse(ctx, types.NewInternalServerError())
			return
		}

		ctx.SetCookieKV(
			"jwt",
			token,
			iris.CookieExpires(24*time.Hour),
			iris.CookieHTTPOnly(true),
		)

		_ = ctx.JSON(loginResponse{
			AccessToken: token,
			Jwt:         token,
			ExpireAt:    expireAt,
			User: authUserInfo{
				ID:    strconv.FormatUint(systemUserID, 10),
				Name:  "Product Admin",
				Email: "admin@example.com",
				Role:  "Admin",
			},
		})
	}
}

func (controller *GatewayControllerImp) Logout() iris.Handler {
	return func(ctx iris.Context) {
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
