package routes

import (
	"github.com/kataras/iris/v12"
	"tmossDev.github.com/eco-system/order-management/backend/app/order-service/controller"
	"tmossDev.github.com/eco-system/order-management/backend/domain/order/service"
)

func Setup(app *iris.Application, orderService service.OrderService) {
	orderController := controller.NewServiceController(orderService)

	app.Get("/health", orderController.Health())
	app.Get("/api/orders", orderController.ListOrders())
	app.Get("/api/orders/{id:uint64}", orderController.OrderDetails())
	app.Put("/api/orders/{id:uint64}/status", orderController.UpdateOrderStatus())
}
