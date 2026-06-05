package controller

import (
	"errors"
	"strings"
	"time"

	"github.com/kataras/iris/v12"
	userClient "tmossDev.github.com/eco-system/order-management/backend/app/order-gateway/client"
	"tmossDev.github.com/eco-system/order-management/backend/domain/order/service"
	sharedConstants "tmossDev.github.com/eco-system/shared-components/backend/package/constants"
	"tmossDev.github.com/eco-system/shared-components/backend/package/types"
)

type GatewayController struct {
	orderService service.OrderService
	authClient   userClient.AuthClient
}

func NewGatewayController(orderService service.OrderService, authClient userClient.AuthClient) *GatewayController {
	return &GatewayController{orderService: orderService, authClient: authClient}
}

func (controller *GatewayController) Health() iris.Handler {
	return func(ctx iris.Context) {
		_ = ctx.JSON(iris.Map{"status": "ok"})
	}
}

func (controller *GatewayController) Login() iris.Handler {
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

func (controller *GatewayController) Logout() iris.Handler {
	return func(ctx iris.Context) {
		jwt, err := controller.getJwtTokenFromSession(ctx)
		if err != nil {
			controller.marshalErrorResponse(ctx, types.NewUnauthorizedError())
			return
		}

		requestID := ctx.Values().GetString(sharedConstants.CTXRequestIdKey)
		if err := controller.authClient.Logout(requestID, jwt); err != nil {
			controller.marshalErrorResponse(ctx, err)
			return
		}

		ctx.RemoveCookie("jwt")
		ctx.StatusCode(iris.StatusOK)
		_ = ctx.JSON(iris.Map{"message": "logged out"})
	}
}

func (controller *GatewayController) ListOrders() iris.Handler {
	return func(ctx iris.Context) {
		orders, err := controller.orderService.ListOrders()
		if err != nil {
			controller.marshalErrorResponse(ctx, err)
			return
		}
		_ = ctx.JSON(orders)
	}
}

func (controller *GatewayController) OrderDetails() iris.Handler {
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

func (controller *GatewayController) UpdateOrderStatus() iris.Handler {
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

func (controller *GatewayController) getJwtTokenFromSession(ctx iris.Context) (string, error) {
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

func (controller *GatewayController) marshalErrorResponse(ctx iris.Context, err error) {
	var typedErr *types.SocketError
	if !errors.As(err, &typedErr) {
		ctx.StatusCode(iris.StatusInternalServerError)
		_ = ctx.JSON(iris.Map{"error": "Internal Server Error", "message": err.Error()})
		return
	}

	ctx.StatusCode(typedErr.StatusCode())
	_ = ctx.JSON(iris.Map{"message": typedErr.Error()})
}
