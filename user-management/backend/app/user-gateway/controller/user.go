package controller

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/kataras/iris/v12"
	sharedConstants "tmossDev.github.com/eco-system/shared-components/backend/package/constants"
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
	DashboardSummary() iris.Handler
	ListUsers() iris.Handler
	UserDetails() iris.Handler
	CreateUser() iris.Handler
	UpdateUser() iris.Handler
	DeleteUser() iris.Handler
	GetSettings() iris.Handler
	UpdateSettings() iris.Handler
	Logout() iris.Handler
	UpdatePassword() iris.Handler
	UpdateInfo() iris.Handler
}

type GatewayControllerImp struct {
	publicService  service.PublicUserService
	privateService service.PrivateUserService
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
	ApplicationName    string `json:"applicationName"`
	DefaultRole        string `json:"defaultRole"`
	EmailNotifications bool   `json:"emailNotifications"`
	RequireApproval    bool   `json:"requireApproval"`
}

type adminUserResponse struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	Status string `json:"status"`
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

		requestId := ctx.Values().GetString(sharedConstants.CTXRequestIdKey)
		loginResponse, err := controller.publicService.Login(requestId, string(body))
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

func (controller *GatewayControllerImp) DashboardSummary() iris.Handler {
	return func(ctx iris.Context) {
		_ = ctx.JSON(dashboardSummaryResponse{
			Stats: []dashboardStatResponse{
				{Label: "Total users", Value: "4", Caption: "Seeded local accounts"},
				{Label: "Active users", Value: "4", Caption: "All local users active"},
				{Label: "Pending invites", Value: "0", Caption: "No pending invites"},
				{Label: "Admin users", Value: "1", Caption: "Privileged accounts"},
			},
			RecentActivity: []string{
				"Local beta users were seeded",
				"Admin permissions are available for admin@test.com",
				"Gateway is serving live API responses",
				"System settings were reviewed",
			},
		})
	}
}

func seededAdminUsers() []adminUserResponse {
	return []adminUserResponse{
		{ID: "1", Name: "System Auto", Email: "system@test.com", Role: "System", Status: "Active"},
		{ID: "2", Name: "Jane Doe", Email: "admin@test.com", Role: "Admin", Status: "Active"},
		{ID: "3", Name: "Michael Rogers", Email: "moderator@test.com", Role: "Moderator", Status: "Active"},
		{ID: "4", Name: "Blikkies Blignaut", Email: "blikkies@test.com", Role: "User", Status: "Active"},
	}
}

func (controller *GatewayControllerImp) ListUsers() iris.Handler {
	return func(ctx iris.Context) {
		_ = ctx.JSON(seededAdminUsers())
	}
}

func (controller *GatewayControllerImp) UserDetails() iris.Handler {
	return func(ctx iris.Context) {
		id := ctx.Params().Get("id")
		for _, user := range seededAdminUsers() {
			if user.ID == id {
				_ = ctx.JSON(user)
				return
			}
		}

		controller.marshalErrorResponse(ctx, types.NewNoTFoundOrNoRecordError())
	}
}

func (controller *GatewayControllerImp) CreateUser() iris.Handler {
	return func(ctx iris.Context) {
		var user adminUserResponse
		if err := ctx.ReadJSON(&user); err != nil {
			controller.marshalErrorResponse(ctx, types.NewInvalidInputError())
			return
		}
		if user.ID == "" {
			user.ID = strconv.FormatInt(time.Now().UnixNano(), 10)
		}

		ctx.StatusCode(iris.StatusCreated)
		_ = ctx.JSON(user)
	}
}

func (controller *GatewayControllerImp) UpdateUser() iris.Handler {
	return func(ctx iris.Context) {
		var user adminUserResponse
		if err := ctx.ReadJSON(&user); err != nil {
			controller.marshalErrorResponse(ctx, types.NewInvalidInputError())
			return
		}
		user.ID = ctx.Params().Get("id")
		_ = ctx.JSON(user)
	}
}

func (controller *GatewayControllerImp) DeleteUser() iris.Handler {
	return func(ctx iris.Context) {
		ctx.StatusCode(iris.StatusNoContent)
	}
}

func (controller *GatewayControllerImp) GetSettings() iris.Handler {
	return func(ctx iris.Context) {
		_ = ctx.JSON(applicationSettingsResponse{
			ApplicationName:    "Admin Web App",
			DefaultRole:        "User",
			EmailNotifications: true,
			RequireApproval:    false,
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
