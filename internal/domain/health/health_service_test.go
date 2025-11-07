package health

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHealthService_Check(t *testing.T) {
	service := NewHealthService()
	ctx := context.Background()

	result := service.Check(ctx)

	assert.Equal(t, "healthy", result["status"])
	assert.Equal(t, "botbooker-api", result["service"])
	assert.Len(t, result, 2)
}
