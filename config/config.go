package config

import (
	"errors"
	"log/slog"
	"os"

	"github.com/ttys3/slogx"

	v10validator "github.com/go-playground/validator/v10"
)

const (
	EnvOtlpGrpcEndpoint = "OTLP_GRPC_ENDPOINT"
)

type Config struct {
	Base        // no tag needed here
	GeoDB GeoDB `toml:"geo_db" json:"geo_db"`
}

type GeoDB struct {
	CountryDBPath string `toml:"country_db_path" json:"country_db_path" validate:"required"`
	CityDBPath    string `toml:"city_db_path" json:"city_db_path" validate:"required"`
}

type Base struct {
	HttpAddr string `toml:"http_addr" yaml:"http_addr" validate:"required"`
	GrpcAddr string `toml:"grpc_addr" yaml:"grpc_addr" validate:"required"`

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

func (b *Config) InitDefault() {
	if b.HttpAddr == "" {
		b.HttpAddr = ":8080"
	}
	if b.GrpcAddr == "" {
		b.GrpcAddr = ":8090"
	}
}

func (b *Config) InitOtlpGrpcEndpointFromEnv() {
	slog.Info("init otlp_grpc_endpoint from env")
	if b.OtlpGrpcEndpoint != "" {
		return
	}
	if tmp := os.Getenv(EnvOtlpGrpcEndpoint); tmp != "" {
		b.OtlpGrpcEndpoint = tmp
		slog.Info("otlp_grpc_endpoint is empty, get from env var", "otlp_grpc_endpoint", tmp)
	} else {
		slog.Warn("otlp_grpc_endpoint is empty")
	}
}

const (
	DefaultLogLevel    = "info"
	DefaultLogEncoding = LogEncodingJSON
	DefaultLogOutput   = "stderr"
)

func (b *Config) InitLogger(serviceName, version string) {
	slog.Info("init logger")
	if b.Log.Level == "" {
		b.Log.Level = DefaultLogLevel
	}

	if b.Log.Encoding == "" {
		b.Log.Encoding = DefaultLogEncoding
	}

	if b.Log.Output == "" {
		b.Log.Output = DefaultLogOutput
	}

	logger := slogx.New(slogx.WithTracing(), slogx.WithLevel(b.Log.Level), slogx.WithOutput(b.Log.Output))
	slog.SetDefault(logger)
}

func (b *Config) Validate() error {
	validate := v10validator.New()
	err := validate.Struct(b)
	if err != nil {
		var validationErrors v10validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			return validationErrors
		}
		return err
	}
	return nil
}
