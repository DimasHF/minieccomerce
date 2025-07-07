package auth

import (
	"context"
	"errors"
	"minieccomerce/utils"

	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	LoginService(ctx context.Context, loginRequest *LoginRequest) (*LoginResponse, error)
}

type AuthServiceImpl struct {
	db             *sqlx.DB
	authRepository AuthRepository
}

func NewAuthService(db *sqlx.DB) *AuthServiceImpl {
	return &AuthServiceImpl{
		db:             db,
		authRepository: NewAuthRepository(db),
	}
}

func (authService *AuthServiceImpl) LoginService(ctx context.Context, loginRequest *LoginRequest) (*LoginResponse, error) {
	tx, err := authService.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	user, err := authService.authRepository.FindUserByUsername(ctx, loginRequest.Username)
	if err != nil || user == nil {
		return nil, errors.New("login failed")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password)); err != nil {
		return nil, errors.New("login failed")
	}

	tokens, err := utils.GenerateTokenPair(user.ID, user.Name, user.Role)
	if err != nil || tokens == nil {
		return nil, errors.New("login failed")
	}

	return &LoginResponse{
		AccessToken: tokens.AccessToken,
	}, nil
}
