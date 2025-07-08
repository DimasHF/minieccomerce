//go:build wireinject
// +build wireinject

package auth

import (
	"minieccomerce/configs"

	"github.com/google/wire"
)

func InitializeAuthController() AuthController {
	wire.Build(
		configs.InitDB,
		NewAuthRepository,
		NewAuthService,
		NewAuthController,
	)
	return nil
}
