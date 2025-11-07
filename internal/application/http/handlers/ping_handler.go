package handlers

import (
	"fmt"
	"net/http"

	observability "github.com/botbooker/botbooker/internal/infrastructure/observability"
	"github.com/gin-gonic/gin"
)

type PingHandler struct{}

func NewPingHandler() *PingHandler {
	return &PingHandler{}
}

func (h *PingHandler) Handle(ctx *gin.Context) {
	traceID, spanID, isSampled := observability.GetTraceInfo(ctx)
	fmt.Printf("traceID: %v; spanID: %v; isSampled: %v\n", traceID, spanID, isSampled)

	ctx.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}
