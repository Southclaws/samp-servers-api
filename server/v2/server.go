package v2

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"

	"github.com/Southclaws/samp-servers-api/types"
)

// serverAdd handles "simple" posts where the only data is the server address which is passed to
// the QueryDaemon which handles pulling the rest of the information from the legacy query API.
func (v *V2) serverAdd(w http.ResponseWriter, r *http.Request) {
	address := r.FormValue("address")
	if address == "" {
		WriteError(w, http.StatusBadRequest, errors.New("no address specified"))
		return
	}

	normalised, errs := types.AddressFromString(address)
	if errs != nil {
		WriteErrors(w, http.StatusBadRequest, errs)
		return
	}

	v.Scraper.Add(normalised)

	w.Header().Set("Location", "/")
	w.WriteHeader(http.StatusFound)
}

// serverPost handles posting a server object
func (v *V2) serverPost(w http.ResponseWriter, r *http.Request) {
	var from string
	if from = r.Header.Get("X-Forwarded-For"); from == "" {
		from = strings.Split(r.RemoteAddr, ":")[0]
	}

	server := types.Server{}
	err := json.NewDecoder(r.Body).Decode(&server)
	if err != nil {
		WriteError(w, http.StatusBadRequest, err)
		return
	}

	if v.Config.VerifyByHost {
		addressIP := strings.Split(server.Core.Address, ":")[0]
		if from != addressIP {
			WriteError(w, http.StatusBadRequest,
				errors.Errorf("request address '%v' does not match declared server address '%s'", from, addressIP))
			return
		}
	}

	errs := server.Validate()
	if errs != nil {
		WriteErrors(w, http.StatusUnprocessableEntity, errs)
		return
	}

	server.Active = true

	err = v.Storage.UpsertServer(server)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err)
	}

	v.Scraper.Add(server.Core.Address)
}

// serverGet handles responding to a request by server address
func (v *V2) serverGet(w http.ResponseWriter, r *http.Request) {
	address, ok := mux.Vars(r)["address"]
	if !ok {
		WriteError(w, http.StatusBadRequest, errors.New("no address specified"))
	}

	var err error

	_, errs := types.AddressFromString(address)
	if errs != nil {
		WriteErrors(w, http.StatusBadRequest, errs)
		return
	}

	server, found, err := v.Storage.GetServer(address)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err)
		return
	}

	if !found {
		WriteError(w, http.StatusNotFound, errors.Errorf("could not find server by address '%s'", address))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(&server)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err)
		return
	}
}
