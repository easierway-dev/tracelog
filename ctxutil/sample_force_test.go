package ctxutil_test

import (
	"context"
	"github.com/google/go-cmp/cmp"
	"gitlab.mobvista.com/mtech/tracelog/ctxutil"
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

func TestNonSampledContextToSampledContext(t *testing.T) {
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
			gotCtx := ctxutil.ContextToRecordingContext(ctx)
			gotSc := trace.SpanContextFromContext(gotCtx)
			if diff := cmp.Diff(gotSc, tt.wantSc, cmp.Comparer(func(sc, other trace.SpanContext) bool { return sc.Equal(other) })); diff != "" {
				t.Errorf("Extract Tracecontext: %s: -got +want %s", tt.name, diff)
			}

		})
	}
}
