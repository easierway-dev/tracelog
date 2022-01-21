package tracelog

/*
   一些参考: https://github.com/aliyun-sls/opentelemetry-go-provider-sls
*/
import (
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"time"
	//    "fmt"
	"strings"

	tracesdk "go.opentelemetry.io/otel/sdk/trace"

	"context"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"os"
)

type Config struct {
	ServiceName string
	// exporter采用log, 进行统一
	// 后续有其他exporter的话, 或collector的话, 再做修改, 代价不大
	TraceLogPath string

	// IsRootSpan 是否有root span的服务

	// 以下两个配置, 进行采样的配置
	// 使用otel/sdk/trace中定义的sampler, ParentBased和RatioBased
	// 只接受上游传来的服务, 应使用ParentBasedSampler,取样的决定
	// 由context传来, 如juno服务
	// 本身是请求的起始服务, 需要初始化RatioBasedSampler, 传入取样比例
	// 注意: 像dsp这种, 则需要配置ParentBased和RatioBased两种取样配置

	// HasRemoteParent 是否初始化ParentBasedSampler
	// 置为true时, 优先从远程上下文中传递取样配置
	// 为了省事, 暂时不暴露配置, 默认为true
	HasRemoteParent bool

	// SampleRatio 取样比例, 用于初始TraceIDRatioSampler
	// HasRemoteParent为true时, 可以设置parentBasedSampler的root为TraceIDRatioSampler
	SampleRatio float64

	JaegerAgentEndpoint string
	JaegerAgentHost     string
	JaegerAgentPort     string
	Resource            *resource.Resource
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
func WithJaegerAgentEndpoint(url string) Option {
	return func(c *Config) {
		if url == ""{
			url ="127.0.0.1:6831"
		}
		host := strings.Split(url, ":")[0]
		port := strings.Split(url, ":")[1]

		c.JaegerAgentEndpoint = url
		c.JaegerAgentHost = host
		c.JaegerAgentPort = port
	}
}

// 设置tracelog的日志路径
func WithTraceLogPath(path string) Option {
	return func(c *Config) {
		c.TraceLogPath = path
	}
}
// 设置是否有RemoteParent
func WithRemoteParent(b bool) Option {
	return func(c *Config) {
		c.HasRemoteParent = b
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

//
// 初始化Exporter
// 如果jaeger endpoint不为空, 则优先使用jaeger作为exporter
// 如果没配jaeger endpoint, 使用stdouttrace将trace信息输出到trace日志中
// 日志使用logrus
func (c *Config) initOtelExporter() (traceExporter tracesdk.SpanExporter, exporterStop func(), initErr error) {

	var exporter tracesdk.SpanExporter
	var err error
	if c.JaegerAgentEndpoint != "" {
		exporter, err = jaeger.New(jaeger.WithAgentEndpoint(
			jaeger.WithAgentHost(c.JaegerAgentHost),
			jaeger.WithAgentPort(c.JaegerAgentPort),
		))
		if err != nil {
			return nil, nil, errors.Wrap(err, "init jaegertrace exporter faild")
		}

	} else {
		logger, err := initLogger(c.TraceLogPath)
		if err != nil {
			return nil, nil, errors.Wrap(err, "init logTracer failed!")
		}
		exporter, err = stdouttrace.New(stdouttrace.WithWriter(logger.Writer()))

		if err != nil {
			return nil, nil, errors.Wrap(err, "init stdouttrace exporter faild")
		}
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
		tracesdk.WithSampler(tracesdk.ParentBased(tracesdk.TraceIDRatioBased(c.SampleRatio))),
		// Record information about this application in an Resource.
		tracesdk.WithResource(c.Resource),
	)

	otel.SetTracerProvider(tp)
	// 设置上下文传递
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	c.stop = append(c.stop, func() {
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
