package routes

import (
	"github.com/kataras/iris/v12"
	"tmossDev.github.com/eco-system/user-management/backend/app/user-gateway/controller"
	"tmossDev.github.com/eco-system/user-management/backend/package/user/service"
)

func Setup(app *iris.Application, publicUserService service.PublicUserService, privateUserService service.PrivateUserService) {

	gatewayController := controller.NewGatewayControllerImp(publicUserService, privateUserService)
	// auth routes
	app.Post("/api/register", gatewayController.Register())
	app.Put("/api/login", gatewayController.Login())
	app.Put("/api/logout", gatewayController.Logout())

	// private routes
	app.Get("/api/users/me", gatewayController.User())
	app.Put("/api/users/info", gatewayController.UpdateInfo())
	app.Put("/api/users/password", gatewayController.UpdatePassword())
}
