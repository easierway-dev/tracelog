package log_event

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"gitlab.mobvista.com/mtech/tracelog"
	"go.opentelemetry.io/otel/trace"
)
const (
	ES string = "ES"
	Kafka string = "Kafka"
	Stdout   string = "Stdout"
)
type LogrusLogEvent struct {
	span       trace.Span
	traceID    trace.TraceID
	spanID     trace.SpanID
	spanFlag   logSpanFlag
	attributes map[string]string
	resource   map[string]string
	logger 	   *log.Logger
}

type LogrusLogEventVec struct {
	logrusLogEvent *LogrusLogEvent
}

func InitExporter(loggingExporter *tracelog.LoggingExporter) (*tracelog.Config ,error){

	switch loggingExporter.ExporterType {
	case ES:
		config, err := tracelog.NewConfig(tracelog.WithExporterType(loggingExporter.ExporterType),
			tracelog.WithElasticSearchUrl(loggingExporter.ElasticSearchUrl))
		return config,err
	case Kafka:
		config, err := tracelog.NewConfig(tracelog.WithExporterType(loggingExporter.ExporterType),
			tracelog.WithKafkaUrl(loggingExporter.KafkaUrl))
		return config,err
	case Stdout:
		config, err := tracelog.NewConfig(tracelog.WithExporterType(loggingExporter.ExporterType))
		return config,err
	}
	return nil,tracelog.GetConsulKvFailed
}
func InitLogger(config *tracelog.Config) *log.Logger{
	switch config.ExporterType {
	case ES:
		logger := AddES(config.ElasticSearchUrl)
		return logger
	case Kafka:
		kafka := AddKafka(config.KafkaUrl)
		return kafka
	case Stdout:
		stdout := AddStdout()
		return stdout
	}
	return log.New()
}
func NewLogrusLogEventVec(ctx context.Context, config *tracelog.Config) logEventVec {
	logger := InitLogger(config)
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
		logger: logger,
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
		"topics": []string{"first_topic"},
	}).Info(msg)
}
