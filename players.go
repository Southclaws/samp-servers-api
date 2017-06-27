package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

// Players returns a JSON encoded array of players for a particular server
func (app *App) Players(w http.ResponseWriter, r *http.Request) {
	address, ok := mux.Vars(r)["address"]
	if !ok {
		logger.Fatal("no address specified in request",
			zap.String("request", r.URL.String()))
	}

	logger.Debug("getting player list for server",
		zap.String("address", address))
}
