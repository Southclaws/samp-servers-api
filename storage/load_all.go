package storage

import (
	"github.com/pkg/errors"
	"gopkg.in/mgo.v2/bson"

	"github.com/Southclaws/samp-servers-api/types"
)

// LoadAllAddresses loads all addresses from the database as a slice of strings for synchronisation
// with the QueryDaemon.
func (mgr *Manager) LoadAllAddresses() (result []string, err error) {
	allServers := []types.Server{}
	err = mgr.collection.Find(bson.M{}).All(&allServers)
	if err != nil {
		err = errors.Wrap(err, "failed to load current addresses for query daemon")
		return
	}
	for i := range allServers {
		result = append(result, allServers[i].Core.Address)
	}
	return
}
