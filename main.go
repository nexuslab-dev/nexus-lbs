package main

import (
	"context"
	"errors"
	"math/rand"
	"net/http"
	"time"

	"github.com/nexuslab-dev/nexus-lbs/config"
	"github.com/nexuslab-dev/nexus-lbs/core"
	"github.com/nexuslab-dev/nexus-lbs/httpapi"
	"github.com/nexuslab-dev/nexus-lbs/metrics"
	"github.com/prometheus/common/version"
	"github.com/ttys3/lgr"
	"github.com/ttys3/tracing"
)

var ServiceName = "undefined"

var cfg config.Config

// @title Nexus LBS Api
// @version 1.0
// @description This is a simple lbs server, currently only query location by IP feature is implemented
// @termsOfService https://github.com/nexuslab-dev/nexus-lbs

// @BasePath /v1/api/
func main() {
	rand.Seed(time.Now().UnixNano())
	metrics.MustRegisterVersionCollector(ServiceName)

	if err := config.NewLoader(ServiceName).Load(&cfg); err != nil {
		lgr.S().Fatal("load config failed", "err", err)
	}
	cfg.InitDefault()
	cfg.InitOtlpGrpcEndpointFromEnv()
	cfg.InitLogger(ServiceName, version.Version)

	if err := cfg.Validate(); err != nil {
		lgr.S().Fatal("config validation failed", "err", err)
	}

	lgr.S().Debug("load config success", "config", cfg)

	// init sentry
	if sentryFlusher, err := cfg.Sentry.SentryInit(ServiceName); err != nil {
		// we can skip this error, we'll start the service even sentry is down
		lgr.S().Error("init sentry failed", "err", err)
	} else {
		lgr.S().Info("init sentry succeeded", "sentry", cfg.Sentry)
		defer func() {
			lgr.S().Info("flushing sentry ...")
			sentryFlusher()
		}()
	}

	// init tracing
	if cfg.Tracing && cfg.OtlpGrpcEndpoint != "" {
		if tpShutdown, err := initTracing(); err != nil {
			lgr.S().Fatal("tracing init failed", "err", err)
		} else {
			// nolint: forbidigo
			defer tpShutdown(context.Background())
			lgr.S().Info("tracing init success", "otlp_grpc_endpoint", cfg.OtlpGrpcEndpoint)
		}
	} else {
		lgr.S().Warn("tracing is disabled by config")
	}

	// open db
	countryQuery, err := core.New(cfg.GeoDB.CountryDBPath)
	if err != nil {
		lgr.S().Fatal("new country query failed", "err", err, "db", cfg.GeoDB.CountryDBPath)
	}
	cityQuery, err := core.New(cfg.GeoDB.CityDBPath)
	if err != nil {
		lgr.S().Fatal("new city query failed", "err", err, "db", cfg.GeoDB.CityDBPath)
	}

	httpserver := httpapi.NewServer(ServiceName, cfg.Profile, cfg.Tracing)
	httpapi.Register(httpserver, ServiceName, countryQuery, cityQuery)

	lgr.S().Info("http server started", "http_addr", cfg.HttpAddr)
	if err := httpserver.Start(cfg.HttpAddr); err != nil && !errors.Is(err, http.ErrServerClosed) {
		lgr.S().Error("http server exited with error", "err", err)
	} else {
		lgr.S().Info("http server exited")
	}
}

func initTracing() (tracing.TpShutdownFunc, error) {
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
	return tpShutdown, err
}
