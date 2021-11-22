package tracelog

import (
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"time"
	//    "fmt"
	"strings"

	tracesdk "go.opentelemetry.io/otel/sdk/trace"

	"context"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"os"
)

func ToRecordingContext(ctx context.Context, sc trace.SpanContext) {
     if !sc.IsValid() {
         return
     }
     prop := propagation.TraceContext{}
     flags := trace.FlagsSampled & trace.FlagsSampled
     th := fmt.Sprintf(
          "%.2x-%s-%s-%s",
          supportedVersion,
          sc.TraceID(),
          sc.SpanID(),
          flags,
      )
     req, _ := http.NewRequest("GET", "http://example.com", nil)
     req.Header.Set("traceparent", th)

     ctx = prop.Extract(ctx, propagation.HeaderCarrier(req.Header))
     return
}
