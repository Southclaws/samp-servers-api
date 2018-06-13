package v2

import (
	"encoding/json"
	"net/http"

	"github.com/dyninc/qstring"
	"github.com/pkg/errors"

	"github.com/Southclaws/samp-servers-api/types"
)

// Servers returns a JSON encoded array of available servers
func (v *V2) serverList(w http.ResponseWriter, r *http.Request) {
	var params types.ServerListParams
	err := qstring.Unmarshal(r.URL.Query(), &params)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, errors.Wrap(err, "invalid parameters"))
		return
	}

	servers, err := v.Storage.GetServers(params.Page, params.Sort, params.By, params.Filters)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, errors.Wrap(err, "failed to get servers"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(servers)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, errors.Wrap(err, "failed to encode response"))
		return
	}
}
