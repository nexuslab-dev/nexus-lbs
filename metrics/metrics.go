package metrics

import (
	"github.com/prometheus/common/version"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

func MustRegisterVersionCollector(serviceName string) {
	promSubsystem := strings.Replace(serviceName, "-", "_", -1)
	prometheus.MustRegister(version.NewCollector(promSubsystem))
}
