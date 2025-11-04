package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/botbooker/botbooker/internal/health"
	observability "github.com/botbooker/botbooker/internal/observability/otel"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func TestPingEndpoint(t *testing.T) {
	// Отключаем логи Gin в тестах
	gin.SetMode(gin.TestMode)

	// Создаём тестовый роутер
	router := setupRouter()

	ctx := context.Background()

	// Создаём запрос
	req, err := http.NewRequestWithContext(ctx, "GET", "/ping", nil)
	if err != nil {
		t.Fatal("не удалось создать запрос:", err)
	}
	w := httptest.NewRecorder()

	// Выполняем запрос
	router.ServeHTTP(w, req)

	// Проверяем статус
	if w.Code != http.StatusOK {
		t.Errorf("Ожидаемый статус %d, получен %d", http.StatusOK, w.Code)
	}

	// Проверяем Content-Type
	contentType := w.Header().Get("Content-Type")
	if !strings.Contains(contentType, "application/json") {
		t.Errorf("Content-Type должен быть application/json, получен %s", contentType)
	}

	// Декодируем ответ
	var response map[string]string
	err = json.NewDecoder(w.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Не удалось декодировать JSON: %v", err)
	}

	// Проверяем поле "message"
	if response["message"] != "pong" {
		t.Errorf("Поле 'message' должно быть 'pong', получено %s", response["message"])
	}
}

func TestHealthEndpoint(t *testing.T) {
	// Отключаем логи Gin в тестах
	gin.SetMode(gin.TestMode)
	// Создаём тестовый роутер
	router := setupRouter()

	ctx := context.Background()

	req, err := http.NewRequestWithContext(ctx, "GET", "/health", nil)
	if err != nil {
		t.Fatal("не удалось создать запрос:", err)
	}
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Ожидаемый статус %d, получен %d", http.StatusOK, w.Code)
	}

	contentType := w.Header().Get("Content-Type")
	if !strings.Contains(contentType, "application/json") {
		t.Errorf("Content-Type должен быть application/json, получен %s", contentType)
	}

	// Проверяем, что ответ не пустой
	if w.Body.Len() == 0 {
		t.Error("Ответ /health пуст")
	}
}

func setupRouter() *gin.Engine {
	api := gin.Default()
	api.ContextWithFallback = true
	api.Use(otelgin.Middleware(applicationName))

	api.GET("/ping", func(ctx *gin.Context) {
		traceID, spanID, isSampled := observability.GetTraceInfo(ctx)
		fmt.Printf("traceID: %v; spanID: %v; isSampled: %v\n", traceID, spanID, isSampled)
		ctx.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	api.GET("/health", health.HealthHandler)

	return api
}

func TestPingLogsTraceInfo(t *testing.T) {
	gin.SetMode(gin.TestMode)

	var logOutput strings.Builder
	observability.LogHook = func(format string, a ...interface{}) {
		logOutput.WriteString(fmt.Sprintf(format, a...))
	}

	router := setupRouter() // ваш существующий setupRouter

	ctx := context.Background()

	// Создаём запрос
	req, err := http.NewRequestWithContext(ctx, "GET", "/ping", nil)
	if err != nil {
		t.Fatal("не удалось создать запрос:", err)
	}
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if !strings.Contains(logOutput.String(), "traceID:") {
		t.Error("В логах отсутствует traceID")
	}
	if !strings.Contains(logOutput.String(), "spanID:") {
		t.Error("В логах отсутствует spanID")
	}

	// Сброс хука после теста
	observability.LogHook = nil
}
