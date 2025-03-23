package delivery

import (
	"context"
	"net/http"
	"sample-project/internal/usecase"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	useCase usecase.AuthUseCase
}

func NewAuthHandler(router *gin.Engine, useCase usecase.AuthUseCase) {
	handler := &AuthHandler{useCase: useCase}

	auth := router.Group("/api/v1/auth")

	auth.POST("/login", handler.Login)
	auth.GET("/me", handler.GetUserProfile)
}

// @Summary      Login user
// @Description  Authenticates a user and returns access & refresh tokens
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body entity.LoginRequest true "Login request payload"
// @Success 200 {object} entity.Tokens
// @Failure 400 {object} entity.ErrorResponse
// @Failure 500 {object} entity.ErrorResponse
// @Router       /api/v1/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req struct {
		Name     string `json:"name"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	accessToken, refreshToken, err := h.useCase.Login(context.Background(), req.Name, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

// @Summary      Get authenticated user data
// @Description  Retrieves user data using the authorization token
// @Tags         auth
// @Security 	 BearerAuth
// @Accept       json
// @Produce      json
// @Success 200 {object} entity.User
// @Failure 401 {object} entity.ErrorResponse
// @Failure 500 {object} entity.ErrorResponse
// @Router       /api/v1/auth/me [get]
func (h *AuthHandler) GetUserProfile(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
		return
	}

	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	user, err := h.useCase.GetUserProfile(context.Background(), token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}
