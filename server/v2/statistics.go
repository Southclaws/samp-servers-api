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

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(stats)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, errors.Wrap(err, "failed to encode response"))
		return
	}
}
