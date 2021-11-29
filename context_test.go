package tracelog_test

import (
	"context"
	"go.opentelemetry.io/otel/trace"
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
func TestIsSampledFromContext(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		ctx context.Context
		isSuccess    bool
	}{
		{
			name: "new invalid context, return false",
			ctx: nil,
			isSuccess: false,
		},
		{
			name: "new invalid context, return false",
			ctx: context.Background(),
			isSuccess: false,
		},
		{
			name: "new invalid context, return false",
			ctx: trace.ContextWithRemoteSpanContext(context.Background(),trace.NewSpanContext(trace.SpanContextConfig{Remote: true})),
			isSuccess: false,
		},
		{
			name: "new valid context, return true",
			ctx: trace.ContextWithRemoteSpanContext(context.Background(),trace.NewSpanContext(trace.SpanContextConfig{
				TraceID:    mustTraceIDFromHex(traceIDStr),
				SpanID:     mustSpanIDFromHex(spanIDStr),
				TraceFlags: trace.FlagsSampled,
				Remote:     true,
			})),
			isSuccess: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCtx := ContextToRecordingContext(tt.ctx)
			if IsSampledFromContext(gotCtx) != tt.isSuccess {
				t.Errorf("Extract Tracecontext: %s: IsSampledFromContext() returned %#v",tt.name,tt.isSuccess)
			}
		})
	}
}