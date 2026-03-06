package otel

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/trace"
)

func TestGetTraceInfo_WithSpanContext(t *testing.T) {
	// Создаем контекст с трассировкой
	ctx := context.Background()
	spanCtx := trace.NewSpanContext(trace.SpanContextConfig{
		TraceID:    trace.TraceID{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16},
		SpanID:     trace.SpanID{17, 18, 19, 20, 21, 22, 23, 24},
		TraceFlags: trace.TraceFlags(0x01),
	})
	ctx = trace.ContextWithSpanContext(ctx, spanCtx)

	// Вызываем функцию
	traceID, spanID, isSampled := GetTraceInfo(ctx)

	// Проверяем результат
	assert.Equal(t, "0102030405060708090a0b0c0d0e0f10", traceID)
	assert.Equal(t, "1112131415161718", spanID)
	assert.True(t, isSampled)
}

func TestGetTraceInfo_WithoutSpanContext(t *testing.T) {
	// Создаем контекст без трассировки
	ctx := context.Background()

	// Вызываем функцию
	traceID, spanID, isSampled := GetTraceInfo(ctx)

	// Проверяем результат
	assert.Equal(t, "", traceID)
	assert.Equal(t, "", spanID)
	assert.False(t, isSampled)
}

func TestGetTraceInfo_LogHook(t *testing.T) {
	// Устанавливаем лог-хук
	logged := false
	LogHook = func(format string, args ...any) {
		logged = true
	}

	// Создаем контекст с трассировкой
	ctx := context.Background()
	spanCtx := trace.NewSpanContext(trace.SpanContextConfig{
		TraceID:    trace.TraceID{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16},
		SpanID:     trace.SpanID{17, 18, 19, 20, 21, 22, 23, 24},
		TraceFlags: trace.TraceFlags(0x01),
	})
	ctx = trace.ContextWithSpanContext(ctx, spanCtx)

	// Вызываем функцию
	GetTraceInfo(ctx)

	// Проверяем, что хук был вызван
	assert.True(t, logged)
}
