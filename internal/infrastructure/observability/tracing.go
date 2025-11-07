package observability

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

var LogHook func(string, ...interface{})

func SetupTracing() {
	otel.SetTracerProvider(sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.ParentBased(sdktrace.AlwaysSample())),
	))
}

func GetTraceInfo(ctx *gin.Context) (traceID string, spanID string, isSampled bool) {
	// Извлекаем контекст из запроса
	spanCtx := trace.SpanContextFromContext(ctx.Request.Context())

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
