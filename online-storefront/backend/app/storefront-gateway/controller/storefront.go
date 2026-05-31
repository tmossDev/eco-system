package controller

import (
	"errors"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/kataras/iris/v12"
	"tmossDev.github.com/eco-system/online-storefront/backend/app/storefront-gateway/client"
	"tmossDev.github.com/eco-system/product-management/backend/domain/product/model"
	sharedConstants "tmossDev.github.com/eco-system/shared-components/backend/package/constants"
	"tmossDev.github.com/eco-system/shared-components/backend/package/types"
	userModel "tmossDev.github.com/eco-system/shared-components/backend/package/user/model"
)

const activeProductStatus = "Active"

type GatewayController struct {
	userClient    client.UserClient
	productClient client.ProductClient
}

func NewGatewayController(userClient client.UserClient, productClient client.ProductClient) *GatewayController {
	return &GatewayController{
		userClient:    userClient,
		productClient: productClient,
	}
}

func (controller *GatewayController) Register() iris.Handler {
	return controller.authenticate(controller.userClient.Register)
}

func (controller *GatewayController) Login() iris.Handler {
	return controller.authenticate(controller.userClient.Login)
}

func (controller *GatewayController) authenticate(action func(string, []byte) (*userModel.LoginResponse, error)) iris.Handler {
	return func(ctx iris.Context) {
		body, err := ctx.GetBody()
		if err != nil {
			controller.marshalErrorResponse(ctx, types.NewInternalServerError())
			return
		}

		loginResponse, err := action(controller.requestID(ctx), body)
		if err != nil {
			controller.marshalErrorResponse(ctx, err)
			return
		}

		ctx.SetCookieKV("jwt", loginResponse.Jwt, iris.CookieExpires(24*time.Hour), iris.CookieHTTPOnly(true))
		_ = ctx.JSON(loginResponse)
	}
}

func (controller *GatewayController) Logout() iris.Handler {
	return func(ctx iris.Context) {
		jwt, err := controller.jwt(ctx)
		if err != nil {
			controller.marshalErrorResponse(ctx, types.NewUnauthorizedError())
			return
		}

		if err := controller.userClient.Logout(controller.requestID(ctx), jwt); err != nil {
			controller.marshalErrorResponse(ctx, err)
			return
		}

		ctx.RemoveCookie("jwt")
		_ = ctx.JSON(iris.Map{"message": "logged out"})
	}
}

func (controller *GatewayController) User() iris.Handler {
	return func(ctx iris.Context) {
		jwt, err := controller.jwt(ctx)
		if err != nil {
			controller.marshalErrorResponse(ctx, types.NewUnauthorizedError())
			return
		}

		user, err := controller.userClient.User(controller.requestID(ctx), jwt)
		if err != nil {
			controller.marshalErrorResponse(ctx, err)
			return
		}

		_ = ctx.JSON(user)
	}
}

func (controller *GatewayController) ListProducts() iris.Handler {
	return func(ctx iris.Context) {
		products, err := controller.productClient.ListProducts(controller.requestID(ctx))
		if err != nil {
			controller.marshalErrorResponse(ctx, err)
			return
		}

		activeProducts := make([]model.ProductResponse, 0, len(products))
		for _, product := range products {
			if product.Status == activeProductStatus {
				activeProducts = append(activeProducts, product)
			}
		}

		_ = ctx.JSON(activeProducts)
	}
}

func (controller *GatewayController) ProductDetails() iris.Handler {
	return func(ctx iris.Context) {
		product, err := controller.productClient.GetProduct(controller.requestID(ctx), ctx.Params().Get("id"))
		if err != nil {
			controller.marshalErrorResponse(ctx, err)
			return
		}
		if product.Status != activeProductStatus {
			controller.marshalErrorResponse(ctx, types.NewNoTFoundOrNoRecordError())
			return
		}

		_ = ctx.JSON(product)
	}
}

func (controller *GatewayController) GetProductMedia() iris.Handler {
	return func(ctx iris.Context) {
		media, err := controller.productClient.GetProductMedia(controller.requestID(ctx), ctx.Params().Get("objectKey"))
		if err != nil {
			controller.marshalErrorResponse(ctx, err)
			return
		}
		defer media.Body.Close()

		ctx.Header("Content-Type", media.ContentType)
		ctx.Header("Cache-Control", "public, max-age=31536000, immutable")
		if media.ContentLength > 0 {
			ctx.Header("Content-Length", strconv.FormatInt(media.ContentLength, 10))
		}
		if _, err := io.Copy(ctx.ResponseWriter(), media.Body); err != nil {
			controller.marshalErrorResponse(ctx, types.NewInternalServerError())
		}
	}
}

func (controller *GatewayController) Health() iris.Handler {
	return func(ctx iris.Context) {
		_ = ctx.JSON(iris.Map{"status": "ok"})
	}
}

func (controller *GatewayController) requestID(ctx iris.Context) string {
	return ctx.Values().GetString(sharedConstants.CTXRequestIdKey)
}

func (controller *GatewayController) jwt(ctx iris.Context) (string, error) {
	authHeader := ctx.GetHeader("Authorization")
	if authHeader != "" {
		parts := strings.Split(authHeader, " ")
		if len(parts) == 2 && parts[0] == "Bearer" {
			return parts[1], nil
		}

		return "", errors.New("authorization header format must be 'Bearer {token}'")
	}

	if cookie := ctx.GetCookie("jwt"); cookie != "" {
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
