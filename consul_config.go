package tracelog

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/hashicorp/consul/api"
	"github.com/pkg/errors"
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
	SampleRatio         float64
	JaegerAgentEndpoint string
	JaegerAgentHost     string
	JaegerAgentPort     string
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
	// 根据consulConfig初始化config
	config, err := NewConfig(WithServiceName(service_name),
		WithSampleRatio(consulConfig.SampleRatio),
		WithJaegerAgentEndpoint(consulConfig.JaegerAgentEndpoint))
	if err != nil {
		fmt.Println("init traceconfig failed:", err.Error())
		return nil, err
	}
	// 初始化OpenTelemetry SDK
	if err := Start(config); err != nil {
		fmt.Println("init tracelog start failed:", err.Error())
		return nil, err
	}
	return config, nil
}
