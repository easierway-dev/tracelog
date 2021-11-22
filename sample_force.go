package tracelog

import (
	"fmt"
	"go.opentelemetry.io/otel/propagation"

	"go.opentelemetry.io/otel/trace"
	"net/http"

	"context"
)

const (
	supportedVersion  = 0
	traceparentHeader = "traceparent"
)

func ContextToRecordingContext(ctx context.Context) context.Context {
	sc := trace.SpanContextFromContext(ctx)

	if !sc.IsValid() {
		return ctx
	}

	prop := propagation.TraceContext{}
	th := fmt.Sprintf(
		"%.2x-%s-%s-%s",
		supportedVersion,
		sc.TraceID(),
		sc.SpanID(),
		trace.FlagsSampled,
	)
	header := make(http.Header)
	header.Set(traceparentHeader, th)

	ctx = prop.Extract(ctx, propagation.HeaderCarrier(header))
	return ctx
}
