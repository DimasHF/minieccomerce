package auth

import (
	"context"
	"errors"
	"minieccomerce/modules/user"
	"minieccomerce/utils"

	"github.com/jmoiron/sqlx"
)

type AuthRepository interface {
	CreateSaveToken(ctx context.Context, tx *sqlx.Tx, token *utils.TokenPair) error
	ReadUserByID(ctx context.Context, userID int) (*user.User, error)
	ReadUserByUsername(ctx context.Context, username string) (*user.User, error)
	ReadRefreshToken(ctx context.Context, userID int, userAgent string) (*utils.TokenPair, error)
}

type AuthRepositoryImpl struct {
	db *sqlx.DB
}

func NewAuthRepository(db *sqlx.DB) AuthRepository {
	return &AuthRepositoryImpl{
		db: db,
	}
}

func (authRepo *AuthRepositoryImpl) CreateSaveToken(ctx context.Context, tx *sqlx.Tx, token *utils.TokenPair) error {
	query := `
		INSERT INTO user_tokens (
			user_id, user_agent, refresh_token, expired_at, last_login_at, created_at
		) VALUES (
			:user_id, :user_agent, :refresh_token, :expired_at, :last_login_at, :created_at
		) ON CONFLICT (user_id, user_agent) DO UPDATE SET
			refresh_token = EXCLUDED.refresh_token,
			expired_at = EXCLUDED.expired_at,
			last_login_at = EXCLUDED.last_login_at
	`

	params := map[string]interface{}{
		"user_id":       token.UserID,
		"user_agent":    token.UserAgent,
		"refresh_token": token.RefreshToken,
		"expired_at":    token.ExpiredAt,
		"last_login_at": token.LastLoginAt,
		"created_at":    token.CreatedAt,
	}

	if _, err := authRepo.db.NamedExecContext(ctx, query, params); err != nil {
		return errors.New("failed to save refresh token: " + err.Error())
	}
	return nil
}

func (authRepo *AuthRepositoryImpl) ReadUserByID(ctx context.Context, userID int) (*user.User, error) {
	var user user.User
	err := authRepo.db.GetContext(ctx, &user, "SELECT * FROM users WHERE id = $1", userID)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (authRepo *AuthRepositoryImpl) ReadUserByUsername(ctx context.Context, username string) (*user.User, error) {
	var user user.User
	err := authRepo.db.GetContext(ctx, &user, "SELECT * FROM users WHERE username = $1", username)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (authRepo *AuthRepositoryImpl) ReadRefreshToken(ctx context.Context, userID int, userAgent string) (*utils.TokenPair, error) {
	var token utils.TokenPair
	err := authRepo.db.GetContext(ctx, &token, "SELECT * FROM user_tokens WHERE user_id = $1 AND user_agent = $2", userID, userAgent)
	if err != nil {
		return nil, err
	}
	return &token, nil
}
