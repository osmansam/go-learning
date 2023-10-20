package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/osmansam/react/controllers"
)

func UserRoutes(baseUrl string, app *fiber.App) {
	userGroup := app.Group(baseUrl)
	userGroup.Post("/register", controllers.Register)
}