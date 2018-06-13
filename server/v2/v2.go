// Package v2 implements version 2 of the public API
package v2

import (
	"net/http"

	"github.com/Southclaws/samp-servers-api/scraper"
	"github.com/Southclaws/samp-servers-api/storage"
	"github.com/Southclaws/samp-servers-api/types"
)

// V2 represents an API endpoint handler
type V2 struct {
	Storage *storage.Manager
	Scraper *scraper.Scraper
	Config  types.Config
}

// Init initialises and returns a handler group
func Init(Storage *storage.Manager, Scraper *scraper.Scraper, Config types.Config) *V2 {
	return &V2{
		Storage: Storage,
		Scraper: Scraper,
		Config:  Config,
	}
}

// Routes returns the version routes
func (v *V2) Routes() []types.Route {
	return []types.Route{
		{
			Name:    "serverAdd",
			Path:    "/server",
			Method:  "POST",
			Handler: v.serverAdd,
		},
		{
			Name:    "serverPost",
			Path:    "/server/{address}",
			Method:  "POST",
			Handler: v.serverPost,
		},
		{
			Name:    "serverGet",
			Path:    "/server/{address}",
			Method:  "GET",
			Handler: v.serverGet,
		},
		{
			Name:    "serverList",
			Path:    "/servers",
			Method:  "GET",
			Handler: v.serverList,
		},
		{
			Name:    "serverStats",
			Path:    "/stats",
			Method:  "GET",
			Handler: v.serverStats,
		},
	}
}

// TODO: replace with handler wrapper

// WriteError is a utility function for logging a request error and writing a response all in one.
func WriteError(w http.ResponseWriter, status int, err error) {
	w.WriteHeader(status)
	w.Write([]byte(err.Error()))
}

// WriteErrors does the same but for groups of errors
func WriteErrors(w http.ResponseWriter, status int, errs []error) {
	w.WriteHeader(status)
	for _, err := range errs {
		w.Write([]byte(err.Error() + ", "))
	}
}
