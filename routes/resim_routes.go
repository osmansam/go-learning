package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/osmansam/react/controllers"
)

func ResimRoutes(baseUrl string, app *fiber.App) {
	resimGroup := app.Group(baseUrl)
	resimGroup.Post("/", controllers.CreateResim)
	resimGroup.Get("/", controllers.GetAllResim)
}