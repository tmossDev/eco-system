package controller

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/kataras/iris/v12"
	"tmossDev.github.com/eco-system/shared-components/backend/package/logger"
	"tmossDev.github.com/eco-system/shared-components/backend/package/types"
	"tmossDev.github.com/eco-system/shared-components/backend/package/utils"
	"tmossDev.github.com/eco-system/user-management/backend/package/user/constants"
	"tmossDev.github.com/eco-system/user-management/backend/package/user/service"
)

type GatewayController interface {
	Register() iris.Handler
	Login() iris.Handler
	User() iris.Handler
	Logout() iris.Handler
	UpdatePassword() iris.Handler
	UpdateInfo() iris.Handler
}

type GatewayControllerImp struct {
	publicService  service.PublicUserService
	privateService service.PrivateUserService
}

func NewGatewayControllerImp(publicService service.PublicUserService, privateService service.PrivateUserService) *GatewayControllerImp {
	return &GatewayControllerImp{
		publicService:  publicService,
		privateService: privateService,
	}
}

func (controller *GatewayControllerImp) getUserIdSession(ctx iris.Context) (uint64, error) {
	jwt, err := controller.getJwtTokenFromSession(ctx)
	if err != nil {
		return 0, err
	}

	issuer, err := utils.GetIssuerFromJwt(jwt, constants.PASSWORD_SECRET_HASHING_KEY)
	if err != nil {
		logger.Errorf("Cant get issuer from jwt '%s'", issuer)
		return 0, types.NewInternalServerError()
	}

	userId, err := strconv.ParseUint(issuer, 10, 64)
	if err != nil {
		logger.Errorf("Cant parse as Uint: '%s'", issuer)
		return 0, types.NewInternalServerError()
	}

	return userId, nil
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

func (controller *GatewayControllerImp) Login() iris.Handler {
	return func(ctx iris.Context) {
		body, err := ctx.GetBody()
		if err != nil {
			controller.marshalErrorResponse(ctx, types.NewInternalServerError())
			return
		}

		loginResponse, err := controller.publicService.Login(string(body))
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

func (controller *GatewayControllerImp) Register() iris.Handler {
	return func(ctx iris.Context) {
		body, err := ctx.GetBody()
		if err != nil {
			controller.marshalErrorResponse(ctx, types.NewInternalServerError())
			return
		}

		loginResponse, err := controller.publicService.Register(string(body))
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

func (controller *GatewayControllerImp) User() iris.Handler {
	return func(ctx iris.Context) {
		jwt, err := controller.getJwtTokenFromSession(ctx)
		if err != nil {
			controller.marshalErrorResponse(ctx, types.NewUnauthorizedError())
			return
		}

		userResponse, err := controller.publicService.User(jwt)
		if err != nil {
			controller.marshalErrorResponse(ctx, err)
			return
		}

		_ = ctx.JSON(userResponse)
	}
}

func (controller *GatewayControllerImp) Logout() iris.Handler {
	return func(ctx iris.Context) {
		jwt, err := controller.getJwtTokenFromSession(ctx)
		if err != nil {
			controller.marshalErrorResponse(ctx, types.NewUnauthorizedError())
			return
		}

		err = controller.publicService.Logout(jwt)
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

func (controller *GatewayControllerImp) UpdateInfo() iris.Handler {
	return func(ctx iris.Context) {
		userId, err := controller.getUserIdSession(ctx)
		if err != nil {
			controller.marshalErrorResponse(ctx, types.NewUnauthorizedError())
			return
		}

		body, err := ctx.GetBody()
		if err != nil {
			controller.marshalErrorResponse(ctx, types.NewInternalServerError())
			return
		}

		userResponse, err := controller.privateService.UpdateUserInfo(userId, string(body), userId)
		if err != nil {
			controller.marshalErrorResponse(ctx, err)
			return
		}

		_ = ctx.JSON(userResponse)
	}
}

func (controller *GatewayControllerImp) UpdatePassword() iris.Handler {
	return func(ctx iris.Context) {
		userId, err := controller.getUserIdSession(ctx)
		if err != nil {
			controller.marshalErrorResponse(ctx, types.NewUnauthorizedError())
			return
		}

		body, err := ctx.GetBody()
		if err != nil {
			controller.marshalErrorResponse(ctx, types.NewInternalServerError())
			return
		}

		updatedResponse, err := controller.privateService.UpdateUserPassword(userId, string(body), userId)
		if err != nil {
			controller.marshalErrorResponse(ctx, err)
			return
		}

		_ = ctx.JSON(updatedResponse)
	}
}
