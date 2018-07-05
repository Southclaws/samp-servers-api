package v3

import (
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
)

// serverStats returns a set of statistics about the indexed servers
func (v *V3) serverStats(w http.ResponseWriter, r *http.Request) {
	stats, err := v.Storage.GetStatistics()
	if err != nil {
		WriteError(w, http.StatusInternalServerError, errors.Wrap(err, "failed to get servers"))
	}

	stats.Metrics = v.Metrics.GetValues()

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(stats)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, errors.Wrap(err, "failed to encode response"))
		return
	}
}
