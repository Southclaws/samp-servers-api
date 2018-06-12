package server

import (
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
	"gopkg.in/mgo.v2/bson"

	"github.com/Southclaws/samp-servers-api/types"
)

// serverStats returns a set of statistics about the indexed servers
func (app *App) serverStats(w http.ResponseWriter, r *http.Request) {
	logger.Debug("getting listing statistics")

	stats, err := app.GetStatistics()
	if err != nil {
		WriteError(w, http.StatusInternalServerError, errors.Wrap(err, "failed to get servers"))
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(stats)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, errors.Wrap(err, "failed to encode response"))
		return
	}
}

// GetStatistics returns the current statistics for the server database
// todo: cache this data
func (app *App) GetStatistics() (statistics types.Statistics, err error) {
	statistics.Servers, err = app.collection.Find(bson.M{"active": true}).Count()
	if err != nil {
		err = errors.Wrap(err, "failed to execute find query on database")
		return
	}

	pipe := app.collection.Pipe([]bson.M{
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
