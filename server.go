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
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

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

// ServerSimple handles "simple" posts where the only data is the server address which is passed to
// the QueryDaemon which handles pulling the rest of the information from the legacy query API.
func (app *App) ServerSimple(w http.ResponseWriter, r *http.Request) {
	raw, err := ioutil.ReadAll(r.Body)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err)
		return
	}

	address := string(raw)

	normalised, errs := ValidateAddress(address)
	if errs != nil {
		WriteErrors(w, http.StatusBadRequest, errs)
		return
	}

	app.qd.Add(normalised)
}

// Server handles either posting a server object or requesting a server object
func (app *App) Server(w http.ResponseWriter, r *http.Request) {
	address, ok := mux.Vars(r)["address"]
	if !ok {
		WriteError(w, http.StatusBadRequest, errors.New("no address specified"))
	}

	switch r.Method {
	case "GET":
		logger.Debug("getting server",
			zap.String("address", address))

		var (
			err error
		)

		_, errs := ValidateAddress(address)
		if errs != nil {
			WriteErrors(w, http.StatusBadRequest, errs)
			return
		}

		server, found, err := app.GetServer(address)
		if err != nil {
			WriteError(w, http.StatusInternalServerError, err)
			return
		}

		if !found {
			WriteError(w, http.StatusNotFound, errors.Errorf("could not find server by address '%s'", address))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(&server)
		if err != nil {
			WriteError(w, http.StatusInternalServerError, err)
			return
		}

	case "POST":
		from := strings.Split(r.RemoteAddr, ":")[0]

		if app.config.VerifyByHost {
			addressIP := strings.Split(address, ":")[0]
			if from != addressIP {
				WriteError(w, http.StatusBadRequest, errors.Errorf("request address '%v' does not match declared server address '%s'", from, addressIP))
				return
			}
		}

		logger.Debug("posting server",
			zap.String("address", address),
			zap.String("from", from))

		server := Server{}
		err := json.NewDecoder(r.Body).Decode(&server)
		if err != nil {
			WriteError(w, http.StatusBadRequest, err)
			return
		}

		if server.Core.Address != address {
			WriteError(w, http.StatusBadRequest, errors.Errorf("route address '%v' does not match payload address '%s'", address, server.Core.Address))
			return
		}

		errs := server.Validate()
		if errs != nil {
			WriteErrors(w, http.StatusUnprocessableEntity, errs)
			return
		}

		err = app.UpsertServer(server)
		if err != nil {
			WriteError(w, http.StatusInternalServerError, err)
		}
	}
}

// GetServer looks up a server via the address
func (app *App) GetServer(address string) (server Server, found bool, err error) {
	err = app.collection.Find(bson.M{"core.address": address, "active": true}).One(&server)
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

// UpsertServer creates or updates a server object in the database, implicitly sets `Active` to true
func (app *App) UpsertServer(server Server) (err error) {
	server.Active = true
	info, err := app.collection.Upsert(bson.M{"core.address": server.Core.Address}, server)
	if err != nil {
		logger.Error("upsert server failed",
			zap.String("address", server.Core.Address))

	} else if info != nil {
		logger.Debug("upsert server",
			zap.String("address", server.Core.Address),
			zap.Int("matched", info.Matched),
			zap.Int("removed", info.Removed),
			zap.Int("updated", info.Updated),
			zap.Any("id", info.UpsertedId))

		app.qd.Add(server.Core.Address)
	}

	return
}

// MarkInactive marks a server as inactive by setting the `Active` field to false
func (app *App) MarkInactive(address string) (err error) {
	return app.collection.Update(bson.M{"core.address": address}, bson.M{"$set": bson.M{"active": false}})
}

// RemoveServer deletes a server from the database
func (app *App) RemoveServer(address string) (err error) {
	return app.collection.Remove(bson.M{"core.address": address})
}
