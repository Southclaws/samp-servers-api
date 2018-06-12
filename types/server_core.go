package types

// ServerCore stores the standard SA:MP 'info' query fields necessary for server lists. The json keys are short to cut down on
// network traffic since these are the objects returned to a listing request which could contain hundreds of objects.
type ServerCore struct {
	Address    string `json:"ip"`
	Hostname   string `json:"hn"`
	Players    int    `json:"pc"`
	MaxPlayers int    `json:"pm"`
	Gamemode   string `json:"gm"`
	Language   string `json:"la"`
	Password   bool   `json:"pa"`
	Version    string `json:"vn"`
}
