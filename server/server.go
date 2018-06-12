package server

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/Southclaws/samp-servers-api/types"
)

// serverAdd handles "simple" posts where the only data is the server address which is passed to
// the QueryDaemon which handles pulling the rest of the information from the legacy query API.
func (app *App) serverAdd(w http.ResponseWriter, r *http.Request) {
	raw, err := ioutil.ReadAll(r.Body)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err)
		return
	}

	address := string(raw)

	normalised, errs := types.ValidateAddress(address)
	if errs != nil {
		WriteErrors(w, http.StatusBadRequest, errs)
		return
	}

	app.qd.Add(normalised)
}

// serverPost handles posting a server object
func (app *App) serverPost(w http.ResponseWriter, r *http.Request) {
	address, ok := mux.Vars(r)["address"]
	if !ok {
		WriteError(w, http.StatusBadRequest, errors.New("no address specified"))
	}

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

	server := types.Server{}
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

// serverGet handles responding to a request by server address
func (app *App) serverGet(w http.ResponseWriter, r *http.Request) {
	address, ok := mux.Vars(r)["address"]
	if !ok {
		WriteError(w, http.StatusBadRequest, errors.New("no address specified"))
	}

	logger.Debug("getting server",
		zap.String("address", address))

	var (
		err error
	)

	_, errs := types.ValidateAddress(address)
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
}

// GetServer looks up a server via the address
func (app *App) GetServer(address string) (server types.Server, found bool, err error) {
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
func (app *App) UpsertServer(server types.Server) (err error) {
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
