package logevent

import (
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/trace"
	"reflect"
)
// 参考链接https://blog.csdn.net/raoxiaoya/article/details/112607719
// 获取span中的Resource
func GetResource(span trace.Span) map[string]string{
	m := make(map[string]string)
	getType := reflect.TypeOf(span)
	getValue := reflect.ValueOf(span)
	if _, ok := getType.MethodByName("Resource"); ok {
		handler := getValue.MethodByName("Resource")
		val := handler.Call([]reflect.Value{})
		intf := val[0].Interface()    // interface{} 类型
		resource1 := intf.(*resource.Resource) // bool 类型
		for _,value:=range resource1.Attributes(){
			m[string(value.Key)] = value.Value.AsString()
		}
	}
	return m
}
// 获取span中的Attributes
func GetAttributes(span trace.Span) map[string]string{
	m := make(map[string]string)
	getType := reflect.TypeOf(span)
	getValue := reflect.ValueOf(span)
	if _, ok := getType.MethodByName("Attributes"); ok {
		handler := getValue.MethodByName("Attributes")
		val := handler.Call([]reflect.Value{})
		intf := val[0].Interface()
		for _,value:=range intf.([]attribute.KeyValue){
			m[string(value.Key)] = value.Value.AsString()
		}
	}
	return m
}