package dic

import (
	"golang-api-auth-template/http/controllers"
	"golang-api-auth-template/logger"
)

var authController *controllers.AuthController

func AuthController() *controllers.AuthController {
	if authController == nil {
		authController = controllers.NewAuthController(
			UserService(),
			JwtService(),
			RoleService(),
		)
		logger.Logger().Debug("AuthController initialized")
	}

	return authController
}
