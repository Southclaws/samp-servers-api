package types

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

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
	_, addrErrs := ValidateAddress(server.Core.Address)
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

// ValidateAddress validates an address field for a server and ensures it contains the correct
// combination of host:port with either "samp://" or an empty scheme. returns an address with the
// :7777 port if absent (this is the default SA:MP port) and strips the "samp:// protocol".
func ValidateAddress(address string) (normalised string, errs []error) {
	if len(address) < 1 {
		errs = append(errs, errors.New("address is empty"))
	}

	if !strings.Contains(address, "://") {
		normalised = fmt.Sprintf("samp://%s", address)
	} else {
		normalised = address
	}

	u, err := url.Parse(normalised)
	if err != nil {
		errs = append(errs, err)
		return
	}

	if u.User != nil {
		errs = append(errs, errors.New("address contains a user:password component"))
	}

	if u.Scheme != "samp" {
		errs = append(errs, errors.Errorf("address contains invalid scheme '%s', must be either empty or 'samp://'", u.Scheme))
	}

	portStr := u.Port()

	if portStr != "" {
		port, err := strconv.Atoi(u.Port())
		if err != nil {
			errs = append(errs, errors.Errorf("invalid port '%s' specified", u.Port()))
			return
		}

		if port < 1024 || port > 49152 {
			errs = append(errs, errors.Errorf("port %d falls within reserved or ephemeral range", port))
			return
		}

		normalised = u.Host
	} else {
		normalised = u.Host + ":7777"
	}

	return
}
