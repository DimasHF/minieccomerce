package auth

import (
	"context"
	"errors"
	"minieccomerce/utils"
	"time"

	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	LoginService(ctx context.Context, loginRequest *LoginRequest, userAgent string) (*LoginResponse, error)
	RefreshTokenService(ctx context.Context, userID int, userAgent string) (*LoginResponse, error)
}

type AuthServiceImpl struct {
	db             *sqlx.DB
	authRepository AuthRepository
}

func NewAuthService(db *sqlx.DB, authRepository AuthRepository) AuthService {
	return &AuthServiceImpl{
		db:             db,
		authRepository: authRepository,
	}
}

func (authService *AuthServiceImpl) LoginService(ctx context.Context, loginRequest *LoginRequest, userAgent string) (*LoginResponse, error) {
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

	user, err := authService.authRepository.ReadUserByUsername(ctx, loginRequest.Username)
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

	saveToken := &utils.TokenPair{
		UserID:       user.ID,
		UserAgent:    userAgent,
		RefreshToken: tokens.RefreshToken,
		ExpiredAt:    tokens.ExpiredAt,
		LastLoginAt:  tokens.LastLoginAt,
		CreatedAt:    tokens.CreatedAt,
	}
	if err := authService.authRepository.CreateSaveToken(ctx, tx, saveToken); err != nil {
		return nil, errors.New("login failed")
	}

	return &LoginResponse{
		AccessToken: tokens.AccessToken,
	}, nil
}

func (authService *AuthServiceImpl) RefreshTokenService(ctx context.Context, userID int, userAgent string) (*LoginResponse, error) {
	tx, err := authService.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, errors.New("failed to begin transaction " + err.Error())
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	token, err := authService.authRepository.ReadRefreshToken(ctx, userID, userAgent)
	if err != nil || token == nil {
		return nil, errors.New("failed to find refresh token: " + err.Error())
	}

	if token.ExpiredAt.Before(time.Now().UTC()) {
		return nil, errors.New("refresh token expired")
	}

	user, err := authService.authRepository.ReadUserByID(ctx, userID)
	if err != nil {
		return nil, errors.New("failed to find user: " + err.Error())
	}
	tokens, err := utils.GenerateTokenPair(user.ID, user.Name, user.Role)
	if err != nil {
		return nil, errors.New("failed to generate token pair: " + err.Error())
	}

	return &LoginResponse{AccessToken: tokens.AccessToken}, nil
}
