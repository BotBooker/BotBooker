package handlers

import (
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestPingHandler_Handle(t *testing.T) {
	// Инициализируем Gin в тестовом режиме
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	fmt.Printf("context: %+v\n", c)
	handler := NewPingHandler()
	handler.Handle(c)
}
