package observability

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

func TestSetupTracing(t *testing.T) {
	originalProvider := otel.GetTracerProvider()
	SetupTracing()
	currentProvider := otel.GetTracerProvider()

	if currentProvider == originalProvider {
		t.Errorf("SetupTracing() не установил новый TracerProvider")
	}

	_, ok := currentProvider.(*sdktrace.TracerProvider)
	if !ok {
		t.Errorf("TracerProvider имеет неверный тип: %T", currentProvider)
	}
}

func TestGetTraceInfo_NoSpanContext(t *testing.T) {
	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatalf("не удалось создать HTTP-запрос: %v", err)
	}
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req

	traceID, spanID, isSampled := GetTraceInfo(ctx)

	if traceID != "" {
		t.Errorf("traceID должен быть пустым, получил: %q", traceID)
	}
	if spanID != "" {
		t.Errorf("spanID должен быть пустым, получил: %q", spanID)
	}
	if isSampled {
		t.Error("isSampled должен быть false при отсутствии SpanContext")
	}
}

func TestGetTraceInfo_WithSpanContext(t *testing.T) {
	spanContext := trace.NewSpanContext(trace.SpanContextConfig{
		TraceID:    trace.TraceID{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10},
		SpanID:     trace.SpanID{0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18},
		TraceFlags: trace.FlagsSampled,
	})

	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatalf("не удалось создать HTTP-запрос: %v", err)
	}
	req = req.WithContext(trace.ContextWithSpanContext(context.Background(), spanContext))

	// Пересоздаём Gin-контекст полностью, чтобы он "увидел" новый контекст запроса
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req

	traceID, spanID, isSampled := GetTraceInfo(ctx)

	expectedTraceID := spanContext.TraceID().String()
	if traceID != expectedTraceID {
		t.Errorf("traceID не соответствует ожидаемому. Ожидалось: %s, получено: %s", expectedTraceID, traceID)
	}

	expectedSpanID := spanContext.SpanID().String()
	if spanID != expectedSpanID {
		t.Errorf("spanID не соответствует ожидаемому. Ожидалось: %s, получено: %s", expectedSpanID, spanID)
	}

	if !isSampled {
		t.Error("isSampled должен быть true для sampled SpanContext")
	}
}

func TestGetTraceInfo_LogHookCalled(t *testing.T) {
	var logHookCalled bool
	var capturedTraceID, capturedSpanID string
	var capturedIsSampled bool

	LogHook = func(format string, args ...interface{}) {
		logHookCalled = true
		capturedTraceID = args[0].(string)
		capturedSpanID = args[1].(string)
		capturedIsSampled = args[2].(bool)
	}

	spanContext := trace.NewSpanContext(trace.SpanContextConfig{
		TraceID:    trace.TraceID{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10},
		SpanID:     trace.SpanID{0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18},
		TraceFlags: trace.FlagsSampled,
	})

	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatalf("не удалось создать HTTP-запрос: %v", err)
	}
	req = req.WithContext(trace.ContextWithSpanContext(context.Background(), spanContext))

	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req

	GetTraceInfo(ctx)

	if !logHookCalled {
		t.Error("LogHook не был вызван")
	}

	expectedTraceID := spanContext.TraceID().String()
	expectedSpanID := spanContext.SpanID().String()

	if capturedTraceID != expectedTraceID {
		t.Errorf("LogHook получил неверный traceID. Ожидалось: %s, получено: %s", expectedTraceID, capturedTraceID)
	}
	if capturedSpanID != expectedSpanID {
		t.Errorf("LogHook получил неверный spanID. Ожидалось: %s, получено: %s", expectedSpanID, capturedSpanID)
	}
	if !capturedIsSampled {
		t.Error("LogHook получил неверный isSampled. Ожидалось: true")
	}
}

func TestGetTraceInfo_LogHookNotCalledWhenNil(t *testing.T) {
	LogHook = nil

	spanContext := trace.NewSpanContext(trace.SpanContextConfig{
		TraceID: trace.TraceID{0x01, 0x02, 0x03, 0x04},
		SpanID:  trace.SpanID{0x05, 0x06, 0x07, 0x08},
	})

	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatalf("не удалось создать HTTP-запрос: %v", err)
	}
	req = req.WithContext(trace.ContextWithSpanContext(context.Background(), spanContext))

	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req

	GetTraceInfo(ctx)
}
