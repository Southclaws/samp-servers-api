package scraper

import "github.com/prometheus/client_golang/prometheus"

// metrics stores rates and guages for monitoring
type metrics struct {
	Errors    prometheus.Counter
	Queries   prometheus.Counter
	Successes prometheus.Counter
	Failures  prometheus.Counter
	Archives  prometheus.Counter
	Removals  prometheus.Counter
	QueryTime prometheus.Summary
}

// newMetricsRecorder initialises a new metrics recorder
func newMetricsRecorder() (m *metrics) {
	m = &metrics{
		Errors: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: "samplist",
			Subsystem: "scraper",
			Name:      "errors",
			Help:      "Total errors",
		}),
		Queries: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: "samplist",
			Subsystem: "scraper",
			Name:      "queries",
			Help:      "Total queries",
		}),
		Successes: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: "samplist",
			Subsystem: "scraper",
			Name:      "successes",
			Help:      "Successfully updated servers",
		}),
		Failures: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: "samplist",
			Subsystem: "scraper",
			Name:      "failures",
			Help:      "Failed queries",
		}),
		Archives: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: "samplist",
			Subsystem: "scraper",
			Name:      "archives",
			Help:      "Archived servers",
		}),
		Removals: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: "samplist",
			Subsystem: "scraper",
			Name:      "removes",
			Help:      "Removed servers",
		}),
		QueryTime: prometheus.NewSummary(prometheus.SummaryOpts{
			Namespace: "samplist",
			Subsystem: "scraper",
			Name:      "query_time",
			Help:      "The length of queries in seconds",
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
