package main

import (
	"context"
	"gitlab.mobvista.com/mvbjqa/otel-practices/example/trace/tracelog/tracer"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"time"
)

func main() {
	traceLogPath := "../trace.log"
	sampleRatio := 1.0
	tracerConfig, err := tracer.NewConfig(tracer.WithServiceName("trace_demo"),
		tracer.WithTraceLogPath(traceLogPath),
		tracer.WithSampleRatio(sampleRatio),
	)
	// 初始化失败, panic
	if err != nil {
		panic(err)
	}
	//
	if err := tracer.Start(tracerConfig); err != nil {
		panic(err)
	}
	defer tracer.Shutdown(tracerConfig)
	trace1()
}
func trace1() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	tr := otel.Tracer("component-main")

	ctx, span := tr.Start(ctx, "foo")
	// fmt.Println("trace on:", span.SpanContext().IsSampled())
	defer span.End()

	bar(ctx)
	time.Sleep(1 * time.Second)
}
func bar(ctx context.Context) {
	// Use the global TracerProvider.
	tr := otel.Tracer("component-bar")
	_, span := tr.Start(ctx, "bar")
	span.SetAttributes(attribute.Key("testset").String("value"))
	defer span.End()

	// Do bar...

}
