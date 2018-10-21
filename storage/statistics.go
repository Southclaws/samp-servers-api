package storage

import (
	"github.com/pkg/errors"
	"gopkg.in/mgo.v2/bson"
)

// GetActiveServers returns the number of active servers
func (mgr *Manager) GetActiveServers() (servers int, err error) {
	servers, err = mgr.collection.Find(bson.M{"active": true}).Count()
	if err != nil {
		err = errors.Wrap(err, "failed to execute find query on database")
		return
	}
	return
}

// GetInactiveServers returns the number of inactive servers
func (mgr *Manager) GetInactiveServers() (servers int, err error) {
	servers, err = mgr.collection.Find(bson.M{"active": false}).Count()
	if err != nil {
		err = errors.Wrap(err, "failed to execute find query on database")
		return
	}
	return
}

// GetTotalPlayers returns the number of total players
func (mgr *Manager) GetTotalPlayers() (players int, err error) {
	pipe := mgr.collection.Pipe([]bson.M{
		bson.M{
			"$group": bson.M{
				"_id": nil,
				"players": bson.M{
					"$sum": "$core.players",
				},
			},
		},
	})

	var tmp map[string]interface{}
	err = pipe.One(&tmp)
	if err != nil {
		err = errors.Wrap(err, "failed to sum core.players")
		return
	}
	players = tmp["players"].(int)
	return
}
