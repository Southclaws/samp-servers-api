package storage

import (
	"fmt"

	"github.com/pkg/errors"
	"gopkg.in/mgo.v2/bson"

	"github.com/Southclaws/samp-servers-api/types"
)

// GetServers returns a slice of Core objects
func (mgr *Manager) GetServers(pageNum int, sort types.SortOrder, by types.SortColumn, filters []types.FilterAttribute) (servers []types.ServerCore, err error) {
	selected := []types.Server{}

	if pageNum <= 0 {
		pageNum = 0
	} else {
		pageNum = pageNum - 1 // subtract 1 so 1 becomes 0, "page 1" makes more sense to users
	}

	var sortBy types.SortOrder

	if sort == "" {
		sortBy = "-"
	} else {
		switch sort {
		case types.SortAsc:
			sortBy = ""
		case types.SortDesc:
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
		case types.ByPlayers:
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
			case types.FilterPassword:
				query["core.password"] = false
			case types.FilterEmpty:
				query["core.players"] = bson.M{"$gt": 0}
			case types.FilterFull:
				query["$where"] = "this.core.players < this.core.maxplayers"
			}
		}
	}

	fmt.Println(query, sortBy, pageNum)

	err = mgr.collection.
		Find(query).
		Sort(string(sortBy)).
		Skip(pageNum * int(types.PageSizeDefault)).
		Limit(int(types.PageSizeDefault)).
		All(&selected)
	if err == nil {
		for i := range selected {
			servers = append(servers, selected[i].Core)
		}
	}
	return
}
