package routes

import (
	"github.com/kataras/iris/v12"
	"tmossDev.github.com/eco-system/user-management/backend/app/user-service/controller"
	"tmossDev.github.com/eco-system/user-management/backend/domain/user/service"
)

func Setup(app *iris.Application, publicUserService service.PublicUserService, privateUserService service.PrivateUserService) {

	gatewayController := controller.NewServiceControllerImp(publicUserService, privateUserService)
	// auth routes
	app.Post("/api/register", gatewayController.Register())
	app.Get("/api/users/me", gatewayController.User())
	app.Put("/api/users/info", gatewayController.UpdateInfo())
	app.Put("/api/users/password", gatewayController.UpdatePassword())
}
