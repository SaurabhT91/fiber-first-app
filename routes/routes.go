package routes

import (
	"fiber-first-app/handlers"

	"github.com/gofiber/fiber/v2"
	fiberSwagger "github.com/swaggo/fiber-swagger"
)

func SetupRoutes(app *fiber.App) {
	// Swagger
	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	// User Routes
	userAPI := app.Group("/api/users")
	userAPI.Get("/", handlers.GetUsers)
	userAPI.Get("/:id", handlers.GetUserByID)
	userAPI.Post("/", handlers.CreateUser)
	userAPI.Put("/:id", handlers.UpdateUser)
	userAPI.Delete("/:id", handlers.DeleteUser)

	// Event Routes
	eventAPI := app.Group("/api/events")
	eventAPI.Post("/", handlers.CreateEvent)
	eventAPI.Get("/", handlers.GetEvents)
	eventAPI.Get("/:id", handlers.GetEventByID)
	eventAPI.Put("/:id", handlers.UpdateEvent)
	eventAPI.Delete("/:id", handlers.DeleteEvent)
}
