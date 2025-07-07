package utils

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type TokenPair struct {
	UserID       int       `json:"user_id" db:"user_id"`
	UserAgent    string    `json:"user_agent" db:"user_agent"`
	RefreshToken string    `json:"refresh_token" db:"refresh_token"`
	ExpiredAt    time.Time `json:"expired_at" db:"expired_at"`
	LastLoginAt  time.Time `json:"last_login_at" db:"last_login_at"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`

	AccessToken          string    `json:"-" db:"-"`
	AccessTokenExpiresAt time.Time `json:"-" db:"-"`
}

type JWTCustomClaims struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Role string `json:"role"`
	jwt.RegisteredClaims
}

type RefreshTokenClaims struct {
	TokenID string `json:"jti"`
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Role    string `json:"role"`
	jwt.RegisteredClaims
}

var JWTSecret = "1234567890"

func generateRandomString(length int) (string, error) {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func GenerateTokenPair(id int, name, role string) (*TokenPair, error) {
	now := time.Now()

	accessTokenExp := now.Add(15 * time.Minute).UTC()
	refreshTokenExp := now.Add(7 * 24 * time.Hour).UTC()

	accessClaims := &JWTCustomClaims{
		ID:   id,
		Name: name,
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "minieccomerce",
			Subject:   "access_token",
			Audience:  []string{"minieccomerce"},
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(accessTokenExp),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString([]byte(JWTSecret))
	if err != nil {
		return nil, errors.New("login failed")
	}

	tokenID, err := generateRandomString(8)
	if err != nil {
		return nil, errors.New("login failed")
	}

	refreshClaims := &RefreshTokenClaims{
		TokenID: tokenID,
		ID:      id,
		Name:    name,
		Role:    role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(refreshTokenExp),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ID:        tokenID,
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(JWTSecret))
	if err != nil {
		return nil, errors.New("login failed")
	}

	return &TokenPair{
		UserID:               id,
		AccessToken:          accessTokenString,
		AccessTokenExpiresAt: accessTokenExp,
		RefreshToken:         refreshTokenString,
		ExpiredAt:            refreshTokenExp,
	}, nil
}

func ParseAccessToken(tokenStr string) (*JWTCustomClaims, error) {
	token, _, err := new(jwt.Parser).ParseUnverified(tokenStr, &JWTCustomClaims{})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTCustomClaims); ok {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

func ParseToken(c echo.Context) (*JWTCustomClaims, error) {
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		return nil, errors.New("authorization header is missing")
	}

	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenStr == "" {
		return nil, errors.New("token is missing")
	}

	token, err := jwt.ParseWithClaims(tokenStr, &JWTCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(JWTSecret), nil
	})
	if err != nil {
		return nil, errors.New("failed to parse token: " + err.Error())
	}

	if claims, ok := token.Claims.(*JWTCustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token claims")
}
