package types

import (
	"github.com/pkg/errors"
)

// Server contains all the information associated with a game server including the core information, the standard SA:MP
// "rules" and "players" lists as well as any additional fields to enhance the server browsing experience.
type Server struct {
	Core        ServerCore        `json:"core"`
	Rules       map[string]string `json:"ru,omitempty"`
	Description string            `json:"description"`
	Banner      string            `json:"banner"`
	Active      bool              `json:"active"`
}

// Validate checks the contents of a Server object to ensure all the required fields are valid.
func (server *Server) Validate() (errs []error) {
	_, addrErrs := AddressFromString(string(server.Core.Address))
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
