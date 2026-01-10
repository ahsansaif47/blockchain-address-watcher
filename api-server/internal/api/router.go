package api

import (
	"github.com/ahsansaif47/blockchain-address-watcher/api-server/internal/repository/postgres"
	"github.com/ahsansaif47/blockchain-address-watcher/api-server/internal/service"
	"github.com/ahsansaif47/blockchain-address-watcher/api-server/utils/validators"
	"github.com/gofiber/fiber/v2"
)

// SetupRoutes configures all API routes
func SetupRoutes(app *fiber.App) {
	// Initialize repository
	userRepo := postgres.NewUserRepository(nil) // TODO: Pass actual database connection

	// Initialize service
	userService := service.NewService(userRepo)

	// Initialize validator with custom validators
	validator := validators.NewValidator()

	// Initialize handler
	userHandler := NewUserHandler(userService, validator)

	// API v1 routes
	api := app.Group("/api/v1")

	// User routes
	users := api.Group("/users")
	{
		// Public routes
		users.Post("/register", userHandler.Register)
		users.Post("/login", userHandler.Login)

		// Protected routes (TODO: Add authentication middleware)
		users.Get("/", userHandler.Login)
		users.Delete("/delete", userHandler.DeleteUser)
	}

	// Health check endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"service": "blockchain-address-watcher-api",
		})
	})

	// Root endpoint
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Blockchain Address Watcher API",
			"version": "1.0.0",
		})
	})
}
