package otelgrpc

import (
    "context"
)

func GetOtelSpenContextFromRpcxContext(ctx context.Context) trace.SpanContext {
    reqMeta, ok := ctx.Value(ReqMetaDataKey).(map[string]string)
    if ! ok {
        reutrn nil
    }
    spanKey := reqMeta[OpenTelemetrySpanRequestKey]
    if spanKey == "" {
        return nil
    }
    th := []byte(spanKey)

    prop := propagation.TraceContext{}
    header := make(http.Header)
    header.Set(traceparentHeader, th)

    return  prop.Extract(ctx, propagation.HeaderCarrier(header))
}
