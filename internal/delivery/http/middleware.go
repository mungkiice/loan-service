package http

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mungkiice/-loan-service/internal/usecase"
)

const BearerPrefix = "Bearer "

func AuthMiddleware(authUseCase *usecase.AuthUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing auth header"})
			c.Abort()
			return
		}

		if !strings.HasPrefix(auth, BearerPrefix) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid auth format"})
			c.Abort()
			return
		}

		token := strings.TrimPrefix(auth, BearerPrefix)
		claims, err := authUseCase.ValidateToken(c.Request.Context(), token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		c.Set("uid", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("utype", claims.UserType)
		c.Set("role", claims.Role)

		c.Next()
	}
}

func RequireUserType(required string) gin.HandlerFunc {
	return func(c *gin.Context) {
		utype, ok := c.Get("utype")
		if !ok || utype != required {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func RequireRole(required string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, ok := c.Get("role")
		if !ok || role != required {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func userID(c *gin.Context) (string, bool) {
	id, ok := c.Get("uid")
	if !ok {
		return "", false
	}
	if s, ok := id.(interface{ String() string }); ok {
		return s.String(), true
	}
	s, ok := id.(string)
	return s, ok
}
