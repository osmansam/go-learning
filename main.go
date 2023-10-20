package main

import (
	"log"
	"os"

	"github.com/osmansam/react/configs"
	"github.com/osmansam/react/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)




const portNumber = ":3002"
func main() {

	app := fiber.New()
	// Create a new directory if not exists to store images
		if _, err := os.Stat("./temp"); os.IsNotExist(err) {
    os.Mkdir("./temp", 0755)
}
	//cors
app.Use(cors.New())
	//run database
	configs.ConnectDB()

	//routes
	routes.HelalRoutes("api/v1/helal", app)
	routes.OsmanRoutes("api/v1/osman", app)
	routes.ResimRoutes("api/v1/resim", app)
	log.Println("Server is running on port: ", portNumber)
	app.Listen(portNumber)
}