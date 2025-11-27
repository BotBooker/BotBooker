// Package main реализует HTTP‑сервер API для сервиса бронирования botbooker.
//
// Функционал:
//   - обработка REST‑запросов;
//   - интеграция с базой данных;
//   - аутентификация пользователей.
package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"

	"github.com/botbooker/botbooker/internal/health"
	observability "github.com/botbooker/botbooker/internal/observability/otel"
)

var applicationName = "botbooker-api"

func main() {
	otel.SetTracerProvider(sdktrace.NewTracerProvider(sdktrace.WithSampler(sdktrace.ParentBased(sdktrace.AlwaysSample()))))
	// Create a Gin router with default middleware (logger and recovery)
	api := gin.Default()
	api.ContextWithFallback = true
	api.Use(otelgin.Middleware(applicationName))

	// Define a simple GET endpoint
	api.GET("/ping", func(ctx *gin.Context) {
		traceID, spanID, isSampled := observability.GetTraceInfo(ctx)
		fmt.Printf("traceID: %v; spanID: %v; isSampled: %v\n", traceID, spanID, isSampled)
		// Return JSON response
		ctx.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	api.GET("/health", health.Handler)

	// Start server on port 8080 (default)
	// Server will listen on 0.0.0.0:8080 (localhost:8080 on Windows)
	if err := api.Run(); err != nil {
		fmt.Printf("cannot start API server: %s", err)
	}
}
