package log_event

import (
	"context"
	"fmt"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"net/http"
	"testing"
)
const (
	traceIDStr = "4bf92f3577b34da6a3ce929d0e0e4736"
	spanIDStr  = "00f067aa0ba902b7"
)
var (
	traceID = mustTraceIDFromHex(traceIDStr)
	spanID  = mustSpanIDFromHex(spanIDStr)
)
func TestJaegerLogger(t *testing.T) {
	tests := []struct {
		name   string
		sc     trace.SpanContext
		wantSc trace.SpanContext
	}{
		{
			name:   "in valid context, non sampled -> non sampled",
			sc:     trace.SpanContext{},
			wantSc: trace.NewSpanContext(trace.SpanContextConfig{Remote: true}),
		},
		{
			name: "valid context, non sampled -> sampled",
			sc: trace.NewSpanContext(trace.SpanContextConfig{
				TraceID: traceID,
				SpanID:  spanID,
			}),
			wantSc: trace.NewSpanContext(trace.SpanContextConfig{
				TraceID:    traceID,
				SpanID:     spanID,
				TraceFlags: trace.FlagsSampled,
				Remote:     true,
			}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// ctx = trace.ContextWithRemoteSpanContext(ctx, sc)
			ctx := context.Background()
			ctx = trace.ContextWithRemoteSpanContext(ctx, tt.sc)
			gotCtx := ContextToRecordingContext(ctx)
			m := make(map[string]string)
			m["Label1"] = "tag1"
			m["Label2"] = "tag2"
			eventVec := WithContext(gotCtx, "testJaegerLogger")
			labelValues := eventVec.WithLabelValues(m)
			labelValues.Log("testlog")
		})
	}
}
