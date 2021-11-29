
package log_event

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
		ctx context.Context
	}{
		{
			ctx: nil,
		},
		{
			ctx: context.Background(),
		},
		{
			ctx: trace.ContextWithRemoteSpanContext(context.Background(),trace.NewSpanContext(trace.SpanContextConfig{
			TraceID:    mustTraceIDFromHex(traceIDStr),
			SpanID:     mustSpanIDFromHex(spanIDStr),
			TraceFlags: trace.FlagsSampled,
			Remote:     true,
		})),
		},
	}
	for _, tt := range tests {
		gotCtx := ContextToRecordingContext(tt.ctx)
		if got, want := IsSampledFromContext(gotCtx), false; got != want {
			t.Errorf("IsSampledFromContext returned %#v, want %#v", got, want)
		}
	}
}