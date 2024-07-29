package routes

import (
	"github.com/AhishekOza/E-commerce/controllers"
	"github.com/AhishekOza/E-commerce/middleware"
	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {

	// ---AUTH ROUTES
	app.Post("/api/v1/register", controllers.Register)
	app.Post("/api/v1/login", controllers.Login)

	// ---CATEGORY ROUTES

	// GET CATEGORIES
	app.Get("/api/v1/get-categories", controllers.GetAllCategories)

	// ---PRODUCT ROUTES
	// --GET PRODUCTS
	app.Get("/api/v1/get-all-products", controllers.GetAllProducts)
	// GET PRODUCTS BY CATEGORY
	app.Get("/api/v1/get-products-category", controllers.GetProductByCategory)
	// GET SINGLE PRODUCT
	app.Get("/api/v1/get-single-product/:id", controllers.GetSingleProduct)
	// GET SINGLE PRODUCT
	app.Get("/api/v1/update-product/:id", controllers.UpdateProduct)

	// AUTH ROUTE----
	authGroup := app.Group("/api/v1", middleware.TokenMiddleware)

	// ---CATEGORY AUTH ROUTES
	authGroup.Post("/create-category", controllers.CreateCategory)
	authGroup.Delete("/delete-category/:id", controllers.DeleteCategory)

	// ---PRODUCT AUTH ROUTES
	// --POST PRODUCT
	authGroup.Post("/create-product", controllers.CreateProduct)
	// --EDIT PRODUCT
	// --DELETE PRODUCT

}
