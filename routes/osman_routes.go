package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/osmansam/react/controllers"
)

func OsmanRoutes(baseUrl string, app *fiber.App) {
	osmanGroup := app.Group(baseUrl)
	osmanGroup.Post("/", controllers.CreateOsman)
	osmanGroup.Get("/", controllers.GetAllOsman)

}