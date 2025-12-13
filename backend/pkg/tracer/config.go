package tracer

import "time"

// Protocol はOTLPエクスポーターの通信方式を示します。
type Protocol string

const (
	ProtocolGRPC Protocol = "grpc"
	ProtocolHTTP Protocol = "http"
)

// Config はOpenTelemetry TracerProviderを初期化するための設定です。
type Config struct {
	ServiceName        string
	ServiceVersion     string
	Environment        string
	Endpoint           string
	Protocol           Protocol
	Insecure           bool
	Headers            map[string]string
	Timeout            time.Duration
	SampleRatio        float64
	ResourceAttributes map[string]string
}

func (c *Config) normalize() {
	if c.Endpoint == "" {
		c.Endpoint = "localhost:4317"
	}
	if c.Protocol == "" {
		c.Protocol = ProtocolGRPC
	}
	if c.Timeout == 0 {
		c.Timeout = 5 * time.Second
	}
	if c.SampleRatio <= 0 || c.SampleRatio > 1 {
		c.SampleRatio = 1
	}
	if c.Headers == nil {
		c.Headers = map[string]string{}
	}
	if c.ResourceAttributes == nil {
		c.ResourceAttributes = map[string]string{}
	}
}
