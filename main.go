package main

import (
	"context"
	"github.com/nexuslab-dev/nexus-lbs/config"
	"github.com/nexuslab-dev/nexus-lbs/metrics"
	"github.com/prometheus/common/version"
	"github.com/ttys3/lgr"
	"github.com/ttys3/tracing"
	"math/rand"
	"time"
)

var ServiceName string = "undefined"

var cfg config.Config

func main() {
	rand.Seed(time.Now().UnixNano())
	metrics.MustRegisterVersionCollector(ServiceName)

	if err := config.NewLoader(ServiceName).Load(&cfg); err != nil {
		lgr.S().Fatal("load config failed", "err", err)
	}
	lgr.S().Debug("load config success", "config", cfg)

	// init sentry
	sentryFlusher, err := cfg.Sentry.SentryInit(ServiceName)
	defer func() {
		lgr.S().Info("flushing sentry ...")
		sentryFlusher()
	}()

	if err != nil {
		// we can skip this error, we'll start the service even sentry is down
		lgr.S().Error("init sentry failed", "err", err)
	} else {
		lgr.S().Info("init sentry succeeded", "sentry", cfg.Sentry)
	}

	// init tracing
	if cfg.Tracing {
		if cfg.OtlpGrpcEndpoint != "" {
			attrs := []string{
				"svc_branch", version.Branch,
				"svc_revision", version.Revision,
				"svc_build_date", version.BuildDate,
				"svc_go_version", version.GoVersion,
			}

			// nolint: forbidigo
			tpShutdown, err := tracing.InitOtlpTracerProvider(context.Background(),
				tracing.WithOtelGrpcEndpoint(cfg.OtlpGrpcEndpoint),
				tracing.WithSerivceName(ServiceName),
				tracing.WithServiceVersion(version.Version),
				tracing.WithAttributes(attrs),
			)
			// nolint: forbidigo
			defer tpShutdown(context.Background())

			if err != nil {
				lgr.S().Fatal("tracing init failed", "err", err)
			} else {
				lgr.S().Info("tracing init success", "otlp_grpc_endpoint", cfg.OtlpGrpcEndpoint)
			}
		} else {
			lgr.S().Error("OtlpGrpcEndpoint is empty, skip init tracing")
		}
	} else {
		lgr.S().Warn("tracing is disabled by config")
	}

}
