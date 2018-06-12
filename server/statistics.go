package server

import (
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
)

// serverStats returns a set of statistics about the indexed servers
func (app *App) serverStats(w http.ResponseWriter, r *http.Request) {
	logger.Debug("getting listing statistics")

	stats, err := app.db.GetStatistics()
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
