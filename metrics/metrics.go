package metrics

import (
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors/version"
)

func MustRegisterVersionCollector(serviceName string) {
	promSubsystem := strings.Replace(serviceName, "-", "_", -1)
	prometheus.MustRegister(version.NewCollector(promSubsystem))
}
