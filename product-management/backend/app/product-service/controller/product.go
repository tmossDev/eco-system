package controller

import (
	"errors"

	"github.com/kataras/iris/v12"
	"tmossDev.github.com/eco-system/product-management/backend/domain/product/service"
	"tmossDev.github.com/eco-system/shared-components/backend/package/types"
)

const systemUserID uint64 = 1

type ServiceController interface {
	ListProducts() iris.Handler
	ProductDetails() iris.Handler
	CreateProduct() iris.Handler
	UpdateProduct() iris.Handler
	DeleteProduct() iris.Handler
}

type ServiceControllerImp struct {
	productService service.ProductService
}

func NewServiceControllerImp(productService service.ProductService) *ServiceControllerImp {
	return &ServiceControllerImp{
		productService: productService,
	}
}

func (controller *ServiceControllerImp) marshalErrorResponse(ctx iris.Context, err error) {
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

func (controller *ServiceControllerImp) ListProducts() iris.Handler {
	return func(ctx iris.Context) {
		products, err := controller.productService.ListProducts()
		if err != nil {
			controller.marshalErrorResponse(ctx, err)
			return
		}

		_ = ctx.JSON(products)
	}
}

func (controller *ServiceControllerImp) ProductDetails() iris.Handler {
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

func (controller *ServiceControllerImp) CreateProduct() iris.Handler {
	return func(ctx iris.Context) {
		body, err := ctx.GetBody()
		if err != nil {
			controller.marshalErrorResponse(ctx, types.NewInternalServerError())
			return
		}

		product, err := controller.productService.CreateProduct(string(body), systemUserID)
		if err != nil {
			controller.marshalErrorResponse(ctx, err)
			return
		}

		ctx.StatusCode(iris.StatusCreated)
		_ = ctx.JSON(product)
	}
}

func (controller *ServiceControllerImp) UpdateProduct() iris.Handler {
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

		product, err := controller.productService.UpdateProduct(productID, string(body), systemUserID)
		if err != nil {
			controller.marshalErrorResponse(ctx, err)
			return
		}

		_ = ctx.JSON(product)
	}
}

func (controller *ServiceControllerImp) DeleteProduct() iris.Handler {
	return func(ctx iris.Context) {
		productID, err := ctx.Params().GetUint64("id")
		if err != nil {
			controller.marshalErrorResponse(ctx, types.NewInvalidInputError())
			return
		}

		if err := controller.productService.DeleteProduct(productID, systemUserID); err != nil {
			controller.marshalErrorResponse(ctx, err)
			return
		}

		ctx.StatusCode(iris.StatusNoContent)
	}
}
