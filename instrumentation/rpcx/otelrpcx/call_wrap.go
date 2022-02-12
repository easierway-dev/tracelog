package otelrpcx

import (
	"context"
	rpcxshare "github.com/smallnest/rpcx/share"
	"gitlab.mobvista.com/mtech/tracelog/ctxutil"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
    "github.com/smallnest/rpcx/client"
)

func CallWithTrace(ctx context.Context, xclient client.XClient, serviceMethod string, args interface{}, reply interface{}) error {
	// drs inject trace to context, rpcx plugin not good for write plugin, so use manual code
	tracer := otel.GetTracerProvider().Tracer(instrumentationName)
	opts := []trace.SpanStartOption{
		trace.WithSpanKind(trace.SpanKindClient),
	}

	ctxDrs, span := tracer.Start(ctx, serviceMethod, opts...)

	defer span.End()

	if traceParent := ctxutil.TraceParenetFromContext(ctxDrs); traceParent != "" {
		ctxDrs = context.WithValue(ctxDrs, rpcxshare.ReqMetaDataKey, map[string]string{"traceparent": traceParent})
	}
	return xclient.Call(ctxDrs, "WrapRank", args, reply)
}
