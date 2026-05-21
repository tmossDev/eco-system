package routes

import (
	"github.com/kataras/iris/v12"
	"tmossDev.github.com/eco-system/user-management/backend/app/user-gateway/controller"
	"tmossDev.github.com/eco-system/user-management/backend/domain/user/service"
)

func Setup(app *iris.Application, publicUserService service.PublicUserService, privateUserService service.PrivateUserService) {

	gatewayController := controller.NewGatewayControllerImp(publicUserService, privateUserService)
	// auth routes
	app.Post("/api/register", gatewayController.Register())
	app.Post("/api/auth/login", gatewayController.Login())
	app.Post("/api/login", gatewayController.Login())
	app.Put("/api/login", gatewayController.Login())
	app.Post("/api/auth/logout", gatewayController.Logout())
	app.Put("/api/logout", gatewayController.Logout())

	// admin app routes
	app.Get("/api/dashboard/summary", gatewayController.DashboardSummary())
	app.Get("/api/settings", gatewayController.GetSettings())
	app.Put("/api/settings", gatewayController.UpdateSettings())
	app.Get("/api/users", gatewayController.ListUsers())
	app.Post("/api/users", gatewayController.CreateUser())
	app.Get("/api/users/{id:uint64}", gatewayController.UserDetails())
	app.Put("/api/users/{id:uint64}", gatewayController.UpdateUser())
	app.Delete("/api/users/{id:uint64}", gatewayController.DeleteUser())

	// private routes
	app.Get("/api/users/me", gatewayController.User())
	app.Put("/api/users/info", gatewayController.UpdateInfo())
	app.Put("/api/users/password", gatewayController.UpdatePassword())
}
