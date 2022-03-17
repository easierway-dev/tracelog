package logevent

import (
	"context"
	"go.opentelemetry.io/otel"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"log"
	"testing"
)

func mustTraceIDFromHex(s string) (t trace.TraceID) {
	var err error
	t, err = trace.TraceIDFromHex(s)
	if err != nil {
		panic(err)
	}
	return
}

func mustSpanIDFromHex(s string) (t trace.SpanID) {
	var err error
	t, err = trace.SpanIDFromHex(s)
	if err != nil {
		panic(err)
	}
	return
}
func TestLogrusLoggerNoSample(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		ctx         context.Context
		logEventVec logEventVec
	}{
		{
			name:        "new invalid context, return nil",
			ctx:         nil,
			logEventVec: nopLogEventVec{},
		},
		{
			name:        "new invalid context, return nil",
			ctx:         context.Background(),
			logEventVec: nopLogEventVec{},
		},
		{
			name:        "new valid context, non sampled return nil",
			ctx:         trace.ContextWithRemoteSpanContext(context.Background(), trace.NewSpanContext(trace.SpanContextConfig{Remote: true})),
			logEventVec: nopLogEventVec{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if NewLogrusLogEventVec(tt.ctx, "testLogrus") != tt.logEventVec {
				t.Errorf("Extract Tracecontext: %s: NewLogrusLogEventVec() returned %#v", tt.name, tt.logEventVec)
			}
		})
	}
}
func TestLogrusLoggerSample(t *testing.T) {
	t.Parallel()
	Logger = AddStdout()
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()))
	otel.SetTracerProvider(tp)
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()
	tr := otel.Tracer("testlogrus")
	_, span := tr.Start(context.Background(), "test")
	tests := []struct {
		name   string
		ctx    context.Context
		sample logSpanFlag
		m      map[string]string
	}{
		{
			name:   "new valid context, parent sampled return 1",
			ctx:    trace.ContextWithSpan(context.Background(), span),
			sample: logSpanUseParent,
			m:      map[string]string{"peer.service": "ExampleService1"},
		},
		{
			name: "new valid context,customize sampled return 2",
			ctx: trace.ContextWithRemoteSpanContext(context.Background(), trace.NewSpanContext(trace.SpanContextConfig{
				TraceID:    mustTraceIDFromHex(traceIDStr),
				SpanID:     mustSpanIDFromHex(spanIDStr),
				TraceFlags: trace.FlagsSampled,
				Remote:     true,
			})),
			sample: logSpanNewSpan,
			m:      map[string]string{"peer.service": "ExampleService2"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lle := NewLogrusLogEventVec(tt.ctx, "test")
			if lle.(*LogrusLogEventVec).logrusLogEvent.spanFlag != tt.sample {
				t.Errorf("Extract Tracecontext: %s: NewLogrusLogEventVec() returned %#v", tt.name, tt.sample)
			}
			le := lle.WithLabelValues(tt.m)
			le.Log("testSuccess")
		})
	}
}
