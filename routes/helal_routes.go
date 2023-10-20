package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/osmansam/react/controllers"
)

func HelalRoutes(baseUrl string, app *fiber.App) {
	helalGroup := app.Group(baseUrl)
	helalGroup.Post("/", controllers.CreateHelal)
	helalGroup.Get("/", controllers.GetAllHelal)
	helalGroup.Get("/search", controllers.SearchHelal)
	helalGroup.Get("/:id", controllers.GetHelal)
	helalGroup.Delete("/:id", controllers.DeleteHelal)
	helalGroup.Patch("/:id", controllers.UpdateHelal)
}
