package auth

import (
	"minieccomerce/utils"

	"github.com/labstack/echo/v4"
)

type AuthController interface {
	AuthRoutes(e *echo.Echo)
	LoginController(c echo.Context) error
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
}

func (authController *AuthControllerImpl) LoginController(c echo.Context) error {
	ctx := c.Request().Context()

	loginReq := &LoginRequest{}
	if err := c.Bind(loginReq); err != nil {
		return c.JSON(400, utils.CommonResponse{Message: err.Error(), Status: "400"})
	}

	access_token, err := authController.authService.LoginService(ctx, loginReq)
	if err != nil {
		return c.JSON(400, utils.CommonResponse{Message: err.Error(), Status: "400"})
	}

	return c.JSON(200, LoginResponse{AccessToken: access_token.AccessToken})
}
