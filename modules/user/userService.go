package user

import (
	"context"
	"errors"

	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	RegisterUserService(ctx context.Context, userReq *UserRequest) error
}

type UserServiceImpl struct {
	db             *sqlx.DB
	userRepository UserRepository
}

func NewUserService(userRepository UserRepository, db *sqlx.DB) UserService {
	return &UserServiceImpl{
		userRepository: userRepository,
		db:             db,
	}
}

func (userService *UserServiceImpl) RegisterUserService(ctx context.Context, userReq *UserRequest) error {
	tx, err := userService.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	exists, err := userService.userRepository.ReadUserByUsername(ctx, userReq.Username)
	if err != nil {
		return err
	}
	if exists != nil {
		return errors.New("username already exists")
	}

	if len(userReq.Password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userReq.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("failed to hash password: " + err.Error())
	}

	user := &User{
		Name:     userReq.Name,
		Username: userReq.Username,
		Password: string(hashedPassword),
		Role:     "user",
	}

	if err := userService.userRepository.CreateUser(ctx, user); err != nil {
		return err
	}

	return nil
}
