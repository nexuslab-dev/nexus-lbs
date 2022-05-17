package config

import (
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/prometheus/common/version"
)

type SentryConfig struct {
	Dsn          string `toml:"dsn" yaml:"dsn" json:"dsn,omitempty"`
	Env          string `toml:"env" yaml:"env" json:"env,omitempty"`
	FlushWaitSec int    `toml:"flush_wait_sec" yaml:"flush_wait_sec" json:"flush_wait_sec,omitempty"`
}

type SentryFlusher func() bool

func (sc *SentryConfig) SentryInit(serviceName string) (SentryFlusher, error) {
	err := sentry.Init(sentry.ClientOptions{
		Dsn:         sc.Dsn,
		Environment: sc.Env,
		Release:     version.Version,
		ServerName:  serviceName,
	})
	if err != nil {
		// return empty flusher and the error
		return func() bool {
			return false
		}, err
	}

	return func() bool {
		// Since sentry emits events in the background we need to make sure
		// they are sent before we shut down
		waitSec := sc.FlushWaitSec
		if waitSec == 0 {
			waitSec = 3
		}
		return sentry.Flush(time.Second * time.Duration(waitSec))
	}, err
}
