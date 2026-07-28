package handler

import "github.com/gin-gonic/gin"

func respond(c *gin.Context, status int, data any) {
	c.JSON(status, gin.H{
		"telegram_id": c.GetString("telegramID"),
		"data":        data,
	})
}

func respondMessage(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{
		"telegram_id": c.GetString("telegramID"),
		"message":     message,
	})
}

func respondError(c *gin.Context, status int, err error) {
	c.JSON(status, gin.H{
		"telegram_id": c.GetString("telegramID"),
		"error":       err.Error(),
	})
}

