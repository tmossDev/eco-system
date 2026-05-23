package routes

import (
	"github.com/kataras/iris/v12"
	userClient "tmossDev.github.com/eco-system/product-management/backend/app/product-gateway/client"
	"tmossDev.github.com/eco-system/product-management/backend/app/product-gateway/controller"
	"tmossDev.github.com/eco-system/product-management/backend/domain/product/service"
)

func Setup(app *iris.Application, productService service.ProductService, authClient userClient.AuthClient) {

	gatewayController := controller.NewGatewayControllerImp(productService, authClient)
	// auth routes
	app.Post("/api/auth/login", gatewayController.Login())
	app.Post("/api/login", gatewayController.Login())
	app.Put("/api/login", gatewayController.Login())
	app.Post("/api/auth/logout", gatewayController.Logout())
	app.Put("/api/logout", gatewayController.Logout())

	// admin app routes
	app.Get("/api/dashboard/summary", gatewayController.DashboardSummary())
	app.Get("/api/settings", gatewayController.GetSettings())
	app.Put("/api/settings", gatewayController.UpdateSettings())
	app.Get("/api/products", gatewayController.ListProducts())
	app.Post("/api/products", gatewayController.CreateProduct())
	app.Get("/api/products/{id:uint64}", gatewayController.ProductDetails())
	app.Put("/api/products/{id:uint64}", gatewayController.UpdateProduct())
	app.Delete("/api/products/{id:uint64}", gatewayController.DeleteProduct())
	app.Post("/api/products/{id:uint64}/photos", gatewayController.UploadProductPhoto())
	app.Get("/api/product-media/{objectKey:path}", gatewayController.GetProductMedia())
}
