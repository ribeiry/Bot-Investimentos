package handler

import (
	"time"

	"portifolio-api/internal/domain"

	"github.com/gin-gonic/gin"
)

const testTelegramID = "123456"
const testUserID int64 = 1

var testUser = &domain.User{
	ID:         testUserID,
	TelegramID: testTelegramID,
	Name:       "Test User",
	APIKey:     "test-key",
	CreatedAt:  time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC),
}

func injectUser(c *gin.Context) {
	c.Set("userID", testUserID)
	c.Set("telegramID", testTelegramID)
	c.Set("user", testUser)
	c.Next()
}

func floatPtr(v float64) *float64 { return &v }
