package health

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {
	// Создаем новый HTTP-запрос
	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Создаем ответ
	w := httptest.NewRecorder()

	// Создаем контекст Gin
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Вызываем обработчик
	Handler(c)

	// Проверяем код состояния
	assert.Equal(t, http.StatusOK, w.Code)

	// Проверяем тело ответа
	assert.JSONEq(t, `{"message": "OK"}`, w.Body.String())
}
