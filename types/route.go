package types

import (
	"net/http"
	"net/url"
)

// Route represents an API route and its associated handler function
type Route struct {
	Name        string           `json:"name"`
	Method      string           `json:"method"`
	Path        string           `json:"path"`
	Description string           `json:"description"`
	Params      url.Values       `json:"params"`
	Accepts     interface{}      `json:"accepts"`
	Returns     interface{}      `json:"returns"`
	Handler     http.HandlerFunc `json:"-"`
}

// RouteHandler represents an version group of API endpoints
type RouteHandler interface {
	Version() string
	Routes() []Route
}

// // EndpointHandler wraps a HTTP handler function with app-specific args/returns
// type EndpointHandler func(w http.ResponseWriter, r *http.Request)

// // ServeHTTP implements the necessary chaining functionality for HTTP middleware
// func (f EndpointHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
// 	f(w, r)
// 	// TODO: errors
// 	// if err != nil {
// 	// 	logger.Error("request handler failed",
// 	// 		zap.Error(err))
// 	// 	w.WriteHeader(http.StatusInternalServerError)
// 	// }

// 	w.Header().Set("Content-Type", "application/json")
// }
