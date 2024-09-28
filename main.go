// main.go
package main

import (
	"bank/config"
	"bank/routes"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	config.SetupDB()
	app := fiber.New()

	routes.SetupUserRoutes(app)

	log.Fatal(app.Listen(":" + port))
}
