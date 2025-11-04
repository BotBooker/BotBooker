package otel

import (
	"context"

	"go.opentelemetry.io/otel/trace"
)

var LogHook func(string, ...interface{})

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
