// Package v3 implements version 3 of the public API
package v3

import (
	"net/http"

	"github.com/Southclaws/samp-servers-api/scraper"
	"github.com/Southclaws/samp-servers-api/storage"
	"github.com/Southclaws/samp-servers-api/types"
)

// V3 represents an API endpoint handler
type V3 struct {
	Storage *storage.Manager
	Scraper *scraper.Scraper
	Metrics *types.Metrics
	Config  types.Config
}

// Init initialises and returns a handler group
func Init(Storage *storage.Manager, Scraper *scraper.Scraper, Metrics *types.Metrics, Config types.Config) *V3 {
	return &V3{
		Storage: Storage,
		Scraper: Scraper,
		Metrics: Metrics,
		Config:  Config,
	}
}

// Version returns the route group version name
func (v *V3) Version() string { return "v3" }

// Routes returns the version routes
func (v *V3) Routes() []types.Route {
	return []types.Route{
		{
			Name:        "serverAdd",
			Path:        "/server/{address}",
			Method:      "POST",
			Description: `Add a server to the index using just the IP address. This endpoint requires no body and no additional information. The IP address is added to an internal queue and will be queried periodically for information via the legacy server API. This allows any server to be added with the basic information provided by SA:MP itself.`,
			Accepts:     nil,
			Returns:     nil,
			Handler:     v.serverAdd,
		},
		{
			Name:        "serverPost",
			Path:        "/server",
			Method:      "POST",
			Description: `Provide additional information for a server such as a description and a banner image. This requires a body to be posted which contains information for the server.`,
			Accepts:     types.Server{}.Example(),
			Returns:     nil,
			Handler:     v.serverPost,
		},
		{
			Name:        "serverGet",
			Path:        "/server/{address}",
			Method:      "GET",
			Description: `Returns a full server object using the specified address.`,
			Accepts:     nil,
			Returns:     types.Server{}.Example(),
			Handler:     v.serverGet,
		},
		{
			Name:        "serverList",
			Path:        "/servers",
			Method:      "GET",
			Description: "Returns a list of servers based on the specified query parameters. Supported query parameters are: `page` `sort` `by` `filters`.",
			Params:      types.ServerListParams{}.Example(),
			Accepts:     nil,
			Returns:     []types.ServerCore{types.Server{}.Example().Core, types.Server{}.Example().Core, types.Server{}.Example().Core},
			Handler:     v.serverList,
		},
		{
			Name:        "serverStats",
			Path:        "/stats",
			Method:      "GET",
			Description: `Returns a some statistics of the server index.`,
			Accepts:     nil,
			Returns:     types.Statistics{}.Example(),
			Handler:     v.serverStats,
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
