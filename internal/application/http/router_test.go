package http

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/botbooker/botbooker/internal/application/http/handlers"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRouter_SetupRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Создаём моковые обработчики
	pingHandler := &handlers.PingHandler{}
	healthHandler := &handlers.HealthHandler{}

	router := NewRouter(pingHandler, healthHandler)
	api := gin.New()

	// Настраиваем маршруты
	router.SetupRoutes(api)

	// Тестируем /ping
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/ping", nil)
	api.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.True(t, strings.Contains(w.Body.String(), `"message":"pong"`))

	// Тестируем /health
	w = httptest.NewRecorder()
	req = httptest.NewRequest("GET", "/health", nil)
	api.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	// Здесь можно проверить ожидаемый JSON, если мокнуть HealthService
}
