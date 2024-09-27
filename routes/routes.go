package routes

import (
	"bank/controller"
	"bank/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupUserRoutes(app *fiber.App) {
	app.Post("/register", controller.RegisterUser)

	app.Post("/login", controller.LoginUser)

	app.Use("/user", middleware.AuthRequired)

	app.Get("/user", controller.GetUsers)
}
