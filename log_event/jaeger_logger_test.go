package log_event

import (
	"context"
	"go.opentelemetry.io/otel/propagation"
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
type traceContextKeyType int

const currentSpanKey traceContextKeyType = iota
func TestJaegerLoggerNoSample(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		ctx context.Context
		logEventVec logEventVec
	}{
		{
			name: "new invalid context, return nil",
			ctx: nil,
			logEventVec: nopLogEventVec{},
		},
		{
			name: "new invalid context, return nil",
			ctx: context.Background(),
			logEventVec: nopLogEventVec{},
		},
		{
			name: "new valid context, non sampled return nil",
			ctx: trace.ContextWithRemoteSpanContext(context.Background(),trace.NewSpanContext(trace.SpanContextConfig{Remote: true})),
			logEventVec: nopLogEventVec{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if NewJaegerLogEventVec(tt.ctx, "test") != tt.logEventVec {
				t.Errorf("Extract Tracecontext: %s: NewJaegerLogEventVec() returned %#v",tt.name,tt.logEventVec)
			}
		})
	}
}

func TestJaegerLoggerSample(t *testing.T) {
	t.Parallel()
	start, _ := otel.Tracer("a").Start(context.Background(), "test")
	tests := []struct {
		name string
		ctx context.Context
		sample logSpanFlag
	}{
		// 这个还有点问题
		{
			name: "new valid context, parent sampled return 1",
			ctx: context.WithValue(
				context.WithValue(context.Background(),currentSpanKey,
					trace.NewSpanContext(trace.SpanContextConfig{
						TraceFlags: trace.FlagsSampled,
						Remote:     true,
					})),currentSpanKey,trace.SpanFromContext(start),
			),
			sample: logSpanUseParent,
		},
		{
			name: "new valid context,customize sampled return 2",
			ctx: trace.ContextWithRemoteSpanContext(context.Background(),trace.NewSpanContext(trace.SpanContextConfig{
				TraceID:    mustTraceIDFromHex(traceIDStr),
				SpanID:     mustSpanIDFromHex(spanIDStr),
				TraceFlags: trace.FlagsSampled,
				Remote:     true,
			})),
			sample: logSpanNewSpan,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vec := NewJaegerLogEventVec(tt.ctx, "test")
			if vec.(*JaegerLogEventVec).jaegerLogEvent.spanFlag != tt.sample {
				t.Errorf("Extract Tracecontext: %s: NewJaegerLogEventVec() returned %#v",tt.name,tt.sample)
			}
		})
	}
}