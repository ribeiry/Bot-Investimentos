package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	usecaseuser "portifolio-api/internal/usecase/user"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupUserMeRouter(h *UserHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(injectUser)
	r.GET("/users/me", h.GetMe)
	return r
}

func TestUserHandler_GetMe_Success(t *testing.T) {
	h := NewUserHandler(usecaseuser.CreateUserUseCase{})
	r := setupUserMeRouter(h)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users/me", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var body map[string]any
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))

	assert.Equal(t, testTelegramID, body["telegram_id"])

	data := body["data"].(map[string]any)
	assert.Equal(t, float64(testUserID), data["id"])
	assert.Equal(t, testTelegramID, data["telegram_id"])
	assert.Equal(t, "Test User", data["name"])
	assert.NotEmpty(t, data["created_at"])

	_, hasAPIKey := data["api_key"]
	assert.False(t, hasAPIKey, "api_key não deve ser exposto no /users/me")
}

func TestUserHandler_GetMe_SemUsuarioNoContexto(t *testing.T) {
	h := NewUserHandler(usecaseuser.CreateUserUseCase{})

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/users/me", h.GetMe)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users/me", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
