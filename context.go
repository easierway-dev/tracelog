func tracelog

import (
    "context"
    "go.opentelemetry.io/otel/trace"
)

func IsSampledFromContext(ctx context.Context) bool {
         // context invalid, not sample
     if ctx == nil {
         return false
     }

     // 从上下文中获取当前span的信息
     span := trace.SpanFromContext(ctx)
     return span.SpanContext().IsSampled()
}


