package types

// Statistics represents a set of simple metrics for the entire listing database
type Statistics struct {
	Servers          int     `json:"servers"`
	Players          int     `json:"players"`
	PlayersPerServer float32 `json:"players_per_server"`
}

// Example returns an example of Statistics
func (s Statistics) Example() Statistics {
	return Statistics{
		Servers:          1000,
		Players:          10000,
		PlayersPerServer: 10,
	}
}
