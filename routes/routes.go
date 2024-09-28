package routes

import (
	"bank/controller"
	"bank/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupUserRoutes(app *fiber.App) {
	app.Post("/users/register", controller.RegisterUser)

	app.Post("/users/login", controller.LoginUser)

	app.Post("/admin/product", controller.CreateProduct)

	app.Use(middleware.AuthRequired)

	app.Get("/user", controller.GetUsers)
	app.Get("/users/products", controller.SearchProducts)
	app.Get("/users/search", controller.SearchProductByQuery)
}
