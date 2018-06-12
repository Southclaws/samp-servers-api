package types

// Statistics represents a set of simple metrics for the entire listing database
type Statistics struct {
	Servers          int     `json:"servers"`
	Players          int     `json:"players"`
	PlayersPerServer float32 `json:"players_per_server"`
}
