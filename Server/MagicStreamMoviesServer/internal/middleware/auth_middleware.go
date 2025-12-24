package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirolad/MagicStreamMovies/Server/MagicStreamMovies/Server/internal/config"
	"github.com/sirolad/MagicStreamMovies/Server/MagicStreamMovies/Server/pkg/utils"
)

const (
	ContextKeyUserID    = "user_id"
	ContextKeyEmail     = "email"
	ContextKeyFirstName = "first_name"
	ContextKeyLastName  = "last_name"
	ContextKeyRole      = "role"
)

func NewAuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := GetAccessToken(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: " + err.Error()})
			c.Abort()
			return
		}

		if token == "" || strings.TrimSpace(token) == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: Token is missing"})
			c.Abort()
			return
		}

		claims, err := utils.ValidateToken(token, cfg.SecretKey)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: " + err.Error()})
			c.Abort()
			return
		}

		c.Set(ContextKeyEmail, claims.Email)
		c.Set(ContextKeyFirstName, claims.FirstName)
		c.Set(ContextKeyLastName, claims.LastName)
		c.Set(ContextKeyRole, claims.Role)
		c.Set(ContextKeyUserID, claims.UserId)

		c.Next()
	}
}

func GetAccessToken(c *gin.Context) (string, error) {
	authHeader := c.Request.Header.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("authorization header is required")
	}

	// Check if the header starts with "Bearer "
	if strings.HasPrefix(authHeader, "Bearer ") {
		parts := strings.Split(authHeader, "Bearer ")
		if len(parts) == 2 {
			return strings.TrimSpace(parts[1]), nil
		}
	}

	// Fallback: assume the whole header is the token if "Bearer " is missing
	return strings.TrimSpace(authHeader), nil
}

func GetUserIdFromContext(c *gin.Context) (string, error) {
	userId, exists := c.Get(ContextKeyUserID)
	if !exists {
		return "", errors.New("user_id not found in context")
	}
	userIdStr, ok := userId.(string)
	if !ok {
		return "", errors.New("user_id is not a string")
	}
	return userIdStr, nil
}

func GetRoleFromContext(c *gin.Context) (string, error) {
	role, exists := c.Get(ContextKeyRole)
	if !exists {
		return "", errors.New("role not found in context")
	}
	roleStr, ok := role.(string)
	if !ok {
		return "", errors.New("role is not a string")
	}
	return roleStr, nil
}
