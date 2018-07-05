package types

import (
	"github.com/rcrowley/go-metrics"
)

// Metrics stores rates and guages for monitoring
type Metrics struct {
	Errors    metrics.Meter
	Queries   metrics.Meter
	Successes metrics.Meter
	Failures  metrics.Meter
	Archives  metrics.Meter
	Removals  metrics.Meter
}

// MetricsValues is for encoding the above object
type MetricsValues struct {
	Errors    float32 `json:"errors_per_second"`
	Queries   float32 `json:"queries_per_second"`
	Successes float32 `json:"successes_per_second"`
	Failures  float32 `json:"failures_per_second"`
	Archives  float32 `json:"archives_per_second"`
	Removals  float32 `json:"removals_per_second"`
}

// NewMetricsRecorder initialises a new metrics recorder
func NewMetricsRecorder() (m *Metrics) {
	return &Metrics{
		Errors:    metrics.NewMeter(),
		Queries:   metrics.NewMeter(),
		Successes: metrics.NewMeter(),
		Failures:  metrics.NewMeter(),
		Archives:  metrics.NewMeter(),
		Removals:  metrics.NewMeter(),
	}
}

// GetValues returns all metrics as Rate5 values
func (m Metrics) GetValues() (mv MetricsValues) {
	return MetricsValues{
		Errors:    float32(m.Errors.Rate5()),
		Queries:   float32(m.Queries.Rate5()),
		Successes: float32(m.Successes.Rate5()),
		Failures:  float32(m.Failures.Rate5()),
		Archives:  float32(m.Archives.Rate5()),
		Removals:  float32(m.Removals.Rate5()),
	}
}
