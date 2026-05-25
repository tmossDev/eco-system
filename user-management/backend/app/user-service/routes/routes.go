package routes

import (
	"github.com/kataras/iris/v12"
	"tmossDev.github.com/eco-system/user-management/backend/app/user-service/controller"
	"tmossDev.github.com/eco-system/user-management/backend/domain/user/service"
)

func Setup(app *iris.Application, userService service.UserService) {

	gatewayController := controller.NewServiceControllerImp(userService)
	// auth routes
	app.Post("/api/auth/login", gatewayController.Login())
	app.Post("/api/login", gatewayController.Login())
	app.Put("/api/login", gatewayController.Login())
	app.Post("/api/auth/logout", gatewayController.Logout())
	app.Put("/api/logout", gatewayController.Logout())
	app.Post("/api/register", gatewayController.Register())
	app.Get("/api/users/me", gatewayController.User())
	app.Put("/api/users/info", gatewayController.UpdateInfo())
	app.Put("/api/users/password", gatewayController.UpdatePassword())
}
