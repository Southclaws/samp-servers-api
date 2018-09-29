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
