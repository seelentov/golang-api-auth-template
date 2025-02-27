package controllers

import (
	"app/http/dtos"
	"app/http/validator"
	"app/models"
	"app/services"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AuthController struct {
	userService services.UserService
	jwtService  services.JwtService
	roleService services.RoleService
}

func (c *AuthController) Register(ctx *gin.Context) {
	var req dtos.RegisterRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": validator.ParseValidationErrors(err)})
		return
	}

	user := &models.User{
		Name:     req.Name,
		Email:    req.Email,
		Number:   req.Number,
		Password: req.Password,
	}

	err := c.userService.Create(user)

	if err != nil {
		if errors.Is(err, services.ErrAlreadyExist) {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "User already exist"})
			return
		}

		if errors.Is(err, services.ErrDuplicate) {
			if errors.Is(err, services.ErrDuplicateEmail) {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": "User with the same email already exists"})
				return
			}

			if errors.Is(err, services.ErrDuplicateNumber) {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": "User with the same number already exists"})
				return
			}

			if errors.Is(err, services.ErrDuplicateName) {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": "User with the same name already exists"})
				return
			}
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	ctx.Status(http.StatusCreated)
}

func (c *AuthController) Me(ctx *gin.Context) {
	anyId, ok := ctx.Get("user_id")

	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	id, ok := anyId.(float64)

	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	user, err := c.userService.GetById(uint(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	user.Password = ""

	ctx.JSON(http.StatusOK, user)
}

func (c *AuthController) Login(ctx *gin.Context) {
	var req dtos.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": validator.ParseValidationErrors(err)})
		return
	}

	user, err := c.userService.GetByCredential(req.Credential)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Wrong login or password"})
		return
	}

	err = c.userService.VerifyPassword(req.Password, user.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Wrong login or password"})
		return
	}

	token, err := c.jwtService.GenerateToken(user.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// Generate refresh token
	refreshToken, err := c.jwtService.GenerateRefreshToken(user.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate refresh token"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"token":         token,
		"refresh_token": refreshToken,
	})
}

func (c *AuthController) Refresh(ctx *gin.Context) {
	var req dtos.RefreshTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": validator.ParseValidationErrors(err)})
		return
	}

	refreshToken := req.RefreshToken

	claims, err := c.jwtService.ValidateRefreshToken(refreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	userID := uint(claims["user_id"].(float64))

	newToken, err := c.jwtService.GenerateToken(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate new token"})
		return
	}

	//newRefreshToken, err := c.jwtService.GenerateRefreshToken(userID)
	//if err != nil {
	//	ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate new refresh token"})
	//	return
	//}

	ctx.JSON(http.StatusOK, gin.H{
		"token": newToken,
		//"refresh_token": newRefreshToken,
	})
}

func NewAuthController(userService services.UserService, jwtService services.JwtService, roleService services.RoleService) *AuthController {
	return &AuthController{userService, jwtService, roleService}
}
