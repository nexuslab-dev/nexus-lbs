package config

import (
	"github.com/ttys3/lgr"
	"os"
)

const (
	EnvOtlpGrpcEndpoint = "OTLP_GRPC_ENDPOINT"

	DefaultLogLevel    = "info"
	DefaultLogEncoding = LogEncodingJSON
	DefaultLogOutput   = "stderr"
)

type Config struct {
	HttpAddr string `toml:"http_addr" yaml:"http_addr"`
	GrpcAddr string `toml:"grpc_addr" yaml:"grpc_addr"`

	Tracing          bool   `toml:"tracing" yaml:"tracing"`
	Profile          bool   `toml:"profile" yaml:"profile"`
	Metric           bool   `toml:"metric" yaml:"metric"`
	OtlpGrpcEndpoint string `toml:"otlp_grpc_endpoint" yaml:"otlp_grpc_endpoint"`

	Log    LogConfig    `toml:"log" yaml:"log"`
	Sentry SentryConfig `toml:"sentry" yaml:"sentry"`
}

// log

type LogEncoding string

// the default is json, if not empty, must be one of: console|json
const (
	LogEncodingJSON    LogEncoding = "json"
	LogEncodingConsole LogEncoding = "console"
)

type LogConfig struct {
	// Level set log level, can be empty, or one of debug|info|warn|error|fatal|panic
	Level string `toml:"level" json:"level" yaml:"level"`
	// Output set output file path, can be filepath or stdout|stderr
	Output string `toml:"output" json:"output" yaml:"output"`
	// Encoding sets the logger's encoding. Valid values are "json" and "console". default: json
	Encoding LogEncoding `toml:"encoding" json:"encoding" yaml:"encoding"`
	// enable stacktrace
	DisableStacktrace bool `toml:"disable_stacktrace" json:"disable_stacktrace" yaml:"disable_stacktrace"`
}

func (b *Config) InitOtlpGrpcEndpointFromEnv() {
	if b.OtlpGrpcEndpoint != "" {
		return
	}
	if tmp := os.Getenv(EnvOtlpGrpcEndpoint); tmp != "" {
		b.OtlpGrpcEndpoint = tmp
		lgr.S().Info("otlp_grpc_endpoint is empty, get from env var", "otlp_grpc_endpoint", tmp)
	} else {
		lgr.S().Warn("otlp_grpc_endpoint is empty")
	}
}
