package router

import (
	"app/dic"
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	router := gin.Default()

	router.Use(dic.AuthMiddleware())

	apiGroup := router.Group("/api")
	{
		authGroup := apiGroup.Group("/auth")
		{
			authGroup.POST("/register", dic.AuthController().Register)
			authGroup.POST("/me", dic.RequiredAuthMiddleware(), dic.AuthController().Me)
			authGroup.POST("/login", dic.AuthController().Login)
			authGroup.POST("/refresh", dic.AuthController().Refresh)
		}
	}

	return router
}
