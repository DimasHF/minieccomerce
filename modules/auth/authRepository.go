package auth

import (
	"context"
	"minieccomerce/modules/user"

	"github.com/jmoiron/sqlx"
)

type AuthRepository interface {
	FindUserByUsername(ctx context.Context, username string) (*user.User, error)
}

type AuthRepositoryImpl struct {
	db *sqlx.DB
}

func NewAuthRepository(db *sqlx.DB) *AuthRepositoryImpl {
	return &AuthRepositoryImpl{
		db: db,
	}
}

func (authRepo *AuthRepositoryImpl) FindUserByUsername(ctx context.Context, username string) (*user.User, error) {
	var user user.User
	err := authRepo.db.GetContext(ctx, &user, "SELECT * FROM users WHERE username = $1", username)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
