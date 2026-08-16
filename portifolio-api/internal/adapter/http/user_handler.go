package handler

import (
	"errors"
	"net/http"
	"portifolio-api/internal/domain"
	usecaseuser "portifolio-api/internal/usecase/user"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	createUser usecaseuser.CreateUserUseCase
}

func NewUserHandler(createUser usecaseuser.CreateUserUseCase) *UserHandler {
	return &UserHandler{createUser: createUser}
}

func (h UserHandler) GetMe(c *gin.Context) {
	value, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	user := value.(*domain.User)
	c.JSON(http.StatusOK, gin.H{
		"telegram_id": user.TelegramID,
		"data": gin.H{
			"id":          user.ID,
			"telegram_id": user.TelegramID,
			"name":        user.Name,
			"created_at":  user.CreatedAt,
		},
	})
}

func (h UserHandler) CreateUser(c *gin.Context) {
	var input usecaseuser.CreateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := h.createUser.Execute(input)
	if err != nil {
		if errors.Is(err, domain.ErrTelegramIDAlreadyExists) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"id":      user.ID,
		"api_key": user.APIKey,
		"name":    user.Name,
	})
}
