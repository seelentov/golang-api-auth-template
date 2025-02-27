package dic

import (
	"github.com/gin-gonic/gin"
	"golang-api-auth-template/http/middlewares"
	"golang-api-auth-template/logger"
	"os"
)

var authMiddleware gin.HandlerFunc

func AuthMiddleware() gin.HandlerFunc {
	if authMiddleware == nil {
		authMiddleware = middlewares.AuthMiddleware(
			logger.Logger(),
			os.Getenv("JWT_SECRET"),
		)
		logger.Logger().Debug("AuthMiddleware initialized")
	}

	return authMiddleware
}

var requiredAuthMiddleware gin.HandlerFunc

func RequiredAuthMiddleware() gin.HandlerFunc {
	if requiredAuthMiddleware == nil {
		requiredAuthMiddleware = middlewares.RequiredAuthMiddleware()
		logger.Logger().Debug("RequiredAuthMiddleware initialized")
	}

	return requiredAuthMiddleware
}

func RequiredRoleMiddleware(roleNames []string) gin.HandlerFunc {
	return middlewares.RequiredRolesMiddleware(
		roleNames,
		logger.Logger(),
		RoleService(),
	)
}
