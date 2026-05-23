package routes

import (
	"github.com/kataras/iris/v12"
	"tmossDev.github.com/eco-system/product-management/backend/app/product-service/controller"
	"tmossDev.github.com/eco-system/product-management/backend/domain/product/service"
)

func Setup(app *iris.Application, productService service.ProductService) {

	productController := controller.NewServiceControllerImp(productService)
	app.Get("/api/products", productController.ListProducts())
	app.Post("/api/products", productController.CreateProduct())
	app.Get("/api/products/{id:uint64}", productController.ProductDetails())
	app.Put("/api/products/{id:uint64}", productController.UpdateProduct())
	app.Delete("/api/products/{id:uint64}", productController.DeleteProduct())
	app.Get("/api/discounts", productController.ListDiscounts())
	app.Post("/api/discounts", productController.CreateDiscount())
	app.Get("/api/discounts/{id:uint64}", productController.DiscountDetails())
	app.Put("/api/discounts/{id:uint64}", productController.UpdateDiscount())
	app.Delete("/api/discounts/{id:uint64}", productController.DeleteDiscount())
}
