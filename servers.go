package main

import (
	"encoding/json"
	"net/http"

	"gopkg.in/mgo.v2/bson"
)

// Servers returns a JSON encoded array of available servers
func (app *App) Servers(w http.ResponseWriter, r *http.Request) {
	logger.Debug("getting server list")

	servers := []Core{}
	err := app.GetServers(&servers)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err)
		return
	}

	err = json.NewEncoder(w).Encode(servers)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err)
		return
	}
}

// GetServers returns a slice of Core objects
func (app *App) GetServers(servers *[]Core) (err error) {
	err = app.db.Find(bson.M{}).All(&servers)
	return
}
