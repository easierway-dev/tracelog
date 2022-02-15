package logevent

import (
	"context"
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"reflect"
	"testing"
)

func TestLogWithContext(t *testing.T) {
	t.Parallel()
	tp:=sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()))
	otel.SetTracerProvider(tp)
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()
	tr := otel.Tracer("testlogrus")
	ctx, span := tr.Start(context.Background(), "test")
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
			name:        "new valid context, sampled return nil",
			ctx:         ctx,
			logEventVec: &LogrusLogEventVec{logrusLogEvent: &LogrusLogEvent{
				span: span,
				spanFlag: logSpanFlag(span.SpanContext().TraceFlags()),
				spanID: span.SpanContext().SpanID(),
				traceID: span.SpanContext().TraceID(),
				eventName: "testLog",
				kafkaTopic: []string{"trace_log"},
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			withContext := WithContext(tt.ctx, "testLog")
			if !reflect.DeepEqual(withContext, tt.logEventVec) {
				t.Errorf("Extract Tracecontext: %s: NewLogrusLogEventVec() returned %#v",tt.name,tt.logEventVec)
			}
		})
	}
}