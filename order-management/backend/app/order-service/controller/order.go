package controller

import (
	"errors"

	"github.com/kataras/iris/v12"
	"tmossDev.github.com/eco-system/order-management/backend/domain/order/service"
	"tmossDev.github.com/eco-system/shared-components/backend/package/types"
)

type ServiceController struct {
	orderService service.OrderService
}

func NewServiceController(orderService service.OrderService) *ServiceController {
	return &ServiceController{orderService: orderService}
}

func (controller *ServiceController) Health() iris.Handler {
	return func(ctx iris.Context) {
		_ = ctx.JSON(iris.Map{"status": "ok"})
	}
}

func (controller *ServiceController) ListOrders() iris.Handler {
	return func(ctx iris.Context) {
		orders, err := controller.orderService.ListOrders()
		if err != nil {
			controller.marshalErrorResponse(ctx, err)
			return
		}
		_ = ctx.JSON(orders)
	}
}

func (controller *ServiceController) OrderDetails() iris.Handler {
	return func(ctx iris.Context) {
		orderID, err := ctx.Params().GetUint64("id")
		if err != nil {
			controller.marshalErrorResponse(ctx, types.NewInvalidInputError())
			return
		}

		order, err := controller.orderService.GetOrder(orderID)
		if err != nil {
			controller.marshalErrorResponse(ctx, err)
			return
		}
		_ = ctx.JSON(order)
	}
}

func (controller *ServiceController) UpdateOrderStatus() iris.Handler {
	return func(ctx iris.Context) {
		orderID, err := ctx.Params().GetUint64("id")
		if err != nil {
			controller.marshalErrorResponse(ctx, types.NewInvalidInputError())
			return
		}
		body, err := ctx.GetBody()
		if err != nil {
			controller.marshalErrorResponse(ctx, types.NewInternalServerError())
			return
		}

		order, err := controller.orderService.UpdateStatus(orderID, string(body))
		if err != nil {
			controller.marshalErrorResponse(ctx, err)
			return
		}
		_ = ctx.JSON(order)
	}
}

func (controller *ServiceController) marshalErrorResponse(ctx iris.Context, err error) {
	var typedErr *types.SocketError
	if !errors.As(err, &typedErr) {
		ctx.StatusCode(iris.StatusInternalServerError)
		_ = ctx.JSON(iris.Map{"error": "Internal Server Error", "message": err.Error()})
		return
	}

	ctx.StatusCode(typedErr.StatusCode())
	_ = ctx.JSON(iris.Map{"message": typedErr.Error()})
}
