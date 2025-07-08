package auth

import (
	"minieccomerce/utils"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

type AuthController interface {
	AuthRoutes(e *echo.Echo)
	LoginController(c echo.Context) error
	RefreshTokenController(c echo.Context) error
}

type AuthControllerImpl struct {
	authService AuthService
}

func NewAuthController(authService AuthService) AuthController {
	return &AuthControllerImpl{
		authService: authService,
	}
}

func (authController *AuthControllerImpl) AuthRoutes(e *echo.Echo) {
	e.POST("/login", authController.LoginController)
	e.GET("/refresh", authController.RefreshTokenController)
}

func (authController *AuthControllerImpl) LoginController(c echo.Context) error {
	ctx := c.Request().Context()

	loginReq := &LoginRequest{}
	if err := c.Bind(loginReq); err != nil {
		return c.JSON(400, utils.CommonResponse{Message: err.Error(), Status: "400"})
	}

	userAgent := c.Request().UserAgent()
	if userAgent == "" {
		userAgent = "unknown"
	}

	access_token, err := authController.authService.LoginService(ctx, loginReq, userAgent)
	if err != nil {
		return c.JSON(400, utils.CommonResponse{Message: err.Error(), Status: "400"})
	}

	return c.JSON(200, LoginResponse{AccessToken: access_token.AccessToken})
}

func (authController *AuthControllerImpl) RefreshTokenController(c echo.Context) error {
	ctx := c.Request().Context()

	userAgent := c.Request().UserAgent()

	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		return c.JSON(http.StatusUnauthorized, echo.Map{"message": "Authorization header is missing"})
	}
	accessToken := strings.TrimPrefix(authHeader, "Bearer ")

	claims, err := utils.ParseAccessToken(accessToken)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"message": "Invalid access token: " + err.Error()})
	}

	userID := claims.ID

	tokens, err := authController.authService.RefreshTokenService(ctx, userID, userAgent)
	if err != nil {
		return c.JSON(500, utils.CommonResponse{Message: "failed to refresh token: " + err.Error(), Status: "500"})
	}

	return c.JSON(200, tokens)
}
