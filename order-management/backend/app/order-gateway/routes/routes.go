package routes

import (
	"github.com/kataras/iris/v12"
	userClient "tmossDev.github.com/eco-system/order-management/backend/app/order-gateway/client"
	"tmossDev.github.com/eco-system/order-management/backend/app/order-gateway/controller"
	"tmossDev.github.com/eco-system/order-management/backend/domain/order/service"
)

func Setup(app *iris.Application, orderService service.OrderService, authClient userClient.AuthClient) {
	orderController := controller.NewGatewayController(orderService, authClient)

	app.Get("/health", orderController.Health())
	app.Post("/api/auth/login", orderController.Login())
	app.Post("/api/login", orderController.Login())
	app.Post("/api/auth/logout", orderController.Logout())
	app.Put("/api/logout", orderController.Logout())
	app.Get("/api/orders", orderController.ListOrders())
	app.Get("/api/orders/{id:uint64}", orderController.OrderDetails())
	app.Put("/api/orders/{id:uint64}/status", orderController.UpdateOrderStatus())
}
