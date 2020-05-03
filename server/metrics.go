package server

import "github.com/prometheus/client_golang/prometheus"

// metrics stores rates and guages for monitoring
type metrics struct {
	Active   prometheus.Gauge
	Inactive prometheus.Gauge
	Players  *prometheus.GaugeVec
}

// newMetricsRecorder initialises a new metrics recorder
func newMetricsRecorder() (m *metrics) {
	m = &metrics{
		Active: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "samplist",
			Subsystem: "index",
			Name:      "active",
			Help:      "Total active servers.",
		}),
		Inactive: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "samplist",
			Subsystem: "index",
			Name:      "inactive",
			Help:      "Total servers that are offline but being given a grace-period to come back online.",
		}),
		Players: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "samplist",
			Subsystem: "index",
			Name:      "players",
			Help:      "Total players across all servers",
		}, []string{"addr"}),
	}
	prometheus.MustRegister(
		m.Active,
		m.Inactive,
		m.Players,
	)
	return m
}
