package routes

import (
	"fiber-first-app/handlers"

	"github.com/gofiber/fiber/v2"
	fiberSwagger "github.com/swaggo/fiber-swagger"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api/users")

	api.Get("/", handlers.GetUsers)
	api.Get("/:id", handlers.GetUserByID)
	api.Post("/", handlers.CreateUser)
	api.Put("/:id", handlers.UpdateUser)
	api.Delete("/:id", handlers.DeleteUser)
	app.Get("/swagger/*", fiberSwagger.WrapHandler)

}
