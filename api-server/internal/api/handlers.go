package api

import (
	"github.com/ahsansaif47/blockchain-address-watcher/api-server/internal/dto"
	"github.com/ahsansaif47/blockchain-address-watcher/api-server/internal/service"
	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	service service.IUserService
}

func NewUserHandler(userService service.IUserService) *UserHandler {
	return &UserHandler{
		service: userService,
	}
}

// Register handles user registration
// @Summary Register a new user
// @Description Create a new user account
// @Tags users
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "User registration details"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/users/register [post]
func (h *UserHandler) Register(c *fiber.Ctx) error {
	var req dto.RegisterUserRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error:   "Invalid request body",
			Details: err.Error(),
		})
	}

	status, userID, err := h.service.RegisterUser(req)
	if err != nil {
		return c.Status(status).JSON(dto.ErrorResponse{
			Error:   "Failed to create user",
			Details: err.Error(),
		})
	}

	return c.Status(status).JSON(fiber.Map{"id": userID})
}

// Login handles user login
// @Summary Login user
// @Description Authenticate user with email and password
// @Tags users
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Login credentials"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/users/login [post]
func (h *UserHandler) Login(c *fiber.Ctx) error {
	var req dto.LoginRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error:   "Invalid request body",
			Details: err.Error(),
		})
	}

	// Service layer handles authentication logic
	// TODO: Implement password verification and JWT token generation in service layer
	status, user, err := h.service.Login(req)
	if err != nil {
		return c.Status(status).JSON(dto.ErrorResponse{
			Error:   "Failed to authenticate",
			Details: err.Error(),
		})
	}

	if user == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(dto.ErrorResponse{
			Error: "Invalid credentials",
		})
	}

	return c.Status(status).JSON(dto.LoginResponse{})
}

// GetUser retrieves a user by email
// @Summary Get user by email
// @Description Get user details
// @Tags users
// @Produce json
// @Param email query string true "User email"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/users [get]
func (h *UserHandler) GetUser(c *fiber.Ctx) error {
	email := c.Query("email")
	if email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error: "Email query parameter is required",
		})
	}

	status, user, err := h.service.GetUser(email)
	if err != nil {
		return c.Status(status).JSON(dto.ErrorResponse{
			Error:   "Failed to get user",
			Details: err.Error(),
		})
	}

	return c.Status(status).JSON(dto.SuccessResponse{
		Message: "User retrieved successfully",
		Data: map[string]interface{}{
			"user": user,
		},
	})
}

// DeleteUser handles user deletion (soft or hard)
// @Summary Delete user
// @Description Delete a user account (soft or hard delete)
// @Tags users
// @Accept json
// @Produce json
// @Param request body dto.DeleteUserRequest true "Deletion details"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/users/delete [delete]
func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	var req dto.DeleteUserRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error:   "Invalid request body",
			Details: err.Error(),
		})
	}

	// TODO: Move validation logic to service layer
	// Service should handle validation of user ID and delete type

	var status int
	var err error

	if req.Type == "soft" {
		status, err = h.service.SoftDeleteUser(req.UserID)
	} else {
		status, err = h.service.HardDeleteUser(req.UserID)
	}

	if err != nil {
		return c.Status(status).JSON(dto.ErrorResponse{
			Error:   "Failed to delete user",
			Details: err.Error(),
		})
	}

	return c.Status(status).JSON(dto.SuccessResponse{
		Message: "User deleted successfully",
	})
}
