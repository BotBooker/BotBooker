package main

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/botbooker/botbooker/internal/application/http"
	"github.com/botbooker/botbooker/internal/application/http/handlers"
	"github.com/botbooker/botbooker/internal/domain/health"
	"github.com/botbooker/botbooker/internal/infrastructure/observability"
)

func main() {
	// 1. Инфраструктура: настройка трассировки
	observability.SetupTracing()

	// 2. Домен: создаём сервисы
	healthService := health.NewHealthService()

	// 3. Приложение: создаём обработчики
	pingHandler := handlers.NewPingHandler()
	healthHandler := handlers.NewHealthHandler(healthService)

	// 4. Приложение: настраиваем роутер
	router := http.NewRouter(pingHandler, healthHandler)

	// 5. Интерфейс: инициализируем Gin
	api := gin.Default()
	api.ContextWithFallback = true

	// 6. Приложение: монтируем маршруты
	router.SetupRoutes(api)

	// 7. Интерфейс: запускаем сервер
	if err := api.Run(); err != nil {
		fmt.Printf("cannot start API server: %s", err)
	}
}
