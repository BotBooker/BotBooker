// Package otel предоставляет утилиты для работы с OpenTelemetry трассировкой.
//
// Основные функции:
//   - GetTraceInfo: извлекает идентификаторы trace/span и флаг sampled из контекста.
//   - LogHook: опциональный хук для логирования данных трассировки.
package otel

import (
	"context"

	"go.opentelemetry.io/otel/trace"
)

// LogHook — опциональный хук для логирования информации о трассировке.
//
// Если установлен, вызывается внутри GetTraceInfo с форматированной строкой,
// содержащей traceID, spanID и isSampled.
//
// Пример установки:
//
//	otel.LogHook = func(format string, args ...any) {
//		log.Printf(format, args...)
//	}
var LogHook func(string, ...any)

// GetTraceInfo извлекает данные трассировки из контекста.
//
// Параметры:
//   - ctx context.Context: контекст, содержащий SpanContext OpenTelemetry.
//
// Возвращает:
//   - traceID string: идентификатор трассировки (пустая строка, если отсутствует);
//   - spanID string: идентификатор спана (пустая строка, если отсутствует);
//   - isSampled bool: флаг, указывающий, был ли спан выбран для сбора (sampled).
//
// Если LogHook установлен, функция вызывает его с отформатированной строкой,
// содержащей извлечённые данные.
func GetTraceInfo(ctx context.Context) (traceID string, spanID string, isSampled bool) {
	spanCtx := trace.SpanContextFromContext(ctx)

	if spanCtx.HasTraceID() {
		traceID = spanCtx.TraceID().String()
	}
	if spanCtx.HasSpanID() {
		spanID = spanCtx.SpanID().String()
	}

	isSampled = spanCtx.IsSampled()

	if LogHook != nil {
		LogHook("traceID: %v; spanID: %v; isSampled: %v", traceID, spanID, isSampled)
	}

	return traceID, spanID, isSampled
}
