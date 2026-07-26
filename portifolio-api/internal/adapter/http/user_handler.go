package handler

import (
	"net/http"
	usecaseuser "portifolio-api/internal/usecase/user"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	createUser usecaseuser.CreateUserUseCase
}

func NewUserHandler(createUser usecaseuser.CreateUserUseCase) *UserHandler {
	return &UserHandler{createUser: createUser}
}

func (h UserHandler) CreateUser(c *gin.Context) {
	var input usecaseuser.CreateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := h.createUser.Execute(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"id":      user.ID,
		"api_key": user.APIKey,
		"name":    user.Name,
	})
}
