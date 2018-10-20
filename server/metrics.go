package server

import (
	"github.com/prometheus/client_golang/prometheus"
)

// Metrics stores a set of metrics recorders for monitoring
type Metrics struct {
	Archives prometheus.Histogram
	Removes  prometheus.Histogram
	Updates  prometheus.Histogram
}

//nolint:lll
func newMetrics() (metrics Metrics) {
	metrics = Metrics{
		Archives: prometheus.NewHistogram(prometheus.HistogramOpts{
			Namespace: "samplist",
			Subsystem: "scraper",
			Name:      "archives",
			Help:      "Archived servers",
		}),
		Removes: prometheus.NewHistogram(prometheus.HistogramOpts{
			Namespace: "samplist",
			Subsystem: "scraper",
			Name:      "removes",
			Help:      "Removed servers",
		}),
		Updates: prometheus.NewHistogram(prometheus.HistogramOpts{
			Namespace: "samplist",
			Subsystem: "scraper",
			Name:      "updates",
			Help:      "Successfully updated servers",
		}),
	}
	prometheus.MustRegister(
		metrics.Archives,
		metrics.Removes,
		metrics.Updates,
	)
	return
}
