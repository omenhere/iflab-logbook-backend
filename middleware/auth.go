package middleware

import (
	"log"
	"net/http"
	"auth-backend/utils"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Retrieve the "jwt" cookie
		cookie, err := c.Cookie("jwt")
		if err != nil {
			log.Printf("Error retrieving JWT: %v", err) // Debug log
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: No token found"})
			c.Abort()
			return
		}

		// Parse and validate the JWT token
		userID, err := utils.ParseToken(cookie)
		if err != nil {
			log.Printf("Error validating token: %v", err) // Debug log
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: Invalid token"})
			c.Abort()
			return
		}

		// Log userID for debugging
		log.Printf("Authenticated user ID: %s", userID)

		// Store userID in the context for use in controllers
		c.Set("user_id", userID)
		c.Next()
	}
}
