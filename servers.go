package main

import (
	"encoding/json"
	"net/http"

	"gopkg.in/mgo.v2/bson"
)

const (
	SORT_ASC     = "asc"
	SORT_DESC    = "desc"
	BY_PLAYERS   = "player"
	FILTER_PASS  = "password"
	FILTER_EMPTY = "empty"
	FILTER_FULL  = "full"
)

// Servers returns a JSON encoded array of available servers
func (app *App) Servers(w http.ResponseWriter, r *http.Request) {
	logger.Debug("getting server list")

	var (
		err    error
		sort   = r.URL.Query().Get("sort")
		by     = r.URL.Query().Get("by")
		filter = r.URL.Query().Get("filter")
	)

	servers, err := app.GetServers(sort, by, filter)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(servers)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err)
		return
	}
}

// GetServers returns a slice of Core objects
func (app *App) GetServers(offset, limit int, sort, by, filter string) (servers []ServerCore, err error) {
	selected := []Server{}
	query := bson.M{}

	switch filter {
	case FILTER_PASS:
		query["core.pa"] = false
	case FILTER_EMPTY:
		query["core.pc"] = bson.M{"$gt": 0}
	case FILTER_FULL:
		query["$where"] = "this.core.pc < this.core.pm"
	}

	switch sort {
	case BY_PLAYERS:
		sort = "core.pc"
	default:
		sort = ""
	}

	err = app.db.Find(query).Sort(sort).Skip(offset).Limit(limit).All(selected)
	if err == nil {
		for i := range selected {
			servers = append(servers, selected[i].Core)
		}
	}
	return
}
