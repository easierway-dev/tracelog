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
func GetResource(span trace.Span) map[string]string{
	//var resource1 *resource.Resource
	m := make(map[string]string)
	getType := reflect.TypeOf(span)
	getValue := reflect.ValueOf(span)
	if _, ok := getType.MethodByName("Resource"); ok {
		handler := getValue.MethodByName("Resource")
		val := handler.Call([]reflect.Value{})
		intf := val[0].Interface()    // interface{} 类型
		resource1 := intf.(*resource.Resource) // bool 类型
		for key,value:=range resource1.Attributes(){
			//value := intf.([]attribute.KeyValue)
			fmt.Printf("Attributes:%d %s %s",key,value.Key,value.Value.AsString())
			m[string(value.Key)] = value.Value.AsString()
		}
		fmt.Printf(" %v:\n",m)
	}
	return m
}

func GetAttributes(span trace.Span) map[string]string{
	//var resource1 *resource.Resource
	m := make(map[string]string)
	getType := reflect.TypeOf(span)
	fmt.Println("类型具体信息:"+getType.String())
	getValue := reflect.ValueOf(span)
	if _, ok := getType.MethodByName("Attributes"); ok {
		handler := getValue.MethodByName("Attributes")
		val := handler.Call([]reflect.Value{})
		intf := val[0].Interface()
		//keyvalue := intf.([]attribute.KeyValue)
		for key,value:=range intf.([]attribute.KeyValue){
			//value := intf.([]attribute.KeyValue)
			fmt.Printf("Attributes:%d %s %s",key,value.Key,value.Value.AsString())
			m[string(value.Key)] = value.Value.AsString()
		}
		fmt.Printf(" %v:\n",m)
	}
	return m
}