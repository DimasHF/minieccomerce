package user

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *User) error
	ReadUserByUsername(ctx context.Context, username string) (*User, error)
}

type UserRepositoryImpl struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) UserRepository {
	return &UserRepositoryImpl{
		db: db,
	}
}

func (userRepo *UserRepositoryImpl) CreateUser(ctx context.Context, user *User) error {
	query := "INSERT INTO users (name, username, password, role) VALUES (:name, :username, :password, :role)"

	params := map[string]interface{}{
		"name":     user.Name,
		"username": user.Username,
		"password": user.Password,
		"role":     user.Role,
	}

	_, err := userRepo.db.NamedExecContext(ctx, query, params)
	return err
}

func (userRepo *UserRepositoryImpl) ReadUserByUsername(ctx context.Context, username string) (*User, error) {
	var user User
	err := userRepo.db.GetContext(ctx, &user, "SELECT * FROM users WHERE username = $1", username)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
