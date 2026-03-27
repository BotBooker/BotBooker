package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/uptrace/bun"

	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"

	"github.com/botbooker/botbooker/internal/apis"
	"github.com/botbooker/botbooker/internal/controllers"
	db "github.com/botbooker/botbooker/internal/database"
	"github.com/botbooker/botbooker/internal/health"
	"github.com/botbooker/botbooker/internal/middleware"
	observability "github.com/botbooker/botbooker/internal/observability/otel"
)

var (
	applicationName = "botbooker-api"
	dbController    *bun.DB
)

func main() {
	otel.SetTracerProvider(sdktrace.NewTracerProvider(sdktrace.WithSampler(sdktrace.ParentBased(sdktrace.AlwaysSample()))))
	dbConfig := db.DBConfig{
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASS"),
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Database: os.Getenv("DB_NAME"),
	}
	dbConn, dbErr := db.ConnectDB(&dbConfig)
	if dbErr != nil {
		log.Fatal("error connecting with database\n")
	}
	dbController = dbConn
	fmt.Println("Connected to database successfully")
	defer dbController.Close()
	db.RegisterModels(dbController)
	userController := controllers.UserControllerInitializer(dbConn)
	router := gin.Default()
	router.Use(middleware.CorsMiddleware())
	router.MaxMultipartMemory = 64 << 20 // 64 MiB
	router.ContextWithFallback = true
	router.Use(otelgin.Middleware(applicationName))
	api := router.Group("/api/v1")
	{
		api.GET("/", func(ctx *gin.Context) {
			traceID, spanID, isSampled := observability.GetTraceInfo(ctx)
			fmt.Printf("traceID: %v; spanID: %v; isSampled: %v\n", traceID, spanID, isSampled)
			ctx.JSON(200, gin.H{
				"message": "server ap1/v1 Running ..  ",
			})
		})
		apis.AddApiUsers(api, userController)
	}
	router.GET("/ping", func(ctx *gin.Context) {
		traceID, spanID, isSampled := observability.GetTraceInfo(ctx)
		fmt.Printf("traceID: %v; spanID: %v; isSampled: %v\n", traceID, spanID, isSampled)
		ctx.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	router.GET("/health", health.Handler)
	if err := router.Run(":8080"); err != nil {
		fmt.Printf("cannot start API server: %s", err)
	}
}
