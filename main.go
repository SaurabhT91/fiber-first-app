package main

import (
	"io"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"

	"fiber-first-app/db"
	"fiber-first-app/routes"

	_ "fiber-first-app/docs" // Swagger docs

	fiberSwagger "github.com/swaggo/fiber-swagger" // Swagger middleware
)

// @title Fiber CRUD API
// @version 1.0
// @description A simple CRUD API built with Fiber, PostgreSQL, and AWS Cognito.
// @host localhost:3000
// @BasePath /api
func main() {
	// 🔧 Setup logging to both a file and stdout
	logFile, err := os.OpenFile("fiber-app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("❌ Failed to open log file: %v", err)
	}

	logrus.SetOutput(io.MultiWriter(os.Stdout, logFile))
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	logrus.SetLevel(logrus.InfoLevel)

	logrus.Info("🔧 Starting application...")

	// Initialize DB
	db.InitDB()

	// Create Fiber app
	app := fiber.New()

	// Swagger endpoint
	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	// Register routes
	routes.SetupRoutes(app)

	// Start server
	port := ":3000"
	logrus.Infof("🚀 Server is running at http://localhost%s", port)
	if err := app.Listen(port); err != nil {
		logrus.Fatalf("❌ Failed to start server: %v", err)
	}
}
