package otelrpcx

import (
	"context"
	rpcxshare "github.com/smallnest/rpcx/share"
	"gitlab.mobvista.com/mtech/tracelog"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

func CallWithTrace(ctx context.Context, xclient XClient, serviceMethod string, args interface{}, reply interface{}) error {
	// drs inject trace to context, rpcx plugin not good for write plugin, so use manual code
	tracer := otel.GetTracerProvider().Tracer(instrumentationName)
	opts := []trace.SpanStartOption{
		trace.WithSpanKind(trace.SpanKindClient),
	}

	ctxDrs, span := tracer.Start(ctx1, serviceMethod, opts...)

	defer span.End()

	if traceParent := tracelog.TraceParenetFromContext(ctxDrs); traceParent != "" {
		ctxDrs = context.WithValue(ctxDrs, rpcxshare.ReqMetaDataKey, map[string]string{"traceparent": traceParent})
	}
	return xclient.Call(ctxDrs, "WrapRank", args, reply)
}
