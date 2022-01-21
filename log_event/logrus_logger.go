package log_event

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"gitlab.mobvista.com/mtech/tracelog"
	"go.opentelemetry.io/otel/trace"
)

type LogrusLogEvent struct {
	span       trace.Span
	traceID    trace.TraceID
	spanID     trace.SpanID
	spanFlag   logSpanFlag
	attributes map[string]string
	resource   map[string]string
	kafkaTopic 	   []string
	eventName string
	logger 	   *log.Logger
}
const (
	ES string = "ES"
	Kafka string = "Kafka"
	Stdout   string = "Stdout"
)
type LogrusLogEventVec struct {
	logrusLogEvent *LogrusLogEvent
}

func InitLogger(loggingExporter *tracelog.LoggingExporter) *log.Logger{
	switch loggingExporter.ExporterType {
	case ES:
		logger := AddES(loggingExporter.ElasticSearchUrl)
		return logger
	case Kafka:
		kafka := AddKafka(loggingExporter.KafkaUrl)
		return kafka
	case Stdout:
		stdout := AddStdout()
		return stdout
	default:
		return log.New()
	}
}
func NewLogrusLogEventVec(ctx context.Context,name string) logEventVec {
	span, spanFlag := logSpanFromContext(ctx)
	if span == nil || spanFlag == logSpanNoSampled {
		return NewNopLogEventVec()
	}
	// setup span
	lle := &LogrusLogEvent{
		span:     span,
		spanFlag: spanFlag,
		traceID:  span.SpanContext().TraceID(),
		spanID:   span.SpanContext().SpanID(),
		logger: tracelog.Logger,
		eventName: name,
		kafkaTopic: []string{"trace_log"},
	}
	lleVec := &LogrusLogEventVec{lle}
	return lleVec
}

func (lev *LogrusLogEventVec) getLogEventWithLabelValues(m map[string]string) (*LogrusLogEvent, error) {
	if lev == nil || lev.logrusLogEvent == nil {
		return nil, fmt.Errorf("invalid logrus log event")
	}
	lev.logrusLogEvent.attributes = m
	return lev.logrusLogEvent, nil
}

func (lev *LogrusLogEventVec) WithLabelValues(m map[string]string) logEvent {
	le, err := lev.getLogEventWithLabelValues(m)
	// when error, return nopLogEvent
	if err != nil {
		return nopLogEvent{}
	}
	return le
}

func (le *LogrusLogEvent) Log(msg string) {
	if le.span == nil {
		return
	}

	if le.spanFlag == logSpanNewSpan {
		defer le.span.End()
	}
	le.logger.WithFields(log.Fields{
		"traceId": le.traceID.String(),
		"spanId": le.spanID.String(),
		"traceFlags":int(le.spanFlag),
		"attributes": le.attributes,
		"resource": le.resource,
		"kafkaTopic":le.kafkaTopic,
	}).Info(msg)
}