package handler

import (
	"github.com/gin-gonic/gin"
)

const testTelegramID = "123456"
const testUserID int64 = 1

func injectUser(c *gin.Context) {
	c.Set("userID", testUserID)
	c.Set("telegramID", testTelegramID)
	c.Next()
}

func floatPtr(v float64) *float64 { return &v }
