package middleware

import (
	"net/http"
	"net/http/httptest"
	"portifolio-api/internal/domain"
	"portifolio-api/internal/mocks"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupAuthRouter(userRepo domain.UserRepository) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(Auth(userRepo))
	r.GET("/test", func(c *gin.Context) {
		userID := c.GetInt64("userID")
		c.JSON(http.StatusOK, gin.H{"userID": userID})
	})
	return r
}

func TestAuth_Success(t *testing.T) {
	userRepo := new(mocks.UserRepository)
	user := &domain.User{ID: 42, TelegramID: "123", APIKey: "valid-key"}
	userRepo.On("FindByAPIKey", "valid-key").Return(user, nil)

	r := setupAuthRouter(userRepo)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("X-API-Key", "valid-key")

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "42")
	userRepo.AssertExpectations(t)
}

func TestAuth_SemHeader(t *testing.T) {
	userRepo := new(mocks.UserRepository)

	r := setupAuthRouter(userRepo)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	userRepo.AssertNotCalled(t, "FindByAPIKey")
}

func TestAuth_ChaveInvalida(t *testing.T) {
	userRepo := new(mocks.UserRepository)
	userRepo.On("FindByAPIKey", "wrong-key").Return(nil, nil)

	r := setupAuthRouter(userRepo)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("X-API-Key", "wrong-key")

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	userRepo.AssertExpectations(t)
}

func TestAuth_RepoError(t *testing.T) {
	userRepo := new(mocks.UserRepository)
	userRepo.On("FindByAPIKey", "some-key").Return(nil, assert.AnError)

	r := setupAuthRouter(userRepo)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("X-API-Key", "some-key")

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	userRepo.AssertExpectations(t)
}
