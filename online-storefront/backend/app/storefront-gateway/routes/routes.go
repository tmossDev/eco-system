package routes

import (
	"github.com/kataras/iris/v12"
	"tmossDev.github.com/eco-system/online-storefront/backend/app/storefront-gateway/client"
	"tmossDev.github.com/eco-system/online-storefront/backend/app/storefront-gateway/controller"
)

func Setup(app *iris.Application, userClient client.UserClient, productClient client.ProductClient) {
	gatewayController := controller.NewGatewayController(userClient, productClient)

	app.Get("/health", gatewayController.Health())
	app.Post("/api/register", gatewayController.Register())
	app.Post("/api/auth/login", gatewayController.Login())
	app.Post("/api/auth/logout", gatewayController.Logout())
	app.Get("/api/users/me", gatewayController.User())
	app.Get("/api/products", gatewayController.ListProducts())
	app.Get("/api/products/{id:uint64}", gatewayController.ProductDetails())
	app.Get("/api/product-media/{objectKey:path}", gatewayController.GetProductMedia())
}
