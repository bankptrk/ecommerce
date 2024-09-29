package routes

import (
	"bank/controller"
	"bank/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupUserRoutes(app *fiber.App) {

	// User Registration and Authentication
	app.Post("/users/register", controller.RegisterUser)
	app.Post("/users/login", controller.LoginUser)

	// Middleware
	secured := app.Group("/", middleware.AuthRequired)

	// Admin Product Management
	secured.Post("/admin/products", controller.CreateProduct)
	secured.Get("/users/products", controller.GetAllProducts)
	secured.Get("/users/search", controller.SearchProductByQuery)

	// Address Management
	secured.Post("/users/addresses", controller.AddAddress)
	secured.Put("/users/addresses/billing", controller.EditBillingAddress)
	secured.Put("/users/addresses/shipping", controller.EditShippingAddress)
	secured.Delete("/users/addresses", controller.DeleteAddress)

	// Cart Management
	secured.Post("/users/cart", controller.AddToCart)
	secured.Delete("/users/cart/items", controller.RemoveItemFromCart)
	secured.Get("/users/cart", controller.GetItemsFromCart)
	secured.Post("/users/cart/purchase", controller.BuyFromCart)
	secured.Post("/users/cart/instant-buy", controller.InstantBuy)
}
