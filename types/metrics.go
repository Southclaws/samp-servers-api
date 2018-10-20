package types

import (
	"github.com/prometheus/client_golang/prometheus"
)

// Metrics stores rates and guages for monitoring
type Metrics struct {
	Errors    prometheus.Histogram
	Queries   prometheus.Histogram
	Successes prometheus.Histogram
	Failures  prometheus.Histogram
	Archives  prometheus.Histogram
	Removals  prometheus.Histogram
}

// NewMetricsRecorder initialises a new metrics recorder
func NewMetricsRecorder() (m *Metrics) {
	m = &Metrics{
		Errors: prometheus.NewHistogram(prometheus.HistogramOpts{
			Namespace: "samplist",
			Subsystem: "scraper",
			Name:      "errors",
			Help:      "Total errors",
		}),
		Queries: prometheus.NewHistogram(prometheus.HistogramOpts{
			Namespace: "samplist",
			Subsystem: "scraper",
			Name:      "queries",
			Help:      "Total queries",
		}),
		Successes: prometheus.NewHistogram(prometheus.HistogramOpts{
			Namespace: "samplist",
			Subsystem: "scraper",
			Name:      "successes",
			Help:      "Successfully updated servers",
		}),
		Failures: prometheus.NewHistogram(prometheus.HistogramOpts{
			Namespace: "samplist",
			Subsystem: "scraper",
			Name:      "failures",
			Help:      "Failed queries",
		}),
		Archives: prometheus.NewHistogram(prometheus.HistogramOpts{
			Namespace: "samplist",
			Subsystem: "scraper",
			Name:      "archives",
			Help:      "Archived servers",
		}),
		Removals: prometheus.NewHistogram(prometheus.HistogramOpts{
			Namespace: "samplist",
			Subsystem: "scraper",
			Name:      "removes",
			Help:      "Removed servers",
		}),
	}
	prometheus.MustRegister(
		m.Errors,
		m.Queries,
		m.Successes,
		m.Failures,
		m.Archives,
		m.Removals,
	)
	return m
}
