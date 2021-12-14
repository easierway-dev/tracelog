package otelrpcx

import (
    "context"
    "net"
    "go.opentelemetry.io/otel/trace"
    	"github.com/smallnest/rpcx/protocol"
        "go.opentelemetry.io/otel"
        "github.com/smallnest/rpcx/share"
        "github.com/smallnest/rpcx/server"
        "go.opentelemetry.io/otel/attribute"
        
)


type OpenTelemetryServerPlugin struct{}

func (p OpenTelemetryServerPlugin)Register(name string, rcvr interface{}, metadata string) error {
	return nil
}

func (p OpenTelemetryServerPlugin)RegisterFunction(serviceName, fname string, fn interface{}, metadata string) error {
	return nil
}

func (p OpenTelemetryServerPlugin) PostConnAccept(conn net.Conn) (net.Conn, bool) {
	return conn, true
}

func (p OpenTelemetryServerPlugin) PreHandleRequest(ctx context.Context, r *protocol.Message) error {
    sc := GetOtelSpanContextFromRpcxContext(ctx)
    if sc == nil {
        return nil
    }

    tracer := otel.GetTracerProvider().Tracer(instrumentationName)
    spanName := "rpcx.service."+r.ServicePath+"."+r.ServiceMethod
    clientConn := ctx.Value(server.RemoteConnContextKey).(net.Conn)
    opts := []trace.SpanStartOption{
        trace.WithAttributes(attribute.String("remote_addr", clientConn.RemoteAddr().String())),
        trace.WithSpanKind(trace.SpanKindServer),
    }

    _,span := tracer.Start(ctx, spanName, opts...)
    if rpcxContext, ok := ctx.(*share.Context); ok {
		rpcxContext.SetValue(OpenTelemetrySpanRequestKey, span)
    }
	return nil

}

func (p OpenTelemetryServerPlugin)  PostWriteResponse(ctx context.Context, req *protocol.Message, res *protocol.Message, err error) error {
	if rpcxContext, ok := ctx.(*share.Context); ok {
		span1 := rpcxContext.Value(OpenTelemetrySpanRequestKey)
		if span1 != nil {
			span1.(trace.Span).End()
		}
	}
	return nil
}

