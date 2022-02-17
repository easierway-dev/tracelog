package tracelog

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/hashicorp/consul/api"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"gitlab.mobvista.com/mtech/tracelog/logevent"
	"time"
	"unsafe"
)

var (
	KvNotFound        = errors.New("kv not found")
	GetConsulKvFailed = errors.New("get consul kv info failed")
)

type Ops struct {
	Type     string
	Address  string
	Path     string
	Interval time.Duration
	TryTimes int
	OnChange func(value interface{}, err error) bool
}
type ConsulConfig struct {
	// SampleRatio 取样比例, 用于初始TraceIDRatioSampler
	// HasRemoteParent为true时, 可以设置parentBasedSampler的root为TraceIDRatioSampler
	SampleRatio         float64		`json:"SampleRatio" toml:"SampleRatio"`
	JaegerAgentEndpoint string 		`json:"JaegerAgentEndpoint" toml:"JaegerAgentEndpoint"`
	JaegerAgentHost     string
	JaegerAgentPort     string
	RootService         []string	`json:"RootService" toml:"RootService"`
	LoggingExporter *LoggingExporter `json:"LoggingExporter" toml:"LoggingExporter"`
}
type LoggingExporter struct {
	ExporterType string 	`json:"ExporterType" toml:"ExporterType"`
	ElasticSearchUrl string `json:"ElasticSearchUrl" toml:"ElasticSearchUrl"`
	KafkaUrl	string		`json:"KafkaUrl" toml:"KafkaUrl"`
}

func getTomlConfig(ops *Ops, value interface{}) error {
	// 获取toml配置文件中的值
	pair, err := GetValue(ops)
	if err != nil {
		return err
	}
	if pair == nil {
		return KvNotFound
	}
	// 将配置文件的值与consulConfig进行绑定
	if _, err = toml.Decode(*(*string)(unsafe.Pointer(&pair.Value)), value); err != nil {
		return err
	}
	return nil
}
func GetValue(ops *Ops) (*api.KVPair, error) {
	config := api.DefaultConfig()
	config.Address = ops.Address
	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}
	kv := client.KV()
	if kv == nil {
		return nil, GetConsulKvFailed
	}
	pair, _, err := kv.Get(ops.Path, &api.QueryOptions{})
	if err != nil {
		return nil, err
	}
	if pair == nil {
		return nil, KvNotFound
	}
	return pair, nil
}

// 通过传入服务名,consul地址,consul_key,获取config配置
func FromConsulConfig(service_name string, consul_addr string, consul_key string) (*Config, error) {
	// 定义consulConfig配置
	var consulConfig ConsulConfig
	// 定义config配置
	var config *Config
	// 初始化Ops配置，传入配置文件格式,consul地址,key
	tomlFormat := &Ops{Type: "toml", Address: consul_addr, Path: consul_key}
	// 获取配置文件并初始化consulConfig
	getTomlConfig(tomlFormat, &consulConfig)
	// 如果rootServer不为空, service_name不在RootServer中时, sampleRatio设置为0
	sampleRatio := consulConfig.SampleRatio
	if len(consulConfig.RootService) > 0 {
		isRootService := false
		for _, svc := range consulConfig.RootService {
			if service_name == svc {
				isRootService = true
			}
		}
		if !isRootService {
			sampleRatio = 0.0
		}
	}

	// 根据consulConfig初始化config
	config, err := NewConfig(WithServiceName(service_name),
		WithSampleRatio(sampleRatio),
		WithJaegerAgentEndpoint(consulConfig.JaegerAgentEndpoint))
	if err != nil {
		fmt.Println("init traceconfig failed:", err.Error())
		return nil, err
	}
	// 给log_event包下的全局变量Logger赋值
	logevent.Logger = InitLogger(consulConfig.LoggingExporter)
	// 初始化OpenTelemetry SDK
	if err := Start(config); err != nil {
		fmt.Println("init tracelog start failed:", err.Error())
		return nil, err
	}
	return config, nil
}
// 初始化日志配置,仅为ES,Kafka,Stdout中的一种
func InitLogger(loggingExporter *LoggingExporter) *log.Logger{
    // 不配置LoggingExporter时，不会panic
    if loggingExporter == nil {
        return nil
    }
	switch loggingExporter.ExporterType {
	case logevent.ES:
		logger := logevent.AddES(loggingExporter.ElasticSearchUrl)
		return logger
	case logevent.Kafka:
		kafka := logevent.AddKafka(loggingExporter.KafkaUrl)
		return kafka
	case logevent.Stdout:
		stdout := logevent.AddStdout()
		return stdout
    // 不配置log exporter时, 不打印日志
	default:
		return nil
	}
}
