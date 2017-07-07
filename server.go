package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Core stores the standard SA:MP 'info' query fields necessary for server lists. The json keys are short to cut down on
// network traffic since these are the objects returned to a listing request which could contain hundreds of objects.
type Core struct {
	Address    string `json:"ip"`
	Hostname   string `json:"hn"`
	Players    int    `json:"pc"`
	MaxPlayers int    `json:"pm"`
	Gamemode   string `json:"gm"`
	Language   string `json:"la"`
	Password   bool   `json:"pa"`
}

// Server contains all the information associated with a game server including the core information, the standard SA:MP
// "rules" and "players" lists as well as any additional fields to enhance the server browsing experience.
type Server struct {
	Core        Core              `json:"core"`
	Rules       map[string]string `json:"ru,omitempty"`
	PlayerList  []string          `json:"pl,omitempty"`
	Description string            `json:"description"`
	Banner      string            `json:"banner"`
}

// Validate checks the contents of a Server object to ensure all the required fields are valid.
func (server *Server) Validate() (errs []error) {
	errs = append(errs, ValidateAddress(server.Address)...)

	if len(server.Hostname) < 1 {
		errs = append(errs, fmt.Errorf("hostname is empty"))
	}

	if server.MaxPlayers == 0 {
		errs = append(errs, fmt.Errorf("maxplayers is empty"))
	}

	if len(server.Gamemode) < 1 {
		errs = append(errs, fmt.Errorf("gamemode is empty"))
	}

	return
}

// ValidateAddress validates an address field for a server and ensures it contains the correct
// combination of host:port with either "samp://" or an empty scheme.
func ValidateAddress(address string) (errs []error) {
	if len(address) < 1 {
		errs = append(errs, fmt.Errorf("address is empty"))
	}

	if !strings.Contains(address, "://") {
		address = fmt.Sprintf("samp://%s", address)
	}

	u, err := url.Parse(address)
	if err != nil {
		errs = append(errs, err)
		return
	}

	if u.User != nil {
		errs = append(errs, fmt.Errorf("address contains a user:password component"))
	}

	if u.Scheme != "samp" && u.Scheme != "" {
		errs = append(errs, fmt.Errorf("address contains invalid scheme '%s', must be either empty or 'samp://'", u.Scheme))
	}

	portStr := u.Port()

	if portStr != "" {
		port, err := strconv.Atoi(u.Port())
		if err != nil {
			errs = append(errs, fmt.Errorf("invalid port '%s' specified", u.Port()))
			return
		}

		if port < 1024 || port > 49152 {
			errs = append(errs, fmt.Errorf("port %d falls within reserved or ephemeral range", port))
		}
	}

	return
}

// ServerSimple handles "simple" posts where the only data is the server address which is passed to
// the QueryDaemon which handles pulling the rest of the information from the legacy query API.
func (app *App) ServerSimple(w http.ResponseWriter, r *http.Request) {
	raw, err := ioutil.ReadAll(r.Body)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err)
		return
	}

	address := string(raw)

	errs := ValidateAddress(address)
	if errs != nil {
		WriteErrors(w, http.StatusBadRequest, errs)
		return
	}
}

// Server handles either posting a server object or requesting a server object
func (app *App) Server(w http.ResponseWriter, r *http.Request) {
	address, ok := mux.Vars(r)["address"]
	if !ok {
		WriteError(w, http.StatusBadRequest, fmt.Errorf("no address specified"))
	}

	switch r.Method {
	case "GET":
		logger.Debug("getting server",
			zap.String("address", address))

		var (
			err error
		)

		errs := ValidateAddress(address)
		if errs != nil {
			WriteErrors(w, http.StatusBadRequest, errs)
			return
		}

		server := Server{}

		found, err := app.GetServer(address, &server)
		if err != nil {
			WriteError(w, http.StatusInternalServerError, err)
			return
		}

		if !found {
			WriteError(w, http.StatusNotFound, fmt.Errorf("could not find server by address '%s'", address))
			return
		}

		err = json.NewEncoder(w).Encode(&server)
		if err != nil {
			WriteError(w, http.StatusInternalServerError, err)
			return
		}

	case "POST":
		logger.Debug("posting server",
			zap.String("address", address))

		server := Server{}
		err := json.NewDecoder(r.Body).Decode(&server)
		if err != nil {
			WriteError(w, http.StatusBadRequest, err)
			return
		}

		errs := server.Validate()
		if errs != nil {
			WriteErrors(w, http.StatusUnprocessableEntity, errs)
		}

		err = app.UpsertServer(server)
		if err != nil {
			WriteError(w, http.StatusInternalServerError, err)
		}
	}
}

// GetServer looks up a server via the address
func (app *App) GetServer(address string, server *Server) (found bool, err error) {
	err = app.db.Find(bson.M{"address": address}).One(server)
	if err == mgo.ErrNotFound {
		found = false
		err = nil // the caller does not need to interpret this as an "error"
	} else if err != nil {
		return
	} else {
		found = true
	}

	return
}

// UpsertServer creates or updates a server object in the database.
func (app *App) UpsertServer(server Server) (err error) {
	info, err := app.db.Upsert(bson.M{"address": server.Address}, server)
	logger.Debug("upsert server",
		zap.Int("matched", info.Matched),
		zap.Int("removed", info.Removed),
		zap.Int("updated", info.Updated),
		zap.Any("id", info.UpsertedId))
	return
}
