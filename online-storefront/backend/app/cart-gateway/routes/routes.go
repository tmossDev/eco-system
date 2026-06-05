package routes

import (
	"github.com/kataras/iris/v12"
	"tmossDev.github.com/eco-system/online-storefront/backend/app/cart-gateway/controller"
	cartService "tmossDev.github.com/eco-system/online-storefront/backend/domain/cart/service"
	orderService "tmossDev.github.com/eco-system/online-storefront/backend/domain/order/service"
)

func Setup(app *iris.Application, cartService cartService.CartService, orderService orderService.OrderService) {
	cartController := controller.NewCartController(cartService, orderService)

	app.Get("/health", cartController.Health())
	app.Get("/api/cart", cartController.GetCurrent())
	app.Post("/api/cart/checkout", cartController.Checkout())
	app.Get("/api/orders", cartController.ListOrders())
	app.Post("/api/cart/items", cartController.AddItem())
	app.Put("/api/cart/items/{productID:uint64}", cartController.UpdateItem())
	app.Delete("/api/cart/items/{productID:uint64}", cartController.RemoveItem())
	app.Delete("/api/cart", cartController.Clear())
}
