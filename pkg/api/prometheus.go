package api

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
	"strconv"
)

// Metrics
var (
	summary = prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Name:      "pluto_scan_total",
		Help:      "Pluto scan details.",
	},
		[]string{
			"name",
			"deprecated",
			"removed",
		})
	count = prometheus.NewGauge(prometheus.GaugeOpts{
		Name:      "pluto_reports_count",
		Help:      "Pluto report.",
	})
)

type Prometheus struct {
	address *string
	instance *Instance
}

func (p Prometheus) marshal() error {
	count.Inc()

	for _, output := range p.instance.Outputs {
		summary.WithLabelValues(output.Name, strconv.FormatBool(output.Deprecated), strconv.FormatBool(output.Removed))
	}

	pusher := newPusher(p.address)
	return pusher.Push()
}

func newPusher(address *string) *push.Pusher {
	registry := prometheus.NewRegistry()
	registry.MustRegister(summary, count)
	return push.New(*address, "pluto").Gatherer(registry)
}
