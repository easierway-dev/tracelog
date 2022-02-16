package logevent

import (
	"fmt"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/trace"
	"reflect"
)
// 参考链接https://blog.csdn.net/raoxiaoya/article/details/112607719
// 获取span中的Attributes
func GetAttributes(span trace.Span) map[string]string{
	attributes1 := make(map[string]string)
	getType := reflect.TypeOf(span)
	fmt.Println("类型具体信息:"+getType.String())
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
	}
	return attributes1
}
// 获取span中的Resource
func GetResource(span trace.Span) map[string]string{
	resources := make(map[string]string)
	getType := reflect.TypeOf(span)
	fmt.Println("类型具体信息:"+getType.String())
	getValue := reflect.ValueOf(span)
	// 反射获取Resource的属性值
	if _, ok := getType.MethodByName("Resource"); ok {
		m := make(map[string]string)
		handler := getValue.MethodByName("Resource")
		val := handler.Call([]reflect.Value{})
		intf := val[0].Interface()
		resource1 := intf.(*resource.Resource)
		attributes := resource1.Attributes()
		for i := 0; i < len(attributes ); i++ {
			m[string(attributes[i].Key)] = attributes[i].Value.AsString()
		}
	}
	return resources
}