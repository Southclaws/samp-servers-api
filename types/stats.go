package types

// Statistics represents a set of simple metrics for the entire listing database
type Statistics struct {
	Servers          int           `json:"servers"`
	Players          int           `json:"players"`
	PlayersPerServer float32       `json:"players_per_server"`
	Metrics          MetricsValues `json:"metrics"`
}

// Example returns an example of Statistics
func (s Statistics) Example() Statistics {
	return Statistics{
		Servers:          1000,
		Players:          10000,
		PlayersPerServer: 10,
		Metrics: MetricsValues{
			Errors:   3.4,
			Archives: 3.4,
			Removals: 3.4,
		},
	}
}
