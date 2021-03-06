package logevent

import (
	"context"
	"fmt"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"reflect"
)

type logSpanFlag int // 0: not-sampled; 1: use parent-span; 2: new span
const (
	logSpanNoSampled logSpanFlag = 0
	logSpanUseParent logSpanFlag = 1
	logSpanNewSpan   logSpanFlag = 2
)

type JaegerLogEvent struct {
	span      trace.Span
	spanFlag  logSpanFlag
	eventName string
	attrs     []attribute.KeyValue
}

type JaegerLogEventVec struct {
	jaegerLogEvent *JaegerLogEvent
}

func NewJaegerLogEventVec(ctx context.Context, name string) logEventVec {
	span, spanFlag := logSpanFromContext(ctx)
	if span == nil || spanFlag == logSpanNoSampled {
		return NewNopLogEventVec()
	}
	// setup span
	jle := &JaegerLogEvent{
		span:      span,
		spanFlag:  spanFlag,
		eventName: name,
	}
	jleVec := &JaegerLogEventVec{jle}
	return jleVec
}

func (lev *JaegerLogEventVec) getLogEventWithLabelValues(m map[string]string) (*JaegerLogEvent, error) {
	if lev == nil || lev.jaegerLogEvent == nil {
		return nil, fmt.Errorf("invalid jaeger log event")
	}
	attrs := make([]attribute.KeyValue, len(m)+1)
	for k, v := range m {
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

func (le *JaegerLogEvent) Log(msg interface{}) {
	if le.span == nil {
		return
	}

	if le.spanFlag == logSpanNewSpan {
		defer le.span.End()
	}
	// to do: LabelValues to attr
	// msg to body
	switch reflect.TypeOf(msg).Kind() {
	case reflect.Int:
		le.attrs = append(le.attrs, attribute.Int("event.message", msg.(int)))
		break
	case reflect.String:
		le.attrs = append(le.attrs, attribute.String("event.message", msg.(string)))
		break
	default:break
	}
	// use jaeger span event as logger writer
	le.span.AddEvent(le.eventName, trace.WithAttributes(le.attrs...))
}

// ??????jeager???????????????
// return value:
// span: ??????log???span, IsSampled=false, ???nil;
// IsNewSpan: ?????????????????????span, ?????????????????????????????????????????????(defered sampled)
// IsSampled: ????????????
func logSpanFromContext(ctx context.Context) (trace.Span, logSpanFlag) {
	// context invalid, not sample
	if ctx == nil {
		return nil, logSpanNoSampled
	}

	// ???????????????????????????span?????????
	span := trace.SpanFromContext(ctx)

	// ??????SpanContext???????????????IsSampled=true
	// IsSampled = false, ?????????
	if !span.SpanContext().IsSampled() {
		return nil, logSpanNoSampled
	}

	// ???span??????,
	if span.IsRecording() {
		return span, logSpanUseParent
	}

	// ???span?????????, ????????????span
	_, span = otel.Tracer("tracelog-log").Start(ctx, "tracelog-log")
	return span, logSpanNewSpan
}
