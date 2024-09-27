// main.go
package main

import (
	"bank/config"
	"bank/routes"
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	config.SetupDB()
	app := fiber.New()

	routes.SetupUserRoutes(app)

	log.Fatal(app.Listen(":8000"))
}
