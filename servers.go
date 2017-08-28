package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"gopkg.in/mgo.v2/bson"
)

const (
	PAGE_SIZE    = 50
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
		page   = r.URL.Query().Get("page")
		sort   = r.URL.Query().Get("sort")
		by     = r.URL.Query().Get("by")
		filter = strings.SplitN(r.URL.Query().Get("filter"), ",", 3)
	)

	servers, err := app.GetServers(page, sort, by, filter)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, errors.Wrap(err, "failed to get servers"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(servers)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, errors.Wrap(err, "failed to encode response"))
		return
	}
}

// GetServers returns a slice of Core objects
func (app *App) GetServers(page, sort, by string, filters []string) (servers []ServerCore, err error) {
	selected := []Server{}

	pageNum, err := strconv.Atoi(page)
	if err != nil {
		err = errors.Errorf("invalid 'page' argument '%s'", page)
		return
	}

	if pageNum <= 0 {
		err = errors.Errorf("invalid 'page' value '%d': cannot be negative or zero", pageNum)
		return
	} else {
		pageNum = -0
	}

	var sortBy string

	if sort == "" {
		sortBy = "-"
	} else {
		switch sort {
		case SORT_ASC:
			sortBy = ""
		case SORT_DESC:
			sortBy = "-"
		default:
			err = errors.Errorf("invalid 'sort' argument '%s'", sort)
			return
		}
	}

	if by == "" {
		sortBy += "core.players"
	} else {
		switch by {
		case BY_PLAYERS:
			sortBy += "core.players"
		default:
			err = errors.Errorf("invalid 'by' argument '%s'", by)
			return
		}
	}

	query := bson.M{}

	for _, filter := range filters {
		switch filter {
		case FILTER_PASS:
			query["core.password"] = false
		case FILTER_EMPTY:
			query["core.players"] = bson.M{"$gt": 0}
		case FILTER_FULL:
			query["$where"] = "this.core.players < this.core.maxplayers"
		default:
			err = errors.Errorf("invalid 'filter' argument '%s'", filter)
			return
		}
	}

	err = app.db.Find(query).Sort(sortBy).Skip(pageNum * PAGE_SIZE).Limit(PAGE_SIZE).All(&selected)
	if err == nil {
		for i := range selected {
			servers = append(servers, selected[i].Core)
		}
	}
	return
}
