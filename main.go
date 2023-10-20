package main

import (
	"github.com/osmansam/react/configs"
	"github.com/osmansam/react/routes"

	"github.com/gofiber/fiber/v2"
)

const portNumber = ":3002"
func main() {
	app := fiber.New()

	//run database
	configs.ConnectDB()

	//routes
	routes.HelalRoutes("api/v1/helal", app)
	routes.OsmanRoutes("api/v1/osman", app)

	app.Listen(portNumber)
}