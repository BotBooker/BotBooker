package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/botbooker/botbooker/internal/domain/health"
)

type HealthHandler struct {
	service *health.HealthService
}

func NewHealthHandler(service *health.HealthService) *HealthHandler {
	return &HealthHandler{service: service}
}

func (h *HealthHandler) Handle(ctx *gin.Context) {
	result := h.service.Check(ctx.Request.Context())
	ctx.JSON(http.StatusOK, result)
}
