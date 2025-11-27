// Package health предоставляет обработчики для проверок работоспособности сервиса.
//
// Включает:
//   - Handler: эндпоинт /health для проверки статуса сервиса.
package health

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Handler — обработчик HTTP‑запроса для проверки здоровья сервиса.
//
// Возвращает JSON‑ответ с статусом 200 OK и сообщением "OK".
// Используется для:
//   - мониторинга доступности сервиса (например, в Kubernetes liveness probe);
//   - базовых проверок работоспособности API.
//
// Параметры:
//   - ctx *gin.Context: контекст HTTP‑запроса от Gin.
func Handler(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "OK",
	})
}
