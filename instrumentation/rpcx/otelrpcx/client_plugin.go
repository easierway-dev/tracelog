package otelrpcx

import (
    "context"
    "github.com/smallnest/rpcx/share"
    "go.opentelemetry.io/otel/trace"
)

// 需要参考opencensus
type OpenTelemetryClientPlugin struct{}

func (p *OpenTelemetryClientPlugin)PreCall(ctx context.Context, servicePath, serviceMethod string, args interface{}) error{
    tracer := otel.GetTracerProvider("otelrpcx")
    spanName := "rpcx.client."+servicePath+"."+serviceMethod
    opsts := []trace.SpanStartOption{
        trace.WithSpanKind(trace.SpanKindServer),
    }

    ctx ,span := tracer.Start(ctx, spanName, opts...)
    if rpcxContext, ok := ctx.(*share.Context); ok {
		rpcxContext.SetValue(share.OpencensusSpanClientKey, span)
	}
	return nil

}

func (p *OpenTelemetryClientPlugin)PostCall(ctx context.Context, servicePath, serviceMethod string, args interface{}, reply interface{}, err error) error {
	if rpcxContext, ok := ctx.(*share.Context); ok {
		span1 := rpcxContext.Value(share.OpencensusSpanClientKey)
		if span1 != nil {
			span1.(*trace.Span).End()
		}
	}
	return nil
}
