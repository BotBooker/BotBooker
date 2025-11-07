package http

import (
	"github.com/botbooker/botbooker/internal/application/http/handlers"
	"github.com/gin-gonic/gin"
	otelgin "go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

type Router struct {
	pingHandler   *handlers.PingHandler
	healthHandler *handlers.HealthHandler
}

func NewRouter(
	pingHandler *handlers.PingHandler,
	healthHandler *handlers.HealthHandler,
) *Router {
	return &Router{
		pingHandler:   pingHandler,
		healthHandler: healthHandler,
	}
}

func (r *Router) SetupRoutes(api *gin.Engine) {
	api.Use(otelgin.Middleware("botbooker-api"))

	api.GET("/ping", r.pingHandler.Handle)
	api.GET("/health", r.healthHandler.Handle)
}
