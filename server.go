package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

// Server handles either posting a server object or requesting a server object
func (app *App) Server(w http.ResponseWriter, r *http.Request) {
	address, ok := mux.Vars(r)["address"]
	if !ok {
		logger.Fatal("no address specified in request",
			zap.String("request", r.URL.String()))
	}

	switch r.Method {
	case "GET":
		logger.Debug("getting server",
			zap.String("address", address))

	case "POST":
		logger.Debug("posting server",
			zap.String("address", address))
	}
}
