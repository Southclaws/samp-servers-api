package types

import (
	"github.com/pkg/errors"
)

// Server contains all the information associated with a game server including the core information, the standard SA:MP
// "rules" and "players" lists as well as any additional fields to enhance the server browsing experience.
type Server struct {
	IP          string            `json:"ip"`
	Core        ServerCore        `json:"core"`
	Rules       map[string]string `json:"ru,omitempty"`
	Description string            `json:"description"`
	Banner      string            `json:"banner"`
	Active      bool              `json:"active"`
}

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

// Validate checks the contents of a Server object to ensure all the required fields are valid.
func (server *Server) Validate() (errs []error) {
	_, addrErrs := AddressFromString(server.Core.Address)
	errs = append(errs, addrErrs...)

	if len(server.Core.Hostname) < 1 {
		errs = append(errs, errors.New("hostname is empty"))
	}

	if server.Core.MaxPlayers == 0 {
		errs = append(errs, errors.New("maxplayers is empty"))
	}

	if len(server.Core.Gamemode) < 1 {
		errs = append(errs, errors.New("gamemode is empty"))
	}

	return
}

// Example returns an example of Server
func (server Server) Example() Server {
	return Server{
		Core: ServerCore{
			Address:    "127.0.0.1:7777",
			Hostname:   "SA-MP SERVER CLAN tdm [NGRP] [GF EDIT] [Y_INI] [RUS] [BASIC] [GODFATHER] [REFUNDING] [STRCMP]",
			Players:    32,
			MaxPlayers: 128,
			Gamemode:   "Grand Larceny",
			Language:   "English",
			Password:   false,
			Version:    "0.3.7-R2",
		},
		Rules: map[string]string{
			"lagcomp":   "On",
			"mapname":   "San Andreas",
			"version":   "0.3.7-R2",
			"weather":   "10",
			"weburl":    "www.sa-mp.com",
			"worldtime": "10:00",
		},
		Description: "An awesome server! Come and play with us.",
		Banner:      "https://i.imgur.com/Juaezhv.jpg",
		Active:      true,
	}
}
