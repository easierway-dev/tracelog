package logevent

import (
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"go.opentelemetry.io/otel/trace"
	"reflect"
)
// 获取span中的AttributesAndResource
func GetAttributesAndResource(span trace.Span) map[string]interface{}{
	attributesAndResource := make(map[string]interface{})
	getType := reflect.TypeOf(span)
	getValue := reflect.ValueOf(span)
	// 反射获取Attributes的属性值
	if _, ok := getType.MethodByName("Attributes"); ok {
		m := make(map[string]string)
		handler := getValue.MethodByName("Attributes")
		val := handler.Call([]reflect.Value{})
		intf := val[0].Interface()
		attributes := intf.([]attribute.KeyValue)
		for i := 0; i < len(attributes); i++ {
			value := intf.([]attribute.KeyValue)
			m[string(value[i].Key)] = value[i].Value.AsString()
		}
		attributesAndResource["Attributes"] = m
	}
	// 反射获取Resource的属性值
	if _, ok := getType.MethodByName("Resource"); ok {
		m := make(map[string]string)
		handler := getValue.MethodByName("Resource")
		val := handler.Call([]reflect.Value{})
		intf := val[0].Interface()
		resource1 := intf.(*resource.Resource)
		attributes := resource1.Attributes()
		for i := 0; i < len(attributes ); i++ {
			if(attributes[i].Key) == semconv.ServiceNameKey{
				m[string(attributes[i].Key)] = attributes[i].Value.AsString()
			}
		}
		attributesAndResource["Resource"] = m
	}
	return attributesAndResource
}