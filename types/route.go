package types

import (
	"net/http"
)

// Route represents an API route and its associated handler function
type Route struct {
	Name    string           `json:"name"`
	Method  string           `json:"method"`
	Path    string           `json:"path"`
	Accepts interface{}      `json:"accepts"`
	Returns interface{}      `json:"returns"`
	Handler http.HandlerFunc `json:"-"`
}

// RouteHandler represents an version group of API endpoints
type RouteHandler interface {
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
