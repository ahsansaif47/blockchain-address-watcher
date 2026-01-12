package postgres

import (
	"context"

	sqlc "github.com/ahsansaif47/blockchain-address-watcher/api-server/db/generated"
	"github.com/google/uuid"
)

type IUserInterface interface {
	CreateNewUser(user sqlc.CreateUserParams) (uuid.UUID, error)
	GetUser(email string) (*sqlc.User, error)
	SoftDeleteUser(id uuid.UUID) error
	HardDeleteUser(id uuid.UUID) error
}

type UserRepo struct {
	ctx context.Context
	db  *sqlc.Queries
}

func NewUserRepository(db sqlc.DBTX) IUserInterface {
	return &UserRepo{
		db:  sqlc.New(db),
		ctx: context.Background(),
	}
}

func (r *UserRepo) CreateNewUser(user sqlc.CreateUserParams) (uuid.UUID, error) {
	id, err := r.db.CreateUser(r.ctx, user)
	if err != nil {
		return uuid.UUID{}, err
	}

	return id, err
}

func (r *UserRepo) GetUser(email string) (*sqlc.User, error) {
	user, err := r.db.SignInUser(r.ctx, email)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepo) SoftDeleteUser(id uuid.UUID) error {
	return r.db.SoftDeleteUser(r.ctx, id)
}

func (r *UserRepo) HardDeleteUser(id uuid.UUID) error {
	return r.db.HardDeleteUser(r.ctx, id)
}
