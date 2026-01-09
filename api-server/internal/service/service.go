package service

import (
	sqlc "github.com/ahsansaif47/blockchain-address-watcher/api-server/db/generated"
	"github.com/ahsansaif47/blockchain-address-watcher/api-server/internal/dto"
	"github.com/ahsansaif47/blockchain-address-watcher/api-server/internal/repository/postgres"
	"github.com/ahsansaif47/blockchain-address-watcher/api-server/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type IUserService interface {
	CreateNewUser(user dto.User) (int, string, error)
	GetUser(email string) (int, *dto.User, error)
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

func (s *UserService) CreateNewUser(user dto.User) (int, string, error) {

	uuid := uuid.New()
	pgUUID := pgtype.UUID{}
	if err := pgUUID.Scan(uuid); err != nil {
		return fiber.StatusBadRequest, "", err
	}

	// TODO: Change the user.Password
	usr := sqlc.CreateUserParams{
		ID:            pgUUID,
		Email:         user.Email,
		PasswordHash:  user.PasswordHash,
		PhoneNumber:   utils.ToPgText(&user.PhoneNo),
		WalletAddress: utils.ToPgText(&user.WalletAddress),
		Subscribed:    false,
	}

	// utils.

	id, err := s.repo.CreateNewUser(usr)
	if err != nil {
		return fiber.StatusInternalServerError, "", err
	}

	userID, err := utils.PgUUIDToUUID(id)
	return fiber.StatusCreated, userID, nil
}

func (s *UserService) GetUser(email string) (int, *dto.User, error) {

	user, err := s.repo.GetUser(email)

	retUser := dto.User{
		Email:         user.Email,
		PasswordHash:  user.PasswordHash,
		PhoneNo:       user.PhoneNumber.String,
		WalletAddress: user.WalletAddress.String,
		Subscribed:    user.Subscribed,
		CreatedAt:     user.CreatedAt.Time,
		UpdatedAt:     user.UpdatedAt.Time,
		DeletedAt:     &user.DeletedAt.Time,
	}
	if err != nil {
		return fiber.StatusInternalServerError, nil, err
	}

	return fiber.StatusOK, &retUser, nil
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
