package controller

import (
	"errors"
	"strconv"
	"strings"

	"github.com/kataras/iris/v12"
	"tmossDev.github.com/eco-system/online-storefront/backend/domain/cart/service"
	"tmossDev.github.com/eco-system/shared-components/backend/package/types"
	userConstants "tmossDev.github.com/eco-system/shared-components/backend/package/user/constants"
	"tmossDev.github.com/eco-system/shared-components/backend/package/utils"
)

type CartController struct {
	cartService service.CartService
}

func NewCartController(cartService service.CartService) *CartController {
	return &CartController{cartService: cartService}
}

func (controller *CartController) Health() iris.Handler {
	return func(ctx iris.Context) {
		_ = ctx.JSON(iris.Map{"status": "ok"})
	}
}

func (controller *CartController) GetCurrent() iris.Handler {
	return controller.withUser(func(ctx iris.Context, userID uint64) error {
		cart, err := controller.cartService.GetCurrent(userID)
		if err != nil {
			return err
		}
		return ctx.JSON(cart)
	})
}

func (controller *CartController) AddItem() iris.Handler {
	return controller.withUser(func(ctx iris.Context, userID uint64) error {
		body, err := ctx.GetBody()
		if err != nil {
			return types.NewInternalServerError()
		}
		cart, err := controller.cartService.AddItem(userID, string(body))
		if err != nil {
			return err
		}
		ctx.StatusCode(iris.StatusCreated)
		return ctx.JSON(cart)
	})
}

func (controller *CartController) UpdateItem() iris.Handler {
	return controller.withUser(func(ctx iris.Context, userID uint64) error {
		productID, err := ctx.Params().GetUint64("productID")
		if err != nil {
			return types.NewInvalidInputError()
		}
		body, err := ctx.GetBody()
		if err != nil {
			return types.NewInternalServerError()
		}
		cart, err := controller.cartService.UpdateItem(userID, productID, string(body))
		if err != nil {
			return err
		}
		return ctx.JSON(cart)
	})
}

func (controller *CartController) RemoveItem() iris.Handler {
	return controller.withUser(func(ctx iris.Context, userID uint64) error {
		productID, err := ctx.Params().GetUint64("productID")
		if err != nil {
			return types.NewInvalidInputError()
		}
		cart, err := controller.cartService.RemoveItem(userID, productID)
		if err != nil {
			return err
		}
		return ctx.JSON(cart)
	})
}

func (controller *CartController) Clear() iris.Handler {
	return controller.withUser(func(ctx iris.Context, userID uint64) error {
		cart, err := controller.cartService.Clear(userID)
		if err != nil {
			return err
		}
		return ctx.JSON(cart)
	})
}

func (controller *CartController) withUser(action func(iris.Context, uint64) error) iris.Handler {
	return func(ctx iris.Context) {
		userID, err := controller.userID(ctx)
		if err != nil {
			controller.marshalErrorResponse(ctx, types.NewUnauthorizedError())
			return
		}
		if err := action(ctx, userID); err != nil {
			controller.marshalErrorResponse(ctx, err)
		}
	}
}

func (controller *CartController) userID(ctx iris.Context) (uint64, error) {
	token := ctx.GetCookie("jwt")
	if authorization := ctx.GetHeader("Authorization"); authorization != "" {
		parts := strings.Split(authorization, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return 0, errors.New("invalid authorization header")
		}
		token = parts[1]
	}
	if token == "" {
		return 0, errors.New("jwt token not found")
	}

	issuer, err := utils.GetIssuerFromJwt(token, userConstants.PASSWORD_SECRET_HASHING_KEY)
	if err != nil {
		return 0, err
	}
	return strconv.ParseUint(issuer, 10, 64)
}

func (controller *CartController) marshalErrorResponse(ctx iris.Context, err error) {
	var typedErr *types.SocketError
	if !errors.As(err, &typedErr) {
		ctx.StatusCode(iris.StatusInternalServerError)
		_ = ctx.JSON(iris.Map{"error": "Internal Server Error", "message": err.Error()})
		return
	}
	ctx.StatusCode(typedErr.StatusCode())
	_ = ctx.JSON(iris.Map{"message": typedErr.Error()})
}
