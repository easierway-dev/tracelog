package otelrpcx

import (
    "fmt"
    "context"
    "github.com/smallnest/rpcx/share"
    "go.opentelemetry.io/otel/trace"
    "go.opentelemetry.io/otel"
)

// 需要参考opencensus
type OpenTelemetryClientPlugin struct{}

func (p *OpenTelemetryClientPlugin)PreCall(ctx context.Context, servicePath, serviceMethod string, args interface{}) error{
    fmt.Println("OpenTelemetryClientPlugin precall")
                 _, okk := ctx.(*share.Context)
             fmt.Println("OpenTelemetryClientPlugin plugin context rpcxcontext:", okk)

    tracer := otel.GetTracerProvider().Tracer("otelrpcx")
    spanName := "rpcx.client."+servicePath+"."+serviceMethod
    opts := []trace.SpanStartOption{
        trace.WithSpanKind(trace.SpanKindServer),
    }

    _,span := tracer.Start(ctx, spanName, opts...)
    if rpcxContext, ok := ctx.(*share.Context); ok {
		rpcxContext.SetValue("otel-rpcx", span)
        fmt.Println("rpcx context:",rpcxContext.Value("otel-rpcx"))
    } else {
        fmt.Println("context not rpcxcontext")
    }
        
	return nil

}

func (p *OpenTelemetryClientPlugin)PostCall(ctx context.Context, servicePath, serviceMethod string, args interface{}, reply interface{}, err error) error {
	if rpcxContext, ok := ctx.(*share.Context); ok {
		span1 := rpcxContext.Value("otel-rpcx")
		if span1 != nil {
			span1.(trace.Span).End()
		}
	}
	return nil
}

