package v2

import (
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"

	"github.com/Southclaws/samp-servers-api/types"
)

// serverStats returns a set of statistics about the indexed servers
func (v *V2) serverStats(w http.ResponseWriter, r *http.Request) {
	var (
		stats types.Statistics
		err   error
	)

	stats.Servers, err = v.Storage.GetActiveServers()
	if err != nil {
		WriteError(w, http.StatusInternalServerError, errors.Wrap(err, "failed to get servers"))
	}
	stats.Players, err = v.Storage.GetTotalPlayers()
	if err != nil {
		WriteError(w, http.StatusInternalServerError, errors.Wrap(err, "failed to get servers"))
	}

	if stats.Servers > 0 {
		stats.PlayersPerServer = float32(stats.Players / stats.Servers)
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(stats)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, errors.Wrap(err, "failed to encode response"))
		return
	}
}
