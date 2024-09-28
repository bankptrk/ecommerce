package routes

import (
	"bank/controller"

	"github.com/gofiber/fiber/v2"
)

func SetupUserRoutes(app *fiber.App) {
	app.Post("/users/register", controller.RegisterUser)

	app.Post("/users/login", controller.LoginUser)

	app.Post("/admin/addproduct", controller.CreateProduct)
	app.Get("/user", controller.GetUsers)
	app.Get("/users/products", controller.SearchProducts)
	app.Get("/users/search", controller.SearchProductByQuery)

	// app.Use(middleware.AuthRequired)

	app.Post("/addaddress", controller.AddAddress)
	app.Put("/editaddress", controller.EditHomeAddress)
	app.Put("/editaddress2", controller.EditHomeAddress2)
	app.Delete("/deleteaddress", controller.DeleteAddress)
}
