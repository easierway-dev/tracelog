package logevent

import (
	"fmt"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

// 获取span中的Attributes
func GetAttributes(span trace.Span) map[string]string {
	m := make(map[string]string)
	if readSpan, ok := span.(tracesdk.ReadWriteSpan); ok {
		for _, value := range readSpan.Attributes() {
			m[string(value.Key)] = value.Value.AsString()
		}
	}
	return m
}

// 获取span中的Resource
func GetResource(span trace.Span) map[string]string {
	m := make(map[string]string)
	if readSpan, ok := span.(tracesdk.ReadWriteSpan); ok {
		for _, value := range readSpan.Resource().Attributes() {
			m[string(value.Key)] = value.Value.AsString()
		}
	}
	if len(m) == 0 {
		fmt.Errorf("resources return nil")
	}
	return m
}
