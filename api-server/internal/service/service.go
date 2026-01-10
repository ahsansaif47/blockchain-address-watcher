package service

import (
	"fmt"

	sqlc "github.com/ahsansaif47/blockchain-address-watcher/api-server/db/generated"
	"github.com/ahsansaif47/blockchain-address-watcher/api-server/internal/dto"
	"github.com/ahsansaif47/blockchain-address-watcher/api-server/internal/repository/postgres"
	"github.com/ahsansaif47/blockchain-address-watcher/api-server/utils"
	"github.com/ahsansaif47/blockchain-address-watcher/api-server/utils/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type IUserService interface {
	RegisterUser(user dto.RegisterUserRequest) (int, string, error)
	Login(req dto.LoginRequest) (int, *dto.LoginResponse, error)
	SoftDeleteUser(id string) (int, error)
	HardDeleteUser(id string) (int, error)
}

type UserService struct {
	repo postgres.IUserInterface
}

func NewService(repo postgres.IUserInterface) IUserService {
	return &UserService{
		repo: repo,
	}
}

func (s *UserService) RegisterUser(user dto.RegisterUserRequest) (int, string, error) {

	uuid := uuid.New()
	pgUUID := pgtype.UUID{}
	if err := pgUUID.Scan(uuid); err != nil {
		return fiber.StatusBadRequest, "", err
	}

	passHash, err := utils.HashPassword(user.Password)
	if err != nil {
		return fiber.StatusInternalServerError, "", err
	}

	usr := sqlc.CreateUserParams{
		ID:            pgUUID,
		Email:         user.Email,
		PasswordHash:  passHash,
		PhoneNumber:   utils.ToPgText(&user.PhoneNo),
		WalletAddress: utils.ToPgText(&user.WalletAddress),
		Subscribed:    false,
	}

	id, err := s.repo.CreateNewUser(usr)
	if err != nil {
		return fiber.StatusInternalServerError, "", err
	}

	userID, err := utils.PgUUIDToUUID(id)
	return fiber.StatusCreated, userID, nil
}

func (s *UserService) Login(req dto.LoginRequest) (int, *dto.LoginResponse, error) {

	user, err := s.repo.GetUser(req.Email)
	if err != nil {
		return fiber.StatusInternalServerError, nil, err
	}

	// Compare the hash here from the utils function..

	status := utils.ComparePasswordHash(req.Password, user.PasswordHash)
	fmt.Println("Status is: ", status)

	// Generate the token if status is true
	token, err := jwt.GenerateJWT(req.Email)
	if err != nil {
		return fiber.StatusInternalServerError, nil, err
	}

	res := dto.LoginResponse{ID: user.ID.String(), Token: token}

	return fiber.StatusOK, &res, nil
}

func (s *UserService) SoftDeleteUser(id string) (int, error) {
	pgUUID := pgtype.UUID{}
	if err := pgUUID.Scan(id); err != nil {
		return fiber.StatusBadRequest, err
	}

	if err := s.repo.SoftDeleteUser(pgUUID); err != nil {
		return fiber.StatusInternalServerError, err
	}

	return fiber.StatusOK, nil
}

func (s *UserService) HardDeleteUser(id string) (int, error) {
	pgUUID := pgtype.UUID{}
	if err := pgUUID.Scan(id); err != nil {
		return fiber.StatusBadRequest, err
	}

	if err := s.repo.HardDeleteUser(pgUUID); err != nil {
		return fiber.StatusInternalServerError, err
	}

	return fiber.StatusOK, nil
}
