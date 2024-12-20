package middleware

import (
	"net/http"
	"strconv"
	"strings"

	"online-learning-golang/utils"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token required"})
			c.Abort()
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		userId, role, err := utils.ValidToken(tokenStr)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		userIdStr := strconv.Itoa(userId)
		c.Set("userId", userIdStr)
		c.Set("role", role)
		c.Next()
	}
}

func OnlyAdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token required"})
			c.Abort()
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		userId, role, err := utils.ValidToken(tokenStr)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		if role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied: Admins only"})
			c.Abort()
			return
		}

		userIdStr := strconv.Itoa(userId)
		c.Set("userId", userIdStr)
		c.Set("role", role)
		c.Next()
	}
}
