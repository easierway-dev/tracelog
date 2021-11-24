package tracelog 

import (
    "go.opentelemetry.io/otel/trace"
)

type logSpanFlag // 0: not-sampled; 1: use parent-span; 2: new span
const (
    logSpanNoSampled logSpanFlag = 0
    logSpanUseParent logSpanFlag = 1
    logSpanNewSpan   logSpanFlag = 2
)

// 适用jeager来记录日志
// return value:
// span: 记录log的span, IsSampled=false, 为nil; 
// IsNewSpan: 是否有采样的父span, 适用在处理中途才决定采样的情况(defered sampled)
// IsSampled: 是否采样
func logSpanFromContext(ctx context.Context) (trace.Span, logSpanFlag) {
    // context invalid, not sample
    if ctx == nil {
        return nil, logSpanNoSampled
    }

    // 从上下文中获取当前span的信息
    span := trace.SpanFromContext(ctx)

    // 使用SpanContext来判断是否IsSampled=true
    // IsSampled = false, 不采样 
    if !span.SpanContext().IsSampled() {
        return nil, logSpanNoSampled
    }
    
    // 父span采样, 
    if span.IsRecording() {
        return span, logSpanUseParent
    }
    
    // 父span不采样, 创建一个span
    _, span = otel.Tracer("tracelog-log").Start(ctx, "tracelog-log")
    return span, logSpanNewSpan
}

// 上下文信息, 因为使用jaeger作为日志收集, 本质仍是trace, 因此需要上下文, 其他日志系统的话, 不需要
// 暴露出去的日志接口, attrs为label支持索引
// 内部使用jaeger的AddEvent方法来记录日志, 与jeager耦合, 后续有其他日志引起再做解耦
// TO-DO: 如果jaeger不适用于记录日志, 则需要更换后面的日志引擎, 
func LogContext(ctx context.Context, logStr string, attrs ...attribute.KeyValue) {
    span, flag := logSpanFromContext(ctx)

    switch flag {
    case logSpanNoSampled:
    case logSpanUseParent:
        span.AddEvent(logStr, trace.WithAttributes(attrs))
    case logSpanNewSpan:
        span.AddEvent(logStr, trace.WithAttributes(attrs))
    default:
    }
    return
}
