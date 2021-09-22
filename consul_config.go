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
	NotSupportType      = errors.New("not support type")
	KvNotFound          = errors.New("kv not found")
	GetConsulKvFailed   = errors.New("get consul kv info failed")
	JsonUnmarshalFailed = errors.New("json unmarshal failed")
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
	pair, err := GetValue(ops)
	if err != nil {
		return err
	}
	if pair == nil {
		return KvNotFound
	}
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
func FromConsulConfig(service_name string, consul_addr string, consul_key string) (*Config, error) {
	var consulConfig ConsulConfig
	var config *Config
	tomlFormat := &Ops{Type: "toml", Address: consul_addr, Path: consul_key}
	getTomlConfig(tomlFormat, &consulConfig)
	config, err := NewConfig(WithServiceName(service_name),
		WithSampleRatio(consulConfig.SampleRatio),
		WithJaegerAgentEndpoint(consulConfig.JaegerAgentEndpoint))
	if err != nil {
		fmt.Println("init traceconfig failed:", err.Error())
		return nil, err
	}
	if err := Start(config); err != nil {
		fmt.Println("init tracelog start failed:", err.Error())
		return nil, err
	}
	return config, nil
}
