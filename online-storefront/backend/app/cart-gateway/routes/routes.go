package routes

import (
	"github.com/kataras/iris/v12"
	"tmossDev.github.com/eco-system/online-storefront/backend/app/cart-gateway/controller"
	"tmossDev.github.com/eco-system/online-storefront/backend/domain/cart/service"
)

func Setup(app *iris.Application, cartService service.CartService) {
	cartController := controller.NewCartController(cartService)

	app.Get("/health", cartController.Health())
	app.Get("/api/cart", cartController.GetCurrent())
	app.Post("/api/cart/items", cartController.AddItem())
	app.Put("/api/cart/items/{productID:uint64}", cartController.UpdateItem())
	app.Delete("/api/cart/items/{productID:uint64}", cartController.RemoveItem())
	app.Delete("/api/cart", cartController.Clear())
}
