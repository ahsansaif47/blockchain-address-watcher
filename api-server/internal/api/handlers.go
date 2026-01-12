package api

import (
	"github.com/ahsansaif47/blockchain-address-watcher/api-server/internal/dto"
	"github.com/ahsansaif47/blockchain-address-watcher/api-server/internal/service"
	"github.com/ahsansaif47/blockchain-address-watcher/api-server/utils/validators"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	service   service.IUserService
	validator *validator.Validate
}

func NewUserHandler(userService service.IUserService, validator *validator.Validate) *UserHandler {
	return &UserHandler{
		service:   userService,
		validator: validator,
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

	if err := h.validator.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error:   "Validation failed",
			Details: "Please check the fields and try again",
			Fields:  validators.GetValidationErrors(err),
		})
	}

	status, userID, err := h.service.RegisterUser(req)
	if err != nil {
		c.Status(status).JSON(dto.ErrorResponse{
			Error:   "Failed to register",
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
	status, res, err := h.service.Login(req)
	if err != nil {
		return c.Status(status).JSON(dto.ErrorResponse{
			Error:   "Failed to authenticate",
			Details: err.Error(),
		})
	}

	if res == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(dto.ErrorResponse{
			Error: "Invalid credentials",
		})
	}

	return c.Status(status).JSON(res)
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

	return c.Status(status).JSON(dto.DeleteUserResponse{
		Message: "User deleted successfully",
	})
}
