package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"

	"fiber-first-app/db"
	"fiber-first-app/routes"

	_ "fiber-first-app/docs" // Needed if you are using Swagger
	"github.com/swaggo/fiber-swagger" // fiber-swagger middleware
)

// @title Fiber CRUD API
// @version 1.0
// @description A simple CRUD API built with Fiber, PostgreSQL, and AWS Cognito.
// @host localhost:3000
// @BasePath /api
func main() {
	// Initialize logger
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	logrus.Info("🔧 Starting application...")

	// Initialize database connection
	db.InitDB()

	// Create a new Fiber app
	app := fiber.New()

	// Optional: Swagger docs endpoint
	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	// Register routes
	routes.SetupRoutes(app)

	// Start the server
	port := ":3000"
	logrus.Infof("🚀 Server is running on http://localhost%s", port)
	if err := app.Listen(port); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
}
