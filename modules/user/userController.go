package user

import (
	"minieccomerce/utils"

	"github.com/labstack/echo/v4"
)

type UserController interface {
	UserRoutes(e *echo.Group)
	RegisterUserController(c echo.Context) error
}

type UserControllerImpl struct {
	userService UserService
}

func NewUserController(userService UserService) UserController {
	return &UserControllerImpl{
		userService: userService,
	}
}

func (userController *UserControllerImpl) UserRoutes(e *echo.Group) {
	e.POST("/register", userController.RegisterUserController)
}

func (userController *UserControllerImpl) RegisterUserController(c echo.Context) error {
	ctx := c.Request().Context()

	userReq := &UserRequest{}
	if err := c.Bind(userReq); err != nil {
		return c.JSON(400, utils.CommonResponse{Message: err.Error(), Status: "400"})
	}

	err := userController.userService.RegisterUserService(ctx, userReq)
	if err != nil {
		return c.JSON(400, utils.CommonResponse{Message: err.Error(), Status: "400"})
	}

	return c.JSON(200, utils.CommonResponse{Message: "success", Status: "200"})
}
