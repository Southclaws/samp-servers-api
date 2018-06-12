package storage

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/Southclaws/samp-servers-api/types"
)

// GetServer looks up a server via the address
func (mgr *Manager) GetServer(address string) (server types.Server, found bool, err error) {
	err = mgr.collection.Find(bson.M{"core.address": address, "active": true}).One(&server)
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
func (mgr *Manager) UpsertServer(server types.Server) (err error) {
	server.Active = true
	_, err = mgr.collection.Upsert(bson.M{"core.address": server.Core.Address}, server)
	return
}

// MarkInactive marks a server as inactive by setting the `Active` field to false
func (mgr *Manager) MarkInactive(address string) (err error) {
	return mgr.collection.Update(bson.M{"core.address": address}, bson.M{"$set": bson.M{"active": false}})
}

// RemoveServer deletes a server from the database
func (mgr *Manager) RemoveServer(address string) (err error) {
	return mgr.collection.Remove(bson.M{"core.address": address})
}
