package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mungkiice/-loan-service/internal/usecase"
)

type AuthHandler struct {
	authUseCase *usecase.AuthUseCase
}

func NewAuthHandler(authUseCase *usecase.AuthUseCase) *AuthHandler {
	return &AuthHandler{authUseCase: authUseCase}
}

type SignInRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type SignInResponse struct {
	Token     string                 `json:"token"`
	User      map[string]interface{} `json:"user"`
	ExpiresIn int64                  `json:"expires_in"`
}

func (h *AuthHandler) SignIn(c *gin.Context) {
	var req SignInRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.authUseCase.SignIn(c.Request.Context(), usecase.SignInRequest{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}
