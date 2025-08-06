package middleware

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the expected auth token from environment variable
		expectedToken := os.Getenv("AUTH_TOKEN")
		if expectedToken == "" {
			expectedToken = "admin-secret-token-2024" // Default token for development
		}

		// Get the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header is required",
			})
			c.Abort()
			return
		}

		// Check if the token matches
		if authHeader != expectedToken {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid authorization token",
			})
			c.Abort()
			return
		}

		// Continue to the next handler
		c.Next()
	}
}