package logevent

import (
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/trace"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
)
// 获取span中的Resource
func GetAttributes(span trace.Span) map[string]string {
	readSpan := span.(tracesdk.ReadWriteSpan)
	m := make(map[string]string)
	for _, value := range readSpan.Attributes() {
		m[string(value.Key)] = value.Value.AsString()
	}
	return m
}

func GetResource(span trace.Span) map[string]string {
	readSpan := span.(tracesdk.ReadWriteSpan)
	m := make(map[string]string)
	for _, value := range readSpan.Resource().Attributes() {
		m[string(value.Key)] = value.Value.AsString()
	}
	return m
}