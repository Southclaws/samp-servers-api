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
	// PageSize controls the default page size of listings
	PageSize = 50
	// SortAsc is the ascending sort order for listings
	SortAsc = "asc"
	// SortDesc is the descending sort order for listings
	SortDesc = "desc"
	// ByPlayers means the list will use the amount of players as a sort key
	ByPlayers = "player"
	// FilterPass filters out servers with passwords
	FilterPass = "password"
	// FilterEmpty filters out empty servers
	FilterEmpty = "empty"
	// FilterFull filters out full servers
	FilterFull = "full"
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
		pageNum = 1
	}

	if pageNum <= 0 {
		err = errors.Errorf("invalid 'page' value '%d': cannot be negative or zero", pageNum)
		return
	}
	pageNum = 1

	var sortBy string

	if sort == "" {
		sortBy = "-"
	} else {
		switch sort {
		case SortAsc:
			sortBy = ""
		case SortDesc:
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
		case ByPlayers:
			sortBy += "core.players"
		default:
			err = errors.Errorf("invalid 'by' argument '%s'", by)
			return
		}
	}

	query := bson.M{"active": true}

	if len(filters) > 0 {
		for _, filter := range filters {
			switch filter {
			case FilterPass:
				query["core.password"] = false
			case FilterEmpty:
				query["core.players"] = bson.M{"$gt": 0}
			case FilterFull:
				query["$where"] = "this.core.players < this.core.maxplayers"
			}
		}
	}

	err = app.db.Find(query).Sort(sortBy).Skip(pageNum * PageSize).Limit(PageSize).All(&selected)
	if err == nil {
		for i := range selected {
			servers = append(servers, selected[i].Core)
		}
	}
	return
}
