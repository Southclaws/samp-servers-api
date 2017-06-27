package main

import (
	"net/http"
)

// Servers returns a JSON encoded array of available servers
func (app *App) Servers(w http.ResponseWriter, r *http.Request) {
	logger.Debug("getting server list")
}
