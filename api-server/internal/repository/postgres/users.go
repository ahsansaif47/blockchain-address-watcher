package postgres

import (
	"context"

	sqlc "github.com/ahsansaif47/blockchain-address-watcher/api-server/db/generated"
	"github.com/jackc/pgx/v5/pgtype"
)

type IUserInterface interface {
}

type UserRepo struct {
	ctx context.Context
	db  *sqlc.Queries
}

func NewRepository(db sqlc.DBTX) IUserInterface {
	return &UserRepo{
		db:  sqlc.New(db),
		ctx: context.Background(),
	}
}

func (r *UserRepo) CreateNewUser(user sqlc.CreateUserParams) (pgtype.UUID, error) {
	id, err := r.db.CreateUser(r.ctx, user)
	if err != nil {
		return pgtype.UUID{}, err
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

func (r *UserRepo) SoftDeleteUser(id pgtype.UUID) error {
	return r.db.SoftDeleteUser(r.ctx, id)
}

func (r *UserRepo) HardDeleteUser(id pgtype.UUID) error {
	return r.db.HardDeleteUser(r.ctx, id)
}
