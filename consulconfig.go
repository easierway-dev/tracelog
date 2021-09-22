package tracelog

type ConsulConfig struct {
	// SampleRatio 取样比例, 用于初始TraceIDRatioSampler
	// HasRemoteParent为true时, 可以设置parentBasedSampler的root为TraceIDRatioSampler
	SampleRatio         float64
	JaegerAgentEndpoint string
	JaegerAgentHost     string
	JaegerAgentPort     string
}

func FromConsulConfig(service_name string, consul_addr string, consul_key string) (*Config, error) {
	var consulConfig ConsulConfig
	var config *Config
	tomlFormat := &Ops{Type: "toml", Address: consul_addr, Path: consul_key}
	getTomlConfig(tomlFormat, &consulConfig)
	config, err := NewConfig(WithServiceName(service_name),
		WithSampleRatio(consulConfig.SampleRatio),
		WithJaegerAgentEndpoint(consulConfig.JaegerAgentEndpoint))
	return config, err
}
