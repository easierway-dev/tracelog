package otelrpcx

import (
    "context"
    "github.com/smallnest/rpcx/share"
    "go.opentelemetry.io/otel/propagation"
    "net/http"
)

func GetOtelSpanContextFromRpcxContext(ctx context.Context) context.Context {
    reqMeta, ok := ctx.Value(share.ReqMetaDataKey).(map[string]string)
    if !ok {
        return nil
    }
    spanKey := reqMeta[OpenTelemetrySpanRequestKey]
    if spanKey == "" {
        return nil
    }

    prop := propagation.TraceContext{}
    header := make(http.Header)
    header.Set(traceparentHeader, spanKey)

    return  prop.Extract(ctx, propagation.HeaderCarrier(header))
}
