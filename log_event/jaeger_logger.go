package log_event

import (
    "context"
    "fmt"
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/trace"
    "go.opentelemetry.io/otel/attribute"
)

type logSpanFlag int // 0: not-sampled; 1: use parent-span; 2: new span
const (
    logSpanNoSampled logSpanFlag = 0
    logSpanUseParent logSpanFlag = 1
    logSpanNewSpan   logSpanFlag = 2
)

type JaegerLogEvent struct {
    span trace.Span
    spanFlag logSpanFlag
    eventName string 
    attrs []attribute.KeyValue
}

type JaegerLogEventVec struct {
    jaegerLogEvent *JaegerLogEvent
}

func NewJaegerLogEventVec(ctx context.Context, name string) logEventVec {
    span, spanFlag := logSpanFromContext(ctx)
    if span == nil  || spanFlag == logSpanNoSampled {
        return NewNopLogEventVec()
    }
    // setup span
    jle := &JaegerLogEvent{
        span: span,
        spanFlag: spanFlag,
        eventName: name,
    }
    jleVec := &JaegerLogEventVec{jle}
    return jleVec
}

func (lev *JaegerLogEventVec)getLogEventWithLabelValues(m map[string]string) (*JaegerLogEvent, error) {
    if lev == nil || lev.jaegerLogEvent == nil {
        return nil, fmt.Errorf("invalid jaeger log event")
    }
    attrs := make([]attribute.KeyValue, len(m)+1)
    for k,v := range m {
        attrs = append(attrs, attribute.String(k, v))
    }
    lev.jaegerLogEvent.attrs = attrs
    return lev.jaegerLogEvent, nil

}

func (lev *JaegerLogEventVec) WithLabelValues(m map[string]string) logEvent {
    le, err := lev.getLogEventWithLabelValues(m)
    // when error, return nopLogEvent
    if err != nil {
        return nopLogEvent{}
    }
    
    return le 
}

func (le *JaegerLogEvent) Log(msg string) {
    if le.span == nil {
        return 
    }

    if le.spanFlag == logSpanNewSpan {
        defer le.span.End()
    }
    
    // to do: LabelValues to attr 
    // msg to body
    le.attrs = append(le.attrs, attribute.String("event.message", msg))
        
    // use jaeger span event as logger writer
    le.span.AddEvent(le.eventName, trace.WithAttributes(le.attrs...))
}

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

