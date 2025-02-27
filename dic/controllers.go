package dic

import (
	"app/http/controllers"
	"app/logger"
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
