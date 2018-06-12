package storage

import (
	"github.com/Southclaws/samp-servers-api/types"
	"github.com/pkg/errors"
	"gopkg.in/mgo.v2/bson"
)

// GetStatistics returns the current statistics for the server database
// todo: cache this data
func (mgr *Manager) GetStatistics() (statistics types.Statistics, err error) {
	statistics.Servers, err = mgr.collection.Find(bson.M{"active": true}).Count()
	if err != nil {
		err = errors.Wrap(err, "failed to execute find query on database")
		return
	}

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
	statistics.Players = tmp["players"].(int)

	if statistics.Servers > 0 {
		statistics.PlayersPerServer = float32(statistics.Players / statistics.Servers)
	}

	return
}
