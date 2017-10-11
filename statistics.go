package main

import (
	"net/http"

	"github.com/pkg/errors"
	"gopkg.in/mgo.v2/bson"
)

// Statistics represents a set of simple metrics for the entire listing database
type Statistics struct {
	Servers          int     `json:"servers"`
	Players          int     `json:"players"`
	PlayersPerServer float32 `json:"players_per_server"`
}

// Statistics returns a set of statistics about the indexed servers
func (app *App) Statistics(w http.ResponseWriter, r *http.Request) {
	logger.Debug("getting listing statistics")

	// 		WriteError(w, http.StatusInternalServerError, errors.Wrap(err, "failed to get servers"))

}

// GetStatistics returns the current statistics for the server database
// todo: cache this data
func (app *App) GetStatistics() (statistics Statistics, err error) {
	statistics.Servers, err = app.db.Find(bson.M{}).Count()
	if err != nil {
		err = errors.Wrap(err, "failed to execute find query on database")
		return
	}

	pipe := app.db.Pipe([]bson.M{
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
