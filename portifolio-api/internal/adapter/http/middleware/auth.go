package middleware

import (
	"net/http"
	"portifolio-api/internal/domain"

	"github.com/gin-gonic/gin"
)

func Auth(userRepo domain.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-Key")
		if apiKey == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		user, err := userRepo.FindByAPIKey(apiKey)
		if err != nil || user == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		c.Set("userID", user.ID)
		c.Set("telegramID", user.TelegramID)
		c.Set("user", user)
		c.Next()
	}
}
