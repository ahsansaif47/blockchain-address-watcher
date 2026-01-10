package main

import (
	"log"

	"github.com/ahsansaif47/blockchain-address-watcher/api-server/config"
	"github.com/ahsansaif47/blockchain-address-watcher/api-server/internal/api"
	"github.com/ahsansaif47/blockchain-address-watcher/api-server/internal/repository/postgres"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	// Load configuration
	cfg := config.GetConfig()

	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName: "Blockchain Address Watcher API",
		// DisableStartupMessage: false,
		// ErrorHandler:          customErrorHandler,
	})

	// Middleware
	app.Use(recover.New()) // Recover from panics
	app.Use(logger.New())  // Request logging
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization",
		AllowCredentials: true,
	}))

	// Initialize database
	// TODO: This needs to be fixed - currently creating both connection and pool
	// The repository should use the pool, but NewUserRepository receives nil
	postgres.GetDatabaseInstance()
	log.Printf("Database connected successfully")

	// Setup routes
	api.SetupRoutes(app)

	// Start server
	port := cfg.Port
	if port == "" {
		port = "3000"
	}

	log.Printf("Server starting on port %s", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// // customErrorHandler handles errors in a standardized way
// func customErrorHandler(c *fiber.Ctx, err error) error {
// 	code := fiber.StatusInternalServerError

// 	if e, ok := err.(*fiber.Error); ok {
// 		code = e.Code
// 	}

// 	return c.Status(code).JSON(fiber.Map{
// 		"error": err.Error(),
// 	})
// }
