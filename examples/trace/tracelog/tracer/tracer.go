package tracer

/*
   一些参考: https://github.com/aliyun-sls/opentelemetry-go-provider-sls
*/
import (
	"github.com/pkg/errors"
	"time"
	//    "fmt"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"

	"context"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"os"
)

type Config struct {
	ServiceName  string
	TraceLogPath string
	SampleRatio  float64

	Resource *resource.Resource

	resourceAttributes map[string]string

	stop []func()
}

type Option func(*Config)

// 设置servicename, 'service.name' resource attribute
func WithServiceName(name string) Option {
	return func(c *Config) {
		c.ServiceName = name
	}
}

// 设置tracelog的日志路径
func WithTraceLogPath(path string) Option {
	return func(c *Config) {
		c.TraceLogPath = path
	}
}

// 设置tracelog取样的比例
func WithSampleRatio(ratio float64) Option {
	return func(c *Config) {
		c.SampleRatio = ratio
	}
}

func (c *Config) IsValid() error {
	if c.ServiceName == "" {
		return errors.New("empty service name")
	}
	if c.TraceLogPath == "" {
		return errors.New("empty trace log path")
	}
	if c.SampleRatio == 0 {
		return errors.New("zero sample ratio")
	}
	return nil

}

func getDefaultResource(c *Config) *resource.Resource {
	hostname, _ := os.Hostname()
	return resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(c.ServiceName),
		semconv.HostNameKey.String(hostname),
	)
}

func mergeResource(c *Config) {
	c.Resource, _ = resource.Merge(getDefaultResource(c), c.Resource)
	var keyValues []attribute.KeyValue
	for key, value := range c.resourceAttributes {
		keyValues = append(keyValues, attribute.KeyValue{
			Key:   attribute.Key(key),
			Value: attribute.StringValue(value),
		})
	}
	newResource := resource.NewWithAttributes(semconv.SchemaURL, keyValues...)
	c.Resource, _ = resource.Merge(c.Resource, newResource)

}

func NewConfig(opts ...Option) (*Config, error) {
	var c Config

	// load config from option function
	for _, opt := range opts {
		opt(&c)
	}

	mergeResource(&c)
	return &c, c.IsValid()
}

// 初始化Exporter, 使用stdouttrace将trace信息输出到trace日志中
// 日志使用logrus
func (c *Config) initOtelExporter() (traceExporter tracesdk.SpanExporter, exporterStop func(), initErr error) {
	logger, err := initLogger(c.TraceLogPath)
	if err != nil {
		return nil, nil, errors.Wrap(err, "init logTracer failed!")
	}
	exporter, err := stdouttrace.New(stdouttrace.WithWriter(logger.Writer()))
	if err != nil {
		return nil, nil, errors.Wrap(err, "init stdouttrace faild")
	}
	traceExporter = exporter
	exporterStop = func() {
		exporter.Shutdown(context.Background())
	}
	return traceExporter, exporterStop, nil
}

// 初始化tracer
func (c *Config) initTracer(traceExporter tracesdk.SpanExporter, stop func()) error {
	if traceExporter == nil {
		return errors.New("no trace exporter")
	}

	tp := tracesdk.NewTracerProvider(
		// Always be sure to batch in production.
		// record trace log one-by-one
		tracesdk.WithBatcher(traceExporter,
			tracesdk.WithMaxExportBatchSize(1),
			tracesdk.WithExportTimeout(time.Second),
		),
		// init sample
		// tracesdk.WithSampler(tracesdk.AlwaysSample()),
		tracesdk.WithSampler(tracesdk.TraceIDRatioBased(c.SampleRatio)),
		// Record information about this application in an Resource.
		tracesdk.WithResource(c.Resource),
	)

	otel.SetTracerProvider(tp)
	c.stop = append(c.stop, func() {
		tp.Shutdown(context.Background())
		stop()
	})
	return nil
}

// Start 初始化OpenTelemetry SDK
func Start(c *Config) error {
	tracerExport, traceExpStop, err := c.initOtelExporter()
	if err != nil {
		return errors.Wrap(err, "start trace failed")
	}
	err = c.initTracer(tracerExport, traceExpStop)
	if err != nil {
		return err
	}
	return err
}

// Shutdown 优雅关闭，将OpenTelemetry SDK内存中的数据发送到服务端
func Shutdown(c *Config) {
	for _, stop := range c.stop {
		stop()

	}

}

/*
func trace1() {
	tp := otel.GetTracerProvider().(*tracesdk.TracerProvider)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	tr := tp.Tracer("component-main")

	ctx, span := tr.Start(ctx, "foo")
	fmt.Println("trace on:", span.SpanContext().IsSampled())
	defer span.End()

	bar(ctx)
	bar(ctx)
	time.Sleep(3 * time.Second)
}
*/
