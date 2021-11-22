package tracelog

import (
	"go.opentelemetry.io/otel/propagation"
	"fmt"

	"go.opentelemetry.io/otel/trace"
    "net/http"

	"context"
)
 const   supportedVersion  = 0

func ToRecordingContext(ctx context.Context, sc trace.SpanContext) {
     if !sc.IsValid() {
         return
     }
     prop := propagation.TraceContext{}
     th := fmt.Sprintf(
          "%.2x-%s-%s-%s",
          supportedVersion,
          sc.TraceID(),
          sc.SpanID(),
          trace.FlagsSampled,
      )
     req, _ := http.NewRequest("GET", "http://example.com", nil)
     req.Header.Set("traceparent", th)

     ctx = prop.Extract(ctx, propagation.HeaderCarrier(req.Header))
     return
}
